BEGIN;

DO $$
DECLARE
    product_id_value BIGINT;
    plan_id_value BIGINT;
    version_id_value BIGINT;
BEGIN
    INSERT INTO catalog_products(code, display_name, status)
    VALUES ('qhpro.local', 'QHPro Local Lab', 'active')
    RETURNING id INTO product_id_value;

    INSERT INTO catalog_plans(product_id, code, status)
    VALUES (product_id_value, 'qhpro.local.smoke', 'active')
    RETURNING id INTO plan_id_value;

    INSERT INTO catalog_plan_versions(
        plan_id,
        version,
        display_name,
        status,
        subject_scope,
        source,
        terms_checksum
    ) VALUES (
        plan_id_value,
        '1.0.0',
        'QHPro Local Smoke Plan',
        'draft',
        'profile',
        'native',
        repeat('0', 64)
    ) RETURNING id INTO version_id_value;

    INSERT INTO catalog_plan_entitlements(
        plan_version_id,
        code,
        kind,
        feature_code,
        amount,
        unlimited,
        period
    ) VALUES (
        version_id_value,
        'workspace.report.generate.access',
        'feature_access',
        'workspace.report.generate',
        0,
        FALSE,
        'none'
    );

    INSERT INTO catalog_plan_entitlements(
        plan_version_id,
        code,
        kind,
        meter_code,
        amount,
        unlimited,
        period
    ) VALUES (
        version_id_value,
        'workspace.report_generation.accepted.allowance',
        'usage_allowance',
        'workspace.report_generation.accepted',
        3,
        FALSE,
        'subscription_cycle'
    );

    UPDATE catalog_plan_versions
       SET status = 'active',
           effective_from = NOW() - INTERVAL '1 day',
           published_at = NOW() - INTERVAL '1 day',
           published_by = 1,
           updated_by = 1,
           terms_checksum = repeat('a', 64)
     WHERE id = version_id_value;

    INSERT INTO catalog_subscriptions(
        subscription_key,
        subject_kind,
        subject_id,
        product_id,
        plan_version_id,
        status,
        started_at,
        current_period_start,
        current_period_end,
        access_until,
        auto_renew
    ) VALUES (
        'qhpro-local-profile-42',
        'profile',
        '42',
        product_id_value,
        version_id_value,
        'active',
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day',
        NOW() + INTERVAL '29 days',
        NOW() + INTERVAL '29 days',
        FALSE
    );
END $$;

COMMIT;
