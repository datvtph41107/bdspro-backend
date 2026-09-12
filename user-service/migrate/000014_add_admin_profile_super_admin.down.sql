BEGIN;

ALTER TABLE admin_profiles
    DROP CONSTRAINT IF EXISTS chk_admin_profiles_super_admin;

ALTER TABLE admin_profiles
    DROP COLUMN IF EXISTS super_admin;

COMMIT;
