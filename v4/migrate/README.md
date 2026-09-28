# User Service schema authority

`user-service/migrate/` is the only production schema-evolution authority for User Service.

- `000001_user_schema.up.sql` builds the current fresh database schema.
- `000001_user_schema.down.sql` is only for disposable/local rollback labs.
- Serving startup always applies pending `migrate/*.up.sql` via
  `RunPendingMigrations` before traffic is accepted (same pattern as hub-service).
- `QHPRO_USER_DB_SCHEMA_MODE=sql` (default) skips GORM AutoMigrate after SQL
  has already been applied; the canonical `migrate/` runner remains schema owner.
- `QHPRO_USER_DB_SCHEMA_MODE=automigrate` explicitly enables the current GORM
  model registry for local/development bootstrap. SQL and AutoMigrate are never
  executed by the same process startup.
- Historical migration/backfill/cutover SQL remains in source-control history
  and the legacy repository for adoption audits. It is not duplicated beside
  the canonical `migrate/` authority.

Do not add a second migration directory. New production schema changes must be created here with the next ordered version.
