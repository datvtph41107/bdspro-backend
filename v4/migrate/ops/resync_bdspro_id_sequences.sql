-- Resync BIGSERIAL/identity sequences for bdspro tables after seed, import, or copy with explicit ids.
-- Run on the v4 user/bdspro DSN (public schema or after domain split — adjust search_path if needed).
--
--   psql "$DATABASE_DSN" -v ON_ERROR_STOP=1 -f v4/migrate/ops/resync_bdspro_id_sequences.sql
--
-- Safe to re-run. Skips missing tables and tables without an id serial.

BEGIN;

CREATE OR REPLACE FUNCTION pg_temp.resync_table_id_serial(p_table regclass) RETURNS void
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
    RAISE NOTICE 'resync % sequence % -> next id after MAX(%)', p_table, seq_name, max_id;
END;
$$;

-- Property graph (create / property/new)
SELECT pg_temp.resync_table_id_serial('property_info'::regclass);
SELECT pg_temp.resync_table_id_serial('property_location'::regclass);
SELECT pg_temp.resync_table_id_serial('property_land_info'::regclass);
SELECT pg_temp.resync_table_id_serial('property_building_info'::regclass);
SELECT pg_temp.resync_table_id_serial('property_statistic'::regclass);
SELECT pg_temp.resync_table_id_serial('property_lineage'::regclass);
SELECT pg_temp.resync_table_id_serial('property_user'::regclass);
SELECT pg_temp.resync_table_id_serial('property_media'::regclass);
SELECT pg_temp.resync_table_id_serial('property_edvidence'::regclass);
SELECT pg_temp.resync_table_id_serial('property_external_ref'::regclass);
SELECT pg_temp.resync_table_id_serial('property_identify'::regclass);
SELECT pg_temp.resync_table_id_serial('property_type'::regclass);
SELECT pg_temp.resync_table_id_serial('property_product'::regclass);
SELECT pg_temp.resync_table_id_serial('property_report'::regclass);
SELECT pg_temp.resync_table_id_serial('property_histories'::regclass);
SELECT pg_temp.resync_table_id_serial('property_audit_logs'::regclass);
SELECT pg_temp.resync_table_id_serial('property_impact_versions'::regclass);
SELECT pg_temp.resync_table_id_serial('property_lineage_history'::regclass);
SELECT pg_temp.resync_table_id_serial('property_identifiers'::regclass);

-- Products
SELECT pg_temp.resync_table_id_serial('products'::regclass);
SELECT pg_temp.resync_table_id_serial('product_price'::regclass);
SELECT pg_temp.resync_table_id_serial('product_user'::regclass);
SELECT pg_temp.resync_table_id_serial('product_notes'::regclass);
SELECT pg_temp.resync_table_id_serial('product_note_mentions'::regclass);
SELECT pg_temp.resync_table_id_serial('product_note_files'::regclass);
SELECT pg_temp.resync_table_id_serial('house_info'::regclass);
SELECT pg_temp.resync_table_id_serial('listings'::regclass);

-- Assets
SELECT pg_temp.resync_table_id_serial('assets'::regclass);
SELECT pg_temp.resync_table_id_serial('asset_legals'::regclass);

-- Deals & transactions
SELECT pg_temp.resync_table_id_serial('deals'::regclass);
SELECT pg_temp.resync_table_id_serial('deal_members'::regclass);
SELECT pg_temp.resync_table_id_serial('deal_products'::regclass);
SELECT pg_temp.resync_table_id_serial('deal_investments'::regclass);
SELECT pg_temp.resync_table_id_serial('deal_internal_notes'::regclass);
SELECT pg_temp.resync_table_id_serial('tx_transaction'::regclass);
SELECT pg_temp.resync_table_id_serial('tx_transaction_action'::regclass);
SELECT pg_temp.resync_table_id_serial('tx_contract_deals'::regclass);
SELECT pg_temp.resync_table_id_serial('tx_cost_types'::regclass);
SELECT pg_temp.resync_table_id_serial('tx_deal_costs'::regclass);

-- Appointments & banking
SELECT pg_temp.resync_table_id_serial('appointments'::regclass);
SELECT pg_temp.resync_table_id_serial('appointment_participants'::regclass);
SELECT pg_temp.resync_table_id_serial('bank_accounts'::regclass);

COMMIT;

DROP FUNCTION IF EXISTS pg_temp.resync_table_id_serial(regclass);
