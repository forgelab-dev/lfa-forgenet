# ForgeNet — Matrice Fonctionnelle Complète (Backlog Dev)

**Base :** [ForgeNet_Document_de_Projet_v2.3.0.md](ForgeNet_Document_de_Projet_v2.3.0.md) — remplace/étend le §3 (Matrice des Fonctionnalités)
**Portée :** architecture à 5 modules (§2) — Core Engine, Dashboard Web, Mobile, IA, Reverse Proxy — + Écosystème
**Statut :** backlog exhaustif, sans priorisation de sprint (voir le document projet §3 pour le P0/P1/P2)
**Révision :** 15 septembre 2026

---

## Légende — Modèle Open-Core

Chaque fonctionnalité est taguée par édition, pour que la matrice serve aussi de grille de pricing produit et reste alignée avec le modèle économique du §5 du document projet.

| Tag | Édition | Description |
|---|---|---|
| 🔓 **CE** | Community Edition | Gratuit, open-source (licence permissive), self-hosted, **sans limite artificielle de conteneurs/serveurs** |
| 🏢 **EE** | Enterprise (add-on) | Fonctions avancées multi-équipe/compliance, payantes, activables sur une instance self-hosted |
| ☁️ **MC** | Managed Cloud | Hébergé par LucidForge Africa, facturé en FCFA (Mobile Money/carte) |

Le cœur du produit — déploiement, monitoring, reverse proxy, IA locale — reste **CE** par défaut : c'est la condition de crédibilité face à Coolify/Dokploy/CapRover, tous gratuits en self-host. EE et MC ne ferment jamais une fonctionnalité déjà en CE, ils ajoutent de l'échelle, du support et de la commodité.

---

## 1. Core Engine (Go, Gin/Echo, gRPC)

### 1.1 Orchestration Docker
- [ ] CRUD conteneurs (créer, démarrer, arrêter, redémarrer, supprimer, kill forcé) `🔓 CE`
- [ ] Support Docker Compose (parsing, déploiement multi-service, dépendances entre services) `🔓 CE`
- [ ] Build d'image depuis Dockerfile `🔓 CE`
- [ ] Détection automatique de buildpacks (Node, Python, Go, PHP, Ruby, statique) façon Nixpacks `🔓 CE`
- [ ] Déploiement depuis Git (push-to-deploy via webhook GitHub/GitLab/Gitea) `🔓 CE`
- [ ] **Déploiements Blue-Green et Canary** avec bascule de trafic progressive et rollback automatique sur échec de health check `🔓 CE`
- [ ] Rollback vers une version précédente d'un déploiement (1 clic) `🔓 CE`
- [ ] Gestion des volumes (création, montage, snapshot) `🔓 CE`
- [ ] Gestion des réseaux Docker (isolation par projet/équipe) `🔓 CE`
- [ ] Health checks configurables (HTTP, TCP, commande custom) avec auto-restart `🔓 CE`
- [ ] Limites de ressources par conteneur (CPU, RAM, disque) et enforcement `🔓 CE`
- [ ] File d'attente de déploiement (anti-conflit sur un même service) `🔓 CE`
- [ ] Scheduler de tâches cron définies par l'utilisateur `🔓 CE`
- [ ] **Boucle GitOps** : réconciliation continue entre l'état déclaré en dépôt (manifeste ForgeNet) et l'état réel des conteneurs, avec dérive détectée et corrigée `🏢 EE`

### 1.2 Métriques & Observabilité
- [ ] Collecte via `gopsutil` (CPU, RAM, disque, réseau, uptime) — hôte + par conteneur `🔓 CE`
- [ ] Historisation des métriques avec rétention configurable (stockage local) `🔓 CE`
- [ ] Endpoint streaming WebSocket pour dashboard/mobile temps réel `🔓 CE`
- [ ] Export Prometheus / OpenTelemetry standard (interopérable avec Grafana existant) `🔓 CE`
- [ ] **Tableau de prévision de capacité** : projection d'usage CPU/RAM/disque à 30/60/90 jours, alerte avant saturation `🏢 EE`

