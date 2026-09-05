-- Gán category_id cho pois chưa có category, dựa point_type
-- khớp LOWER(TRIM(point_type)) = LOWER(TRIM(poi_categories.name)).

UPDATE pois p
SET category_id = m.category_id,
    updated_at = NOW()
FROM (
    SELECT DISTINCT ON (LOWER(TRIM(name)))
        id AS category_id,
        LOWER(TRIM(name)) AS name_key
    FROM poi_categories
    WHERE deleted_at IS NULL
      AND COALESCE(TRIM(name), '') <> ''
    ORDER BY LOWER(TRIM(name)), id ASC
) m
WHERE (p.category_id IS NULL OR p.category_id = 0)
  AND p.deleted_at IS NULL
  AND COALESCE(TRIM(p.point_type), '') <> ''
  AND LOWER(TRIM(p.point_type)) = m.name_key;

-- Sync lại poi_count cho các category vừa được gán
UPDATE poi_categories c
SET poi_count = COALESCE(sub.cnt, 0),
    updated_at = NOW()
FROM (
    SELECT category_id AS id, COUNT(*)::int AS cnt
    FROM pois
    WHERE deleted_at IS NULL
      AND category_id IS NOT NULL
      AND category_id > 0
    GROUP BY category_id
) sub
WHERE c.id = sub.id
  AND c.deleted_at IS NULL;
