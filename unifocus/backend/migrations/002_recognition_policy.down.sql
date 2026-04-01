-- 002_recognition_policy.down.sql

-- Drop new tables (Order matters due to FK constraints)
DROP TABLE IF EXISTS opportunity_recognitions;
DROP TABLE IF EXISTS recognition_rules;
DROP TABLE IF EXISTS recognition_profiles;

-- Revert opportunities table changes
DROP INDEX IF EXISTS idx_opportunities_requirements_gin;

ALTER TABLE opportunities RENAME COLUMN default_points_value TO points_value;
