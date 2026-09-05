-- GIN FTS cho public POI search (không cần cột adr_tsv).
-- Query app: to_tsvector('simple', coalesce(adr_search,'')) @@ to_tsquery('simple', ...)
--
-- Có thể chạy CONCURRENTLY ngoài transaction nếu bảng lớn:
--   CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pois_adr_tsv ...

CREATE INDEX IF NOT EXISTS idx_pois_adr_tsv
    ON pois USING GIN (to_tsvector('simple', coalesce(adr_search, '')))
    WHERE deleted_at IS NULL
      AND is_active = TRUE
      AND COALESCE(adr_search, '') <> '';
