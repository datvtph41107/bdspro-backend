BEGIN;

-- Legacy db_payment instances already had wallet/bank tables before
-- 000004_wallet_and_bank. Bring their columns in line without recreating tables.

ALTER TABLE banks
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE wallet_transactions
    ADD COLUMN IF NOT EXISTS payment_method VARCHAR(96) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS created_by BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS approved_by BIGINT NOT NULL DEFAULT 0;

COMMIT;
