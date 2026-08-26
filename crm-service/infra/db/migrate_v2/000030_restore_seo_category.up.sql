BEGIN;

-- 000024 từng DROP seo_category nhầm. Khôi phục bảng chủ đề độc lập
-- (không FK / không phụ thuộc bảng khác) cho seo_domain.category_id.

CREATE TABLE IF NOT EXISTS seo_category (
    id BIGSERIAL PRIMARY KEY,
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
    CONSTRAINT uq_seo_category_code UNIQUE (code),
    CONSTRAINT uq_seo_category_slug UNIQUE (slug),
    CONSTRAINT chk_seo_category_status
        CHECK (status IN ('active', 'inactive', 'archived'))
);

-- Schema cũ (000012/000019) có domain_id + unique theo domain → gỡ hết.
DO $$
BEGIN
    IF to_regclass('public.seo_category') IS NULL THEN
        RETURN;
    END IF;

    ALTER TABLE public.seo_category
        DROP CONSTRAINT IF EXISTS seo_content_category_domain_id_fkey;
    ALTER TABLE public.seo_category
        DROP CONSTRAINT IF EXISTS seo_category_domain_id_fkey;
    ALTER TABLE public.seo_category
        DROP CONSTRAINT IF EXISTS uq_seo_content_category_domain_code;
    ALTER TABLE public.seo_category
        DROP CONSTRAINT IF EXISTS uq_seo_content_category_domain_slug;
    ALTER TABLE public.seo_category
        DROP CONSTRAINT IF EXISTS uq_seo_category_domain_code;
    ALTER TABLE public.seo_category
        DROP CONSTRAINT IF EXISTS uq_seo_category_domain_slug;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'seo_category'
          AND column_name = 'domain_id'
    ) THEN
        ALTER TABLE public.seo_category DROP COLUMN domain_id;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'public.seo_category'::regclass
          AND conname = 'uq_seo_category_code'
    ) THEN
        ALTER TABLE public.seo_category
            ADD CONSTRAINT uq_seo_category_code UNIQUE (code);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'public.seo_category'::regclass
          AND conname = 'uq_seo_category_slug'
    ) THEN
        ALTER TABLE public.seo_category
            ADD CONSTRAINT uq_seo_category_slug UNIQUE (slug);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_seo_category_active
    ON seo_category(status, sort_order, id)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE seo_category IS
'Standalone topic taxonomy. seo_domain.category_id references this by id only.';

INSERT INTO seo_category(code, slug, name, description, sort_order)
VALUES
    ('planning_news', 'tin-quy-hoach', 'Tin quy hoạch', 'Tin tức và cập nhật quy hoạch.', 10),
    ('planning_project', 'do-an-quy-hoach', 'Đồ án quy hoạch', 'Nội dung về đồ án và hồ sơ quy hoạch.', 20),
    ('planning_report', 'bao-cao-quy-hoach', 'Báo cáo quy hoạch', 'Báo cáo và phân tích quy hoạch.', 30),
    ('administrative_update', 'cap-nhat-dia-gioi', 'Cập nhật địa giới', 'Thông tin đơn vị hành chính và địa giới.', 40)
ON CONFLICT (code) DO UPDATE
SET slug = EXCLUDED.slug,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    sort_order = EXCLUDED.sort_order,
    status = 'active',
    updated_at = NOW(),
    deleted_at = NULL;

COMMIT;
