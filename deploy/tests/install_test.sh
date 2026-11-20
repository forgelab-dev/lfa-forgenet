#!/usr/bin/env bash
#
# Tests d'install.sh pour l'ADR-0014 (réglages SMTP, premier administrateur
# sans terminal), sans root ni Docker : les fonctions sont chargées
# (FORGENET_INSTALL_LIB=1) et écrivent dans un dossier temporaire ; compose est
# remplacé par une doublure. Lancement : bash deploy/tests/install_test.sh

set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT="$HERE/../install.sh"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

pass=0
fail=0

# run lance write_smtp dans un sous-shell, avec les options données, sur le
# .env de $WORK/install ; sa sortie (erreurs comprises) va dans $WORK/out.
# setsid retire le terminal de contrôle : /dev/tty ne s'ouvre plus, comme
# sous cloud-init, et rien n'est demandé sur le terminal du développeur.
run() {
	# Les apostrophes sont voulues : « $0 » et « $@ » sont ceux du bash lancé.
	# shellcheck disable=SC2016
	FORGENET_INSTALL_LIB=1 FORGENET_DIR="$WORK/install" \
		setsid -w bash -c '. "$0" "$@" && write_smtp' "$SCRIPT" "$@" </dev/null >"$WORK/out" 2>&1
}

value() { sed -n "s/^$1=//p" "$WORK/install/.env" | tail -n 1; }

check() {
	local name="$1" got="$2" want="$3"
	if [ "$got" = "$want" ]; then
		pass=$((pass + 1))
		printf 'ok    %s\n' "$name"
	else
		fail=$((fail + 1))
		printf 'ÉCHEC %s\n      attendu : %s\n      obtenu  : %s\n' "$name" "$want" "$got"
	fi
}

mkdir -p "$WORK/install"
: >"$WORK/install/.env"

# Le mot de passe vient de la variable, jamais d'un argument, et « $ » y
# reste tel quel grâce aux apostrophes.
FORGENET_SMTP_PASSWORD='pa$$ word' run --smtp-host SMTP.Exemple.org --smtp-user forgenet@exemple.org
check "relais écrit" "$(value FORGENET_SMTP_HOST)" "smtp.exemple.org"
check "port par défaut" "$(value FORGENET_SMTP_PORT)" "587"
check "STARTTLS sur 587" "$(value FORGENET_SMTP_TLS)" "starttls"
check "expéditeur pris de l'identifiant" "$(value FORGENET_SMTP_FROM)" "forgenet@exemple.org"
check "mot de passe entre apostrophes" "$(value FORGENET_SMTP_PASSWORD)" "'pa\$\$ word'"
check "mot de passe jamais affiché" "$(grep -c 'pa\$\$ word' "$WORK/out")" "0"
check ".env en 0600" "$(stat -c %a "$WORK/install/.env")" "600"

# Sans --smtp-host, les réglages en place sont gardés.
run
check "réglages gardés à la relance" "$(value FORGENET_SMTP_HOST)" "smtp.exemple.org"

FORGENET_SMTP_PASSWORD=secret run --smtp-host smtp.exemple.org --smtp-port 465 --smtp-user forgenet --smtp-from envois@exemple.org
check "TLS direct sur 465" "$(value FORGENET_SMTP_TLS)" "implicit"
check "expéditeur donné" "$(value FORGENET_SMTP_FROM)" "envois@exemple.org"

FORGENET_SMTP_PASSWORD="avec'apostrophe" run --smtp-host smtp.exemple.org --smtp-user forgenet --smtp-from envois@exemple.org
check "apostrophe refusée" "$(grep -c "apostrophe" "$WORK/out")" "1"
check "réglages précédents intacts après un refus" "$(value FORGENET_SMTP_PASSWORD)" "'secret'"

run --smtp-host smtp.exemple.org --smtp-user forgenet --smtp-from envois@exemple.org
check "pas de terminal ni de variable : refus" "$(grep -c "FORGENET_SMTP_PASSWORD" "$WORK/out")" "1"

run --smtp-host ""
check "envoi coupé" "$(value FORGENET_SMTP_HOST)$(value FORGENET_SMTP_USER)$(value FORGENET_SMTP_PASSWORD)" ""

run --smtp-host smtp.exemple.org --smtp-from envois@exemple.org
check "relais sans identifiant" "$(value FORGENET_SMTP_HOST):$(value FORGENET_SMTP_USER)" "smtp.exemple.org:"

# --- Premier administrateur sans terminal -----------------------------------

