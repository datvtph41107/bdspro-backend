-- Module 1: SEO Objective & Value Foundation.
--
-- Objective giải thích tại sao một chương trình SEO tồn tại.
-- Không đặt keyword, canonical, title hoặc robots trong module này.

CREATE TABLE IF NOT EXISTS seo_objective (
    id BIGSERIAL PRIMARY KEY,

    -- Semantic identity ổn định dùng xuyên backend/admin.
    code VARCHAR(120) NOT NULL,
    name TEXT NOT NULL,

    -- Mô tả cấp chiến lược, chưa phải audience model chi tiết.
    audience_summary TEXT NOT NULL DEFAULT '',
    situation TEXT NOT NULL DEFAULT '',
    problem_statement TEXT NOT NULL DEFAULT '',

    -- Phân tách giá trị cho ba đối tượng khác nhau.
    user_outcome TEXT NOT NULL DEFAULT '',
    product_outcome TEXT NOT NULL DEFAULT '',
    business_outcome TEXT NOT NULL DEFAULT '',

    -- Giữ phạm vi rõ để tránh objective mở rộng không kiểm soát.
    non_goals TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    constraints TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],

    status VARCHAR(24) NOT NULL DEFAULT 'draft',
    confidence VARCHAR(24) NOT NULL DEFAULT 'exploratory',

    -- Owner là mã tổ chức/team, không phải display label.
    owner_code VARCHAR(120) NOT NULL DEFAULT '',
    review_due_at TIMESTAMPTZ NULL,

    -- Optimistic locking, ngăn hai admin ghi đè lẫn nhau.
    lock_version INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT chk_seo_objective_code
        CHECK (LENGTH(TRIM(code)) > 0),

    CONSTRAINT chk_seo_objective_name
        CHECK (LENGTH(TRIM(name)) > 0),

    CONSTRAINT chk_seo_objective_status
        CHECK (status IN (
            'draft',
            'active',
            'paused',
            'archived'
        )),

    CONSTRAINT chk_seo_objective_confidence
        CHECK (confidence IN (
            'exploratory',
            'supported',
            'validated'
        )),

    CONSTRAINT chk_seo_objective_lock_version
        CHECK (lock_version > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_seo_objective_code_active
ON seo_objective(code)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_objective_status
ON seo_objective(status)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_objective_confidence
ON seo_objective(confidence)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_objective_owner
ON seo_objective(owner_code)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_objective_review_due
ON seo_objective(review_due_at)
WHERE deleted_at IS NULL
  AND status IN ('active', 'paused');


CREATE TABLE IF NOT EXISTS seo_objective_evidence (
    id BIGSERIAL PRIMARY KEY,

    objective_id BIGINT NOT NULL
        REFERENCES seo_objective(id)
        ON DELETE CASCADE,

    evidence_type VARCHAR(40) NOT NULL,
    stance VARCHAR(20) NOT NULL DEFAULT 'supports',

    source_system VARCHAR(80) NOT NULL DEFAULT '',
    source_reference TEXT NOT NULL DEFAULT '',

    finding TEXT NOT NULL,
    confidence VARCHAR(24) NOT NULL DEFAULT 'exploratory',

    observed_at TIMESTAMPTZ NULL,
    valid_until TIMESTAMPTZ NULL,

    status VARCHAR(24) NOT NULL DEFAULT 'active',

    -- Chỉ chứa metadata riêng của connector.
    -- Không đưa các field nghiệp vụ chính vào đây.
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_seo_objective_evidence_type
        CHECK (evidence_type IN (
            'assumption',
            'user_interview',
            'search_console',
            'analytics',
            'internal_search',
            'support_request',
            'crm',
            'market_research',
            'experiment',
            'expert_review'
        )),

    CONSTRAINT chk_seo_objective_evidence_stance
        CHECK (stance IN (
            'supports',
            'contradicts',
            'context'
        )),

    CONSTRAINT chk_seo_objective_evidence_confidence
        CHECK (confidence IN (
            'exploratory',
            'supported',
            'validated'
        )),

    CONSTRAINT chk_seo_objective_evidence_status
        CHECK (status IN (
            'active',
            'superseded',
            'rejected'
        )),

    CONSTRAINT chk_seo_objective_evidence_finding
        CHECK (LENGTH(TRIM(finding)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_seo_objective_evidence_objective
ON seo_objective_evidence(objective_id);

CREATE INDEX IF NOT EXISTS idx_seo_objective_evidence_source
ON seo_objective_evidence(source_system, evidence_type);

CREATE INDEX IF NOT EXISTS idx_seo_objective_evidence_validity
ON seo_objective_evidence(valid_until)
WHERE status = 'active';


CREATE TABLE IF NOT EXISTS seo_objective_metric (
    id BIGSERIAL PRIMARY KEY,

    objective_id BIGINT NOT NULL
        REFERENCES seo_objective(id)
        ON DELETE CASCADE,

    metric_code VARCHAR(120) NOT NULL,
    metric_name TEXT NOT NULL,

    metric_layer VARCHAR(24) NOT NULL,
    description TEXT NOT NULL DEFAULT '',

    source_system VARCHAR(80) NOT NULL DEFAULT '',
    calculation_method TEXT NOT NULL DEFAULT '',

    unit VARCHAR(32) NOT NULL DEFAULT 'count',
    target_direction VARCHAR(20) NOT NULL DEFAULT 'increase',

    baseline_value NUMERIC(20, 6) NULL,
    target_value NUMERIC(20, 6) NULL,

    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    instrumentation_status VARCHAR(24)
        NOT NULL DEFAULT 'unknown',

    status VARCHAR(24) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_seo_objective_metric
        UNIQUE (objective_id, metric_code),

    CONSTRAINT chk_seo_objective_metric_code
        CHECK (LENGTH(TRIM(metric_code)) > 0),

    CONSTRAINT chk_seo_objective_metric_layer
        CHECK (metric_layer IN (
            'search',
            'technical',
            'user',
            'product',
            'business',
            'trust'
        )),

    CONSTRAINT chk_seo_objective_metric_direction
        CHECK (target_direction IN (
            'increase',
            'decrease',
            'maintain'
        )),

    CONSTRAINT chk_seo_objective_metric_instrumentation
        CHECK (instrumentation_status IN (
            'unknown',
            'planned',
            'instrumented',
            'verified'
        )),

    CONSTRAINT chk_seo_objective_metric_status
        CHECK (status IN (
            'active',
            'deprecated'
        ))
);

-- Mỗi objective chỉ có một primary metric active.
CREATE UNIQUE INDEX IF NOT EXISTS uq_seo_objective_primary_metric
ON seo_objective_metric(objective_id)
WHERE is_primary = TRUE
  AND status = 'active';

CREATE INDEX IF NOT EXISTS idx_seo_objective_metric_objective
ON seo_objective_metric(objective_id);

CREATE INDEX IF NOT EXISTS idx_seo_objective_metric_layer
ON seo_objective_metric(metric_layer);