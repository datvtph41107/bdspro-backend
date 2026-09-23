BEGIN;

CREATE TABLE IF NOT EXISTS demo_contact_log (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    contact_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_demo_contact_log_profile_contact UNIQUE (profile_id, contact_id)
);

CREATE INDEX IF NOT EXISTS idx_demo_contact_log_profile_id ON demo_contact_log(profile_id);

COMMIT;
