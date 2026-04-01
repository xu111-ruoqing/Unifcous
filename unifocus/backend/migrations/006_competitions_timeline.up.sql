-- Add timeline fields to competitions table
ALTER TABLE competitions
    ADD COLUMN IF NOT EXISTS category TEXT,
    ADD COLUMN IF NOT EXISTS typical_time_window TEXT,
    ADD COLUMN IF NOT EXISTS timeline_hint TEXT,
    ADD COLUMN IF NOT EXISTS note TEXT;
