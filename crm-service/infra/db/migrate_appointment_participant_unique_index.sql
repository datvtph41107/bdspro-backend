-- Migration: Tạo partial unique index cho appointment_participants
-- Description: Tạo unique constraint cho (appointment_id, contact_id) chỉ áp dụng cho các record chưa bị xóa
-- Date: 2025-01-XX

-- Xóa unique index cũ nếu có (từ GORM tag)
DROP INDEX IF EXISTS idx_appointment_contact_unique;

-- Tạo partial unique index chỉ áp dụng cho các record chưa bị xóa
CREATE UNIQUE INDEX IF NOT EXISTS idx_appointment_contact_unique 
ON appointment_participants (appointment_id, contact_id) 
WHERE deleted_at IS NULL;

-- Comment
COMMENT ON INDEX idx_appointment_contact_unique IS 'Unique constraint cho appointment_id và contact_id, chỉ áp dụng cho các record chưa bị xóa (deleted_at IS NULL)';

