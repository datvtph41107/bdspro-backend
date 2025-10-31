-- Migration: Create provinces and districts tables
-- Created: 2025-10-23

-- Create provinces table
CREATE TABLE IF NOT EXISTS provinces (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type INT NOT NULL,
    type_text VARCHAR(50) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

-- Create index for provinces
CREATE INDEX IF NOT EXISTS idx_provinces_deleted_at ON provinces(deleted_at);
CREATE INDEX IF NOT EXISTS idx_provinces_slug ON provinces(slug);
CREATE INDEX IF NOT EXISTS idx_provinces_type ON provinces(type);

-- Create districts table
CREATE TABLE IF NOT EXISTS districts (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    province_id VARCHAR(10) NOT NULL,
    type INT NOT NULL,
    type_text VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE
);

-- Create index for districts
CREATE INDEX IF NOT EXISTS idx_districts_deleted_at ON districts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_districts_province_id ON districts(province_id);
CREATE INDEX IF NOT EXISTS idx_districts_type ON districts(type);
CREATE INDEX IF NOT EXISTS idx_districts_name ON districts(name);

-- Comments
COMMENT ON TABLE provinces IS 'Stores provinces/cities data of Vietnam';
COMMENT ON TABLE districts IS 'Stores districts data of Vietnam';
COMMENT ON COLUMN provinces.type IS '1: Trung ương, 2: Tỉnh';
COMMENT ON COLUMN districts.type IS '3: Thành phố thuộc tỉnh, 4: Quận, 5: Thị xã, 6: Huyện';

