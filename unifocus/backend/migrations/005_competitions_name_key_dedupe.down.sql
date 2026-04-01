DROP INDEX IF EXISTS ux_competitions_name_key;
ALTER TABLE competitions ALTER COLUMN name_key DROP NOT NULL;
ALTER TABLE competitions DROP COLUMN IF EXISTS name_key;

