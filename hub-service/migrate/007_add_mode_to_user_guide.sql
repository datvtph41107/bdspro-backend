-- Migration: Add mode column to user guide table
-- Created: 2025-01-XX
-- Description: Thêm cột mode để phân biệt chế độ đơn (10) và các bước (20)

-- Add mode column to tb_user_guide table
ALTER TABLE tb_user_guide 
ADD COLUMN IF NOT EXISTS mode INT NOT NULL DEFAULT 10;

-- Add comment
COMMENT ON COLUMN tb_user_guide.mode IS 'Chế độ hiển thị: 10 = chế độ đơn, 20 = các bước';

-- Update existing records to default mode (10: chế độ đơn)
UPDATE tb_user_guide SET mode = 10 WHERE mode IS NULL;

