DROP INDEX IF EXISTS idx_qh_layers_ward_id;
DROP INDEX IF EXISTS idx_qh_layers_province_id;

ALTER TABLE qh_layers
    DROP COLUMN IF EXISTS ward_id,
    DROP COLUMN IF EXISTS province_id;
