-- Create poi_categories table
CREATE TABLE IF NOT EXISTS poi_categories (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    parent_id BIGINT REFERENCES poi_categories(id) ON DELETE CASCADE,
    icon VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed data for service_categories table
-- Insert parent categories first
INSERT INTO poi_categories (
    id, code, name, icon, is_active, sort_order, description, created_at, updated_at
) VALUES 
(1, 'RESIDENTIAL', 'Khu dân cư', 'HomeOutlined', true, 1, 'Các khu dân cư, chung cư, nhà ở', '2024-01-15', '2024-01-20'),
(2, 'COMMERCIAL', 'Thương mại', 'ShopOutlined', true, 2, 'Các khu thương mại, trung tâm mua sắm', '2024-01-10', '2024-01-18'),
(3, 'OFFICE', 'Văn phòng', 'BuildOutlined', true, 3, 'Tòa nhà văn phòng, công ty', '2024-01-08', '2024-01-16'),
(4, 'RETAIL', 'Bán lẻ', 'ShopOutlined', true, 4, 'Cửa hàng bán lẻ, shop', '2024-01-05', '2024-01-15'),
(5, 'PARKING', 'Bãi đỗ xe', 'CarOutlined', true, 5, 'Bãi đỗ xe, garage', '2024-01-12', '2024-01-16'),
(6, 'BANK', 'Ngân hàng', 'BankOutlined', true, 6, 'Ngân hàng, tài chính', '2024-01-20', '2024-01-25'),
(7, 'HOSPITAL', 'Bệnh viện', 'HeartOutlined', true, 7, 'Bệnh viện, phòng khám', '2024-01-08', '2024-01-12'),
(8, 'SCHOOL', 'Trường học', 'BookOutlined', true, 8, 'Trường học, trung tâm giáo dục', '2024-01-12', '2024-01-16');

-- Insert child categories
INSERT INTO service_categories (
    id, code, name, parent_id, icon, is_active, sort_order, created_at, updated_at
) VALUES 
(11, 'APARTMENT', 'Chung cư', 1, 'HomeOutlined', true, 1, '2024-01-15', '2024-01-20'),
(12, 'HOUSE', 'Nhà riêng', 1, 'HomeOutlined', true, 2, '2024-01-15', '2024-01-20'),
(21, 'MALL', 'Trung tâm thương mại', 2, 'ShopOutlined', true, 1, '2024-01-10', '2024-01-18'),
(22, 'MARKET', 'Chợ', 2, 'ShopOutlined', true, 2, '2024-01-10', '2024-01-18');

