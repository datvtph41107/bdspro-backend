BEGIN;

CREATE TABLE banks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(192) NOT NULL,
    logo TEXT NOT NULL DEFAULT '',
    code VARCHAR(64) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT NOT NULL DEFAULT '',
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT banks_code_uq UNIQUE (code)
);

CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL CHECK (user_id > 0),
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency VARCHAR(8) NOT NULL CHECK (currency <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT wallets_user_uq UNIQUE (user_id)
);

CREATE TABLE payment_methods (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(192) NOT NULL,
    code VARCHAR(96) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT payment_methods_org_code_uq UNIQUE (organization_id, code)
);

CREATE TABLE wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    type VARCHAR(32) NOT NULL CHECK (type <> ''),
    amount BIGINT NOT NULL CHECK (amount <> 0),
    status VARCHAR(32) NOT NULL CHECK (status <> ''),
    related_service VARCHAR(128) NOT NULL DEFAULT '',
    related_id VARCHAR(192) NOT NULL DEFAULT '',
    external_payment_id VARCHAR(192) NOT NULL DEFAULT '',
    payment_method VARCHAR(96) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_by BIGINT NOT NULL DEFAULT 0,
    approved_by BIGINT NOT NULL DEFAULT 0,
    transaction_code VARCHAR(192) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT wallet_transactions_code_uq UNIQUE (transaction_code)
);
CREATE INDEX wallet_transactions_wallet_created_idx ON wallet_transactions(wallet_id, created_at DESC);
CREATE INDEX wallet_transactions_status_idx ON wallet_transactions(status, created_at DESC);

CREATE TABLE wallet_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    old_balance BIGINT NOT NULL CHECK (old_balance >= 0),
    new_balance BIGINT NOT NULL CHECK (new_balance >= 0),
    reason TEXT NOT NULL,
    changed_by BIGINT NOT NULL DEFAULT 0,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX wallet_audit_logs_wallet_idx ON wallet_audit_logs(wallet_id, changed_at DESC);

CREATE TABLE withdrawal_requests (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    status VARCHAR(32) NOT NULL CHECK (status <> ''),
    payment_method VARCHAR(96) NOT NULL DEFAULT '',
    bank_account_info TEXT NOT NULL DEFAULT '',
    requested_by BIGINT NOT NULL DEFAULT 0,
    approved_by BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX withdrawal_requests_wallet_status_idx ON withdrawal_requests(wallet_id, status, created_at DESC);

CREATE TABLE transaction_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL,
    category VARCHAR(96) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT transaction_types_code_uq UNIQUE (code)
);

CREATE TABLE dashboard_metrics (
    id BIGSERIAL PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL,
    processing_count BIGINT NOT NULL DEFAULT 0 CHECK (processing_count >= 0),
    completed_count BIGINT NOT NULL DEFAULT 0 CHECK (completed_count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT dashboard_metrics_time_uq UNIQUE (time)
);

COMMIT;
