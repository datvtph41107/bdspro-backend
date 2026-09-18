BEGIN;

DROP INDEX IF EXISTS idx_user_profile_last_login_at;

ALTER TABLE user_profile
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS total_login;

COMMIT;
