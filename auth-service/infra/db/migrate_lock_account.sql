-- Migration: Add account lock fields to auth_method table
-- Description: Thêm các trường hỗ trợ khóa tạm thời và vĩnh viễn tài khoản

-- Thêm các cột mới vào bảng auth_method
ALTER TABLE auth_method 
ADD COLUMN status TINYINT UNSIGNED DEFAULT 1 COMMENT 'Trạng thái tài khoản: 1=active, 0=inactive, 2=temporarily_locked, 3=permanently_locked',
ADD COLUMN locked_at TIMESTAMP NULL COMMENT 'Thời gian bị khóa',
ADD COLUMN locked_until TIMESTAMP NULL COMMENT 'Thời gian khóa đến khi nào (cho khóa tạm thời)',
ADD COLUMN lock_reason VARCHAR(500) DEFAULT '' COMMENT 'Lý do khóa tài khoản',
ADD COLUMN locked_by BIGINT UNSIGNED NULL COMMENT 'ID admin khóa tài khoản',
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Thời gian cập nhật';

-- Tạo index cho các trường mới để tối ưu query
CREATE INDEX idx_auth_method_status ON auth_method(status);
CREATE INDEX idx_auth_method_locked_until ON auth_method(locked_until);
CREATE INDEX idx_auth_method_locked_by ON auth_method(locked_by);

-- Cập nhật các tài khoản hiện tại về trạng thái active
UPDATE auth_method SET status = 1 WHERE status IS NULL;
