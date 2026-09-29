-- Migration: Tạo unique index cho appointment_participants
-- Description: Tạo unique constraint cho (appointment_id, user_id) để hỗ trợ ON CONFLICT
-- Date: 2025-01-XX

-- Xóa unique index cũ nếu có (từ GORM tag hoặc migration cũ)
DROP INDEX IF EXISTS idx_appointment_contact_unique;
DROP INDEX IF EXISTS idx_appointment_user_unique;

-- Tạo unique index cho (appointment_id, user_id)
-- Note: Không dùng partial index vì ON CONFLICT chỉ hoạt động với unique constraint/index không có điều kiện
-- Logic restore soft-deleted records được xử lý trong code bằng cách set deleted_at = NULL và is_deleted = false
CREATE UNIQUE INDEX IF NOT EXISTS idx_appointment_user_unique
ON appointment_participants (appointment_id, user_id);

-- Comment
COMMENT ON INDEX idx_appointment_user_unique IS 'Unique constraint cho appointment_id và user_id';
