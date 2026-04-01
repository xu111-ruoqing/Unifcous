-- Drop competitions master table (phase-1 minimal storage)

DROP INDEX IF EXISTS idx_competitions_official_url;
DROP INDEX IF EXISTS idx_competitions_name_unique;
DROP TABLE IF EXISTS competitions;

