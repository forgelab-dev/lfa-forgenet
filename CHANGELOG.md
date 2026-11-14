# Changelog

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
