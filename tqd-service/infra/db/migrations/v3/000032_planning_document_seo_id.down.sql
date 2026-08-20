DROP INDEX IF EXISTS idx_qh_planning_documents_seo_id;

ALTER TABLE qh_planning_documents
    DROP COLUMN IF EXISTS seo_id;
