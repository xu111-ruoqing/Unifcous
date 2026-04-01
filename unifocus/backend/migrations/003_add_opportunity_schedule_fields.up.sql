-- 003_add_opportunity_schedule_fields.up.sql

ALTER TABLE opportunities
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS deadline TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS event_start_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS event_end_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS location TEXT,
    ADD COLUMN IF NOT EXISTS time_info_raw TEXT,
    ADD COLUMN IF NOT EXISTS location_raw TEXT;

CREATE INDEX IF NOT EXISTS idx_opportunities_deadline_ts ON opportunities(deadline);
CREATE INDEX IF NOT EXISTS idx_opportunities_event_start_ts ON opportunities(event_start_at);
