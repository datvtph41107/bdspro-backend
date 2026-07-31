-- Seed ADMIN_KD (Quản lý kinh doanh) → Xem, Tạo, Sửa, Xóa
-- Cây 3 cấp như HO_TRO (edit-role UI cần mid tab có children):
--   KINH_DOANH (root)
--     ADMIN_KD (mid)
--       ADMIN_KD_XEM / TAO / SUA / XOA (leaves)
-- Idempotent. After run: flush Redis ROLE_PERMS_* rồi login lại.

-- 1) Root KINH_DOANH
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT 'KINH DOANH', 'KINH_DOANH', '', 'ORGANIZATION', '', NULL, 30, '', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = 'KINH_DOANH' AND p.deleted_at IS NULL
);

UPDATE permissions
SET parent_id = NULL, name = 'KINH DOANH', priority = 30, updated_at = NOW(), deleted_at = NULL
WHERE key = 'KINH_DOANH';

-- 2) Mid tab ADMIN_KD under KINH_DOANH
WITH root AS (
  SELECT id FROM permissions WHERE key = 'KINH_DOANH' AND deleted_at IS NULL LIMIT 1
)
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT 'Quản lý kinh doanh', 'ADMIN_KD', '', 'ORGANIZATION', '', root.id, 20, '', NOW(), NOW()
FROM root
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = 'ADMIN_KD' AND p.deleted_at IS NULL
);

UPDATE permissions p
SET parent_id = (SELECT id FROM permissions WHERE key = 'KINH_DOANH' AND deleted_at IS NULL LIMIT 1),
    name = 'Quản lý kinh doanh',
    priority = 20,
    updated_at = NOW(),
    deleted_at = NULL
WHERE p.key = 'ADMIN_KD';

-- 3) Leaves under ADMIN_KD
WITH parent AS (
  SELECT id FROM permissions WHERE key = 'ADMIN_KD' AND deleted_at IS NULL LIMIT 1
)
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT v.name, v.key, '', 'ORGANIZATION', '', parent.id, 10, '', NOW(), NOW()
FROM parent
CROSS JOIN (VALUES
  ('Xem', 'ADMIN_KD_XEM'),
  ('Tạo', 'ADMIN_KD_TAO'),
  ('Sửa', 'ADMIN_KD_SUA'),
  ('Xóa', 'ADMIN_KD_XOA')
) AS v(name, key)
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = v.key AND p.deleted_at IS NULL
);

UPDATE permissions p
SET parent_id = (SELECT id FROM permissions WHERE key = 'ADMIN_KD' AND deleted_at IS NULL LIMIT 1),
    name = CASE p.key
      WHEN 'ADMIN_KD_XEM' THEN 'Xem'
      WHEN 'ADMIN_KD_TAO' THEN 'Tạo'
      WHEN 'ADMIN_KD_SUA' THEN 'Sửa'
      WHEN 'ADMIN_KD_XOA' THEN 'Xóa'
      ELSE p.name
    END,
    priority = 10,
    updated_at = NOW(),
    deleted_at = NULL
WHERE p.key IN ('ADMIN_KD_XEM', 'ADMIN_KD_TAO', 'ADMIN_KD_SUA', 'ADMIN_KD_XOA');

-- 4) Attach to SUPER_ADMIN
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.deleted_at IS NULL
  AND p.deleted_at IS NULL
  AND (r.key = 'ROLE_SUPER_ADMIN' OR r.role_key = 30)
  AND p.key IN (
    'KINH_DOANH', 'ADMIN_KD', 'ADMIN_KD_XEM', 'ADMIN_KD_TAO', 'ADMIN_KD_SUA', 'ADMIN_KD_XOA'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );
