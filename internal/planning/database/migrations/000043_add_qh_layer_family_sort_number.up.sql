ALTER TABLE qh_layer_families
    ADD COLUMN IF NOT EXISTS sort_number INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_qh_layer_families_active_sort
    ON qh_layer_families (sort_number, id)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN qh_layer_families.sort_number IS
    'Stable display order for planning layer families; lower values are shown first to clients';
