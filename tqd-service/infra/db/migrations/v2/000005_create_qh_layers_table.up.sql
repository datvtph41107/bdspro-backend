-- =====================================================
-- MIGRATION UP: CREATE qh_layers
-- =====================================================

CREATE TABLE IF NOT EXISTS qh_layers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,

    type INT NOT NULL,
    status INT NOT NULL DEFAULT 10,

    display_order INT NOT NULL DEFAULT 0,
    visible INT NOT NULL DEFAULT 0,

    min_zoom INT NOT NULL DEFAULT 0,
    max_zoom INT NOT NULL DEFAULT 22,

    layer_url TEXT NOT NULL,
    avatar TEXT,
    image_url TEXT,
    thumbnail_url TEXT,

    style_config JSONB NOT NULL DEFAULT '{}',

    legal_doc TEXT,
    legal_doc_url TEXT,
    document_number VARCHAR(100),
    issuing_authority VARCHAR(255),

    effective_date DATE,
    expiry_date DATE,

    created_by BIGINT NOT NULL DEFAULT 0,
    updated_by BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- =====================================================
-- INDEXES
-- =====================================================

-- CREATE UNIQUE INDEX IF NOT EXISTS idx_qh_layers_name
-- ON qh_layers(name)
-- WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_layers_client
ON qh_layers(display_order, id)
WHERE status = 10 AND visible = 1 AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_layers_admin_status
ON qh_layers(status, type, visible)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_layers_dates
ON qh_layers(effective_date, expiry_date)
WHERE status = 10;

CREATE INDEX IF NOT EXISTS idx_qh_layers_search
ON qh_layers USING GIN (
    to_tsvector('simple', name || ' ' || display_name || ' ' || COALESCE(description, ''))
)
WHERE deleted_at IS NULL;

-- =====================================================
-- FUNCTION + TRIGGER
-- =====================================================

CREATE OR REPLACE FUNCTION update_qh_layers_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_qh_layers_updated_at ON qh_layers;

CREATE TRIGGER trg_qh_layers_updated_at
BEFORE UPDATE ON qh_layers
FOR EACH ROW
EXECUTE FUNCTION update_qh_layers_updated_at();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE qh_layers IS 'Bảng lưu trữ thông tin các layer quy hoạch';
COMMENT ON COLUMN qh_layers.name IS 'Tên kỹ thuật của layer (unique)';
COMMENT ON COLUMN qh_layers.display_name IS 'Tên hiển thị cho người dùng';
COMMENT ON COLUMN qh_layers.type IS '1-Quy hoạch tổng thể, 2-Phân khu, 3-Sử dụng đất, 4-Hạ tầng, 5-Môi trường';
COMMENT ON COLUMN qh_layers.status IS '1-Nháp, 10-Hiệu lực, 20-Hết hạn, 30-Lưu trữ';
COMMENT ON COLUMN qh_layers.visible IS '0-Ẩn, 1-Hiển thị';
COMMENT ON COLUMN qh_layers.layer_url IS 'URL GeoJSON hoặc PMTiles';
COMMENT ON COLUMN qh_layers.style_config IS 'JSON config style';