# ADR-0004 : exécution des déploiements et service `deploy`

- Statut : proposé
- Date : 2026-10-29
- Contexte : ticket #7, troisième PR (service gRPC `deploy`)

## Contexte

`CreateDeployment` doit répondre tout de suite, avec l'état `PENDING` : un
`docker compose up` qui attend des sondes de santé dure des secondes, parfois des
minutes. Il faut donc exécuter les déploiements ailleurs que dans l'appel, et
décider de l'ordre, de la reprise après un arrêt et de qui peut déployer. Le
contrat ne sait pas non plus supprimer une stack, ce que demande le ticket.

## Décision

- **La file est la table `deployments`.** L'exécuteur (`internal/deploy`) prend
  les déploiements `pending`, du plus ancien au plus récent ; `Create` le réveille,
  et il relit la table toutes les 5 secondes. Pas de file en mémoire : rien ne se
  perd à un redémarrage.
- **Un seul déploiement s'exécute à la fois**, sur tout le serveur. Deux mises à
  jour d'une même stack ne se croisent jamais. C'est plus lent quand plusieurs
  projets déploient en même temps : on parallélisera par projet quand le besoin
  sera mesuré.
- **Les étapes d'une source Compose** : `building` vérifie le fichier (il n'y a
  rien à construire), `deploying` lance `up --wait`, `health_checking` relit les
  conteneurs (tous en marche, et sains s'ils ont une sonde), puis `success` ou
  `failed` avec le message de Compose. Un déploiement est borné à 15 minutes.
- **Au démarrage**, un déploiement surpris en cours par un arrêt passe en
  `failed` (« interrompu par l'arrêt du serveur ») ; ceux en attente sont exécutés.
  À l'arrêt, celui en cours a le délai de grâce du serveur pour finir.
- **Annuler** fonctionne en `pending`, `building` et `deploying`, comme le dit le
  contrat. Une annulation pendant `up` peut laisser la stack à moitié à jour : un
  nouveau déploiement la remet d'aplomb.
- **`DeleteStack` s'ajoute au contrat** (changement compatible) : il supprime les
  conteneurs et les réseaux de la stack, les volumes sur demande, et garde
  l'historique. Il est refusé tant qu'un déploiement du projet attend ou
  s'exécute, sinon ce déploiement recréerait la stack juste après.
- **Tout le service est réservé aux administrateurs** (décision du 29/10), lecture
  comprise : un fichier Compose donne un accès complet à l'hôte et peut contenir
  des secrets. La politique d'`auth.Guard` le fait déjà par défaut ; un test le
  vérifie.
- Une stack s'appelle `forgenet-<projet>`, ce qui la distingue des projets Compose
  de l'hôte. Un identifiant de projet suit la règle de Compose (minuscules,
  chiffres, `_` et `-`, 63 caractères au plus).

## Conséquences

- Les sources Git et image, les stratégies autres que la recréation, le retour en
  arrière et les flux (`StreamDeployment`, `StreamDeploymentLogs`) répondent
  `UNIMPLEMENTED` : ils viendront avec leurs tickets.
- Un service sans sonde de santé qui s'arrête de lui-même (une tâche ponctuelle)
  fait échouer le déploiement : le contrôle exige des conteneurs en marche.
- Ouvrir la lecture aux membres demandera de masquer le fichier Compose dans les
  réponses, et une nouvelle décision.
