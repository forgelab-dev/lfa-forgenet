# ADR-0008 : brancher les applications sur forgenet-edge sans les redémarrer

- Statut : proposé
- Date : 2026-11-13
- Contexte : ticket #10, deuxième PR (domaines d'une application déployée) ;
  suite de l'ADR-0007

## Contexte

L'ADR-0007 fait joindre les applications par Caddy sur un réseau Docker
partagé, `forgenet-edge`, où seul le service qui reçoit un domaine est branché.
Reste à dire comment ce branchement est fait, et comment il survit aux
redéploiements, aux redémarrages du serveur et aux erreurs.

## Options

1. **Branchement direct.** ForgeNet branche les conteneurs du service sur le
   réseau (`docker network connect`) dès l'ajout du domaine, puis les rebranche
   après chaque déploiement réussi du projet et au démarrage du serveur.
2. **Fichier Compose complémentaire.** ForgeNet passe à `docker compose up` un
   second fichier qui déclare le réseau pour le service. Le conteneur naît
   branché, mais il faut lire la configuration de l'application pour ne pas
   lui retirer son réseau par défaut, modifier le moteur Compose (ADR-0003), et
   recréer le service, donc le couper, à chaque ajout ou retrait de domaine.

## Décision

Option 1, retenue avant le code.

- **Un domaine** relie un nom à un service Compose d'un projet déployé et au
  port où ce service écoute dans son conteneur. Le nom est normalisé avant
  toute vérification : minuscules, sans point final, nom internationalisé
  converti en punycode (`golang.org/x/net/idna`, déjà dans les dépendances).
  Il passe ensuite la liste blanche de l'ADR-0007.
- **Le service doit être déployé** : au moins un conteneur de ce service dans
  la stack du projet. Un service en `network_mode: host` ou `none` ne peut pas
  rejoindre un réseau : Docker refuse, et l'ajout est refusé avec son message.
- **Alias DNS** du service sur le réseau : `<service>.<projet>.forgenet.internal`.
  Il est unique par stack et par service, lisible dans le Caddyfile, et
  `.internal` est réservé aux réseaux privés : il ne peut pas désigner un
  domaine public. Chaque étiquette fait 63 caractères au plus.
- **Le réseau** est créé par ForgeNet s'il manque, avec l'étiquette
  `forgenet.managed=edge`. Un réseau du même nom sans cette étiquette est
  refusé : ForgeNet ne branche pas d'applications sur un réseau qu'il n'a pas
  créé. Le conteneur de Caddy, s'il est indiqué dans la configuration, y est
  branché lui aussi.
- **La base est la source de vérité, et une modification ne laisse jamais de
  demi-état.**
  - Ajout : enregistrement, branchement, rechargement de Caddy. Si une étape
    échoue, ce qui a été fait est défait, et l'erreur revient à l'appelant.
  - Retrait : Caddy est d'abord rechargé sans le domaine. S'il refuse, rien ne
    change. Ensuite seulement le domaine est effacé, et le service débranché
    si plus aucun domaine ne le vise.
- **Rebranchement** après chaque déploiement réussi du projet, et
  **réconciliation complète au démarrage** : réseau, conteneur de Caddy,
  conteneurs des services, configuration de Caddy. Une modification, un
  rebranchement et une réconciliation ne se croisent jamais : un seul verrou.
- **Supprimer une stack garde ses domaines.** La route répond alors 502 jusqu'à
  un nouveau déploiement, qui la rebranche, ou jusqu'au retrait du domaine.
  Effacer la configuration d'un domaine reste un geste explicite.
- **Droits** : lister les domaines est ouvert aux membres, comme la liste des
  conteneurs ; ajouter et retirer sont réservés aux administrateurs. Chaque
  ajout et retrait est inscrit au journal d'audit.
- **Proxy désactivé** (`proxy.enabled`, désactivé par défaut) : le service
  répond `FAILED_PRECONDITION` avec un message clair, plutôt que d'enregistrer
  des domaines que rien ne sert.

## Conséquences

- Après un redéploiement qui recrée le conteneur, la route est indisponible
  entre le démarrage du nouveau conteneur et son rebranchement, une fraction
  de seconde, en plus de la coupure du redéploiement lui-même.
- Un conteneur recréé hors de ForgeNet (`docker compose up` lancé à la main)
  perd son branchement jusqu'au prochain déploiement ou redémarrage du
  serveur.
- Le moteur Compose et l'ADR-0003 ne changent pas.
