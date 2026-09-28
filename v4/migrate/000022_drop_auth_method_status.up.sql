-- Drop credential-level account status/lock columns from auth_method.
-- Account lock/status truth is user_profile only.

BEGIN;

ALTER TABLE auth_method
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS locked_at,
    DROP COLUMN IF EXISTS locked_until,
    DROP COLUMN IF EXISTS lock_reason,
    DROP COLUMN IF EXISTS locked_by;

COMMIT;
