-- Seed data cho bảng POI Hours (Giờ mở cửa của địa điểm)
-- Bảng này chứa thông tin giờ hoạt động của các POI theo từng ngày trong tuần

-- Xóa bảng nếu đã tồn tại
DROP TABLE IF EXISTS open_hour;

-- Tạo bảng open_hour
CREATE TABLE open_hour (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    poi_id BIGINT,
    day_of_week INT NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    note TEXT,
    is_open BOOLEAN DEFAULT true,
    status INT DEFAULT 0,
    type INT DEFAULT 0,
    
    -- Các trường từ BaseEntity
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    -- Các trường từ AuditBase
    created_by BIGINT,
    updated_by BIGINT
);

-- Tạo indexes
CREATE INDEX idx_deleted_at ON open_hour(deleted_at);
CREATE INDEX idx_day_of_week ON open_hour(day_of_week);
CREATE INDEX idx_poi_id ON open_hour(poi_id);

-- Insert seed data
INSERT INTO open_hour (name, poi_id, day_of_week, open_time, close_time, note, is_open, status, type, created_at, updated_at, deleted_at, created_by, updated_by) VALUES
-- Giờ làm việc bình thường (Thứ 2-6: 8h-22h)
('Giờ làm việc thứ 2', NULL, 2, '08:00:00', '22:00:00', 'Giờ làm việc thứ 2 - Ca thường', true, 0, 0, NOW(), NOW(), NULL, 1, 1),
('Giờ làm việc thứ 3', NULL, 3, '08:00:00', '22:00:00', 'Giờ làm việc thứ 3 - Ca thường', true, 0, 0, NOW(), NOW(), NULL, 1, 1),
('Giờ làm việc thứ 4', NULL, 4, '08:00:00', '22:00:00', 'Giờ làm việc thứ 4 - Ca thường', true, 0, 0, NOW(), NOW(), NULL, 1, 1),
('Giờ làm việc thứ 5', NULL, 5, '08:00:00', '22:00:00', 'Giờ làm việc thứ 5 - Ca thường', true, 0, 0, NOW(), NOW(), NULL, 1, 1),
('Giờ làm việc thứ 6', NULL, 6, '08:00:00', '22:00:00', 'Giờ làm việc thứ 6 - Ca thường', true, 0, 0, NOW(), NOW(), NULL, 1, 1),

-- Giờ cuối tuần (Thứ 7, CN: 8h-23h)
('Giờ cuối tuần thứ 7', NULL, 7, '08:00:00', '23:00:00', 'Giờ mở cửa thứ 7 - Cuối tuần', true, 0, 0, NOW(), NOW(), NULL, 1, 1),
('Giờ Chủ nhật', NULL, 8, '08:00:00', '23:00:00', 'Giờ mở cửa Chủ nhật - Cuối tuần', true, 0, 0, NOW(), NOW(), NULL, 1, 1),

-- Mở cửa sớm (6h-14h) - Ca sáng
('Ca sáng thứ 2', NULL, 2, '06:00:00', '14:00:00', 'Ca sáng thứ 2 - Mở cửa sớm', true, 0, 1, NOW(), NOW(), NULL, 1, 1),
('Ca sáng thứ 3', NULL, 3, '06:00:00', '14:00:00', 'Ca sáng thứ 3 - Quán ăn sáng', true, 0, 1, NOW(), NOW(), NULL, 1, 1),

-- Mở cửa chiều tối (14h-02h) - Ca chiều/tối
('Ca tối thứ 5', NULL, 5, '14:00:00', '02:00:00', 'Ca chiều thứ 5 - Đến tận khuya', true, 0, 1, NOW(), NOW(), NULL, 1, 1),
('Ca tối thứ 6', NULL, 6, '14:00:00', '02:00:00', 'Ca chiều thứ 6 - Bar/Pub', true, 0, 1, NOW(), NOW(), NULL, 1, 1),

-- Mở cửa cả ngày (24/7)
('24/7 thứ 2', NULL, 2, '00:00:00', '23:59:59', 'Thứ 2 - Mở cửa 24/7 (Cửa hàng tiện lợi)', true, 0, 0, NOW(), NOW(), NULL, 1, 1),
('24/7 thứ 7', NULL, 7, '00:00:00', '23:59:59', 'Thứ 7 - Siêu thị mở cửa cả ngày', true, 0, 0, NOW(), NOW(), NULL, 1, 1),

-- Giờ đặc biệt
('Giờ hành chính', NULL, 4, '10:00:00', '18:00:00', 'Thứ 4 - Giờ hành chính (Văn phòng)', true, 0, 1, NOW(), NOW(), NULL, 1, 1),
('Giờ ngắn CN', NULL, 8, '09:00:00', '17:00:00', 'Chủ nhật - Giờ làm việc ngắn', true, 0, 1, NOW(), NOW(), NULL, 1, 1),
('Giờ gym', NULL, 3, '07:00:00', '21:00:00', 'Thứ 3 - Phòng gym buổi sáng', true, 0, 0, NOW(), NOW(), NULL, 1, 1);

-- Lưu ý:
-- - day_of_week: 2 = Thứ 2, 3 = Thứ 3, ..., 7 = Thứ 7, 8 = Chủ nhật
-- - open_time, close_time: định dạng TIME (HH:MM:SS)
-- - note: ghi chú về giờ mở cửa
-- - created_at, updated_at: timestamp tự động
-- - deleted_at: NULL (không bị xóa mềm)
-- - created_by, updated_by: ID của user tạo/cập nhật (mặc định = 1)

