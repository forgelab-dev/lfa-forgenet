#!/usr/bin/env bash
#
# Installation de ForgeNet sur un serveur Debian ou Ubuntu (ADR-0009).
#
#   curl -fsSLO https://raw.githubusercontent.com/forgelab-dev/lfa-forgenet/v0.2.0/deploy/install.sh # x-release-please-version
#   less install.sh            # relisez avant de lancer en root
#   sudo bash install.sh [--domain forgenet.exemple.org --email vous@exemple.org]
#
# Le script peut être relancé sans risque : il garde la clé de session et le
# compte administrateur déjà créés, et applique seulement ce qui change
# (version, domaine). Options : bash install.sh --help
#
# Ce qu'il fait :
#   1. vérifie le système (Debian/Ubuntu, x86_64 ou ARM64, root, ports 80/443) ;
#   2. installe Docker depuis le dépôt officiel s'il manque (clé vérifiée) ;
#   3. écrit /opt/forgenet : docker-compose.yml, Caddyfile.base et .env (0600),
#      avec une clé de session aléatoire et le groupe du socket Docker ;
#   4. démarre le serveur et Caddy, puis crée le premier administrateur.
#
# Il ne touche ni au pare-feu, ni à SSH, ni aux autres conteneurs.

set -euo pipefail
umask 077

readonly REPO_RAW="https://raw.githubusercontent.com/forgelab-dev/lfa-forgenet"
# Empreinte publiée de la clé de signature des paquets Docker.
readonly DOCKER_KEY_FPR="9DC858229FC7DD38854AE2D88D81803C0EBFCD88"

INSTALL_DIR="${FORGENET_DIR:-/opt/forgenet}"
VERSION="${FORGENET_VERSION:-0.2.0}" # x-release-please-version
DOMAIN="${FORGENET_DOMAIN-}"
ACME_EMAIL="${FORGENET_ACME_EMAIL-}"
ADMIN_USER="${FORGENET_ADMIN_USER-}"
ADMIN_EMAIL="${FORGENET_ADMIN_EMAIL-}"
SOURCE_DIR=""
DOMAIN_SET=0
UPDATE=0
VERSION_SET=0
EMAIL_SET=0
[ -n "${FORGENET_DOMAIN+x}" ] && DOMAIN_SET=1
[ -n "${FORGENET_ACME_EMAIL+x}" ] && EMAIL_SET=1

info() { printf '\033[1;33m==>\033[0m %s\n' "$*"; }
ok() { printf '\033[1;32m ok\033[0m %s\n' "$*"; }
warn() { printf '\033[1;35m /!\\\033[0m %s\n' "$*" >&2; }
die() {
	printf '\033[1;31merreur :\033[0m %s\n' "$*" >&2
	exit 1
}

usage() {
	cat <<'EOF'
Installation de ForgeNet (Debian, Ubuntu).

Usage : sudo bash install.sh [options]

  --domain NOM        Domaine du dashboard (enregistrement DNS vers ce serveur).
                      Sans domaine, le dashboard s'ouvre par un tunnel SSH.
                      --domain "" retire un domaine déjà configuré.
  --email ADRESSE     Adresse des avis Let's Encrypt (facultative).
  --version X.Y.Z     Version de ForgeNet (défaut : celle publiée avec ce script).
  --update            Met à jour une installation existante : sauvegarde, bascule,
                      retour arrière automatique en cas d'échec. Sans --version,
                      vise la dernière version publiée.
  --admin-user NOM    Nom du premier administrateur (sinon demandé).
  --admin-email ADR   Adresse du premier administrateur (sinon demandée).
  --dir CHEMIN        Dossier d'installation (défaut : /opt/forgenet).
  --source DOSSIER    Prendre docker-compose.yml et Caddyfile.base dans ce
                      dossier au lieu de les télécharger (tests, hors ligne).
  -h, --help          Cette aide.

Les mêmes réglages sont lus dans FORGENET_DOMAIN, FORGENET_ACME_EMAIL,
FORGENET_VERSION, FORGENET_ADMIN_USER, FORGENET_ADMIN_EMAIL et FORGENET_DIR.
EOF
}

