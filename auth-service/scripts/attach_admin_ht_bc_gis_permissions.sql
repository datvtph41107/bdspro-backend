-- Seed top-level HỖ TRỢ (HO_TRO) → Yêu cầu / Báo cáo → Xem,Tạo,Xử lý,Đóng
-- Idempotent. Reparents if previously under HE_THONG. Clears old *_GAN keys.
-- After run: flush Redis ROLE_PERMS_* + PERM_KEY_* (or restart auth) rồi login lại.

-- 0) Soft-delete obsolete Gán keys (replaced by Đóng / Xử lý)
UPDATE permissions
SET deleted_at = NOW(), updated_at = NOW()
WHERE key IN ('ADMIN_HT_GAN', 'ADMIN_BC_GIS_GAN')
  AND deleted_at IS NULL;

-- 0b) Rename Sửa → Xử lý, Xóa → Đóng.
-- Đổi key tại chỗ để giữ nguyên permission id, role nào đang được gán vẫn giữ quyền.
UPDATE permissions p
SET key = CASE p.key
      WHEN 'ADMIN_HT_SUA' THEN 'ADMIN_HT_XU_LY'
      WHEN 'ADMIN_HT_XOA' THEN 'ADMIN_HT_DONG'
      WHEN 'ADMIN_BC_GIS_SUA' THEN 'ADMIN_BC_GIS_XU_LY'
      WHEN 'ADMIN_BC_GIS_XOA' THEN 'ADMIN_BC_GIS_DONG'
    END,
    name = CASE
      WHEN p.key IN ('ADMIN_HT_SUA', 'ADMIN_BC_GIS_SUA') THEN 'Xử lý'
      ELSE 'Đóng'
    END,
    updated_at = NOW()
WHERE p.key IN ('ADMIN_HT_SUA', 'ADMIN_HT_XOA', 'ADMIN_BC_GIS_SUA', 'ADMIN_BC_GIS_XOA')
  AND NOT EXISTS (
    SELECT 1 FROM permissions p2
    WHERE p2.deleted_at IS NULL
      AND p2.key = CASE p.key
        WHEN 'ADMIN_HT_SUA' THEN 'ADMIN_HT_XU_LY'
        WHEN 'ADMIN_HT_XOA' THEN 'ADMIN_HT_DONG'
        WHEN 'ADMIN_BC_GIS_SUA' THEN 'ADMIN_BC_GIS_XU_LY'
        WHEN 'ADMIN_BC_GIS_XOA' THEN 'ADMIN_BC_GIS_DONG'
      END
  );

-- Nếu key mới đã tồn tại từ trước thì key cũ là bản trùng → soft-delete
UPDATE permissions
SET deleted_at = NOW(), updated_at = NOW()
WHERE key IN ('ADMIN_HT_SUA', 'ADMIN_HT_XOA', 'ADMIN_BC_GIS_SUA', 'ADMIN_BC_GIS_XOA')
  AND deleted_at IS NULL;

-- 1) Root HO_TRO
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT 'HỖ TRỢ', 'HO_TRO', '', 'ORGANIZATION', '', NULL, 30, '', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = 'HO_TRO' AND p.deleted_at IS NULL
);

-- Ensure parent_id is null (top-level tab); priority <= 30 để không hiện đỏ trên edit-role
UPDATE permissions
SET parent_id = NULL, name = 'HỖ TRỢ', priority = 30, updated_at = NOW(), deleted_at = NULL
WHERE key = 'HO_TRO';

-- 2) Mid tabs under HO_TRO
WITH root AS (
  SELECT id FROM permissions WHERE key = 'HO_TRO' AND deleted_at IS NULL LIMIT 1
)
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT v.name, v.key, '', 'ORGANIZATION', '', root.id, 20, '', NOW(), NOW()
FROM root
CROSS JOIN (VALUES
  ('Yêu cầu hỗ trợ', 'ADMIN_HT'),
  ('Báo cáo sai sót dữ liệu', 'ADMIN_BC_GIS')
) AS v(name, key)
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = v.key AND p.deleted_at IS NULL
);

-- Reparent existing mid tabs under HO_TRO
UPDATE permissions p
SET parent_id = (SELECT id FROM permissions WHERE key = 'HO_TRO' AND deleted_at IS NULL LIMIT 1),
    name = CASE
      WHEN p.key = 'ADMIN_HT' THEN 'Yêu cầu hỗ trợ'
      WHEN p.key = 'ADMIN_BC_GIS' THEN 'Báo cáo sai sót dữ liệu'
      ELSE p.name
    END,
    updated_at = NOW(),
    deleted_at = NULL
