# ADR-0007 : ForgeNet pilote Caddy par son API d'administration

- Statut : proposé
- Date : 2026-11-13
- Contexte : ticket #10 (reverse proxy Caddy), à trancher avant son code

## Contexte

Le #10 demande que `internal/proxy` génère et recharge la configuration de
Caddy, qu'un domaine s'ajoute et se retire pour une application déployée, et
qu'un écran du dashboard les gère. Caddy obtient et renouvelle seul les
certificats : c'est son HTTPS automatique.

Deux questions restent : comment ForgeNet donne sa configuration à Caddy, et
comment Caddy atteint les applications. Toutes deux touchent à la sécurité.
L'API d'administration de Caddy n'a aucune authentification : qui la joint
contrôle tout le trafic du serveur. Et chaque nom de domaine saisi dans le
dashboard finit dans un fichier de configuration.

## Options pour recharger Caddy

1. **API d'administration, `POST /load` avec un Caddyfile**
   (`Content-Type: text/caddyfile`). Caddy adapte, valide et applique la
   configuration d'un bloc. Si elle échoue, l'ancienne reste en place sans
   coupure, et l'erreur revient dans la réponse (statut 400 et message JSON).
2. **Écrire le fichier puis lancer `caddy reload`.** ForgeNet doit alors
   disposer du binaire de Caddy et du même disque que lui, alors que
   `caddy reload` passe de toute façon par l'API d'administration.
3. **Configuration JSON native par l'API.** Elle est à l'abri de l'injection
   par construction, mais illisible pour un administrateur, et sans lien avec
   l'éditeur de Caddyfile brut prévu par la matrice (§5).

## Options pour joindre les applications

- **A. Réseau Docker partagé `forgenet-edge`.** Caddy et le seul service
  exposé y sont branchés ; Caddy l'appelle par son nom de conteneur. Aucun
  port n'est publié sur l'hôte.
- **B. Caddy sur le réseau de l'hôte, vers un port publié sur `127.0.0.1`.**
  Chaque application doit publier un port que ForgeNet attribue, et un
  `127.0.0.1:` oublié expose l'application en clair sur Internet.

## Décision

Options 1 et A, validées avant le code avec une exigence : un moteur propre et
sûr.

- **Caddyfile généré en entier par ForgeNet** à partir des routes (un domaine
  vers un conteneur et un port), trié par domaine : les mêmes routes donnent
  toujours le même fichier. Une copie de la configuration appliquée est écrite
  sur disque, de façon atomique, pour l'inspection. ForgeNet ne relit jamais
  cette copie : la base de données reste la source de vérité.
- **Aucune injection possible.** Rien n'entre dans le Caddyfile sans passer
  par une liste blanche :
  - domaine : étiquettes de lettres minuscules, chiffres et tirets (63
    caractères au plus, sans tiret au bord), au moins deux étiquettes,
    253 caractères au plus, dernière étiquette non numérique ; ni joker ni
    adresse IP ;
  - cible : nom de conteneur Docker et port de 1 à 65535 ;
  - courriel ACME : adresse simple, sans espace, guillemet ni accolade.

  Un test de fuzzing vérifie qu'aucune entrée acceptée ne produit un bloc ou
  une directive de plus.
- **L'API d'administration n'est jamais exposée.** ForgeNet n'accepte qu'un
  socket Unix (`unix//chemin`, recommandé) ou une adresse de bouclage
  (`127.0.0.1`, `::1`, `localhost`). Toute autre adresse est refusée avant le
  moindre appel. Le Caddyfile généré reprend toujours cette adresse dans
  l'option `admin` : une configuration qui l'omettrait ramènerait Caddy sur
  `localhost:2019`, son adresse par défaut.
- **Un rechargement à la fois**, sous verrou, avec un délai borné et une
  réponse lue jusqu'à une taille maximale. Le client ignore les variables de
  proxy HTTP de l'environnement. Les avertissements de Caddy sont journalisés.
- **Réseau `forgenet-edge`** : seul le service qui reçoit un domaine y est
  branché, pas le reste de sa stack. Sa base de données reste sur le réseau
  privé de la stack. Le branchement, fait de façon à survivre aux
  redéploiements, arrive avec la deuxième PR du #10, celle des domaines.
- **Image de Caddy épinglée** par version et par empreinte dans les tests.
  L'installation (#15) reprendra la même règle.

## Vérification

Un test d'intégration, sauté avec `-short` ou sans Docker, lance Caddy 2.11.9
dans un conteneur jetable, son API d'administration sur un socket Unix. Il
utilise les certificats internes de Caddy (`local_certs`), car Let's Encrypt
exige un domaine public. Il vérifie que :

- la route appliquée répond en HTTPS, certificat validé par la racine de Caddy ;
- une configuration cassée est refusée et l'ancienne reste en service ;
- la route retirée ne répond plus.

## Conséquences

- Les services branchés sur `forgenet-edge` peuvent se joindre entre eux. Ce
  sont seulement les services exposés publiquement ; un réseau par
  application reste possible plus tard sans changer le moteur.
- L'obtention d'un vrai certificat Let's Encrypt ne se prouve que sur un
  serveur réel avec un domaine pointé : elle sera vérifiée au #15.
- Le socket d'administration est la seule clé de Caddy. Ses droits (mode du
  socket, dossier parent) sont fixés par l'installation (#15), pour que seul
  le compte du serveur ForgeNet y accède.
- Restent hors du #10 : le statut des certificats, les domaines jokers, les
  redirections et l'éditeur de Caddyfile brut (matrice §2.2 et §5).
