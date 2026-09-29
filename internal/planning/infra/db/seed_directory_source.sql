-- Seed data for directory_sources table
INSERT INTO directory_sources (
    name, code, description, type, icon, color, is_active, sort_order,
    expected_amount, actual_amount, frequency, start_date, is_recurring,
    payment_method, notes, raw_tags, created_at, updated_at
) VALUES 
(
    'Thuê căn hộ A1', 'RENT-A1', 'Thu nhập từ việc cho thuê căn hộ A1', 
    'recurring', 'HomeOutlined', '#1890ff', true, 1,
    15000000, 15000000, 'monthly', '2024-01-01', true,
    'bank_transfer', 'Khách thuê ổn định', '["Thu nhập chính", "Ổn định"]',
    NOW(), NOW()
),
(
    'Dịch vụ bảo trì', 'MAINT-SVC', 'Dịch vụ bảo trì tòa nhà', 
    'recurring', 'ToolOutlined', '#52c41a', true, 2,
    5000000, 4800000, 'monthly', '2024-01-15', true,
    'bank_transfer', 'Dịch vụ theo hợp đồng', '["Dịch vụ", "Dài hạn"]',
    NOW(), NOW()
),
(
    'Bán sản phẩm', 'SALES-PROD', 'Bán các sản phẩm bất động sản', 
    'one-time', 'ShoppingOutlined', '#faad14', true, 3,
    2000000000, 1500000000, 'custom', '2024-01-01', false,
    'bank_transfer', 'Bán căn hộ dự án', '["Bán hàng", "Biến động"]',
    NOW(), NOW()
),
(
    'Hoa hồng môi giới', 'COMM-BROKER', 'Hoa hồng từ dịch vụ môi giới', 
    'commission', 'PercentageOutlined', '#f5222d', true, 4,
    10000000, 7500000, 'monthly', '2024-01-01', true,
    'bank_transfer', 'Hoa hồng theo giao dịch', '["Hoa hồng", "Theo mùa"]',
    NOW(), NOW()
),
(
    'Đầu tư chứng khoán', 'INVEST-STOCK', 'Lợi nhuận từ đầu tư chứng khoán', 
    'recurring', 'RiseOutlined', '#722ed1', false, 5,
    2000000, 1500000, 'monthly', '2024-01-01', true,
    'bank_transfer', 'Đầu tư rủi ro cao', '["Đầu tư", "Rủi ro cao"]',
    NOW(), NOW()
);
