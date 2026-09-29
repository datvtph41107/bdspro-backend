-- Migration: Add image_id column to products table
-- Description: Thêm trường image_id để lưu ảnh mặc định của sản phẩm

ALTER TABLE products ADD COLUMN IF NOT EXISTS image_id BIGINT;

-- Add comment to column
COMMENT ON COLUMN products.image_id IS 'ID của ảnh mặc định của sản phẩm';

