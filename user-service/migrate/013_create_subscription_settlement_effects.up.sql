BEGIN;

-- User owns the durable outcome of applying one Payment settlement to
-- Subscription. Rejected same-tier/downgrade outcomes are recorded too, so a
-- retried EffectKey cannot become live later after unrelated state changes.
CREATE TABLE IF NOT EXISTS commercial_subscription_effects (
    effect_key VARCHAR(128) PRIMARY KEY,
    fingerprint CHAR(64) NOT NULL,
    order_id BIGINT NOT NULL UNIQUE,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    product_code VARCHAR(128) NOT NULL,
    plan_code VARCHAR(128) NOT NULL,
    plan_version_id BIGINT NOT NULL REFERENCES catalog_plan_versions(id),
    plan_version VARCHAR(32) NOT NULL,
    tier_rank INTEGER NOT NULL,
    subscription_term_days INTEGER NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    action VARCHAR(32) NOT NULL,
    subscription_id BIGINT REFERENCES catalog_subscriptions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_commercial_subscription_effect_key
        CHECK (effect_key ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_commercial_subscription_effect_fingerprint
        CHECK (fingerprint ~ '^[0-9a-f]{64}$'),
    CONSTRAINT chk_commercial_subscription_effect_subject_kind
        CHECK (subject_kind IN ('profile', 'organization')),
    CONSTRAINT chk_commercial_subscription_effect_subject_id
        CHECK (subject_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_commercial_subscription_effect_tier
        CHECK (tier_rank > 0),
    CONSTRAINT chk_commercial_subscription_effect_term
        CHECK (subscription_term_days > 0),
    CONSTRAINT chk_commercial_subscription_effect_checksum
        CHECK (terms_checksum ~ '^[0-9a-f]{64}$'),
    CONSTRAINT chk_commercial_subscription_effect_action
        CHECK (action IN ('activated', 'upgraded', 'rejected_same_tier', 'rejected_downgrade'))
);

CREATE INDEX IF NOT EXISTS idx_commercial_subscription_effect_subject
    ON commercial_subscription_effects(subject_kind, subject_id, product_code, created_at DESC);

COMMIT;
