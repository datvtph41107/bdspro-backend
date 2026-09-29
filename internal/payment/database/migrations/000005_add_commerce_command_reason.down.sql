BEGIN;

ALTER TABLE commerce_command_effects
    DROP COLUMN IF EXISTS reason;

COMMIT;
