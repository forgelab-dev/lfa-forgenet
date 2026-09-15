<div align="center">

# ForgeNet

**Le cloud souverain, self-hosted, qui ne quitte jamais vos serveurs.**

Déployez, surveillez et orchestrez vos services — avec une IA d'administration qui tourne *en local*, sans qu'une seule ligne de log ne parte vers un cloud tiers.

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![React](https://img.shields.io/badge/React-Vite-61DAFB?logo=react&logoColor=white)](https://react.dev)
[![Flutter](https://img.shields.io/badge/Flutter-Mobile-02569B?logo=flutter&logoColor=white)](https://flutter.dev)
[![Status](https://img.shields.io/badge/statut-pré--alpha-orange.svg)](#-état-du-projet)

**Français** · [English](README.en.md)

[Fonctionnalités](#-fonctionnalités) · [Installation](#-installation) · [Comparatif](#-pourquoi-pas-coolify-ou-dokploy-) · [Architecture](#-architecture) · [Contribuer](#-contribuer) · [Soutenir](#-soutenir-le-projet)

</div>

---

## 🚧 État du projet

> **Pré-alpha — en cours de construction.** Le code n'est pas encore utilisable en production et les commandes d'installation ci-dessous ne fonctionneront pas tant que la première release n'est pas publiée. Ce dépôt contient aujourd'hui la spécification, la matrice fonctionnelle et l'architecture cible.
>
> ⭐ Mettez une étoile pour suivre l'avancement, ou consultez la [feuille de route](docs/ForgeNet_Document_de_Projet_v2.3.0.md#6-feuille-de-route--suivi-de-projet).

---

## Le problème

Vos données de production — logs, métriques, secrets, bases — finissent chez AWS, GCP ou Azure. Vos factures sont en dollars, votre support est à 5 000 km et 6 heures de décalage, et le moindre outil de diagnostic « intelligent » envoie vos logs à une API tierce.

**ForgeNet part de l'hypothèse inverse : rien ne sort de vos serveurs.** Pas les données, pas les logs, pas les prompts de l'assistant IA. Installable sur un VPS, un serveur nu ou une machine de bureau.

## ✨ Fonctionnalités

### Le socle

- **Déploiement en un clic** — Docker, Docker Compose, ou détection automatique de buildpack (Node, Python, Go, PHP, Ruby, statique). Pas de Dockerfile requis.
- **Push-to-deploy** — connectez un dépôt GitHub/GitLab/Gitea, chaque push déclenche le déploiement de la branche configurée.
- **Blue-Green & Canary** — bascule de trafic progressive, rollback automatique si le health check échoue. Zéro downtime.
- **HTTPS automatique** — Caddy génère et renouvelle les certificats Let's Encrypt sans un octet de configuration.
- **Marketplace 1-click** — Postgres, MySQL, MongoDB, Redis, MinIO, et un catalogue communautaire extensible.
- **Monitoring temps réel** — CPU, RAM, disque, réseau par hôte et par conteneur, en streaming WebSocket.
- **Backups vers MinIO self-hosted** — planifiés, vérifiés par checksum, restaurables en un clic. Aucun stockage cloud tiers.

### Ce qui n'existe nulle part ailleurs

| | |
| --- | --- |
| 🧠 **IA sysadmin 100 % locale** | Un assistant propulsé par [Ollama](https://ollama.com) qui analyse vos logs, audite votre sécurité et suggère un right-sizing — **sans appel à une API externe**. Vos logs ne quittent jamais la machine. Décisif en contexte bancaire ou administratif. |
| 📶 **Conçu pour le bas débit** | Mode dégradé avec file d'attente locale, rejeu automatique au retour du réseau, synchronisation différentielle et flux compressés. Pensé pour une connectivité réelle, pas pour un datacenter européen. |
| 📱 **App mobile native (Flutter)** | Surveillance temps réel, terminal SSH, alertes push et redémarrage d'urgence — avec confirmation biométrique. Le cas « incident à 2 h du matin », traité sérieusement. |
| 🔄 **Migration depuis vos outils actuels** | Import automatique de vos projets Coolify, CapRover, Dokploy et Portainer. Basculer prend des minutes, pas un week-end. |
| 🌍 **Souveraineté et facturation locales** | Résidence des données en zone CEDEAO/OHADA, facturation en FCFA (Mobile Money, carte), support en français. |

> La liste exhaustive — 150+ fonctionnalités réparties sur 5 modules — est dans la **[matrice fonctionnelle](docs/ForgeNet_Matrice_Fonctionnelle_Complete_Dev.md)**.

## 🚀 Installation

> ⚠️ Non fonctionnel avant la première release — voir [État du projet](#-état-du-projet).

```bash
curl -fsSL https://forgenet.lucidforgeafrica.com/install.sh | sh
```

Le script installe le Core Engine, Caddy et le dashboard. Ouvrez ensuite `https://<votre-serveur>` et suivez l'assistant de configuration.

**Prérequis** — Linux x86_64 ou ARM64, Docker 24+, 2 Go de RAM (4 Go recommandés avec le moteur IA activé).

<details>
<summary>Installation via Docker Compose</summary>

```bash
git clone https://github.com/forgelab-dev/lfa-forgenet.git
cd lfa-forgenet/deploy
docker compose up -d
```

</details>

## 🤔 Pourquoi pas Coolify ou Dokploy ?

Ce sont d'excellents outils, matures et gratuits. Si vous êtes en Europe ou en Amérique du Nord et que le diagnostic assisté par IA ne vous intéresse pas, ils répondront très bien à votre besoin — nous ne prétendons pas le contraire.

ForgeNet existe pour ce qu'ils ne traitent pas :

| | ForgeNet | Coolify | Dokploy | CapRover |
| --- | :---: | :---: | :---: | :---: |
| Self-hosted gratuit, sans limite | ✅ | ✅ | ✅ | ✅ |
| HTTPS automatique | ✅ | ✅ | ✅ | ✅ |
| **IA d'administration 100 % locale** | ✅ | ❌ | ❌ | ❌ |
| **App mobile native** | ✅ | ❌ | ❌ | ❌ |
| **Mode dégradé / bas débit** | ✅ | ❌ | ❌ | ❌ |
| **Import depuis les autres PaaS** | ✅ | ❌ | ❌ | ❌ |
| **Facturation FCFA / Mobile Money** | ✅ | ❌ | ❌ | ❌ |
| Support en français | ✅ | ❌ | ❌ | ❌ |
| Blue-Green / Canary | ✅ | partiel | partiel | ❌ |

## 🏗 Architecture

| Composant | Stack | Rôle |
| --- | --- | --- |
| **Core Engine** | Go, Gin/Echo, gRPC | Orchestration Docker, métriques (`gopsutil`), pilotage Caddy, appels IA |
| **Dashboard Web** | React, Vite, Tailwind | Administration complète : marketplace, logs, domaines, bases |
| **App Mobile** | Flutter | Supervision temps réel, alerting push, gestion de crise |
| **Moteur IA** | Ollama (Llama 3, Mistral, CodeLlama) | Analyse de logs, audit sécurité, assistance — en local |
| **Reverse Proxy** | Caddy | Routage, TLS automatique, load balancing |

<details>
<summary>Structure du dépôt (monorepo)</summary>

```
├── api/          # Contrats gRPC + OpenAPI — source de vérité, codegen Go/TS/Dart
├── cmd/          # Binaires : forgenet-server, forgenet-agent, forgenet-cli
├── internal/     # Cœur métier Go (orchestrator, proxy, metrics, ai, auth…)
├── pkg/          # SDK Go public
├── web/          # Dashboard React + Vite
├── mobile/       # Application Flutter
├── templates/    # Catalogue Marketplace 1-click
├── deploy/       # install.sh, docker-compose, Caddyfile
└── docs/         # Spécification, matrice fonctionnelle, ADRs
```

</details>

## 🧑‍💻 Contribuer

Les contributions sont bienvenues — code, templates de marketplace, traductions, documentation ou simple retour d'usage.

```bash
git clone https://github.com/forgelab-dev/lfa-forgenet.git
cd lfa-forgenet
task setup     # dépendances Go, Node et Flutter
task dev       # stack complète en local
task test      # suite de tests
```

**Conventions** — [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `chore:`), SemVer sur le Core Engine, revue par Pull Request obligatoire avant merge sur `main`. Consultez [CONTRIBUTING.md](CONTRIBUTING.md) avant votre première PR.

Vous cherchez par où commencer ? Les issues étiquetées [`good first issue`](https://github.com/forgelab-dev/lfa-forgenet/labels/good%20first%20issue) sont faites pour ça.

## 📦 Éditions

ForgeNet suit un modèle **open-core** : le cœur du produit est gratuit et le restera.

| Édition | Contenu | Prix |
| --- | --- | --- |
| **Community** | Déploiement, monitoring, reverse proxy, IA locale, backups — sans limite de conteneurs ni de serveurs | Gratuit, Apache 2.0 |
| **Enterprise** | SSO/SAML, multi-tenant & marque blanche, HA, rotation de secrets, export de conformité | Sur devis |
| **Managed Cloud** | ForgeNet hébergé et opéré par LucidForge Africa, facturé en FCFA | À partir de ~5 000 FCFA/mois |

Aucune fonctionnalité Community n'est retirée ou bridée pour pousser vers une édition payante.

## ❤️ Soutenir le projet

ForgeNet est développé en Afrique de l'Ouest avec des moyens modestes. Si le projet vous est utile :

- ⭐ **Mettez une étoile** — c'est gratuit et c'est ce qui nous rend visibles
- 💜 **[Faire un don](https://forgenet.lucidforgeafrica.com/soutenir)** — FedaPay ou Kkiapay en FCFA (Orange Money, MTN MoMo, Moov, carte), ou PayPal depuis l'international
- 🗣 **Parlez-en** — un article, un tweet, une démo à votre équipe

## 📄 Licence

[Apache 2.0](LICENSE) — © [LucidForge Africa](https://lucidforgeafrica.com)

<div align="center">
<sub>Construit en Afrique de l'Ouest 🌍 — parce que la souveraineté numérique ne devrait pas être un produit d'importation.</sub>
</div>
