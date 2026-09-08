-- Add role_id and role_key columns to user_profile table
-- Migration: Add role support to user_profile

ALTER TABLE user_profile 
ADD COLUMN IF NOT EXISTS role_id BIGINT,
ADD COLUMN IF NOT EXISTS role_key INTEGER;

-- Add index for role_id for better query performance
CREATE INDEX IF NOT EXISTS idx_user_profile_role_id ON user_profile(role_id);

-- Add comment
COMMENT ON COLUMN user_profile.role_id IS 'ID của role trong auth-service';
COMMENT ON COLUMN user_profile.role_key IS 'Key của role để xác định nhanh quyền hạn';

