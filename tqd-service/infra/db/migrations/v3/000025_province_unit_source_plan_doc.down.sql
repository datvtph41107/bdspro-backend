DROP INDEX IF EXISTS idx_province_v2_plan_doc_id;
DROP INDEX IF EXISTS idx_province_v2_data_source;
DROP INDEX IF EXISTS idx_province_v2_unit_type;

ALTER TABLE province_v2
    DROP COLUMN IF EXISTS plan_doc_id,
    DROP COLUMN IF EXISTS data_source,
    DROP COLUMN IF EXISTS unit_type;
