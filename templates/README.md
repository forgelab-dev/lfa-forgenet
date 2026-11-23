# Templates du Marketplace

Chaque dossier de `templates/` est une application que ForgeNet installe en
un clic depuis son dashboard. Ce document décrit le format d'un template ; les
raisons de chaque règle sont dans l'[ADR-0018](../docs/decisions/0018-format-des-templates-marketplace.md).

Le plus simple pour commencer : copier [`_exemple/`](_exemple/), qui passe la
vérification, puis l'adapter.

## Le dossier

```
templates/postgres/
├── manifest.yaml        ce que le dashboard affiche et demande
├── docker-compose.yml   la stack déployée
├── icon.svg             l'icône de la carte du catalogue
└── README.md            le service, ses variables, comment s'y connecter
```

- Le **nom du dossier** identifie le template : minuscules, chiffres et `-`,
  31 caractères au plus (`postgres`, `uptime-kuma`).
- Ces quatre fichiers, et aucun autre : le fichier Compose ne peut pas lire
  ses voisins.
- Un dossier qui commence par `_` est vérifié mais n'apparaît pas au
  catalogue.

## `manifest.yaml`

```yaml
schema: 1
name: PostgreSQL
description: Base de données relationnelle, la plus répandue des bases libres.
category: database
version: "17.2"
service: postgres
license: PostgreSQL
homepage: https://www.postgresql.org

variables:
  - name: POSTGRES_USER
    label: Utilisateur créé à l'installation
    default: app
    pattern: "^[a-z_][a-z0-9_]{0,62}$"
  - name: POSTGRES_DB
    label: Base créée à l'installation
    default: app
    pattern: "^[a-z_][a-z0-9_]{0,62}$"
  - name: POSTGRES_PASSWORD
    label: Mot de passe de l'utilisateur
    secret: true
    generate: 32

ports:
  - service: postgres
    port: 5432
    protocol: tcp
    description: Connexions des autres applications

volumes:
  - name: data
    description: Fichiers de la base
```

| Champ | Obligatoire | Règle |
|---|---|---|
| `schema` | oui | `1` |
| `name` | oui | Le nom affiché, 40 caractères au plus. |
| `description` | oui | Une phrase, 160 caractères au plus. |
| `category` | oui | `database`, `cache`, `storage`, `cms`, `devtools`, `monitoring`, `communication` ou `other`. |
| `version` | oui | La version de l'application, entre guillemets. Elle doit figurer dans l'étiquette de l'image du service principal : `17.2` dans `postgres:17.2-alpine`. |
| `service` | oui | Le service principal du fichier Compose. |
| `license` | oui | L'identifiant [SPDX](https://spdx.org/licenses/) de la licence du logiciel : `MIT`, `AGPL-3.0-only`… |
| `homepage` | oui | Le site du logiciel, en `https://`. |
| `variables` | non | Les valeurs demandées à l'installation. |
| `ports` | non | Les ports où les services écoutent. |
| `volumes` | non | Les volumes nommés, chacun avec une `description`. |

Un champ inconnu est refusé : une faute de frappe ne passe pas.

### Les variables

| Champ | Règle |
|---|---|
| `name` | Majuscules, chiffres et `_`. Interdits : `FORGENET_…`, `COMPOSE_…`, `DOCKER_…`, `PATH`, `HOME`. |
| `label` | Le libellé du champ du formulaire d'installation. |
| `help` | Facultatif : une phrase d'aide. |
| `default` | Facultatif : la valeur proposée. Jamais pour un secret. |
| `pattern` | Facultatif : expression régulière que la valeur respecte. |
| `optional` | Facultatif : `true` accepte une valeur vide. |
| `secret` | Facultatif : `true` masque la valeur, la chiffre sur le serveur et ne la réaffiche jamais. |
| `generate` | Pour un secret : longueur (16 à 128) d'une valeur tirée au hasard si l'administrateur ne saisit rien. Lettres et chiffres seulement, utilisables sans échappement dans une URL. |

Le fichier Compose utilise la variable sous la forme `${NOM}`, et seulement
celle-là. La valeur par défaut va dans le manifeste, pas dans un
`${NOM:-défaut}`. Un `$` littéral s'écrit `$$` (`$${NOM}` est lu par le shell
du conteneur, pas par Compose).

