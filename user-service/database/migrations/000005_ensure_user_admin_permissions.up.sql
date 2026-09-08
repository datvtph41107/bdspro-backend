BEGIN;

-- User administration is a User/IAM capability. These permission identities
-- replace the historical profile-id/role-name shortcut and are deliberately
-- separate from commercial projection permissions.
INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    ('Xem người dùng quản trị', 'USER_ADMIN_VIEW', 'Đọc danh sách, chi tiết và thống kê người dùng trong Admin', 'ORGANIZATION', 'user_admin', 10),
    ('Quản lý người dùng quản trị', 'USER_ADMIN_MANAGE', 'Tạo, sửa, khóa, duyệt hoặc xóa người dùng trong Admin', 'ORGANIZATION', 'user_admin', 30)
) AS seed(name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (SELECT 1 FROM permissions existing WHERE existing.key = seed.key);

COMMIT;
