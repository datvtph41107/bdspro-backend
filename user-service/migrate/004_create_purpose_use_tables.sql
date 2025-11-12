-- Migration: Tạo bảng purpose_use và purpose_use_profile để lưu mục đích sử dụng của user
-- Date: 2025-11-11

-- Bảng purpose_use lưu danh sách mục đích sử dụng
CREATE TABLE IF NOT EXISTS purpose_use (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

COMMENT ON TABLE purpose_use IS 'Danh mục mục đích sử dụng gắn với người dùng';
COMMENT ON COLUMN purpose_use.name IS 'Tên mục đích sử dụng';
COMMENT ON COLUMN purpose_use.description IS 'Mô tả ngắn cho mục đích sử dụng';

CREATE UNIQUE INDEX IF NOT EXISTS uq_purpose_use_name ON purpose_use (LOWER(name));

-- Bảng nối giữa user và purpose_use
CREATE TABLE IF NOT EXISTS purpose_use_profile (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    purpose_use_id BIGINT NOT NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_purpose_use_profile_profile FOREIGN KEY (profile_id) REFERENCES user_profile(profile_id),
    CONSTRAINT fk_purpose_use_profile_purpose FOREIGN KEY (purpose_use_id) REFERENCES purpose_use(id)
);

COMMENT ON TABLE purpose_use_profile IS 'Bảng nối giữa user_profile và purpose_use';

CREATE UNIQUE INDEX IF NOT EXISTS uq_purpose_use_profile_active
    ON purpose_use_profile (profile_id, purpose_use_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_purpose_use_profile_profile_id
    ON purpose_use_profile (profile_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_purpose_use_profile_purpose_use_id
    ON purpose_use_profile (purpose_use_id)
    WHERE deleted_at IS NULL;

