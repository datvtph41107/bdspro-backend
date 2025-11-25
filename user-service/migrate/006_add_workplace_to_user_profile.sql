-- Migration: Thêm workplace (nơi làm việc) vào bảng user_profile
-- Date: 2025-01-XX

-- Thêm cột workplace
ALTER TABLE user_profile 
ADD COLUMN IF NOT EXISTS workplace VARCHAR(255) DEFAULT NULL;

-- Thêm comment cho cột
COMMENT ON COLUMN user_profile.workplace IS 'Nơi làm việc của người dùng';

