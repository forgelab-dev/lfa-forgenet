# ADR-0012 : relancer les applications après un redémarrage du serveur

- Statut : proposé
- Date : 2026-11-16
- Contexte : ticket #70, constat du test réel du 14/11 sur un VPS ; complète
  l'ADR-0003

## Contexte

ForgeNet applique le `docker-compose.yml` de l'utilisateur tel quel
(ADR-0003). Un service sans `restart:` reçoit la politique `no` de Docker, qui
ne le relance pas au démarrage du démon. Après un redémarrage du VPS, le
14/11, le serveur ForgeNet et Caddy sont revenus seuls en une vingtaine de
secondes, mais toutes les applications sont restées arrêtées et leurs domaines
répondaient 502. Une maintenance de l'hébergeur ou une coupure de courant met
donc hors ligne les sites d'un client jusqu'à une intervention à la main.

Trois exigences : une valeur écrite par l'utilisateur est respectée ; le
fichier enregistré et affiché reste le sien ; un service arrêté depuis le
dashboard reste arrêté.

## Options

1. **Réécrire le fichier de l'utilisateur** (lecture YAML, ajout de `restart:`,
   écriture). Il faut alors reproduire les règles de Compose : ancres et clé de
   fusion `<<`, `extends`, profils, `deploy.restart_policy`. Un service qui
   hérite sa politique d'une ancre serait vu à tort comme sans politique.
2. **`docker update --restart` après le déploiement.** Le fichier ne change
   pas, mais la politique ne fait pas partie de la configuration Compose :
   chaque service recréé par un déploiement la perd jusqu'au passage suivant,
   et l'état du conteneur ne distingue pas `no` écrit de `no` par défaut.
3. **Fichier Compose complémentaire, d'après le modèle résolu.** Avant `up`,
   ForgeNet lit le modèle par `docker compose config --format json`, où
   Compose a déjà appliqué ancres, `extends` et profils. Il écrit un second
   fichier qui ne donne `restart: unless-stopped` qu'aux services sans
   politique, et le passe à `up` après l'entrée standard.

## Décision

Option 3.

- **Services concernés.** Un service reçoit `unless-stopped` s'il n'a ni
  `restart:` ni `deploy.restart_policy` (que Compose applique quand
  `restart:` manque). Une tâche ponctuelle, qu'un autre service attend avec la
  condition `service_completed_successfully` (une migration, par exemple), est
  laissée à `no` : relancée sans fin, elle tournerait en boucle. Un service
  d'un profil inactif n'est pas dans le modèle, donc pas dans le complément.
- **`unless-stopped` plutôt que `always`.** Docker mémorise un arrêt
  volontaire : un service arrêté depuis le dashboard ne repart pas au
  démarrage suivant.
- **Le fichier de l'utilisateur ne change pas.** Il passe toujours par l'entrée
  standard et c'est lui qui est enregistré et affiché. Le complément est du
  JSON produit par l'encodeur (les noms de service viennent de Compose et sont
  échappés), écrit en `0600` dans un dossier `0700` distinct du dossier de
  travail, qui reste vide, et supprimé après la commande, y compris en cas
  d'échec.
- **Dossier du projet fixé.** L'entrée standard n'ayant pas de dossier, Compose
  prendrait le dossier du complément comme dossier du projet, et un chemin
  relatif comme `./donnees` changerait d'emplacement. `up` reçoit donc
  `--project-directory` avec le dossier de travail : les chemins se résolvent
  exactement comme avant.
- **Une seule lecture.** La lecture du modèle remplace la validation que
  `Deploy` faisait déjà (`config --quiet`) : un fichier refusé reste un
  `ErrInvalidFile`, avant que la stack en place soit touchée.
- **Transparence.** Le journal du déploiement dit à quels services la politique
  a été ajoutée.
- **Stacks déjà en place.** Elles reçoivent la politique à leur prochain
  déploiement, qui recrée une fois les services concernés (leur configuration
  change). La documentation dit de les redéployer une fois après la mise à
  jour ; aucun rattrapage automatique, qui couperait les applications sans
  qu'on le demande.

## Pourquoi pas le branchement direct, comme pour les réseaux

L'ADR-0008 a écarté un fichier complémentaire pour `forgenet-edge`, parce
qu'ajouter ou retirer un domaine aurait recréé, donc coupé, le service. Ici la
politique se fixe une fois pour toutes au déploiement, moment où Compose crée
ou recrée déjà les conteneurs : le complément ne provoque aucune coupure de
plus. Et une politique posée après coup (option 2) serait perdue à chaque
recréation.

## Conséquences

- Un déploiement lance `config --format json` au lieu de `config --quiet` ;
  `Validate` garde `--quiet`.
- L'étiquette `com.docker.compose.project.config_files` d'un conteneur cite le
  chemin du complément, supprimé depuis. Elle citait déjà `-`, l'entrée
  standard ; ForgeNet ne relit jamais ces fichiers (`ps` et `down` passent par
  le nom du projet).
- Testé sur un vrai démon (`TestEngine_RealComposeRestartPolicy`) : politique
  de chaque service, relance par Docker d'un service dont le processus
  s'arrête, aucun redémarrage pour `restart: "no"`, redéploiement sans
  recréation. Le redémarrage complet d'un serveur se vérifie sur le VPS de
  test une fois la version publiée.
