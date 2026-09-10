BEGIN;

DELETE FROM role_permissions
 WHERE permission_id IN (
    SELECT id FROM permissions WHERE key = 'PAYMENT_FUNDS_CONFIRM'
 );
DELETE FROM permissions WHERE key = 'PAYMENT_FUNDS_CONFIRM';

COMMIT;