### 1.3 API & Intégrations
- [ ] API gRPC interne (dashboard ↔ core, mobile ↔ core via gateway REST/gRPC-Web) `🔓 CE`
- [ ] API REST publique documentée (OpenAPI 3) pour intégrations tierces `🔓 CE`
- [ ] Authentification par token API (scoped, révocable, expirable) `🔓 CE`
- [ ] Rate limiting sur l'API publique `🔓 CE`
- [ ] Webhooks sortants (Slack, Discord, Telegram, email, SMS) `🔓 CE`
- [ ] CLI `forgenet-cli` (login, deploy, logs, status, scale) pour scripting/CI externe `🔓 CE`
- [ ] **Provider Terraform / Pulumi officiel** — gérer ForgeNet en Infrastructure-as-Code `🏢 EE`
- [ ] **Assistant de migration** : import automatique de projets depuis Coolify, CapRover, Dokploy et Portainer (lecture de leurs configs/Compose existants) `🔓 CE`

### 1.4 Sécurité & Secrets
- [ ] Intégration HashiCorp Vault (secrets de build/déploiement : tokens registry, clés API) `🔓 CE`
- [ ] Chiffrement des variables d'environnement au repos `🔓 CE`
- [ ] RBAC enforcement côté moteur (vérification systématique par appel, pas seulement côté UI) `🔓 CE`
- [ ] Audit log horodaté de toutes les actions sensibles (start/stop/deploy/delete), attribué à un utilisateur `🔓 CE`
- [ ] Scan de vulnérabilités des images (intégration Trivy, open-source) `🔓 CE`
- [ ] Détection de ports exposés à risque + alerte `🔓 CE`
- [ ] **Rotation automatique des secrets** (tokens registry, clés DB) avec fenêtre de grâce sans downtime `🏢 EE`
- [ ] Politiques réseau par projet (firewall applicatif, allow-list d'IP sortantes/entrantes) `🏢 EE`
- [ ] **Export de conformité** (journal d'audit signé, format exploitable pour due diligence OHADA/bancaire) `🏢 EE`

### 1.5 Backup & Restauration
- [ ] Backup planifié (bases de données + volumes) vers MinIO self-hosted `🔓 CE`
- [ ] Restauration en un clic depuis une sauvegarde donnée `🔓 CE`
- [ ] Rétention configurable (nombre de copies, durée) `🔓 CE`
- [ ] Vérification d'intégrité des backups (checksum) `🔓 CE`
- [ ] Réplication géographique des backups entre deux nœuds MinIO (DR inter-site) `☁️ MC`

### 1.6 Multi-serveur / Fleet
- [ ] Mode agent : un Core Engine par serveur, orchestré depuis un control plane `🔓 CE`
- [ ] Vue fédérée multi-serveurs `🔓 CE`
- [ ] Placement de service sur le serveur le moins chargé (bin-packing basique) `🔓 CE`
- [ ] Haute disponibilité du control plane (mode cluster, pas de SPOF) `🏢 EE`

### 1.7 Résilience Réseau & Contexte Bas-Débit
> Différenciateur direct pour la cible ouest-africaine — aucun concurrent (Coolify/Dokploy/CapRover) ne conçoit pour la connectivité intermittente.
- [ ] **Mode dégradé** : file d'attente locale des commandes/déploiements en cas de coupure, rejeu automatique au retour du réseau `🔓 CE`
- [ ] Compression et agrégation des flux de logs/metrics avant transmission (économie de bande passante) `🔓 CE`
- [ ] Synchronisation différentielle (delta uniquement) entre agents et control plane `🔓 CE`
- [ ] Mode "faible bande passante" pour le dashboard (désactive le streaming temps réel, polling espacé) `🔓 CE`

### 1.8 Intégration IA & Télémétrie
- [ ] Endpoint interne exposant logs/metrics formatés au moteur Ollama `🔓 CE`
- [ ] File d'attente des requêtes IA (éviter de saturer la machine hôte) `🔓 CE`
- [ ] Télémétrie opt-in uniquement, anonymisée, désactivable à tout moment `🔓 CE`

---

## 2. Dashboard Web (React + Vite + Tailwind)

### 2.1 Marketplace / App Store
- [ ] Catalogue de templates 1-click (BDD, CMS, outils devops, etc.) `🔓 CE`
- [ ] Soumission de template par la communauté `🔓 CE`
- [ ] Notation/avis sur les templates `🔓 CE`
- [ ] **SDK d'extension** (plugins tiers : actions custom, nouvelles intégrations, hooks de déploiement) `🔓 CE`
- [ ] Page vendeur avec suivi des revenus (commission 70/30 sur templates premium) `☁️ MC`

### 2.2 Reverse Proxy & Domain Manager
- [ ] Ajout/suppression de domaines et sous-domaines par service `🔓 CE`
- [ ] Statut visible des certificats SSL (génération auto via Caddy) `🔓 CE`
- [ ] Éditeur de règles de redirection/réécriture `🔓 CE`
- [ ] Support domaines wildcard `🔓 CE`
- [ ] Load balancing visuel entre instances d'un même service `🔓 CE`
- [ ] **Environnements de preview éphémères** : URL unique auto-générée et auto-expirée par pull request `🔓 CE`

### 2.3 Gestionnaire de Fichiers & Bases de Données
- [ ] Explorateur de fichiers par conteneur/volume (upload/download/édition) `🔓 CE`
- [ ] Provisioning BDD en un clic (Postgres, MySQL, MongoDB, Redis) `🔓 CE`
- [ ] Client SQL/NoSQL intégré (requêtes, visualisation de tables) `🔓 CE`
- [ ] Backup/restore manuel ou planifié par base `🔓 CE`

### 2.4 AI SysAdmin Assistant (interface)
- [ ] Chat intégré (langage naturel, FR/EN) `🔓 CE`
- [ ] Résumé automatique des logs d'erreur d'un déploiement en échec `🔓 CE`
- [ ] Suggestions d'optimisation de ressources `🔓 CE`
- [ ] Audit de sécurité à la demande `🔓 CE`
- [ ] Historique des conversations par projet `🔓 CE`
- [ ] **Auto-remédiation suggérée en un clic** (ex : redémarrer un conteneur OOM-killed, augmenter une limite mémoire, régénérer un certificat expiré) avec action réversible et journalisée `🏢 EE`

### 2.5 Gestion d'Équipe & RBAC
- [ ] Rôles prédéfinis (Owner, Admin, Développeur, Lecture seule) `🔓 CE`
- [ ] RBAC granulaire par projet/service `🔓 CE`
- [ ] Invitations par email `🔓 CE`
- [ ] Journal d'activité par membre `🔓 CE`
- [ ] SSO/SAML `🏢 EE`
- [ ] **Mode multi-tenant / marque blanche** : un intégrateur/MSP héberge ForgeNet et revend des instances à ses propres clients sous sa marque `🏢 EE`

### 2.6 Monitoring & Logs
- [ ] Logs temps réel (stream WebSocket, filtrage par service/niveau) `🔓 CE`
- [ ] Graphiques de métriques historiques (CPU/RAM/réseau/disque) `🔓 CE`
- [ ] Alertes configurables (seuils, canaux de notification) `🔓 CE`
- [ ] Timeline des déploiements `🔓 CE`
- [ ] **Tableau de coût & de ressources** : estimation du coût d'hébergement par projet/service, comparatif vs plans cloud tiers `☁️ MC`

### 2.7 Déploiement & CI/CD
- [ ] Connexion dépôt Git (GitHub, GitLab, Gitea) — GitHub en priorité d'implémentation `🔓 CE`
- [ ] Déploiement automatique sur push (branche configurable) `🔓 CE`
- [ ] Environnements de prévisualisation avant merge `🔓 CE`
- [ ] Historique des déploiements avec rollback en un clic `🔓 CE`
- [ ] **Assistant d'import** : détection et migration guidée d'un projet Coolify/CapRover/Dokploy existant (parcours "quitter mon PaaS actuel" en moins de 10 minutes) `🔓 CE`

### 2.8 Facturation & Abonnement (Managed Cloud)
- [ ] Intégration **FedaPay** (Mobile Money Orange/MTN/Moov, carte, FCFA) — passerelle principale `☁️ MC`
- [ ] Intégration **Kkiapay** (Mobile Money, carte, FCFA) — passerelle de redondance, bascule automatique si FedaPay est indisponible `☁️ MC`
- [ ] Intégration **PayPal** — rail international pour la diaspora et les clients hors zone FCFA `☁️ MC`
- [ ] Couche d'abstraction paiement (un seul port `PaymentProvider`, trois adaptateurs) — ne jamais coder un prestataire en dur dans la logique de facturation `☁️ MC`
- [ ] Gestion d'abonnement (upgrade/downgrade Starter/Pro) `☁️ MC`
- [ ] Factures téléchargeables `☁️ MC`
- [ ] Page "Soutenir ForgeNet" sur `forgenet.lucidforgeafrica.com/soutenir` (GitHub Sponsors, FedaPay, Kkiapay, PayPal) — page et profil Sponsors référencés depuis `.github/FUNDING.yml` `🔓 CE`

### 2.9 Paramètres & Compte
- [ ] Profil, tokens API, préférences de notification `🔓 CE`
- [ ] Thème clair/sombre `🔓 CE`
- [ ] Localisation FR/EN `🔓 CE`
- [ ] Personnalisation de marque (logo, couleurs, domaine du dashboard) pour usage marque blanche `🏢 EE`

---

## 3. Application Mobile (Flutter)

### 3.1 Tableau de Bord Santé
- [ ] Vue temps réel CPU/RAM/disque/réseau par serveur/service `🔓 CE`
- [ ] Sélecteur multi-serveurs `🔓 CE`
- [ ] Widgets personnalisables (favoris) `🔓 CE`

### 3.2 Terminal / Quick Console
- [ ] Accès terminal type SSH depuis le mobile `🔓 CE`
- [ ] Historique de commandes `🔓 CE`
- [ ] Raccourcis de commandes fréquentes `🔓 CE`

### 3.3 Notifications Push
- [ ] Alertes serveur down / haute utilisation / échec de déploiement `🔓 CE`
- [ ] Notification de succès de déploiement `🔓 CE`
- [ ] Configuration granulaire des types d'alerte `🔓 CE`
- [ ] **Résumé vocal/texte généré par l'IA locale** de l'incident en cours, envoyé en notification enrichie ("quoi, depuis quand, cause probable") `🏢 EE`

### 3.4 Contrôle à Distance d'Urgence
- [ ] Start/stop/restart d'un service depuis le mobile `🔓 CE`
- [ ] Confirmation obligatoire (biométrie/double-tap) avant action destructive en prod `🔓 CE`
- [ ] Journal des actions effectuées depuis le mobile `🔓 CE`

### 3.5 Authentification & Sessions
- [ ] Connexion biométrique / PIN `🔓 CE`
- [ ] Session multi-serveurs avec bascule rapide `🔓 CE`
- [ ] Déconnexion à distance des sessions actives `🔓 CE`

### 3.6 Mode Hors-ligne
- [ ] Cache du dernier état connu (metrics, statut services) en cas de perte de connexion `🔓 CE`
- [ ] File d'actions en attente (start/stop demandés hors-ligne, exécutés au retour du réseau, avec confirmation) `🔓 CE`

---

## 4. Moteur IA Souverain (Ollama)

- [ ] Sélection de modèle selon la tâche (Llama 3, Mistral, CodeLlama) `🔓 CE`
- [ ] Analyse de logs (détection d'anomalies, résumé d'erreurs) `🔓 CE`
- [ ] Audit sécurité automatisé (configs exposées, secrets en clair, ports ouverts) `🔓 CE`
- [ ] Assistant conversationnel sysadmin (dashboard + mobile) `🔓 CE`
- [ ] Suggestions de dimensionnement (right-sizing CPU/RAM) `🔓 CE`
- [ ] **Traitement 100 % local — zéro appel à une API IA tierce, zéro fuite de données hors du serveur du client** `🔓 CE`
- [ ] File d'attente / gestion de charge des requêtes (éviter de saturer la machine hôte) `🔓 CE`
- [ ] **Détection proactive d'anomalies** (dérive de comportement avant l'incident, pas seulement après) par apprentissage de la baseline normale par service `🏢 EE`
- [ ] Mode "incognito" strict : aucune télémétrie, aucun log de prompt conservé au-delà de la session, activable en un clic pour secteurs réglementés `🔓 CE`

---

## 5. Reverse Proxy & SSL (Caddy)

- [ ] Auto-HTTPS (Let's Encrypt) sans configuration manuelle `🔓 CE`
- [ ] Génération automatique de Caddyfile à partir des règles définies dans le dashboard `🔓 CE`
- [ ] Éditeur de configuration avancée (Caddyfile brut) pour utilisateurs experts `🔓 CE`
- [ ] Support des domaines wildcard `🔓 CE`
- [ ] Load balancing round-robin / least-connections entre instances `🔓 CE`
- [ ] Redirection HTTP→HTTPS forcée par défaut `🔓 CE`
- [ ] **Pare-feu applicatif (WAF) intégré** : règles anti-injection/anti-scan de base, activables par service sans dépendance tierce `🏢 EE`

---

## 6. Fonctionnalités Différenciantes — Vue d'Ensemble

Ces 8 points sont la synthèse "pitch" de la matrice : ce qu'aucun des trois concurrents directs (Coolify, Dokploy, CapRover) n'offre aujourd'hui, et qui justifie l'existence de ForgeNet plutôt qu'un fork de l'un d'eux.

| # | Fonctionnalité | Pourquoi c'est fort |
|---|---|---|
| 1 | **IA souveraine 100 % locale (Ollama)** | Diagnostic et audit sans qu'une seule ligne de log ne quitte le serveur — argument décisif pour banques/administrations |
| 2 | **Résilience bas-débit / mode dégradé** | Conçu pour la connectivité ouest-africaine réelle, pas seulement pour un datacenter européen bien connecté |
| 3 | **Assistant de migration depuis Coolify/CapRover/Dokploy** | Réduit à quelques minutes le coût de bascule depuis un concurrent — lève le principal frein à l'adoption |
| 4 | **Facturation FCFA native (FedaPay + Kkiapay), PayPal pour l'international** | Zéro friction de change ni carte bancaire internationale requise en local, et deux passerelles FCFA plutôt qu'une — une panne d'agrégateur ne bloque pas l'encaissement |
| 5 | **Mode multi-tenant / marque blanche** | Ouvre un canal de revenu B2B2X : intégrateurs et MSP locaux revendent ForgeNet sous leur marque |
| 6 | **Déploiements Blue-Green/Canary + auto-remédiation IA** | Fiabilité de niveau plateforme cloud majeure, dans un produit self-hosted gratuit en CE |
| 7 | **Stockage 100 % self-hosted (MinIO), zéro cloud tiers** | Coût prévisible, aucune dépendance à un fournisseur étranger — cohérent de bout en bout avec la promesse de souveraineté |
| 8 | **SDK de plugins + marketplace communautaire rémunérée** | Écosystème extensible dès la Community Edition, avec un vrai modèle de revenu pour les contributeurs |

---

## Notes d'usage de cette matrice

- Les cases à cocher servent de suivi de backlog brut ; le document projet (§3, §6) reste la source de vérité pour la priorisation (P0/P1/P2) et le calendrier de sprint.
- Le tag d'édition (`CE`/`EE`/`MC`) est une proposition de cadrage produit — à valider en comité avant tout engagement commercial ou communication publique de pricing.
- Toute fonctionnalité marquée `EE` ou `MC` doit avoir un chemin d'usage dégradé mais fonctionnel en `CE` : ForgeNet ne doit jamais rendre un service self-hosté inutilisable pour pousser vers le payant.
