-- Migration script for user_pin table
-- Tạo bảng user_pin để lưu trữ mã PIN của người dùng

CREATE TABLE IF NOT EXISTS user_pin (
    auth_id BIGINT PRIMARY KEY REFERENCES auth_method(id) ON DELETE CASCADE,
    pin VARCHAR(255) NOT NULL,                     -- Mã PIN đã được hash (bcrypt)
    pin_check_time INT DEFAULT 0 NOT NULL,         -- Số lần nhập sai PIN
    pin_date TIMESTAMPTZ,                          -- Thời gian tạo/cập nhật PIN
    is_active BOOLEAN DEFAULT TRUE NOT NULL,       -- PIN có đang hoạt động không
    locked_until TIMESTAMPTZ,                      -- Thời gian khóa tài khoản (nếu nhập sai nhiều lần)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Index cho việc tìm kiếm nhanh
CREATE INDEX IF NOT EXISTS idx_user_pin_auth_id ON user_pin(auth_id);
CREATE INDEX IF NOT EXISTS idx_user_pin_is_active ON user_pin(is_active);
CREATE INDEX IF NOT EXISTS idx_user_pin_locked_until ON user_pin(locked_until);

-- Trigger để tự động cập nhật updated_at
CREATE OR REPLACE FUNCTION update_user_pin_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_user_pin_updated_at
BEFORE UPDATE ON user_pin
FOR EACH ROW
EXECUTE FUNCTION update_user_pin_updated_at();

-- Comments cho các trường
COMMENT ON TABLE user_pin IS 'Bảng lưu trữ mã PIN của người dùng';
COMMENT ON COLUMN user_pin.auth_id IS 'ID của auth_method (Foreign Key)';
COMMENT ON COLUMN user_pin.pin IS 'Mã PIN đã được hash bằng bcrypt';
COMMENT ON COLUMN user_pin.pin_check_time IS 'Số lần nhập sai PIN liên tiếp';
COMMENT ON COLUMN user_pin.pin_date IS 'Thời gian tạo hoặc cập nhật PIN';
COMMENT ON COLUMN user_pin.is_active IS 'PIN có đang hoạt động hay bị vô hiệu hóa';
COMMENT ON COLUMN user_pin.locked_until IS 'Thời gian khóa PIN do nhập sai quá nhiều lần';

