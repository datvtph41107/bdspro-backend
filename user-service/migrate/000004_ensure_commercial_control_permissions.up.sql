BEGIN;

-- Existing production databases may already have applied the original user
-- schema before Catalog and TQD usage control-plane permissions were added.
-- Keep permission ownership in User/IAM and make the forward repair idempotent.
INSERT INTO permissions (name, key, description, permission_type, module, priority)
SELECT seed.name, seed.key, seed.description, seed.permission_type, seed.module, seed.priority
FROM (VALUES
    ('Xem phiên bản catalog', 'CATALOG_PLAN_VIEW', 'Đọc plan version và điều khoản thương mại', 'ORGANIZATION', 'commercial', 10),
    ('Xuất bản phiên bản catalog', 'CATALOG_PLAN_PUBLISH', 'Xuất bản draft plan version', 'ORGANIZATION', 'commercial', 20),
    ('Xem usage thương mại', 'COMMERCIAL_USAGE_VIEW', 'Đọc durable/runtime usage projection của TQD', 'ORGANIZATION', 'commercial', 10),
    ('Reconcile usage thương mại', 'COMMERCIAL_USAGE_RECONCILE', 'Yêu cầu TQD đồng bộ runtime projection từ durable usage', 'ORGANIZATION', 'commercial', 20)
) AS seed(name, key, description, permission_type, module, priority)
WHERE NOT EXISTS (SELECT 1 FROM permissions existing WHERE existing.key = seed.key);

COMMIT;
