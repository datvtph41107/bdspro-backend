-- Migration: Thêm các trường Diffie-Hellman key exchange vào bảng auth_method
-- Date: 2025-10-27

-- Thêm cột private_key, public_key, auth_key vào bảng auth_method
ALTER TABLE auth_method
ADD COLUMN IF NOT EXISTS private_key VARCHAR(512) DEFAULT NULL COMMENT 'Diffie-Hellman private key (base64 encoded)',
ADD COLUMN IF NOT EXISTS public_key VARCHAR(512) DEFAULT NULL COMMENT 'Diffie-Hellman public key (base64 encoded)',
ADD COLUMN IF NOT EXISTS auth_key VARCHAR(512) DEFAULT NULL COMMENT 'Auth key được mã hóa khi xác thực';

-- Tạo index cho public_key để tăng tốc độ truy vấn (nếu cần)
CREATE INDEX IF NOT EXISTS idx_auth_method_public_key ON auth_method(public_key);

-- Rollback commands (nếu cần revert migration):
-- ALTER TABLE auth_method DROP COLUMN IF EXISTS private_key;
-- ALTER TABLE auth_method DROP COLUMN IF EXISTS public_key;
-- ALTER TABLE auth_method DROP COLUMN IF EXISTS auth_key;
-- DROP INDEX IF EXISTS idx_auth_method_public_key ON auth_method;

