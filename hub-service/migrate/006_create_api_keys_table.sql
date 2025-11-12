-- Migration for Api Keys table
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    app_name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    api_key VARCHAR(255) NOT NULL,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT uq_api_keys_name UNIQUE (name),
    CONSTRAINT uq_api_keys_value UNIQUE (api_key)
);

CREATE INDEX IF NOT EXISTS idx_api_keys_deleted_at ON api_keys (deleted_at);
CREATE INDEX IF NOT EXISTS idx_api_keys_expired_at ON api_keys (expired_at);
CREATE INDEX IF NOT EXISTS idx_api_keys_app_name ON api_keys (app_name);

CREATE OR REPLACE FUNCTION update_api_keys_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_api_keys_updated_at
BEFORE UPDATE ON api_keys
FOR EACH ROW
EXECUTE FUNCTION update_api_keys_updated_at();

COMMENT ON TABLE api_keys IS 'Bảng lưu trữ API key dùng cho giao tiếp dịch vụ';
COMMENT ON COLUMN api_keys.name IS 'Tên hiển thị của API key';
COMMENT ON COLUMN api_keys.app_name IS 'Tên ứng dụng hoặc dịch vụ sử dụng key';
COMMENT ON COLUMN api_keys.description IS 'Mô tả ứng dụng hoặc mục đích sử dụng';
COMMENT ON COLUMN api_keys.api_key IS 'Giá trị API key';
