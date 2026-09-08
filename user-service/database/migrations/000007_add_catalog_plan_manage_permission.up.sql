BEGIN;

-- Draft lifecycle is separate from publication so an operator can delegate
-- data entry without granting the authority to make terms sellable.
INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT 'Quản lý bản nháp gói', 'CATALOG_PLAN_MANAGE',
       'Tạo, sửa, kiểm tra và xóa phiên bản gói chưa xuất bản',
       'ORGANIZATION', 'commercial', 20
WHERE NOT EXISTS (
    SELECT 1 FROM permissions WHERE key = 'CATALOG_PLAN_MANAGE' AND deleted_at IS NULL
);

COMMIT;
