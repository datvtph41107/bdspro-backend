ALTER TABLE qh_planning_documents
    ADD COLUMN IF NOT EXISTS seo_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_qh_planning_documents_seo_id
    ON qh_planning_documents(seo_id)
    WHERE seo_id IS NOT NULL AND seo_id > 0;

COMMENT ON COLUMN qh_planning_documents.seo_id IS 'ID bản ghi SEO gắn với tài liệu quy hoạch.';
