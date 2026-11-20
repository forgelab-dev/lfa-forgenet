#!/usr/bin/env bash
#
# Tests de la commande deploy/forgenet (#71), sans Docker ni root : docker, id
# et la copie d'install.sh sont remplacés par des doublures qui écrivent ce
# qu'on leur demande. Lancement : bash deploy/tests/forgenet_test.sh

set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CMD="$HERE/../forgenet"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

mkdir -p "$WORK/bin" "$WORK/install/sauvegardes"

# Doublure de docker : son comportement vient des variables FAKE_*.
cat >"$WORK/bin/docker" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$*" >>"$FAKE_CALLS"
case "$1" in
info) exit "${FAKE_DOCKER_DOWN:-0}" ;;
inspect)
	name="${!#}"
	case "$name" in
	forgenet-server) v="${FAKE_SERVER-running healthy}" ;;
	forgenet-caddy) v="${FAKE_CADDY-running }" ;;
	forgenet-updater) v="${FAKE_UPDATER-}" ;;
	*) v="" ;;
	esac
	[ -n "$v" ] || exit 1
	printf '%s\n' "$v"
	;;
exec)
	case "$*" in
	*" -version") printf 'forgenet-server %s (566bd31, 2026-11-16T12:24:33Z, go1.26.11)\n' "${FAKE_VERSION:-0.4.1}" ;;
	*"etat.json") [ -n "${FAKE_STATE-}" ] && printf '%s\n' "$FAKE_STATE" || exit 1 ;;
	*) printf 'exec %s\n' "$*" ;;
	esac
	;;
compose)
	# Sortie JSON de « docker compose ls --all --format json » ; FAKE_APPS
	# donne « nom état » par ligne.
	[ "$*" = "compose ls --all --format json" ] || exit 1
	printf '['
	sep=""
	while read -r name state; do
		[ -n "$name" ] || continue
		printf '%s{"Name":"%s","Status":"%s","ConfigFiles":"/x.yml"}' "$sep" "$name" "$state"
		sep=","
	done <<<"${FAKE_APPS-}"
	printf ']\n'
	;;
logs) printf 'logs %s\n' "$*" ;;
esac
EOF
cat >"$WORK/bin/id" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "${FAKE_UID:-0}"
EOF
chmod +x "$WORK/bin/docker" "$WORK/bin/id"

# Copie d'install.sh de l'installation : elle écrit ses arguments.
cat >"$WORK/install/install.sh" <<'EOF'
printf 'install.sh %s\n' "$*"
EOF

export PATH="$WORK/bin:$PATH"
export FORGENET_DIR="$WORK/install"
export FAKE_CALLS="$WORK/calls"

fails=0
runs=0

# check NOM CODE_ATTENDU MOTIF -- ARGS… : lance forgenet ARGS, vérifie le code
# de sortie et la présence de MOTIF (expression grep -E) dans la sortie.
check() {
	local name="$1" want="$2" pattern="$3"
	shift 4
	: >"$FAKE_CALLS"
	local out code
	out="$(bash "$CMD" "$@" 2>&1)"
	code=$?
	runs=$((runs + 1))
	if [ "$code" -ne "$want" ] || ! printf '%s\n' "$out" | grep -qE -- "$pattern"; then
		fails=$((fails + 1))
		printf 'ÉCHEC %s : code %s (attendu %s)\n--- sortie ---\n%s\n--------------\n' "$name" "$code" "$want" "$out"
	else
		printf 'ok    %s\n' "$name"
	fi
}

# called MOTIF : vrai si un appel à docker du dernier test correspond.
called() {
	runs=$((runs + 1))
	if grep -qE -- "$1" "$FAKE_CALLS"; then
		printf 'ok    appel docker : %s\n' "$1"
	else
		fails=$((fails + 1))
		printf 'ÉCHEC appel docker attendu : %s\nappels :\n%s\n' "$1" "$(cat "$FAKE_CALLS")"
	fi
}

check "aide sans argument" 0 "Usage : forgenet COMMANDE" --
check "help" 0 "restore DOSSIER" -- help
check "commande inconnue" 1 "commande inconnue : danse" -- danse
FORGENET_DIR=relatif check "FORGENET_DIR relatif" 1 "FORGENET_DIR attend un chemin absolu" -- status

check "version" 0 "ForgeNet 0.4.1 \(révision 566bd31, compilé le 2026-11-16T12:24:33Z\)" -- version
check "version de la commande" 0 "Commande forgenet [0-9]+\.[0-9]+\.[0-9]+" -- version
FAKE_VERSION=9.9.9 check "commande d'une autre version" 0 "sudo forgenet update" -- version
FAKE_UID=1000 check "user set-email, transmis au serveur" 0 "exec exec -i forgenet-server forgenet-server user set-email --username ana --email ana@exemple.org" -- user set-email --username ana --email ana@exemple.org
FAKE_SERVER="exited " check "version, serveur arrêté" 1 "serveur ForgeNet est arrêté" -- version
FAKE_SERVER="" check "version, ForgeNet absent" 1 "n'existe pas" -- version
FAKE_DOCKER_DOWN=1 check "Docker injoignable" 1 "Docker ne répond pas" -- version

FAKE_STATE='{
  "de": "0.4.0",
  "vers": "0.4.1",
  "resultat": "reussie",
}' FAKE_APPS=$'forgenet running(2)\nforgenet-boutique running(1)\nforgenet-blog exited(2)\nautre running(1)\n' \
	check "status, tout va bien" 0 "Serveur : en marche, sain" -- status
