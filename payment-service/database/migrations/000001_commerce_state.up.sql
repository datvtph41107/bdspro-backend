BEGIN;

CREATE TABLE commerce_orders (
    id BIGSERIAL PRIMARY KEY,
    subject_kind VARCHAR(32) NOT NULL CHECK (subject_kind IN ('profile', 'organization')),
    subject_id VARCHAR(128) NOT NULL CHECK (subject_id <> ''),
    product_code VARCHAR(128) NOT NULL CHECK (product_code <> ''),
    plan_code VARCHAR(128) NOT NULL CHECK (plan_code <> ''),
    plan_version_id BIGINT NOT NULL CHECK (plan_version_id > 0),
    plan_version VARCHAR(64) NOT NULL CHECK (plan_version <> ''),
    tier_rank INTEGER NOT NULL CHECK (tier_rank > 0),
    subscription_term_days INTEGER NOT NULL CHECK (subscription_term_days > 0),
    terms_checksum CHAR(64) NOT NULL CHECK (terms_checksum ~ '^[0-9a-fA-F]{64}$'),
    currency VARCHAR(8) NOT NULL CHECK (currency <> ''),
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    command_key VARCHAR(128) NOT NULL CHECK (command_key <> ''),
    command_fingerprint CHAR(64) NOT NULL CHECK (command_fingerprint ~ '^[0-9a-f]{64}$'),
    reference VARCHAR(64) NOT NULL CHECK (reference <> ''),
    status VARCHAR(32) NOT NULL CHECK (status IN ('pending_funds', 'funds_confirmed', 'requires_review')),
    expires_at TIMESTAMPTZ NOT NULL,
    funds_confirmed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT commerce_orders_command_scope_uq
        UNIQUE (subject_kind, subject_id, product_code, command_key),
    CONSTRAINT commerce_orders_reference_uq UNIQUE (reference)
);

CREATE TABLE commerce_settlements (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NULL REFERENCES commerce_orders(id),
    provider VARCHAR(64) NOT NULL CHECK (provider <> ''),
    provider_transaction_id VARCHAR(192) NOT NULL CHECK (provider_transaction_id <> ''),
    reference VARCHAR(128) NOT NULL CHECK (reference <> ''),
    currency VARCHAR(8) NOT NULL CHECK (currency <> ''),
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    occurred_at TIMESTAMPTZ NOT NULL,
    evidence_hash CHAR(64) NOT NULL CHECK (evidence_hash ~ '^[0-9a-f]{64}$'),
    status VARCHAR(32) NOT NULL CHECK (status IN ('funds_confirmed', 'unmatched', 'requires_review', 'late')),
    review_reason TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT commerce_settlements_provider_tx_uq
        UNIQUE (provider, provider_transaction_id)
);

CREATE INDEX commerce_settlements_order_idx ON commerce_settlements(order_id);

CREATE TABLE commerce_fulfillments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES commerce_orders(id),
    status VARCHAR(32) NOT NULL CHECK (status IN ('pending', 'running', 'completed', 'retry', 'requires_review')),
    attempt_count INTEGER NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    available_at TIMESTAMPTZ NOT NULL,
    locked_by VARCHAR(128) NOT NULL DEFAULT '',
    claim_version BIGINT NOT NULL DEFAULT 0 CHECK (claim_version >= 0),
    lease_until TIMESTAMPTZ NULL,
    last_error TEXT NOT NULL DEFAULT '',
    completed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT commerce_fulfillments_order_uq UNIQUE (order_id)
);

CREATE INDEX commerce_fulfillments_claim_idx
    ON commerce_fulfillments(status, available_at, lease_until, id);

CREATE TABLE commerce_command_effects (
    id BIGSERIAL PRIMARY KEY,
    effect_type VARCHAR(96) NOT NULL CHECK (effect_type <> ''),
    scope_id BIGINT NOT NULL CHECK (scope_id > 0),
    command_key VARCHAR(128) NOT NULL CHECK (command_key <> ''),
    actor_id VARCHAR(128) NOT NULL CHECK (actor_id <> ''),
    reason TEXT NOT NULL DEFAULT '',
    outcome VARCHAR(32) NOT NULL DEFAULT 'accepted',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT commerce_command_effects_scope_uq
        UNIQUE (effect_type, scope_id, command_key)
);

COMMIT;
