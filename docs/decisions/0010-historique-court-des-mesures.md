# ADR-0010 : historique court des mesures, gardé en mémoire

- Statut : accepté (revue de la PR #55, 14/11/2026)
- Date : 2026-11-14
- Contexte : ticket #49 (écran Surveillance du dashboard)

## Contexte

L'ADR-0005 diffuse les mesures en direct par WebSocket, un collecteur par flux.
Le contrat prévoit aussi l'historique (`QuerySystemMetrics`,
`QueryContainerMetrics`), mais le serveur ne garde aucune mesure : l'ADR-0005
l'avait laissé hors de son ticket. L'écran Surveillance a besoin :

- de courbes sur les 15 dernières minutes ou la dernière heure, dès l'ouverture ;
- de combler le trou laissé par une coupure de la liaison (connexion lente ou
  instable) ;
- de suivre **tous** les conteneurs à la fois. Avec l'ADR-0005 seul, cela veut
  dire une liaison WebSocket par conteneur : coûteux sur une connexion lente, et
  autant d'appels au démon Docker que de conteneurs et d'onglets ouverts.

## Décision

- **Un relevé de fond** (`metrics.Recorder`) mesure lui-même l'hôte toutes les
  5 s et chaque conteneur en cours d'exécution toutes les 10 s, au plus quatre
  appels simultanés au démon. Il démarre et s'arrête avec les autres serveurs.
- **La dernière heure est gardée en mémoire**, dans des tampons circulaires de
  taille fixe : 721 mesures pour l'hôte, 361 par conteneur. L'ordre de grandeur
  est de quelques centaines de Ko, quelques Mo au plus avec le plafond de 200
  conteneurs suivis. Un conteneur arrêté garde son historique une heure, puis
  il est oublié.
- **Rien n'est écrit sur le disque** : l'historique repart de zéro au
  redémarrage du serveur, et l'écran le dit. Une base SQLite garderait
  l'historique au redémarrage, mais elle écrirait toutes les 5 s pour un besoin
  de quelques minutes. Un historique long (jours, semaines) relèvera d'un
  ticket à part, avec sa rétention et son agrégation.
- **Les lectures passent par le contrat**, servies en gRPC et en Connect
  (ADR-0006) : `GetSystemMetrics` et `GetContainerMetrics` (dernière mesure),
  `QuerySystemMetrics` et `QueryContainerMetrics` (historique). Elles sont
  ouvertes aux membres, comme les flux. Les règles des requêtes :
  - sans début, la plage couvre la fenêtre entière ; sans fin, elle s'arrête à
    maintenant. Un début postérieur à la fin est refusé ;
  - le pas suit le rythme du relevé, sans dépasser `max_points` (1 000 au
    plus) ; un pas demandé plus fin est relevé d'autant ;
  - les tranches sont alignées sur l'époque Unix, pour que deux requêtes
    successives se raccordent. Une valeur instantanée y prend sa moyenne, un
    compteur cumulé sa dernière valeur ;
  - les métriques d'extension (`custom`) sont refusées tant qu'aucune extension
    n'existe ; les métriques propres à l'hôte (disque, uptime) sont refusées
    pour un conteneur ;
  - un conteneur se désigne par son identifiant complet, ou par un préfixe
    d'au moins 12 caractères qui n'en désigne qu'un, comme dans `docker ps`.
- **Une route WebSocket pour tous les conteneurs** :
  `GET /ws/v1/metrics/containers` envoie chaque passage du relevé de fond en un
  message `StreamAllContainerMetricsResponse` (ajouté au contrat). Un conteneur
  absent du message n'est plus en cours d'exécution. Le rythme est celui du
  relevé ; le paramètre `interval` y est refusé. Authentification, refus en HTTP
  et codes de fermeture sont ceux de l'ADR-0005. Le démon Docker n'est pas
  sollicité de plus, quel que soit le nombre d'abonnés.
- **Les flux existants ne changent pas** : `/ws/v1/metrics/system` et
  `/ws/v1/metrics/containers/{id}` gardent leurs collecteurs par flux et leurs
  intervalles réglables. `StreamSystemMetrics` et `StreamContainerMetrics`
  restent non implémentés en gRPC et en Connect, avec un message qui renvoie
  vers la route WebSocket.

## Conséquences

- Le dashboard ouvre deux liaisons pour l'écran Surveillance (l'hôte, tous les
  conteneurs), plus une pour le conteneur dont il montre le détail.
- Après une coupure, il comble le trou par une requête d'historique, puis
  reprend le direct.
- Le premier relevé d'un conteneur attend une seconde, car le démon prend deux
  mesures pour calculer le processeur. Chaque passage a donc 5 s de marge en
  plus de son rythme. Les passages ne se chevauchent jamais : un démon lent
  espace les mesures au lieu de les empiler.
- Un serveur qui redémarre perd son historique : c'est le prix d'une solution
  sans écriture, accepté par Mike le 14/11/2026.
