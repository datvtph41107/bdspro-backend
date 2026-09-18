BEGIN;

CREATE TABLE IF NOT EXISTS user_status (
    auth_id BIGINT PRIMARY KEY,
    active BOOLEAN,
    verified BOOLEAN,
    locked_until TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

INSERT INTO user_status (auth_id, active, verified, locked_until, updated_at)
SELECT auth_id, active, verified, locked_until, updated_at
FROM auth_status
ON CONFLICT (auth_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS user_info (
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    profile_id BIGINT PRIMARY KEY,
    status SMALLINT DEFAULT 10,
    locked_at TIMESTAMPTZ,
    locked_until TIMESTAMPTZ,
    lock_reason VARCHAR(500),
    locked_by BIGINT,
    role_id BIGINT
);

INSERT INTO user_info (
    profile_id, status, locked_at, locked_until, lock_reason, locked_by, created_at, updated_at
)
SELECT
    profile_id,
    COALESCE(status, 10),
    locked_at,
    locked_until,
    COALESCE(lock_reason, ''),
    locked_by,
    COALESCE(created_at, NOW()),
    COALESCE(updated_at, NOW())
FROM user_profile
WHERE deleted_at IS NULL
ON CONFLICT (profile_id) DO NOTHING;

DROP TABLE IF EXISTS auth_status CASCADE;

ALTER TABLE user_profile
    DROP COLUMN IF EXISTS locked_at,
    DROP COLUMN IF EXISTS locked_until,
    DROP COLUMN IF EXISTS lock_reason,
    DROP COLUMN IF EXISTS locked_by;

COMMIT;
