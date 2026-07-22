BEGIN;

-- seo_domain already means a public SEO page registry in CRM. The following
-- tables intentionally use seo_content_* names so editorial taxonomy does not
-- become a second, conflicting meaning of seo_domain.
CREATE TABLE IF NOT EXISTS seo_content_domain (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(160) NOT NULL,
    description TEXT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_seo_content_domain_status
        CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_seo_content_domain_active
    ON seo_content_domain(status, sort_order, id)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS seo_content_category (
    id BIGSERIAL PRIMARY KEY,
    domain_id BIGINT NOT NULL REFERENCES seo_content_domain(id),
    code VARCHAR(100) NOT NULL,
    slug VARCHAR(140) NOT NULL,
    name VARCHAR(180) NOT NULL,
    description TEXT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    sort_order INTEGER NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_seo_content_category_domain_code UNIQUE(domain_id, code),
    CONSTRAINT uq_seo_content_category_domain_slug UNIQUE(domain_id, slug),
    CONSTRAINT chk_seo_content_category_status
        CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_seo_content_category_domain
    ON seo_content_category(domain_id, status, sort_order, id)
    WHERE deleted_at IS NULL;

-- Generic assignment keeps CRM independent from the physical news/article
-- table. source_system + source_type + source_id identify the real owner.
CREATE TABLE IF NOT EXISTS seo_content_category_assignment (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES seo_content_category(id),
    seo_domain_id BIGINT NULL REFERENCES seo_domain(id),
    source_system VARCHAR(60) NOT NULL,
    source_type VARCHAR(80) NOT NULL,
    source_id VARCHAR(160) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_seo_content_category_assignment
        UNIQUE(category_id, source_system, source_type, source_id)
);

CREATE INDEX IF NOT EXISTS idx_seo_content_category_assignment_source
    ON seo_content_category_assignment(source_system, source_type, source_id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_seo_content_category_assignment_primary
    ON seo_content_category_assignment(source_system, source_type, source_id)
    WHERE is_primary = TRUE AND deleted_at IS NULL;

-- Stores only explicit editorial decisions. latest, most_read and trending are
-- computed feeds and must not be persisted as mutually exclusive enum values.
CREATE TABLE IF NOT EXISTS seo_content_placement (
    id BIGSERIAL PRIMARY KEY,
    seo_domain_id BIGINT NULL REFERENCES seo_domain(id),
    source_system VARCHAR(60) NOT NULL,
    source_type VARCHAR(80) NOT NULL,
    source_id VARCHAR(160) NOT NULL,
    placement_code VARCHAR(40) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by VARCHAR(120) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_seo_content_placement_code
        CHECK (placement_code IN ('featured', 'attention', 'breaking', 'pinned')),
    CONSTRAINT chk_seo_content_placement_status
        CHECK (status IN ('draft', 'active', 'paused', 'expired', 'archived')),
    CONSTRAINT chk_seo_content_placement_window
        CHECK (ends_at IS NULL OR starts_at IS NULL OR ends_at > starts_at),
    CONSTRAINT uq_seo_content_placement
        UNIQUE(source_system, source_type, source_id, placement_code, starts_at)
);

CREATE INDEX IF NOT EXISTS idx_seo_content_placement_feed
    ON seo_content_placement(placement_code, status, priority DESC, starts_at, ends_at)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_seo_content_placement_source
    ON seo_content_placement(source_system, source_type, source_id)
    WHERE deleted_at IS NULL;

INSERT INTO seo_content_domain(code, name, description, sort_order)
VALUES
    ('planning', 'Quy hoạch', 'Nội dung quy hoạch, địa giới, đồ án và dữ liệu không gian.', 10),
    ('infrastructure', 'Hạ tầng', 'Nội dung hạ tầng và tác động phát triển đô thị.', 20),
    ('legal', 'Pháp lý', 'Văn bản, thủ tục và diễn giải pháp lý.', 30),
    ('market', 'Thị trường', 'Phân tích và biến động thị trường.', 40)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    sort_order = EXCLUDED.sort_order,
    updated_at = NOW(),
    deleted_at = NULL;

INSERT INTO seo_content_category(domain_id, code, slug, name, description, sort_order)
SELECT d.id, v.code, v.slug, v.name, v.description, v.sort_order
FROM seo_content_domain d
JOIN (VALUES
    ('planning_news', 'tin-quy-hoach', 'Tin quy hoạch', 'Tin tức và cập nhật quy hoạch.', 10),
    ('planning_project', 'do-an-quy-hoach', 'Đồ án quy hoạch', 'Nội dung về đồ án và hồ sơ quy hoạch.', 20),
    ('planning_report', 'bao-cao-quy-hoach', 'Báo cáo quy hoạch', 'Báo cáo và phân tích quy hoạch.', 30),
    ('administrative_update', 'cap-nhat-dia-gioi', 'Cập nhật địa giới', 'Thông tin đơn vị hành chính và địa giới.', 40)
) AS v(code, slug, name, description, sort_order) ON TRUE
WHERE d.code = 'planning'
ON CONFLICT (domain_id, code) DO UPDATE
SET slug = EXCLUDED.slug,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    sort_order = EXCLUDED.sort_order,
    status = 'active',
    updated_at = NOW(),
    deleted_at = NULL;

COMMIT;
