DROP INDEX IF EXISTS idx_ward_v2_data_source;
DROP INDEX IF EXISTS idx_ward_v2_unit_type;

ALTER TABLE ward_v2 DROP COLUMN IF EXISTS data_source;
ALTER TABLE ward_v2 DROP COLUMN IF EXISTS unit_type;
