-- Migration: Thêm province_id và ward_id vào bảng user_profile
-- Date: 2025-11-04

-- Thêm cột province_id
ALTER TABLE user_profile 
ADD COLUMN IF NOT EXISTS province_id BIGINT DEFAULT NULL;

-- Thêm cột ward_id
ALTER TABLE user_profile 
ADD COLUMN IF NOT EXISTS ward_id BIGINT DEFAULT NULL;

-- Thêm comment cho các cột
COMMENT ON COLUMN user_profile.province_id IS 'ID của tỉnh/thành phố';
COMMENT ON COLUMN user_profile.ward_id IS 'ID của phường/xã';

-- Tạo index cho province_id (nếu cần query theo tỉnh/thành)
CREATE INDEX IF NOT EXISTS idx_user_profile_province_id ON user_profile(province_id);

-- Tạo index cho ward_id (nếu cần query theo phường/xã)
CREATE INDEX IF NOT EXISTS idx_user_profile_ward_id ON user_profile(ward_id);

