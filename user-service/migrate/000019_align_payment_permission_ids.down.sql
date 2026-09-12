BEGIN;

DELETE FROM role_permissions
 WHERE permission_id IN (305, 306, 307, 308);

DELETE FROM group_permissions
 WHERE permission_id IN (305, 306, 307, 308);

DELETE FROM permissions
 WHERE id IN (305, 306, 307, 308)
   AND key IN (
       'PAYMENT_ORDER_VIEW',
       'PAYMENT_FULFILLMENT_VIEW',
       'PAYMENT_FULFILLMENT_REDRIVE',
       'PAYMENT_FUNDS_CONFIRM'
   );

COMMIT;
