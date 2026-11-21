# Page d'exemple

Le plus petit template qui passe la vérification : nginx sert une page dont le
titre est demandé à l'installation. Il sert de modèle à copier, et n'apparaît
pas dans le catalogue (son dossier commence par `_`).

## Variables

| Variable | Rôle | Par défaut |
|---|---|---|
| `SITE_TITRE` | Titre affiché sur la page | `Bonjour depuis ForgeNet` |

## Accès

Le service `web` écoute sur le port 80 de son conteneur (`http`) : reliez-le à
un domaine depuis l'écran **Domaines**, et Caddy le sert en HTTPS.

## Données

Le volume `site` garde les fichiers de la page.
