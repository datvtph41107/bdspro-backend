BEGIN;

-- Role grants already reference enum IDs 299-301. Ensure the durable rows exist
-- even when an earlier align migration was recorded without leaving the rows.

INSERT INTO permissions (id, name, key, description, permission_type, module, priority)
SELECT seed.id, seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    (299::bigint, 'Xem usage thương mại', 'COMMERCIAL_USAGE_VIEW', 'Đọc durable/runtime usage projection của TQD', 'ORGANIZATION', 'commercial', 10),
    (300::bigint, 'Reconcile usage thương mại', 'COMMERCIAL_USAGE_RECONCILE', 'Yêu cầu TQD đồng bộ runtime projection từ durable usage', 'ORGANIZATION', 'commercial', 20),
    (301::bigint, 'Xem subscription thương mại', 'COMMERCIAL_SUBSCRIPTION_VIEW', 'Đọc subscription và commercial projection của người dùng', 'ORGANIZATION', 'commercial', 10)
) AS seed(id, name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (
    SELECT 1 FROM permissions existing
     WHERE (existing.key = seed.key OR existing.id = seed.id)
       AND existing.deleted_at IS NULL
);

SELECT setval(
    pg_get_serial_sequence('permissions', 'id'),
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM permissions), 304)
);

COMMIT;
