DROP INDEX IF EXISTS idx_qh_layers_replaced_by_layer_id;
ALTER TABLE qh_layers DROP COLUMN IF EXISTS replaced_by_layer_id;
