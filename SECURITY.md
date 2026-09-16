# Politique de sécurité

ForgeNet gère des secrets, des accès serveur et des données de production. Les vulnérabilités sont traitées en priorité, en toute confidentialité.

**Français** · [English](#english)

## Versions prises en charge

ForgeNet est en **pré-alpha** : aucune version n'est encore publiée pour la production. Jusqu'à la première release, seule la branche `main` est prise en charge. Ce tableau sera complété à la sortie de la v0.1.0.

| Version | Prise en charge |
|---|---|
| `main` | ✅ |

## Signaler une vulnérabilité

**N'ouvrez jamais d'issue, de discussion ou de pull request publique pour une vulnérabilité.**

1. **Canal privilégié** : le signalement privé de GitHub. Onglet **Security** du dépôt, puis **Report a vulnerability**, ou directement [ce formulaire](https://github.com/forgelab-dev/lfa-forgenet/security/advisories/new). Seule l'équipe sécurité de ForgeLab y a accès.
2. **Si vous n'avez pas de compte GitHub** : écrivez à **contact@lucidforgeafrica.com** en indiquant « Sécurité ForgeNet » dans l'objet, sans détail technique dans ce premier message. Nous vous répondrons avec un moyen d'échange confidentiel.

Indiquez si possible :

- la version ou le commit concerné ;
- le composant touché (Core Engine, dashboard web, application mobile, agent, CLI…) ;
- les étapes de reproduction ou une preuve de concept ;
- l'impact estimé : ce qu'un attaquant peut obtenir.

## Ce qui se passe ensuite

| Étape | Délai |
|---|---|
| Accusé de réception | sous 72 heures ouvrées |
| Première évaluation (confirmation, gravité) | sous 7 jours |
| Correctif et avis de sécurité publié | selon la gravité, en vous tenant informé |

Nous pratiquons la **divulgation coordonnée** : merci de ne rien publier avant la sortie du correctif, ou au plus tard 90 jours après votre signalement si nous n'avons pas pu corriger d'ici là. Avec votre accord, vous serez crédité dans l'avis de sécurité.

## Périmètre

**Dans le périmètre** : le code de ce dépôt (`cmd/`, `internal/`, `pkg/`, `api/`, `web/`, `mobile/`, `deploy/`, `templates/`).

**Hors périmètre** :

- les vulnérabilités d'une dépendance tierce : signalez-les à son projet d'origine (vous pouvez nous prévenir en parallèle si ForgeNet est exposé) ;
- les problèmes qui supposent un accès déjà administrateur au serveur hôte ;
- les attaques par déni de service volumétrique et l'ingénierie sociale.

---

## English

ForgeNet is **pre-alpha**: only the `main` branch is supported until the first release.

**Never report a vulnerability through a public issue, discussion or pull request.** Use GitHub's private reporting (**Security** tab → **Report a vulnerability**, or [this form](https://github.com/forgelab-dev/lfa-forgenet/security/advisories/new)). Without a GitHub account, email **contact@lucidforgeafrica.com** with "ForgeNet security" as the subject and no technical details in that first message; we will reply with a confidential channel.

We acknowledge reports within 72 business hours and give a first assessment within 7 days. We follow coordinated disclosure: please do not publish anything before the fix is released, or 90 days after your report at the latest. With your consent, you will be credited in the advisory.
