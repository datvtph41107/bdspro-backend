BEGIN;

CREATE TABLE IF NOT EXISTS demo_appointment_log (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    appointment_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_demo_appointment_log_profile_appointment UNIQUE (profile_id, appointment_id)
);

CREATE INDEX IF NOT EXISTS idx_demo_appointment_log_profile_id ON demo_appointment_log(profile_id);

CREATE TABLE IF NOT EXISTS demo_deal_log (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    deal_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_demo_deal_log_profile_deal UNIQUE (profile_id, deal_id)
);

CREATE INDEX IF NOT EXISTS idx_demo_deal_log_profile_id ON demo_deal_log(profile_id);

COMMIT;
