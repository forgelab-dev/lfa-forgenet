# ADR-0018 : format des templates du Marketplace

- Statut : accepté (revue de la PR #120 par le porteur du projet, 2026-11-21)
- Date : 2026-11-21
- Contexte : ticket #111, découpé de #12 ; prépare les templates #112 à #115
  et l'installation en un clic (#116) ; s'appuie sur les ADR-0003 (stacks
  Compose), 0008 (domaines), 0012 (relance après redémarrage) et 0015
  (secrets chiffrés en base)

## Contexte

`CONTRIBUTING.md` dit qu'un template est un dossier de `templates/` avec un
`manifest.yaml` et un `docker-compose.yml`, mais rien ne dit ce que contient
le manifeste. Quatre templates (PostgreSQL, MySQL, Redis, MinIO) sont
proposés à des contributeurs qui découvrent le projet : sans format écrit et
vérifié, chacun inventerait le sien, et l'installation en un clic (#116)
devrait deviner.

Trois questions se posent avant le format lui-même :

1. **Comment les valeurs saisies à l'installation arrivent dans le fichier
   Compose.** L'ADR-0003 coupe la commande `docker compose` de
   l'environnement du serveur : un `${VAR}` du fichier ne vaut rien tant que
   ForgeNet ne le fournit pas.
2. **Où le serveur trouve le catalogue.** Le dépôt, ou une copie livrée
   avec ForgeNet.
3. **Ce qu'un template a le droit de faire.** Un fichier Compose donne un
   accès complet à l'hôte (ADR-0003) ; un template, écrit par n'importe qui
   et installé d'un clic, doit en faire beaucoup moins.

## Décision

### Le dossier d'un template

```
templates/postgres/
├── manifest.yaml        ce que le dashboard affiche et demande
├── docker-compose.yml   la stack, déployée comme toute stack Compose
├── icon.svg             l'icône de la carte du catalogue
└── README.md            le service, ses variables, comment s'y connecter
```

- Le **nom du dossier** est l'identifiant du template : minuscules, chiffres
  et `-`, 31 caractères au plus, commençant par une lettre ou un chiffre.
- Un dossier commençant par `_` n'entre pas au catalogue mais passe la même
  vérification : c'est la place du template d'exemple (`templates/_exemple/`).
- L'icône est un fichier fixe plutôt qu'un champ du manifeste : un seul
  endroit, que la vérification contrôle (SVG carré, 16 Kio au plus, ni
  `<script>`, ni attribut `on…`, ni lien externe).

### Le manifeste

Exemple complet, celui que le template PostgreSQL (#112) pourra reprendre :

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
| `schema` | oui | `1`. Un changement incompatible du format passera à `2` ; un champ facultatif ajouté ne change pas le numéro. |
| `name` | oui | Le nom affiché, 40 caractères au plus. |
| `description` | oui | Une phrase pour la carte du catalogue, 160 caractères au plus. |
| `category` | oui | Une des catégories de la liste fermée : `database`, `cache`, `storage`, `cms`, `devtools`, `monitoring`, `communication`, `other`. Une catégorie s'ajoute par une PR qui modifie la liste et la vérification. |
| `version` | oui | La version de l'application, entre guillemets. Elle doit figurer dans l'étiquette de l'image du service principal (`17.2` dans `postgres:17.2-alpine`) : la version affichée est celle qui tourne. |
| `service` | oui | Le service principal de la stack : son état et sa version représentent l'application dans le dashboard. |
| `license` | oui | Identifiant SPDX de la licence du logiciel (`MIT`, `AGPL-3.0-only`, `PostgreSQL`…). |
| `homepage` | oui | Le site du logiciel, en `https://`. |
| `variables` | non | Les valeurs demandées à l'installation, détaillées plus bas. |
| `ports` | non | Les ports où les services écoutent, détaillés plus bas. |
| `volumes` | non | Les volumes nommés de la stack, chacun avec une `description` de ce qu'il garde. |

Un champ inconnu est une erreur : une faute de frappe (`generat: 32`) ne
passe pas en silence.

Le texte est en français, comme le dashboard. Les traductions viendront avec
celles du dashboard, par des champs facultatifs.

### Les variables

| Champ | Règle |
|---|---|
| `name` | Majuscules, chiffres et `_`, 64 caractères au plus. Refusés : les préfixes `FORGENET_`, `COMPOSE_` et `DOCKER_`, et `PATH` ou `HOME`, parce que ces variables changeraient le comportement de la commande `docker compose` elle-même (projet, démon visé). |
| `label` | Le libellé du champ du formulaire d'installation. |
| `help` | Facultatif : une phrase d'aide sous le champ. |
| `default` | Facultatif : la valeur proposée. |
| `pattern` | Facultatif : expression régulière que la valeur doit respecter, vérifiée par le serveur et par le formulaire. |
| `optional` | Facultatif, `false` par défaut : sans lui, une valeur vide est refusée. |
| `secret` | Facultatif, `false` par défaut : la valeur est masquée à la saisie, chiffrée en base, jamais réaffichée ni écrite au journal. |
| `generate` | Réservé aux secrets : longueur (16 à 128) d'une valeur tirée au hasard si l'administrateur laisse le champ vide. L'alphabet est `A-Z a-z 0-9` : la valeur passe dans une URL ou une chaîne de connexion sans échappement, et 32 caractères donnent environ 190 bits. Un autre alphabet s'ajoutera le jour où un template en aura besoin. |

**Le chemin des valeurs.** Le fichier Compose garde ses `${POSTGRES_PASSWORD}`.
À l'installation, ForgeNet enregistre les valeurs avec le déploiement (les
secrets chiffrés par la clé des secrets du serveur, comme le mot de passe du
relais, ADR-0015) et les ajoute à l'environnement de **chaque** commande
Compose de cette stack : déploiement, relance après un redémarrage
(ADR-0012), suppression. Seules les variables déclarées par le manifeste
passent ; l'environnement du serveur reste coupé (ADR-0003).

Deux autres chemins sont écartés :

- **Remplacer les `${VAR}` dans le contenu du fichier** à l'installation : le
  mot de passe serait alors écrit en clair dans la base, dans les sauvegardes
  et dans le fichier affiché par le dashboard.
- **Écrire un fichier `.env`** dans le dossier de travail : un secret en
  clair sur le disque, alors que l'ADR-0003 n'y écrit rien.

### Les ports

| Champ | Règle |
|---|---|
| `service`, `port` | Le service et le port où il écoute **dans son conteneur**. |
| `protocol` | `http` : le formulaire d'installation peut lui donner un domaine, servi par Caddy en HTTPS (ADR-0008). `tcp` : pas de domaine possible, le service sert d'autres applications du serveur. |
| `description` | À quoi sert ce port. |

Un template ne **publie aucun port sur l'hôte** (`ports:` interdit dans son
`docker-compose.yml`). Docker ouvre un port publié en contournant le
pare-feu du serveur : une base installée d'un clic serait joignable depuis
Internet. Une application web passe par un domaine, une base reste interne.

