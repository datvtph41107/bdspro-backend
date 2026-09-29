-- Seed data for contact_labels table
INSERT INTO contact_labels (
    name, code, description, color, icon, is_active, sort_order, 
    contact_count, is_system, status, type, created_at, updated_at
) VALUES 
(
    'VIP', 'VIP', 'Khách hàng VIP', '#722ed1', 'CrownOutlined', 
    true, 1, 25, true, 'active', 'system', 
    NOW(), NOW()
),
(
    'Quan trọng', 'IMPORTANT', 'Liên hệ quan trọng', '#faad14', 'StarOutlined', 
    true, 2, 18, true, 'active', 'system', 
    NOW(), NOW()
),
(
    'Khẩn cấp', 'URGENT', 'Liên hệ khẩn cấp', '#f5222d', 'ThunderboltOutlined', 
    true, 3, 8, true, 'active', 'system', 
    NOW(), NOW()
),
(
    'Khách hàng', 'CUSTOMER', 'Danh sách khách hàng', '#1890ff', 'UserOutlined', 
    true, 4, 45, false, 'active', 'custom', 
    NOW(), NOW()
),
(
    'Nhà cung cấp', 'SUPPLIER', 'Danh sách nhà cung cấp', '#13c2c2', 'ShopOutlined', 
    true, 5, 12, false, 'active', 'custom', 
    NOW(), NOW()
),
(
    'Đối tác', 'PARTNER', 'Đối tác kinh doanh', '#52c41a', 'TeamOutlined', 
    true, 6, 15, false, 'active', 'custom', 
    NOW(), NOW()
),
(
    'Gia đình', 'FAMILY', 'Liên hệ gia đình', '#fa8c16', 'HomeOutlined', 
    false, 7, 5, false, 'inactive', 'custom', 
    NOW(), NOW()
),
(
    'Bạn bè', 'FRIENDS', 'Danh sách bạn bè', '#eb2f96', 'HeartOutlined', 
    true, 8, 30, false, 'active', 'custom', 
    NOW(), NOW()
);
