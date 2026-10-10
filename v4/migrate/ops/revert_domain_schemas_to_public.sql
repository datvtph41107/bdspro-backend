-- =============================================================================
-- Ops script: revert split_public_into_domain_schemas.sql (db_v4)
-- =============================================================================
-- After the split, the app still resolved tables through search_path=public and
-- v4 SQL migrations re-created empty copies in public. Reads then hit the empty
-- copies (e.g. public.user_profile = 0 rows, "user".user_profile = real data).
--
-- What this does, in one transaction:
--   1) Public copies of moved tables -> schema split_backup_20261009 (not dropped;
--      FKs between copies move with them).
--   2) Every table in the domain schemas -> public (data, indexes, owned
--      sequences, constraints move with the table).
--   3) Adds columns that only existed on the public copies (added by v4
--      migrations after the split).
--   4) Re-creates missing primary keys from the public copies. Rows that are
--      exact duplicates are removed first; tables whose keys still collide
--      (same key, different content) are skipped with a WARNING.
--   Aborts if any table is left behind in a domain schema.
--
-- Undo: move the restored tables back to their domain schemas, then
--   ALTER TABLE split_backup_20261009.<t> SET SCHEMA public for each copy.
--
-- Run:
--   psql "$DSN" -v ON_ERROR_STOP=1 -f revert_domain_schemas_to_public.sql
-- =============================================================================

BEGIN;

CREATE SCHEMA IF NOT EXISTS split_backup_20261009;

DO $$
DECLARE
  r record;
  domain_schemas text[] := ARRAY['auth','user','org','hub','bdspro','crm','seo',
                                 'notification','catalog','tqd','tx','legacy'];
BEGIN
  FOR r IN
    SELECT n.nspname AS sch, c.relname AS tbl
    FROM pg_class c
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE c.relkind IN ('r', 'p')
      AND NOT c.relispartition
      AND n.nspname = ANY (domain_schemas)
    ORDER BY n.nspname, c.relname
  LOOP
    IF EXISTS (
      SELECT 1 FROM pg_class pc
      WHERE pc.relnamespace = 'public'::regnamespace
        AND pc.relname = r.tbl
        AND pc.relkind IN ('r', 'p')
    ) THEN
      EXECUTE format('ALTER TABLE public.%I SET SCHEMA split_backup_20261009', r.tbl);
    END IF;
    EXECUTE format('ALTER TABLE %I.%I SET SCHEMA public', r.sch, r.tbl);
  END LOOP;

  IF EXISTS (
    SELECT 1 FROM pg_class c
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE c.relkind IN ('r', 'p') AND n.nspname = ANY (domain_schemas)
  ) THEN
    RAISE EXCEPTION 'tables still left in domain schemas; aborting';
  END IF;
END $$;

-- Columns present only on the public copies (v4 migrations after the split).
-- NOT NULL without a default cannot be added to populated tables, so those stay nullable.
ALTER TABLE public.click_urls
    ADD COLUMN IF NOT EXISTS id BIGSERIAL;

ALTER TABLE public.provinces
    ADD COLUMN IF NOT EXISTS codename VARCHAR(100),
    ADD COLUMN IF NOT EXISTS short_codename VARCHAR(100);

ALTER TABLE public.districts
    ADD COLUMN IF NOT EXISTS codename VARCHAR(100),
    ADD COLUMN IF NOT EXISTS short_codename VARCHAR(100);

ALTER TABLE public.wards
    ADD COLUMN IF NOT EXISTS codename VARCHAR(100),
    ADD COLUMN IF NOT EXISTS short_codename VARCHAR(100),
    ADD COLUMN IF NOT EXISTS type INTEGER,
    ADD COLUMN IF NOT EXISTS type_text VARCHAR(50);

ALTER TABLE public.pack_price_table
    ADD COLUMN IF NOT EXISTS price DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE public.purpose_use_profile
    ADD COLUMN IF NOT EXISTS id BIGSERIAL,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS created_by BIGINT,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Most restored tables were copied without constraints. Re-create the primary key
-- the public copy had, but only when the data is non-null and (after dropping
-- exact duplicate rows) unique.
DO $$
DECLARE
  r record;
  cols text;
  null_check text;
  missing int;
  bad bigint;
  removed bigint;
BEGIN
  FOR r IN
    SELECT pc.relname AS tbl, k.conkey, bc.oid AS backup_oid, pc.oid AS restored_oid
    FROM pg_class bc
    JOIN pg_constraint k ON k.conrelid = bc.oid AND k.contype = 'p'
    JOIN pg_class pc ON pc.relname = bc.relname
                    AND pc.relnamespace = 'public'::regnamespace
                    AND pc.relkind IN ('r', 'p')
    WHERE bc.relnamespace = 'split_backup_20261009'::regnamespace
      AND NOT EXISTS (
        SELECT 1 FROM pg_constraint k2 WHERE k2.conrelid = pc.oid AND k2.contype = 'p'
      )
    ORDER BY pc.relname
  LOOP
    SELECT string_agg(quote_ident(a.attname), ', ' ORDER BY array_position(r.conkey, a.attnum)),
           string_agg(quote_ident(a.attname) || ' IS NULL', ' OR '),
           count(*) FILTER (WHERE NOT EXISTS (
             SELECT 1 FROM pg_attribute ra
             WHERE ra.attrelid = r.restored_oid AND ra.attname = a.attname AND NOT ra.attisdropped))
      INTO cols, null_check, missing
    FROM pg_attribute a
    WHERE a.attrelid = r.backup_oid AND a.attnum = ANY (r.conkey);

    IF missing > 0 THEN
      RAISE WARNING 'skip PK on %: column(s) % missing', r.tbl, cols;
      CONTINUE;
    END IF;

    EXECUTE format('SELECT count(*) FROM public.%I WHERE %s', r.tbl, null_check) INTO bad;
    IF bad > 0 THEN
      RAISE WARNING 'skip PK on % (%): % row(s) with NULL', r.tbl, cols, bad;
      CONTINUE;
    END IF;

    EXECUTE format(
      'SELECT count(*) FROM (SELECT 1 FROM public.%I GROUP BY %s HAVING count(*) > 1) d',
      r.tbl, cols) INTO bad;
    IF bad > 0 THEN
      EXECUTE format(
        'DELETE FROM public.%1$I WHERE ctid IN (
           SELECT ctid FROM (
             SELECT ctid, row_number() OVER (PARTITION BY x::text ORDER BY ctid) AS rn
             FROM public.%1$I x) s
           WHERE rn > 1)', r.tbl);
      GET DIAGNOSTICS removed = ROW_COUNT;
      IF removed > 0 THEN
        RAISE NOTICE 'removed % exact duplicate row(s) from %', removed, r.tbl;
      END IF;
      EXECUTE format(
        'SELECT count(*) FROM (SELECT 1 FROM public.%I GROUP BY %s HAVING count(*) > 1) d',
        r.tbl, cols) INTO bad;
    END IF;
    IF bad > 0 THEN
      RAISE WARNING 'skip PK on % (%): % duplicate key(s)', r.tbl, cols, bad;
      CONTINUE;
    END IF;

    EXECUTE format('ALTER TABLE public.%I ADD PRIMARY KEY (%s)', r.tbl, cols);
  END LOOP;
END $$;

COMMIT;
