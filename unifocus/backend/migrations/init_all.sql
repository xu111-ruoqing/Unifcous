-- 001_init_schema.up.sql (reconstructed)

CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name TEXT,
    avatar_url TEXT,
    role TEXT NOT NULL DEFAULT 'student',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS user_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    school TEXT,
    college TEXT,
    major TEXT,
    grade TEXT,
    interests TEXT[],
    skills TEXT[],
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS competition_level_rules (
    id BIGSERIAL PRIMARY KEY,
    competition_name TEXT NOT NULL UNIQUE,
    short_name TEXT,
    level TEXT,
    certification_source TEXT,
    keywords TEXT[],
    url_patterns TEXT[],
    points_value INT DEFAULT 0,
    difficulty_level INT DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_competition_rules_updated_at BEFORE UPDATE ON competition_level_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS opportunities (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'competition',
    description TEXT,
    source_url TEXT,
    source_type TEXT,
    competition_level TEXT,
    certification_type TEXT,
    organizer TEXT,
    organizer_type TEXT,
    award_level TEXT,
    points_value INT DEFAULT 0,
    is_official BOOLEAN NOT NULL DEFAULT false,
    start_date TIMESTAMPTZ,
    deadline TIMESTAMPTZ,
    event_date TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    location TEXT,
    time_info_raw TEXT,
    location_raw TEXT,
    requirements TEXT,
    eligibility_rules TEXT,
    target_majors TEXT[],
    tags TEXT[],
    attachments JSONB,
    description_vector REAL[],
    is_active BOOLEAN NOT NULL DEFAULT true,
    view_count INT NOT NULL DEFAULT 0,
    save_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TRIGGER update_opportunities_updated_at BEFORE UPDATE ON opportunities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE INDEX IF NOT EXISTS idx_opportunities_type ON opportunities(type);
CREATE INDEX IF NOT EXISTS idx_opportunities_deadline ON opportunities(deadline);
CREATE INDEX IF NOT EXISTS idx_opportunities_is_active ON opportunities(is_active);
CREATE INDEX IF NOT EXISTS idx_opportunities_tags ON opportunities USING GIN(tags);

CREATE TABLE IF NOT EXISTS user_opportunities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    opportunity_id BIGINT REFERENCES opportunities(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'saved',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, opportunity_id)
);
CREATE TRIGGER update_user_opportunities_updated_at BEFORE UPDATE ON user_opportunities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS crawl_tasks (
    id BIGSERIAL PRIMARY KEY,
    target_url TEXT NOT NULL UNIQUE,
    site_name TEXT,
    frequency TEXT NOT NULL DEFAULT 'daily',
    status TEXT NOT NULL DEFAULT 'pending',
    last_crawl_at TIMESTAMPTZ,
    next_crawl_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS schedules (
    id BIGSERIAL PRIMARY KEY,
    opportunity_id BIGINT REFERENCES opportunities(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    remind_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS nlp_tasks (
    id BIGSERIAL PRIMARY KEY,
    opportunity_id BIGINT REFERENCES opportunities(id) ON DELETE CASCADE,
    task_type TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 002_recognition_policy.up.sql

-- ============================================
-- 1. Modify opportunities table
-- ============================================
ALTER TABLE opportunities RENAME COLUMN points_value TO default_points_value;

CREATE INDEX IF NOT EXISTS idx_opportunities_requirements_gin
ON opportunities USING GIN (to_tsvector('simple', COALESCE(requirements, '')));

COMMENT ON COLUMN opportunities.default_points_value IS '默认分值，不代表各专业最终认定分值';

-- ============================================
-- 2. New Table: recognition_profiles
-- ============================================
CREATE TABLE recognition_profiles (
    id BIGSERIAL PRIMARY KEY,
    school TEXT NOT NULL,
    college TEXT,
    major TEXT,
    entry_year INT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TRIGGER update_recognition_profiles_updated_at BEFORE UPDATE ON recognition_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 3. New Table: recognition_rules
-- ============================================
CREATE TABLE recognition_rules (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT REFERENCES recognition_profiles(id) ON DELETE CASCADE, -- NULL for global default rules
    enabled BOOLEAN NOT NULL DEFAULT true,
    priority INT NOT NULL DEFAULT 100, -- Smaller value means higher priority
    
    -- Match Conditions
    match_title_regex TEXT,
    match_name_keywords TEXT[],
    match_organizer_regex TEXT,
    only_tracks TEXT[],
    
    -- Output results
    output_competition_level TEXT,
    output_grade TEXT, -- A/B/C
    points_by_award INT[], -- e.g., {100, 80, 60}
    rationale_template TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX idx_recognition_rules_profile ON recognition_rules(profile_id);
CREATE INDEX idx_recognition_rules_enabled_priority ON recognition_rules(enabled, priority);

CREATE TRIGGER update_recognition_rules_updated_at BEFORE UPDATE ON recognition_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 4. New Table: opportunity_recognitions (Cache)
-- ============================================
CREATE TABLE opportunity_recognitions (
    opportunity_id BIGINT REFERENCES opportunities(id) ON DELETE CASCADE,
    profile_id BIGINT REFERENCES recognition_profiles(id) ON DELETE CASCADE,
    
    computed_level TEXT,
    computed_grade TEXT,
    computed_points INT,
    confidence NUMERIC,
    rationale JSONB,
    
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    
    PRIMARY KEY (opportunity_id, profile_id)
);
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
-- Add competitions.name_key for stable de-duplication (quotes/whitespace-insensitive)
-- and clean existing data.

ALTER TABLE competitions ADD COLUMN IF NOT EXISTS name_key TEXT;

UPDATE competitions
SET name_key = lower(
    regexp_replace(
        regexp_replace(
            regexp_replace(name, '["“”‘’`]', '', 'g'),
            '[（）()]',
            '',
            'g'
        ),
        '\s+',
        '',
        'g'
    )
)
WHERE name_key IS NULL OR name_key = '';

-- Drop rows without an official URL (frontend needs a clickable link).
DELETE FROM competitions WHERE official_url IS NULL OR btrim(official_url) = '';

-- De-duplicate by name_key: keep the smallest id.
WITH ranked AS (
    SELECT id,
           name_key,
           row_number() OVER (PARTITION BY name_key ORDER BY id) AS rn
    FROM competitions
    WHERE name_key IS NOT NULL AND name_key <> ''
)
DELETE FROM competitions c
USING ranked r
WHERE c.id = r.id AND r.rn > 1;

ALTER TABLE competitions ALTER COLUMN name_key SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_competitions_name_key ON competitions (name_key);
-- Add timeline fields to competitions table
ALTER TABLE competitions
    ADD COLUMN IF NOT EXISTS category TEXT,
    ADD COLUMN IF NOT EXISTS typical_time_window TEXT,
    ADD COLUMN IF NOT EXISTS timeline_hint TEXT,
    ADD COLUMN IF NOT EXISTS note TEXT;
