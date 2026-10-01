-- Align seo_category with CRM canonical schema (metadata jsonb).
-- Safe for tables that already have rows: add nullable → backfill → NOT NULL + default.

ALTER TABLE public.seo_category
    ADD COLUMN IF NOT EXISTS metadata jsonb;

UPDATE public.seo_category
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

ALTER TABLE public.seo_category
    ALTER COLUMN metadata SET DEFAULT '{}'::jsonb;

ALTER TABLE public.seo_category
    ALTER COLUMN metadata SET NOT NULL;

COMMENT ON COLUMN public.seo_category.metadata IS
    'JSON mở rộng cho chủ đề (tương thích crm-service seo_category.metadata).';
