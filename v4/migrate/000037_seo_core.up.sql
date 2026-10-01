-- 000037_seo_core.up.sql
-- Independent v4 clone of CRM SEO tables — v4 DSN, no crm-service dependency.

-- ── seo_domain ────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_domain (
    id                  BIGSERIAL PRIMARY KEY,
    slug                VARCHAR(255),
    origin_url          TEXT,
    canonical_url       TEXT NOT NULL DEFAULT '',

    ref_type            SMALLINT NOT NULL DEFAULT 0,
    ref_id              BIGINT,
    ref_source          VARCHAR(80) NOT NULL DEFAULT 'manual',
    ref_label           TEXT,
    ref_url             TEXT,
    source_status       VARCHAR(40) NOT NULL DEFAULT 'manual',
    ref_missing         BOOLEAN NOT NULL DEFAULT FALSE,
    ref_snapshot_json   JSONB NOT NULL DEFAULT '{}'::JSONB,
    ref_hash            VARCHAR(64),
    ref_last_synced_at  TIMESTAMPTZ,

    scope               VARCHAR(50) NOT NULL DEFAULT 'public',
    page_status         VARCHAR(30) NOT NULL DEFAULT 'draft',

    title               TEXT,
    description         TEXT,
    content             TEXT,
    summary             TEXT,

    visible_on          SMALLINT NOT NULL DEFAULT 30,
    category_id         BIGINT,
    province_id         BIGINT,
    image_url           TEXT,

    published           BOOLEAN NOT NULL DEFAULT FALSE,
    published_at        TIMESTAMPTZ,

    is_site_map         BOOLEAN NOT NULL DEFAULT FALSE,
    is_index            BOOLEAN NOT NULL DEFAULT FALSE,
    is_robot            BOOLEAN DEFAULT FALSE,
    site_map_lasted_at  TIMESTAMPTZ,
    sitemap_priority    REAL NOT NULL DEFAULT 0.5,
    sitemap_change_freq VARCHAR(20) DEFAULT 'daily',

    need_generate       BOOLEAN NOT NULL DEFAULT TRUE,
    generated_at        TIMESTAMPTZ,
    source_updated_at   TIMESTAMPTZ,

    render_status       VARCHAR(30) NOT NULL DEFAULT 'none',
    rendered_html       TEXT,
    last_render_error   TEXT,
    static_html_path    TEXT,
    static_html_hash    VARCHAR(64),

    template_key        VARCHAR(120),
    template_version    VARCHAR(80),
    deep_link           TEXT,
    metadata            JSONB DEFAULT '{}'::JSONB,
    note                TEXT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_seo_domain_page_status   ON seo_domain (page_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_published      ON seo_domain (published)   WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_slug           ON seo_domain (slug)        WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_ref            ON seo_domain (ref_type, ref_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_category       ON seo_domain (category_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_province       ON seo_domain (province_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_visible_on     ON seo_domain (visible_on)  WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_domain_deleted_at     ON seo_domain (deleted_at);
CREATE INDEX IF NOT EXISTS idx_seo_domain_canonical_url  ON seo_domain (canonical_url) WHERE deleted_at IS NULL;

-- ── seo_links (internal links) ────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_links (
    id              BIGSERIAL PRIMARY KEY,
    parent_seo_id   BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    child_seo_id    BIGINT REFERENCES seo_domain(id) ON DELETE SET NULL,
    href            TEXT NOT NULL DEFAULT '',
    anchor_text     TEXT,
    priority        INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_seo_links_parent ON seo_links (parent_seo_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seo_links_child  ON seo_links (child_seo_id)  WHERE deleted_at IS NULL;

-- ── seo_relative (related pages) ──────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_relative (
    id              BIGSERIAL PRIMARY KEY,
    parent_seo_id   BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    child_seo_id    BIGINT REFERENCES seo_domain(id) ON DELETE SET NULL,
    relation_type   VARCHAR(60),
    priority        INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_seo_relative_parent ON seo_relative (parent_seo_id) WHERE deleted_at IS NULL;

-- ── seo_generation_log ────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_generation_log (
    id              BIGSERIAL PRIMARY KEY,
    seo_domain_id   BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    status          VARCHAR(30) NOT NULL DEFAULT 'pending',
    message         TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_seo_gen_log_domain ON seo_generation_log (seo_domain_id);
CREATE INDEX IF NOT EXISTS idx_seo_gen_log_status ON seo_generation_log (status);

-- ── seo_count ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_count (
    id              BIGSERIAL PRIMARY KEY,
    seo_domain_id   BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    click_count     BIGINT NOT NULL DEFAULT 0,
    view_count      BIGINT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (seo_domain_id)
);

-- ── seo_count_log ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_count_log (
    id              BIGSERIAL PRIMARY KEY,
    seo_domain_id   BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    event_type      VARCHAR(20) NOT NULL DEFAULT 'click',
    user_id         BIGINT,
    ip_addr         VARCHAR(45),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_seo_count_log_domain ON seo_count_log (seo_domain_id);
CREATE INDEX IF NOT EXISTS idx_seo_count_log_date   ON seo_count_log (created_at);

-- ── seo_category ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS seo_category (
    id          BIGSERIAL PRIMARY KEY,
    code        VARCHAR(80) NOT NULL,
    slug        VARCHAR(120) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    sort_order  INT NOT NULL DEFAULT 0,
    metadata    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS idx_seo_category_status ON seo_category (status) WHERE deleted_at IS NULL;

-- Seed default categories
INSERT INTO seo_category (code, slug, name, description, status, sort_order) VALUES
    ('planning_news',    'tin-tuc-quy-hoach',    'Tin tức quy hoạch',     'Tin tức và thông báo liên quan đến quy hoạch', 'active', 10),
    ('planning_project', 'du-an-quy-hoach',       'Đồ án quy hoạch',       'Các đồ án quy hoạch đô thị và nông thôn',      'active', 20),
    ('legal_document',   'van-ban-phap-ly',        'Văn bản pháp lý',       'Văn bản và tài liệu pháp lý quy hoạch',        'active', 30),
    ('analysis_report',  'bao-cao-phan-tich',      'Báo cáo & Phân tích',   'Báo cáo nghiên cứu và phân tích thị trường',   'active', 40),
    ('planning_region',  'vung-quy-hoach',         'Vùng quy hoạch',        'Thông tin quy hoạch theo vùng địa lý',         'active', 50)
ON CONFLICT (code) DO NOTHING;

-- ── user_saved_planning_news ──────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_saved_planning_news (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    seo_domain_id   BIGINT NOT NULL REFERENCES seo_domain(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE (user_id, seo_domain_id)
);

CREATE INDEX IF NOT EXISTS idx_saved_news_user   ON user_saved_planning_news (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_saved_news_domain ON user_saved_planning_news (seo_domain_id);
