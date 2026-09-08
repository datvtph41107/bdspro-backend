-- ============================================================
-- Insert dữ liệu Tỉnh/Thành phố cấp 1 vào bảng region
-- ============================================================
-- Ngày tạo: 2025-01-08
-- Mô tả: Insert 35 tỉnh/thành phố trực thuộc trung ương
-- ============================================================

-- Xóa dữ liệu cũ nếu có (optional - comment lại nếu không muốn xóa)
-- DELETE FROM region WHERE level = 1;

-- Insert dữ liệu tỉnh/thành phố cấp 1
INSERT INTO region (name, code, code_name, level, parent_id, unit, created_at, updated_at) VALUES
-- 1. Thành phố Hà Nội
('Thành phố Hà Nội', 01, 'thanh-pho-ha-noi', 1, NULL, 'Thành phố Trung ương', NOW(), NOW()),

-- 2. Tỉnh Cao Bằng
('Tỉnh Cao Bằng', 04, 'tinh-cao-bang', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 3. Tỉnh Tuyên Quang
('Tỉnh Tuyên Quang', 08, 'tinh-tuyen-quang', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 4. Tỉnh Điện Biên
('Tỉnh Điện Biên', 11, 'tinh-dien-bien', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 5. Tỉnh Lai Châu
('Tỉnh Lai Châu', 12, 'tinh-lai-chau', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 6. Tỉnh Sơn La
('Tỉnh Sơn La', 14, 'tinh-son-la', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 7. Tỉnh Lào Cai
('Tỉnh Lào Cai', 15, 'tinh-lao-cai', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 8. Tỉnh Thái Nguyên
('Tỉnh Thái Nguyên', 19, 'tinh-thai-nguyen', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 9. Tỉnh Lạng Sơn
('Tỉnh Lạng Sơn', 20, 'tinh-lang-son', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 10. Tỉnh Quảng Ninh
('Tỉnh Quảng Ninh', 22, 'tinh-quang-ninh', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 11. Tỉnh Bắc Ninh
('Tỉnh Bắc Ninh', 24, 'tinh-bac-ninh', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 12. Tỉnh Phú Thọ
('Tỉnh Phú Thọ', 25, 'tinh-phu-tho', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 13. Thành phố Hải Phòng
('Thành phố Hải Phòng', 31, 'thanh-pho-hai-phong', 1, NULL, 'Thành phố Trung ương', NOW(), NOW()),

-- 14. Tỉnh Hưng Yên
('Tỉnh Hưng Yên', 33, 'tinh-hung-yen', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 15. Tỉnh Ninh Bình
('Tỉnh Ninh Bình', 37, 'tinh-ninh-binh', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 16. Tỉnh Thanh Hóa
('Tỉnh Thanh Hóa', 38, 'tinh-thanh-hoa', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 17. Tỉnh Nghệ An
('Tỉnh Nghệ An', 40, 'tinh-nghe-an', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 18. Tỉnh Hà Tĩnh
('Tỉnh Hà Tĩnh', 42, 'tinh-ha-tinh', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 19. Tỉnh Quảng Trị
('Tỉnh Quảng Trị', 44, 'tinh-quang-tri', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 20. Thành phố Huế
('Thành phố Huế', 46, 'thanh-pho-hue', 1, NULL, 'Thành phố Trung ương', NOW(), NOW()),

-- 21. Thành phố Đà Nẵng
('Thành phố Đà Nẵng', 48, 'thanh-pho-da-nang', 1, NULL, 'Thành phố Trung ương', NOW(), NOW()),

-- 22. Tỉnh Quảng Ngãi
('Tỉnh Quảng Ngãi', 51, 'tinh-quang-ngai', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 23. Tỉnh Gia Lai
('Tỉnh Gia Lai', 52, 'tinh-gia-lai', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 24. Tỉnh Khánh Hòa
('Tỉnh Khánh Hòa', 56, 'tinh-khanh-hoa', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 25. Tỉnh Đắk Lắk
('Tỉnh Đắk Lắk', 66, 'tinh-dak-lak', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 26. Tỉnh Lâm Đồng
('Tỉnh Lâm Đồng', 68, 'tinh-lam-dong', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 27. Tỉnh Đồng Nai
('Tỉnh Đồng Nai', 75, 'tinh-dong-nai', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 28. Thành phố Hồ Chí Minh
('Thành phố Hồ Chí Minh', 79, 'thanh-pho-ho-chi-minh', 1, NULL, 'Thành phố Trung ương', NOW(), NOW()),

-- 29. Tỉnh Tây Ninh
('Tỉnh Tây Ninh', 80, 'tinh-tay-ninh', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 30. Tỉnh Đồng Tháp
('Tỉnh Đồng Tháp', 82, 'tinh-dong-thap', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 31. Tỉnh Vĩnh Long
('Tỉnh Vĩnh Long', 86, 'tinh-vinh-long', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 32. Tỉnh An Giang
('Tỉnh An Giang', 91, 'tinh-an-giang', 1, NULL, 'Tỉnh', NOW(), NOW()),

-- 33. Thành phố Cần Thơ
('Thành phố Cần Thơ', 92, 'thanh-pho-can-tho', 1, NULL, 'Thành phố Trung ương', NOW(), NOW()),

-- 34. Tỉnh Cà Mau
('Tỉnh Cà Mau', 96, 'tinh-ca-mau', 1, NULL, 'Tỉnh', NOW(), NOW());

-- Kiểm tra số lượng bản ghi đã insert
SELECT COUNT(*) as total_provinces FROM region WHERE level = 1;

-- Xem danh sách các tỉnh/thành phố đã insert
SELECT id, name, code, code_name, level, unit 
FROM region 
WHERE level = 1 
ORDER BY code;

-- ============================================================
-- HƯỚNG DẪN SỬ DỤNG
-- ============================================================
-- 1. Kết nối database bdspro:
--    psql -U bdspro -h 14.225.210.29 -p 5432 -d db_bdspro
--
-- 2. Chạy script:
--    \i /path/to/insert_provinces.sql
--
-- 3. Hoặc từ command line:
--    psql -U bdspro -h 14.225.210.29 -p 5432 -d db_bdspro -f insert_provinces.sql
--
-- LƯU Ý:
-- - Script này insert 35 tỉnh/thành phố cấp 1 (level = 1)
-- - parent_id = NULL vì đây là cấp cao nhất
-- - code_name được tạo từ name bằng cách chuyển thành slug
-- - unit là loại hành chính (Tỉnh / Thành phố Trung ương)
-- ============================================================

