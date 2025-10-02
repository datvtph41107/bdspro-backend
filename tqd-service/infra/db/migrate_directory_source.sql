-- Migration for directory_sources table
CREATE TABLE IF NOT EXISTS directory_sources (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    icon VARCHAR(100),
    color VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    expected_amount BIGINT DEFAULT 0,
    actual_amount BIGINT DEFAULT 0,
    frequency VARCHAR(50),
    start_date DATE,
    is_recurring BOOLEAN DEFAULT false,
    payment_method VARCHAR(50),
    notes TEXT,
    tags TEXT, -- JSON array string
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_directory_sources_code ON directory_sources(code);
CREATE INDEX IF NOT EXISTS idx_directory_sources_category ON directory_sources(category);
CREATE INDEX IF NOT EXISTS idx_directory_sources_type ON directory_sources(type);
CREATE INDEX IF NOT EXISTS idx_directory_sources_is_active ON directory_sources(is_active);
CREATE INDEX IF NOT EXISTS idx_directory_sources_is_recurring ON directory_sources(is_recurring);
CREATE INDEX IF NOT EXISTS idx_directory_sources_sort_order ON directory_sources(sort_order);
CREATE INDEX IF NOT EXISTS idx_directory_sources_deleted_at ON directory_sources(deleted_at);

-- Insert sample data
INSERT INTO directory_sources (name, code, description, category, type, icon, color, is_active, sort_order, expected_amount, actual_amount, frequency, start_date, is_recurring, payment_method, notes, tags) VALUES
('Dịch vụ bảo trì', 'MAINT-SVC', 'Dịch vụ bảo trì tòa nhà', 'SERVICE', 'recurring', 'ToolOutlined', '#52c41a', true, 2, 5000000, 4800000, 'monthly', '2024-01-15', true, 'bank_transfer', 'Dịch vụ theo hợp đồng', '["Dịch vụ", "Dài hạn"]'),
('Dịch vụ vệ sinh', 'CLEAN-SVC', 'Dịch vụ vệ sinh hàng ngày', 'SERVICE', 'recurring', 'CleanOutlined', '#1890ff', true, 1, 2000000, 1950000, 'daily', '2024-01-01', true, 'cash', 'Dịch vụ vệ sinh tòa nhà', '["Vệ sinh", "Hàng ngày"]'),
('Dịch vụ bảo vệ', 'SECURITY-SVC', 'Dịch vụ bảo vệ 24/7', 'SERVICE', 'recurring', 'SecurityScanOutlined', '#f5222d', true, 3, 8000000, 7800000, 'monthly', '2024-01-01', true, 'bank_transfer', 'Dịch vụ bảo vệ an ninh', '["Bảo vệ", "24/7"]'),
('Sửa chữa điện', 'ELEC-REPAIR', 'Sửa chữa hệ thống điện', 'REPAIR', 'one-time', 'ThunderboltOutlined', '#faad14', true, 4, 1500000, 1500000, NULL, '2024-01-20', false, 'bank_transfer', 'Sửa chữa khẩn cấp', '["Điện", "Sửa chữa"]'),
('Nâng cấp hệ thống', 'SYSTEM-UPGRADE', 'Nâng cấp hệ thống IT', 'UPGRADE', 'one-time', 'ApiOutlined', '#722ed1', true, 5, 12000000, 11500000, NULL, '2024-02-01', false, 'bank_transfer', 'Nâng cấp hệ thống quản lý', '["IT", "Nâng cấp"]'),
('Dịch vụ internet', 'INTERNET-SVC', 'Dịch vụ internet doanh nghiệp', 'SERVICE', 'recurring', 'WifiOutlined', '#13c2c2', true, 6, 3000000, 3000000, 'monthly', '2024-01-01', true, 'bank_transfer', 'Gói internet doanh nghiệp', '["Internet", "Doanh nghiệp"]'),
('Bảo trì thang máy', 'ELEVATOR-MAINT', 'Bảo trì thang máy định kỳ', 'SERVICE', 'recurring', 'VerticalAlignTopOutlined', '#fa8c16', true, 7, 4000000, 3800000, 'quarterly', '2024-01-01', true, 'bank_transfer', 'Bảo trì thang máy 3 tháng/lần', '["Thang máy", "Bảo trì"]'),
('Dịch vụ cây xanh', 'GARDEN-SVC', 'Chăm sóc cây xanh', 'SERVICE', 'recurring', 'EnvironmentOutlined', '#52c41a', false, 8, 1000000, 950000, 'weekly', '2024-01-01', true, 'cash', 'Chăm sóc cây xanh khu vực', '["Cây xanh", "Chăm sóc"]')
ON CONFLICT (code) DO NOTHING;
