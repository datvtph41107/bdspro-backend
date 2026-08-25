-- Metadata pháp lý / hồ sơ hiển thị trên tài liệu quy hoạch.

ALTER TABLE qh_planning_documents
    ADD COLUMN IF NOT EXISTS indicators JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS authority VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS sector VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS field VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS signer VARCHAR(255) NOT NULL DEFAULT '';

COMMENT ON COLUMN qh_planning_documents.document_type IS 'Loại tài liệu (ví dụ: Bản đồ, Pháp lý, Thuyết minh).';
COMMENT ON COLUMN qh_planning_documents.issue_date IS 'Ngày ban hành.';
COMMENT ON COLUMN qh_planning_documents.effective_date IS 'Ngày hiệu lực.';
COMMENT ON COLUMN qh_planning_documents.validity_status IS 'Tình trạng hiệu lực (còn hiệu lực / hết hiệu lực / ...).';
COMMENT ON COLUMN qh_planning_documents.indicators IS 'Số liệu / chỉ tiêu của tài liệu, lưu JSON array.';
COMMENT ON COLUMN qh_planning_documents.authority IS 'Cơ quan ban hành.';
COMMENT ON COLUMN qh_planning_documents.sector IS 'Ngành (phân loại chuyên môn của tài liệu).';
COMMENT ON COLUMN qh_planning_documents.field IS 'Lĩnh vực (phân loại nội dung của tài liệu).';
COMMENT ON COLUMN qh_planning_documents.signer IS 'Người ký.';
