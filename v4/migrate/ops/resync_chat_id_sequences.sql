-- Resync chat BIGSERIAL/identity sequences after importing rows with explicit ids.
-- Chat tables are created by the chat module's GORM AutoMigrate and live in the
-- database's active schema (normally public).
--
--   psql "$DATABASE_DSN" -v ON_ERROR_STOP=1 -f v4/migrate/ops/resync_chat_id_sequences.sql
--
-- Safe to re-run. Missing tables and tables without an id sequence are skipped.

BEGIN;

CREATE OR REPLACE FUNCTION pg_temp.resync_chat_table_id_serial(p_table regclass) RETURNS void
LANGUAGE plpgsql
AS $$
DECLARE
    seq_name text;
    max_id bigint;
BEGIN
    IF p_table IS NULL OR to_regclass(p_table::text) IS NULL THEN
        RETURN;
    END IF;

    seq_name := pg_get_serial_sequence(p_table::text, 'id');
    IF seq_name IS NULL THEN
        RAISE NOTICE 'skip % (no serial on id)', p_table;
        RETURN;
    END IF;

    EXECUTE format('SELECT COALESCE(MAX(id), 0) FROM %s', p_table) INTO max_id;
    PERFORM setval(seq_name::regclass, GREATEST(max_id, 1), max_id > 0);
    RAISE NOTICE 'resynced % sequence % after id %', p_table, seq_name, max_id;
END;
$$;

SELECT pg_temp.resync_chat_table_id_serial(to_regclass('conversations'));
SELECT pg_temp.resync_chat_table_id_serial(to_regclass('messages'));
SELECT pg_temp.resync_chat_table_id_serial(to_regclass('message_reactions'));
SELECT pg_temp.resync_chat_table_id_serial(to_regclass('background_images'));

COMMIT;

DROP FUNCTION IF EXISTS pg_temp.resync_chat_table_id_serial(regclass);
