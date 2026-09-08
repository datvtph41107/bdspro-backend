-- Migration: Tạo các bảng liên quan đến Property (Hồ sơ BĐS)
-- property – Hồ sơ lõi BĐS
-- property_land_info – Thông tin đất (mở rộng)
-- property_building_info – Nhà/công trình (mở rộng; chỉ có khi có nhà)
-- property_media – Ảnh hồ sơ BĐS (không phải ảnh marketing)
-- property_tag_link – Tag tiện ích/đặc điểm (dạng Tag)
-- property_relation – Liên kết vòng đời BĐS với Tài sản/Sản phẩm
-- 1. Tạo bảng property (Hồ sơ lõi BĐS)
CREATE TABLE IF NOT EXISTS property (
    id BIGSERIAL PRIMARY KEY,
    property_type_id BIGINT NULL,
    title VARCHAR(255) NULL,
    location_id BIGINT NULL,
    address_text VARCHAR(500) NULL,
    latitude DECIMAL(10,8) NULL,
    longitude DECIMAL(11,8) NULL,
    map_url VARCHAR(500) NULL,
    project_id BIGINT NULL,
    building_block VARCHAR(100) NULL,
    level VARCHAR(100) NULL,
    unit_code VARCHAR(50) NULL,
    area_total DECIMAL(15,2) NULL,
    area_land DECIMAL(15,2) NULL,
    area_residential DECIMAL(15,2) NULL,
    legal_status VARCHAR(100) NULL,
    legal_note TEXT NULL,
    has_building BOOLEAN DEFAULT FALSE,
    record_status VARCHAR(50) NULL,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

-- Indexes cho bảng property
CREATE INDEX IF NOT EXISTS idx_property_property_type_id ON property(property_type_id);
CREATE INDEX IF NOT EXISTS idx_property_location_id ON property(location_id);
CREATE INDEX IF NOT EXISTS idx_property_project_id ON property(project_id);
CREATE INDEX IF NOT EXISTS idx_property_record_status ON property(record_status);

-- 2. Tạo bảng property_land_info (Thông tin đất - mở rộng)
CREATE TABLE IF NOT EXISTS property_land_info (
    id BIGSERIAL PRIMARY KEY,
    property_id BIGINT NOT NULL,
    frontage DECIMAL(10,2) NULL,
    depth DECIMAL(10,2) NULL,
    road_width DECIMAL(10,2) NULL,
    land_note TEXT NULL,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    
    FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE
);

-- Indexes cho bảng property_land_info
CREATE INDEX IF NOT EXISTS idx_property_land_info_property_id ON property_land_info(property_id);

-- 3. Tạo bảng property_building_info (Nhà/công trình - mở rộng, chỉ có khi có nhà)
CREATE TABLE IF NOT EXISTS property_building_info (
    id BIGSERIAL PRIMARY KEY,
    property_id BIGINT NOT NULL,
    building_type VARCHAR(100) NULL,
    construction_area DECIMAL(15,2) NULL,
    floor_area DECIMAL(15,2) NULL,
    floors SMALLINT NULL,
    bedrooms SMALLINT NULL,
    bathrooms SMALLINT NULL,
    direction VARCHAR(50) NULL,
    balcony_direction VARCHAR(50) NULL,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    
    FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE
);

-- Indexes cho bảng property_building_info
CREATE INDEX IF NOT EXISTS idx_property_building_info_property_id ON property_building_info(property_id);

-- 4. Tạo bảng property_media (Ảnh hồ sơ BĐS - không phải ảnh marketing)
CREATE TABLE IF NOT EXISTS property_media (
    id BIGSERIAL PRIMARY KEY,
    property_id BIGINT NOT NULL,
    media_type VARCHAR(20) NOT NULL,
    url VARCHAR(500) NOT NULL,
    thumb_url VARCHAR(500) NULL,
    is_cover BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    created_by BIGINT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE
);

-- Indexes cho bảng property_media
CREATE INDEX IF NOT EXISTS idx_property_media_property_id ON property_media(property_id);
CREATE INDEX IF NOT EXISTS idx_property_media_is_cover ON property_media(is_cover);

-- 5. Tạo bảng tag (Tag tiện ích/đặc điểm)
CREATE TABLE IF NOT EXISTS tag (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NULL,
    description TEXT NULL,
    icon VARCHAR(255) NULL,
    active BOOLEAN DEFAULT TRUE,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

-- Indexes cho bảng tag
CREATE INDEX IF NOT EXISTS idx_tag_type ON tag(type);
CREATE INDEX IF NOT EXISTS idx_tag_active ON tag(active);

-- 6. Tạo bảng property_tag_link (Tag tiện ích/đặc điểm - liên kết)
CREATE TABLE IF NOT EXISTS property_tag_link (
    property_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    
    PRIMARY KEY (property_id, tag_id),
    FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag(id) ON DELETE CASCADE
);

-- Indexes cho bảng property_tag_link
CREATE INDEX IF NOT EXISTS idx_property_tag_link_property_id ON property_tag_link(property_id);
CREATE INDEX IF NOT EXISTS idx_property_tag_link_tag_id ON property_tag_link(tag_id);

-- 7. Tạo bảng property_relation (Liên kết vòng đời BĐS với Tài sản/Sản phẩm)
-- Lưu ý: asset_id và primary_product_id không có foreign key vì chúng là optional và có thể tham chiếu đến các bảng khác
CREATE TABLE IF NOT EXISTS property_relation (
    property_id BIGINT NOT NULL,
    asset_id BIGINT NULL,
    primary_product_id BIGINT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (property_id),
    FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE
    -- Note: asset_id và primary_product_id không có foreign key constraint vì chúng là optional
    -- và có thể tham chiếu đến các bảng asset/product trong các service khác
);

-- Indexes cho bảng property_relation
CREATE INDEX IF NOT EXISTS idx_property_relation_asset_id ON property_relation(asset_id);
CREATE INDEX IF NOT EXISTS idx_property_relation_primary_product_id ON property_relation(primary_product_id);

-- Create triggers to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_property_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER property_updated_at_trigger
BEFORE UPDATE ON property
FOR EACH ROW
EXECUTE FUNCTION update_property_updated_at();

CREATE TRIGGER property_land_info_updated_at_trigger
BEFORE UPDATE ON property_land_info
FOR EACH ROW
EXECUTE FUNCTION update_property_updated_at();

CREATE TRIGGER property_building_info_updated_at_trigger
BEFORE UPDATE ON property_building_info
FOR EACH ROW
EXECUTE FUNCTION update_property_updated_at();

CREATE TRIGGER tag_updated_at_trigger
BEFORE UPDATE ON tag
FOR EACH ROW
EXECUTE FUNCTION update_property_updated_at();

CREATE TRIGGER property_relation_updated_at_trigger
BEFORE UPDATE ON property_relation
FOR EACH ROW
EXECUTE FUNCTION update_property_updated_at();

