DROP INDEX IF EXISTS idx_qh_layer_families_active_sort;

ALTER TABLE qh_layer_families
    DROP COLUMN IF EXISTS sort_number;
