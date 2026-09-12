BEGIN;

-- Seed catalog-plan permissions with shared enum IDs:
-- CATALOG_PLAN_VIEW=302, CATALOG_PLAN_MANAGE=303, CATALOG_PLAN_PUBLISH=304.

INSERT INTO permissions (id, name, key, description, permission_type, module, priority)
SELECT seed.id, seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    (302::bigint, 'Xem phiên bản catalog', 'CATALOG_PLAN_VIEW', 'Đọc plan version và điều khoản thương mại', 'ORGANIZATION', 'commercial', 10),
    (303::bigint, 'Quản lý bản nháp gói', 'CATALOG_PLAN_MANAGE', 'Tạo/sửa/xóa/validate draft plan version', 'ORGANIZATION', 'commercial', 20),
    (304::bigint, 'Xuất bản phiên bản catalog', 'CATALOG_PLAN_PUBLISH', 'Xuất bản draft plan version', 'ORGANIZATION', 'commercial', 20)
) AS seed(id, name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (
    SELECT 1 FROM permissions existing
     WHERE existing.key = seed.key AND existing.deleted_at IS NULL
);

SELECT setval(
    pg_get_serial_sequence('permissions', 'id'),
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM permissions), 304)
);

COMMIT;