need_value() { [ $# -ge 2 ] || die "l'option $1 attend une valeur (voir --help)"; }

while [ $# -gt 0 ]; do
	case "$1" in
	--domain) need_value "$@"; DOMAIN="$2"; DOMAIN_SET=1; shift 2 ;;
	--email) need_value "$@"; ACME_EMAIL="$2"; EMAIL_SET=1; shift 2 ;;
	--version) need_value "$@"; VERSION="$2"; VERSION_SET=1; shift 2 ;;
	--admin-user) need_value "$@"; ADMIN_USER="$2"; shift 2 ;;
	--admin-email) need_value "$@"; ADMIN_EMAIL="$2"; shift 2 ;;
	--dir) need_value "$@"; INSTALL_DIR="$2"; shift 2 ;;
	--source) need_value "$@"; SOURCE_DIR="$2"; shift 2 ;;
	--update) UPDATE=1; shift ;;
	-h | --help) usage; exit 0 ;;
	*) die "option inconnue : $1 (voir --help)" ;;
	esac
done

# --- Validation des entrées -------------------------------------------------
# Tout ce qui finit dans .env est vérifié : une valeur avec un saut de ligne
# ou un caractère spécial pourrait y injecter une autre variable.

DOMAIN="$(printf '%s' "$DOMAIN" | tr '[:upper:]' '[:lower:]')"
DOMAIN="${DOMAIN%.}"
if [ -n "$DOMAIN" ] && ! [[ "$DOMAIN" =~ ^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$ ]]; then
	die "domaine invalide : « $DOMAIN » (exemple : forgenet.exemple.org)"
fi
email_re='^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
if [ -n "$ACME_EMAIL" ] && ! [[ "$ACME_EMAIL" =~ $email_re ]]; then
	die "adresse Let's Encrypt invalide : « $ACME_EMAIL »"
