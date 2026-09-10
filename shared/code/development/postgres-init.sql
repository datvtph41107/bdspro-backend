-- Local physical PostgreSQL/PostGIS bootstrap only.
--
-- Each application still owns its database and versioned migrations. This
-- script creates empty ownership boundaries; it contains no business schema
-- and is executed only when the development postgres-data volume is new.

SELECT 'CREATE ROLE qhpro WITH LOGIN SUPERUSER PASSWORD ''qhpro'''
WHERE NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'qhpro')
\gexec

SELECT 'CREATE DATABASE user_service OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'user_service')
\gexec

SELECT 'CREATE DATABASE organization OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'organization')
\gexec

SELECT 'CREATE DATABASE payment_service OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'payment_service')
\gexec

SELECT 'CREATE DATABASE qhpro_tqd OWNER qhpro'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'qhpro_tqd')
\gexec

SELECT 'CREATE DATABASE qhpro_notification OWNER qhpro'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'qhpro_notification')
\gexec

SELECT 'CREATE DATABASE file_service OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'file_service')
\gexec

SELECT 'CREATE DATABASE hub_service OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'hub_service')
\gexec

-- Source-health modules are not standing applications in the default Compose
-- topology, but their canonical migrations must be runnable against the same
-- local PostgreSQL when a developer works on those modules.
SELECT 'CREATE DATABASE db_bdspro OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'db_bdspro')
\gexec

SELECT 'CREATE DATABASE db_crm OWNER postgres'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'db_crm')
\gexec
