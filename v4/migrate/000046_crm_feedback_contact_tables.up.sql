-- CRM tables used by modules/crm (contact products/tags, appointment stats/reminders, feedback)
-- that db_v4 did not absorb from db_crm. DDL mirrors db_crm; idempotent for DBs that already have them.

CREATE TABLE IF NOT EXISTS tb_contact_product (
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    contact_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    CONSTRAINT tb_contact_product_pkey PRIMARY KEY (contact_id, product_id)
);
CREATE INDEX IF NOT EXISTS idx_tb_contact_product_deleted_at ON tb_contact_product (deleted_at);

CREATE TABLE IF NOT EXISTS tb_contact_tag (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    owner_id BIGINT,
    owner_of INTEGER,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_default_active ON tb_contact_tag (is_default, is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_deleted_at ON tb_contact_tag (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_owner ON tb_contact_tag (owner_id, owner_of) WHERE deleted_at IS NULL;

INSERT INTO tb_contact_tag (id, name, is_default, is_active) VALUES
    (1, 'Khách hàng', true, true),
    (2, 'Môi giới', true, true),
    (3, 'Đối tác', true, true),
    (4, 'Chủ nhà', true, true)
ON CONFLICT (id) DO NOTHING;
SELECT setval(pg_get_serial_sequence('tb_contact_tag', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM tb_contact_tag), 1));

CREATE TABLE IF NOT EXISTS tb_contact_tag_relation (
    contact_id BIGINT NOT NULL REFERENCES tb_contact (id),
    tag_id BIGINT NOT NULL REFERENCES tb_contact_tag (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tb_contact_tag_relation_pkey PRIMARY KEY (contact_id, tag_id)
);
CREATE INDEX IF NOT EXISTS idx_tb_contact_tag_relation_tag_id ON tb_contact_tag_relation (tag_id);

CREATE TABLE IF NOT EXISTS appointment_reminders (
    id BIGSERIAL PRIMARY KEY,
    appointment_id BIGINT,
    remind_at TIMESTAMPTZ,
    is_sent BOOLEAN,
    sent_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS appointment_stats (
    id BIGSERIAL PRIMARY KEY,
    calculate_time TIMESTAMPTZ NOT NULL,
    count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS feedback_registry (
    registry_key SERIAL PRIMARY KEY,
    registry_name VARCHAR(20) NOT NULL,
    port BIGINT NOT NULL,
    service VARCHAR(100) NOT NULL
);

INSERT INTO feedback_registry (registry_key, registry_name, port, service) VALUES
    (1, 'user', 8201, 'user')
ON CONFLICT (registry_key) DO NOTHING;
SELECT setval(pg_get_serial_sequence('feedback_registry', 'registry_key'), GREATEST((SELECT COALESCE(MAX(registry_key), 1) FROM feedback_registry), 1));

CREATE TABLE IF NOT EXISTS feedback_rates (
    created_by BIGINT,
    updated_by BIGINT,
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    anonymous BOOLEAN NOT NULL,
    owner_id BIGINT NOT NULL,
    owner_of BIGINT NOT NULL DEFAULT 10,
    score SMALLINT NOT NULL,
    comment TEXT,
    parent_id BIGINT,
    root_id BIGINT,
    hidden BOOLEAN DEFAULT false,
    approved_at TIMESTAMPTZ,
    approved_by BIGINT
);
CREATE INDEX IF NOT EXISTS idx_feedback_rates_deleted_at ON feedback_rates (deleted_at);

CREATE TABLE IF NOT EXISTS feedback_rate_attachments (
    created_by BIGINT,
    updated_by BIGINT,
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    rate_id BIGINT NOT NULL CONSTRAINT fk_feedback_rates_attachs REFERENCES feedback_rates (id),
    file_url VARCHAR(500),
    file_name VARCHAR(255),
    file_type VARCHAR(50)
);
CREATE INDEX IF NOT EXISTS idx_feedback_rate_attachments_deleted_at ON feedback_rate_attachments (deleted_at);

CREATE TABLE IF NOT EXISTS feedback_report_reasons (
    created_by BIGINT,
    updated_by BIGINT,
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true
);
CREATE INDEX IF NOT EXISTS idx_feedback_report_reasons_deleted_at ON feedback_report_reasons (deleted_at);

CREATE TABLE IF NOT EXISTS feedback_reports (
    created_by BIGINT,
    updated_by BIGINT,
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    reason_id BIGINT NOT NULL CONSTRAINT fk_feedback_reports_reason REFERENCES feedback_report_reasons (id),
    report_status SMALLINT NOT NULL DEFAULT 10,
    user_id BIGINT NOT NULL,
    response VARCHAR(255),
    admin_note VARCHAR(255),
    owner_id BIGINT NOT NULL,
    owner_of BIGINT NOT NULL,
    content TEXT
);
CREATE INDEX IF NOT EXISTS idx_feedback_reports_deleted_at ON feedback_reports (deleted_at);

CREATE TABLE IF NOT EXISTS feedback_report_proofs (
    created_by BIGINT,
    updated_by BIGINT,
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    report_id BIGINT NOT NULL CONSTRAINT fk_feedback_reports_proof_docs REFERENCES feedback_reports (id),
    file_url VARCHAR(500),
    file_name VARCHAR(255),
    file_type VARCHAR(50)
);
CREATE INDEX IF NOT EXISTS idx_feedback_report_proofs_deleted_at ON feedback_report_proofs (deleted_at);
