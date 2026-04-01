-- Revert timeline fields from competitions table
ALTER TABLE competitions
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS typical_time_window,
    DROP COLUMN IF EXISTS timeline_hint,
    DROP COLUMN IF EXISTS note;
