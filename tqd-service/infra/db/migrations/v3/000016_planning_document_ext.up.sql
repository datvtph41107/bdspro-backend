-- 000016_planning_document_ext.up.sql
-- Lưu định dạng file (pdf, dwg, ...) trên qh_planning_documents.

ALTER TABLE qh_planning_documents
    ADD COLUMN IF NOT EXISTS ext VARCHAR(20) NOT NULL DEFAULT '';

UPDATE qh_planning_documents
SET ext = LOWER(regexp_replace(
        COALESCE(NULLIF(BTRIM(relative_path), ''), NULLIF(BTRIM(filepath), ''), ''),
        '^.*\.', ''
    ))
WHERE deleted_at IS NULL
  AND ext = ''
  AND COALESCE(NULLIF(BTRIM(relative_path), ''), NULLIF(BTRIM(filepath), '')) ~ '\.[^./\\]+$';
