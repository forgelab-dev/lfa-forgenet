#!/usr/bin/env bash
#
# Tests des réglages SMTP d'install.sh (ADR-0014), sans root ni Docker : les
# fonctions sont chargées (FORGENET_INSTALL_LIB=1) et écrivent dans un dossier
# temporaire. Lancement : bash deploy/tests/install_test.sh

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

printf '\n%d vérifications, %d en échec\n' "$((pass + fail))" "$fail"
[ "$fail" -eq 0 ]
