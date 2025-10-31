-- Migration: Create profile_deleted table
-- Description: Bảng lưu thông tin profile đã bị xóa

CREATE TABLE IF NOT EXISTS profile_deleted (
    id BIGSERIAL PRIMARY KEY,
    profile_id BIGINT NOT NULL,
    phone VARCHAR(20),
    phone2 VARCHAR(20),
    email VARCHAR(50),
    tax_code VARCHAR(13),
    full_name VARCHAR(255),
    address VARCHAR(255),
    gender SMALLINT,
    avatar VARCHAR(255),
    tick_verified BOOLEAN DEFAULT false,
    background_image VARCHAR(255),
    referral_code VARCHAR(255),
    referral_by VARCHAR(255),
    role_type INTEGER DEFAULT 10,
    role_key INTEGER DEFAULT 10,
    role_real_estate INTEGER DEFAULT 10,
    position VARCHAR(255),
    department_id BIGINT,
    birth TIMESTAMP,
    start_date TIMESTAMP,
    front_identify VARCHAR(255),
    back_identify VARCHAR(255),
    plan_id BIGINT,
    plan_at TIMESTAMP,
    slogan TEXT,
    website VARCHAR(255),
    facebook VARCHAR(255),
    instagram VARCHAR(255),
    twitter VARCHAR(255),
    linkedin VARCHAR(255),
    youtube VARCHAR(255),
    introduction TEXT,
    visibility SMALLINT,
    status_online INTEGER,
    profile_visibility SMALLINT,
    tab_default INTEGER,
    status INTEGER DEFAULT 10,
    warning INTEGER DEFAULT 0,
    role_id BIGINT,
    
    -- Metadata về việc xóa
    deleted_by BIGINT,
    deleted_reason TEXT,
    deleted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    original_created_at TIMESTAMP,
    original_updated_at TIMESTAMP,
    
    -- Indexes
    CONSTRAINT idx_profile_deleted_profile_id UNIQUE (profile_id)
);

-- Create index for faster lookup
CREATE INDEX IF NOT EXISTS idx_profile_deleted_deleted_at ON profile_deleted(deleted_at DESC);
CREATE INDEX IF NOT EXISTS idx_profile_deleted_deleted_by ON profile_deleted(deleted_by);
CREATE INDEX IF NOT EXISTS idx_profile_deleted_email ON profile_deleted(email);
CREATE INDEX IF NOT EXISTS idx_profile_deleted_phone ON profile_deleted(phone);

-- Add comments
COMMENT ON TABLE profile_deleted IS 'Bảng lưu thông tin profile đã bị xóa';
COMMENT ON COLUMN profile_deleted.profile_id IS 'ID của profile đã bị xóa';
COMMENT ON COLUMN profile_deleted.deleted_by IS 'ID của người thực hiện xóa';
COMMENT ON COLUMN profile_deleted.deleted_reason IS 'Lý do xóa profile';
COMMENT ON COLUMN profile_deleted.deleted_at IS 'Thời điểm xóa profile';
COMMENT ON COLUMN profile_deleted.original_created_at IS 'Thời điểm tạo profile ban đầu';
COMMENT ON COLUMN profile_deleted.original_updated_at IS 'Thời điểm cập nhật cuối cùng trước khi xóa';

