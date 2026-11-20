# ADR-0014 : envoi d'e-mails et liens à usage unique

- Statut : accepté (choix validés par le porteur du projet le 2026-11-18)
- Date : 2026-11-18
- Contexte : ticket #73 ; s'appuie sur le changement de mot de passe (#72,
  v0.6.0) et complète l'ADR-0009 (installation)

## Contexte

ForgeNet n'envoie aucun e-mail. L'adresse demandée par `install.sh` est
seulement enregistrée. Trois besoins en dépendent :

- **le premier accès de l'administrateur** : installé sans terminal
  (cloud-init par exemple), ForgeNet tire aujourd'hui un mot de passe et le
  range dans `/opt/forgenet/admin-password`, un secret en clair sur le disque ;
- **« mot de passe oublié »** : seul `forgenet user set-password`, sur la
  machine, rend aujourd'hui l'accès (#72) ;
- **les invitations** d'un membre depuis la page Équipe (matrice §2.5).

Contraintes :

- un mot de passe envoyé par e-mail resterait en clair dans la boîte : il
  n'en part jamais. Tout passe par un **lien** qui fait choisir le mot de
  passe ;
- un VPS neuf envoie mal par lui-même : le port 25 est souvent bloqué et ses
  messages finissent en indésirables. Il faut passer par un **relais SMTP**
  (celui du fournisseur de messagerie, ou un service d'envoi) ;
- ForgeNet doit rester complet **sans e-mail** : beaucoup d'installations
  n'auront pas de relais ;
- aucun jeton, lien ou mot de passe ne doit finir dans un journal : ni celui
  du serveur, ni celui de Caddy, ni celui de cloud-init.

## Options

### Pour envoyer

1. **Un serveur de messagerie local** (Postfix dans l'image ou à côté) :
   écarté, pour le port 25 et la réputation d'un VPS neuf.
2. **L'API HTTP d'un service d'envoi** (Brevo, Mailgun…) : lie ForgeNet à un
   fournisseur, alors que tous ces services proposent aussi SMTP.
3. **SMTP vers un relais, par la bibliothèque standard Go** (`net/smtp`,
   `crypto/tls`) : STARTTLS sur le port 587, TLS direct sur le 465,
   authentification PLAIN seulement sous TLS. Une bibliothèque comme
   `go-mail` apporterait pièces jointes et HTML que nous n'utilisons pas :
   moins de 200 lignes à écrire nous-mêmes, sans dépendance.

### Pour garder les réglages

1. **Dans la base, modifiables depuis le dashboard** : le mot de passe SMTP
   devrait alors être chiffré par une clé qu'il faudrait elle-même garder,
   sauvegarder et faire tourner.
2. **Dans `.env`**, comme le secret JWT, écrits par `install.sh` et lus au
   démarrage ; le dashboard les montre sans le mot de passe et envoie un
   e-mail d'essai.

## Décision

Envoi : option 3. Réglages : option 2.

### Réglages SMTP

- Section `smtp` de la configuration : `host`, `port`, `username`,
  `password`, `from`, `tls` (`starttls` par défaut, `implicit`, ou `none`
  accepté seulement vers `localhost`). Sans `host`, l'envoi est éteint et
  tout le reste marche.
- `install.sh --smtp-host HÔTE [--smtp-port 587] [--smtp-user NOM]
  [--smtp-from ADRESSE]` les écrit dans `.env` (0600). Le mot de passe SMTP
  est lu sur `/dev/tty`, ou dans la variable `FORGENET_SMTP_PASSWORD` sans
  terminal, jamais en argument : il resterait dans l'historique du shell et
  dans `ps`. Relancer `install.sh` avec ces options change les réglages.
- Le dashboard (administrateurs) montre l'état de l'envoi, le relais et
  l'expéditeur, jamais le mot de passe, et envoie un **e-mail d'essai** à
  l'administrateur connecté, en rendant l'erreur du relais telle quelle.

### Les messages

- Texte brut UTF-8, en français comme le dashboard tant que la localisation
  (§2.9) n'existe pas. Pas de HTML : rien à assainir, et mieux reçu par les
  filtres.
- Envoyés en tâche de fond, trois essais espacés, sauf l'e-mail d'essai, qui
  attend la réponse du relais. Le journal note le type de message, le compte
  concerné (son identifiant, pas son adresse) et l'issue ; jamais le lien, le
  jeton, le contenu ni les identifiants du relais.

### Les liens à usage unique

- Table `action_tokens` : compte, usage (`first_login`, `password_reset`,
  `invitation`), **empreinte SHA-256** du jeton (le jeton lui-même n'est
  jamais stocké), expiration, date d'usage, auteur.
- Jeton : 32 octets aléatoires (`crypto/rand`), en base64url.
- Durées : première connexion 24 h, mot de passe oublié 1 h, invitation
  7 jours. Un nouveau lien du même usage pour le même compte annule les
  précédents.
- Le jeton voyage dans le **fragment** de l'adresse :
  `https://<dashboard>/acces#<jeton>`. Le navigateur n'envoie jamais le
  fragment au serveur, donc il n'apparaît dans aucun journal HTTP, ni
  celui de Caddy ni celui du serveur. La page le lit, puis l'envoie dans le
  corps d'un appel `AuthService` avec le nouveau mot de passe.
- L'usage est **atomique** : une seule écriture marque le jeton utilisé s'il
  ne l'est pas déjà et n'a pas expiré, puis fixe le mot de passe par
  `SetPassword`, dans la même transaction. Les sessions ouvertes sont
  fermées comme pour un changement (#72), et l'événement est inscrit au
  journal d'audit.
- L'adresse du lien vient de `proxy.dashboard_host` (`https://…`). Sans
  domaine, c'est `http://127.0.0.1:8080`, à ouvrir par le tunnel SSH
  (ADR-0009), et le message le dit.

### Premier accès de l'administrateur

- **Avec un terminal**, rien ne change : `install.sh` demande le mot de passe.
- **Sans terminal**, le compte est créé sans mot de passe utilisable, et un
  lien de première connexion est fait :
  - avec SMTP, il part à `--admin-email` ;
  - sans SMTP, il est écrit dans `/opt/forgenet/premier-acces` (0600, root),
    à la place de `admin-password`. Il n'est jamais affiché, puisque sans
    terminal la sortie part dans le journal de cloud-init.
- `forgenet user login-link --username NOM`, sur la machine, refait un lien de
  première connexion et l'affiche dans le terminal. Comme `set-password`, être
  sur la machine en est la preuve.

### Mot de passe oublié

- Un lien « Mot de passe oublié ? » sur l'écran de connexion, seulement quand
  l'envoi est configuré. Sinon, l'écran renvoie à `forgenet user
  set-password`.
- La réponse est **la même, qu'un compte corresponde ou non** (« si un compte
  correspond, un e-mail vient de partir »), et dans le même délai, puisque
  l'envoi se fait en tâche de fond. Frein : 5 demandes par heure par compte,
  et par adresse IP quand le serveur la voit (derrière Caddy, il voit celle du
  proxy, comme pour la connexion).

### Invitations

- Un administrateur invite depuis la page Équipe avec un nom d'utilisateur,
  une adresse et un rôle. Le compte est créé sans mot de passe utilisable,
  avec un lien d'invitation.
- Avec SMTP, le lien part à l'adresse. Sans SMTP, le dashboard l'affiche
  **une seule fois** à l'administrateur, qui le transmet lui-même.

## Conséquences

- Livré en quatre PR après celle-ci : réglages et envoi avec l'e-mail
  d'essai ; lien de première connexion à l'installation ; mot de passe
  oublié ; invitations.
- Les installations en place gardent leur fichier `admin-password` s'il existe
  encore ; `install.sh` ne le crée plus.
- Changer de relais demande de relancer `install.sh` avec ses options, puis
  le serveur redémarre. Le régler depuis le dashboard, avec un secret chiffré
  en base, reste possible plus tard sans changer les liens (fait par
  l'ADR-0015, #102).
- Testé contre un faux relais SMTP en Go dans les tests, puis sur un VPS avec
  un vrai relais, et sans relais.
