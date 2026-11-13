# ADR-0009 : service du dashboard en production

- Statut : accepté (revue de la PR #50, 13/11/2026)
- Date : 2026-11-13
- Contexte : ticket #15 (installation en une commande), à trancher avant son
  `docker-compose.yml`

## Contexte

Le dashboard (#9) est une application Vite. En développement, le serveur Vite
le sert et relaie l'API vers ForgeNet (ADR-0006). En production, rien ne le
sert encore :

- le serveur HTTP de ForgeNet ne sert que l'API, les sondes et les flux
  WebSocket ;
- Caddy reçoit de ForgeNet un Caddyfile **complet** à chaque changement de
  domaine (ADR-0007). Un site ajouté à la main serait effacé au rechargement
  suivant.

L'ADR-0006 exige en outre que le dashboard et l'API partagent une seule
origine, sans CORS.

Reste aussi le cas de l'utilisateur qui installe ForgeNet sans nom de domaine :
il doit pouvoir ouvrir le dashboard, et la page de connexion ne doit jamais
passer en clair sur Internet.

## Options

**Qui sert les fichiers ?**

1. **Le serveur ForgeNet**, depuis un dossier (`web.root`) que l'image place à
   côté du binaire. Caddy n'a qu'un `reverse_proxy` à faire, et l'accès sans
   domaine (option C ci-dessous) obtient le dashboard sans rien de plus.
2. **Caddy**, par `file_server` sur un volume partagé : deux services à
   accorder (le volume, plus un relais des chemins de l'API vers le serveur),
   et aucun dashboard sans Caddy.
3. **Les fichiers intégrés au binaire** (`go:embed`) : `go build` dépendrait
   d'un build web préalable, ce qui casse la compilation du serveur seul en CI
   et chez les contributeurs.

**Comment Caddy connaît le dashboard ?**

1. **Un site réservé dans le Caddyfile généré** (`proxy.Settings.Dashboard`),
   réécrit à chaque rechargement comme les autres.
2. Un fichier importé par le Caddyfile : un second fichier à garder valide, que
   le rendu en liste blanche de l'ADR-0007 ne vérifie plus.

**Accès sans nom de domaine** (choix de Mike, 13/11/2026) :

- A. Domaine obligatoire : l'installation refuse de continuer. Sûr, mais cela
  bloque ceux qui veulent juste essayer.
- B. HTTPS sur l'adresse IP avec l'autorité interne de Caddy. L'avertissement
  du navigateur habitue à passer outre les erreurs de certificat.
- C. **Tunnel SSH** : sans domaine, le dashboard n'écoute que sur
  `127.0.0.1:8080` du serveur, et on l'ouvre par
  `ssh -L 8080:localhost:8080 serveur`. Rien n'est exposé tant que le HTTPS
  n'est pas en place.

## Décision

Option 1 pour le service des fichiers, option 1 pour Caddy, option C sans
domaine.

- **`web.root`** désigne le dossier du dashboard compilé. Vide (par défaut,
  comme en développement), le serveur ne sert que l'API. Le dossier est lu par
  `os.Root` : ni `..` ni un lien symbolique n'en sortent. Les fichiers cachés
  (`.env`, `.git`…) ne sont jamais servis, aucun dossier n'est listé, et seuls
  `GET` et `HEAD` sont acceptés. Une adresse de page (`/domaines`) reçoit
  `index.html`, que le routeur du dashboard interprète. Une ressource absente
  (`/assets/x.js`) reste une 404, jamais une page.
- **Le dashboard est installé sous `/`** du serveur HTTP : les routes de l'API,
  des sondes et des flux, plus précises, gardent la priorité.
- **En-têtes de sécurité sur chaque réponse du dashboard.**
  - CSP `default-src 'self'` sans aucune exception `unsafe-*` : le build Vite
    n'a ni script ni style en ligne, et les polices sont servies avec lui.
    `connect-src 'self'` couvre l'API Connect et les WebSocket de la même
    origine.
  - `frame-ancestors 'none'` et `X-Frame-Options: DENY` contre l'incrustation
    dans une autre page.
  - `nosniff`, `Referrer-Policy: no-referrer` et `Permissions-Policy`, qui
    coupe la caméra, le micro, la géolocalisation, le paiement et l'USB.
- **Cache.** Les fichiers de `assets/` portent l'empreinte de leur contenu dans
  leur nom : ils sont mis en cache un an. `index.html` est toujours revalidé,
  pour qu'une mise à jour soit vue au rechargement suivant.
- **`proxy.dashboard_host` et `proxy.dashboard_upstream`.** Avec un nom de
  domaine, Caddy reçoit un site réservé : `reverse_proxy` vers le serveur
  (`forgenet-server:8080` dans l'installation standard), en HTTPS seulement,
  avec `Strict-Transport-Security: max-age=31536000`. Sans `includeSubDomains`,
  car les sous-domaines peuvent servir des applications que ForgeNet ne gère
  pas. Les deux options vont ensemble et demandent `proxy.enabled`. Leurs
  valeurs passent par les mêmes listes blanches que les domaines des
  applications.
- **Le nom du dashboard est réservé.** Le service des domaines refuse de le
  relier à une application (`INVALID_ARGUMENT`), après normalisation, donc
  aussi en majuscules ou avec un point final. Le rendu refuse lui aussi toute
  route qui le reprend.
- **Sans domaine**, l'installation (#15) publie le port HTTP du serveur sur
  `127.0.0.1` seulement, et indique la commande du tunnel SSH.

## Conséquences

- L'image du serveur contient le dashboard compilé et fixe `web.root`.
  `docker-compose.yml` et `install.sh` (#15) en découlent.
- Le serveur sert désormais aussi des fichiers. Ce code est court, et des tests
  couvrent les tentatives de sortie du dossier : `..`, liens symboliques,
  fichiers cachés, méthodes d'écriture.
- Le HSTS s'applique au seul nom du dashboard. Revenir en HTTP sur ce nom
  demanderait d'attendre l'expiration d'un an : c'est voulu.
- L'accès par tunnel SSH demande un compte SSH sur le serveur, ce qui est déjà
  le cas de quiconque l'installe.
