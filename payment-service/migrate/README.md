# Payment Service schema authority

`payment-service/migrate/` is the only production schema evolution authority.
Serving processes default to `QHPRO_PAYMENT_DB_SCHEMA_MODE=sql` and execute no
DDL. Explicit `automigrate` mode runs the complete Bank, Wallet, and Commerce
GORM registry under an advisory lock for local/development bootstrap.

Current durable owners:
- Order / Settlement / Fulfillment / command effects
- Payment Attempt
- Payment completion Outbox

Historical SQL remains available in source-control history and the legacy
repository for audit/adoption comparison; it is not a competing local path.
