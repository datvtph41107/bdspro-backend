BEGIN;

-- Role and permission policy are IAM capabilities owned by User Service.
-- These canonical keys replace the undeclared ADMIN_VT_* runtime dependency
-- so a database created only from versioned migrations is operable.
INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT seed.name, seed.key, seed.description, 'SYSTEM', 'iam', seed.priority
FROM (VALUES
    ('Xem vai trò và quyền', 'IAM_ROLE_VIEW', 'Đọc vai trò và ma trận quyền quản trị', 10),
    ('Quản lý vai trò và quyền', 'IAM_ROLE_MANAGE', 'Tạo, sửa, xóa, gán vai trò và cập nhật ma trận quyền', 30)
) AS seed(name, key, description, priority)
WHERE NOT EXISTS (
    SELECT 1 FROM permissions permission
    WHERE permission.key = seed.key AND permission.deleted_at IS NULL
);

COMMIT;
