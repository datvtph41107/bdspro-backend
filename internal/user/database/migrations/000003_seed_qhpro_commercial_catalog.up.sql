BEGIN;

-- Catalog develop tối thiểu để chạy trọn luồng checkout -> payment -> quota.
-- User Service là owner duy nhất của giá, điều khoản, feature và allowance.
INSERT INTO catalog_products (
    code, display_name, status, created_at, updated_at
) VALUES (
    'qhpro', 'QHPRO', 'active', TIMESTAMPTZ '2026-01-01 00:00:00+00', TIMESTAMPTZ '2026-01-01 00:00:00+00'
) ON CONFLICT (code) DO NOTHING;

INSERT INTO catalog_plans (
    product_id, code, tier_rank, status, created_at, updated_at
)
SELECT product.id, seed.code, seed.tier_rank, 'active',
       TIMESTAMPTZ '2026-01-01 00:00:00+00', TIMESTAMPTZ '2026-01-01 00:00:00+00'
FROM catalog_products product
CROSS JOIN (VALUES
    ('qhpro.basic', 10),
    ('qhpro.pro', 20)
) AS seed(code, tier_rank)
WHERE product.code = 'qhpro'
ON CONFLICT (product_id, code) DO NOTHING;

INSERT INTO catalog_plan_versions (
    plan_id, version, display_name, status, subject_scope,
    subscription_term_days, effective_from, published_at, source,
    terms_checksum, created_at, updated_at
)
SELECT plan.id, '1.0.0', seed.display_name, 'active', 'any',
       30, TIMESTAMPTZ '2026-01-01 00:00:00+00', TIMESTAMPTZ '2026-01-01 00:00:00+00',
       'migration', seed.terms_checksum,
       TIMESTAMPTZ '2026-01-01 00:00:00+00', TIMESTAMPTZ '2026-01-01 00:00:00+00'
FROM catalog_plans plan
JOIN catalog_products product ON product.id = plan.product_id
JOIN (VALUES
    ('qhpro.basic', 'Cơ bản', '13cf5445e5cc2c7bdca81eedbb7ee8e960d12c0a3eef5d184c4ae57339f692e2'),
    ('qhpro.pro', 'Chuyên nghiệp', 'fd6d9c47d4d68be465d28001cdc3ce81322e26ad09e1fc5300a66589d0339fc0')
) AS seed(plan_code, display_name, terms_checksum) ON seed.plan_code = plan.code
WHERE product.code = 'qhpro'
ON CONFLICT (plan_id, version) DO NOTHING;

INSERT INTO catalog_plan_entitlements (
    plan_version_id, code, kind, feature_code, meter_code,
    amount, unlimited, period, created_at
)
SELECT version.id, seed.code, seed.kind, seed.feature_code, seed.meter_code,
       seed.amount, false, seed.period, TIMESTAMPTZ '2026-01-01 00:00:00+00'
FROM catalog_plan_versions version
JOIN catalog_plans plan ON plan.id = version.plan_id
JOIN catalog_products product ON product.id = plan.product_id
CROSS JOIN LATERAL (VALUES
    ('workspace.report.feature', 'feature_access', 'workspace.report.generate', NULL, 0::BIGINT, 'none'),
    ('workspace.report.allowance', 'usage_allowance', NULL, 'workspace.report_generation.accepted',
        CASE plan.code WHEN 'qhpro.basic' THEN 3::BIGINT ELSE 100::BIGINT END, 'subscription_cycle')
) AS seed(code, kind, feature_code, meter_code, amount, period)
WHERE product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro')
  AND version.version = '1.0.0'
ON CONFLICT (plan_version_id, code) DO NOTHING;

INSERT INTO catalog_plan_operation_policies (
    plan_version_id, operation_code, feature_code, meter_code,
    units_per_action, created_at
)
SELECT version.id, 'workspace.report.generate', 'workspace.report.generate',
       'workspace.report_generation.accepted', 1,
       TIMESTAMPTZ '2026-01-01 00:00:00+00'
FROM catalog_plan_versions version
JOIN catalog_plans plan ON plan.id = version.plan_id
JOIN catalog_products product ON product.id = plan.product_id
WHERE product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro')
  AND version.version = '1.0.0'
ON CONFLICT (plan_version_id, operation_code) DO NOTHING;

INSERT INTO catalog_price_items (
    plan_version_id, code, kind, currency, amount_minor,
    billing_unit, meter_code, quantity, created_at
)
SELECT version.id, plan.code || '.monthly', 'recurring', 'VND',
       CASE plan.code WHEN 'qhpro.basic' THEN 99000::BIGINT ELSE 199000::BIGINT END,
       'subscription.cycle', NULL, 1,
       TIMESTAMPTZ '2026-01-01 00:00:00+00'
FROM catalog_plan_versions version
JOIN catalog_plans plan ON plan.id = version.plan_id
JOIN catalog_products product ON product.id = plan.product_id
WHERE product.code = 'qhpro'
  AND plan.code IN ('qhpro.basic', 'qhpro.pro')
  AND version.version = '1.0.0'
ON CONFLICT (plan_version_id, code) DO NOTHING;

COMMIT;
