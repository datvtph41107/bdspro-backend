BEGIN;

CREATE TABLE commerce_payment_attempts (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES commerce_orders(id),
    command_key VARCHAR(128) NOT NULL CHECK (command_key <> ''),
    command_fingerprint CHAR(64) NOT NULL CHECK (command_fingerprint ~ '^[0-9a-f]{64}$'),
    method VARCHAR(32) NOT NULL CHECK (method IN ('card', 'bank_transfer', 'qr')),
    provider VARCHAR(64) NOT NULL CHECK (provider <> ''),
    provider_reference VARCHAR(192) NOT NULL CHECK (provider_reference <> ''),
    status VARCHAR(32) NOT NULL CHECK (status IN (
        'created', 'pending_action', 'processing', 'authorized', 'succeeded',
        'declined', 'expired', 'canceled', 'failed'
    )),
    next_action_kind VARCHAR(32) NOT NULL DEFAULT 'none' CHECK (next_action_kind IN ('none', 'redirect', 'display_qr', 'wait')),
    redirect_url TEXT NOT NULL DEFAULT '',
    qr_payload TEXT NOT NULL DEFAULT '',
    action_expires_at TIMESTAMPTZ NULL,
    expires_at TIMESTAMPTZ NULL,
    failure_code VARCHAR(96) NOT NULL DEFAULT '',
    failure_detail TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT commerce_payment_attempts_command_uq UNIQUE (order_id, command_key),
    CONSTRAINT commerce_payment_attempts_provider_ref_uq UNIQUE (provider, provider_reference)
);

CREATE INDEX commerce_payment_attempts_order_idx
    ON commerce_payment_attempts(order_id, created_at DESC);

COMMIT;
