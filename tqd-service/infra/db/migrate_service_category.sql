-- Migration for directory_categories table
CREATE TABLE IF NOT EXISTS directory_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    code VARCHAR(50) NOT NULL UNIQUE,
    icon VARCHAR(100),
    color VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    path VARCHAR(255),
    service_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_directory_categories_code ON directory_categories(code);
CREATE INDEX IF NOT EXISTS idx_directory_categories_is_active ON directory_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_directory_categories_level ON directory_categories(level);
CREATE INDEX IF NOT EXISTS idx_directory_categories_sort_order ON directory_categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_directory_categories_deleted_at ON directory_categories(deleted_at);

-- Insert sample data
INSERT INTO directory_categories (name, description, code, icon, color, is_active, sort_order, level, path, service_count) VALUES
('Dịch vụ nhà ở', 'Các dịch vụ liên quan đến nhà ở và bất động sản', 'HOUSING', 'HomeOutlined', '#1890ff', true, 1, 1, '/housing', 15),
('Dịch vụ mua sắm', 'Các dịch vụ mua sắm và thương mại', 'SHOPPING', 'ShopOutlined', '#52c41a', true, 2, 1, '/shopping', 8),
('Dịch vụ giao thông', 'Các dịch vụ vận chuyển và giao thông', 'TRANSPORT', 'CarOutlined', '#faad14', true, 3, 1, '/transport', 12),
('Dịch vụ tài chính', 'Các dịch vụ ngân hàng và tài chính', 'FINANCE', 'BankOutlined', '#f5222d', true, 4, 1, '/finance', 6),
('Dịch vụ y tế', 'Các dịch vụ chăm sóc sức khỏe', 'HEALTH', 'HeartOutlined', '#722ed1', false, 5, 1, '/health', 4),
('Dịch vụ giáo dục', 'Các dịch vụ học tập và đào tạo', 'EDUCATION', 'BookOutlined', '#13c2c2', true, 6, 1, '/education', 9),
('Dịch vụ xây dựng', 'Các dịch vụ xây dựng và sửa chữa', 'CONSTRUCTION', 'BuildOutlined', '#fa8c16', true, 7, 1, '/construction', 11),
('Dịch vụ công nghệ', 'Các dịch vụ công nghệ thông tin', 'TECHNOLOGY', 'GlobalOutlined', '#eb2f96', true, 8, 1, '/technology', 7)
ON CONFLICT (code) DO NOTHING;
