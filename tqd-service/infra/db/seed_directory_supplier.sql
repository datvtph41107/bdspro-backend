-- Seed data for directory_suppliers table
INSERT INTO directory_suppliers (
    name, code, description, contact_person, phone, email, address, website,
    tax_code, business_license, rating, is_active, service_count, contract_count,
    total_value, join_date, last_contact_date, notes, raw_tags, raw_categories,
    created_at, updated_at
) VALUES 
(
    'Công ty TNHH Xây dựng ABC', 'ABC-CONSTRUCTION', 
    'Chuyên cung cấp dịch vụ xây dựng và sửa chữa',
    'Nguyễn Văn A', '028 1234 5678', 'contact@abc-construction.com',
    '123 Đường ABC, Quận 1, TP.HCM', 'https://abc-construction.com',
    '0123456789', 'BL-2024-001', 4.5, true, 25, 8,
    5000000000, '2024-01-15', '2024-01-20',
    'Đối tác uy tín, thực hiện đúng tiến độ',
    '["Uy tín cao", "Giá cả hợp lý", "Dịch vụ tốt"]',
    '["CONSTRUCTION", "MAINTENANCE"]',
    NOW(), NOW()
),
(
    'Công ty Dịch vụ Vệ sinh XYZ', 'XYZ-CLEANING',
    'Dịch vụ vệ sinh chuyên nghiệp cho tòa nhà',
    'Trần Thị B', '028 8765 4321', 'info@xyz-cleaning.com',
    '456 Đường XYZ, Quận 2, TP.HCM', 'https://xyz-cleaning.com',
    '0987654321', 'BL-2024-002', 4.2, true, 18, 5,
    2000000000, '2024-01-10', '2024-01-18',
    'Dịch vụ vệ sinh chất lượng cao',
    '["Chuyên nghiệp", "Phản hồi nhanh"]',
    '["CLEANING"]',
    NOW(), NOW()
),
(
    'Công ty An ninh DEF', 'DEF-SECURITY',
    'Dịch vụ bảo vệ và an ninh 24/7',
    'Lê Văn C', '028 5555 6666', 'security@def-security.com',
    '789 Đường DEF, Quận 3, TP.HCM', 'https://def-security.com',
    '0555666777', 'BL-2024-003', 4.8, true, 12, 3,
    3000000000, '2024-01-08', '2024-01-16',
    'Đội ngũ bảo vệ chuyên nghiệp',
    '["Đáng tin cậy", "Kinh nghiệm lâu năm", "Hỗ trợ 24/7"]',
    '["SECURITY"]',
    NOW(), NOW()
),
(
    'Công ty Công nghệ GHI', 'GHI-TECH',
    'Giải pháp công nghệ thông tin và phần mềm',
    'Phạm Thị D', '028 7777 8888', 'tech@ghi-tech.com',
    '321 Đường GHI, Quận 7, TP.HCM', 'https://ghi-tech.com',
    '0777888999', 'BL-2024-004', 4.0, false, 8, 2,
    1500000000, '2024-01-05', '2024-01-12',
    'Công nghệ hiện đại nhưng giá cao',
    '["Công nghệ hiện đại"]',
    '["TECHNOLOGY"]',
    NOW(), NOW()
),
(
    'Công ty Catering JKL', 'JKL-CATERING',
    'Dịch vụ cung cấp suất ăn và tiệc',
    'Hoàng Văn E', '028 9999 0000', 'catering@jkl-catering.com',
    '654 Đường JKL, Quận 10, TP.HCM', 'https://jkl-catering.com',
    '0999000111', 'BL-2024-005', 4.3, true, 15, 4,
    800000000, '2024-01-12', '2024-01-19',
    'Thực phẩm tươi ngon, giá cả hợp lý',
    '["Thực phẩm tươi", "Giá cả hợp lý"]',
    '["CATERING"]',
    NOW(), NOW()
);
