-- Migration: Thêm role_title (chức danh) vào bảng user_profile
-- Date: 2025-01-XX

-- Thêm cột role_title
ALTER TABLE user_profile 
ADD COLUMN IF NOT EXISTS role_title VARCHAR(255) DEFAULT NULL;

-- Thêm comment cho cột
COMMENT ON COLUMN user_profile.role_title IS 'Chức danh của người dùng';

