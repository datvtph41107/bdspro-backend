-- Add login tracking columns expected by UserProfileEntity / CreateProfile.

BEGIN;

ALTER TABLE user_profile
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS total_login INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_user_profile_last_login_at ON user_profile(last_login_at);

COMMIT;
