BEGIN;

-- Align user-admin permission primary keys with shared/common/domain/enum:
-- USER_ADMIN_VIEW=309, USER_ADMIN_MANAGE=310.

INSERT INTO permissions (id, name, key, description, permission_type, module, priority)
SELECT seed.id, seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    (309::bigint, 'Xem người dùng quản trị', 'USER_ADMIN_VIEW', 'Đọc danh sách, chi tiết và thống kê người dùng trong Admin', 'ORGANIZATION', 'user_admin', 10),
    (310::bigint, 'Quản lý người dùng quản trị', 'USER_ADMIN_MANAGE', 'Tạo, sửa, khóa, duyệt hoặc xóa người dùng trong Admin', 'ORGANIZATION', 'user_admin', 30)
) AS seed(id, name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (
    SELECT 1 FROM permissions existing
     WHERE (existing.key = seed.key OR existing.id = seed.id)
       AND existing.deleted_at IS NULL
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
              ('USER_ADMIN_VIEW', 309::bigint),
              ('USER_ADMIN_MANAGE', 310::bigint)
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
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM permissions), 310)
);

COMMIT;
