BEGIN;

-- Align payment commercial permission primary keys with shared/common/domain/enum:
-- PAYMENT_ORDER_VIEW=305, PAYMENT_FULFILLMENT_VIEW=306,
-- PAYMENT_FULFILLMENT_REDRIVE=307, PAYMENT_FUNDS_CONFIRM=308.

INSERT INTO permissions (id, name, key, description, permission_type, module, priority)
SELECT seed.id, seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    (305::bigint, 'Xem payment thương mại', 'PAYMENT_ORDER_VIEW', 'Đọc order, payment attempt và settlement evidence', 'ORGANIZATION', 'commercial', 10),
    (306::bigint, 'Xem fulfillment thương mại', 'PAYMENT_FULFILLMENT_VIEW', 'Đọc trạng thái cấp quyền và hàng đợi fulfillment', 'ORGANIZATION', 'commercial', 10),
    (307::bigint, 'Redrive fulfillment thương mại', 'PAYMENT_FULFILLMENT_REDRIVE', 'Đưa fulfillment cần review trở lại hàng đợi xử lý', 'ORGANIZATION', 'commercial', 30),
    (308::bigint, 'Xác nhận tiền thủ công', 'PAYMENT_FUNDS_CONFIRM', 'Xác nhận funds confirmation thủ công khi provider evidence thiếu', 'ORGANIZATION', 'commercial', 20)
) AS seed(id, name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (
    SELECT 1 FROM permissions existing
     WHERE (existing.key = seed.key OR existing.id = seed.id)
       AND existing.deleted_at IS NULL
);

DO $$
DECLARE
    rec RECORD;
    occupant_key TEXT;
BEGIN
    FOR rec IN
        SELECT p.id AS old_id, seed.target_id, seed.permission_key
          FROM permissions p
          JOIN (VALUES
              ('PAYMENT_ORDER_VIEW', 305::bigint),
              ('PAYMENT_FULFILLMENT_VIEW', 306::bigint),
              ('PAYMENT_FULFILLMENT_REDRIVE', 307::bigint),
              ('PAYMENT_FUNDS_CONFIRM', 308::bigint)
          ) AS seed(permission_key, target_id) ON seed.permission_key = p.key
         WHERE p.deleted_at IS NULL
           AND p.id <> seed.target_id
    LOOP
        SELECT key INTO occupant_key
          FROM permissions
         WHERE id = rec.target_id
           AND deleted_at IS NULL
           AND key IS DISTINCT FROM rec.permission_key;

        IF occupant_key IS NOT NULL THEN
            RAISE EXCEPTION
                'cannot remap % to id %: id already owned by key %',
                rec.permission_key, rec.target_id, occupant_key;
        END IF;

        DELETE FROM permissions
         WHERE id = rec.target_id
           AND deleted_at IS NOT NULL;

        UPDATE role_permissions
           SET permission_id = rec.target_id
         WHERE permission_id = rec.old_id;

        UPDATE group_permissions
           SET permission_id = rec.target_id
         WHERE permission_id = rec.old_id;

        UPDATE permissions
           SET id = rec.target_id, updated_at = NOW()
         WHERE id = rec.old_id;
    END LOOP;
END $$;

-- Roles that already operate the commercial admin plane get the payment view/control
-- grants so Auth RequiredPermissions(id) matches existing operator assignments.
INSERT INTO role_permissions (role_id, permission_id)
SELECT DISTINCT rp.role_id, seed.permission_id
  FROM role_permissions rp
  JOIN (VALUES (305::bigint), (306::bigint), (307::bigint), (308::bigint)) AS seed(permission_id)
    ON TRUE
 WHERE rp.permission_id = 301
   AND NOT EXISTS (
       SELECT 1 FROM role_permissions existing
        WHERE existing.role_id = rp.role_id
          AND existing.permission_id = seed.permission_id
   );

SELECT setval(
    pg_get_serial_sequence('permissions', 'id'),
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM permissions), 308)
);

COMMIT;
