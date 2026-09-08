-- Nhóm layer (family) + FK từ qh_layers
CREATE TABLE IF NOT EXISTS qh_layer_families (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL DEFAULT '',
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_qh_layer_families_deleted
    ON qh_layer_families (deleted_at);

COMMENT ON TABLE qh_layer_families IS 'Nhóm layer quy hoạch (gán nhiều layer vào cùng family)';

ALTER TABLE qh_layers
    ADD COLUMN IF NOT EXISTS family_id BIGINT REFERENCES qh_layer_families (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_qh_layers_family_id
    ON qh_layers (family_id)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN qh_layers.family_id IS 'FK → qh_layer_families.id (nhóm layer)';
