-- Migration for Versions table
CREATE TABLE IF NOT EXISTS versions (
    id BIGSERIAL PRIMARY KEY,
    app_name VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    version_name VARCHAR(50) NOT NULL,
    build_number VARCHAR(50) NOT NULL,
    force_update BOOLEAN DEFAULT FALSE,
    download_url VARCHAR(255),
    release_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT uq_versions_app_platform_version UNIQUE(app_name, platform, version_name)
);

CREATE INDEX IF NOT EXISTS idx_versions_app_platform ON versions (app_name, platform);
CREATE INDEX IF NOT EXISTS idx_versions_platform ON versions (platform);
CREATE INDEX IF NOT EXISTS idx_versions_deleted_at ON versions (deleted_at);

CREATE OR REPLACE FUNCTION update_versions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_versions_updated_at
BEFORE UPDATE ON versions
FOR EACH ROW
EXECUTE FUNCTION update_versions_updated_at();

COMMENT ON TABLE versions IS 'Bảng lưu thông tin phiên bản ứng dụng theo từng nền tảng';
COMMENT ON COLUMN versions.app_name IS 'Tên ứng dụng (VD: hub-app, agent-app, ...)';
COMMENT ON COLUMN versions.platform IS 'Nền tảng ứng dụng (ios, android, web...)';
COMMENT ON COLUMN versions.version_name IS 'Tên hiển thị phiên bản (VD: 1.0.0)';
COMMENT ON COLUMN versions.build_number IS 'Mã build, dùng cho so sánh nâng cấp';
COMMENT ON COLUMN versions.force_update IS 'Có bắt buộc cập nhật hay không';
COMMENT ON COLUMN versions.download_url IS 'Đường dẫn tải phiên bản';


