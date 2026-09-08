# User Service schema authority

`user-service/database/migrations/` is the only production schema-evolution authority for User Service.

- `000001_user_schema.up.sql` builds the current fresh database schema.
- `000001_user_schema.down.sql` is only for disposable/local rollback labs.
- Serving startup defaults to `QHPRO_USER_DB_SCHEMA_MODE=sql`, so the process
  performs no DDL and the canonical runner remains the schema owner.
- `QHPRO_USER_DB_SCHEMA_MODE=automigrate` explicitly enables the current GORM
  model registry for local/development bootstrap. SQL and AutoMigrate are never
  executed by the same process startup.
- Historical migration/backfill/cutover SQL remains in source-control history
  and the legacy repository for adoption audits. It is not duplicated beside
  the canonical `database/migrations/` authority.

Do not add a second migration directory. New production schema changes must be created here with the next ordered version.
