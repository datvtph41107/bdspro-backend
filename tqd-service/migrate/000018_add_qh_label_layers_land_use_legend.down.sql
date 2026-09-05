DROP INDEX IF EXISTS idx_qh_label_layers_legend_id;
DROP INDEX IF EXISTS idx_qh_label_layers_land_use_id;

ALTER TABLE qh_label_layers
    DROP COLUMN IF EXISTS legend_id,
    DROP COLUMN IF EXISTS land_use_id;
