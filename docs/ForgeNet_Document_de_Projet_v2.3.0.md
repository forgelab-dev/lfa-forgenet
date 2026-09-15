# ForgeNet — Cloud Souverain, Local & Hybride (Open-Core)

**Version du document :** 2.3.0
**Statut :** Spécification, Modèle Économique & Feuille de Route
**Auteur :** LucidForge Africa
**Révision :** 15 septembre 2026 — révision du document v1.0.0

## Changelog v1.0.0 → v2.3.0

- Modèle économique restructuré en **4 phases**, pensé pour un lancement à budget zéro (aucun CapEx avant le premier revenu)
- Ajout d'une section **Standards de Développement** pour un suivi d'ingénierie plus propre et reproductible
- **Benchmark concurrentiel chiffré** (Coolify, Dokploy, CapRover) pour recalibrer le pricing — l'offre v1.0.0 était 3 à 10× au-dessus du marché
- Priorisation P0/P1/P2 de la matrice fonctionnelle pour un suivi de sprint plus clair
- Ajout d'indicateurs de suivi (KPIs), d'une feuille de route datée, et d'une gouvernance légère
- Clarification de la stratégie d'infrastructure : partenariat plutôt que datacenter propre (hors de portée financièrement)
- **v2.1.0** — Ajout d'un système de don/sponsoring (§5.4) en canal complémentaire, actif dès la Phase 0
- **v2.2.0** — Correction de principe : stockage 100 % self-hosted sur MinIO, plus de revente de Backblaze B2/Wasabi. Le stockage est un coût inclus dans l'abonnement, jamais une ligne facturée à part — cohérent avec la promesse "zéro-cloud tiers" du produit
- **v2.3.0** — **Bascule de Gitea vers GitHub** pour l'hébergement du dépôt, la CI/CD et le suivi de projet. Motif : la traction d'un projet open-source se joue là où la communauté se trouve déjà — un dépôt Gitea auto-hébergé est invisible pour les développeurs qui découvrent Coolify (52k étoiles) ou Dokploy. GitHub apporte la découvrabilité, les contributions externes et des Actions **gratuites et illimitées sur dépôt public** — le budget zéro est préservé. (L'argument "GitHub Sponsors en natif", avancé dans une première rédaction de cette entrée, a dû être retiré : LFA n'y est pas éligible — voir §5.3 bis. La bascule reste justifiée par les trois autres motifs.) Distinction à garder en tête : ce changement concerne l'outillage interne de LFA, **pas** le produit, qui continue de supporter Gitea comme source Git pour ses utilisateurs self-hosted
- **v2.3.0** — Élargissement des moyens de paiement : **FedaPay, Kkiapay et PayPal** (§5.3 bis, §5.3, §3). FedaPay et Kkiapay couvrent le rail FCFA/Mobile Money en redondance l'un de l'autre ; PayPal ouvre le rail international (diaspora, clients hors zone FCFA) que le Mobile Money ne peut pas servir. Point de vigilance consigné : l'encaissement PayPal est restreint dans plusieurs pays d'Afrique de l'Ouest — à confirmer avant toute communication publique

---

## 1. Résumé Exécutif & Vision Stratégique

ForgeNet est une plateforme open-source de gestion d'infrastructure cloud souveraine, conçue pour permettre aux développeurs, PME et entreprises de déployer, surveiller et orchestrer leurs services informatiques en toute autonomie. Face à la dépendance vis-à-vis des géants du cloud (AWS, GCP, Azure), ForgeNet apporte une alternative installable en local, sur VPS ou sur serveur nu.

**Positionnement révisé.** Le marché du "PaaS self-hosted" n'est pas vide : Coolify (52k+ étoiles GitHub), Dokploy et CapRover sont matures, gratuits en self-host, et déjà utilisés par des dizaines de milliers de développeurs. ForgeNet ne peut pas gagner en répliquant leurs fonctionnalités seules — il doit gagner sur un axe qu'aucun d'eux n'adresse sérieusement :

1. **Souveraineté réelle des données** en Afrique de l'Ouest (zone OHADA/CEDEAO), utile aux banques, institutions et administrations soumises à des contraintes de résidence des données.
2. **Facturation en FCFA, support en français**, sans friction de change ni décalage horaire/culturel avec un support basé en Europe/Asie.
3. **IA souveraine locale (Ollama)** : diagnostic sans fuite de données vers un cloud tiers — argument fort pour le secteur bancaire/gouvernemental.

