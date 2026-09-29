-- Layer thay thế: lưu id layer mới trên bản ghi cũ để truy vết
ALTER TABLE qh_layers ADD COLUMN IF NOT EXISTS replaced_by_layer_id BIGINT REFERENCES qh_layers (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_qh_layers_replaced_by_layer_id
ON qh_layers (replaced_by_layer_id)
WHERE deleted_at IS NULL AND replaced_by_layer_id IS NOT NULL;

COMMENT ON COLUMN qh_layers.replaced_by_layer_id IS 'ID layer mới thay thế layer này (khi legal_status = replaced)';
