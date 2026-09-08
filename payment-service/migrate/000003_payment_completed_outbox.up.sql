BEGIN;

CREATE TABLE commerce_outbox_events (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(192) NOT NULL CHECK (event_id <> ''),
    event_type VARCHAR(128) NOT NULL CHECK (event_type <> ''),
    schema_version INTEGER NOT NULL CHECK (schema_version > 0),
    routing_key VARCHAR(128) NOT NULL CHECK (routing_key <> ''),
    payload JSONB NOT NULL,
    status VARCHAR(32) NOT NULL CHECK (status IN ('pending', 'running', 'published', 'retry')),
    attempt_count INTEGER NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    available_at TIMESTAMPTZ NOT NULL,
    locked_by VARCHAR(128) NOT NULL DEFAULT '',
    claim_version BIGINT NOT NULL DEFAULT 0 CHECK (claim_version >= 0),
    lease_until TIMESTAMPTZ NULL,
    last_error TEXT NOT NULL DEFAULT '',
    published_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT commerce_outbox_events_event_id_uq UNIQUE (event_id)
);

CREATE INDEX commerce_outbox_events_claim_idx
    ON commerce_outbox_events(status, available_at, lease_until, id);

COMMIT;
