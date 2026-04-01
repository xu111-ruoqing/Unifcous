-- 002_recognition_policy.up.sql

-- ============================================
-- 1. Modify opportunities table
-- ============================================
ALTER TABLE opportunities RENAME COLUMN points_value TO default_points_value;

CREATE INDEX IF NOT EXISTS idx_opportunities_requirements_gin
ON opportunities USING GIN (requirements);

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
