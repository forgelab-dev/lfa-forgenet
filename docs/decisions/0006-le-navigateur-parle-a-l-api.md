# ADR-0006 : le navigateur parle à l'API par Connect

- Statut : proposé
- Date : 2026-11-09
- Contexte : ticket #9 (dashboard web), à trancher avant son code

## Contexte

Le serveur expose ses services en gRPC natif (grpc-go, port 9090) : le CLI et
l'application Flutter s'en servent. Un navigateur ne sait pas parler ce
protocole : `fetch` ne donne accès ni aux trames HTTP/2 ni aux trailers. Le
dashboard (#9) doit pourtant appeler les mêmes services, y compris les flux
(journaux d'un conteneur, d'un déploiement), avec un client issu du codegen de
`api/` et sans client écrit à la main.

Les services sont écrits pour grpc-go (`forgenetv1.*Server`), et le contrôle
d'accès est un intercepteur grpc-go (`auth.Guard`).

## Options

1. **Connect par vanguardgrpc** (`connectrpc.com/vanguard/vanguardgrpc` sur
   `connectrpc.com/connect/v2`). Les implémentations grpc-go existantes
   s'enregistrent telles quelles sur un `connect.Server`, monté sur le serveur
   HTTP par `connecthttp.Mount` : il répond en Connect, en gRPC-Web et en gRPC.
   Les en-têtes de la requête deviennent les métadonnées gRPC que les services
   lisent déjà.
2. **connect-go v1 avec des adaptateurs écrits à la main** pour chaque service :
   une seconde couche à maintenir à chaque RPC ajoutée.
3. **Proxy gRPC-Web (Envoy)** devant le serveur : un service de plus à
   installer, configurer et surveiller chez chaque utilisateur, ce qui va
   contre l'installation en une commande (#15).
4. **improbable-eng/grpc-web** : projet archivé, écarté.

## Décision

Option 1.

- **Les services gRPC restent la seule implémentation**, et le port gRPC (9090)
  ne change pas pour le CLI et l'application mobile.
- **Le serveur HTTP sert aussi l'API** sous les chemins du contrat
  (`/forgenet.v1.ContainerService/ListContainers`…), à côté des sondes et des
  flux WebSocket (ADR-0005).
- **Le contrôle d'accès est commun aux deux portes.** `auth.Guard` expose sa
  décision (méthode et jeton donnent une identité, ou un refus), et la même
  politique (`DefaultPolicy`) l'applique par un intercepteur grpc-go et par un
  intercepteur Connect. Les intercepteurs grpc-go ne sont pas rejoués par
  vanguardgrpc : sans ce second branchement, la porte HTTP serait ouverte.
- **Le dashboard utilise `@connectrpc/connect-web`** avec le protocole Connect
  et le code protobuf-es déjà généré dans `web/src/gen` : pas de client écrit
  à la main. Le protocole Connect transporte les flux serveur en HTTP/1.1,
  derrière n'importe quel proxy. Le contrat n'a aucun flux client, que les
  navigateurs ne savent pas envoyer.
- **Une seule origine, pas de CORS.** En développement, le serveur Vite relaie
  `/forgenet.v1.*` et `/ws` vers le serveur ForgeNet ; en production, le
  dashboard et l'API sont servis sous le même nom de domaine. Les cookies et
  les en-têtes CORS restent donc hors du sujet, et le refus des origines
  étrangères de l'ADR-0005 s'applique sans exception.

## Vérification préalable

Essai jetable du 09/11 (non commité) : `ContainerService` du contrat,
implémentation grpc-go, enregistrée par `vanguardgrpc.NewServiceRegistrar` et
montée par `connecthttp.Mount`.

- appel Connect JSON sans jeton : `{"code":"unauthenticated"}` ; avec
  `Authorization: Bearer …` lu dans les métadonnées gRPC : la réponse ;
- `StreamContainerLogs` en Connect (`application/connect+json`) puis en
  gRPC-Web : trois messages et la fin de flux, en HTTP/1.1.

## Conséquences

- Deux dépendances : `connectrpc.com/vanguard` v0.5.0 (avant la 1.0) et
  `connectrpc.com/connect/v2` v2.0.x, toutes deux de Buf, l'éditeur de `buf`
  déjà utilisé pour le contrat. Une API encore jeune peut bouger : elle reste
  cantonnée au câblage du serveur, les services n'en dépendent pas.
- Un test de bout en bout par porte (gRPC et Connect) vérifie que la politique
  d'accès s'applique à chacune.
- La connexion depuis le dashboard demande un appel `Login` dans le contrat :
  il arrive avec la seconde PR de #9.
