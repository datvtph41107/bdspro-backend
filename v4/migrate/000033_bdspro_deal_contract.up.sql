-- Deal contract process tables for GET /v4/transaction/contract/deal/{id}/process.
-- Idempotent for DBs that already absorbed bdspro dumps.

CREATE TABLE IF NOT EXISTS tx_contract_deals (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deal_id BIGINT NOT NULL,
    type SMALLINT NOT NULL DEFAULT 10,
    note TEXT,
    product_id BIGINT,
    transaction_id BIGINT,
    customer_id BIGINT NOT NULL DEFAULT 0,
    approved_by BIGINT,
    approved_at TIMESTAMPTZ,
    reject_reason TEXT
);
CREATE INDEX IF NOT EXISTS idx_tx_contract_deals_deal_id ON tx_contract_deals (deal_id);
CREATE INDEX IF NOT EXISTS idx_tx_contract_deals_transaction_id ON tx_contract_deals (transaction_id);
CREATE INDEX IF NOT EXISTS idx_tx_contract_deals_deleted_at ON tx_contract_deals (deleted_at);

CREATE TABLE IF NOT EXISTS tx_transaction_action (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    tx_id BIGINT NOT NULL DEFAULT 0,
    from_id BIGINT NOT NULL DEFAULT 0,
    from_of SMALLINT NOT NULL DEFAULT 10,
    to_id BIGINT NOT NULL DEFAULT 0,
    to_of SMALLINT NOT NULL DEFAULT 10,
    action SMALLINT NOT NULL DEFAULT 10,
    value DOUBLE PRECISION NOT NULL DEFAULT 0,
    note TEXT,
    message TEXT,
    timestamp TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_tx_transaction_action_tx_id ON tx_transaction_action (tx_id);
CREATE INDEX IF NOT EXISTS idx_tx_transaction_action_deleted_at ON tx_transaction_action (deleted_at);
