BEGIN;

DELETE FROM catalog_price_items
WHERE plan_version_id IN (
    SELECT version.id
    FROM catalog_plan_versions version
    JOIN catalog_plans plan ON plan.id = version.plan_id
    WHERE plan.code IN ('qhpro.basic', 'qhpro.pro') AND version.version = '1.0.0'
);
DELETE FROM catalog_plan_operation_policies
WHERE plan_version_id IN (
    SELECT version.id
    FROM catalog_plan_versions version
    JOIN catalog_plans plan ON plan.id = version.plan_id
    WHERE plan.code IN ('qhpro.basic', 'qhpro.pro') AND version.version = '1.0.0'
);
DELETE FROM catalog_plan_entitlements
WHERE plan_version_id IN (
    SELECT version.id
    FROM catalog_plan_versions version
    JOIN catalog_plans plan ON plan.id = version.plan_id
    WHERE plan.code IN ('qhpro.basic', 'qhpro.pro') AND version.version = '1.0.0'
);
DELETE FROM catalog_plan_versions
WHERE plan_id IN (SELECT id FROM catalog_plans WHERE code IN ('qhpro.basic', 'qhpro.pro'))
  AND version = '1.0.0';
DELETE FROM catalog_plans WHERE code IN ('qhpro.basic', 'qhpro.pro');
DELETE FROM catalog_products product
WHERE product.code = 'qhpro'
  AND NOT EXISTS (SELECT 1 FROM catalog_plans plan WHERE plan.product_id = product.id);

COMMIT;
