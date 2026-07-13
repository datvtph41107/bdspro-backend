-- Module 2: Audience, Situation & User Need.
--
-- Module này định nghĩa ai đang gặp vấn đề gì và trong hoàn cảnh nào.
-- Không chứa keyword, entity, URL, canonical hoặc index state.

CREATE TABLE seo_audience_segment (
    id BIGSERIAL PRIMARY KEY,

    -- Semantic identity dùng xuyên backend và admin.
    code VARCHAR(120) NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',

    -- Vai trò trong quyết định giúp phân biệt người tự quyết định,
    -- người tư vấn và người nghiên cứu.
    decision_role VARCHAR(40) NOT NULL DEFAULT 'unknown',

    -- Ảnh hưởng đến thuật ngữ, mức giải thích và độ sâu answer sau này.
    knowledge_level VARCHAR(32) NOT NULL DEFAULT 'unknown',

    status VARCHAR(24) NOT NULL DEFAULT 'draft',

    owner_code VARCHAR(120) NOT NULL DEFAULT '',
    review_due_at TIMESTAMPTZ,

    lock_version INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_seo_audience_code
        CHECK (LENGTH(TRIM(code)) > 0),

    CONSTRAINT chk_seo_audience_name
        CHECK (LENGTH(TRIM(name)) > 0),

    CONSTRAINT chk_seo_audience_decision_role
        CHECK (decision_role IN (
            'unknown',
            'self_decision_maker',
            'advisor',
            'researcher',
            'operator'
        )),

    CONSTRAINT chk_seo_audience_knowledge_level
        CHECK (knowledge_level IN (
            'unknown',
            'beginner',
            'intermediate',
            'advanced',
            'expert'
        )),

    CONSTRAINT chk_seo_audience_status
        CHECK (status IN (
            'draft',
            'active',
            'paused',
            'archived'
        )),

    CONSTRAINT chk_seo_audience_lock_version
        CHECK (lock_version > 0)
);

CREATE UNIQUE INDEX uq_seo_audience_code_active
ON seo_audience_segment(code)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_audience_status
ON seo_audience_segment(status)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_audience_role
ON seo_audience_segment(decision_role)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_audience_owner
ON seo_audience_segment(owner_code)
WHERE deleted_at IS NULL;


CREATE TABLE seo_user_need (
    id BIGSERIAL PRIMARY KEY,

    objective_id BIGINT NOT NULL
        REFERENCES seo_objective(id)
        ON DELETE RESTRICT,

    audience_segment_id BIGINT NOT NULL
        REFERENCES seo_audience_segment(id)
        ON DELETE RESTRICT,

    code VARCHAR(120) NOT NULL,
    name TEXT NOT NULL,

    -- Giai đoạn quyết định giúp phân biệt tìm hiểu,
    -- so sánh, xác minh, giao dịch và theo dõi.
    decision_stage VARCHAR(32) NOT NULL DEFAULT 'unknown',

    situation TEXT NOT NULL DEFAULT '',
    trigger_statement TEXT NOT NULL DEFAULT '',
    job_to_be_done TEXT NOT NULL DEFAULT '',

    barrier_statement TEXT NOT NULL DEFAULT '',
    current_alternative TEXT NOT NULL DEFAULT '',

    desired_outcome TEXT NOT NULL DEFAULT '',
    failure_impact TEXT NOT NULL DEFAULT '',

    risk_level VARCHAR(24) NOT NULL DEFAULT 'unknown',
    priority VARCHAR(24) NOT NULL DEFAULT 'unassessed',

    status VARCHAR(24) NOT NULL DEFAULT 'draft',
    confidence VARCHAR(24) NOT NULL DEFAULT 'exploratory',

    owner_code VARCHAR(120) NOT NULL DEFAULT '',
    review_due_at TIMESTAMPTZ,

    lock_version INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Cần cho composite FK bảo đảm evidence và need
    -- thuộc cùng một objective.
    CONSTRAINT uq_seo_user_need_id_objective
        UNIQUE (id, objective_id),

    CONSTRAINT chk_seo_user_need_code
        CHECK (LENGTH(TRIM(code)) > 0),

    CONSTRAINT chk_seo_user_need_name
        CHECK (LENGTH(TRIM(name)) > 0),

    CONSTRAINT chk_seo_user_need_decision_stage
        CHECK (decision_stage IN (
            'unknown',
            'exploration',
            'comparison',
            'verification',
            'transaction',
            'ownership',
            'monitoring'
        )),

    CONSTRAINT chk_seo_user_need_risk
        CHECK (risk_level IN (
            'unknown',
            'low',
            'medium',
            'high',
            'critical'
        )),

    CONSTRAINT chk_seo_user_need_priority
        CHECK (priority IN (
            'unassessed',
            'low',
            'medium',
            'high',
            'critical'
        )),

    CONSTRAINT chk_seo_user_need_status
        CHECK (status IN (
            'draft',
            'active',
            'paused',
            'archived'
        )),

    CONSTRAINT chk_seo_user_need_confidence
        CHECK (confidence IN (
            'exploratory',
            'supported',
            'validated'
        )),

    CONSTRAINT chk_seo_user_need_lock_version
        CHECK (lock_version > 0)
);

CREATE UNIQUE INDEX uq_seo_user_need_code_active
ON seo_user_need(objective_id, code)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_user_need_objective
ON seo_user_need(objective_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_user_need_audience
ON seo_user_need(audience_segment_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_user_need_status
ON seo_user_need(status)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_user_need_priority
ON seo_user_need(priority, risk_level)
WHERE deleted_at IS NULL;


-- Cho phép database xác minh evidence và user need
-- cùng thuộc một objective.
ALTER TABLE seo_objective_evidence
ADD CONSTRAINT uq_seo_objective_evidence_id_objective
UNIQUE (id, objective_id);


CREATE TABLE seo_user_need_evidence (
    objective_id BIGINT NOT NULL,

    user_need_id BIGINT NOT NULL,
    evidence_id BIGINT NOT NULL,

    -- Relation nằm ở link vì cùng một finding có thể hỗ trợ
    -- một need nhưng chỉ cung cấp context cho need khác.
    relation_type VARCHAR(20) NOT NULL DEFAULT 'supports',

    relevance_note TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_need_id, evidence_id),

    CONSTRAINT chk_seo_user_need_evidence_relation
        CHECK (relation_type IN (
            'supports',
            'contradicts',
            'context'
        )),

    CONSTRAINT fk_seo_user_need_evidence_need
        FOREIGN KEY (user_need_id, objective_id)
        REFERENCES seo_user_need(id, objective_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_seo_user_need_evidence_evidence
        FOREIGN KEY (evidence_id, objective_id)
        REFERENCES seo_objective_evidence(id, objective_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_seo_user_need_evidence_objective
ON seo_user_need_evidence(objective_id);

CREATE INDEX idx_seo_user_need_evidence_evidence
ON seo_user_need_evidence(evidence_id);