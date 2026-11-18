# Changelog

## [0.10.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.9.0...v0.10.0) (2026-11-18)


### Fonctionnalités

* **grpc:** mot de passe oublié, GetLoginOptions et RequestPasswordReset ([892b047](https://github.com/forgelab-dev/lfa-forgenet/commit/892b047f9b645ecb586e5e03a24a04294912a03b)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **grpc:** TeamService, membres et invitations ([72a6629](https://github.com/forgelab-dev/lfa-forgenet/commit/72a6629aeed56eaebc64444cacc95ab695aa0ca7)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **mail:** e-mail d'invitation d'un membre ([3d0f7e2](https://github.com/forgelab-dev/lfa-forgenet/commit/3d0f7e26d87b6fed113cb87e9e06863007b46300))
* **mail:** file d'envoi en tâche de fond, trois essais espacés ([2b8ebef](https://github.com/forgelab-dev/lfa-forgenet/commit/2b8ebef445cad0e16541177ac691da1fdd06f443)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **web:** « Mot de passe oublié ? » sur l'écran de connexion ([ef5e1ab](https://github.com/forgelab-dev/lfa-forgenet/commit/ef5e1ab932ba13bbb2e9a8326c404dca57970357)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **web:** page Équipe, membres et invitations ([ee0e364](https://github.com/forgelab-dev/lfa-forgenet/commit/ee0e364716b33004e5d5242370e710c239bd0a81)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)


### Corrections

* **deploy:** Cancel rend la main une fois le projet libéré ([1e23f7f](https://github.com/forgelab-dev/lfa-forgenet/commit/1e23f7f4ad8b37a7bedda5ee6fd07f47ca5caaed))

## [0.9.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.8.0...v0.9.0) (2026-11-18)


### Fonctionnalités

* **auth:** comptes sans mot de passe et liens d'accès ([3452413](https://github.com/forgelab-dev/lfa-forgenet/commit/3452413956b5e30ec15de5e6133ae68c70b4bb16)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **cli:** user create --no-password et user login-link ([b214ab0](https://github.com/forgelab-dev/lfa-forgenet/commit/b214ab03c91d4e6396f3b9f964790b4878a9fb09)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **grpc:** InspectAccessLink et RedeemAccessLink, publics ([7ae2e6e](https://github.com/forgelab-dev/lfa-forgenet/commit/7ae2e6ec922202d62f40ed7aa0a841050f194e67)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **install:** lien de première connexion au lieu d'admin-password ([06f9a6f](https://github.com/forgelab-dev/lfa-forgenet/commit/06f9a6f0acdf6c5e35550f83f71d761dfbabb89d)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **store:** table action_tokens des liens à usage unique ([44feb80](https://github.com/forgelab-dev/lfa-forgenet/commit/44feb80b113abfb1eb3e8fc73317352311d2298e)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **web:** page /acces pour choisir son mot de passe par un lien ([d4e5ce9](https://github.com/forgelab-dev/lfa-forgenet/commit/d4e5ce9f052e370f5f406d6eb562c41a0334d372)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)

## [0.8.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.7.0...v0.8.0) (2026-11-18)


### Fonctionnalités

* **api:** MailService, état de l'envoi et e-mail d'essai ([25365d8](https://github.com/forgelab-dev/lfa-forgenet/commit/25365d8ec72dfb0f408d75936b2784d46d5dadf7)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **config:** section smtp pour l'envoi d'e-mails ([eb6d3fe](https://github.com/forgelab-dev/lfa-forgenet/commit/eb6d3feef0069ed1bbaea5c55f32a0b0b2f52fd4)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **grpc:** servir MailService aux administrateurs ([f6fe26e](https://github.com/forgelab-dev/lfa-forgenet/commit/f6fe26e0105059bc6fdcf324e4bf7247af6b1738)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **install:** options --smtp-* pour le relais d'e-mails ([7fa6990](https://github.com/forgelab-dev/lfa-forgenet/commit/7fa6990d051cdb9bc154f8fe4866ce31c881bede)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **mail:** envoi par relais SMTP, bibliothèque standard seule ([9dcbe59](https://github.com/forgelab-dev/lfa-forgenet/commit/9dcbe592352df43d6a184b43a0bdcfe05e317628)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)
* **web:** page Réglages, état de l'envoi et e-mail d'essai ([9010148](https://github.com/forgelab-dev/lfa-forgenet/commit/901014812ec4f62e19bcd6c5ca4599ee6ba975fc)), closes [#73](https://github.com/forgelab-dev/lfa-forgenet/issues/73)

## [0.7.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.6.0...v0.7.0) (2026-11-18)


### Fonctionnalités

* **api:** Container.system pour les conteneurs de ForgeNet ([f7042f4](https://github.com/forgelab-dev/lfa-forgenet/commit/f7042f43592cf46c864a07220b636fb350b2a53f)), closes [#74](https://github.com/forgelab-dev/lfa-forgenet/issues/74)
* **deploy:** étiqueter le serveur et Caddy comme conteneurs de ForgeNet ([9456ba0](https://github.com/forgelab-dev/lfa-forgenet/commit/9456ba024249ef5beabdd4838a0e900d7c692802)), closes [#74](https://github.com/forgelab-dev/lfa-forgenet/issues/74)
* **orchestrator:** reconnaître l'étiquette forgenet.system ([69a10af](https://github.com/forgelab-dev/lfa-forgenet/commit/69a10af510d993d959a0bfbee5d0dd6c4330e0f6)), closes [#74](https://github.com/forgelab-dev/lfa-forgenet/issues/74)


### Corrections

* **web:** « Vérifier les mises à jour » sur sa propre ligne ([cc77021](https://github.com/forgelab-dev/lfa-forgenet/commit/cc7702193767d81d1aaec562c5ece2bb35562af8)), closes [#74](https://github.com/forgelab-dev/lfa-forgenet/issues/74)
* **web:** compter les conteneurs de ForgeNet à part des applications ([20160cc](https://github.com/forgelab-dev/lfa-forgenet/commit/20160cc49aad58ee5e66714a2cb27ff825914e14)), closes [#74](https://github.com/forgelab-dev/lfa-forgenet/issues/74)
* **web:** relire l'état des mises à jour à la fin du suivi ([72a9e99](https://github.com/forgelab-dev/lfa-forgenet/commit/72a9e99e415212c39e7df08afdc6b850d73696c5)), closes [#74](https://github.com/forgelab-dev/lfa-forgenet/issues/74)

## [0.6.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.5.0...v0.6.0) (2026-11-18)


### Fonctionnalités

* **api:** ChangePassword dans AuthService ([8279c7b](https://github.com/forgelab-dev/lfa-forgenet/commit/8279c7b737175e3eb040a148a0a86524923b664e)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)
* **auth:** changer ou réinitialiser un mot de passe ([9b261e9](https://github.com/forgelab-dev/lfa-forgenet/commit/9b261e9b2381a0ccc6912a693caec2737a8b654a)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)
* **cli:** forgenet-server user set-password ([8300a76](https://github.com/forgelab-dev/lfa-forgenet/commit/8300a764b951fca3274b4aa45891174aa75e1e86)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)
* **grpc:** servir ChangePassword aux membres connectés ([b53b4cc](https://github.com/forgelab-dev/lfa-forgenet/commit/b53b4cc7be07bec6801f4787c7602456a68958a8)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)
* **store:** SetPassword change l'empreinte et ferme les sessions ([5acb2eb](https://github.com/forgelab-dev/lfa-forgenet/commit/5acb2eb70c4e81bf07363058ea1247d7954d1665)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)
* **web:** page « Mon compte » et changement de mot de passe ([fe88c27](https://github.com/forgelab-dev/lfa-forgenet/commit/fe88c2791db7a594c713c7921933827607f05233)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)


### Corrections

* **install:** ne plus faire supprimer le mot de passe avant de l'avoir changé ([1123c7c](https://github.com/forgelab-dev/lfa-forgenet/commit/1123c7c2bf56aabcb676132a61e1babf11fe4012)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)
* **web:** case et commande de secours lisibles sur la page du compte ([8d80a00](https://github.com/forgelab-dev/lfa-forgenet/commit/8d80a00baa62838504d4ee08d5bb4aa809fd2514)), closes [#72](https://github.com/forgelab-dev/lfa-forgenet/issues/72)

## [0.5.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.4.1...v0.5.0) (2026-11-16)


### Fonctionnalités

* **deploy:** commande forgenet d'administration sur le serveur ([256f014](https://github.com/forgelab-dev/lfa-forgenet/commit/256f014d08fb8c828c11285aa92430fa0fdbe034)), closes [#71](https://github.com/forgelab-dev/lfa-forgenet/issues/71)


### Corrections

* **version:** afficher la vraie date de compilation ([0ff83e0](https://github.com/forgelab-dev/lfa-forgenet/commit/0ff83e09c2063e2cf0b818983c517d0e2c45c859)), closes [#71](https://github.com/forgelab-dev/lfa-forgenet/issues/71)

## [0.4.1](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.4.0...v0.4.1) (2026-11-16)


### Corrections

* **compose:** relancer les applications après un redémarrage du serveur ([50f19f8](https://github.com/forgelab-dev/lfa-forgenet/commit/50f19f86d25a24ab0744e9d16bbd3edb0f3446a8)), closes [#70](https://github.com/forgelab-dev/lfa-forgenet/issues/70)

## [0.4.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.3.0...v0.4.0) (2026-11-14)


### Fonctionnalités

* **web:** logo ForgeNet, l'enclume ([f23e6fd](https://github.com/forgelab-dev/lfa-forgenet/commit/f23e6fd33dd33a6fd875fa2da19bfa2d98a887c5))

## [0.3.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.2.0...v0.3.0) (2026-11-14)


### Fonctionnalités

* **deploy:** sauvegarder les 0.1.x avant réinstallation, et --restore ([2da6ee4](https://github.com/forgelab-dev/lfa-forgenet/commit/2da6ee4bbef59a0bed34e619de4ee550dbd74eeb))

## [0.2.0](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.1.1...v0.2.0) (2026-11-14)


### Fonctionnalités

* **api:** contrat UpdateService ([fd10501](https://github.com/forgelab-dev/lfa-forgenet/commit/fd105015d594ce3604dc4e80a9d0c9b24f6cc3f7)), closes [#62](https://github.com/forgelab-dev/lfa-forgenet/issues/62)
* **deploy:** install.sh --update et image épinglée par empreinte ([271a077](https://github.com/forgelab-dev/lfa-forgenet/commit/271a077a52e500bbaa69bf7f37c082db9b7e04ae)), closes [#62](https://github.com/forgelab-dev/lfa-forgenet/issues/62)
* **update:** annonce des versions et mise à jour par l'API ([5063a9b](https://github.com/forgelab-dev/lfa-forgenet/commit/5063a9b5a4f53528ef1a78b02a9092337c3764a7)), closes [#62](https://github.com/forgelab-dev/lfa-forgenet/issues/62)
* **update:** moteur de mise à jour et sous-commandes ([f1533ea](https://github.com/forgelab-dev/lfa-forgenet/commit/f1533ea1391fa31a885d7ad7e36f146d287a7ac8)), closes [#62](https://github.com/forgelab-dev/lfa-forgenet/issues/62)
* **web:** bandeau de mise à jour du dashboard ([e2067b2](https://github.com/forgelab-dev/lfa-forgenet/commit/e2067b2f409176ff2925a7945f47bff902465005)), closes [#62](https://github.com/forgelab-dev/lfa-forgenet/issues/62)


### Corrections

* **update:** reprendre un verrou abandonné ([be33730](https://github.com/forgelab-dev/lfa-forgenet/commit/be337307a82241cfe948b3ec52947f833f66d772)), closes [#62](https://github.com/forgelab-dev/lfa-forgenet/issues/62)

## [0.1.1](https://github.com/forgelab-dev/lfa-forgenet/compare/v0.1.0...v0.1.1) (2026-11-14)


### Corrections

* **web:** un nom d'image long n'élargit plus le tableau ([b5a948b](https://github.com/forgelab-dev/lfa-forgenet/commit/b5a948b6bdb5330089053f92f6fd3d7fde4793b5))

## 0.1.0 (2026-11-14)

Première version publiée de ForgeNet, en **alpha** : à essayer sur un serveur de test, pas encore en production. L'installation, ce que fait cette version et ce qui manque sont décrits dans le [README](https://github.com/forgelab-dev/lfa-forgenet/blob/v0.1.0/README.md) ; l'état détaillé dans la [matrice fonctionnelle](https://github.com/forgelab-dev/lfa-forgenet/blob/v0.1.0/docs/ForgeNet_Matrice_Fonctionnelle_Complete_Dev.md).

### Fonctionnalités

* **installation :** `deploy/install.sh` pour Debian et Ubuntu, serveur et Caddy en conteneurs durcis ; image `ghcr.io/forgelab-dev/forgenet-server` pour amd64 et arm64, signée par cosign, avec provenance et SBOM ([#15](https://github.com/forgelab-dev/lfa-forgenet/issues/15), [#16](https://github.com/forgelab-dev/lfa-forgenet/issues/16))
* **déploiements :** stacks Docker Compose déployées, mises à jour et supprimées depuis le dashboard, mises en file et reprises au redémarrage du serveur ; suivi en direct avec leur journal, état réel de chaque stack ([#7](https://github.com/forgelab-dev/lfa-forgenet/issues/7), [#48](https://github.com/forgelab-dev/lfa-forgenet/issues/48))
* **domaines et HTTPS :** un domaine par service, Caddyfile généré et rechargé par l'API d'administration de Caddy, certificats Let's Encrypt ([#10](https://github.com/forgelab-dev/lfa-forgenet/issues/10))
* **surveillance :** processeur, mémoire, disque et réseau de l'hôte et de chaque conteneur, en direct par WebSocket compressé, avec la dernière heure gardée en mémoire ; reconnexion propre sur une connexion lente ou coupée ([#11](https://github.com/forgelab-dev/lfa-forgenet/issues/11), [#49](https://github.com/forgelab-dev/lfa-forgenet/issues/49))
* **conteneurs :** lister, consulter, démarrer, arrêter, redémarrer, lire les journaux ([#6](https://github.com/forgelab-dev/lfa-forgenet/issues/6))
* **comptes :** connexion par identifiant et mot de passe, rôles administrateur et membre vérifiés à chaque appel, tentatives limitées, journal d'audit ([#8](https://github.com/forgelab-dev/lfa-forgenet/issues/8), [#9](https://github.com/forgelab-dev/lfa-forgenet/issues/9))
* **dashboard :** connexion, vue d'ensemble, conteneurs, déploiements, domaines et surveillance, servi par le serveur lui-même à la même origine que l'API ([#9](https://github.com/forgelab-dev/lfa-forgenet/issues/9), [#15](https://github.com/forgelab-dev/lfa-forgenet/issues/15))
* **api :** contrat gRPC v1 (sessions, conteneurs, déploiements, domaines, métriques), servi aussi en Connect et gRPC-Web, avec des clients générés en TypeScript et en Dart
* **serveur :** configuration, transports HTTP et gRPC, sondes de santé, stockage SQLite avec migrations versionnées ([#4](https://github.com/forgelab-dev/lfa-forgenet/issues/4), [#5](https://github.com/forgelab-dev/lfa-forgenet/issues/5))

### Corrections

* **sécurité :** golang.org/x/net 0.62.0 (GO-2026-6866, GO-2026-6860) ; minimatch 11.0.8 dans l'outillage du dashboard (GHSA-g3hv-75wg-m249)
* **auth :** `/readyz` ne détaille l'état des dépendances qu'aux appelants autorisés ; un appel abandonné par le client n'est plus compté comme une panne
* **serveur :** flux gRPC protégés contre les paniques ; horodatages SQLite fiabilisés
* **dashboard :** l'heure des dernières mesures ne s'affiche qu'après une coupure
