-- Migration: Thêm cột device_name vào bảng user_session
-- Date: 2025-10-20

ALTER TABLE user_session 
ADD COLUMN IF NOT EXISTS device_name VARCHAR(255) DEFAULT '';

-- Index để tìm kiếm theo device_name
CREATE INDEX IF NOT EXISTS idx_user_session_device_name ON user_session(device_name);

COMMENT ON COLUMN user_session.device_name IS 'Tên thiết bị: iPhone 13 Pro, Samsung Galaxy S21, MacBook Pro, etc.';

