ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS unit_type INTEGER NOT NULL DEFAULT 0;
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS data_source INTEGER NOT NULL DEFAULT 0;
ALTER TABLE province_v2 ADD COLUMN IF NOT EXISTS plan_doc_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_province_v2_unit_type
    ON province_v2(unit_type)
    WHERE unit_type > 0;

CREATE INDEX IF NOT EXISTS idx_province_v2_data_source
    ON province_v2(data_source)
    WHERE data_source > 0;

CREATE INDEX IF NOT EXISTS idx_province_v2_plan_doc_id
    ON province_v2(plan_doc_id)
    WHERE plan_doc_id IS NOT NULL AND plan_doc_id > 0;

COMMENT ON COLUMN province_v2.unit_type IS 'Loại đơn vị: 10 Tỉnh, 20 Thành phố, 30 Đặc khu';
COMMENT ON COLUMN province_v2.data_source IS 'Nguồn dữ liệu: 10 Nội bộ, 20 Chính thức, 30 GIS, 40 Import';
COMMENT ON COLUMN province_v2.plan_doc_id IS 'Tài liệu quy hoạch liên kết (qh_planning_documents.id)';
