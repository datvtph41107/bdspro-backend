BEGIN;

-- Migration bổ sung cho database develop đã áp dụng schema nền trước khi
-- Commercial Admin hoàn tất. Fresh database vẫn an toàn nhờ key idempotent.
INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    ('Xem subscription thương mại', 'COMMERCIAL_SUBSCRIPTION_VIEW', 'Đọc subscription và commercial projection của người dùng', 'ORGANIZATION', 'commercial', 10),
    ('Xem payment thương mại', 'PAYMENT_ORDER_VIEW', 'Đọc order, payment attempt và settlement evidence', 'ORGANIZATION', 'commercial', 10),
    ('Xem fulfillment thương mại', 'PAYMENT_FULFILLMENT_VIEW', 'Đọc trạng thái cấp quyền và hàng đợi fulfillment', 'ORGANIZATION', 'commercial', 10),
    ('Redrive fulfillment thương mại', 'PAYMENT_FULFILLMENT_REDRIVE', 'Đưa fulfillment cần review trở lại hàng đợi xử lý', 'ORGANIZATION', 'commercial', 30)
) AS seed(name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (SELECT 1 FROM permissions existing WHERE existing.key = seed.key);

COMMIT;
