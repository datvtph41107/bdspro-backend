BEGIN;

ALTER TABLE admin_profiles
    ADD COLUMN IF NOT EXISTS super_admin SMALLINT NOT NULL DEFAULT 0;

ALTER TABLE admin_profiles
    DROP CONSTRAINT IF EXISTS chk_admin_profiles_super_admin;

ALTER TABLE admin_profiles
    ADD CONSTRAINT chk_admin_profiles_super_admin
    CHECK (super_admin IN (0, 1));

COMMIT;
