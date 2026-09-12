-- Move the versioned Vietnamese administrative reference catalog from the
-- historical tables into the tables owned by the current Location runtime.
-- This is a source-data cutover, not a development/demo fixture.
--
-- Legacy `provinces` rows on some workstations only populate `name` and leave
-- `full_name` / `code` empty. Prefer non-empty display names and skip rows that
-- cannot satisfy province_v2 constraints. Legacy `wards` may be absent when
-- ward_v2 is already the live catalog — skip that half in that case.

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
    COALESCE(NULLIF(btrim(full_name), ''), NULLIF(btrim(name), '')),
    COALESCE(NULLIF(btrim(short_name), ''), NULLIF(btrim(name), ''), ''),
    lat,
    lng,
    NULLIF(btrim(code::text), ''),
    COALESCE(created_at, CURRENT_TIMESTAMP),
    COALESCE(created_at, CURRENT_TIMESTAMP)
FROM provinces
WHERE COALESCE(NULLIF(btrim(full_name), ''), NULLIF(btrim(name), '')) IS NOT NULL
ON CONFLICT DO NOTHING;

-- Never overwrite a populated canonical name with a null/blank legacy value.
UPDATE province_v2 AS canonical
SET
    full_name = COALESCE(
        NULLIF(btrim(legacy.full_name), ''),
        NULLIF(btrim(legacy.name), ''),
        canonical.full_name
    ),
    short_name = COALESCE(
        NULLIF(btrim(legacy.short_name), ''),
        NULLIF(btrim(legacy.name), ''),
        canonical.short_name
    ),
    lat = COALESCE(legacy.lat, canonical.lat),
    lng = COALESCE(legacy.lng, canonical.lng),
    updated_at = CURRENT_TIMESTAMP
FROM provinces AS legacy
WHERE NULLIF(btrim(legacy.code::text), '') IS NOT NULL
  AND canonical.code = NULLIF(btrim(legacy.code::text), '');

DO $$
BEGIN
    IF to_regclass('public.wards') IS NULL THEN
        RAISE NOTICE 'legacy wards table missing; skipping ward cutover';
        RETURN;
    END IF;

    EXECUTE $ward$
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
            COALESCE(
                NULLIF(btrim(legacy_ward.full_name), ''),
                NULLIF(btrim(legacy_ward.name), '')
            ),
            COALESCE(
                NULLIF(btrim(legacy_ward.short_name), ''),
                NULLIF(btrim(legacy_ward.name), ''),
                ''
            ),
            legacy_ward.lat,
            legacy_ward.lng,
            NULLIF(btrim(legacy_ward.code::text), ''),
            canonical_province.id,
            canonical_province.code,
            COALESCE(legacy_ward.created_at, CURRENT_TIMESTAMP),
            COALESCE(legacy_ward.created_at, CURRENT_TIMESTAMP)
        FROM wards AS legacy_ward
        JOIN provinces AS legacy_province
            ON legacy_province.id = legacy_ward.province_id
        JOIN province_v2 AS canonical_province
            ON canonical_province.code = NULLIF(btrim(legacy_province.code::text), '')
        WHERE COALESCE(
                NULLIF(btrim(legacy_ward.full_name), ''),
                NULLIF(btrim(legacy_ward.name), '')
            ) IS NOT NULL
        ON CONFLICT DO NOTHING
    $ward$;

    EXECUTE $ward$
        UPDATE ward_v2 AS canonical
        SET
            full_name = COALESCE(
                NULLIF(btrim(legacy_ward.full_name), ''),
                NULLIF(btrim(legacy_ward.name), ''),
                canonical.full_name
            ),
            short_name = COALESCE(
                NULLIF(btrim(legacy_ward.short_name), ''),
                NULLIF(btrim(legacy_ward.name), ''),
                canonical.short_name
            ),
            lat = COALESCE(legacy_ward.lat, canonical.lat),
            lng = COALESCE(legacy_ward.lng, canonical.lng),
            province_id = canonical_province.id,
            province_code = canonical_province.code,
            updated_at = CURRENT_TIMESTAMP
        FROM wards AS legacy_ward
        JOIN provinces AS legacy_province
            ON legacy_province.id = legacy_ward.province_id
        JOIN province_v2 AS canonical_province
            ON canonical_province.code = NULLIF(btrim(legacy_province.code::text), '')
        WHERE NULLIF(btrim(legacy_ward.code::text), '') IS NOT NULL
          AND canonical.code = NULLIF(btrim(legacy_ward.code::text), '')
    $ward$;
END $$;

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
