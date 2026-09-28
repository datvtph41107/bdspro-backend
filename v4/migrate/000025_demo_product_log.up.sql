BEGIN;

CREATE TABLE IF NOT EXISTS demo_product_log (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_demo_product_log_profile_product UNIQUE (profile_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_demo_product_log_profile_id ON demo_product_log(profile_id);

COMMIT;
