DROP INDEX IF EXISTS idx_qh_labels_standard_at;

ALTER TABLE qh_labels
    DROP COLUMN IF EXISTS standard_at;
