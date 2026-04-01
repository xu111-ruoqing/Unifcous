-- Create competitions master table (phase-1 minimal storage)
-- One row per national competition (no news/province items).

CREATE TABLE IF NOT EXISTS competitions (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    official_url TEXT,
    level TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- If a previous version created official_url as UNIQUE, drop it (multiple competitions can share the same site).
ALTER TABLE competitions DROP CONSTRAINT IF EXISTS competitions_official_url_key;

-- If a previous version required official_url NOT NULL, relax it (some competitions may not have an official URL yet).
ALTER TABLE competitions ALTER COLUMN official_url DROP NOT NULL;

-- Ensure idempotent upsert by competition name.
CREATE UNIQUE INDEX IF NOT EXISTS idx_competitions_name_unique ON competitions (name);

-- Query optimization for frontend lookups.
CREATE INDEX IF NOT EXISTS idx_competitions_official_url ON competitions (official_url);
