ALTER TABLE qh_label_layers
    ADD COLUMN IF NOT EXISTS land_use_id BIGINT,
    ADD COLUMN IF NOT EXISTS legend_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_qh_label_layers_land_use_id
    ON qh_label_layers (land_use_id);

CREATE INDEX IF NOT EXISTS idx_qh_label_layers_legend_id
    ON qh_label_layers (legend_id);
