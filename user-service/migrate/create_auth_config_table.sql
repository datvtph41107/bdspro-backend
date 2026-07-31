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
('ZNS_TOKEN', 'RRoFMY3vhLDil81HA_QkKcQuZLyXwlfKFCw90d7zvWrsnAq0FktTRnF4W3yUxCXf4ToX1ZgDqmahjxf5BBZ_Cp29hMeinjG23ENzNZ3OjYK_mlLQ3_IPDpJZuNeArueNHlEW64w_pt1dzPaXRhdFSKlnfI9i-zSlUUcWL7xmprasxu007_h_0oxyg7OhvkGb1g-EJXITuKHUq8yjJV_nNNllsWPDzxG4UEIOQc7SkYbjbC958z-eDHINo7TPaQS_TAdxTbIzaYzRYznpPPw01dAKgdv2qPqyTi7rK46FX28CXErh7ysH8WJZsaWNpuue2TAaFp_qyc8E_uGADR3JGWkEdX83YjDb4goOSXwJtI0arvTv5-NwDHpNk4rIm-H0O_6E04MW-7DJkQzyJQRN0LgvZ5Dqj-K4Gtzy6ZecxgTO', true),
('ZNS_REFRESH', '3AIiBOb4LqKDfACcWajKH0UWhYY8Lo4mIDUi1iXKE2r6t9SEoNKSR0tMdLRPJZP5MlQGKTKPQp8livqDpLau7Yd7c2FQNmqxEEg90knw4smBbhWFw1ev6mwfio_UFoa74CkFA_Tq0HSPrV8M-7Lf82B0sHVqQIWJSFpPVerrO4jOoDfytrPBSJBUnndrJdmo9ARL1VaAN1eRjAuivGaLEGMFkGNJ92W13x6gAE4P3Hm_jQaMuWeYDH-7fIRGPZn8Bu-cHjSyCtWEiQ1eu0yiGsE0btpaDWTBQjUMMB9c0qurowH2qMChR5Rqtd2MUon1UvAPAFy_A5uzgDezpY8i1H6kh1JbHGO14VUE8l5nUHGXtUyMu7nA56ULlX-8EpGtJBEaEu8E1IfJWO54XJiTJ7kCvcQQRKuBURVPmrkFKrrE', true),
('ZNS_TEMPLATE_ID', '469663', true),
('ZNS_OTP_TEMPLATE_ID', '497058', true)

ON CONFLICT (config_key) DO NOTHING;

