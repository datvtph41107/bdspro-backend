DROP INDEX IF EXISTS idx_pois_adr_search_trgm;

ALTER TABLE pois
    DROP COLUMN IF EXISTS adr_search;