### Les ports

| Champ | Règle |
|---|---|
| `service`, `port` | Le service et le port où il écoute **dans son conteneur**. |
| `protocol` | `http` : l'administrateur peut le relier à un domaine, servi en HTTPS. `tcp` : interne, pour d'autres applications du serveur. |
| `description` | À quoi sert ce port. |

Un template qui déclare un port `tcp` déclare aussi au moins une variable
`secret: true` : un service joignable par d'autres applications demande
toujours une authentification.

## `docker-compose.yml`

Ce qui est vérifié, en plus d'un fichier que `docker compose` accepte :

- **Image épinglée** pour chaque service : `postgres:17.2-alpine`, ou une
  empreinte `@sha256:…`. Jamais `latest`, ni une image sans étiquette. Pas de
  `build:`.
- **Aucun secret en dur** : une variable d'environnement dont le nom contient
  `PASS`, `SECRET`, `TOKEN`, `KEY` ou `PRIVATE` vaut `${NOM}`, où `NOM` est une
  variable `secret: true` du manifeste.
- **Variables cohérentes** : chaque `${NOM}` est déclaré dans le manifeste, et
  chaque variable du manifeste sert dans le fichier. Une entrée
  d'environnement a toujours une valeur.
- **Volumes nommés** et déclarés dans le manifeste. Aucun montage de l'hôte
  (`/chemin:/…`, `./chemin:/…`), aucun volume anonyme, aucun volume ou réseau
  `external`.
- **Aucun port publié** (`ports:`) : il s'ouvrirait sur Internet en
  contournant le pare-feu. Déclarez le port dans le manifeste.
- **Rien qui sorte du conteneur** : pas de `privileged`, `cap_add`, `devices`,
  `network_mode: host`, `pid: host`, `ipc: host`, ni de `security_opt` qui
  désactive une protection.
- **Pas de `container_name`** : deux installations du même template se
  disputeraient le nom.
- **Une sonde de santé** (`healthcheck`) sur le service principal :
  l'installation attend qu'il soit sain.

`restart` est facultatif : ForgeNet donne `unless-stopped` aux services qui
n'en ont pas.

## `icon.svg`

Un SVG carré (`viewBox="0 0 64 64"`), de 16 Kio au plus, sans `<script>`, sans
attribut `on…` et sans lien hors du fichier.

## `README.md`

Court : ce que fait le service, le tableau de ses variables, comment s'y
connecter depuis une autre application, et ce que gardent ses volumes.

## Vérifier son template

La CI vérifie chaque dossier à chaque pull request. Pour le faire avant, avec
Go et Docker installés, depuis la racine du dépôt :

```bash
go test ./templates/
```

Chaque erreur nomme le template, le fichier et la règle, par exemple :

```
templates/postgres/docker-compose.yml : services.postgres : image "postgres:latest" : l'étiquette "latest" change sous le même nom ; épinglez une version
```

La vérification ne remplace pas l'essai réel. Tant que la page Marketplace du
dashboard n'existe pas (#116), essayez la stack avec `docker compose`, en
donnant les variables comme ForgeNet le fait (dans l'environnement de la
commande, ADR-0019) :

```bash
cd templates/postgres
POSTGRES_USER=app POSTGRES_DB=app POSTGRES_PASSWORD=essai-local \
  docker compose -p essai-postgres up -d --wait
docker compose -p essai-postgres ps        # le service est « healthy »
docker compose -p essai-postgres restart   # les données restent
docker compose -p essai-postgres down -v   # retire la stack et ses volumes
```

Un ForgeNet lancé sur votre machine peut aussi lire votre dossier à la place
du catalogue embarqué : démarrez `forgenet-server` avec
`FORGENET_MARKETPLACE_DIR="$PWD/templates"` (un chemin absolu). Votre template
apparaît alors au catalogue et s'installe comme ceux de la version publiée.
Les dossiers en `_` restent écartés, et le serveur refuse de démarrer si un
manifeste ou une icône est invalide.

Joignez une capture de l'essai à la pull request.
