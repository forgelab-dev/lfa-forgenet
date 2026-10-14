# ADR-0002 : détails de /readyz selon la provenance de la requête

- Statut : accepté
- Date : 2026-10-14
- Contexte : revue de la PR #28 (ticket #6)

## Contexte

`/readyz` détaille l'état de chaque dépendance. Les messages d'erreur peuvent
révéler des chemins internes (socket Docker, fichier de base). C'est utile pour
diagnostiquer, risqué une fois le serveur exposé.

## Décision

Réglage `health.expose_details` (`FORGENET_HEALTH_EXPOSE_DETAILS`) :

| Valeur | Détails renvoyés |
|---|---|
| `always` | toujours |
| `local` (par défaut) | seulement à une requête locale directe |
| `never` | jamais |

- Une requête est « locale directe » si elle vient d'une adresse de bouclage
  (127.0.0.1, ::1) **et** ne porte aucun en-tête de proxy (`Forwarded`,
  `X-Forwarded-For`, `X-Real-Ip`). Derrière Caddy sur la même machine, toutes
  les requêtes arrivent de 127.0.0.1 : sans cette seconde condition, `local`
  exposerait les détails à Internet.
- Le statut (`SERVING`/`NOT_SERVING`) et le code HTTP (200/503) ne changent
  pas d'un mode à l'autre.
- La cause d'une indisponibilité est toujours écrite dans le journal du
  serveur, quel que soit le mode.

## Conséquences

- Un opérateur connecté à la machine diagnostique avec `curl localhost:8080/readyz`.
- Un proxy qui retirerait les en-têtes de transfert ferait passer ses requêtes
  pour locales : la configuration Caddy (#10) doit les conserver, ce qu'elle
  fait par défaut, ou le serveur doit passer en `never`.
