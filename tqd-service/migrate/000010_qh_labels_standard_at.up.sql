-- Thời điểm đánh dấu nhãn là tiêu chuẩn hệ thống (NULL = không phải nhãn chuẩn)
ALTER TABLE qh_labels
    ADD COLUMN IF NOT EXISTS standard_at TIMESTAMPTZ;

COMMENT ON COLUMN qh_labels.standard_at IS 'Thời điểm gán làm nhãn tiêu chuẩn hệ thống';

CREATE INDEX IF NOT EXISTS idx_qh_labels_standard_at
    ON qh_labels(standard_at)
    WHERE deleted_at IS NULL AND standard_at IS NOT NULL;
