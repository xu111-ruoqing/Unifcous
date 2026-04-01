-- 003_opportunities_minimal.down.sql

DROP VIEW IF EXISTS v_opportunities_minimal;

DROP INDEX IF EXISTS idx_opportunities_level;
DROP INDEX IF EXISTS idx_opportunities_event_start_at;
DROP INDEX IF EXISTS idx_opportunities_registration_deadline;
DROP INDEX IF EXISTS idx_opportunities_notice_type;
DROP INDEX IF EXISTS idx_opportunities_source_url;

ALTER TABLE opportunities
    DROP COLUMN IF EXISTS raw_text,
    DROP COLUMN IF EXISTS event_end_at,
    DROP COLUMN IF EXISTS event_start_at,
    DROP COLUMN IF EXISTS registration_deadline,
    DROP COLUMN IF EXISTS registration_start_at,
    DROP COLUMN IF EXISTS competition_level,
    DROP COLUMN IF EXISTS notice_type;
