# ADR-0005 : métriques en direct par WebSocket

- Statut : accepté (revue de la PR #38, 09/11/2026)
- Date : 2026-11-09
- Contexte : ticket #11, métriques par hôte et par conteneur

## Contexte

Le dashboard web et l'application mobile affichent le processeur, la mémoire, le
disque et le réseau en temps réel. Le contrat `metrics.proto` décrit déjà les
messages et des flux gRPC (`StreamSystemMetrics`, `StreamContainerMetrics`). Mais
un navigateur ne sait pas ouvrir un flux gRPC, et le ticket demande un flux
WebSocket compressé. Il faut aussi transmettre le jeton : l'API WebSocket des
navigateurs ne permet pas de poser un en-tête `Authorization`.

## Décision

- **Deux routes sur le serveur HTTP** : `GET /ws/v1/metrics/system` et
  `GET /ws/v1/metrics/containers/{id}`, paramètre `interval` en secondes (1 à
  60 ; absent ou 0 : 5 s pour l'hôte, 2 s pour un conteneur, comme le contrat).
- **Les messages restent ceux du contrat** : `StreamSystemMetricsResponse` et
  `StreamContainerMetricsResponse`, en JSON (`protojson`), un message texte par
  mesure. Le dashboard et l'application les décodent avec leur code généré : pas
  de second format à maintenir.
- **Compression permessage-deflate** avec reprise du contexte : les messages se
  répètent d'une mesure à l'autre et se compressent bien.
- **Le jeton** passe par l'en-tête `Authorization: Bearer` (CLI, mobile) ou par
  le sous-protocole `forgenet.bearer.<jeton>` (navigateur), à côté de
  `forgenet.v1`, seul sous-protocole renvoyé. Jamais dans l'URL, qui finit dans
  les journaux d'accès et l'historique.
- **Les refus sont des réponses HTTP**, avant la poignée de main : 401, 403, 400,
  404 (conteneur inconnu, la première mesure est prise avant), 503 (plus de
  64 flux, ou vérification impossible). Une origine étrangère reçoit 403 : un
  site tiers ne peut pas lire les métriques avec la session d'un visiteur.
- **Accès** : membres et admins, comme la lecture des conteneurs en gRPC.
- **Pendant le flux**, le jeton est revérifié chaque minute. Codes de
  fermeture : `4401` jeton révoqué ou expiré, `4404` conteneur disparu, `1001`
  arrêt du serveur, `1011` mesure impossible.
- **Les compteurs réseau et disque sont cumulés** (comme le contrat) : le client
  calcule les débits. Seul le processeur est un instantané ; chaque flux a son
  propre relevé précédent.
- **Bibliothèque** `github.com/coder/websocket` : compression et contextes
  intégrés, maintenue. `gorilla/websocket` demande de gérer soi-même la
  compression et l'annulation.

## Conséquences

- Les RPC `GetSystemMetrics` et `StreamSystemMetrics` (et leurs pendants pour un
  conteneur) ne sont pas encore servis en gRPC ; ils réutiliseront
  `internal/metrics` le jour où un client gRPC en aura besoin.
- `QuerySystemMetrics` et `QueryContainerMetrics` (historique) demandent de
  stocker les mesures : hors de ce ticket.
- Lancé dans un conteneur, le serveur ne voit que son propre `/proc` : il faut y
  monter celui de l'hôte et poser `HOST_PROC` (variable de gopsutil), et régler
  `metrics.storage_path` sur le disque à surveiller.
- Le serveur de développement du dashboard (Vite, autre port) sera une origine
  étrangère : il faudra l'autoriser par `ws.Config.OriginPatterns` dans le
  ticket #9.
- En cgroup v2, le démon Docker ne donne pas le nombre d'opérations disque d'un
  conteneur : `read_ops_total` et `write_ops_total` restent à 0.
