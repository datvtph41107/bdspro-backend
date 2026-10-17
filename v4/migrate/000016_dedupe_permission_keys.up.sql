BEGIN;

-- Auth permission catalog requires unique codes. Historical UI tree folders
-- reused the same key across section roots and nested nodes. Align with
-- shared/common/domain/enum:
--   TO_CHUC=14, TO_CHUC_2=38
--   CRM=259, CRM_2=41
--   BANG_TIN_NHOM=274, BANG_TIN_NHOM_2=154
-- Section roots 220/239 are not in the enum; give them distinct folder keys.

UPDATE permissions
   SET key = 'TO_CHUC_2', updated_at = NOW()
 WHERE id = 38
   AND key = 'TO_CHUC'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'CRM_2', updated_at = NOW()
 WHERE id = 41
   AND key = 'CRM'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'BANG_TIN_NHOM_2', updated_at = NOW()
 WHERE id = 154
   AND key = 'BANG_TIN_NHOM'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'TO_CHUC_SECTION', updated_at = NOW()
 WHERE id = 220
   AND key = 'TO_CHUC'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'CRM_SECTION', updated_at = NOW()
 WHERE id = 239
   AND key = 'CRM'
   AND deleted_at IS NULL;

-- 000007 could insert catalog-plan rows on auto-ids before 000015 pinned 302-304.
DO $$
DECLARE
    rec RECORD;
    occupant_key TEXT;
BEGIN
    FOR rec IN
        SELECT p.id AS legacy_id, seed.canon_id, seed.perm_key
          FROM permissions p
          JOIN (VALUES
              ('CATALOG_PLAN_VIEW', 302::bigint),
              ('CATALOG_PLAN_MANAGE', 303::bigint),
              ('CATALOG_PLAN_PUBLISH', 304::bigint)
          ) AS seed(perm_key, canon_id) ON seed.perm_key = p.key
         WHERE p.deleted_at IS NULL
           AND p.id <> seed.canon_id
    LOOP
        IF EXISTS (
            SELECT 1 FROM permissions
             WHERE id = rec.canon_id AND deleted_at IS NULL
        ) THEN
            INSERT INTO role_permissions (role_id, permission_id)
            SELECT rp.role_id, rec.canon_id
              FROM role_permissions rp
             WHERE rp.permission_id = rec.legacy_id
            ON CONFLICT DO NOTHING;

            DELETE FROM role_permissions WHERE permission_id = rec.legacy_id;

            INSERT INTO group_permissions (role_group_id, permission_id)
            SELECT gp.role_group_id, rec.canon_id
              FROM group_permissions gp
             WHERE gp.permission_id = rec.legacy_id
            ON CONFLICT DO NOTHING;

            DELETE FROM group_permissions WHERE permission_id = rec.legacy_id;

            UPDATE permissions
               SET deleted_at = NOW(), updated_at = NOW()
             WHERE id = rec.legacy_id;
            CONTINUE;
        END IF;

        SELECT key INTO occupant_key
          FROM permissions
         WHERE id = rec.canon_id
           AND deleted_at IS NULL
           AND key IS DISTINCT FROM rec.perm_key;

        IF occupant_key IS NOT NULL THEN
            RAISE EXCEPTION
                'cannot remap % to id %: id already owned by key %',
                rec.perm_key, rec.canon_id, occupant_key;
        END IF;

        DELETE FROM permissions
         WHERE id = rec.canon_id
           AND deleted_at IS NOT NULL;

        UPDATE role_permissions
           SET permission_id = rec.canon_id
         WHERE permission_id = rec.legacy_id;

        UPDATE group_permissions
           SET permission_id = rec.canon_id
         WHERE permission_id = rec.legacy_id;

        UPDATE permissions
           SET id = rec.canon_id, updated_at = NOW()
         WHERE id = rec.legacy_id;
    END LOOP;
END $$;

DO $$
DECLARE
    dup RECORD;
BEGIN
    SELECT key, array_agg(id ORDER BY id) AS ids
      INTO dup
      FROM permissions
     WHERE deleted_at IS NULL
       AND key IS NOT NULL
       AND key <> ''
     GROUP BY key
    HAVING COUNT(DISTINCT id) > 1
     LIMIT 1;

    IF FOUND THEN
        RAISE EXCEPTION
            'permission key still maps to multiple IDs after dedupe: % -> %',
            dup.key, dup.ids;
    END IF;
END $$;

COMMIT;
