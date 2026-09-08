DROP INDEX IF EXISTS idx_qh_layers_family_id;
ALTER TABLE qh_layers DROP COLUMN IF EXISTS family_id;

DROP INDEX IF EXISTS idx_qh_layer_families_deleted;
DROP TABLE IF EXISTS qh_layer_families;