### Ce que le fichier Compose d'un template a le droit de faire

En plus des règles de toute stack (ADR-0003 : images publiées, volumes
nommés, aucun fichier voisin) :

- **Images épinglées** : chaque service a une `image` avec une étiquette qui
  contient un chiffre, ou une empreinte `@sha256:`. Refusés : pas
  d'étiquette, `latest`, `stable`, `edge`, `main`, `master`, `nightly`.
  Pas de `build:`.
- **Aucun secret en dur** : une variable d'environnement dont le nom évoque
  un secret (`PASSWORD`, `PASSWD`, `SECRET`, `TOKEN`, `KEY`, `PRIVATE`) a pour
  valeur une référence `${NOM}` à une variable du manifeste marquée
  `secret: true`.
- **Variables cohérentes** : chaque `${NOM}` du fichier est déclaré dans le
  manifeste, et chaque variable du manifeste sert dans le fichier. Seule la
  forme `${NOM}` est admise : la valeur par défaut vit dans le manifeste,
  pas dans un `${NOM:-défaut}` qui en ferait deux. Un `$` littéral s'écrit
  `$$`, comme dans tout fichier Compose.
- **Volumes déclarés** : les volumes nommés du fichier et ceux du manifeste
  sont les mêmes. Aucun montage de l'hôte (chemin absolu, relatif ou `~`),
  donc ni le socket Docker ni les fichiers du serveur.
- **Rien qui sorte du conteneur** : refusés `privileged`, `cap_add`,
  `devices`, `network_mode`, `pid`, `ipc` et `userns_mode` à `host`,
  `security_opt` qui désactive une protection.
- **Pas de `container_name`** : deux installations du même template se
  disputeraient le nom.
- **Une sonde de santé sur le service principal** (`healthcheck`) :
  `up --wait` attend qu'il soit sain (ADR-0003), et l'installation ne se dit
  réussie qu'à ce moment.
- `restart` est facultatif : ForgeNet donne `unless-stopped` aux services qui
  n'en précisent pas (ADR-0012).

### La vérification

