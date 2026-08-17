ALTER TABLE qh_layers
    ADD COLUMN IF NOT EXISTS province_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS ward_id VARCHAR(50);

CREATE INDEX IF NOT EXISTS idx_qh_layers_province_id
    ON qh_layers(province_id)
    WHERE province_id IS NOT NULL AND province_id <> '';

CREATE INDEX IF NOT EXISTS idx_qh_layers_ward_id
    ON qh_layers(ward_id)
    WHERE ward_id IS NOT NULL AND ward_id <> '';

COMMENT ON COLUMN qh_layers.province_id IS 'Đơn vị hành chính tỉnh/thành (province_v2.id)';
COMMENT ON COLUMN qh_layers.ward_id IS 'Đơn vị hành chính phường/xã (ward_v2.id)';
