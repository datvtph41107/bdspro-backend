ALTER TABLE qh_regions
    ADD COLUMN IF NOT EXISTS land_use_id BIGINT,
    ADD COLUMN IF NOT EXISTS legend_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_qh_regions_land_use_id
    ON qh_regions (land_use_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_regions_legend_id
    ON qh_regions (legend_id)
    WHERE deleted_at IS NULL;
