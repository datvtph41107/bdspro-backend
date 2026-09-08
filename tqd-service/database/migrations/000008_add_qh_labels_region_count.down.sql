DROP TRIGGER IF EXISTS trg_qh_regions_label_region_count ON qh_regions;
DROP FUNCTION IF EXISTS trg_qh_regions_maintain_label_region_count();
DROP FUNCTION IF EXISTS qh_label_recount_region_count(bigint);
ALTER TABLE qh_labels DROP COLUMN IF EXISTS region_count;
