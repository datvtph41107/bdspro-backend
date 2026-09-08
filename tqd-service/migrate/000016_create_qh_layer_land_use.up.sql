-- The legacy v2 history created the join table before it ever materialized
-- the GORM-owned land-use table.  Keep the canonical sequence runnable on a
-- fresh database by establishing that prerequisite here.
CREATE TABLE IF NOT EXISTS qh_land_use (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    note TEXT,
    name VARCHAR(100),
    code VARCHAR(10),
    display_order BIGINT DEFAULT 0,
    description TEXT,
    is_visible BOOLEAN DEFAULT TRUE,
    priority BIGINT DEFAULT 50,
    build_condition TEXT,
    color VARCHAR(10),
    can_build BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    warn_level BIGINT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_qh_land_use_deleted_at
    ON qh_land_use (deleted_at);

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
