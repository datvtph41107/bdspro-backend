-- Merge user_status + registration-facing user_info into auth_status (auth_id).
-- Post-profile account lock/status lives on user_profile.

BEGIN;

CREATE TABLE IF NOT EXISTS auth_status (
    auth_id BIGINT PRIMARY KEY,
    active BOOLEAN,
    verified BOOLEAN,
    locked_until TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

INSERT INTO auth_status (auth_id, active, verified, locked_until, updated_at)
SELECT
    us.auth_id,
    us.active,
    us.verified,
    us.locked_until,
    COALESCE(us.updated_at, NOW())
FROM user_status us
ON CONFLICT (auth_id) DO NOTHING;

ALTER TABLE user_profile
    ADD COLUMN IF NOT EXISTS locked_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS lock_reason VARCHAR(500),
    ADD COLUMN IF NOT EXISTS locked_by BIGINT;

-- Copy lock audit + align status from user_info onto user_profile.
UPDATE user_profile up
SET
    locked_at = ui.locked_at,
    locked_until = ui.locked_until,
    lock_reason = ui.lock_reason,
    locked_by = ui.locked_by,
    status = CASE
        WHEN ui.status IN (20, 30, 40) AND (up.status IS NULL OR up.status = 10 OR up.status = 0)
            THEN ui.status
        WHEN ui.status = 10 AND up.status IS NULL
            THEN 10
        ELSE up.status
    END,
    updated_at = NOW()
FROM user_info ui
WHERE ui.profile_id = up.profile_id
  AND ui.deleted_at IS NULL
  AND up.deleted_at IS NULL;

DROP TABLE IF EXISTS user_info CASCADE;
DROP TABLE IF EXISTS user_status CASCADE;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'user_owner') THEN
        GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE auth_status TO user_owner;
    END IF;
END $$;

COMMIT;
