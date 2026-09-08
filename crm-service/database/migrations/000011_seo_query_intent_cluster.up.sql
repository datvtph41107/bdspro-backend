-- Module 3: Query Evidence, Intent Hypothesis & Query Cluster.
--
-- Module này chuyển Active User Need thành:
-- Query Evidence -> Intent Hypothesis -> Query Cluster -> Search Promise.
--
-- Không chứa URL, title, meta description, canonical, robots,
-- sitemap hoặc index state.

BEGIN;

-- ============================================================
-- 1. QUERY EVIDENCE
-- ============================================================
CREATE TABLE seo_query_evidence (
    id BIGSERIAL PRIMARY KEY,

    objective_id BIGINT NOT NULL,
    user_need_id BIGINT NOT NULL,

    raw_query TEXT NOT NULL,
    normalized_query TEXT NOT NULL,
    query_hash VARCHAR(64) NOT NULL,

    source_system VARCHAR(64) NOT NULL,
    source_kind VARCHAR(24) NOT NULL,
    observation_type VARCHAR(24) NOT NULL,
    source_record_key VARCHAR(200) NOT NULL,
    source_reference TEXT NOT NULL DEFAULT '',

    locale VARCHAR(16) NOT NULL DEFAULT 'vi-VN',
    country_code VARCHAR(2) NOT NULL DEFAULT 'VN',
    device_category VARCHAR(16) NOT NULL DEFAULT 'unknown',

    observed_from TIMESTAMPTZ,
    observed_to TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,

    -- Metric phải nullable để phân biệt:
    -- NULL = source không cung cấp/không đo;
    -- 0    = source có đo và kết quả thực sự bằng 0.
    impressions BIGINT,
    clicks BIGINT,
    ctr NUMERIC(9, 6),
    avg_position NUMERIC(10, 4),
    search_volume BIGINT,

    confidence VARCHAR(24) NOT NULL DEFAULT 'exploratory',
    status VARCHAR(24) NOT NULL DEFAULT 'active',

    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_seo_query_evidence_identity
        UNIQUE (id, user_need_id, objective_id),

    CONSTRAINT fk_seo_query_evidence_need
        FOREIGN KEY (user_need_id, objective_id)
        REFERENCES seo_user_need(id, objective_id)
        ON DELETE CASCADE,

    CONSTRAINT chk_seo_query_evidence_raw
        CHECK (LENGTH(TRIM(raw_query)) > 0),

    CONSTRAINT chk_seo_query_evidence_normalized
        CHECK (LENGTH(TRIM(normalized_query)) > 0),

    CONSTRAINT chk_seo_query_evidence_hash
        CHECK (query_hash ~ '^[0-9a-f]{64}$'),

    CONSTRAINT chk_seo_query_evidence_source_system
        CHECK (source_system ~ '^[a-z][a-z0-9_]{1,63}$'),

    CONSTRAINT chk_seo_query_evidence_source_key
        CHECK (LENGTH(TRIM(source_record_key)) > 0),

    CONSTRAINT chk_seo_query_evidence_source_kind
        CHECK (source_kind IN (
            'first_party',
            'third_party',
            'qualitative',
            'manual'
        )),

    CONSTRAINT chk_seo_query_evidence_observation
        CHECK (observation_type IN (
            'observed',
            'reported',
            'suggested',
            'assumed'
        )),

    CONSTRAINT chk_seo_query_evidence_device
        CHECK (device_category IN (
            'unknown',
            'mobile',
            'desktop',
            'tablet',
            'all'
        )),

    CONSTRAINT chk_seo_query_evidence_country
        CHECK (country_code ~ '^[A-Z]{2}$'),

    CONSTRAINT chk_seo_query_evidence_period
        CHECK (
            observed_from IS NULL
            OR observed_to IS NULL
            OR observed_to >= observed_from
        ),

    CONSTRAINT chk_seo_query_evidence_impressions
        CHECK (impressions IS NULL OR impressions >= 0),

    CONSTRAINT chk_seo_query_evidence_clicks
        CHECK (clicks IS NULL OR clicks >= 0),

    CONSTRAINT chk_seo_query_evidence_clicks_vs_impressions
        CHECK (
            impressions IS NULL
            OR clicks IS NULL
            OR clicks <= impressions
        ),

    CONSTRAINT chk_seo_query_evidence_ctr
        CHECK (ctr IS NULL OR (ctr >= 0 AND ctr <= 1)),

    CONSTRAINT chk_seo_query_evidence_position
        CHECK (avg_position IS NULL OR avg_position > 0),

    CONSTRAINT chk_seo_query_evidence_volume
        CHECK (search_volume IS NULL OR search_volume >= 0),

    CONSTRAINT chk_seo_query_evidence_confidence
        CHECK (confidence IN (
            'exploratory',
            'supported',
            'validated'
        )),

    CONSTRAINT chk_seo_query_evidence_status
        CHECK (status IN (
            'active',
            'superseded',
            'rejected'
        ))
);

