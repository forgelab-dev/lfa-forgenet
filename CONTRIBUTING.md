# Contribuer à ForgeNet

Merci de l'intérêt que vous portez au projet. ForgeNet est développé en Afrique de l'Ouest avec des moyens modestes — chaque contribution compte réellement.

> **English speakers welcome.** Issues and pull requests in English are perfectly fine. The project documentation is in French today; translating it is itself a valuable contribution.

> ⚠️ **Le projet est en pré-alpha.** L'architecture bouge encore. Avant d'investir du temps sur une fonctionnalité conséquente, ouvrez une issue pour en discuter — cela évite le scénario pénible d'une PR de 800 lignes qui part dans une direction déjà écartée.

## Sommaire

- [Manières de contribuer](#manières-de-contribuer)
- [Code de conduite](#code-de-conduite)
- [Mettre en place l'environnement](#mettre-en-place-lenvironnement)
- [Se repérer dans le monorepo](#se-repérer-dans-le-monorepo)
- [Le flux de travail](#le-flux-de-travail)
- [Conventions de commit](#conventions-de-commit)
- [Style de code et tests](#style-de-code-et-tests)
- [La règle du contrat `api/`](#la-règle-du-contrat-api)
- [Proposer un template Marketplace](#proposer-un-template-marketplace)
- [Signaler une faille de sécurité](#signaler-une-faille-de-sécurité)
- [Licence de vos contributions](#licence-de-vos-contributions)

## Manières de contribuer

Le code n'est pas la seule porte d'entrée :

| | |
|---|---|
| 🐛 **Signaler un bug** | Ouvrez une [issue](https://github.com/forgelab-dev/lfa-forgenet/issues) avec les étapes de reproduction, la version, l'OS et les logs pertinents |
| 💡 **Proposer une fonctionnalité** | Ouvrez une issue en décrivant **le problème** avant la solution — le « pourquoi » oriente mieux la discussion que le « comment » |
| 📦 **Ajouter un template Marketplace** | La contribution la plus accessible et l'une des plus utiles — voir [plus bas](#proposer-un-template-marketplace) |
| 🌍 **Traduire** | Documentation et interface, en priorité vers l'anglais |
| 📖 **Améliorer la documentation** | Une section confuse corrigée vaut souvent mieux qu'une fonctionnalité de plus |
| 💬 **Raconter votre usage** | Retours terrain, cas d'usage, contraintes réseau : ce sont eux qui orientent la roadmap |

Vous cherchez un point de départ ? Les issues étiquetées [`good first issue`](https://github.com/forgelab-dev/lfa-forgenet/labels/good%20first%20issue) sont choisies pour être abordables sans connaître tout le projet.

## Code de conduite

Soyez respectueux et constructif. Nous partons du principe que chacun agit de bonne foi et débute quelque part. Sont exclus : le harcèlement, les attaques personnelles, les propos discriminatoires, et le mépris envers ceux qui posent des questions de débutant.

Un comportement problématique peut être signalé en privé à **contact@lucidforgeafrica.com**. Les mainteneurs peuvent modérer, éditer ou rejeter toute contribution, et bloquer les comptes en cas de manquement répété.

## Mettre en place l'environnement

**Prérequis**

| Outil | Version | Nécessaire pour |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.23+ | Core Engine, agent, CLI |
| [Node.js](https://nodejs.org) | 20+ | Dashboard web |
| [Flutter](https://flutter.dev) | 3.x | Application mobile |
| [Docker](https://docs.docker.com/engine/install/) | 24+ | Exécution et tests d'orchestration |
| [Task](https://taskfile.dev) | 3+ | Lancement des commandes du projet |
| [buf](https://buf.build/docs/installation) | 1.x | Génération du code depuis les `.proto` |

Vous n'avez besoin que de la chaîne correspondant à votre contribution : une PR sur le dashboard ne requiert ni Flutter ni Go.

```bash
git clone https://github.com/forgelab-dev/lfa-forgenet.git
cd lfa-forgenet

task setup     # installe les dépendances des trois chaînes
task dev       # démarre la stack complète en local
task test      # exécute la suite de tests
task lint      # vérifie le style sur les trois langages
task gen       # régénère les clients depuis api/proto
```

## Se repérer dans le monorepo

```
api/         Contrats gRPC + OpenAPI — la source de vérité
cmd/         Binaires Go : forgenet-server, forgenet-agent, forgenet-cli
internal/    Cœur métier Go — c'est ici que vit la logique
pkg/         SDK Go public, utilisable par des projets tiers
web/         Dashboard React + Vite + Tailwind
mobile/      Application Flutter
templates/   Catalogue Marketplace 1-click
deploy/      install.sh, docker-compose, Caddyfile
docs/        Spécification, matrice fonctionnelle, ADRs
```

Deux règles d'architecture à connaître avant de déplacer du code :

- **`internal/transport/` ne contient aucune logique métier.** Les handlers gRPC, HTTP et WebSocket valident, appellent le domaine, et sérialisent la réponse. Rien d'autre.
- **`internal/` est privé, `pkg/` est un engagement.** Tout ce qui entre dans `pkg/` devient une API publique qu'on ne peut plus casser sans bump majeur. Dans le doute, mettez-le dans `internal/`.

## Le flux de travail

1. **Ouvrez une issue** (ou commentez une issue existante) avant de commencer un travail conséquent
2. **Forkez** le dépôt et créez une branche depuis `main`
3. **Développez**, en ajoutant des tests pour ce qui peut casser
4. **Vérifiez** avec `task lint && task test` avant de pousser
5. **Ouvrez une Pull Request** vers `main`, en liant l'issue concernée

**Nommage des branches** — `<type>/<description-courte>`, par exemple :

```
feat/canary-deployments
fix/websocket-reconnect-loop
docs/contributing-guide
```

**Ce qu'on regarde en revue** — que le comportement soit correct et testé, que le code reste lisible par quelqu'un qui découvre le fichier, et que la PR fasse une seule chose. Une PR qui corrige un bug *et* reformate 40 fichiers est très difficile à relire : séparez-les.

Les Pull Requests nécessitent une revue et une CI verte avant merge sur `main`.

## Conventions de commit

Nous suivons [Conventional Commits](https://www.conventionalcommits.org) — le changelog en est généré automatiquement.

```
<type>(<scope>): <description à l'impératif>
```

**Types** — `feat`, `fix`, `docs`, `refactor`, `test`, `perf`, `build`, `ci`, `chore`

**Scopes** — `core`, `web`, `mobile`, `ai`, `proxy`, `api`, `cli`, `deploy`, `templates`

```bash
feat(core): ajouter la stratégie de déploiement canary
fix(web): corriger la boucle de reconnexion du stream de logs
docs(api): documenter les endpoints de métriques
```

Un changement cassant se signale par `!` après le scope (`feat(api)!: ...`) et une section `BREAKING CHANGE:` dans le corps du message. Le Core Engine suit le SemVer strict : toute rupture d'API publique implique un bump majeur.

## Style de code et tests

| Langage | Formatage | Lint | Tests |
|---|---|---|---|
| Go | `gofmt` | `golangci-lint` | `testify` |
| TypeScript/React | `prettier` | `eslint` | `vitest` |
| Dart/Flutter | `dart format` | `flutter analyze` | `flutter test` |

Des hooks pre-commit lancent ces vérifications automatiquement ; `task lint` fait la même chose à la demande.

**Sur les tests** — inutile de viser un pourcentage de couverture. Testez ce qui peut réellement casser en production : l'orchestration Docker, le parsing Compose, la génération de Caddyfile, la logique de rollback, l'enforcement RBAC. Un test qui vérifie qu'un getter retourne bien la valeur passée au constructeur ne protège de rien.

## La règle du contrat `api/`

**Aucun client HTTP ou gRPC n'est écrit à la main.** Le dossier `api/proto/` est la source de vérité ; les clients Go, TypeScript et Dart en sont générés.

Pour modifier une API :

1. Éditez le `.proto` correspondant dans `api/proto/forgenet/v1/`
2. Lancez `task gen` — les trois clients sont régénérés
3. Adaptez les appelants, committez le code généré avec le changement de contrat

Écrire un appel `fetch()` à la main peut sembler plus rapide sur le moment. À trois langages, c'est ainsi que les clients se désynchronisent silencieusement du serveur — et le bug ne se manifeste qu'en production.

## Proposer un template Marketplace

Un template = un dossier dans `templates/` contenant un `manifest.yaml` et un `docker-compose.yml`.

Critères d'acceptation : un logiciel open-source ou à licence claire, une image officielle ou reconnue avec **tag de version épinglé** (jamais `latest`), aucun secret en dur, des volumes déclarés pour toute donnée persistante, et un `README.md` court décrivant le service et ses variables d'environnement.

C'est la meilleure première contribution : elle n'exige aucune connaissance du code Go et elle est immédiatement utile aux utilisateurs.

## Signaler une faille de sécurité

**N'ouvrez pas d'issue publique pour une vulnérabilité.**

Écrivez à **security@lucidforgeafrica.com** avec la description, les étapes de reproduction et l'impact estimé. Nous accusons réception sous 72 heures et vous tenons informé du correctif. ForgeNet gère des secrets, des accès serveur et des données de production : ce canal est pris au sérieux.

## Licence de vos contributions

En proposant une contribution, vous acceptez qu'elle soit distribuée sous la [licence Apache 2.0](LICENSE), comme le reste du projet.

Gardez à l'esprit le modèle **open-core** de ForgeNet : le cœur du produit reste libre et gratuit, et certaines fonctions avancées (SSO/SAML, multi-tenant, HA) relèvent de l'édition Enterprise. Cette frontière est documentée dans la [matrice fonctionnelle](docs/ForgeNet_Matrice_Fonctionnelle_Complete_Dev.md) — n'hésitez pas à la questionner en issue si elle vous paraît mal placée.

---

Une question qui n'est pas couverte ici ? Ouvrez une [Discussion](https://github.com/forgelab-dev/lfa-forgenet/discussions) — il n'y a pas de question bête sur un projet en pré-alpha.
