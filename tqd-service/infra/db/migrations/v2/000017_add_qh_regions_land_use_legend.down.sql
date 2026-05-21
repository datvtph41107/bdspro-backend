DROP INDEX IF EXISTS idx_qh_regions_legend_id;
DROP INDEX IF EXISTS idx_qh_regions_land_use_id;

ALTER TABLE qh_regions
    DROP COLUMN IF EXISTS legend_id,
    DROP COLUMN IF EXISTS land_use_id;
