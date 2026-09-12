BEGIN;

-- Align commercial permission primary keys with shared/common/domain/enum:
-- COMMERCIAL_USAGE_VIEW=299, COMMERCIAL_USAGE_RECONCILE=300,
-- COMMERCIAL_SUBSCRIPTION_VIEW=301 so Auth RequiredPermissions(id) matches role grants.

INSERT INTO permissions (id, name, key, description, permission_type, module, priority)
SELECT seed.id, seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    (299::bigint, 'Xem usage thương mại', 'COMMERCIAL_USAGE_VIEW', 'Đọc durable/runtime usage projection của TQD', 'ORGANIZATION', 'commercial', 10),
    (300::bigint, 'Reconcile usage thương mại', 'COMMERCIAL_USAGE_RECONCILE', 'Yêu cầu TQD đồng bộ runtime projection từ durable usage', 'ORGANIZATION', 'commercial', 20),
    (301::bigint, 'Xem subscription thương mại', 'COMMERCIAL_SUBSCRIPTION_VIEW', 'Đọc subscription và commercial projection của người dùng', 'ORGANIZATION', 'commercial', 10)
) AS seed(id, name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (
    SELECT 1 FROM permissions existing
     WHERE existing.key = seed.key AND existing.deleted_at IS NULL
);

DO $$
DECLARE
    rec RECORD;
    occupant_key TEXT;
BEGIN
    FOR rec IN
        SELECT p.id AS old_id, seed.target_id, seed.permission_key
          FROM permissions p
          JOIN (VALUES
              ('COMMERCIAL_USAGE_VIEW', 299::bigint),
              ('COMMERCIAL_USAGE_RECONCILE', 300::bigint),
              ('COMMERCIAL_SUBSCRIPTION_VIEW', 301::bigint)
          ) AS seed(permission_key, target_id) ON seed.permission_key = p.key
         WHERE p.deleted_at IS NULL
           AND p.id <> seed.target_id
    LOOP
        SELECT key INTO occupant_key
          FROM permissions
         WHERE id = rec.target_id
           AND deleted_at IS NULL
           AND key IS DISTINCT FROM rec.permission_key;

        IF occupant_key IS NOT NULL THEN
            RAISE EXCEPTION
                'cannot remap % to id %: id already owned by key %',
                rec.permission_key, rec.target_id, occupant_key;
        END IF;

        DELETE FROM permissions
         WHERE id = rec.target_id
           AND deleted_at IS NOT NULL;

        UPDATE role_permissions
           SET permission_id = rec.target_id
         WHERE permission_id = rec.old_id;

        UPDATE group_permissions
           SET permission_id = rec.target_id
         WHERE permission_id = rec.old_id;

        UPDATE permissions
           SET id = rec.target_id, updated_at = NOW()
         WHERE id = rec.old_id;
    END LOOP;
END $$;

SELECT setval(
    pg_get_serial_sequence('permissions', 'id'),
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM permissions), 301)
);

COMMIT;
