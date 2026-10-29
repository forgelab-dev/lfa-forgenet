# ADR-0003 : les stacks Compose passent par la commande `docker compose`

- Statut : proposé
- Date : 2026-10-29
- Contexte : ticket #7, deuxième PR (moteur `internal/orchestrator/compose`)

## Contexte

Un déploiement de source `compose` reçoit le contenu d'un fichier Docker Compose
(`ComposeSource.compose_content`) et doit déployer, mettre à jour et supprimer
la stack qu'il décrit. Trois façons de faire :

1. **Réimplémenter Compose** sur `orchestrator.Engine` avec `compose-go` pour
   lire le fichier : réseaux, volumes, ordre de `depends_on`, sondes de santé,
   recréation des seuls services modifiés… C'est des milliers de lignes, et
   chaque écart avec Compose surprendra l'utilisateur.
2. **Embarquer `github.com/docker/compose`** comme bibliothèque. Il ramène
   `github.com/docker/docker` et ses failles sans correctif, retirés à la PR #28,
   et un arbre de dépendances très lourd.
3. **Piloter la commande `docker compose`** (plugin officiel, livré avec Docker
   Engine).

## Décision

- Le moteur Compose lance `docker compose`, avec un **nom de projet par stack**
  (`-p`). Le fichier est passé sur l'entrée standard (`-f -`) : rien n'est écrit
  sur le disque, et la source de vérité reste le contenu stocké avec le
  déploiement.
- **Déployer et mettre à jour sont la même opération** :
  `up --detach --wait --remove-orphans`. Compose ne recrée que les services
  modifiés et retire ceux qui ont disparu du fichier. `--wait` attend que les
  services tournent, et qu'ils soient sains quand ils ont une sonde de santé.
- **Supprimer** : `down --remove-orphans`. Les volumes nommés sont gardés, sauf
  demande explicite : ce sont les données de l'application.
- **Le fichier n'a pas accès à l'environnement du serveur.** Compose remplace
  `${VAR}` par les variables de son propre environnement : la commande ne reçoit
  que celles dont elle a besoin (`PATH`, `HOME`, les `DOCKER_*`), sinon un
  fichier pourrait recopier un secret du serveur dans un conteneur.
- **Le dossier de travail est un dossier vide propre à ForgeNet.** Un chemin
  relatif (`build: .`, `env_file`, montage `./data`) ne lit donc jamais les
  fichiers du serveur : le fichier doit se suffire à lui-même (images publiées,
  volumes nommés).
- La commande vise le même démon que le moteur de conteneurs (`docker.host`).

## Conséquences

- Le comportement est exactement celui de Compose, versions du format comprises,
  sans dépendance Go supplémentaire.
- L'hôte doit avoir le plugin `docker compose` (v2 ou plus). Son absence est une
  erreur claire (`ErrUnavailable`), qui pourra s'afficher dans `/readyz`.
- Les erreurs viennent de la sortie de la commande : le moteur les renvoie
  telles quelles, sans les interpréter.
- Un fichier Compose donne un accès complet à l'hôte (montages, mode
  privilégié). Seuls les administrateurs pourront déployer une stack (service
  `deploy`, PR suivante) ; une politique plus fine (refuser `privileged`, les
  montages de l'hôte) demande une nouvelle décision.
- Les sources dont le fichier s'appuie sur des fichiers voisins (`build`,
  `env_file`) ne sont pas prises en charge : elles viendront avec les sources Git.
