BEGIN;

-- Restore pre-dedupe UI tree keys (lossy only if rows were renamed again later).

UPDATE permissions
   SET key = 'TO_CHUC', updated_at = NOW()
 WHERE id = 38
   AND key = 'TO_CHUC_2'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'CRM', updated_at = NOW()
 WHERE id = 41
   AND key = 'CRM_2'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'BANG_TIN_NHOM', updated_at = NOW()
 WHERE id = 154
   AND key = 'BANG_TIN_NHOM_2'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'TO_CHUC', updated_at = NOW()
 WHERE id = 220
   AND key = 'TO_CHUC_SECTION'
   AND deleted_at IS NULL;

UPDATE permissions
   SET key = 'CRM', updated_at = NOW()
 WHERE id = 239
   AND key = 'CRM_SECTION'
   AND deleted_at IS NULL;

COMMIT;
