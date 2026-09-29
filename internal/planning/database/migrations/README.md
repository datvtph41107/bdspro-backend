# TQD schema authority

`tqd-service/database/migrations/` is the only production schema-evolution authority for TQD.
Versions 000001-000026 preserve the prior v2 ordered history; versions
000027-000040 preserve the prior v3 history with one monotonic sequence.
Version 000041 closes the legacy gap where reachable GORM models depended on
ad-hoc AutoMigrate state not represented by the old SQL directories.
Versions 000042-000044 preserve the newer current Git history after that
canonicalization boundary; they must not be dropped merely because historical
checkpoint v13 stopped at version 000042.
Serving startup defaults to `QHPRO_TQD_DB_SCHEMA_MODE=sql`; the process then
executes no DDL and expects the canonical runner to have completed. Explicit
`QHPRO_TQD_DB_SCHEMA_MODE=automigrate` restores the GORM development bootstrap
registry (including PostGIS/unaccent setup) under a PostgreSQL advisory lock.
The two mutation mechanisms never run in the same startup. Historical loose
SQL and pre-canonical layouts remain available in source-control history and
the legacy repository; they are not copied into a second service-local schema
path.
