-- Move the versioned Vietnamese administrative reference catalog from the
-- historical tables into the tables owned by the current Location runtime.
-- This is a source-data cutover, not a development/demo fixture.

INSERT INTO province_v2 (
    id,
    full_name,
    short_name,
    lat,
    lng,
    code,
    created_at,
    updated_at
)
SELECT
    id,
    full_name,
    short_name,
    lat,
    lng,
    code,
    COALESCE(created_at, CURRENT_TIMESTAMP),
    COALESCE(created_at, CURRENT_TIMESTAMP)
FROM provinces
ON CONFLICT DO NOTHING;

-- Keep an already-canonical identity when a deployment previously imported
-- the same province code under another ID, while refreshing its public data.
UPDATE province_v2 AS canonical
SET
    full_name = legacy.full_name,
    short_name = legacy.short_name,
    lat = legacy.lat,
    lng = legacy.lng,
    updated_at = CURRENT_TIMESTAMP
FROM provinces AS legacy
WHERE canonical.code = legacy.code;

INSERT INTO ward_v2 (
    id,
    full_name,
    short_name,
    lat,
    lng,
    code,
    province_id,
    province_code,
    created_at,
    updated_at
)
SELECT
    legacy_ward.id,
    legacy_ward.full_name,
    legacy_ward.short_name,
    legacy_ward.lat,
    legacy_ward.lng,
    legacy_ward.code,
    canonical_province.id,
    canonical_province.code,
    COALESCE(legacy_ward.created_at, CURRENT_TIMESTAMP),
    COALESCE(legacy_ward.created_at, CURRENT_TIMESTAMP)
FROM wards AS legacy_ward
JOIN provinces AS legacy_province
    ON legacy_province.id = legacy_ward.province_id
JOIN province_v2 AS canonical_province
    ON canonical_province.code = legacy_province.code
ON CONFLICT DO NOTHING;

-- Refresh an existing canonical ward by its stable public code without
-- changing its identity or breaking downstream foreign keys.
UPDATE ward_v2 AS canonical
SET
    full_name = legacy_ward.full_name,
    short_name = legacy_ward.short_name,
    lat = legacy_ward.lat,
    lng = legacy_ward.lng,
    province_id = canonical_province.id,
    province_code = canonical_province.code,
    updated_at = CURRENT_TIMESTAMP
FROM wards AS legacy_ward
JOIN provinces AS legacy_province
    ON legacy_province.id = legacy_ward.province_id
JOIN province_v2 AS canonical_province
    ON canonical_province.code = legacy_province.code
WHERE canonical.code = legacy_ward.code;

UPDATE province_v2 AS province
SET
    ward_count = counts.ward_count,
    updated_at = CURRENT_TIMESTAMP
FROM (
    SELECT province_id, COUNT(*)::BIGINT AS ward_count
    FROM ward_v2
    GROUP BY province_id
) AS counts
WHERE province.id = counts.province_id;
