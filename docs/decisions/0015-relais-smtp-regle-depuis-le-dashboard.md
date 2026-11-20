# ADR-0015 : relais SMTP réglé depuis le dashboard

- Statut : accepté (choix validés par le porteur du projet le 2026-11-20)
- Date : 2026-11-20
- Contexte : ticket #102 ; complète l'ADR-0014 (envoi d'e-mails), qui laissait
  ce point pour plus tard

## Contexte

Depuis l'ADR-0014, le relais SMTP se règle en relançant `install.sh` avec ses
options `--smtp-*`, qui écrivent `.env`, puis le serveur redémarre. Changer de
relais, ou corriger son mot de passe, demande donc une connexion SSH, alors que
la page **Réglages** montre déjà le relais et sait l'essayer. Les invitations,
le « mot de passe oublié » et les avis de changement d'adresse ne partent
qu'avec un relais : c'est l'administrateur sans terminal qui en a le plus
besoin.

L'ADR-0014 avait écarté la base pour une raison : le mot de passe du relais y
serait chiffré par une clé qu'il faudrait elle-même garder, sauvegarder et
faire tourner. Cet ADR dit comment.

## Décision

### Où vit le réglage

- Dans la **base**, table `mail_settings` à une seule ligne : hôte, port,
  protection (`starttls`, `implicit`, `none` vers `localhost` seulement),
  identifiant, mot de passe chiffré, expéditeur, envoi coupé ou non, auteur
  et date du dernier changement.
- `.env`, écrit par `install.sh`, reste le **réglage de départ** : tant
  qu'aucun administrateur n'a rien enregistré depuis le dashboard, c'est lui
  qui s'applique, comme avant.

### Le mot de passe du relais

- Chiffré au repos en **AES-256-GCM** (`crypto/aes`, `crypto/cipher` de la
  bibliothèque standard, sans dépendance). Le chiffré porte un octet de
  version, puis le nonce aléatoire de 12 octets ; l'usage
  (`mail_settings.password`) est lié au chiffré comme donnée associée, pour
  qu'il ne puisse pas être déplacé vers un autre champ.
- La clé est un **secret dédié**, `FORGENET_SECRET_KEY` : 32 octets aléatoires
  écrits en hexadécimal dans `.env` (0600) par `install.sh`, comme le secret
  JWT. Un serveur déjà installé la reçoit à sa prochaine mise à jour : par
  `install.sh`, ou depuis le dashboard (ADR-0011), qui l'ajoute à `.env` s'il
  n'en a pas.
- Pas le secret JWT : on le change pour fermer toutes les sessions d'un coup,
  et ce geste ne doit pas faire perdre le mot de passe du relais.
- La base et `.env` sont **sauvegardés ensemble** avant chaque mise à jour
  (ADR-0011), et `install.sh --restore` les remet ensemble : une restauration
  garde la clé qui ouvre le mot de passe.
- Le mot de passe n'est **jamais renvoyé au navigateur**. Le formulaire dit
  seulement s'il y en a un ; on le **garde** ou on le **remplace**.
- **Clé absente ou changée** : le mot de passe enregistré ne s'ouvre plus. Le
  serveur ne s'arrête pas : il reprend les réglages de `.env`, le note au
  journal, et la page Réglages demande de ressaisir le mot de passe. Sans clé
  du tout (serveur lancé hors `install.sh`), la page montre le relais et
  l'essaie, mais ne permet pas de le changer, et dit pourquoi.

### Qui l'emporte : la dernière action explicite

- **Un enregistrement depuis le dashboard** l'emporte sur `.env`.
- **`install.sh --smtp-host …` relancé ensuite** efface les réglages de la
  base (par `forgenet-server mail reset`, une fois le serveur démarré), puis
  redémarre le serveur, qui les gardait en mémoire : celui qui relance la
  commande sur le serveur s'attend à ce qu'elle s'applique.
  Une mise à jour, qui relance `install.sh` sans ces options, n'y touche pas.
- **« Revenir aux réglages du serveur »**, sur la page, efface aussi les
  réglages de la base : `.env` s'applique de nouveau.

### Enregistrer, c'est essayer

- **« Enregistrer et essayer »** envoie un e-mail d'essai à l'administrateur
  connecté avec les **nouveaux** réglages. Ils ne sont enregistrés et
  appliqués que si le relais accepte le message ; sinon les anciens restent,
  et la réponse du relais est montrée telle quelle, comme l'essai de
  l'ADR-0014. Le même frein s'applique aux deux essais.
- **« Couper l'envoi »** est un réglage enregistré, qui l'emporte aussi sur
  `.env` : sans lui, couper depuis la page serait défait au prochain
  démarrage. Les valeurs restent dans le formulaire ; rouvrir l'envoi repasse
  par « Enregistrer et essayer ».

### Sans redémarrage

- Le serveur garde les réglages en vigueur en mémoire et les remplace d'un
  bloc à l'enregistrement. Chaque e-mail lit ceux du moment où il part : un
  message en attente d'un nouvel essai prend le nouveau relais.
- « Mot de passe oublié ? », les invitations par e-mail, la page Équipe et
  les avis de changement d'adresse demandent l'état de l'envoi à chaque fois,
  au lieu de le lire une fois au démarrage.
- Les commandes `forgenet user …` du serveur (`login-link --send`,
  `set-email`) lisent les mêmes réglages : base d'abord, puis `.env`.

### Journal d'audit

- `mail.settings_change` (enregistré, coupé, rendu au serveur) avec l'hôte,
  le port, la protection et `password` : `changé`, `gardé` ou `aucun`.
- **Jamais le mot de passe**, ni l'identifiant du relais ou l'expéditeur, qui
  sont souvent des adresses e-mail : le journal est lu par tous les
  administrateurs, et conservé.
- Un essai refusé est noté `mail.settings_change_failed`, avec l'hôte seul.

## Retour arrière

- La migration ne fait qu'**ajouter** la table `mail_settings`. Une version
  antérieure de ForgeNet, remise par `install.sh --restore` ou installée à la
  main, l'ignore et reprend `.env`.
- Pour abandonner les réglages du dashboard sans changer de version :
  « Revenir aux réglages du serveur », ou `forgenet-server mail reset` dans le
  conteneur du serveur, suivi d'un redémarrage. Le bouton est aussi proposé
  quand le mot de passe enregistré ne s'ouvre plus.

## Conséquences

- `install.sh` crée `FORGENET_SECRET_KEY` s'il manque, et le passe au serveur
  par `docker-compose.yml`. Perdre `.env` sans sauvegarde fait perdre le mot
  de passe du relais, pas le reste : il se ressaisit depuis la page.
- La rotation de la clé n'est pas outillée : la changer revient à ressaisir le
  mot de passe du relais. L'octet de version laisse la place à une rotation
  plus tard sans changer le format.
- Testé en Go (chiffrement, ordre des sources, essai refusé qui ne change
  rien, file d'envoi qui prend le nouveau relais) et au navigateur contre un
  faux relais SMTP local.
