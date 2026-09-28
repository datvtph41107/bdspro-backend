BEGIN;

-- Commercial terms are operator-owned data. The original Basic/Pro rows were
-- useful while reconstructing the value chain, but leaving them active makes a
-- fresh installation sell source-defined prices before an operator approves
-- them. Preserve the immutable rows as history and remove only their ability
-- to appear in the public selling catalog.
UPDATE catalog_plan_versions AS version
SET status = 'retired',
    effective_until = COALESCE(version.effective_until, CURRENT_TIMESTAMP),
    updated_at = CURRENT_TIMESTAMP
FROM catalog_plans AS plan
JOIN catalog_products AS product ON product.id = plan.product_id
WHERE version.plan_id = plan.id
  AND product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro')
  AND version.version = '1.0.0'
  AND version.source = 'migration';

UPDATE catalog_plans AS plan
SET status = 'retired', updated_at = CURRENT_TIMESTAMP
FROM catalog_products AS product
WHERE product.id = plan.product_id
  AND product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro')
  AND NOT EXISTS (
      SELECT 1 FROM catalog_plan_versions AS version
      WHERE version.plan_id = plan.id AND version.status = 'active'
  );

COMMIT;
