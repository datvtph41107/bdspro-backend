BEGIN;

-- Down only drops tables created by the repair migration when they are empty
-- of operator data. Prefer leaving durable commercial evidence in place.
DROP TABLE IF EXISTS subscription_settlement_receipts;
DROP TABLE IF EXISTS subscription_checkout_snapshots;
DROP TABLE IF EXISTS catalog_subscription_events;
DROP TABLE IF EXISTS catalog_subscriptions;
DROP TABLE IF EXISTS catalog_price_items;
DROP TABLE IF EXISTS catalog_plan_operation_policies;
DROP TABLE IF EXISTS catalog_plan_entitlements;
DROP TABLE IF EXISTS catalog_plan_versions;
DROP TABLE IF EXISTS catalog_plans;
DROP TABLE IF EXISTS catalog_products;

COMMIT;
