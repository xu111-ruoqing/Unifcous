-- 003_opportunities_minimal.up.sql

ALTER TABLE opportunities
    ADD COLUMN IF NOT EXISTS notice_type TEXT NOT NULL DEFAULT 'notice',
    ADD COLUMN IF NOT EXISTS competition_level TEXT,
    ADD COLUMN IF NOT EXISTS registration_start_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS registration_deadline TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS event_start_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS event_end_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS raw_text TEXT;

CREATE INDEX IF NOT EXISTS idx_opportunities_source_url ON opportunities(source_url);
CREATE INDEX IF NOT EXISTS idx_opportunities_notice_type ON opportunities(notice_type);
CREATE INDEX IF NOT EXISTS idx_opportunities_registration_deadline ON opportunities(registration_deadline);
CREATE INDEX IF NOT EXISTS idx_opportunities_event_start_at ON opportunities(event_start_at);
CREATE INDEX IF NOT EXISTS idx_opportunities_level ON opportunities(competition_level);

CREATE OR REPLACE VIEW v_opportunities_minimal AS
SELECT
    id,
    title,
    source_url,
    notice_type,
    competition_level,
    registration_start_at,
    registration_deadline,
    event_start_at,
    event_end_at,
    raw_text,
    created_at
FROM opportunities;