# admin lance create_admin sans terminal, avec une doublure de compose qui
# note ses appels dans $WORK/appels et répond comme le serveur : login-link
# affiche un faux lien, --send réussit sauf si ENVOI_ECHOUE=1.
admin() {
	rm -f "$WORK/install/.admin-created" "$WORK/install/premier-acces" "$WORK/appels"
	# shellcheck disable=SC2016
	FORGENET_INSTALL_LIB=1 FORGENET_DIR="$WORK/install" APPELS="$WORK/appels" \
		FORGENET_ADMIN_USER=chef FORGENET_ADMIN_EMAIL=chef@exemple.org \
		setsid -w bash -c '. "$0" && compose() {
			printf "%s\n" "$*" >>"$APPELS"
			case "$*" in
			*"login-link"*"--send"*) [ "${ENVOI_ECHOUE:-0}" != 1 ] || { echo "relais injoignable" >&2; return 1; } ;;
			*"login-link"*) echo "  http://127.0.0.1:8080/acces#JETON-SECRET" ;;
			esac
		} && create_admin' "$SCRIPT" </dev/null >"$WORK/out" 2>&1
}

: >"$WORK/install/.env"
admin
check "compte créé sans mot de passe" "$(grep -c -- 'user create --no-password --username chef' "$WORK/appels")" "1"
check "lien écrit dans premier-acces" "$(grep -c 'acces#JETON-SECRET' "$WORK/install/premier-acces")" "1"
check "premier-acces en 0600" "$(stat -c %a "$WORK/install/premier-acces")" "600"
check "lien jamais affiché" "$(grep -c 'JETON-SECRET' "$WORK/out")" "0"
check "plus de admin-password" "$([ -e "$WORK/install/admin-password" ] && echo présent || echo absent)" "absent"

printf 'FORGENET_SMTP_HOST=smtp.exemple.org\n' >"$WORK/install/.env"
admin
check "avec relais : lien envoyé" "$(grep -c -- 'login-link --username chef --send' "$WORK/appels")" "1"
check "avec relais : pas de fichier" "$([ -e "$WORK/install/premier-acces" ] && echo présent || echo absent)" "absent"

ENVOI_ECHOUE=1 admin
check "envoi refusé : lien écrit dans le fichier" "$(grep -c 'acces#JETON-SECRET' "$WORK/install/premier-acces")" "1"
check "envoi refusé : raison donnée" "$(grep -c 'relais injoignable' "$WORK/out")" "1"
check "envoi refusé : lien jamais affiché" "$(grep -c 'JETON-SECRET' "$WORK/out")" "0"

# --- Réglages du dashboard face à --smtp-host (ADR-0015) --------------------

# smtp_choice lance apply_smtp_choice avec une doublure de compose : « mail
# reset » répond EFFACE (réglages du dashboard effacés) ou non, et chaque appel
# est noté ; wait_healthy est remplacé, faute de Docker.
smtp_choice() {
	rm -f "$WORK/appels"
	# shellcheck disable=SC2016
	FORGENET_INSTALL_LIB=1 FORGENET_DIR="$WORK/install" APPELS="$WORK/appels" \
		setsid -w bash -c '. "$0" "$@" && compose() {
			printf "%s\n" "$*" >>"$APPELS"
			case "$*" in
			*"mail reset"*) if [ "${EFFACE:-0}" = 1 ]; then echo "Réglages d'"'"'envoi du dashboard effacés"; else echo "Aucun réglage"; fi ;;
			esac
		} && wait_healthy() { echo "attente" >>"$APPELS"; } && apply_smtp_choice' "$SCRIPT" "$@" </dev/null >"$WORK/out" 2>&1
}

FORGENET_SMTP_PASSWORD=secret EFFACE=1 smtp_choice --smtp-host smtp.exemple.org --smtp-user forgenet --smtp-from envois@exemple.org
check "--smtp-host efface les réglages du dashboard" "$(grep -c 'exec -T server forgenet-server mail reset' "$WORK/appels")" "1"
check "puis le serveur redémarre pour relire .env" "$(grep -c '^restart server$' "$WORK/appels")" "1"

FORGENET_SMTP_PASSWORD=secret smtp_choice --smtp-host smtp.exemple.org --smtp-user forgenet --smtp-from envois@exemple.org
check "rien à effacer : pas de redémarrage" "$(grep -c '^restart server$' "$WORK/appels")" "0"

smtp_choice
check "sans --smtp-host, les réglages du dashboard restent" "$([ -e "$WORK/appels" ] && cat "$WORK/appels")" ""

printf '\n%d vérifications, %d en échec\n' "$((pass + fail))" "$fail"
[ "$fail" -eq 0 ]
