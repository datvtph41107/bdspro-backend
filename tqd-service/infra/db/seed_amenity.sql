-- Seed data for amenities table
INSERT INTO amenities (
    id, name, description, icon, created_at, updated_at
) VALUES 
(1, 'Hồ bơi', 'Hồ bơi ngoài trời với hệ thống lọc nước hiện đại', 'HeartOutlined', '2024-01-15', '2024-01-20'),
(2, 'Gym', 'Phòng tập gym với đầy đủ thiết bị hiện đại', 'BuildOutlined', '2024-01-10', '2024-01-18'),
(3, 'Bãi đỗ xe', 'Bãi đỗ xe rộng rãi với hệ thống bảo vệ 24/7', 'CarOutlined', '2024-01-08', '2024-01-16'),
(4, 'Siêu thị', 'Siêu thị mini tiện lợi phục vụ cư dân', 'ShopOutlined', '2024-01-05', '2024-01-15'),
(5, 'Ngân hàng', 'Chi nhánh ngân hàng phục vụ giao dịch tài chính', 'BankOutlined', '2024-01-12', '2024-01-16'),
(6, 'Trường học', 'Trường mầm non và tiểu học chất lượng cao', 'BookOutlined', '2024-01-20', '2024-01-25')
ON CONFLICT (id) DO NOTHING;

-- Reset sequence for amenities id
SELECT setval('amenities_id_seq', (SELECT MAX(id) FROM amenities));

