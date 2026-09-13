-- Fresh reconstruction baseline formerly materialized by legacy CRM startup AutoMigrate.
CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    contact_id BIGINT,
    source BIGINT,
    assign_note TEXT,
    pipeline_id BIGINT,
    stage_id BIGINT,
    priority BIGINT,
    charge_person_id BIGINT,
    charge_person_type BIGINT,
    stage_note TEXT,
    note TEXT
);

CREATE INDEX IF NOT EXISTS idx_customers_deleted_at
    ON customers (deleted_at);

-- IV.10.9 opportunity fields on customers (LeadEntity)
ALTER TABLE customers ADD COLUMN IF NOT EXISTS code VARCHAR(64);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS title VARCHAR(255);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS customer_type INT NOT NULL DEFAULT 10;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS need_summary TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS interested_plan VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS opportunity_status INT NOT NULL DEFAULT 10;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS admin_source INT NOT NULL DEFAULT 70;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS expected_value NUMERIC(18,2);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS probability INT NOT NULL DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS expected_close_date DATE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS next_follow_up_at TIMESTAMPTZ;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS churn_risk INT NOT NULL DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS upgrade_signal BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS renewal_signal BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS owner_team VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS related_user_id BIGINT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS related_business_id BIGINT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS proposal_ref VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS payment_request_ref VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS subscription_ref VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS ticket_ref VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS segment VARCHAR(128);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS region VARCHAR(255);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tags TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS closed_at TIMESTAMPTZ;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS close_reason TEXT;

CREATE INDEX IF NOT EXISTS idx_customers_opportunity_status ON customers (opportunity_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_customers_next_follow_up_at ON customers (next_follow_up_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_customers_code ON customers (code) WHERE deleted_at IS NULL;
