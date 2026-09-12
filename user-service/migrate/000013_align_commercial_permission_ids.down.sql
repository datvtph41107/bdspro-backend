BEGIN;

-- Down migration leaves rows in place (IAM grants must not be deleted) but does
-- not attempt to restore prior auto-assigned ids.

COMMIT;
