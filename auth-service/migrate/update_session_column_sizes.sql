-- Migration: Cập nhật kích thước các cột trong bảng user_session để hỗ trợ IPv6 và version dài hơn
-- Date: 2026-06-16

ALTER TABLE user_session 
    ALTER COLUMN ip_request TYPE VARCHAR(50),
    ALTER COLUMN platform TYPE VARCHAR(50),
    ALTER COLUMN version TYPE VARCHAR(50),
    ALTER COLUMN os TYPE VARCHAR(50);

COMMENT ON COLUMN user_session.ip_request IS 'Địa chỉ IP của request (hỗ trợ IPv6)';
