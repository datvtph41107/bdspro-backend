BEGIN;

ALTER TABLE wallet_transactions
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS payment_method;

ALTER TABLE banks
    DROP COLUMN IF EXISTS is_deleted;

COMMIT;
