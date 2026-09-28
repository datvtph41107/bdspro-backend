BEGIN;

-- Repair databases whose schema_migrations advanced past 000001 without the
-- commercial catalog tables. Safe on fresh installs (IF NOT EXISTS).

CREATE TABLE IF NOT EXISTS catalog_products (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(128) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_plans (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    code VARCHAR(128) NOT NULL,
    tier_rank INTEGER,
    status VARCHAR(20) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_plan_versions (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL,
    version VARCHAR(32) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    subject_scope VARCHAR(32) NOT NULL,
    subscription_term_days INTEGER,
    effective_from TIMESTAMPTZ,
    effective_until TIMESTAMPTZ,
    published_at TIMESTAMPTZ,
    source VARCHAR(32) NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    published_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_plan_entitlements (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL,
    code VARCHAR(128) NOT NULL,
    kind VARCHAR(32) NOT NULL,
    feature_code VARCHAR(128),
    meter_code VARCHAR(128),
    amount BIGINT NOT NULL DEFAULT 0,
    unlimited BOOLEAN NOT NULL DEFAULT false,
    period VARCHAR(32) NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_plan_operation_policies (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL,
    operation_code VARCHAR(128) NOT NULL,
    feature_code VARCHAR(128) NOT NULL,
    meter_code VARCHAR(128),
    units_per_action BIGINT NOT NULL DEFAULT 0,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_price_items (
    id BIGSERIAL PRIMARY KEY,
    plan_version_id BIGINT NOT NULL,
    code VARCHAR(128) NOT NULL,
    kind VARCHAR(20) NOT NULL,
    currency CHAR(3) NOT NULL,
    amount_minor BIGINT NOT NULL,
    billing_unit VARCHAR(128) NOT NULL,
    meter_code VARCHAR(128),
    quantity BIGINT NOT NULL,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    subscription_key VARCHAR(128) NOT NULL,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    product_id BIGINT NOT NULL,
    plan_version_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end TIMESTAMPTZ NOT NULL,
    access_until TIMESTAMPTZ,
    auto_renew BOOLEAN NOT NULL DEFAULT false,
    order_reference VARCHAR(128),
    pending_plan_version_id BIGINT,
    pending_effective_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS catalog_subscription_events (
    id BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    from_plan_version_id BIGINT,
    to_plan_version_id BIGINT,
    effective_at TIMESTAMPTZ NOT NULL,
    reason TEXT NOT NULL,
    actor_kind VARCHAR(32) NOT NULL,
    actor_id VARCHAR(128) NOT NULL,
    request_id VARCHAR(128) NOT NULL,
    operation_id VARCHAR(128) NOT NULL,
    metadata JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription_checkout_snapshots (
    id BIGSERIAL PRIMARY KEY,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    command_key TEXT NOT NULL,
    request_hash CHAR(64) NOT NULL,
    product_id BIGINT,
    product_code VARCHAR(128) NOT NULL,
    plan_id BIGINT NOT NULL,
    plan_code VARCHAR(128) NOT NULL,
    plan_version_id BIGINT NOT NULL,
    plan_version VARCHAR(128) NOT NULL,
    tier_rank INTEGER NOT NULL CHECK (tier_rank > 0),
    subscription_term_days INTEGER NOT NULL CHECK (subscription_term_days > 0),
    terms_checksum CHAR(64) NOT NULL,
    subject_scope VARCHAR(32) NOT NULL,
    currency CHAR(3) NOT NULL,
    amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription_settlement_receipts (
    effect_key VARCHAR(128) PRIMARY KEY,
    fingerprint CHAR(64) NOT NULL,
    order_id BIGINT NOT NULL,
    subject_kind VARCHAR(32) NOT NULL,
    subject_id VARCHAR(128) NOT NULL,
    product_code VARCHAR(64) NOT NULL,
    plan_code VARCHAR(64) NOT NULL,
    plan_version_id BIGINT NOT NULL,
    plan_version VARCHAR(64) NOT NULL,
    tier_rank INTEGER NOT NULL,
    subscription_term_days INTEGER NOT NULL,
    terms_checksum CHAR(64) NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    action VARCHAR(32) NOT NULL,
    subscription_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_subscription_settlement_receipts_order_id
    ON subscription_settlement_receipts(order_id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_products_code ON catalog_products(code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_plans_product_code ON catalog_plans(product_id, code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_plan_versions_plan_version ON catalog_plan_versions(plan_id, version);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_entitlement_code ON catalog_plan_entitlements(plan_version_id, code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_operation_code ON catalog_plan_operation_policies(plan_version_id, operation_code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_price_code ON catalog_price_items(plan_version_id, code);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_subscription_key ON catalog_subscriptions(subscription_key);
CREATE INDEX IF NOT EXISTS ix_catalog_subscription_subject ON catalog_subscriptions(subject_kind, subject_id, product_id, status);
CREATE INDEX IF NOT EXISTS ix_catalog_subscription_events_subscription ON catalog_subscription_events(subscription_id, created_at);
CREATE UNIQUE INDEX IF NOT EXISTS ux_checkout_command ON subscription_checkout_snapshots(subject_kind, subject_id, command_key);
CREATE INDEX IF NOT EXISTS ix_catalog_plans_product_id ON catalog_plans(product_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_plans_code ON catalog_plans(code);
CREATE INDEX IF NOT EXISTS ix_catalog_plan_versions_plan_id ON catalog_plan_versions(plan_id);
CREATE INDEX IF NOT EXISTS ix_catalog_plan_entitlements_plan_version_id ON catalog_plan_entitlements(plan_version_id);
CREATE INDEX IF NOT EXISTS ix_catalog_plan_operation_policies_plan_version_id ON catalog_plan_operation_policies(plan_version_id);
CREATE INDEX IF NOT EXISTS ix_catalog_price_items_plan_version_id ON catalog_price_items(plan_version_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_subscriptions_subscription_key ON catalog_subscriptions(subscription_key);
CREATE INDEX IF NOT EXISTS idx_catalog_subscription_subject ON catalog_subscriptions(subject_kind, subject_id);
CREATE INDEX IF NOT EXISTS ix_catalog_subscriptions_product_id ON catalog_subscriptions(product_id);
CREATE INDEX IF NOT EXISTS ix_catalog_subscriptions_plan_version_id ON catalog_subscriptions(plan_version_id);
CREATE INDEX IF NOT EXISTS ix_catalog_subscription_events_subscription_id ON catalog_subscription_events(subscription_id);

COMMIT;
