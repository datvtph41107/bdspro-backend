-- Migration: Tạo bảng area_active và area_active_profile
-- Date: 2025-11-11

CREATE TABLE IF NOT EXISTS area_active (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(512),
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

COMMENT ON TABLE area_active IS 'Danh sách nhãn khu vực hoạt động của người dùng';
COMMENT ON COLUMN area_active.name IS 'Tên hiển thị của nhãn khu vực';
COMMENT ON COLUMN area_active.description IS 'Mô tả thêm cho nhãn khu vực';

CREATE TABLE IF NOT EXISTS area_active_profile (
    profile_id BIGINT NOT NULL,
    area_active_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (profile_id, area_active_id),
    CONSTRAINT fk_area_active_profile_profile FOREIGN KEY (profile_id) REFERENCES user_profile(profile_id) ON DELETE CASCADE,
    CONSTRAINT fk_area_active_profile_area FOREIGN KEY (area_active_id) REFERENCES area_active(id) ON DELETE CASCADE
);

COMMENT ON TABLE area_active_profile IS 'Bảng nối profile và nhãn khu vực hoạt động';

CREATE INDEX IF NOT EXISTS idx_area_active_profile_profile_id ON area_active_profile(profile_id);
CREATE INDEX IF NOT EXISTS idx_area_active_profile_area_active_id ON area_active_profile(area_active_id);

