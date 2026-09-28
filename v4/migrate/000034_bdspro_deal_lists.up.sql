-- Supporting tables for deal list APIs (contracts, investments, costs).
-- Idempotent for DBs that already absorbed bdspro dumps.

CREATE TABLE IF NOT EXISTS tx_transaction (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    transaction_name VARCHAR(255),
    from_id BIGINT NOT NULL DEFAULT 0,
    from_of SMALLINT NOT NULL DEFAULT 10,
    to_id BIGINT NOT NULL DEFAULT 0,
    to_of SMALLINT NOT NULL DEFAULT 10,
    method SMALLINT NOT NULL DEFAULT 10,
    status SMALLINT NOT NULL DEFAULT 10,
    timestamp TIMESTAMPTZ,
    last_action_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_tx_transaction_deleted_at ON tx_transaction (deleted_at);

CREATE TABLE IF NOT EXISTS deal_investments (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    deal_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL DEFAULT 0,
    invest_type SMALLINT NOT NULL DEFAULT 10,
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    transfer_time TIMESTAMPTZ,
    note TEXT,
    transfer_proof_image TEXT,
    status SMALLINT NOT NULL DEFAULT 10,
    changed_at TIMESTAMPTZ,
    confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    proxy_user_id BIGINT
);
CREATE INDEX IF NOT EXISTS idx_deal_investments_deal_id ON deal_investments (deal_id);
CREATE INDEX IF NOT EXISTS idx_deal_investments_member_id ON deal_investments (member_id);
CREATE INDEX IF NOT EXISTS idx_deal_investments_deleted_at ON deal_investments (deleted_at);

CREATE TABLE IF NOT EXISTS tx_cost_types (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    type_name VARCHAR(255),
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    cost_type SMALLINT NOT NULL DEFAULT 10,
    description TEXT
);
CREATE INDEX IF NOT EXISTS idx_tx_cost_types_deleted_at ON tx_cost_types (deleted_at);

CREATE TABLE IF NOT EXISTS tx_deal_costs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    cost_name VARCHAR(255),
    deal_id BIGINT NOT NULL,
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    confirmed_at TIMESTAMPTZ,
    cost_type_id BIGINT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 10,
    note TEXT
);
CREATE INDEX IF NOT EXISTS idx_tx_deal_costs_deal_id ON tx_deal_costs (deal_id);
CREATE INDEX IF NOT EXISTS idx_tx_deal_costs_cost_type_id ON tx_deal_costs (cost_type_id);
CREATE INDEX IF NOT EXISTS idx_tx_deal_costs_deleted_at ON tx_deal_costs (deleted_at);
