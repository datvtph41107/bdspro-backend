DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
         WHERE table_name = 'catalog_plans' AND column_name = 'tier_rank'
    ) THEN
        RAISE EXCEPTION 'catalog_plans.tier_rank missing';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
         WHERE table_name = 'catalog_plan_versions' AND column_name = 'subscription_term_days'
    ) THEN
        RAISE EXCEPTION 'catalog_plan_versions.subscription_term_days missing';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
         WHERE table_name = 'catalog_subscriptions' AND column_name = 'pending_plan_version_id'
    ) THEN
        RAISE EXCEPTION 'catalog_subscriptions.pending_plan_version_id missing';
    END IF;
    IF to_regclass('commercial_checkout_commands') IS NULL THEN
        RAISE EXCEPTION 'commercial_checkout_commands missing';
    END IF;
END;
$$;
