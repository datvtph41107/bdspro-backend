CREATE TABLE IF NOT EXISTS qh_layer_land_use (
    layer_id    BIGINT NOT NULL REFERENCES qh_layers(id) ON DELETE CASCADE,
    land_use_id BIGINT NOT NULL REFERENCES qh_land_use(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (layer_id, land_use_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_layer_land_use_layer_id
    ON qh_layer_land_use (layer_id);

CREATE INDEX IF NOT EXISTS idx_qh_layer_land_use_land_use_id
    ON qh_layer_land_use (land_use_id);

COMMENT ON TABLE qh_layer_land_use IS 'Liên kết nhiều-nhiều: layer ↔ loại đất (qh_land_use)';
