<div align="center">

# ForgeNet

**The sovereign, self-hosted cloud that never leaves your servers.**

Deploy, monitor and orchestrate your services — with a sysadmin AI that runs *locally*, without a single line of your logs ever reaching a third-party cloud.

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![React](https://img.shields.io/badge/React-Vite-61DAFB?logo=react&logoColor=white)](https://react.dev)
[![Flutter](https://img.shields.io/badge/Flutter-Mobile-02569B?logo=flutter&logoColor=white)](https://flutter.dev)
[![Status](https://img.shields.io/badge/status-pre--alpha-orange.svg)](#-project-status)

**English** · [Français](README.md)

[Features](#-features) · [Install](#-installation) · [Comparison](#-why-not-coolify-or-dokploy) · [Architecture](#-architecture) · [Contributing](#-contributing) · [Support](#%EF%B8%8F-support-the-project)

</div>

---

## 🚧 Project status

> **Pre-alpha — under construction.** The code is not production-ready yet, and the install commands below will not work until the first release ships. Today this repository holds the specification, the feature matrix and the target architecture.
>
> ⭐ Star the repo to follow along, or read the [roadmap](docs/ForgeNet_Document_de_Projet_v2.3.0.md#6-feuille-de-route--suivi-de-projet) (in French).

---

## The problem

Your production data — logs, metrics, secrets, databases — ends up on AWS, GCP or Azure. Your invoices are in dollars, your support desk is 5,000 km and six time zones away, and every "smart" diagnostic tool ships your logs off to a third-party API.

**ForgeNet starts from the opposite assumption: nothing leaves your servers.** Not the data, not the logs, not even the AI assistant's prompts. Installs on a VPS, bare metal, or a machine under your desk.

## ✨ Features

### The foundation
- **One-click deploys** — Docker, Docker Compose, or automatic buildpack detection (Node, Python, Go, PHP, Ruby, static). No Dockerfile required.
- **Push-to-deploy** — connect a GitHub/GitLab/Gitea repository; every push deploys the branch you configured.
- **Blue-Green & Canary** — progressive traffic shifting with automatic rollback when a health check fails. Zero downtime.
- **Automatic HTTPS** — Caddy issues and renews Let's Encrypt certificates without a byte of configuration.
- **1-click marketplace** — Postgres, MySQL, MongoDB, Redis, MinIO, plus an extensible community catalogue.
- **Real-time monitoring** — CPU, RAM, disk and network per host and per container, streamed over WebSocket.
- **Backups to self-hosted MinIO** — scheduled, checksum-verified, restorable in one click. No third-party object storage.

### What you won't find elsewhere

| | |
|---|---|
| 🧠 **100% local sysadmin AI** | An [Ollama](https://ollama.com)-powered assistant that analyses your logs, audits your security posture and suggests right-sizing — **with zero external API calls**. Your logs never leave the machine. Decisive in banking and government contexts. |
| 📶 **Built for low bandwidth** | Degraded mode with a local command queue, automatic replay when connectivity returns, differential sync and compressed streams. Designed for real-world networks, not for a well-connected European datacentre. |
| 📱 **Native mobile app (Flutter)** | Live monitoring, SSH terminal, push alerts and emergency restarts — behind biometric confirmation. The "3 a.m. incident" case, taken seriously. |
| 🔄 **Migrate from your current tools** | Automatic import of your Coolify, CapRover, Dokploy and Portainer projects. Switching takes minutes, not a weekend. |
| 🌍 **Local sovereignty and billing** | Data residency in the ECOWAS/OHADA region, billing in XOF (mobile money, card) or via PayPal, and French-language support. |

> The exhaustive list — 150+ features across 5 modules — lives in the **[feature matrix](docs/ForgeNet_Matrice_Fonctionnelle_Complete_Dev.md)** (in French).

## 🚀 Installation

> ⚠️ Not functional before the first release — see [Project status](#-project-status).

```bash
curl -fsSL https://forgenet.lucidforgeafrica.com/install.sh | sh
```

The script installs the Core Engine, Caddy and the dashboard. Then open `https://<your-server>` and follow the setup wizard.

**Requirements** — Linux x86_64 or ARM64, Docker 24+, 2 GB RAM (4 GB recommended with the AI engine enabled).

<details>
<summary>Install with Docker Compose</summary>

```bash
git clone https://github.com/forgelab-dev/lfa-forgenet.git
cd lfa-forgenet/deploy
docker compose up -d
```

</details>

## 🤔 Why not Coolify or Dokploy?

They are excellent tools — mature, free, and widely deployed. If you are in Europe or North America and AI-assisted diagnostics don't interest you, they will serve you well; we won't pretend otherwise.

ForgeNet exists for what they don't address:

| | ForgeNet | Coolify | Dokploy | CapRover |
|---|:---:|:---:|:---:|:---:|
| Free self-hosted, no limits | ✅ | ✅ | ✅ | ✅ |
| Automatic HTTPS | ✅ | ✅ | ✅ | ✅ |
| **100% local sysadmin AI** | ✅ | ❌ | ❌ | ❌ |
| **Native mobile app** | ✅ | ❌ | ❌ | ❌ |
| **Low-bandwidth / degraded mode** | ✅ | ❌ | ❌ | ❌ |
| **Import from other PaaS** | ✅ | ❌ | ❌ | ❌ |
| **XOF / mobile money billing** | ✅ | ❌ | ❌ | ❌ |
| French-language support | ✅ | ❌ | ❌ | ❌ |
| Blue-Green / Canary | ✅ | partial | partial | ❌ |

## 🏗 Architecture

| Component | Stack | Role |
|---|---|---|
| **Core Engine** | Go, Gin/Echo, gRPC | Docker orchestration, metrics (`gopsutil`), Caddy control, AI calls |
| **Web Dashboard** | React, Vite, Tailwind | Full administration: marketplace, logs, domains, databases |
| **Mobile App** | Flutter | Live supervision, push alerting, incident response |
| **AI Engine** | Ollama (Llama 3, Mistral, CodeLlama) | Log analysis, security audit, assistance — locally |
| **Reverse Proxy** | Caddy | Routing, automatic TLS, load balancing |

<details>
<summary>Repository layout (monorepo)</summary>

```
├── api/          # gRPC + OpenAPI contracts — source of truth, codegen Go/TS/Dart
├── cmd/          # Binaries: forgenet-server, forgenet-agent, forgenet-cli
├── internal/     # Go domain logic (orchestrator, proxy, metrics, ai, auth…)
├── pkg/          # Public Go SDK
├── web/          # React + Vite dashboard
├── mobile/       # Flutter application
├── templates/    # 1-click marketplace catalogue
├── deploy/       # install.sh, docker-compose, Caddyfile
└── docs/         # Specification, feature matrix, ADRs
```

</details>

## 🧑‍💻 Contributing

Contributions are welcome — code, marketplace templates, translations, documentation, or simply telling us how it went.

```bash
git clone https://github.com/forgelab-dev/lfa-forgenet.git
cd lfa-forgenet
task setup     # Go, Node and Flutter dependencies
task dev       # full local stack
task test      # test suite
```

**Conventions** — [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `chore:`), SemVer on the Core Engine, mandatory pull-request review before merging to `main`. Read [CONTRIBUTING.md](CONTRIBUTING.md) before your first PR.

Looking for a place to start? Issues labelled [`good first issue`](https://github.com/forgelab-dev/lfa-forgenet/labels/good%20first%20issue) are meant for exactly that.

> Project documentation is currently written in French. Contributions in English are welcome, and translating the docs is itself a valuable first contribution.

## 📦 Editions

ForgeNet follows an **open-core** model: the heart of the product is free and will stay that way.

| Edition | What you get | Price |
|---|---|---|
| **Community** | Deployments, monitoring, reverse proxy, local AI, backups — no container or server limits | Free, Apache 2.0 |
| **Enterprise** | SSO/SAML, multi-tenant & white-label, HA, secret rotation, compliance export | On request |
| **Managed Cloud** | ForgeNet hosted and operated by LucidForge Africa, billed in XOF | From ~5,000 XOF/month (≈ $8) |

No Community feature is ever removed or crippled to push you toward a paid edition.

## ❤️ Support the project

ForgeNet is built in West Africa on a modest budget. If it's useful to you:

- ⭐ **Star the repo** — free, and it's what makes us visible
- 💜 **[Donate](https://forgenet.lucidforgeafrica.com/soutenir)** — PayPal for international supporters, or FedaPay/Kkiapay in XOF (Orange Money, MTN MoMo, Moov, card)
- 🗣 **Spread the word** — a blog post, a tweet, a demo to your team

## 📄 License

[Apache 2.0](LICENSE) — © [LucidForge Africa](https://lucidforgeafrica.com)

<div align="center">
<sub>Built in West Africa 🌍 — because digital sovereignty shouldn't have to be imported.</sub>
</div>
