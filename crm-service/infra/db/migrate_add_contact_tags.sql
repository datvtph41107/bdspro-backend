-- Migration: Thêm cột contact_tags (mảng integer) vào bảng tb_contact
-- Date: 2025-01-XX

-- Thêm cột contact_tags với kiểu integer[]
ALTER TABLE tb_contact 
ADD COLUMN IF NOT EXISTS contact_tags INTEGER[] DEFAULT NULL;

-- Thêm comment cho cột
COMMENT ON COLUMN tb_contact.contact_tags IS 'Mảng các tag contact (10: Khách hàng, 20: Môi giới, 30: Đối tác, 40: Chủ nhà)';
