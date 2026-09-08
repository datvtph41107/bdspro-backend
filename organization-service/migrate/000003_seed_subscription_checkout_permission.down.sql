DELETE FROM role_permissions
WHERE organization_permission_model_id IN (
    SELECT id
    FROM organization_permissions
    WHERE key = 'CHECKOUT_ORGANIZATION_SUBSCRIPTION'
      AND name = 'Thanh toán gói dịch vụ tổ chức'
);

DELETE FROM organization_permissions
WHERE key = 'CHECKOUT_ORGANIZATION_SUBSCRIPTION'
  AND name = 'Thanh toán gói dịch vụ tổ chức';
