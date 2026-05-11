DROP INDEX IF EXISTS idx_qh_layers_import_batch_id;
ALTER TABLE qh_layers DROP COLUMN IF EXISTS import_batch_id;
ALTER TABLE qh_layers DROP COLUMN IF EXISTS import_status;
