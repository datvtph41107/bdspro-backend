CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE pois
    ADD COLUMN IF NOT EXISTS adr_search TEXT;

UPDATE pois
SET adr_search = TRIM(
    REGEXP_REPLACE(
        LOWER(UNACCENT(COALESCE(address, ''))),
        '[^a-z0-9]+',
        ' ',
        'g'
    )
)
WHERE COALESCE(adr_search, '') = ''
  AND COALESCE(address, '') <> '';

CREATE INDEX IF NOT EXISTS idx_pois_adr_search_trgm
    ON pois USING GIN (LOWER(adr_search) gin_trgm_ops)
    WHERE COALESCE(adr_search, '') <> '';