-- Không ingest cùng một canonical source record nhiều lần.
-- Rejected record được phép thay thế bằng một observation mới.
CREATE UNIQUE INDEX uq_seo_query_evidence_source_record
ON seo_query_evidence (
    user_need_id,
    source_system,
    source_record_key
)
WHERE status <> 'rejected';

CREATE INDEX idx_seo_query_evidence_need
ON seo_query_evidence(user_need_id, status);

CREATE INDEX idx_seo_query_evidence_objective
ON seo_query_evidence(objective_id, status);

CREATE INDEX idx_seo_query_evidence_query_hash
ON seo_query_evidence(query_hash);

CREATE INDEX idx_seo_query_evidence_source
ON seo_query_evidence(source_system, source_kind, observation_type);

CREATE INDEX idx_seo_query_evidence_observed
ON seo_query_evidence(observed_to DESC);


-- ============================================================
-- 2. INTENT HYPOTHESIS
-- ============================================================
CREATE TABLE seo_intent_hypothesis (
    id BIGSERIAL PRIMARY KEY,

    objective_id BIGINT NOT NULL,
    user_need_id BIGINT NOT NULL,

    code VARCHAR(120) NOT NULL,
    name TEXT NOT NULL,

    goal_statement TEXT NOT NULL DEFAULT '',
    primary_question TEXT NOT NULL DEFAULT '',
    expected_result TEXT NOT NULL DEFAULT '',

    search_action VARCHAR(24) NOT NULL DEFAULT 'unknown',
    subject_type VARCHAR(32) NOT NULL DEFAULT 'unknown',
    scope_type VARCHAR(24) NOT NULL DEFAULT 'unknown',

    status VARCHAR(24) NOT NULL DEFAULT 'draft',
    confidence VARCHAR(24) NOT NULL DEFAULT 'exploratory',

    owner_code VARCHAR(120) NOT NULL DEFAULT '',
    review_due_at TIMESTAMPTZ,
    lock_version INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uq_seo_intent_identity
        UNIQUE (id, user_need_id, objective_id),

    CONSTRAINT fk_seo_intent_need
        FOREIGN KEY (user_need_id, objective_id)
        REFERENCES seo_user_need(id, objective_id)
        ON DELETE CASCADE,

    CONSTRAINT chk_seo_intent_code
        CHECK (LENGTH(TRIM(code)) > 0),

    CONSTRAINT chk_seo_intent_name
        CHECK (LENGTH(TRIM(name)) > 0),

    CONSTRAINT chk_seo_intent_search_action
        CHECK (search_action IN (
            'unknown',
            'learn',
            'find',
            'compare',
            'verify',
            'transact',
            'monitor',
            'navigate'
        )),

    CONSTRAINT chk_seo_intent_subject_type
        CHECK (subject_type IN (
            'unknown',
            'parcel',
            'administrative_unit',
            'planning_region',
            'planning_document',
            'map_layer',
            'project',
            'report',
            'mixed'
        )),

    CONSTRAINT chk_seo_intent_scope_type
        CHECK (scope_type IN (
            'unknown',
            'specific_entity',
            'area',
            'topic',
            'task',
            'mixed'
        )),

    CONSTRAINT chk_seo_intent_status
        CHECK (status IN ('draft', 'active', 'paused', 'archived')),

    CONSTRAINT chk_seo_intent_confidence
        CHECK (confidence IN ('exploratory', 'supported', 'validated')),

    CONSTRAINT chk_seo_intent_lock_version
        CHECK (lock_version > 0)
);

CREATE UNIQUE INDEX uq_seo_intent_code_active
ON seo_intent_hypothesis(user_need_id, code)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_intent_need
ON seo_intent_hypothesis(user_need_id, status)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_intent_objective
ON seo_intent_hypothesis(objective_id, status)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_intent_action
ON seo_intent_hypothesis(search_action, subject_type, scope_type)
WHERE deleted_at IS NULL;


