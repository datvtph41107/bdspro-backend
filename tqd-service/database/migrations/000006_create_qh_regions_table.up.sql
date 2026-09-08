-- =====================================================
-- MIGRATION: TẠO BẢNG QH_REGIONS
-- =====================================================

CREATE TABLE IF NOT EXISTS qh_regions (
    -- Identifiers
    id bigserial PRIMARY KEY,
    layer_id bigint NOT NULL REFERENCES qh_layers(id) ON DELETE CASCADE,
    
    -- Basic info
    name varchar(255),
    display_name varchar(255),
    description text,
    
    -- Land use classification (quan trọng)
    land_use_code varchar(20),
    land_use_name varchar(255),
    
    -- Geometry (dữ liệu vẽ chính)
    geometry geometry(MultiPolygon, 4326) NOT NULL,
    centroid geometry(Point, 4326),
    area_sqm float DEFAULT 0,
    area_ha float DEFAULT 0,
    bbox geometry(Polygon, 4326),
    
    -- Colors from NDJSON
    color_red smallint DEFAULT 200,
    color_green smallint DEFAULT 200,
    color_blue smallint DEFAULT 200,
    
    -- Style override
    style_config jsonb DEFAULT '{}',
    
    -- Metadata
    raw_properties jsonb,
    
    -- Status
    status int DEFAULT 10,
    version int DEFAULT 1,
    is_latest boolean DEFAULT true,
    
    -- Audit
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    deleted_at timestamptz,
    
    -- Constraints
    CONSTRAINT chk_rgb CHECK (
        color_red BETWEEN 0 AND 255 AND
        color_green BETWEEN 0 AND 255 AND
        color_blue BETWEEN 0 AND 255
    ),
    CONSTRAINT chk_area_positive CHECK (area_sqm >= 0)
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Spatial index (quan trọng nhất)
CREATE INDEX IF NOT EXISTS idx_qh_regions_geometry
ON qh_regions USING GIST(geometry) WHERE deleted_at IS NULL;

-- Index for layer lookup
CREATE INDEX IF NOT EXISTS idx_qh_regions_layer_id
ON qh_regions(layer_id) WHERE deleted_at IS NULL;

-- Index for land use filter
CREATE INDEX IF NOT EXISTS idx_qh_regions_land_use
ON qh_regions(land_use_code) WHERE deleted_at IS NULL;

-- Composite index cho client query
CREATE INDEX IF NOT EXISTS idx_qh_regions_client
ON qh_regions(layer_id, land_use_code, status) 
WHERE status = 10 AND is_latest = true AND deleted_at IS NULL;

-- =====================================================
-- TRIGGERS
-- =====================================================

-- Auto update updated_at
CREATE OR REPLACE FUNCTION update_qh_regions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_qh_regions_updated_at
BEFORE UPDATE ON qh_regions
FOR EACH ROW
EXECUTE FUNCTION update_qh_regions_updated_at();

-- Auto calculate metrics from geometry
CREATE OR REPLACE FUNCTION calculate_qh_region_metrics()
RETURNS TRIGGER AS $$
BEGIN
    -- Tính centroid
    NEW.centroid := ST_Centroid(NEW.geometry);
    
    -- Tính diện tích (chuyển sang Web Mercator để lấy mét vuông)
    NEW.area_sqm := ST_Area(ST_Transform(NEW.geometry, 3857));
    NEW.area_ha := NEW.area_sqm / 10000;
    
    -- Tính bounding box
    NEW.bbox := ST_Envelope(NEW.geometry);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_qh_regions_calculate_metrics
BEFORE INSERT OR UPDATE OF geometry ON qh_regions
FOR EACH ROW
EXECUTE FUNCTION calculate_qh_region_metrics();

-- =====================================================
-- COMMENTS
-- =====================================================

COMMENT ON TABLE qh_regions IS 'Khu vực quy hoạch chi tiết, chứa hình học thực tế';
COMMENT ON COLUMN qh_regions.layer_id IS 'Thuộc layer quy hoạch nào';
COMMENT ON COLUMN qh_regions.land_use_code IS 'Loại đất theo quy hoạch (ODT, TMD, DGT...)';
COMMENT ON COLUMN qh_regions.geometry IS 'Hình học MultiPolygon của khu vực';
COMMENT ON COLUMN qh_regions.area_ha IS 'Diện tích tính sẵn (hecta)';
COMMENT ON COLUMN qh_regions.style_config IS 'Style override cho region này';


-- =====================================================
-- FUNCTION: IMPORT_REGION_FROM_NDJSON
-- =====================================================

CREATE OR REPLACE FUNCTION import_region_from_geojson(
    p_feature_json jsonb,
    p_layer_id bigint
)
RETURNS bigint AS $$
DECLARE
    v_geometry geometry;
    v_properties jsonb;
    v_land_use_code varchar(20);
    v_land_use_name text;
    v_red smallint;
    v_green smallint;
    v_blue smallint;
    v_region_id bigint;
    v_geometry_text text;
BEGIN
    -- Extract geometry
    v_geometry_text := p_feature_json->'geometry'::text;
    v_geometry := ST_SetSRID(ST_GeomFromGeoJSON(v_geometry_text), 4326);
    
    -- Extract properties
    v_properties := COALESCE(p_feature_json->'properties', '{}');
    
    -- Extract land use info
    v_land_use_name := TRIM(COALESCE(v_properties->>'loai_dat_quy_hoach', 'UNKNOWN'));
    v_land_use_code := COALESCE(v_properties->>'land_use_code', 'UNKNOWN');
    
    -- Extract colors (mặc định nếu không có)
    v_red := COALESCE((v_properties->>'red')::smallint, 200);
    v_green := COALESCE((v_properties->>'green')::smallint, 200);
    v_blue := COALESCE((v_properties->>'blue')::smallint, 200);
    
    -- Insert region
    INSERT INTO qh_regions (
        layer_id,
        name,
        display_name,
        land_use_code,
        land_use_name,
        geometry,
        color_red,
        color_green,
        color_blue,
        raw_properties,
        status
    ) VALUES (
        p_layer_id,
        v_land_use_name,
        INITCAP(v_land_use_name),
        v_land_use_code,
        v_land_use_name,
        v_geometry,
        v_red,
        v_green,
        v_blue,
        v_properties,
        10
    )
    RETURNING id INTO v_region_id;
    
    RETURN v_region_id;
    
EXCEPTION WHEN OTHERS THEN
    RAISE WARNING 'Failed to import region: %', SQLERRM;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
