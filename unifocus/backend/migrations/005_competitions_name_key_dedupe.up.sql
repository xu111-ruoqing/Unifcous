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
