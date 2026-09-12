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
