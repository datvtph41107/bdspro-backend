-- =====================================================
-- MIGRATION UP: CREATE qh_labels
-- =====================================================

CREATE TABLE IF NOT EXISTS qh_labels (
    id          BIGSERIAL PRIMARY KEY,
    layer_id    BIGINT       NOT NULL REFERENCES qh_layers(id),
    name        VARCHAR(255) NOT NULL,
    color       VARCHAR(20)  NOT NULL DEFAULT '#1890ff',
    description TEXT,
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_qh_labels_layer_id
ON qh_labels(layer_id, sort_order)
WHERE deleted_at IS NULL;

-- =====================================================
-- TRIGGER updated_at
-- =====================================================

CREATE OR REPLACE FUNCTION update_qh_labels_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_qh_labels_updated_at ON qh_labels;

CREATE TRIGGER trg_qh_labels_updated_at
BEFORE UPDATE ON qh_labels
FOR EACH ROW
EXECUTE FUNCTION update_qh_labels_updated_at();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE  qh_labels            IS 'Nhãn phân loại gắn vào layer quy hoạch';
COMMENT ON COLUMN qh_labels.layer_id   IS 'FK → qh_layers.id';
COMMENT ON COLUMN qh_labels.name       IS 'Tên nhãn hiển thị';
COMMENT ON COLUMN qh_labels.color      IS 'Mã màu hex (vd: #1890ff)';
COMMENT ON COLUMN qh_labels.sort_order IS 'Thứ tự hiển thị, nhỏ hơn lên trước';
