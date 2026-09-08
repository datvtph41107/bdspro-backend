-- Permission nghiệp vụ cho organization subscription. Admin organization được
-- policy owner cấp toàn quyền; custom role phải được gán permission này rõ ràng.
INSERT INTO organization_permissions (name, key, description, permission_type, created_at, updated_at, is_deleted)
SELECT
    'Thanh toán gói dịch vụ tổ chức',
    'CHECKOUT_ORGANIZATION_SUBSCRIPTION',
    'Cho phép tạo order và payment attempt dùng quota của tổ chức',
    10,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    FALSE
WHERE NOT EXISTS (
    SELECT 1
    FROM organization_permissions
    WHERE key = 'CHECKOUT_ORGANIZATION_SUBSCRIPTION' AND is_deleted = FALSE
);
