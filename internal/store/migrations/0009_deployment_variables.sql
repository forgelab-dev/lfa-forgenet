-- Migration 0009 : installations du Marketplace (ADR-0019, #116).
-- template garde l'identifiant du template installé ; vide pour un fichier
-- Compose déployé à la main.
-- variables_sealed garde toutes les valeurs de l'installation, chiffrées
-- ensemble par le serveur (AES-256-GCM, clé FORGENET_SECRET_KEY) : la base
-- n'en garde jamais la forme claire. NULL sans variable.

ALTER TABLE deployments ADD COLUMN template TEXT NOT NULL DEFAULT '';
ALTER TABLE deployments ADD COLUMN variables_sealed BLOB;
