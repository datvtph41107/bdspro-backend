-- Migration: Thêm last_login_at và total_login vào bảng user_profile
-- Date: 2026-09-03
-- Login Create/First profile sẽ fail nếu thiếu các cột này.

ALTER TABLE user_profile
ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ DEFAULT NULL;

ALTER TABLE user_profile
ADD COLUMN IF NOT EXISTS total_login INTEGER NOT NULL DEFAULT 0;

COMMENT ON COLUMN user_profile.last_login_at IS 'Thời điểm đăng nhập gần nhất';
COMMENT ON COLUMN user_profile.total_login IS 'Tổng số lần đăng nhập';

CREATE INDEX IF NOT EXISTS idx_user_profile_last_login_at ON user_profile(last_login_at);
