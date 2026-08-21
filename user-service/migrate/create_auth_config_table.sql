-- Migration: Tạo bảng auth_config để lưu cấu hình động
-- Date: 2025-10-24
-- Description: Bảng này lưu trữ các cấu hình động của hệ thống auth như ZNS tokens, JWT settings, OTP settings

-- Tạo bảng auth_config
CREATE TABLE IF NOT EXISTS auth_config (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(10) NOT NULL UNIQUE,
    config_value TEXT,
    description VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    config_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT
);

-- Index để tìm kiếm nhanh theo config_key
CREATE INDEX IF NOT EXISTS idx_auth_config_key ON auth_config(config_key) WHERE deleted_at IS NULL;

-- Index để filter theo config_type
CREATE INDEX IF NOT EXISTS idx_auth_config_type ON auth_config(config_type) WHERE deleted_at IS NULL AND is_active = true;

-- Comments
COMMENT ON TABLE auth_config IS 'Bảng lưu trữ cấu hình động của hệ thống auth';
COMMENT ON COLUMN auth_config.config_key IS 'Khóa cấu hình (unique): zns.access_token, zns.refresh_token, jwt.secret_key, etc.';
COMMENT ON COLUMN auth_config.config_value IS 'Giá trị cấu hình dạng text';
COMMENT ON COLUMN auth_config.description IS 'Mô tả cấu hình';
COMMENT ON COLUMN auth_config.is_active IS 'Trạng thái kích hoạt của cấu hình';
COMMENT ON COLUMN auth_config.config_type IS 'Loại cấu hình: ZNS, JWT, OTP, SMS';

-- Insert dữ liệu mẫu từ config hiện tại
INSERT INTO auth_config (config_key, config_value, is_active) VALUES
-- ZNS Config
('ZNS_TOKEN', 'CHANGE_ME', true),
('ZNS_REFRESH', 'CHANGE_ME', true),
('ZNS_TEMPLATE_ID', '469663', true),
('ZNS_OTP_TEMPLATE_ID', '497058', true)

ON CONFLICT (config_key) DO NOTHING;

