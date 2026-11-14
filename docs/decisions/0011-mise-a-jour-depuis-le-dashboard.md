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
  "notes_url": "https://github.com/forgelab-dev/lfa-forgenet/releases/tag/v0.1.1"
}
```

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
service), qui lui survit le temps de la mise à jour puis s'efface. Il reçoit le
socket Docker, le dossier d'installation et le volume de données, rien de plus ;
le serveur pilote déjà Docker, ce conteneur n'ajoute donc aucun droit.

Étapes, journalisées dans le volume de données (`maj/etat.json`) pour que le
dashboard les suive et que le serveur relancé affiche le résultat :

1. **Verrou** : une seule mise à jour à la fois.
2. **Sauvegarde** : le serveur est arrêté, puis la base SQLite (avec son journal
   WAL) et `.env` sont copiés dans `/opt/forgenet/sauvegardes/<date>-<version>/`.
   Les cinq dernières sont gardées.
3. **Bascule** : `.env` reçoit la nouvelle version et son empreinte ;
   `docker compose up -d` recrée le serveur sur l'image exacte.
4. **Santé** : le nouveau serveur doit être « healthy » sous 60 secondes.
5. **Retour arrière** en cas d'échec : `.env` et la base sauvegardés sont remis
   en place, l'ancienne version redémarre, et la raison est notée.

`install.sh --update` lance le même conteneur : un seul code pour les deux
chemins.

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
- `docker-compose.yml` désignera l'image par version et empreinte.
- Les installations en 0.1.x passent une dernière fois par `install.sh` pour
  obtenir la version qui contient ce système ; ensuite, le dashboard suffit.
