-- Migration 0008 : relais SMTP réglé depuis le dashboard (ADR-0015, #102).
-- Une seule ligne (id = 1) ; sans elle, les réglages de .env s'appliquent.
-- Le mot de passe du relais est chiffré par le serveur (AES-256-GCM, clé
-- FORGENET_SECRET_KEY) : la base n'en garde jamais la forme claire.

CREATE TABLE mail_settings (
    id              INTEGER PRIMARY KEY CHECK (id = 1),
    disabled        INTEGER NOT NULL DEFAULT 0,
    host            TEXT NOT NULL DEFAULT '',
    port            INTEGER NOT NULL DEFAULT 587,
    tls             TEXT NOT NULL DEFAULT 'starttls',
    username        TEXT NOT NULL DEFAULT '',
    password_sealed BLOB,
    from_address    TEXT NOT NULL DEFAULT '',
    updated_by      TEXT NOT NULL DEFAULT '',
    updated_at      TIMESTAMP NOT NULL
);
