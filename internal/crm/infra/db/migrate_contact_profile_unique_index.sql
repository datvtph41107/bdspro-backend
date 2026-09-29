-- Migration: Tạo partial unique index cho tb_contact
-- Description: Tạo unique constraint cho (profile_id, owner_id, owner_of) chỉ áp dụng cho các record chưa bị xóa và có profile_id
-- Date: 2025-12-25

-- Xóa unique index cũ nếu có
DROP INDEX IF EXISTS idx_contact_profile_owner_unique;

-- Tạo partial unique index chỉ áp dụng cho các record chưa bị xóa và có profile_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_profile_owner_unique 
ON tb_contact (profile_id, owner_id, owner_of) 
WHERE profile_id IS NOT NULL AND deleted_at IS NULL;

-- Comment
COMMENT ON INDEX idx_contact_profile_owner_unique IS 'Unique constraint cho profile_id, owner_id và owner_of, chỉ áp dụng cho các record chưa bị xóa (deleted_at IS NULL) và có profile_id';