-- ============================================================
-- 3. INTENT <-> QUERY EVIDENCE
-- ============================================================
CREATE TABLE seo_intent_evidence_link (
    objective_id BIGINT NOT NULL,
    user_need_id BIGINT NOT NULL,

    intent_hypothesis_id BIGINT NOT NULL,
    query_evidence_id BIGINT NOT NULL,

    relation_type VARCHAR(20) NOT NULL DEFAULT 'supports',
    relevance_note TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (intent_hypothesis_id, query_evidence_id),

    CONSTRAINT chk_seo_intent_evidence_relation
        CHECK (relation_type IN ('supports', 'contradicts', 'context')),

    CONSTRAINT fk_seo_intent_evidence_intent
        FOREIGN KEY (intent_hypothesis_id, user_need_id, objective_id)
        REFERENCES seo_intent_hypothesis(id, user_need_id, objective_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_seo_intent_evidence_query
        FOREIGN KEY (query_evidence_id, user_need_id, objective_id)
        REFERENCES seo_query_evidence(id, user_need_id, objective_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_seo_intent_evidence_query
ON seo_intent_evidence_link(query_evidence_id);

CREATE INDEX idx_seo_intent_evidence_need
ON seo_intent_evidence_link(user_need_id);


-- ============================================================
-- 4. QUERY CLUSTER
-- ============================================================
CREATE TABLE seo_query_cluster (
    id BIGSERIAL PRIMARY KEY,

    objective_id BIGINT NOT NULL,
    user_need_id BIGINT NOT NULL,
    intent_hypothesis_id BIGINT NOT NULL,

    code VARCHAR(120) NOT NULL,
    name TEXT NOT NULL,

    search_promise TEXT NOT NULL DEFAULT '',
    inclusion_rule TEXT NOT NULL DEFAULT '',
    exclusion_rule TEXT NOT NULL DEFAULT '',

    status VARCHAR(24) NOT NULL DEFAULT 'draft',
    confidence VARCHAR(24) NOT NULL DEFAULT 'exploratory',

    owner_code VARCHAR(120) NOT NULL DEFAULT '',
    review_due_at TIMESTAMPTZ,
    lock_version INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uq_seo_query_cluster_identity
        UNIQUE (id, user_need_id, objective_id),

    CONSTRAINT fk_seo_query_cluster_intent
        FOREIGN KEY (intent_hypothesis_id, user_need_id, objective_id)
        REFERENCES seo_intent_hypothesis(id, user_need_id, objective_id)
        ON DELETE CASCADE,

    CONSTRAINT chk_seo_query_cluster_code
        CHECK (LENGTH(TRIM(code)) > 0),

    CONSTRAINT chk_seo_query_cluster_name
        CHECK (LENGTH(TRIM(name)) > 0),

    CONSTRAINT chk_seo_query_cluster_status
        CHECK (status IN ('draft', 'active', 'paused', 'archived')),

    CONSTRAINT chk_seo_query_cluster_confidence
        CHECK (confidence IN ('exploratory', 'supported', 'validated')),

    CONSTRAINT chk_seo_query_cluster_lock_version
        CHECK (lock_version > 0)
);

CREATE UNIQUE INDEX uq_seo_query_cluster_code_active
ON seo_query_cluster(user_need_id, code)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_query_cluster_intent
ON seo_query_cluster(intent_hypothesis_id, status)
WHERE deleted_at IS NULL;

CREATE INDEX idx_seo_query_cluster_need
ON seo_query_cluster(user_need_id, status)
WHERE deleted_at IS NULL;


-- ============================================================
-- 5. QUERY CLUSTER MEMBERS
-- ============================================================
CREATE TABLE seo_query_cluster_member (
    objective_id BIGINT NOT NULL,
    user_need_id BIGINT NOT NULL,

    cluster_id BIGINT NOT NULL,
    query_evidence_id BIGINT NOT NULL,

    membership_type VARCHAR(20) NOT NULL,
    rationale TEXT NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (cluster_id, query_evidence_id),

    CONSTRAINT chk_seo_query_cluster_member_type
        CHECK (membership_type IN ('primary', 'included', 'excluded')),

    CONSTRAINT chk_seo_query_cluster_member_status
        CHECK (status IN ('active', 'removed')),

    CONSTRAINT chk_seo_cluster_excluded_rationale
        CHECK (
            membership_type <> 'excluded'
            OR LENGTH(TRIM(rationale)) > 0
        ),

    CONSTRAINT fk_seo_cluster_member_cluster
        FOREIGN KEY (cluster_id, user_need_id, objective_id)
        REFERENCES seo_query_cluster(id, user_need_id, objective_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_seo_cluster_member_query
        FOREIGN KEY (query_evidence_id, user_need_id, objective_id)
        REFERENCES seo_query_evidence(id, user_need_id, objective_id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_seo_cluster_one_primary
ON seo_query_cluster_member(cluster_id)
WHERE membership_type = 'primary' AND status = 'active';

CREATE UNIQUE INDEX uq_seo_query_one_active_cluster
ON seo_query_cluster_member(query_evidence_id)
WHERE membership_type IN ('primary', 'included') AND status = 'active';

CREATE INDEX idx_seo_cluster_member_query
ON seo_query_cluster_member(query_evidence_id, status);

CREATE INDEX idx_seo_cluster_member_need
ON seo_query_cluster_member(user_need_id, status);

COMMIT;
