-- Chú giải: gắn nhãn (label) với layer kèm ghi chú (note)
CREATE TABLE IF NOT EXISTS qh_legends (
    id         BIGSERIAL PRIMARY KEY,
    layer_id   BIGINT NOT NULL REFERENCES qh_layers(id) ON DELETE CASCADE,
    label_id   BIGINT NOT NULL REFERENCES qh_labels(id) ON DELETE CASCADE,
    note       TEXT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_qh_legends_layer_id
    ON qh_legends (layer_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_legends_label_id
    ON qh_legends (label_id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_qh_legends_layer_label_active
    ON qh_legends (layer_id, label_id)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE qh_legends IS 'Chú giải: đính label với layer kèm note';
