-- Resync BIGSERIAL/identity sequences for bdspro + CRM tables after seed, import, or copy with explicit ids.
-- Run on the v4 user/bdspro DSN (public schema or after domain split — adjust search_path if needed).
--
--   psql "$DATABASE_DSN" -v ON_ERROR_STOP=1 -f v4/migrate/ops/resync_bdspro_id_sequences.sql
--
-- Safe to re-run. Skips missing tables (to_regclass → NULL) and tables without an id serial.

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
SELECT pg_temp.resync_table_id_serial(to_regclass('property_info'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_location'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_land_info'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_building_info'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_statistic'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_lineage'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_user'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_media'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_edvidence'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_external_ref'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_identify'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_type'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_product'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_report'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_histories'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_audit_logs'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_impact_versions'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_lineage_history'));
SELECT pg_temp.resync_table_id_serial(to_regclass('property_identifiers'));

-- Products
SELECT pg_temp.resync_table_id_serial(to_regclass('products'));
SELECT pg_temp.resync_table_id_serial(to_regclass('product_price'));
SELECT pg_temp.resync_table_id_serial(to_regclass('product_user'));
SELECT pg_temp.resync_table_id_serial(to_regclass('product_notes'));
SELECT pg_temp.resync_table_id_serial(to_regclass('product_note_mentions'));
SELECT pg_temp.resync_table_id_serial(to_regclass('product_note_files'));
SELECT pg_temp.resync_table_id_serial(to_regclass('house_info'));
SELECT pg_temp.resync_table_id_serial(to_regclass('listings'));

-- Assets
SELECT pg_temp.resync_table_id_serial(to_regclass('assets'));
SELECT pg_temp.resync_table_id_serial(to_regclass('asset_legals'));

-- Deals & transactions
SELECT pg_temp.resync_table_id_serial(to_regclass('deals'));
SELECT pg_temp.resync_table_id_serial(to_regclass('deal_members'));
SELECT pg_temp.resync_table_id_serial(to_regclass('deal_products'));
SELECT pg_temp.resync_table_id_serial(to_regclass('deal_investments'));
SELECT pg_temp.resync_table_id_serial(to_regclass('deal_internal_notes'));
SELECT pg_temp.resync_table_id_serial(to_regclass('tx_transaction'));
SELECT pg_temp.resync_table_id_serial(to_regclass('tx_transaction_action'));
SELECT pg_temp.resync_table_id_serial(to_regclass('tx_contract_deals'));
SELECT pg_temp.resync_table_id_serial(to_regclass('tx_cost_types'));
SELECT pg_temp.resync_table_id_serial(to_regclass('tx_deal_costs'));

-- Appointments & banking
SELECT pg_temp.resync_table_id_serial(to_regclass('appointments'));
SELECT pg_temp.resync_table_id_serial(to_regclass('appointment_participants'));
SELECT pg_temp.resync_table_id_serial(to_regclass('bank_accounts'));

-- CRM pipeline (modules/crm: PipelineEntity, StageEntity, LeadEntity)
SELECT pg_temp.resync_table_id_serial(to_regclass('pipeline'));
SELECT pg_temp.resync_table_id_serial(to_regclass('stage'));
SELECT pg_temp.resync_table_id_serial(to_regclass('customers'));

COMMIT;

DROP FUNCTION IF EXISTS pg_temp.resync_table_id_serial(regclass);
