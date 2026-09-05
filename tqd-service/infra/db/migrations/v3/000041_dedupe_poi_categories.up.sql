-- Gộp bản ghi poi_categories trùng name (case-insensitive / trim).
-- Giữ 1 bản ghi / name_key, remap pois.category_id (+ parent_id), rồi xóa bản trùng.
--
-- keep_id ưu tiên: chưa soft-delete → poi_count cao hơn → id nhỏ hơn.

CREATE TEMP TABLE tmp_poi_category_dedupe (
    dupe_id BIGINT PRIMARY KEY,
    keep_id BIGINT NOT NULL
) ON COMMIT DROP;

WITH ranked AS (
    SELECT
        id,
        LOWER(TRIM(name)) AS name_key,
        ROW_NUMBER() OVER (
            PARTITION BY LOWER(TRIM(name))
            ORDER BY
                CASE WHEN deleted_at IS NULL THEN 0 ELSE 1 END,
                COALESCE(poi_count, 0) DESC,
                id ASC
        ) AS rn
    FROM poi_categories
),
keepers AS (
    SELECT id AS keep_id, name_key
    FROM ranked
    WHERE rn = 1
)
INSERT INTO tmp_poi_category_dedupe (dupe_id, keep_id)
SELECT r.id, k.keep_id
FROM ranked r
JOIN keepers k ON k.name_key = r.name_key
WHERE r.rn > 1;

-- 1) Remap pois.category_id: id bị xóa → id giữ lại
UPDATE pois p
SET category_id = d.keep_id,
    updated_at = NOW()
FROM tmp_poi_category_dedupe d
WHERE p.category_id = d.dupe_id;

-- 2) Remap parent_id trỏ tới bản trùng (tránh FK khi DELETE)
UPDATE poi_categories c
SET parent_id = CASE
        WHEN d.keep_id = c.id THEN NULL  -- tránh self-parent
        ELSE d.keep_id
    END,
    updated_at = NOW()
FROM tmp_poi_category_dedupe d
WHERE c.parent_id = d.dupe_id;

-- 3) directory_sources.category_id (nếu có)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'directory_sources'
          AND column_name = 'category_id'
    ) THEN
        EXECUTE $sql$
            UPDATE directory_sources ds
            SET category_id = d.keep_id
            FROM tmp_poi_category_dedupe d
            WHERE ds.category_id = d.dupe_id
        $sql$;
    END IF;
END $$;

-- 4) Sync lại poi_count cho các category được giữ lại
UPDATE poi_categories c
SET poi_count = COALESCE(sub.cnt, 0),
    updated_at = NOW()
FROM (
    SELECT DISTINCT keep_id FROM tmp_poi_category_dedupe
) k
LEFT JOIN (
    SELECT category_id AS id, COUNT(*)::int AS cnt
    FROM pois
    WHERE deleted_at IS NULL
      AND category_id IS NOT NULL
    GROUP BY category_id
) sub ON sub.id = k.keep_id
WHERE c.id = k.keep_id;

-- 5) Xóa bản trùng
DELETE FROM poi_categories c
USING tmp_poi_category_dedupe d
WHERE c.id = d.dupe_id;