FAKE_STATE='{
  "de": "0.4.0",
  "vers": "0.4.1",
  "resultat": "reussie",
}' check "status, dernière mise à jour" 0 "Dernière mise à jour : 0\.4\.0 → 0\.4\.1, réussie" -- status
FAKE_APPS=$'forgenet running(2)\nforgenet-boutique running(1)\nforgenet-blog exited(2)\nautre running(1)\n' \
	check "status, applications" 0 "blog +2 arrêté" -- status
FAKE_APPS=$'forgenet-boutique running(1)\nautre running(1)\n' \
	check "status, piles étrangères ignorées" 0 "^  boutique +1 en marche$" -- status
if FAKE_APPS=$'autre running(1)\n' bash "$CMD" status 2>&1 | grep -q autre; then
	fails=$((fails + 1))
	printf 'ÉCHEC status : une pile hors ForgeNet est listée\n'
fi
FAKE_APPS=$'forgenet-agenda exited(1), running(1)\n' \
	check "status, pile à moitié arrêtée" 0 "agenda +1 arrêté\\(s\\), 1 en marche" -- status
FAKE_APPS=$'forgenet running(2)\n' check "status, aucune application" 0 "Applications : aucune" -- status
FAKE_SERVER="running unhealthy" check "status, serveur malade" 1 "ne répond pas" -- status
FAKE_UID=1000 check "user set-email, transmis au serveur" 0 "exec exec -i forgenet-server forgenet-server user set-email --username ana --email ana@exemple.org" -- user set-email --username ana --email ana@exemple.org
FAKE_SERVER="exited " check "status, serveur arrêté" 1 "Serveur : exited" -- status
FAKE_CADDY="" check "status, Caddy absent" 1 "Caddy : conteneur forgenet-caddy absent" -- status
check "status sans option" 1 "ne prend pas d'option" -- status --tout

check "logs" 0 "logs logs --tail 100 forgenet-server" -- logs
check "logs -f -n 20" 0 "logs logs --tail 20 -f forgenet-server" -- logs -f -n 20
check "logs caddy" 0 "logs logs --tail 100 forgenet-caddy" -- logs caddy
FAKE_UPDATER="exited " check "logs maj" 0 "logs logs --tail 100 forgenet-updater" -- logs maj
check "logs maj sans mise à jour" 1 "aucune mise à jour lancée" -- logs maj
check "logs -n invalide" 1 "-n attend un nombre" -- logs -n beaucoup
check "logs option inconnue" 1 "option inconnue" -- logs --tout

check "update check" 0 "exec exec forgenet-server forgenet-server update check" -- update check
check "update check --pre" 0 "update check --pre" -- update check --pre
check "update" 0 "install.sh --update --dir $WORK/install" -- update
check "update --version" 0 "install.sh --update --dir $WORK/install --version 0.5.0" -- update --version 0.5.0
FAKE_UID=1000 check "update sans root" 1 "demande root : sudo forgenet update" -- update
check "update option inconnue" 1 "option inconnue" -- update --vite

check "backups vides" 0 "Aucune sauvegarde" -- backups
for b in 20261114-150834-0.1.1 20261116-121202-0.4.0; do
	mkdir -p "$WORK/install/sauvegardes/$b"
	: >"$WORK/install/sauvegardes/$b/forgenet.db"
done
printf 'FORGENET_VERSION=0.4.0\n' >"$WORK/install/sauvegardes/20261116-121202-0.4.0/.env"
mkdir -p "$WORK/install/sauvegardes/pas-une-sauvegarde"
check "backups, plus récente d'abord" 0 "^20261116-121202-0\.4\.0 +2026-11-16 12:12:02 +0\.4\.0" -- backups
check "backups, version inconnue" 0 "20261114-150834-0\.1\.1 +2026-11-14 15:08:34 +\?" -- backups
check "backups, dossier étranger ignoré" 0 "^SAUVEGARDE" -- backups
if bash "$CMD" backups 2>&1 | grep -q pas-une-sauvegarde; then
	fails=$((fails + 1))
	printf 'ÉCHEC backups : un dossier sans base est listé\n'
fi
FAKE_UID=1000 check "backups sans root" 1 "demande root" -- backups

check "restore" 0 "install.sh --dir $WORK/install --restore 20261116-121202-0.4.0" -- restore 20261116-121202-0.4.0
check "restore sans dossier" 1 "usage : sudo forgenet restore" -- restore
check "restore avec option" 1 "usage : sudo forgenet restore" -- restore --tout
FAKE_UID=1000 check "restore sans root" 1 "demande root" -- restore 20261116-121202-0.4.0
mv "$WORK/install/install.sh" "$WORK/install.sh.de-cote"
check "restore sans copie d'install.sh" 1 "relancez install.sh une fois" -- restore 20261116-121202-0.4.0
mv "$WORK/install.sh.de-cote" "$WORK/install/install.sh"

FAKE_UID=1000 check "user, membre du groupe docker" 0 "exec exec -i forgenet-server forgenet-server user create --username ana" -- user create --username ana
check "user sans sous-commande" 1 "usage : forgenet user create" -- user
FAKE_UID=1000 check "user set-password, transmis au serveur" 0 "exec exec -i forgenet-server forgenet-server user set-password --username ana" -- user set-password --username ana
FAKE_UID=1000 check "user set-email, transmis au serveur" 0 "exec exec -i forgenet-server forgenet-server user set-email --username ana --email ana@exemple.org" -- user set-email --username ana --email ana@exemple.org
FAKE_SERVER="exited " check "user, serveur arrêté" 1 "serveur ForgeNet est arrêté" -- user create

printf '\n%d vérifications, %d en échec\n' "$runs" "$fails"
[ "$fails" -eq 0 ]