WHERE p.key IN ('ADMIN_HT', 'ADMIN_BC_GIS');

-- 3) Leaves under ADMIN_HT
WITH parent AS (
  SELECT id FROM permissions WHERE key = 'ADMIN_HT' AND deleted_at IS NULL LIMIT 1
)
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT v.name, v.key, '', 'ORGANIZATION', '', parent.id, 10, '', NOW(), NOW()
FROM parent
CROSS JOIN (VALUES
  ('Xem', 'ADMIN_HT_XEM'),
  ('Tạo', 'ADMIN_HT_TAO'),
  ('Xử lý', 'ADMIN_HT_XU_LY'),
  ('Đóng', 'ADMIN_HT_DONG')
) AS v(name, key)
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = v.key AND p.deleted_at IS NULL
);

UPDATE permissions p
SET parent_id = (SELECT id FROM permissions WHERE key = 'ADMIN_HT' AND deleted_at IS NULL LIMIT 1),
    name = CASE p.key
      WHEN 'ADMIN_HT_XEM' THEN 'Xem'
      WHEN 'ADMIN_HT_TAO' THEN 'Tạo'
      WHEN 'ADMIN_HT_XU_LY' THEN 'Xử lý'
      WHEN 'ADMIN_HT_DONG' THEN 'Đóng'
      ELSE p.name
    END,
    updated_at = NOW(),
    deleted_at = NULL
WHERE p.key IN ('ADMIN_HT_XEM', 'ADMIN_HT_TAO', 'ADMIN_HT_XU_LY', 'ADMIN_HT_DONG');

-- 4) Leaves under ADMIN_BC_GIS
WITH parent AS (
  SELECT id FROM permissions WHERE key = 'ADMIN_BC_GIS' AND deleted_at IS NULL LIMIT 1
)
INSERT INTO permissions (name, key, description, permission_type, module, parent_id, priority, required_key, created_at, updated_at)
SELECT v.name, v.key, '', 'ORGANIZATION', '', parent.id, 10, '', NOW(), NOW()
FROM parent
CROSS JOIN (VALUES
  ('Xem', 'ADMIN_BC_GIS_XEM'),
  ('Tạo', 'ADMIN_BC_GIS_TAO'),
  ('Xử lý', 'ADMIN_BC_GIS_XU_LY'),
  ('Đóng', 'ADMIN_BC_GIS_DONG')
) AS v(name, key)
WHERE NOT EXISTS (
  SELECT 1 FROM permissions p WHERE p.key = v.key AND p.deleted_at IS NULL
);

UPDATE permissions p
SET parent_id = (SELECT id FROM permissions WHERE key = 'ADMIN_BC_GIS' AND deleted_at IS NULL LIMIT 1),
    name = CASE p.key
      WHEN 'ADMIN_BC_GIS_XEM' THEN 'Xem'
      WHEN 'ADMIN_BC_GIS_TAO' THEN 'Tạo'
      WHEN 'ADMIN_BC_GIS_XU_LY' THEN 'Xử lý'
      WHEN 'ADMIN_BC_GIS_DONG' THEN 'Đóng'
      ELSE p.name
    END,
    updated_at = NOW(),
    deleted_at = NULL
WHERE p.key IN ('ADMIN_BC_GIS_XEM', 'ADMIN_BC_GIS_TAO', 'ADMIN_BC_GIS_XU_LY', 'ADMIN_BC_GIS_DONG');

-- 5) Attach to SUPER_ADMIN
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.deleted_at IS NULL
  AND p.deleted_at IS NULL
  AND (r.key = 'ROLE_SUPER_ADMIN' OR r.role_key = 30)
  AND p.key IN (
    'HO_TRO',
    'ADMIN_HT', 'ADMIN_HT_XEM', 'ADMIN_HT_TAO', 'ADMIN_HT_XU_LY', 'ADMIN_HT_DONG',
    'ADMIN_BC_GIS', 'ADMIN_BC_GIS_XEM', 'ADMIN_BC_GIS_TAO', 'ADMIN_BC_GIS_XU_LY', 'ADMIN_BC_GIS_DONG'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- Drop obsolete Gán from role_permissions
DELETE FROM role_permissions rp
USING permissions p
WHERE rp.permission_id = p.id
  AND p.key IN ('ADMIN_HT_GAN', 'ADMIN_BC_GIS_GAN');
