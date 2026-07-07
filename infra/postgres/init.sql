-- Initial schema placeholder for the shareit database.
-- Add migrations here or move to a proper migration tool (e.g. golang-migrate) as the schema grows.

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS photos (
    id UUID PRIMARY KEY,
    original_key TEXT NOT NULL,
    thumb_key TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
