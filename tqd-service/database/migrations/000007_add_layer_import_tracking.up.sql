-- Layer import async tracking (align with GORM model qh_layers)
ALTER TABLE qh_layers ADD COLUMN IF NOT EXISTS import_status INT NOT NULL DEFAULT 20;
ALTER TABLE qh_layers ADD COLUMN IF NOT EXISTS import_batch_id VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_qh_layers_import_batch_id
ON qh_layers(import_batch_id)
WHERE import_batch_id IS NOT NULL AND deleted_at IS NULL;

COMMENT ON COLUMN qh_layers.import_status IS '10=processing, 20=done, 30=failed';
COMMENT ON COLUMN qh_layers.import_batch_id IS 'Batch id của lần import đang/đã chạy';
