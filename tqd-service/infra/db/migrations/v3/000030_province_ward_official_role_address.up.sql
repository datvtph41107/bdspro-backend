-- Bổ sung official_role, official_address cho province_v2 và ward_v2

-- province_v2
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_role    VARCHAR(500) NOT NULL DEFAULT '';
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS official_address VARCHAR(500) NOT NULL DEFAULT '';

COMMENT ON COLUMN province_v2.official_role    IS 'Vai trò chính thức (ví dụ: Tỉnh ủy, UBND tỉnh...)';
COMMENT ON COLUMN province_v2.official_address IS 'Địa chỉ cơ quan chính quyền cấp tỉnh';

-- ward_v2
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_role    VARCHAR(500) NOT NULL DEFAULT '';
ALTER TABLE ward_v2 ADD COLUMN IF NOT EXISTS official_address VARCHAR(500) NOT NULL DEFAULT '';

COMMENT ON COLUMN ward_v2.official_role    IS 'Vai trò chính thức (ví dụ: Đảng ủy phường, UBND phường...)';
COMMENT ON COLUMN ward_v2.official_address IS 'Địa chỉ cơ quan chính quyền cấp phường/xã';