- Le format est un type Go lu en mode strict ; c'est lui qui fait foi, et
  `templates/README.md` le décrit pour les contributeurs. Un JSON Schema
  pour l'autocomplétion des éditeurs pourra suivre ; il ne remplacerait pas
  la vérification, qui croise le manifeste et le fichier Compose.
- Le fichier Compose est lu sous sa forme résolue par
  `docker compose config --no-interpolate --format json`, comme pour la
  relance (ADR-0012) : Compose y a déjà appliqué ancres, `extends` et
  syntaxes courtes, et les `${NOM}` restent visibles. La vérification ne
  réimplémente pas la lecture d'un fichier Compose.
- Elle tourne dans `go test` sur chaque dossier de `templates/`, donc à
  chaque PR qui touche `templates/` (ajouté au filtre du job Go). Sans
  Docker, elle est sautée sur un poste de développeur, mais **échoue** dans
  la CI plutôt que d'être sautée.
- Ses messages nomment le template, le fichier et la règle : un contributeur
  qui découvre Go doit pouvoir corriger sans lire le code.

### Le catalogue

- Les templates sont **embarqués dans le binaire** (`go:embed`) : le
  catalogue d'un serveur est celui de sa version de ForgeNet, vérifié par la
  CI et disponible sans réseau.
- Écarté pour l'instant : télécharger le catalogue depuis GitHub. Le
  dashboard dépendrait du réseau, et lirait à chaud un contenu que rien ne
  signe, qui changerait sous une version en service.
- À l'installation, le fichier Compose du template est **copié dans le
  déploiement**. Une nouvelle version du template ne touche pas les
  applications déjà installées ; les mettre à jour sera une fonction à part.

## Points tranchés à la revue (PR #120, 2026-11-21)

### Joindre un service `tcp` depuis une autre application

Chaque stack a son propre réseau : une application ne voit pas la base d'une
autre stack. Décision : **un réseau géré par ForgeNet, et sécurisé**.

- **Un réseau par installation** qui déclare un port `tcp`, nommé
  `forgenet-svc-<projet>`, plutôt qu'un réseau commun à toutes les bases :
  une application reliée à une base ne voit pas les autres.
- Seuls les **services qui portent un port `tcp`** y sont branchés, pas le
  reste de la stack, avec l'alias `<service>.<projet>.forgenet.internal`
  (même forme que l'ADR-0008). Le README de chaque template dit cette
  adresse.
- Le réseau est **interne** (`internal: true` de Docker) : il ne mène nulle
  part hors du serveur. Chaque conteneur garde le réseau de sa propre stack
  pour le reste.
- **Caddy n'y est jamais branché**, et ce réseau n'est pas `forgenet-edge` :
  une base n'est jamais à portée du proxy.
- Comme pour l'ADR-0008, le réseau porte une étiquette
  (`forgenet.managed=service`), et un réseau du même nom sans elle est
  refusé.
- **Une application n'y entre que par un geste de l'administrateur**
  (« relier à une base »), inscrit au journal d'audit, et réversible.
  Rien ne s'y branche seul.
- **Pas de service `tcp` sans authentification** : un template qui déclare
  un port `tcp` déclare au moins une variable `secret: true`. La vérification
  le contrôle. Le template Redis (#114) devra donc exiger un mot de passe.

Chaque base reliable consomme un réseau Docker. La réserve d'adresses par
défaut de Docker en permet une trentaine, réseaux des stacks compris : le
dashboard dira clairement l'erreur si elle s'épuise, et `install.sh` pourra
l'élargir si le besoin apparaît.

### Changer une variable après l'installation

Faire tourner un mot de passe, par exemple : **plus tard**, validé tel quel.
Le premier jet de #116 installe ; la modification viendra ensuite.

## Conséquences

- Les contributeurs des templates #112 à #115 ont un format écrit, un exemple
  qui passe, et une vérification qui leur dit quoi corriger.
- Le dashboard (#116) lit le catalogue sans le deviner : carte, formulaire
  et domaines proposés viennent du manifeste.
- Un nouveau template arrive sur les serveurs avec la version suivante de
  ForgeNet, que la mise à jour en un clic installe (ADR-0011).
- Les règles de sécurité des templates sont plus strictes que celles d'une
  stack déployée à la main : un administrateur garde la liberté de l'ADR-0003
  pour ses propres fichiers, le catalogue ne propose que des stacks bornées.
- Le moteur Compose doit accepter des variables propres à un déploiement, et
  le déploiement les garder (secrets chiffrés) : un changement du contrat
  d'API, fait avec #116.
- La CI dépend de Docker pour vérifier les templates : c'est déjà le cas des
  tests d'intégration du moteur Compose.
