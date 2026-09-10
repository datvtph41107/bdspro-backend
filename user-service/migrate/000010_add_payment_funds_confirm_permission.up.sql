BEGIN;

INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT 'Xác nhận tiền thủ công', 'PAYMENT_FUNDS_CONFIRM',
       'Ghi nhận bằng chứng đối soát thủ công để Payment xác nhận tiền và kích hoạt gói',
       'SYSTEM', 'commercial', 30
WHERE NOT EXISTS (
    SELECT 1 FROM permissions
     WHERE key = 'PAYMENT_FUNDS_CONFIRM' AND deleted_at IS NULL
);

COMMIT;