Le reste de ce document part de ce positionnement pour recalibrer l'architecture, les standards de dev, et surtout le modèle économique.

---

## 2. Architecture Technique Globale

Architecture inchangée dans ses grandes lignes — elle est saine et cohérente avec les choix déjà faits sur l'infrastructure LFA (Vault sur `s0`, dépôt et CI sur GitHub).

| Composant | Stack | Rôle |
|---|---|---|
| **ForgeNet Core Engine** | Go, Gin/Echo, gRPC | Daemon léger : orchestration Docker, communication avec Caddy, lecture métriques (`gopsutil`), appels IA |
| **Web Management Dashboard** | React + Vite, Tailwind CSS, WebSockets | Administration complète (App Store, éditeur config, logs, reverse-proxy, BDD) |
| **Mobile Companion App** | Flutter (Android/iOS/Desktop) | Surveillance temps réel, alerting push, gestion de crise |
| **Moteur IA Souverain** | Ollama (Llama 3, Mistral, CodeLlama) | Analyse de logs, audit sécurité, assistance sysadmin zéro-cloud |
| **Reverse Proxy & SSL** | Caddy | Routage, certificats SSL automatiques |

**Note sur Caddy vs Traefik.** L'infrastructure interne LFA (`s0`, "le cerveau") tourne sur Traefik v3.3. Le choix de Caddy pour ForgeNet est néanmoins le bon : Caddy est plus simple à configurer pour un utilisateur final non-expert (auto-HTTPS sans config Let's Encrypt manuelle), ce qui correspond à la cible (PME, freelances). Traefik reste justifié en interne pour un usage avancé avec équipe SRE dédiée. Les deux choix sont cohérents avec leurs contextes respectifs — aucun changement recommandé ici.

