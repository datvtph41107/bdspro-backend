-- Seed idempotent FAQ public planning dùng bởi React Native.
INSERT INTO faqs (question, answer, group_key, created_at, updated_at)
SELECT v.question, v.answer, 'planning-project-detail', NOW(), NOW()
FROM (VALUES
 ('Làm thế nào để tìm một thửa đất?', 'Nhập số tờ, số thửa, địa chỉ hoặc tọa độ trong ô tìm kiếm trên bản đồ. Chọn kết quả để mở phần xem nhanh rồi bấm Xem chi tiết khi cần.'),
 ('Dữ liệu quy hoạch được cập nhật khi nào?', 'Ngày cập nhật được hiển thị trên phần xem nhanh và màn chi tiết. Những bản ghi chưa có ngày cập nhật sẽ được ghi rõ là chưa xác định, không tự tạo ngày giả.'),
 ('Làm thế nào để theo dõi một đồ án quy hoạch?', 'Mở phần xem nhanh hoặc chi tiết đồ án, sau đó chọn Theo dõi. Chức năng này yêu cầu người dùng đăng nhập.'),
 ('Tôi có thể chia sẻ đồ án quy hoạch không?', 'Có. Chọn Chia sẻ để sao chép hoặc gửi liên kết. Liên kết mở trực tiếp đúng đồ án trong khu vực Thư viện quy hoạch.')
) AS v(question, answer)
WHERE NOT EXISTS (
 SELECT 1 FROM faqs f
 WHERE f.deleted_at IS NULL AND LOWER(COALESCE(f.group_key,''))='planning-project-detail'
   AND LOWER(TRIM(f.question))=LOWER(TRIM(v.question))
);
CREATE INDEX IF NOT EXISTS idx_faqs_group_key_active
 ON faqs(LOWER(group_key), updated_at DESC) WHERE deleted_at IS NULL;
