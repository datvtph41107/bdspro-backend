-- Ensure PostgreSQL type-providing extensions for TQD (geometry/geography + text search).
-- Idempotent: safe to re-run on every environment that opens database.tqd.dsn.
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
