BEGIN;

-- Commercial ranking/term fields are deliberately nullable for existing rows.
-- Migration must not invent business truth for legacy plans. Checkout fails
-- closed until reviewed terms exist. A published legacy PlanVersion with no term
-- must receive a reviewed successor version rather than an in-place guessed term.
ALTER TABLE catalog_plans
    ADD COLUMN IF NOT EXISTS tier_rank INTEGER;

ALTER TABLE catalog_plans
    DROP CONSTRAINT IF EXISTS chk_catalog_plans_tier_rank;
ALTER TABLE catalog_plans
    ADD CONSTRAINT chk_catalog_plans_tier_rank
        CHECK (tier_rank IS NULL OR tier_rank > 0);

ALTER TABLE catalog_plan_versions
    ADD COLUMN IF NOT EXISTS subscription_term_days INTEGER;

ALTER TABLE catalog_plan_versions
    DROP CONSTRAINT IF EXISTS chk_catalog_plan_versions_subscription_term;
ALTER TABLE catalog_plan_versions
    ADD CONSTRAINT chk_catalog_plan_versions_subscription_term
        CHECK (subscription_term_days IS NULL OR subscription_term_days > 0);

ALTER TABLE catalog_subscriptions
    ADD COLUMN IF NOT EXISTS pending_plan_version_id BIGINT
        REFERENCES catalog_plan_versions(id);
ALTER TABLE catalog_subscriptions
    ADD COLUMN IF NOT EXISTS pending_effective_at TIMESTAMPTZ;

ALTER TABLE catalog_subscriptions
    DROP CONSTRAINT IF EXISTS chk_catalog_subscriptions_pending_change;
ALTER TABLE catalog_subscriptions
    ADD CONSTRAINT chk_catalog_subscriptions_pending_change
        CHECK (
            (pending_plan_version_id IS NULL AND pending_effective_at IS NULL)
            OR
            (pending_plan_version_id IS NOT NULL AND pending_effective_at IS NOT NULL)
        );

CREATE INDEX IF NOT EXISTS idx_catalog_subscriptions_pending_plan
    ON catalog_subscriptions(pending_plan_version_id)
    WHERE pending_plan_version_id IS NOT NULL;

-- User owns the durable resolution of one Checkout command. This snapshot is
-- written before the Payment call so a network retry does not re-resolve a
-- successor PlanVersion or reclassify an already-changed Subscription.
CREATE TABLE IF NOT EXISTS commercial_checkout_commands (
    id BIGSERIAL PRIMARY KEY,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    command_key VARCHAR(128) NOT NULL,
    product_id BIGINT NOT NULL REFERENCES catalog_products(id),
    product_code VARCHAR(128) NOT NULL,
    plan_id BIGINT NOT NULL REFERENCES catalog_plans(id),
    plan_code VARCHAR(128) NOT NULL,
    plan_version_id BIGINT NOT NULL REFERENCES catalog_plan_versions(id),
    plan_version VARCHAR(32) NOT NULL,
    tier_rank INTEGER NOT NULL,
    subscription_term_days INTEGER NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    subject_scope VARCHAR(32) NOT NULL,
    currency CHAR(3) NOT NULL,
    amount_minor BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_commercial_checkout_command
        UNIQUE (subject_kind, subject_id, command_key),
    CONSTRAINT chk_commercial_checkout_subject_kind
        CHECK (subject_kind IN ('profile', 'organization')),
    CONSTRAINT chk_commercial_checkout_subject_id
        CHECK (subject_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_commercial_checkout_command_key
        CHECK (command_key ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_commercial_checkout_codes
        CHECK (
            product_code ~ '^[a-z][a-z0-9_]*(\.[a-z0-9_]+)*$'
            AND plan_code ~ '^[a-z][a-z0-9_]*(\.[a-z0-9_]+)*$'
            AND plan_version ~ '^[0-9]+\.[0-9]+\.[0-9]+$'
        ),
    CONSTRAINT chk_commercial_checkout_tier
        CHECK (tier_rank > 0),
    CONSTRAINT chk_commercial_checkout_term
        CHECK (subscription_term_days > 0),
    CONSTRAINT chk_commercial_checkout_checksum
        CHECK (terms_checksum ~ '^[0-9a-f]{64}$'),
    CONSTRAINT chk_commercial_checkout_scope
        CHECK (subject_scope IN ('profile', 'organization', 'any')),
    CONSTRAINT chk_commercial_checkout_currency
        CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT chk_commercial_checkout_amount
        CHECK (amount_minor > 0)
);

CREATE INDEX IF NOT EXISTS idx_commercial_checkout_plan_version
    ON commercial_checkout_commands(plan_version_id);

-- Existing published rows with NULL term remain blocked for Checkout. This
-- trigger prevents an ordinary application update from silently rewriting the
-- newly introduced immutable term on any published PlanVersion.
CREATE OR REPLACE FUNCTION qhpro_catalog_protect_subscription_term()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.published_at IS NOT NULL
       AND NEW.subscription_term_days IS DISTINCT FROM OLD.subscription_term_days THEN
        RAISE EXCEPTION
            'published plan version % subscription term is immutable', OLD.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_protect_subscription_term
    ON catalog_plan_versions;
CREATE TRIGGER trg_catalog_protect_subscription_term
BEFORE UPDATE OF subscription_term_days ON catalog_plan_versions
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_protect_subscription_term();

COMMIT;
