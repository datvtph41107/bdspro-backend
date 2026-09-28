-- Bank accounts linked from deals.bank_account_id.
-- Idempotent for DBs that already absorbed bdspro dumps.

CREATE TABLE IF NOT EXISTS bank_accounts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    owner_id BIGINT NOT NULL DEFAULT 0,
    owner_type SMALLINT NOT NULL DEFAULT 10,
    bank_name VARCHAR(255) NOT NULL DEFAULT '',
    bank_id BIGINT NOT NULL DEFAULT 0,
    account_number VARCHAR(255) NOT NULL DEFAULT '',
    account_name VARCHAR(255) NOT NULL DEFAULT '',
    account_type INTEGER NOT NULL DEFAULT 1,
    account_status BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_deleted_at ON bank_accounts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_owner_id ON bank_accounts (owner_id);
