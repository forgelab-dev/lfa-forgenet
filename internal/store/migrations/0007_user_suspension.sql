-- Migration 0007 : suspension d'un compte (#103). Vide : le compte est actif.
-- Un compte suspendu garde ses données et son historique ; il ne peut plus se
-- connecter, et ses sessions, jetons et liens sont fermés au moment de la
-- suspension.

ALTER TABLE users ADD COLUMN suspended_at TIMESTAMP;
