UPDATE faqs
SET answer = 'Có. Chọn Chia sẻ để sao chép hoặc gửi liên kết. Liên kết mở trực tiếp trang chi tiết hồ sơ đồ án tại /ban-do/do-an/:id.',
    updated_at = NOW()
WHERE deleted_at IS NULL
  AND LOWER(COALESCE(group_key, '')) = 'planning-project-detail'
  AND LOWER(TRIM(question)) = LOWER(TRIM('Tôi có thể chia sẻ đồ án quy hoạch không?'));