**Principe d'infrastructure : 100 % self-hosted, zéro dépendance cloud tierce.** Décision de cadrage : ForgeNet n'externalise aucune brique de son propre fonctionnement vers un cloud tiers payant (Backblaze, Wasabi, AWS S3, etc.). Le stockage objet utilise **MinIO** self-hosted — déjà présent dans la stack Marketplace 1-Click, donc aucune nouvelle compétence à acquérir — sur disque dédié plutôt que sur du block storage cloud facturé au Go (le block storage cloud reste sensiblement plus cher au To que le stockage objet spécialisé type Backblaze ; l'argument ici n'est pas l'économie brute mais la cohérence avec la promesse de souveraineté du produit). Conséquence directe sur le pricing : le stockage est un **coût inclus** dans l'abonnement Managed Cloud, jamais une ligne facturée séparément au client — ForgeNet ne fait payer que pour son propre produit (la plateforme, le support, les missions), pas pour de l'infrastructure tierce revendue avec marge. C'est aussi un argument de vente direct face aux plateformes dont la facture "grimpe" avec des frais annexes (bande passante, add-ons, stockage) — un abonnement ForgeNet reste prévisible.

---

## 3. Matrice des Fonctionnalités : Priorisée pour le Suivi

Chaque fonctionnalité est taguée par phase de développement (voir §6 pour le calendrier), pour transformer cette matrice en backlog suivable.

### 3.1. Dashboard Web (React + Vite)

| Fonctionnalité | Priorité | Justification |
|---|---|---|
| Marketplace 1-Click (App Store) | **P0** | Cœur de la proposition de valeur — sans ça, pas de produit |
| Reverse Proxy & Domain Manager | **P0** | Indispensable dès la Community Edition |
| Gestionnaire de Fichiers & Bases de Données | **P1** | Nécessaire pour la Managed Cloud (backup/restore vers stockage géré) |
| AI SysAdmin Assistant | **P1** | Différenciateur fort, mais coûteux en dev — après le socle |
| Gestion d'Équipe & RBAC (SSO/SAML) | **P2** | RBAC basique en P0 ; SSO/SAML est un besoin Enterprise, pas MVP |

### 3.2. Application Mobile (Flutter)

| Fonctionnalité | Priorité | Justification |
|---|---|---|
| Tableau de Bord de Santé (Metrics Live) | **P0** | Différenciateur immédiat vs dashboards web-only de CapRover |
| Terminal / Quick Console | **P0** | Cas d'usage "incident à 2h du matin" — valeur perçue élevée |
| Push Notifications | **P1** | Utile mais peut suivre de quelques semaines le P0 |
| Contrôle à Distance d'Urgence | **P1** | Sensible (start/stop en prod) — mérite un P0 solide en logs/metrics d'abord |

---

## 4. Standards de Développement (nouveau)

Pour un dev "plus propre" et traçable dans le temps :

**Structure du dépôt (monorepo)** — `git@github.com:forgelab-dev/lfa-forgenet.git`

Le monorepo s'impose ici : trois langages (Go, TypeScript, Dart) consomment le **même contrat gRPC/REST**. Un changement d'API se fait en un commit sur les trois clients à la fois ; en multi-dépôt, chaque évolution devient une synchronisation de versions entre trois repos — coût insoutenable à 1-3 développeurs.

```
lfa-forgenet/
├── api/                      # ⭐ SOURCE DE VÉRITÉ — le contrat avant le code
│   ├── proto/forgenet/v1/    #   container, deploy, metrics, auth, ai, proxy
│   ├── openapi/              #   spec REST publique (§3 — API & Intégrations)
│   └── buf.gen.yaml          #   codegen → Go + TypeScript + Dart
├── cmd/                      # binaires Go
│   ├── forgenet-server/      #   control plane
│   ├── forgenet-agent/       #   agent par serveur (mode fleet)
│   └── forgenet-cli/         #   CLI de scripting/CI
├── internal/                 # cœur métier Go (privé)
│   ├── orchestrator/         #   docker, compose, buildpack, strategy, queue
│   ├── proxy/                #   génération Caddyfile
│   ├── metrics/              #   gopsutil, rétention, streaming
│   ├── ai/                   #   client Ollama, prompts, file d'attente
│   ├── auth/                 #   RBAC, tokens API, sessions
│   ├── secrets/              #   Vault + chiffrement des variables d'env
│   ├── backup/               #   MinIO, rétention, checksum
│   ├── fleet/                #   agents, control plane, placement
│   ├── resilience/           #   mode dégradé, sync différentielle (bas débit)
│   ├── audit/                #   journal horodaté des actions sensibles
│   ├── migration/            #   import Coolify / CapRover / Dokploy
│   ├── store/                #   persistance + migrations SQL
│   └── transport/            #   grpc, http, ws — aucune logique métier ici
├── pkg/                      # Go public (SDK tiers, types partagés)
├── web/                      # Dashboard React + Vite + Tailwind
│   └── src/features/         #   découpage par domaine, pas par type de fichier
├── mobile/                   # Application Flutter
│   └── lib/{features,core}/  #   core/offline = cache + file d'actions
├── templates/                # catalogue Marketplace 1-click (1 dossier = 1 app)
├── deploy/                   # install.sh, docker-compose.yml, Caddyfile.base
├── docs/                     # spécification, matrice fonctionnelle, ADRs
└── .github/workflows/        # CI (lint, test, build, release)
```

**Trois décisions structurantes.**

1. **`api/` en premier, codegen obligatoire** — aucun client HTTP écrit à la main. À trois langages, des clients manuels divergent du serveur en quelques semaines sans que personne ne s'en aperçoive avant un bug en production.
2. **SQLite par défaut, Postgres en option** — Coolify impose Postgres + Redis, ce qui alourdit son installation. Une Community Edition installable en une commande sans dépendance externe est un argument d'adoption direct. `internal/store/` derrière une interface permet de basculer sur Postgres pour le mode HA Enterprise sans réécrire le métier.
3. **`forgenet-server` et `forgenet-agent` séparés dès maintenant** — deux entrées `cmd/` partageant `internal/` coûtent quelques heures aujourd'hui ; extraire un agent d'un monolithe six mois plus tard coûte des semaines. Le mode fleet est au backlog, autant ne pas se fermer la porte.

**Contrepartie assumée du monorepo** : la CI doit filtrer par chemin (ne pas rebuilder Flutter quand seul le Go change) et le dépôt grossit — du réglage de pipeline, largement moins coûteux que la synchronisation de trois dépôts.

**Versioning & commits** — SemVer strict sur le Core Engine (API breaking changes = major) ; Conventional Commits (`feat:`, `fix:`, `chore:`) pour générer le changelog automatiquement plutôt qu'à la main.

**CI/CD** — GitHub Actions, avec des minutes **gratuites et illimitées sur dépôt public** : pipelines lint → test → build → release sans aucun coût, budget zéro préservé. Deux garde-fous à poser dès le départ : (1) un runner self-hosted sur `s0` reste branchable en complément si un job lourd (build Flutter, image Docker multi-arch) devient trop lent sur les runners GitHub ; (2) filtrage par chemin (`paths:`) dans les workflows pour ne pas rebuilder les 3 modules du monorepo à chaque commit.

**Tests** — Go : `testify` sur le Core Engine (surtout l'orchestration Docker et le parsing `gopsutil`). React : Vitest sur les composants critiques (App Store, éditeur reverse-proxy). Un unique pipeline e2e léger (Playwright) sur les deux parcours qui cassent tout si bugués : déploiement 1-click et restauration de backup.

**Secrets & config** — Réutiliser le HashiCorp Vault déjà en place sur `s0` (architecture PKI "Mairie") pour les secrets de build et de déploiement (tokens registry, clés API), plutôt que des `.env` en clair dans le repo.

**Qualité** — Hooks pre-commit (`golangci-lint`, `eslint`/`prettier`, `dart format`) ; revue de code obligatoire via Pull Request GitHub avant merge sur `main`, avec branch protection (checks CI verts obligatoires) et `CODEOWNERS` pour router automatiquement les revues par module.

---

## 5. Modèle Économique Révisé — Approche "Budget Zéro"

### 5.1 Pourquoi revoir le pricing v1.0.0

| Offre | Prix v1.0.0 | Comparable marché (2026) | Écart |
|---|---|---|---|
| Managed Cloud SaaS | 15 – 49 $/mois | Coolify Cloud : 5 $/mois · Coolify Pro : 19 $/mois · Dokploy Cloud Managed : 4,50 $/serveur | **3 à 10× le marché** |
| Enterprise/On-Premise | 500 $+/mois | (pas de comparable direct — pricing Enterprise généralement négocié, donc acceptable) | correct, à garder en "sur devis" |

Le fond du problème : à prix identique ou supérieur à des outils gratuits en self-host et déjà installés chez des dizaines de milliers de devs (Coolify a plus de 52 000 étoiles GitHub), ForgeNet n'a aucune raison d'être choisi sur la seule base des fonctionnalités. Le prix doit refléter le vrai différenciateur (souveraineté, FCFA, support FR) et non un ticket d'entrée générique de "PaaS cloud".

Autre point structurant : l'offre v1.0.0 promet un stockage S3 "Cloud LucidForge" géré en interne — c'est le bon réflexe (voir §2, principe 100 % self-hosted), à condition de le construire sur **MinIO** plutôt que d'improviser une solution propriétaire. MinIO est déjà dans la stack Marketplace de ForgeNet, S3-compatible, et évite toute dépendance à un cloud tiers payant — cohérent avec la promesse de souveraineté. Le stockage devient un coût d'exploitation inclus dans l'abonnement, jamais une ligne facturée séparément au client.

Même logique à l'échelle continentale : construire un datacenter Tier III coûte de l'ordre de 11 millions de dollars par mégawatt — hors de portée. Le mouvement observé sur le continent (Gabon avec ST Digital, Sénégal avec son cloud national, ou le partenariat ST Digital × OVHcloud pour une infrastructure "on-prem" en Côte d'Ivoire/Gabon/Cameroun avec facturation en monnaie locale) montre la voie réaliste : **s'appuyer sur un partenaire d'infrastructure existant plutôt que construire**, au moins pour les 2-3 premières années.

### 5.2 Modèle en 4 phases

**Phase 0 — Lancement Community (Mois 0-1, budget : 0 FCFA)**
- Community Edition open-source publiée (Apache 2.0), dépôt public GitHub — README soigné, démo en GIF, install en une commande : c'est la vitrine principale du projet, pas un simple espace de stockage de code
- Canaux de don/sponsoring ouverts en parallèle (voir §5.4) — coût de mise en place nul, revenu d'appoint possible dès le premier jour
- Objectif : traction, retours terrain, premières étoiles — pas de revenu significatif attendu à ce stade (les dons restent un complément, pas un pilier)

**Phase 1 — Revenu immédiat sans infrastructure (Mois 0-3, budget : 0 FCFA)**
Génère du cash sans opérer le moindre serveur pour un client :
- **Prestations "Lancement & Migration"** : installation, migration et formation facturées à la mission (forfait ~150 000 – 500 000 FCFA / ~250-800 $ par PME) — c'est le client qui paie l'infra, LFA vend l'expertise
- **Commission Marketplace** (70/30, inchangée) sur les templates premium — zéro coût d'infra, juste la logique déjà prévue dans le dashboard
- **1-2 pilotes Enterprise "fondateurs"** : un client institutionnel (banque, administration) qui paie un acompte pour un déploiement on-prem sur mesure — ce revenu finance la Phase 2

**Phase 2 — Managed Cloud Lite (Mois 3-6, budget : ~15 000-30 000 FCFA/mois, financé par la Phase 1)**
- Lancement d'un tarif recalibré et compétitif :
  - **Starter : ~5 000-8 000 FCFA/mois** (≈ 8-13 $) — se positionne juste au-dessus de Coolify Cloud (5 $), justifié par le support FR et la facturation locale
  - **Pro : ~15 000-25 000 FCFA/mois** (≈ 25-40 $) — comparable à Coolify Pro (19 $), avec backups inclus
- Infrastructure : un seul VPS mutualisé multi-tenant type Hetzner (5-20 €/mois pour démarrer), même logique que Coolify/Dokploy
- Backups : cluster **MinIO** self-hosted sur disque dédié (pas de block storage cloud facturé au Go, pas de revente de Backblaze/Wasabi) — coût inclus dans l'abonnement, jamais facturé à part

**Phase 3 — Enterprise & souveraineté réelle (Mois 6+, financé par les revenus cumulés)**
- Contrats Enterprise formalisés, facturés **annuellement et à l'avance** (le dépôt finance le déploiement, pas l'inverse)
- Exploration d'un partenariat d'hébergement souverain réel en Afrique de l'Ouest (sur le modèle ST Digital × OVHcloud) une fois le volume de clients Enterprise le justifiant — c'est à ce moment que la "souveraineté" du document v1.0.0 devient une réalité physique et pas seulement un argument marketing

### 5.3 Bis — Système de Don & Sponsoring (canal complémentaire)

À traiter comme un **complément de trésorerie, pas comme un pilier** : même pour des projets open-source établis, les dons individuels tournent typiquement autour de quelques dizaines à quelques centaines de dollars par sponsor et par an — utile pour couvrir les petits coûts récurrents (nom de domaine, VPS de démo, etc.), pas pour financer un salaire.

**❌ GitHub Sponsors — écarté, LFA n'est pas éligible**
- Les versements GitHub Sponsors transitent par **Stripe Connect**, qui ne couvre pas le pays d'immatriculation de LucidForge Africa. Le canal est donc fermé, quelle que soit la qualité du dépôt ou du profil
- **La même contrainte élimine la plupart des plateformes de don grand public** — Ko-fi, Liberapay, Buy Me a Coffee, Patreon s'appuient toutes sur Stripe ou PayPal. Inutile de les évaluer une par une : elles buteront sur le même mur
- Conséquence sur le §5.2 : la ligne "canal de don ouvert dès la Phase 0" reste valable, mais elle repose désormais **entièrement sur les passerelles locales**, pas sur un canal international clé en main
- Conséquence sur le choix GitHub (§4) : la bascule Gitea → GitHub reste justifiée par la découvrabilité, les contributions externes et la CI gratuite — mais **l'argument "GitHub Sponsors en natif" tombe** et ne doit plus être avancé

**Canal principal — FCFA via Mobile Money (FedaPay + Kkiapay)**
- **FedaPay** — passerelle principale, déjà intégrée dans LucidPay : Orange Money, MTN MoMo, Moov, carte. Coût marginal quasi nul puisque l'intégration existe déjà dans l'écosystème LFA
- **Kkiapay** — seconde passerelle, en redondance : même couverture Mobile Money/carte en FCFA, mais frais et disponibilité opérateur différents. Deux agrégateurs valent mieux qu'un : une panne de passerelle ou un opérateur défaillant ne bloque alors pas l'encaissement
- Pertinent pour capter le soutien d'un public local (devs, PME, diaspora) plus à l'aise en FCFA/Mobile Money qu'en carte internationale

**Canal international — PayPal**
- Seul canal international restant après l'exclusion de GitHub Sponsors : couvre la diaspora et les contributeurs hors zone FCFA, sans Mobile Money
- Sert aussi les achats ponctuels (mission, licence, don unique) auprès de clients étrangers
- ⚠️ **À vérifier avant annonce publique** : PayPal restreint la *réception* de fonds dans plusieurs pays d'Afrique de l'Ouest (l'envoi y est souvent ouvert, l'encaissement non). Confirmer l'éligibilité du compte marchand LFA avant d'afficher ce canal sur le site ou le README. **Si l'encaissement PayPal est lui aussi fermé, il ne reste plus aucun canal international direct** — voir Open Collective ci-dessous

**Canal de repli — Open Collective (à évaluer sérieusement)**
- Passe d'option secondaire à **alternative principale pour l'international** : Open Collective fonctionne via un **hôte fiscal** qui reçoit les fonds à la place du projet. C'est précisément ce qui permet de contourner le blocage Stripe/PayPal — l'hôte est domicilié dans un pays supporté, pas le projet
- Contrepartie : commission de 5 à 15 % des fonds entrants, et une gouvernance plus lourde (budget public, justification des dépenses)
- Bénéfice secondaire réel : la transparence budgétaire est un argument face aux sponsors institutionnels et aux prospects Enterprise qui évaluent la pérennité du projet
- **Action** : si PayPal s'avère également fermé, c'est la voie à instruire en priorité pour ne pas se priver totalement du soutien international

**Mise en œuvre** : une page de don **`https://forgenet.lucidforgeafrica.com/soutenir`** regroupant FedaPay, Kkiapay et PayPal, référencée depuis un fichier `.github/FUNDING.yml` via la clé `custom:` — ce qui affiche quand même le bouton "Sponsor" en haut de la page GitHub, sans passer par GitHub Sponsors. La clé `github:` est volontairement absente (voir ci-dessus).

Un seul lien maîtrisé plutôt que trois liens de prestataires : changer de passerelle, en ajouter une ou en retirer une ne demandera aucune modification du dépôt ni de nouvelle release. Aucun développement lourd, activable dès la Phase 0 — c'est d'ailleurs l'une des raisons de la bascule vers GitHub : ce canal y est natif et sans frais.

### 5.3 Tableau récapitulatif

| Offre | Cible | Tarif révisé | Coût infra pour LFA |
|---|---|---|---|
| Community Edition | Devs, étudiants | Gratuit | Nul |
| Dons / Sponsoring | Soutien local, diaspora | Libre — FedaPay et Kkiapay (FCFA/Mobile Money), PayPal (int'l, sous réserve d'éligibilité) | Nul |
| Prestations Lancement & Migration | PME, agences | 150k-500k FCFA / mission | Nul (client héberge) |
| Managed Cloud Lite — Starter | Freelances, PME | ~5-8k FCFA/mois | ~5-20 €/mois mutualisé |
| Managed Cloud Lite — Pro | PME, agences | ~15-25k FCFA/mois | Idem + stockage MinIO self-hosted inclus (aucune facturation séparée) |
| Enterprise / On-Premise | Banques, institutions | Sur devis, facturé à l'avance | Financé par le client |
| Marketplace | Écosystème devs | Commission 70/30 | Nul |

---

## 6. Feuille de Route & Suivi de Projet

| Phase | Période | Jalon clé | Budget requis | Revenu visé |
|---|---|---|---|---|
| 0 — Community | M0-M1 | Publication open-source, premiers retours, canaux de don ouverts | 0 FCFA | Dons d'appoint |
| 1 — Services | M0-M3 | 3-5 missions de migration signées, 1 pilote Enterprise | 0 FCFA | Autofinancement de la Phase 2 |
| 2 — Managed Cloud Lite | M3-M6 | Lancement Starter/Pro, 20-50 abonnés payants | ~15-30k FCFA/mois | MRR positif |
| 3 — Enterprise & souveraineté | M6+ | 1er contrat Enterprise annuel, étude partenariat hébergement local | Financé par revenu | Croissance MRR + contrats annuels |

**Outil de suivi.** Pas besoin d'outil externe payant : GitHub Issues + Projects (gratuits, illimités sur dépôt public) suffisent pour ce stade — un board par phase, une issue par fonctionnalité de la matrice §3, des labels par module (`core`, `web`, `mobile`, `ai`, `proxy`) et par édition (`CE`, `EE`, `MC`). Bénéfice secondaire par rapport à un tracker auto-hébergé : la roadmap est publique, donc lisible par les contributeurs potentiels et par les prospects Enterprise qui évaluent la vitalité du projet avant de s'engager.

**KPIs à suivre dès la Phase 1**
- MRR (revenu mensuel récurrent) cumulé
- Nombre d'installations Community actives (télémétrie opt-in uniquement, jamais forcée — cohérent avec la promesse "zéro-cloud")
- Taux de conversion Community → Managed Cloud
- Nombre de missions "Lancement & Migration" signées / mois
- Tickets de support (volume et délai de réponse, argument de vente du support FR)
- Dons/sponsoring cumulés par mois, ventilés par canal (FedaPay, Kkiapay, PayPal) — à suivre séparément du MRR, sans jamais le remplacer dans les projections. La ventilation par canal permet aussi d'arbitrer au bout de quelques mois : si Kkiapay ou PayPal ne capte rien, autant retirer le canal que maintenir une intégration pour rien

---

## 7. Différenciation Réelle & Risques

**Le vrai moat n'est pas la fonctionnalité.** Coolify, Dokploy et CapRover sont gratuits, matures et déjà largement adoptés. ForgeNet perdrait une bataille de fonctionnalités pures. Le moat réaliste est : souveraineté des données en zone CEDEAO/OHADA, facturation FCFA sans friction de change, support français, et IA locale sans fuite de données — un argumentaire qui parle directement aux banques et administrations, peu servi par les outils existants.

**Risque principal.** Lancer d'abord la Managed Cloud (comme en v1.0.0) expose LFA à des coûts d'infrastructure récurrents avant tout revenu prouvé. Le séquencement révisé (services d'abord, infra ensuite) élimine ce risque : chaque euro d'infra dépensé en Phase 2 est déjà couvert par du cash généré en Phase 1.

**Recommandation de positionnement.** Ne pas vendre ForgeNet comme "un panel PaaS de plus" face à des outils gratuits et installés — le vendre comme "le seul panel cloud self-hosted avec résidence des données en Afrique de l'Ouest, facturation FCFA et support français", en menant la vente avec l'argument conformité/souveraineté pour les grands comptes, et simplicité/coût pour les PME.

---

## Sources & Repères (benchmark, septembre 2026)

- Coolify — self-hosted gratuit (Apache 2.0), Cloud à partir de 5 $/mois, Cloud Pro à 19 $/mois
- Dokploy — self-hosted gratuit, Cloud Managed à 4,50 $/serveur (3,50 $ les suivants)
- Backblaze B2 — stockage S3-compatible à partir de ~7 $/To/mois — évalué en v2.0.0 pour la revente, écarté en v2.2.0 au profit d'un stockage 100 % self-hosted (MinIO)
- Hetzner Storage Box — ~6,49 €/mois pour 1 To — même arbitrage : écarté au profit de MinIO self-hosted, pour rester cohérent avec la promesse "zéro-cloud tiers"
- MinIO — stockage objet S3-compatible open-source, déjà présent dans la stack Marketplace 1-Click de ForgeNet ; le block storage cloud (nécessaire pour l'héberger) reste plus cher au To qu'un service comme Backblaze, mais élimine toute dépendance à un cloud tiers payant
- Coût de construction d'un datacenter Tier III — ~11,3 M$/MW (rapport Africa Data Centres Association, 2026)
- Partenariat ST Digital × OVHcloud (On-Prem Cloud Platform) — Côte d'Ivoire, Gabon, Cameroun, facturation en monnaie locale
- GitHub Sponsors — 0 % de frais de plateforme, mais **versements via Stripe Connect uniquement** : canal inaccessible à LFA, et par extension à la plupart des structures d'Afrique de l'Ouest. Même verrou sur Ko-fi, Liberapay, Buy Me a Coffee et Patreon
- Open Collective — fonctionnement via hôte fiscal, commission typique de 5 à 15 % des fonds entrants
- FedaPay — passerelle de paiement déjà intégrée dans LucidPay (Mobile Money Orange/MTN/Moov, cartes, FCFA)
- Kkiapay — passerelle FCFA alternative (Mobile Money, cartes), retenue en redondance de FedaPay ; grille tarifaire et couverture opérateur à comparer avant mise en production
- PayPal — canal international pour la diaspora et les clients hors zone FCFA ; éligibilité à la *réception* de fonds à confirmer pour le pays d'immatriculation de LFA
