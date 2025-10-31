-- Migration: Create wards table
-- Created: 2025-10-23

-- Create wards table
CREATE TABLE IF NOT EXISTS wards (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    district_id VARCHAR(10) NOT NULL,
    type INT NOT NULL,
    type_text VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE CASCADE
);

-- Create index for wards
CREATE INDEX IF NOT EXISTS idx_wards_deleted_at ON wards(deleted_at);
CREATE INDEX IF NOT EXISTS idx_wards_district_id ON wards(district_id);
CREATE INDEX IF NOT EXISTS idx_wards_type ON wards(type);
CREATE INDEX IF NOT EXISTS idx_wards_name ON wards(name);

-- Comments
COMMENT ON TABLE wards IS 'Stores wards/communes data of Vietnam';
COMMENT ON COLUMN wards.type IS '7: Phường, 8: Xã, 9: Thị trấn';

