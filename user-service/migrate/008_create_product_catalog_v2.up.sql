BEGIN;

CREATE TABLE IF NOT EXISTS catalog_products (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(128) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_catalog_products_code
        CHECK (code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_products_status
        CHECK (status IN ('draft', 'active', 'retired'))
);

CREATE TABLE IF NOT EXISTS catalog_plans (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES catalog_products(id),
    code VARCHAR(128) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_catalog_plans_code
        CHECK (code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_plans_status
        CHECK (status IN ('draft', 'active', 'retired'))
);

CREATE INDEX IF NOT EXISTS idx_catalog_plans_product_id
    ON catalog_plans(product_id);

CREATE TABLE IF NOT EXISTS catalog_plan_versions (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES catalog_plans(id),
    version VARCHAR(32) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    subject_scope VARCHAR(32) NOT NULL,
    effective_from TIMESTAMPTZ,
    effective_until TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    source VARCHAR(32) NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    published_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_catalog_plan_versions UNIQUE (plan_id, version),
    CONSTRAINT chk_catalog_plan_versions_version
        CHECK (version ~ '^[0-9]+[.][0-9]+[.][0-9]+$'),
    CONSTRAINT chk_catalog_plan_versions_status
        CHECK (status IN ('draft', 'active', 'retired')),
    CONSTRAINT chk_catalog_plan_versions_subject_scope
        CHECK (subject_scope IN ('profile', 'organization', 'any')),
    CONSTRAINT chk_catalog_plan_versions_source
        CHECK (source IN ('native', 'legacy_backfill', 'import')),
    CONSTRAINT chk_catalog_plan_versions_checksum
        CHECK (terms_checksum ~ '^[0-9a-f]{64}$'),
    CONSTRAINT chk_catalog_plan_versions_effective_window
        CHECK (
            (status = 'draft' AND effective_from IS NULL)
            OR
            (status <> 'draft' AND effective_from IS NOT NULL)
        ),
    CONSTRAINT chk_catalog_plan_versions_effective_until
        CHECK (
            effective_until IS NULL
            OR effective_from IS NULL
            OR effective_until > effective_from
        ),
    CONSTRAINT chk_catalog_plan_versions_publish_state
        CHECK (published_at IS NULL OR status <> 'draft')
);

CREATE INDEX IF NOT EXISTS idx_catalog_plan_versions_plan_id
    ON catalog_plan_versions(plan_id);
CREATE INDEX IF NOT EXISTS idx_catalog_plan_versions_effective
    ON catalog_plan_versions(plan_id, effective_from, effective_until)
    WHERE status = 'active';

CREATE TABLE IF NOT EXISTS catalog_plan_entitlements (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL REFERENCES catalog_plan_versions(id),
    code VARCHAR(128) NOT NULL,
    kind VARCHAR(32) NOT NULL,
    feature_code VARCHAR(128),
    meter_code VARCHAR(128),
    amount BIGINT NOT NULL DEFAULT 0,
    unlimited BOOLEAN NOT NULL DEFAULT FALSE,
    period VARCHAR(32) NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_catalog_plan_entitlement_code
        UNIQUE (plan_version_id, code),
    CONSTRAINT chk_catalog_plan_entitlements_code
        CHECK (code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_plan_entitlements_kind
        CHECK (kind IN ('feature_access', 'usage_allowance', 'capacity_limit')),
    CONSTRAINT chk_catalog_plan_entitlements_feature_code
        CHECK (
            feature_code IS NULL
            OR feature_code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'
        ),
    CONSTRAINT chk_catalog_plan_entitlements_meter_code
        CHECK (
            meter_code IS NULL
            OR meter_code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'
        ),
    CONSTRAINT chk_catalog_plan_entitlements_period
        CHECK (
            period IN (
                'none',
                'day',
                'calendar_month',
                'subscription_cycle',
                'lifetime'
            )
        ),
    CONSTRAINT chk_catalog_plan_entitlements_shape
        CHECK (
            (
                kind = 'feature_access'
                AND feature_code IS NOT NULL
                AND meter_code IS NULL
                AND amount = 0
                AND unlimited = FALSE
                AND period = 'none'
            )
            OR
            (
                kind = 'usage_allowance'
                AND feature_code IS NULL
                AND meter_code IS NOT NULL
                AND (
                    (unlimited = TRUE AND amount = 0)
                    OR (unlimited = FALSE AND amount > 0)
                )
                AND period IN (
                    'day',
                    'calendar_month',
                    'subscription_cycle',
                    'lifetime'
                )
            )
            OR
            (
                kind = 'capacity_limit'
                AND feature_code IS NULL
                AND meter_code IS NOT NULL
                AND (
                    (unlimited = TRUE AND amount = 0)
                    OR (unlimited = FALSE AND amount > 0)
                )
                AND period = 'none'
            )
        )
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_entitlement_feature_target
    ON catalog_plan_entitlements(plan_version_id, feature_code)
    WHERE kind = 'feature_access';
CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_entitlement_usage_target
    ON catalog_plan_entitlements(plan_version_id, meter_code, period)
    WHERE kind = 'usage_allowance';
CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_entitlement_capacity_target
    ON catalog_plan_entitlements(plan_version_id, meter_code)
    WHERE kind = 'capacity_limit';

CREATE TABLE IF NOT EXISTS catalog_price_items (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL REFERENCES catalog_plan_versions(id),
    code VARCHAR(128) NOT NULL,
    kind VARCHAR(20) NOT NULL,
    currency CHAR(3) NOT NULL,
    amount_minor BIGINT NOT NULL,
    billing_unit VARCHAR(128) NOT NULL,
    meter_code VARCHAR(128),
    quantity BIGINT NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_catalog_price_item_code
        UNIQUE (plan_version_id, code),
    CONSTRAINT chk_catalog_price_items_code
        CHECK (code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_price_items_kind
        CHECK (kind IN ('recurring', 'add_on', 'overage')),
    CONSTRAINT chk_catalog_price_items_currency
        CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT chk_catalog_price_items_amount
        CHECK (amount_minor >= 0),
    CONSTRAINT chk_catalog_price_items_billing_unit
        CHECK (billing_unit ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_price_items_meter_code
        CHECK (
            meter_code IS NULL
            OR meter_code ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'
        ),
    CONSTRAINT chk_catalog_price_items_quantity
        CHECK (quantity > 0),
    CONSTRAINT chk_catalog_price_items_shape
        CHECK (
            (kind = 'recurring' AND meter_code IS NULL)
            OR
            (kind IN ('add_on', 'overage') AND meter_code IS NOT NULL)
        )
);

CREATE TABLE IF NOT EXISTS catalog_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    subscription_key VARCHAR(128) NOT NULL UNIQUE,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    product_id BIGINT NOT NULL REFERENCES catalog_products(id),
    plan_version_id BIGINT NOT NULL REFERENCES catalog_plan_versions(id),
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end TIMESTAMPTZ NOT NULL,
    access_until TIMESTAMPTZ,
    auto_renew BOOLEAN NOT NULL DEFAULT FALSE,
    order_reference VARCHAR(128) NOT NULL DEFAULT '',
    canceled_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_catalog_subscriptions_key
        CHECK (subscription_key ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_catalog_subscriptions_subject_kind
        CHECK (subject_kind IN ('profile', 'organization')),
    CONSTRAINT chk_catalog_subscriptions_subject_id
        CHECK (subject_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_catalog_subscriptions_status
        CHECK (status IN ('pending', 'active', 'past_due', 'canceled', 'expired')),
    CONSTRAINT chk_catalog_subscriptions_period
        CHECK (current_period_end > current_period_start),
    CONSTRAINT chk_catalog_subscriptions_access_until
        CHECK (access_until IS NULL OR access_until >= started_at),
    CONSTRAINT chk_catalog_subscriptions_canceled_at
        CHECK (canceled_at IS NULL OR canceled_at >= started_at),
    CONSTRAINT chk_catalog_subscriptions_ended_at
        CHECK (ended_at IS NULL OR ended_at >= started_at)
);

CREATE INDEX IF NOT EXISTS idx_catalog_subscriptions_subject
    ON catalog_subscriptions(subject_kind, subject_id);
CREATE INDEX IF NOT EXISTS idx_catalog_subscriptions_plan_version
    ON catalog_subscriptions(plan_version_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_catalog_subscriptions_active_product
    ON catalog_subscriptions(subject_kind, subject_id, product_id)
    WHERE status IN ('pending', 'active', 'past_due');

CREATE TABLE IF NOT EXISTS catalog_subscription_events (
    id BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL REFERENCES catalog_subscriptions(id),
    event_type VARCHAR(64) NOT NULL,
    from_plan_version_id BIGINT REFERENCES catalog_plan_versions(id),
    to_plan_version_id BIGINT REFERENCES catalog_plan_versions(id),
    effective_at TIMESTAMPTZ NOT NULL,
    reason TEXT NOT NULL,
    actor_kind VARCHAR(32) NOT NULL,
    actor_id VARCHAR(128) NOT NULL DEFAULT '',
    request_id VARCHAR(128) NOT NULL DEFAULT '',
    operation_id VARCHAR(128) NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_catalog_subscription_events_type
        CHECK (event_type ~ '^[a-z][a-z0-9_]*([.][a-z0-9_]+)*$'),
    CONSTRAINT chk_catalog_subscription_events_reason
        CHECK (LENGTH(BTRIM(reason)) > 0),
    CONSTRAINT chk_catalog_subscription_events_actor_kind
        CHECK (actor_kind IN ('user', 'internal_service', 'migration', 'system')),
    CONSTRAINT chk_catalog_subscription_events_actor_id
        CHECK (actor_id = '' OR actor_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_catalog_subscription_events_request_id
        CHECK (request_id = '' OR request_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$'),
    CONSTRAINT chk_catalog_subscription_events_operation_id
        CHECK (operation_id = '' OR operation_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$')
);

CREATE INDEX IF NOT EXISTS idx_catalog_subscription_events_subscription
    ON catalog_subscription_events(subscription_id, effective_at, id);

CREATE TABLE IF NOT EXISTS legacy_plan_catalog_mappings (
    legacy_plan_id BIGINT PRIMARY KEY,
    plan_version_id BIGINT UNIQUE REFERENCES catalog_plan_versions(id),
    source_checksum CHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    decision JSONB NOT NULL DEFAULT '{}'::jsonb,
    reviewed_by BIGINT,
    reviewed_at TIMESTAMPTZ,
    migrated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_legacy_plan_catalog_mapping_checksum
        CHECK (source_checksum ~ '^[0-9a-f]{64}$'),
    CONSTRAINT chk_legacy_plan_catalog_mapping_status
        CHECK (
            status IN (
                'pending_review',
                'draft_created',
                'shadow_read',
                'cutover',
                'rolled_back',
                'blocked'
            )
        )
);

CREATE TABLE IF NOT EXISTS legacy_membership_subscription_mappings (
    legacy_membership_id BIGINT PRIMARY KEY,
    subscription_id BIGINT UNIQUE REFERENCES catalog_subscriptions(id),
    source_checksum CHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    migrated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_legacy_membership_subscription_mapping_checksum
        CHECK (source_checksum ~ '^[0-9a-f]{64}$'),
    CONSTRAINT chk_legacy_membership_subscription_mapping_status
        CHECK (
            status IN (
                'pending_review',
                'draft_created',
                'shadow_read',
                'cutover',
                'rolled_back',
                'blocked'
            )
        )
);

CREATE OR REPLACE FUNCTION qhpro_catalog_validate_plan_version_window()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_advisory_xact_lock(
        hashtext('qhpro_catalog_plan_version'),
        (NEW.plan_id % 2147483647)::INTEGER
    );

    IF NEW.status = 'active' AND EXISTS (
        SELECT 1
        FROM catalog_plan_versions existing
        WHERE existing.plan_id = NEW.plan_id
          AND existing.id <> COALESCE(NEW.id, 0)
          AND existing.status = 'active'
          AND tstzrange(
                existing.effective_from,
                existing.effective_until,
                '[)'
              ) && tstzrange(
                NEW.effective_from,
                NEW.effective_until,
                '[)'
              )
    ) THEN
        RAISE EXCEPTION
            'active plan version window overlaps plan_id=%', NEW.plan_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_validate_plan_version_window
    ON catalog_plan_versions;
CREATE TRIGGER trg_catalog_validate_plan_version_window
BEFORE INSERT OR UPDATE OF plan_id, status, effective_from, effective_until
ON catalog_plan_versions
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_validate_plan_version_window();

CREATE OR REPLACE FUNCTION qhpro_catalog_protect_plan_version_terms()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        IF OLD.published_at IS NOT NULL THEN
            RAISE EXCEPTION
                'published plan version % cannot be deleted', OLD.id;
        END IF;
        RETURN OLD;
    END IF;

    IF OLD.status = 'draft' AND NEW.status NOT IN ('draft', 'active') THEN
        RAISE EXCEPTION
            'draft plan version % may only become active', OLD.id;
    END IF;

    IF OLD.status = 'active' AND NEW.status NOT IN ('active', 'retired') THEN
        RAISE EXCEPTION
            'active plan version % may only become retired', OLD.id;
    END IF;

    IF OLD.status = 'retired' AND NEW.status <> 'retired' THEN
        RAISE EXCEPTION
            'retired plan version % cannot be reactivated', OLD.id;
    END IF;

    IF OLD.published_at IS NULL AND NEW.published_at IS NOT NULL
       AND NEW.status <> 'active' THEN
        RAISE EXCEPTION
            'plan version % must be active when published', OLD.id;
    END IF;

    IF OLD.published_at IS NOT NULL AND (
        NEW.plan_id IS DISTINCT FROM OLD.plan_id
        OR NEW.version IS DISTINCT FROM OLD.version
        OR NEW.display_name IS DISTINCT FROM OLD.display_name
        OR NEW.subject_scope IS DISTINCT FROM OLD.subject_scope
        OR NEW.effective_from IS DISTINCT FROM OLD.effective_from
        OR NEW.source IS DISTINCT FROM OLD.source
        OR NEW.terms_checksum IS DISTINCT FROM OLD.terms_checksum
        OR NEW.published_at IS DISTINCT FROM OLD.published_at
    ) THEN
        RAISE EXCEPTION
            'published plan version % commercial terms are immutable', OLD.id;
    END IF;

    IF OLD.published_at IS NOT NULL
       AND NEW.effective_until IS DISTINCT FROM OLD.effective_until THEN
        IF NEW.effective_until IS NULL THEN
            RAISE EXCEPTION
                'published plan version % availability cannot be reopened', OLD.id;
        END IF;
        IF OLD.effective_until IS NOT NULL
           AND NEW.effective_until > OLD.effective_until THEN
            RAISE EXCEPTION
                'published plan version % availability cannot be extended', OLD.id;
        END IF;
        IF NEW.effective_until <= transaction_timestamp() THEN
            RAISE EXCEPTION
                'published plan version % must close in the future', OLD.id;
        END IF;
        IF NEW.updated_by IS NULL THEN
            RAISE EXCEPTION
                'published plan version % availability change requires updated_by',
                OLD.id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_protect_plan_version_terms
    ON catalog_plan_versions;
CREATE TRIGGER trg_catalog_protect_plan_version_terms
BEFORE UPDATE OR DELETE ON catalog_plan_versions
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_protect_plan_version_terms();

CREATE OR REPLACE FUNCTION qhpro_catalog_protect_version_child()
RETURNS TRIGGER AS $$
DECLARE
    old_parent_published_at TIMESTAMPTZ;
    new_parent_published_at TIMESTAMPTZ;
BEGIN
    IF TG_OP IN ('UPDATE', 'DELETE') THEN
        SELECT published_at
          INTO old_parent_published_at
          FROM catalog_plan_versions
         WHERE id = OLD.plan_version_id
         FOR UPDATE;

        IF old_parent_published_at IS NOT NULL THEN
            RAISE EXCEPTION
                'published plan version % child records are immutable',
                OLD.plan_version_id;
        END IF;
    END IF;

    IF TG_OP IN ('INSERT', 'UPDATE') THEN
        SELECT published_at
          INTO new_parent_published_at
          FROM catalog_plan_versions
         WHERE id = NEW.plan_version_id
         FOR UPDATE;

        IF new_parent_published_at IS NOT NULL THEN
            RAISE EXCEPTION
                'published plan version % child records are immutable',
                NEW.plan_version_id;
        END IF;
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_protect_entitlement
    ON catalog_plan_entitlements;
CREATE TRIGGER trg_catalog_protect_entitlement
BEFORE INSERT OR UPDATE OR DELETE ON catalog_plan_entitlements
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_protect_version_child();

DROP TRIGGER IF EXISTS trg_catalog_protect_price_item
    ON catalog_price_items;
CREATE TRIGGER trg_catalog_protect_price_item
BEFORE INSERT OR UPDATE OR DELETE ON catalog_price_items
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_protect_version_child();

CREATE OR REPLACE FUNCTION qhpro_catalog_validate_subscription_product()
RETURNS TRIGGER AS $$
DECLARE
    version_product_id BIGINT;
    version_scope VARCHAR(32);
    version_status VARCHAR(20);
    version_effective_from TIMESTAMPTZ;
    version_effective_until TIMESTAMPTZ;
    version_published_at TIMESTAMPTZ;
BEGIN
    SELECT
        plan.product_id,
        version.subject_scope,
        version.status,
        version.effective_from,
        version.effective_until,
        version.published_at
      INTO
        version_product_id,
        version_scope,
        version_status,
        version_effective_from,
        version_effective_until,
        version_published_at
      FROM catalog_plan_versions version
      JOIN catalog_plans plan ON plan.id = version.plan_id
     WHERE version.id = NEW.plan_version_id
     FOR SHARE OF version, plan;

    IF version_product_id IS NULL OR version_product_id <> NEW.product_id THEN
        RAISE EXCEPTION
            'subscription product does not match plan version';
    END IF;

    IF version_published_at IS NULL OR version_status <> 'active' THEN
        RAISE EXCEPTION
            'subscription requires a published active plan version';
    END IF;

    IF version_scope <> 'any' AND version_scope <> NEW.subject_kind THEN
        RAISE EXCEPTION
            'subscription subject kind does not match plan scope';
    END IF;

    IF NEW.status IN ('active', 'past_due') AND (
        NEW.started_at < version_effective_from
        OR (
            version_effective_until IS NOT NULL
            AND NEW.started_at >= version_effective_until
        )
    ) THEN
        RAISE EXCEPTION
            'subscription start is outside the plan effective window';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_validate_subscription_product
    ON catalog_subscriptions;
CREATE TRIGGER trg_catalog_validate_subscription_product
BEFORE INSERT OR UPDATE OF product_id, plan_version_id
ON catalog_subscriptions
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_validate_subscription_product();

CREATE OR REPLACE FUNCTION qhpro_catalog_protect_subscription_event()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'subscription events are append-only';
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_catalog_protect_subscription_event
    ON catalog_subscription_events;
CREATE TRIGGER trg_catalog_protect_subscription_event
BEFORE UPDATE OR DELETE ON catalog_subscription_events
FOR EACH ROW
EXECUTE FUNCTION qhpro_catalog_protect_subscription_event();

COMMENT ON TABLE catalog_plan_versions IS
    'Immutable commercial terms with an audited, narrowing-only availability end.';
COMMENT ON TABLE catalog_subscriptions IS
    'Subject ownership of one immutable plan version and its explicit billing period.';
COMMENT ON TABLE legacy_plan_catalog_mappings IS
    'Migration state only; legacy plans remain the runtime source until cutover.';

COMMIT;
