# ADR-0017 : télémétrie opt-in, compter les installations actives

- Statut : accepté pour les principes (validés par le porteur du projet le
  2026-11-20) ; les points ouverts, en fin de document, se tranchent avant
  le code
- Date : 2026-11-20
- Contexte : ticket #109 ; matrice fonctionnelle §1.8 ; document de projet
  §6 (indicateurs de la phase 1)

## Contexte

On ne sait pas combien de serveurs font tourner ForgeNet. La vérification
des mises à jour (ADR-0011) interroge l'API de GitHub directement : nous ne
voyons pas ces appels. Les statistiques de GitHub ne servent pas non plus :
les clones viennent surtout de la CI, les téléchargements surtout de nos
serveurs de test.

Le document de projet fait pourtant du nombre d'installations actives un
indicateur de la phase 1, à une condition : « télémétrie opt-in uniquement,
jamais forcée », parce que ForgeNet promet que rien ne quitte le serveur
sans l'accord de son administrateur.

## Décision

### Le consentement

- **Désactivée par défaut**, sur une nouvelle installation comme sur un
  serveur mis à jour. Sans accord explicite, rien ne part.
- **Proposée deux fois, jamais imposée** :
  - `install.sh` pose la question, réponse « non » par défaut ; sans
    terminal (cloud-init), elle reste désactivée ; `--telemetry on` ou
    `--telemetry off` répond d'avance ;
  - la page **Réglages** (administrateurs) l'active ou la désactive.
- **Désactivable à tout moment**, effet immédiat : plus aucun envoi.
- Chaque changement est noté au journal d'audit (`telemetry.enable`,
  `telemetry.disable`), avec le compte qui l'a fait.

### Ce qui part, et rien d'autre

Un seul message, en JSON :

```json
{
  "schema": 1,
  "id": "3f6c0b0e9a1d4c2f8e5b7a6d1c0e9f42",
  "version": "0.12.0",
  "arch": "amd64"
}
```

- `id` : 16 octets tirés au hasard à l'activation, en hexadécimal. Il ne
  dérive de rien (ni machine, ni domaine, ni compte). Désactiver l'efface ;
  réactiver en tire un nouveau, qu'on ne peut pas relier à l'ancien.
- `version` et `arch` : la version de ForgeNet et l'architecture
  (`amd64`, `arm64`).
- **Jamais** : domaine, adresse IP du serveur, nom ou nombre de comptes,
  applications déployées, journaux, métriques, configuration.

### La transparence

- La page Réglages affiche, avant l'activation, le message exact qui
  partira et l'adresse qui le reçoit, avec un lien vers cet ADR.
- Le code qui construit le message est le seul à l'envoyer, et un test
  vérifie qu'il ne contient que ces champs.

### Le destinataire et la fréquence

- Un point de collecte à nous, sur un domaine de LucidForge Africa
  (`https://telemetrie.lucidforgeafrica.com/v1/ping`). GitHub ne peut pas
  servir : il ne nous montre pas qui l'interroge.
- **Une fois par jour**, à une heure tirée au hasard pour étaler les
  envois, en HTTPS seulement. Un échec est ignoré, sans nouvel essai avant
  le lendemain : un serveur hors ligne ou derrière un pare-feu ne remarque
  rien.
- À la désactivation, un dernier message demande l'effacement de l'`id`
  côté collecte.

### Ce que le point de collecte garde

- Par `id` : version, architecture, date du premier et du dernier message.
- **Pas l'adresse IP** : il la voit le temps de la connexion, comme tout
  serveur web, mais ne l'écrit ni en base ni dans ses journaux.
- Une installation est **active** si son `id` s'est manifesté dans les
  30 derniers jours. Un `id` silencieux depuis 90 jours est effacé.
- Les chiffres sortent agrégés : nombre d'installations actives, par
  version et par architecture.

## Retour arrière

- Côté serveur ForgeNet : désactiver suffit, et une version antérieure
  n'envoie rien.
- Côté collecte : l'arrêter n'a aucun effet sur les serveurs, dont les
  envois échouent en silence.

## Points ouverts, à trancher avant le code

1. **L'hébergement du point de collecte** : petit service Go à part, ou
   page sur un hébergement existant. Il doit pouvoir ne rien journaliser
   des adresses IP.
2. **Publier les chiffres** : une page publique du nombre d'installations
   actives et de leur répartition par version, ou un usage interne seul.
3. **La question d'`install.sh`** : son texte exact, court et neutre, qui
   dit ce qui part et que la réponse se change plus tard.

## Conséquences

- L'indicateur « installations Community actives » devient mesurable, sur
  la seule part des administrateurs qui l'acceptent : c'est un plancher,
  pas un total.
- Un nouveau service à faire tourner et à surveiller.
- La promesse « rien ne quitte le serveur sans votre accord » reste vraie et
  vérifiable dans le code.
