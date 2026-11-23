# ADR-0019 : valeurs des installations du Marketplace

- Statut : proposé (PR du lot 2 de #116)
- Date : 2026-11-23
- Contexte : issue #116 (installer un template en un clic) ; complète
  l'ADR-0018 (format des templates) sans le contredire ; s'appuie sur les
  ADR-0003 (stacks Compose), 0004 (déploiements) et 0015 (secrets chiffrés
  en base)

## Contexte

L'ADR-0018 dit d'où viennent les valeurs d'une installation (le formulaire,
ou un tirage au hasard pour un secret laissé vide) et par où elles passent
(l'environnement de la commande `docker compose`). Il restait à décider où
le serveur les garde, ce qui se passe quand l'application est redéployée,
et comment l'administrateur apprend un mot de passe que le serveur a tiré
lui-même alors qu'un secret n'est « jamais réaffiché ».

## Décision

### Une installation est un déploiement ordinaire

`MarketplaceService.InstallTemplate` vérifie les valeurs (`Manifest.Resolve`),
tire les secrets laissés vides, puis crée un déploiement du fichier Compose
**du catalogue** : le navigateur n'envoie jamais ce fichier, il ne peut donc
pas le modifier au passage. Le déploiement suit ensuite la file, les étapes
et le journal de l'ADR-0004, et la page Déploiements le montre comme les
autres, avec le template installé (`Deployment.template`).

Une installation refuse un nom d'application déjà pris, par un déploiement ou
par une stack en place (`ALREADY_EXISTS`) : `CreateDeployment` met à jour une
stack existante, une installation ne doit jamais la remplacer.

### Les valeurs vivent avec chaque déploiement, chiffrées ensemble

La colonne `deployments.variables_sealed` (migration 0009) garde toutes les
valeurs de l'installation, en JSON, chiffrées par `sealed.Box` avec l'usage
`deployment.variables`. Les valeurs qui ne sont pas secrètes sont chiffrées
aussi : c'est plus simple qu'un mélange de clair et de chiffré, et rien ne
les affiche après l'installation.

Sans clé des secrets (`FORGENET_SECRET_KEY`, écrite par `install.sh`), un
template qui demande des variables ne s'installe pas
(`FAILED_PRECONDITION`) : les valeurs ne sont jamais gardées en clair.

Elles sont gardées par déploiement plutôt que par projet : le retour en
arrière (`RollbackDeployment`, pas encore codé) reprendra celles de la
version visée.

### Un redéploiement reprend les valeurs

Un fichier redéployé à la main depuis la page Déploiements n'a pas de
valeurs : sans elles, `${POSTGRES_PASSWORD}` deviendrait vide. Un déploiement
créé sans variables reprend donc celles du déploiement précédent du projet,
et son template. Le chiffré est copié tel quel ; il n'est ouvert que pour
vérifier le fichier avec ses valeurs.

### Seules deux commandes lisent le fichier

`config` (vérification) et `up` (mise en service) reçoivent les valeurs dans
leur environnement. Les autres commandes du moteur (`ps`, `ls`, `down`) ne
lisent pas le fichier, et la relance après un redémarrage est faite par
Docker lui-même (ADR-0012) : elles n'en ont pas besoin. Le moteur refuse un
nom réservé (`PATH`, `HOME`, préfixes `DOCKER_`, `COMPOSE_`, `FORGENET_`)
même si le manifeste l'interdit déjà.

### Les secrets ne s'écrivent nulle part en clair

- Le journal du déploiement, son message d'échec, l'erreur rendue au client
  et le journal du serveur **masquent** les valeurs secrètes (`••••••`) :
  Compose peut citer une valeur dans une erreur.
- L'audit (`deployment.create`) garde le template et les **noms** des
  variables, jamais une valeur.

### Un secret généré n'est montré qu'une fois

`InstallTemplate` renvoie les secrets tirés au hasard dans sa réponse, et
seulement quand l'installation vient d'être créée. Ni une relecture du
déploiement, ni la même demande renvoyée avec sa clé d'idempotence ne les
redonnent, et aucun appel ne permet de les réafficher, même après avoir
retapé son mot de passe (décision du porteur du projet, 23/11). Le dashboard
les montre à la fin de l'installation, avec la consigne de les noter.

### Un catalogue de test

Le réglage `marketplace.dir` (`FORGENET_MARKETPLACE_DIR`, chemin absolu)
remplace le catalogue embarqué par un dossier, lu au démarrage avec les
mêmes vérifications que le catalogue embarqué. Il sert aux tests de bout en
bout, le catalogue réel pouvant être vide, et aux contributeurs qui essaient
leur template dans leur ForgeNet avant la PR. Le serveur l'annonce dans son
journal au démarrage.

## Options écartées

- **Réafficher un secret sur demande**, après avoir retapé son mot de passe :
  plus confortable, mais il faut alors pouvoir le déchiffrer pour le
  montrer, ce qui contredit « jamais réaffiché ».
- **Un champ de plus dans `CreateDeployment`** : le navigateur enverrait le
  fichier du template et pourrait le modifier ; la vérification des valeurs
  et le tirage des secrets se feraient hors du serveur.
- **Garder les valeurs avec le projet** : il n'y a pas de table des projets,
  et un retour en arrière perdrait les valeurs de la version visée.

## Conséquences

- Une base installée garde son mot de passe à chaque redéploiement.
- Changer la clé des secrets rend les installations impossibles à
  redéployer : leur déploiement échoue avec un message qui nomme
  `FORGENET_SECRET_KEY`. Les conteneurs en place continuent de tourner.
- Un administrateur qui n'a pas noté un mot de passe généré doit le changer
  dans l'application elle-même ; la modification d'une variable depuis le
  dashboard viendra plus tard (ADR-0018).
