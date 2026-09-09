BEGIN;

UPDATE catalog_plans AS plan
SET status = 'active', updated_at = CURRENT_TIMESTAMP
FROM catalog_products AS product
WHERE product.id = plan.product_id
  AND product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro');

UPDATE catalog_plan_versions AS version
SET status = 'active', effective_until = NULL, updated_at = CURRENT_TIMESTAMP
FROM catalog_plans AS plan
JOIN catalog_products AS product ON product.id = plan.product_id
WHERE version.plan_id = plan.id
  AND product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro')
  AND version.version = '1.0.0'
  AND version.source = 'migration';

COMMIT;
