-- Migration: Tạo 3 bảng provinces, districts, wards riêng biệt
-- Thay thế bảng region cũ

-- 1. Tạo bảng provinces
CREATE TABLE IF NOT EXISTS provinces (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code INTEGER,
    code_name VARCHAR(255),
    unit VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

-- 2. Tạo bảng districts
CREATE TABLE IF NOT EXISTS districts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    province_id BIGINT NOT NULL,
    code INTEGER,
    code_name VARCHAR(255),
    unit VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE
);

-- 3. Tạo bảng wards
CREATE TABLE IF NOT EXISTS wards (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    district_id BIGINT NOT NULL,
    code INTEGER,
    code_name VARCHAR(255),
    unit VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE CASCADE
);

-- 4. Migrate dữ liệu từ bảng region cũ (nếu tồn tại)
-- Migrate provinces (level = 1)
INSERT INTO provinces (id, name, code, code_name, unit, created_at, updated_at, deleted_at, created_by, updated_by)
SELECT id, name, code, code_name, unit, created_at, updated_at, deleted_at, created_by, updated_by
FROM region
WHERE level = 1 AND deleted_at IS NULL
ON CONFLICT (id) DO NOTHING;

-- Migrate districts (level = 2)
INSERT INTO districts (id, name, province_id, code, code_name, unit, created_at, updated_at, deleted_at, created_by, updated_by)
SELECT id, name, parent_id, code, code_name, unit, created_at, updated_at, deleted_at, created_by, updated_by
FROM region
WHERE level = 2 AND deleted_at IS NULL AND parent_id IS NOT NULL
ON CONFLICT (id) DO NOTHING;

-- Migrate wards (level = 3)
INSERT INTO wards (id, name, district_id, code, code_name, unit, created_at, updated_at, deleted_at, created_by, updated_by)
SELECT id, name, parent_id, code, code_name, unit, created_at, updated_at, deleted_at, created_by, updated_by
FROM region
WHERE level = 3 AND deleted_at IS NULL AND parent_id IS NOT NULL
ON CONFLICT (id) DO NOTHING;

-- 5. Tạo indexes cho performance
CREATE INDEX IF NOT EXISTS idx_provinces_deleted_at ON provinces(deleted_at);
CREATE INDEX IF NOT EXISTS idx_provinces_name ON provinces(name);
CREATE INDEX IF NOT EXISTS idx_provinces_code ON provinces(code);

CREATE INDEX IF NOT EXISTS idx_districts_deleted_at ON districts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_districts_province_id ON districts(province_id);
CREATE INDEX IF NOT EXISTS idx_districts_name ON districts(name);
CREATE INDEX IF NOT EXISTS idx_districts_code ON districts(code);

CREATE INDEX IF NOT EXISTS idx_wards_deleted_at ON wards(deleted_at);
CREATE INDEX IF NOT EXISTS idx_wards_district_id ON wards(district_id);
CREATE INDEX IF NOT EXISTS idx_wards_name ON wards(name);
CREATE INDEX IF NOT EXISTS idx_wards_code ON wards(code);

-- 6. Comments
COMMENT ON TABLE provinces IS 'Bảng tỉnh/thành phố';
COMMENT ON TABLE districts IS 'Bảng quận/huyện';
COMMENT ON TABLE wards IS 'Bảng phường/xã';

-- Note: Có thể drop bảng region cũ sau khi chạy migration thành công
-- DROP TABLE IF EXISTS region CASCADE;

