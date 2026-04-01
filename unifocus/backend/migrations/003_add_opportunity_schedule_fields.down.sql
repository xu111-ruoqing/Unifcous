-- 003_add_opportunity_schedule_fields.down.sql

DROP INDEX IF EXISTS idx_opportunities_deadline_ts;
DROP INDEX IF EXISTS idx_opportunities_event_start_ts;

ALTER TABLE opportunities
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS deadline,
    DROP COLUMN IF EXISTS event_start_at,
    DROP COLUMN IF EXISTS event_end_at,
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS time_info_raw,
    DROP COLUMN IF EXISTS location_raw;
