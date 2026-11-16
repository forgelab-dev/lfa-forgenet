# ADR-0013 : commande `forgenet` d'administration sur le serveur

- Statut : proposé
- Date : 2026-11-16
- Contexte : ticket #71 ; complète l'ADR-0009 (installation) et l'ADR-0011
  (mise à jour)

## Contexte

Administrer ForgeNet en terminal demande aujourd'hui de connaître ses
conteneurs et ses chemins : `docker exec forgenet-server forgenet-server
update check`, `docker logs forgenet-server`, `sudo bash install.sh --restore
…` avec un `install.sh` qu'il faut avoir gardé, `docker compose -f
/opt/forgenet/docker-compose.yml …`. Le but : une commande `forgenet` courte,
en français, qui fait ces gestes à la place de l'administrateur.

Trois contraintes :

- la commande doit marcher quand le serveur est arrêté ou malade : c'est là
  qu'on en a le plus besoin (`status`, `logs`) ;
- la mise à jour et le retour à une sauvegarde existent déjà dans
  `install.sh`, avec leurs garde-fous : ils ne doivent pas être recopiés ;
- une installation en place doit pouvoir se mettre à jour vers la version qui
  apporte la commande.

## Options

1. **Commande dans l'image du serveur**, appelée par un petit relais sur
   l'hôte. Elle suivrait l'image, donc chaque mise à jour, mais ne marcherait
   plus quand le conteneur est arrêté, et l'image n'a ni Docker ni le dossier
   d'installation.
2. **Commande listée dans la fiche de release** (`files` de
   `forgenet-release.json`), que le moteur de mise à jour écrirait comme
   `docker-compose.yml`. Le moteur des versions en place refuse tout nom de
   fichier qu'il ne connaît pas (ADR-0011) : la fiche de la version suivante
   serait refusée par toutes les installations en 0.4.x, qui ne pourraient
   plus se mettre à jour.
3. **Script installé par `install.sh`**, qui range aussi une copie de
   lui-même dans le dossier d'installation.

## Décision

Option 3.

- **Un script bash sur l'hôte.** `deploy/forgenet` est installé en
  `/usr/local/bin/forgenet` (0755). Il ne parle qu'à Docker : `status` et
  `logs` marchent serveur arrêté. Le dossier d'installation est inscrit dans
  le script par `install.sh` (option `--dir`), ce qui justifie la validation
  plus stricte de `--dir` : lettres, chiffres, `. _ - /` seulement.
- **La mise à jour et le retour à une sauvegarde passent par `install.sh`.**
  `install.sh` range une copie de lui-même dans le dossier d'installation
  (0700). `forgenet update` et `forgenet restore` la relancent avec
  `--update` ou `--restore` : un seul code pour ces deux gestes.
- **Droits.** Ce qui écrit ou lit le dossier d'installation (`update`,
  `backups`, `restore`) demande root ; le reste passe par Docker et marche
  aussi pour un membre du groupe `docker`.
- **Remise à niveau.** `install.sh` installe la commande de sa version.
  `install.sh --update`, donc `forgenet update`, installe celle de la version
  en place une fois la mise à jour faite, ou quand il n'y a rien à mettre à
  jour. Une mise à jour lancée depuis le dashboard ne touche pas à la
  commande : `forgenet version` signale l'écart et propose `sudo forgenet
  update`. Le script relancé est remplacé par un renommage, jamais réécrit
  sous lui pendant qu'il tourne.
- **Sans la commande, tout marche.** Une version publiée avant elle ne la
  contient pas : `install.sh` le signale et continue.

## Conséquences

- Les fichiers de la commande sont lus depuis le tag de la version, comme
  `docker-compose.yml` lors d'une installation : même confiance que le script
  d'installation lui-même.
- Une installation antérieure reçoit la commande en lançant une fois le
  `install.sh` d'une version qui la contient.
- Testée sans Docker ni root par `deploy/tests/forgenet_test.sh` (doublures de
  `docker` et `id`), lancé par la CI avec shellcheck.