fi
[[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]] || die "version invalide : « $VERSION » (exemple : 0.1.0)"
[[ "$INSTALL_DIR" = /* ]] || die "--dir attend un chemin absolu"

# --- 1. Système ---------------------------------------------------------------

# os_field lit un champ de /etc/os-release sans rien définir dans le script.
os_field() {
	(
		# shellcheck disable=SC1091
		. /etc/os-release
		eval "printf '%s' \"\${$1:-}\""
	)
}

check_system() {
	info "Vérification du système"
	[ "$(id -u)" -eq 0 ] || die "lancez le script en root : sudo bash install.sh"
	[ -r /etc/os-release ] || die "système non reconnu (/etc/os-release absent)"
	# os-release est lu dans un sous-shell : il définit VERSION, NAME, ID…
	# qui écraseraient les variables du script.
	local os_id os_like os_name os_codename ubuntu_codename
	os_id="$(os_field ID)"
	os_like="$(os_field ID_LIKE)"
	os_name="$(os_field PRETTY_NAME)"
	os_codename="$(os_field VERSION_CODENAME)"
	ubuntu_codename="$(os_field UBUNTU_CODENAME)"
	case "$os_id" in
	debian | ubuntu) DISTRO="$os_id" ;;
	*)
		case " $os_like " in
		*" ubuntu "*) DISTRO=ubuntu ;;
		*" debian "*) DISTRO=debian ;;
		*) die "seuls Debian et Ubuntu sont pris en charge (système : ${os_name:-inconnu})" ;;
		esac
		;;
	esac
	# Une dérivée (Mint, Pop!_OS…) utilise les paquets Docker de sa base.
	if [ "$DISTRO" = ubuntu ]; then
		CODENAME="${ubuntu_codename:-$os_codename}"
	else
		CODENAME="$os_codename"
	fi
	[ -n "$CODENAME" ] || die "nom de version du système introuvable dans /etc/os-release"

	case "$(uname -m)" in
	x86_64 | aarch64 | arm64) ;;
	*) die "architecture non prise en charge : $(uname -m) (x86_64 ou ARM64)" ;;
	esac

	local mem_kb
	mem_kb="$(awk '/^MemTotal:/ {print $2}' /proc/meminfo)"
	if [ "${mem_kb:-0}" -lt 1800000 ]; then
		warn "moins de 2 Go de mémoire : ForgeNet fonctionnera, mais à l'étroit"
	fi
	ok "${os_name:-$DISTRO} ($(uname -m))"
}

check_ports() {
	# Sur une relance, les ports sont tenus par notre propre Caddy.
	if command -v docker >/dev/null 2>&1 && [ "$(docker inspect -f '{{.State.Running}}' forgenet-caddy 2>/dev/null || true)" = true ]; then
		return
	fi
	command -v ss >/dev/null 2>&1 || return 0
	local busy
	busy="$(ss -ltnH '( sport = :80 or sport = :443 )' 2>/dev/null || true)"
	if [ -n "$busy" ]; then
		die "les ports 80 ou 443 sont déjà utilisés (un autre serveur web ?) ; Caddy en a besoin :
$busy"
	fi
}

# --- 2. Docker ----------------------------------------------------------------

install_docker() {
	info "Installation de Docker (dépôt officiel download.docker.com)"
	export DEBIAN_FRONTEND=noninteractive
	# needrestart (Ubuntu) ne doit ni questionner ni redémarrer de services ici.
	export NEEDRESTART_SUSPEND=1
	apt-get update -qq
	apt-get install -y -qq ca-certificates curl gnupg >/dev/null

	install -m 0755 -d /etc/apt/keyrings
	local key
	key="$(mktemp)"
	curl -fsSL "https://download.docker.com/linux/${DISTRO}/gpg" -o "$key"
	# On n'ajoute la clé qu'après avoir vérifié son empreinte.
	if ! gpg --show-keys --with-colons "$key" 2>/dev/null | grep -q "^fpr:::::::::${DOCKER_KEY_FPR}:"; then
		rm -f "$key"
		die "la clé de signature téléchargée n'a pas l'empreinte attendue ; installation de Docker arrêtée"
	fi
	install -m 0644 "$key" /etc/apt/keyrings/docker.asc
	rm -f "$key"

	cat >/etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/${DISTRO}
Suites: ${CODENAME}
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF
	chmod 0644 /etc/apt/sources.list.d/docker.sources

	apt-get update -qq
	apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin >/dev/null
	systemctl enable --now docker >/dev/null 2>&1 || true
}

ensure_docker() {
	if ! command -v docker >/dev/null 2>&1; then
		install_docker
	fi
	docker info >/dev/null 2>&1 || die "Docker est installé mais ne répond pas (systemctl status docker)"

	local server major
	server="$(docker version --format '{{.Server.Version}}')"
	major="${server%%.*}"
	[ "${major:-0}" -ge 24 ] 2>/dev/null || die "Docker $server trop ancien : ForgeNet demande Docker 24 ou plus"
	docker compose version >/dev/null 2>&1 || die "le plugin Docker Compose manque (paquet docker-compose-plugin)"
	ok "Docker $server, $(docker compose version --short 2>/dev/null || echo 'Compose')"
}

# --- 3. Fichiers ----------------------------------------------------------------

fetch_file() {
	local name="$1" dest="$2" tmp
	tmp="$(mktemp)"
	if [ -n "$SOURCE_DIR" ]; then
		cp "$SOURCE_DIR/$name" "$tmp" || die "$SOURCE_DIR/$name introuvable"
	elif ! curl -fsSL "$REPO_RAW/v${VERSION}/deploy/$name" -o "$tmp"; then
		rm -f "$tmp"
		die "téléchargement de $name impossible pour la version $VERSION (version publiée ?)"
	fi
	install -m 0644 "$tmp" "$dest"
	rm -f "$tmp"
}

# set_env écrit CLÉ=valeur dans .env, en remplaçant la ligne existante.
set_env() {
	local key="$1" value="$2" file="$INSTALL_DIR/.env" tmp
	tmp="$(mktemp "$INSTALL_DIR/.env.XXXXXX")"
	grep -v "^${key}=" "$file" >"$tmp" || true
	printf '%s=%s\n' "$key" "$value" >>"$tmp"
	chmod 0600 "$tmp"
	mv "$tmp" "$file"
}

get_env() { sed -n "s/^$1=//p" "$INSTALL_DIR/.env" | tail -n 1; }

write_files() {
	info "Fichiers d'installation dans $INSTALL_DIR"
	install -d -m 0750 "$INSTALL_DIR"
	fetch_file docker-compose.yml "$INSTALL_DIR/docker-compose.yml"
	fetch_file Caddyfile.base "$INSTALL_DIR/Caddyfile.base"

	local env="$INSTALL_DIR/.env"
	if [ ! -f "$env" ]; then
		: >"$env"
		chmod 0600 "$env"
	fi
	chown root:root "$env"
	chmod 0600 "$env"

	# La clé de session n'est créée qu'une fois : la changer déconnecterait
	# tout le monde à chaque relance.
	if [ -z "$(get_env FORGENET_JWT_SECRET)" ]; then
		set_env FORGENET_JWT_SECRET "$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')"
		ok "clé de session générée"
	fi
	set_env FORGENET_VERSION "$VERSION"
	set_env FORGENET_IMAGE_DIGEST "$(release_digest)"
	set_env FORGENET_INSTALL_DIR "$INSTALL_DIR"
	set_env DOCKER_GID "$(stat -c %g /var/run/docker.sock)"

	# Domaine et adresse : ceux de la ligne de commande, sinon ceux déjà en place.
	if [ "$DOMAIN_SET" -eq 1 ]; then
		set_env FORGENET_DOMAIN "$DOMAIN"
	else
		DOMAIN="$(get_env FORGENET_DOMAIN)"
		[ -n "$(get_env FORGENET_DOMAIN)" ] || set_env FORGENET_DOMAIN ""
	fi
	if [ "$EMAIL_SET" -eq 1 ]; then
		set_env FORGENET_ACME_EMAIL "$ACME_EMAIL"
	else
		[ -n "$(get_env FORGENET_ACME_EMAIL)" ] || set_env FORGENET_ACME_EMAIL ""
	fi
	ok "docker-compose.yml, Caddyfile.base, .env (0600)"
}

# release_digest lit l'empreinte de l'image dans la fiche de la release
# (ADR-0011) : l'installation tire alors l'image exacte signée par la CI, pas
# seulement une étiquette. Sans fiche (v0.1.0), l'étiquette seule est gardée.
release_digest() {
	local url="https://github.com/forgelab-dev/lfa-forgenet/releases/download/v${VERSION}/forgenet-release.json"
	local fiche digest
	if ! fiche="$(curl -fsSL --max-time 30 "$url" 2>/dev/null)"; then
		warn "fiche de la release v${VERSION} introuvable : image désignée par son étiquette seule"
		return 0
	fi
	digest="$(printf '%s\n' "$fiche" | sed -n 's/.*"digest": *"\(sha256:[0-9a-f]\{64\}\)".*/\1/p' | head -n 1)"
	if [ -z "$digest" ] || ! printf '%s' "$fiche" | grep -q '"image": *"ghcr.io/forgelab-dev/forgenet-server"' ||
		! printf '%s' "$fiche" | grep -q "\"version\": *\"${VERSION}\""; then
		die "fiche de la release v${VERSION} invalide : installation arrêtée"
	fi
	printf '%s' "$digest"
}

# --- 4. Démarrage et administrateur ---------------------------------------------

compose() { docker compose --project-directory "$INSTALL_DIR" -f "$INSTALL_DIR/docker-compose.yml" "$@"; }

start_stack() {
	info "Démarrage de ForgeNet $VERSION"
	if ! compose up -d --pull missing --remove-orphans; then
		die "démarrage impossible. Si l'image ghcr.io/forgelab-dev/forgenet-server:$VERSION est introuvable, cette version n'est pas publiée."
	fi

	local status=""
	for _ in $(seq 1 60); do
		status="$(docker inspect -f '{{.State.Health.Status}}' forgenet-server 2>/dev/null || true)"
		[ "$status" = healthy ] && break
		sleep 2
	done
	if [ "$status" != healthy ]; then
		compose logs --tail 30 server >&2 || true
		die "le serveur ne répond pas après 2 minutes (état : ${status:-inconnu}) ; journaux ci-dessus"
	fi
	ok "serveur et Caddy démarrés"
}

# ask lit une réponse sur le terminal, même quand le script arrive par un tube.
ask() {
	local prompt="$1" reply
	printf '%s' "$prompt" >/dev/tty
	IFS= read -r reply </dev/tty
	printf '%s' "$reply"
}

create_admin() {
	local marker="$INSTALL_DIR/.admin-created"
	if [ -f "$marker" ]; then
		ok "administrateur déjà créé ($(cat "$marker"))"
		return
	fi
	info "Création du premier administrateur"

	local interactive=0
	if [ -r /dev/tty ] && [ -w /dev/tty ] && (: </dev/tty) 2>/dev/null; then
		interactive=1
	fi
	if [ -z "$ADMIN_USER" ]; then
		[ "$interactive" -eq 1 ] || die "pas de terminal : passez --admin-user et --admin-email"
		ADMIN_USER="$(ask "Nom d'utilisateur [admin] : ")"
		ADMIN_USER="${ADMIN_USER:-admin}"
	fi
	if [ -z "$ADMIN_EMAIL" ]; then
		[ "$interactive" -eq 1 ] || die "pas de terminal : passez --admin-email"
		ADMIN_EMAIL="$(ask "Adresse email : ")"
	fi
	[[ "$ADMIN_USER" =~ ^[A-Za-z0-9._-]{3,32}$ ]] || die "nom d'utilisateur invalide : 3 à 32 caractères parmi lettres, chiffres, . _ -"
	[[ "$ADMIN_EMAIL" =~ $email_re ]] || die "adresse email invalide : « $ADMIN_EMAIL »"

	local password confirm password_file=""
	if [ "$interactive" -eq 1 ]; then
		while :; do
			printf 'Mot de passe (8 caractères au moins) : ' >/dev/tty
			IFS= read -rs password </dev/tty
			printf '\nConfirmation : ' >/dev/tty
			IFS= read -rs confirm </dev/tty
			printf '\n' >/dev/tty
			if [ "${#password}" -lt 8 ]; then
				warn "mot de passe trop court"
			elif [ "$password" != "$confirm" ]; then
				warn "les deux saisies diffèrent"
			else
				break
			fi
		done
	else
		# Sans terminal, le mot de passe va dans un fichier root, jamais à l'écran
		# ni dans un journal.
		password="$(head -c 18 /dev/urandom | base64 | tr -d '/+=\n')"
		password_file="$INSTALL_DIR/admin-password"
		printf '%s\n' "$password" >"$password_file"
		chmod 0600 "$password_file"
	fi

	if ! printf '%s\n' "$password" | compose exec -T server forgenet-server user create \
		--username "$ADMIN_USER" --email "$ADMIN_EMAIL" --role admin; then
		die "création de l'administrateur impossible (message ci-dessus)"
	fi
	unset password confirm
	printf '%s\n' "$ADMIN_USER" >"$marker"
	chmod 0600 "$marker"
	if [ -n "$password_file" ]; then
		warn "mot de passe généré dans $password_file : notez-le, puis supprimez ce fichier"
	fi
}

check_dns() {
	[ -n "$DOMAIN" ] || return 0
	if ! getent ahosts "$DOMAIN" >/dev/null 2>&1; then
		warn "$DOMAIN ne se résout pas encore : créez l'enregistrement DNS (A ou AAAA) vers ce serveur ; Caddy obtiendra le certificat dès qu'il pointera ici"
	fi
}

summary() {
	echo
	info "ForgeNet $VERSION est installé"
	if [ -n "$DOMAIN" ]; then
		echo "    Dashboard : https://$DOMAIN"
		echo "    Le certificat HTTPS est obtenu automatiquement dès que le DNS pointe ici."
	else
		echo "    Pas de domaine : le dashboard n'écoute que sur ce serveur (127.0.0.1:8080)."
		echo "    Depuis votre ordinateur, ouvrez un tunnel SSH :"
		echo "        ssh -N -L 8080:127.0.0.1:8080 ${SUDO_USER:-utilisateur}@<adresse-du-serveur>"
		echo "    puis ouvrez http://localhost:8080"
		echo "    Pour un domaine plus tard : sudo bash install.sh --domain forgenet.exemple.org"
	fi
	echo "    Fichiers : $INSTALL_DIR (gardez .env secret : il contient la clé de session)"
	echo "    Pare-feu : seuls 22 (SSH), 80 et 443 ont besoin d'être ouverts."
}

# Première version dont le serveur sait se mettre à jour lui-même (ADR-0011).
readonly SELF_UPDATE_SINCE="0.2.0"

# version_ge A B : vrai si A >= B (versions X.Y.Z, suffixe de pré-version ignoré).
version_ge() {
	[ "$(printf '%s\n%s\n' "${1%%-*}" "${2%%-*}" | sort -V | head -n 1)" = "${2%%-*}" ]
}

# run_update confie la mise à jour au serveur en place quand il sait la faire :
# sauvegarde, bascule, santé, retour arrière automatique. Sinon, ce script
# réinstalle simplement la version qu'il porte, comme avant.
run_update() {
	[ -f "$INSTALL_DIR/.env" ] || die "aucune installation dans $INSTALL_DIR : lancez install.sh sans --update"
	local installed
	installed="$(docker exec forgenet-server forgenet-server -version 2>/dev/null | awk '{print $2}')"
	if [ -z "$installed" ] || ! version_ge "$installed" "$SELF_UPDATE_SINCE"; then
		warn "la version installée (${installed:-inconnue}) ne sait pas encore se mettre à jour : réinstallation vers $VERSION par ce script"
		return 1
	fi
	local to="latest"
	[ "$VERSION_SET" -eq 1 ] && to="$VERSION"
	info "Mise à jour de ForgeNet $installed (vers : $to)"
	local started
	started="$(docker exec forgenet-server forgenet-server update start --to "$to")" || die "mise à jour refusée (message ci-dessus)"
	printf '%s\n' "$started"
	# Déjà à jour : rien n'a été lancé.
	case "$started" in *"Mise à jour lancée"*) ;; *) exit 0 ;; esac
	sleep 2
	# Le conteneur reste après la fin : ses journaux se lisent jusqu'au bout.
	docker logs -f forgenet-updater 2>&1 || true
	local state reason
	state="$(docker exec forgenet-server cat /var/lib/forgenet/maj/etat.json 2>/dev/null || true)"
	reason="$(printf '%s\n' "$state" | sed -n 's/^ *"raison": "\(.*\)",*$/\1/p')"
	case "$state" in
	*'"resultat": "reussie"'*) ok "mise à jour terminée : $(docker exec forgenet-server forgenet-server -version)" ;;
	*'"resultat": "annulee"'*) die "mise à jour annulée, retour à $installed : ${reason:-raison inconnue}" ;;
	*) die "mise à jour en échec : ${reason:-état illisible}" ;;
	esac
	exit 0
}

main() {
	if [ "$UPDATE" -eq 1 ]; then
		run_update || true
	fi
	check_system
	ensure_docker
	check_ports
	write_files
	start_stack
	create_admin
	check_dns
	summary
}

main
