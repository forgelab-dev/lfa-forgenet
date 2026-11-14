# Publier une version

Une version de ForgeNet, c'est trois choses publiées ensemble : le tag `vX.Y.Z`
avec sa release GitHub et ses notes, l'image du serveur
`ghcr.io/forgelab-dev/forgenet-server:X.Y.Z`, et le script `deploy/install.sh`
de ce tag, qui installe cette image par défaut.

## Le processus

1. **Les commits suivent les Conventional Commits** (`feat:`, `fix:`, `perf:`…).
   Ce sont eux qui forment les notes de version.
2. **release-please tient la PR de version à jour.** À chaque fusion sur `main`,
   le workflow `.github/workflows/release.yml` ouvre ou met à jour une PR
   « chore: publier la version X.Y.Z ». Elle contient :
   - le `CHANGELOG.md`, rangé en Fonctionnalités, Corrections, Performances,
     Refactorisation ;
   - la nouvelle version dans `internal/version/version.go`, `deploy/install.sh`
     et `deploy/.env.example` (lignes marquées `x-release-please-version`).
3. **Fusionner cette PR publie la version.** Le même workflow :
   - crée le tag et la release GitHub, marquée « pré-version » tant que
     ForgeNet est en alpha (`prerelease` dans `release-please-config.json`) ;
   - construit l'image pour **amd64 et arm64** ;
   - la publie avec deux étiquettes, `X.Y.Z` et `X.Y`, sans `latest` ;
   - y joint sa **provenance** (SLSA, `mode=max`) et son **SBOM** ;
   - la **signe avec cosign** sans clé, par l'identité OIDC du workflow, puis
     vérifie la signature.
   - joint à la release sa **fiche** `forgenet-release.json` : version, image et
     empreinte exacte de l'index signé, adresse des notes. Les installations
     existantes s'en servent pour se mettre à jour (ADR-0011).

Avant la 1.0, une fonctionnalité ou un changement cassant fait monter la version
mineure (0.1 → 0.2), une correction la version de correctif (0.1.0 → 0.1.1).

## Vérifier une image publiée

```sh
cosign verify ghcr.io/forgelab-dev/forgenet-server:0.1.0 \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  --certificate-identity https://github.com/forgelab-dev/lfa-forgenet/.github/workflows/release.yml@refs/heads/main
```

La provenance et le SBOM se lisent avec
`docker buildx imagetools inspect ghcr.io/forgelab-dev/forgenet-server:0.1.0 --format '{{ json .Provenance }}'`
(et `.SBOM`).

## À la première publication

Un nouveau paquet sur ghcr.io est **privé** par défaut : `install.sh` ne pourrait
pas télécharger l'image. Un administrateur de l'organisation le rend public une
fois, dans les réglages du paquet :
<https://github.com/orgs/forgelab-dev/packages/container/forgenet-server/settings>
(« Change visibility » → Public).

## Les choix

- **Pas de `latest`** : une installation se fait toujours à une version précise
  (`install.sh --version`), et une mise à jour reste un geste volontaire.
- **Actions épinglées par empreinte de commit** dans ce workflow : il publie et
  signe, une étiquette déplacée chez un tiers ne doit rien y changer.
- **Compilation croisée** : le dashboard et le binaire Go (SQLite en Go pur,
  sans cgo) sont compilés sur la machine de build pour la cible. Seule l'étape
  finale de l'image arm64 passe par l'émulation QEMU.
