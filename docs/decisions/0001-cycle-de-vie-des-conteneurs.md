# ADR-0001 : création et suppression des conteneurs par DeploymentService

- Statut : accepté
- Date : 2026-10-14
- Contexte : revue de la PR #28 (ticket #6)

## Contexte

`orchestrator.Engine` sait créer et supprimer un conteneur, mais le contrat
`ContainerService` ne l'expose pas. Il fallait choisir entre ajouter ces deux
méthodes à `ContainerService` ou les réserver au service de déploiement.

## Décision

- `ContainerService` n'expose que la lecture et le cycle d'exécution :
  `ListContainers`, `GetContainer`, `StreamContainerLogs`, `StartContainer`,
  `StopContainer`, `RestartContainer`.
- La création et la suppression sont des opérations métier de
  `DeploymentService` (#7). `Engine.Create` et `Engine.Remove` restent des
  primitives internes qu'il est seul à appeler.
- Toute création et toute suppression sont inscrites au journal d'audit
  (événement en base et ligne de journal structuré), réussies ou non.

## Conséquences

- La logique de déploiement (stratégies, quotas, approbations) vit à un seul
  endroit, au lieu d'être répartie entre deux services.
- Un client qui veut un conteneur passe par un déploiement ; il n'existe pas de
  conteneur « hors déploiement » créé par l'API.
- Ajouter `CreateContainer` ou `DeleteContainer` au contrat demande une
  nouvelle décision.
