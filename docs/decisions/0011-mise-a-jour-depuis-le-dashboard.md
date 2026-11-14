# ADR-0011 : mise à jour d'une installation depuis le dashboard

- Statut : proposé
- Date : 2026-11-14
- Contexte : ticket #62, après la publication de la v0.1.0 (#16)

## Contexte

Relancer `deploy/install.sh` d'une nouvelle version met déjà ForgeNet à jour :
le script réécrit `docker-compose.yml` et la version dans `.env`, télécharge
l'image et recrée les conteneurs ; le serveur applique ses migrations SQL au
démarrage. Mais cette mise à jour est sans filet :

- aucune sauvegarde de la base avant de la migrer ;
- aucun retour arrière si la nouvelle version ne démarre pas ;
- rien n'empêche d'installer une version plus ancienne sur une base déjà
  migrée, ce qui la casserait ;
- l'administrateur n'apprend l'existence d'une version qu'en allant la chercher,
  et doit taper une commande en root à chaque fois.

Le but : une mise à jour en un clic depuis le dashboard, sans risque pour les
données, et la même chose en une commande pour qui préfère le terminal.

## Décision

### La release publie sa fiche

Le workflow de release joint à chaque release GitHub un fichier
`forgenet-release.json` :

```json
{
  "schema": 1,
  "version": "0.1.1",
  "image": "ghcr.io/forgelab-dev/forgenet-server",
  "digest": "sha256:…",
  "notes_url": "https://github.com/forgelab-dev/lfa-forgenet/releases/tag/v0.2.0",
  "files": {
    "docker-compose.yml": "<sha256 hexadécimal>",
    "Caddyfile.base": "<sha256 hexadécimal>"
  }
}
```

`files` (depuis la 0.2.0) donne l'empreinte des fichiers d'installation du tag :
une version qui fait évoluer `docker-compose.yml` l'apporte avec elle, et le
moteur refuse un fichier qui ne correspond pas. Seuls ces deux noms sont admis.
Les fiches sans `files` (0.1.1) laissent les fichiers en place.

L'empreinte est celle de l'index multi-architectures que la même tâche vient de
construire, de signer avec cosign et de vérifier. Une installation tire ainsi
**cette image-là**, pas seulement une étiquette qui pourrait être déplacée.

### Le serveur détecte, l'administrateur décide

- Une fois par jour, le serveur lit la dernière release de
  `forgelab-dev/lfa-forgenet` sur l'API GitHub, en HTTPS, avec un User-Agent
  générique (`forgenet/<version>`). C'est le seul appel sortant ; il ne transmet
  rien sur l'installation. `updates.check: false` le coupe.
- Une version plus récente (comparaison semver) est annoncée aux seuls
  administrateurs, avec ses notes. **Rien ne se met à jour tout seul.**
- Une version plus ancienne n'est jamais proposée : la base est déjà migrée et
  les migrations ne se défont pas. Les pré-versions (`-rc`…) ne sont proposées
  qu'à qui les a demandées.

### Un conteneur de mise à jour, à part

Un serveur ne peut pas se remplacer lui-même : il serait arrêté au milieu de
l'opération. Au clic, le serveur lance donc un conteneur séparé
(`forgenet-updater`), à partir de **sa propre image** (un code connu, déjà en
service), qui lui survit le temps de la mise à jour. Il reste ensuite, arrêté,
pour ses journaux, et la mise à jour suivante le retire. Il reçoit le
socket Docker, le dossier d'installation et le volume de données, rien de plus ;
le serveur pilote déjà Docker, ce conteneur n'ajoute donc aucun droit.

Étapes, journalisées dans le volume de données (`maj/etat.json`) pour que le
dashboard les suive et que le serveur relancé affiche le résultat :

1. **Verrou** : une seule mise à jour à la fois.
2. **Téléchargement**, avant tout arrêt : l'image `version@empreinte` et les
   fichiers d'installation, vérifiés contre la fiche. Un registre injoignable ou
   un fichier altéré annule tout sans coupure.
3. **Sauvegarde** : le serveur est arrêté, puis la base SQLite (avec son journal
   WAL), `.env` et les fichiers d'installation sont copiés dans
   `/opt/forgenet/sauvegardes/<date>-<version>/`. Les cinq dernières sont
   gardées.
4. **Bascule** : les nouveaux fichiers sont écrits, `.env` reçoit la nouvelle
   version et son empreinte, `docker compose up -d` recrée le serveur sur
   l'image exacte (et les autres services si `docker-compose.yml` a changé).
5. **Santé** : le nouveau serveur doit être « healthy » sous 90 secondes (son
   premier contrôle arrive après 30). Un conteneur arrêté ou « unhealthy »
   échoue tout de suite.
6. **Retour arrière** en cas d'échec : fichiers, `.env` et base sauvegardés sont
   remis en place, l'ancienne version redémarre (avec au moins deux minutes pour
   redevenir saine), et la raison est notée.

En ligne de commande : `forgenet-server update check` dit si une version est
disponible, `forgenet-server update start` lance le conteneur de mise à jour, qui
exécute `forgenet-server update apply`. `install.sh --update` passe par
`update start` dès que la version installée le connaît : un seul code pour le
terminal et le dashboard. `--force` réinstalle une version égale ou plus
ancienne, sur décision explicite.

### Sécurité

- Réservé aux administrateurs, et inscrit au journal d'audit.
- Le navigateur ne choisit ni la version ni l'image : il demande « mettre à
  jour », le serveur installe la dernière release officielle à l'empreinte de sa
  fiche.
- La fiche vient de l'API GitHub en HTTPS, et l'image désignée par son empreinte
  est signée par cosign (vérifiable avec la commande de `docs/RELEASING.md`).

## Alternatives écartées

- **Mise à jour automatique (façon Watchtower)** : un serveur qui change seul,
  sans sauvegarde ni décision humaine, n'est pas acceptable pour des données de
  production.
- **Le serveur qui se remplace lui-même** : il serait tué en cours de route,
  sans personne pour revenir en arrière.
- **Une connexion SSH depuis le dashboard** : il faudrait y confier des
  identifiants root.
- **Signer la fiche elle-même** : utile, mais elle ne fait que désigner une image
  déjà signée ; vérifier une signature Sigstore côté serveur demanderait une
  dépendance lourde. À reconsidérer si la fiche porte un jour plus qu'une
  empreinte.

## Conséquences

- Le workflow de release obtient le droit d'écrire dans la release
  (`contents: write`) pour y joindre la fiche.
- `docker-compose.yml` désigne l'image par version et empreinte
  (`FORGENET_IMAGE_DIGEST`, facultative pour les installations antérieures).
- Les installations en 0.1.x passent une dernière fois par `install.sh` pour
  obtenir la version qui contient ce système ; ensuite, le dashboard suffit.
