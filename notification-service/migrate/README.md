# Notification schema authority

`notification-service/migrate/` is the only production schema-evolution authority.
Serving startup defaults to `QHPRO_NOTIFICATION_DB_SCHEMA_MODE=sql` and executes
no DDL. Explicit `automigrate` mode runs the complete current domain + durable
eventing GORM registry under an advisory lock for local bootstrap. Historical
one-off SQL, when needed for audit/adoption, is recovered from source history
and is not executed by normal `make migrate-up`.
