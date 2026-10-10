-- =============================================================================
-- Ops script: split public tables into domain schemas (db_v4 inventory)
-- =============================================================================
-- Mục đích: gom bảng theo owner/domain, không để hết trong public.
--
-- CHÚ Ý TRƯỚC KHI CHẠY:
-- 1) Đây KHÔNG phải migration golang-migrate tự chạy. Chạy thủ công trên DB
--    staging trước; app phải cập nhật TableName / search_path / GORM schema
--    trước hoặc ngay sau khi move (nếu không sẽ 404/relation does not exist).
-- 2) FK cross-schema vẫn hoạt động trong PostgreSQL.
-- 3) Views/materialized views / functions tham chiếu public.X cần sửa sau.
-- 4) Giữ schema_migrations trong public.
-- 5) Backup trước. Script idempotent nhờ CREATE SCHEMA IF NOT EXISTS +
--    chỉ SET SCHEMA khi table còn ở public.
--
-- Cách chạy (ví dụ):
--   psql "$DSN" -v ON_ERROR_STOP=1 -f split_public_into_domain_schemas.sql
-- =============================================================================

BEGIN;

CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS "user";
CREATE SCHEMA IF NOT EXISTS org;
CREATE SCHEMA IF NOT EXISTS hub;
CREATE SCHEMA IF NOT EXISTS bdspro;
CREATE SCHEMA IF NOT EXISTS crm;
CREATE SCHEMA IF NOT EXISTS seo;
CREATE SCHEMA IF NOT EXISTS notification;
CREATE SCHEMA IF NOT EXISTS catalog;
CREATE SCHEMA IF NOT EXISTS tqd;
CREATE SCHEMA IF NOT EXISTS tx;
CREATE SCHEMA IF NOT EXISTS legacy; -- bảng cũ / backup / chưa rõ owner

-- Helper: move table only if it still lives in public.
CREATE OR REPLACE FUNCTION pg_temp.move_if_public(p_table text, p_schema text)
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_tables
    WHERE schemaname = 'public'
      AND tablename = p_table
  ) THEN
    EXECUTE format('ALTER TABLE public.%I SET SCHEMA %I', p_table, p_schema);
  END IF;
END;
$$;

-- -----------------------------------------------------------------------------
-- auth: đăng nhập / OTP / session / device / PIN
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('auth_config', 'auth');
SELECT pg_temp.move_if_public('auth_devices', 'auth');
SELECT pg_temp.move_if_public('auth_method', 'auth');
SELECT pg_temp.move_if_public('auth_status', 'auth');
SELECT pg_temp.move_if_public('history_auth', 'auth');
SELECT pg_temp.move_if_public('login_log', 'auth');
SELECT pg_temp.move_if_public('user_auth', 'auth');
SELECT pg_temp.move_if_public('user_otp', 'auth');
SELECT pg_temp.move_if_public('user_pin', 'auth');
SELECT pg_temp.move_if_public('user_session', 'auth');
SELECT pg_temp.move_if_public('user_status', 'auth');
SELECT pg_temp.move_if_public('api_keys', 'auth'); -- hub api keys cũng có thể để hub; tạm auth nếu dùng cho service auth

-- -----------------------------------------------------------------------------
-- user: profile / IAM role-permission / KYC / bookmark / catalog cá nhân
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('user_profile', 'user');
SELECT pg_temp.move_if_public('user_info', 'user');
SELECT pg_temp.move_if_public('user_stats', 'user');
SELECT pg_temp.move_if_public('user_activity_log', 'user');
SELECT pg_temp.move_if_public('profile_deleted', 'user');
SELECT pg_temp.move_if_public('profile_media', 'user');
SELECT pg_temp.move_if_public('profile_transfer', 'user');
SELECT pg_temp.move_if_public('admin_profiles', 'user');
SELECT pg_temp.move_if_public('admin_access_control', 'user');
SELECT pg_temp.move_if_public('admin_access_logs', 'user');
SELECT pg_temp.move_if_public('admin_histories', 'user');
SELECT pg_temp.move_if_public('roles', 'user');
SELECT pg_temp.move_if_public('role_groups', 'user');
SELECT pg_temp.move_if_public('role_permissions', 'user');
SELECT pg_temp.move_if_public('role_profiles', 'user');
SELECT pg_temp.move_if_public('permissions', 'user');
SELECT pg_temp.move_if_public('group_permissions', 'user');
SELECT pg_temp.move_if_public('kyc', 'user');
SELECT pg_temp.move_if_public('bookmark_user', 'user');
SELECT pg_temp.move_if_public('profession', 'user');
SELECT pg_temp.move_if_public('certification', 'user');
SELECT pg_temp.move_if_public('purpose_use', 'user');
SELECT pg_temp.move_if_public('purpose_use_profile', 'user');
SELECT pg_temp.move_if_public('main_area', 'user');
SELECT pg_temp.move_if_public('main_area_profile', 'user');
SELECT pg_temp.move_if_public('tag', 'user');
SELECT pg_temp.move_if_public('tag_user', 'user');
SELECT pg_temp.move_if_public('tags', 'user');
SELECT pg_temp.move_if_public('account_warnings', 'user');
SELECT pg_temp.move_if_public('account_warning_logs', 'user');
SELECT pg_temp.move_if_public('account_warning_templates', 'user');
SELECT pg_temp.move_if_public('demo_appointment_log', 'user');
SELECT pg_temp.move_if_public('demo_contact_log', 'user');
SELECT pg_temp.move_if_public('demo_deal_log', 'user');
SELECT pg_temp.move_if_public('demo_product_log', 'user');
SELECT pg_temp.move_if_public('person_configs', 'user');
SELECT pg_temp.move_if_public('bank_accounts', 'user');
SELECT pg_temp.move_if_public('user_saved_planning_news', 'user');

-- -----------------------------------------------------------------------------
-- org: tổ chức / membership
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('organizations', 'org');
SELECT pg_temp.move_if_public('organization_members', 'org');
SELECT pg_temp.move_if_public('memberships', 'org');

-- -----------------------------------------------------------------------------
-- hub: location v2, stats, FAQ, guides, system config, event queue
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('province_v2', 'hub');
SELECT pg_temp.move_if_public('province_v3', 'hub');
SELECT pg_temp.move_if_public('district_v2', 'hub');
SELECT pg_temp.move_if_public('ward_v2', 'hub');
SELECT pg_temp.move_if_public('ward_v3', 'hub');
SELECT pg_temp.move_if_public('provinces', 'hub');
SELECT pg_temp.move_if_public('districts', 'hub');
SELECT pg_temp.move_if_public('wards', 'hub');
SELECT pg_temp.move_if_public('request_stats_log', 'hub');
SELECT pg_temp.move_if_public('request_audit_log', 'hub');
SELECT pg_temp.move_if_public('user_interaction_log', 'hub');
SELECT pg_temp.move_if_public('keyword_search_log', 'hub');
SELECT pg_temp.move_if_public('search_latency_log', 'hub');
SELECT pg_temp.move_if_public('error_logs', 'hub');
SELECT pg_temp.move_if_public('faqs', 'hub');
SELECT pg_temp.move_if_public('tb_user_guide', 'hub');
SELECT pg_temp.move_if_public('user_guide_steps', 'hub');
SELECT pg_temp.move_if_public('applink', 'hub');
SELECT pg_temp.move_if_public('system_config', 'hub');
SELECT pg_temp.move_if_public('versions', 'hub');
SELECT pg_temp.move_if_public('update_data', 'hub');
SELECT pg_temp.move_if_public('click_urls', 'hub');
SELECT pg_temp.move_if_public('interactive_events', 'hub');
SELECT pg_temp.move_if_public('tb_event_queue', 'hub');
SELECT pg_temp.move_if_public('activity_log', 'hub');
SELECT pg_temp.move_if_public('observations', 'hub');

-- -----------------------------------------------------------------------------
-- bdspro: sản phẩm / deal / property / asset / post / appointment
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('products', 'bdspro');
SELECT pg_temp.move_if_public('product_amenity', 'bdspro');
SELECT pg_temp.move_if_public('product_asset', 'bdspro');
SELECT pg_temp.move_if_public('product_history', 'bdspro');
SELECT pg_temp.move_if_public('product_media', 'bdspro');
SELECT pg_temp.move_if_public('product_notes', 'bdspro');
SELECT pg_temp.move_if_public('product_note_files', 'bdspro');
SELECT pg_temp.move_if_public('product_note_mentions', 'bdspro');
SELECT pg_temp.move_if_public('product_organization', 'bdspro');
SELECT pg_temp.move_if_public('product_price', 'bdspro');
SELECT pg_temp.move_if_public('product_private', 'bdspro');
SELECT pg_temp.move_if_public('product_stats', 'bdspro');
SELECT pg_temp.move_if_public('product_user', 'bdspro');
SELECT pg_temp.move_if_public('amenity', 'bdspro');
SELECT pg_temp.move_if_public('deals', 'bdspro');
SELECT pg_temp.move_if_public('deal_history', 'bdspro');
SELECT pg_temp.move_if_public('deal_internal_notes', 'bdspro');
SELECT pg_temp.move_if_public('deal_investments', 'bdspro');
SELECT pg_temp.move_if_public('deal_members', 'bdspro');
SELECT pg_temp.move_if_public('deal_milestones', 'bdspro');
SELECT pg_temp.move_if_public('deal_products', 'bdspro');
SELECT pg_temp.move_if_public('deposite', 'bdspro');
SELECT pg_temp.move_if_public('distributions', 'bdspro');
SELECT pg_temp.move_if_public('assets', 'bdspro');
SELECT pg_temp.move_if_public('asset_costs', 'bdspro');
SELECT pg_temp.move_if_public('asset_cost_types', 'bdspro');
SELECT pg_temp.move_if_public('asset_exploitation', 'bdspro');
SELECT pg_temp.move_if_public('asset_income_types', 'bdspro');
SELECT pg_temp.move_if_public('asset_legals', 'bdspro');
SELECT pg_temp.move_if_public('asset_organization', 'bdspro');
SELECT pg_temp.move_if_public('asset_shares', 'bdspro');
SELECT pg_temp.move_if_public('asset_split_history', 'bdspro');
SELECT pg_temp.move_if_public('asset_stats', 'bdspro');
SELECT pg_temp.move_if_public('asset_user', 'bdspro');
SELECT pg_temp.move_if_public('cost_documents', 'bdspro');
SELECT pg_temp.move_if_public('income_documents', 'bdspro');
SELECT pg_temp.move_if_public('attachs', 'bdspro');
SELECT pg_temp.move_if_public('attach_documents', 'bdspro');
SELECT pg_temp.move_if_public('property_info', 'bdspro');
SELECT pg_temp.move_if_public('property_amenity', 'bdspro');
SELECT pg_temp.move_if_public('property_area_region', 'bdspro');
SELECT pg_temp.move_if_public('property_audit_logs', 'bdspro');
SELECT pg_temp.move_if_public('property_building_info', 'bdspro');
SELECT pg_temp.move_if_public('property_edvidence', 'bdspro');
SELECT pg_temp.move_if_public('property_external_ref', 'bdspro');
SELECT pg_temp.move_if_public('property_histories', 'bdspro');
SELECT pg_temp.move_if_public('property_identify', 'bdspro');
SELECT pg_temp.move_if_public('property_impact_versions', 'bdspro');
SELECT pg_temp.move_if_public('property_land_info', 'bdspro');
SELECT pg_temp.move_if_public('property_lineage', 'bdspro');
SELECT pg_temp.move_if_public('property_location', 'bdspro');
SELECT pg_temp.move_if_public('property_media', 'bdspro');
SELECT pg_temp.move_if_public('property_relation', 'bdspro');
SELECT pg_temp.move_if_public('property_statistic', 'bdspro');
SELECT pg_temp.move_if_public('property_type', 'bdspro');
SELECT pg_temp.move_if_public('property_user', 'bdspro');
SELECT pg_temp.move_if_public('lineage_media', 'bdspro');
SELECT pg_temp.move_if_public('house_info', 'bdspro');
SELECT pg_temp.move_if_public('apartment', 'bdspro');
SELECT pg_temp.move_if_public('apartment_attribute', 'bdspro');
SELECT pg_temp.move_if_public('buildings', 'bdspro');
SELECT pg_temp.move_if_public('build_room', 'bdspro');
SELECT pg_temp.move_if_public('projects', 'bdspro');
SELECT pg_temp.move_if_public('project_build', 'bdspro');
SELECT pg_temp.move_if_public('project_inventory', 'bdspro');
SELECT pg_temp.move_if_public('developer', 'bdspro');
SELECT pg_temp.move_if_public('area_regions', 'bdspro');
SELECT pg_temp.move_if_public('region', 'bdspro');
SELECT pg_temp.move_if_public('regions', 'bdspro');
SELECT pg_temp.move_if_public('posts', 'bdspro');
SELECT pg_temp.move_if_public('post_media', 'bdspro');
SELECT pg_temp.move_if_public('post_organization', 'bdspro');
SELECT pg_temp.move_if_public('post_stats', 'bdspro');
SELECT pg_temp.move_if_public('post_user', 'bdspro');
SELECT pg_temp.move_if_public('count_post_by_times', 'bdspro');
SELECT pg_temp.move_if_public('appointments', 'bdspro');
SELECT pg_temp.move_if_public('appointment_participants', 'bdspro');
SELECT pg_temp.move_if_public('reports', 'bdspro');
SELECT pg_temp.move_if_public('report_proofs', 'bdspro');
SELECT pg_temp.move_if_public('report_reasons', 'bdspro');
SELECT pg_temp.move_if_public('rates', 'bdspro');
SELECT pg_temp.move_if_public('colors', 'bdspro');
SELECT pg_temp.move_if_public('cate_value_name', 'bdspro');
SELECT pg_temp.move_if_public('doc_type', 'bdspro');
SELECT pg_temp.move_if_public('doc_type_items', 'bdspro');
SELECT pg_temp.move_if_public('action', 'bdspro');
SELECT pg_temp.move_if_public('plans', 'bdspro');
SELECT pg_temp.move_if_public('pack_items', 'bdspro');
SELECT pg_temp.move_if_public('pack_price_table', 'bdspro');

-- -----------------------------------------------------------------------------
-- crm: contact / social graph
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('tb_contact', 'crm');
SELECT pg_temp.move_if_public('db_contact', 'crm');
SELECT pg_temp.move_if_public('db_group', 'crm');
SELECT pg_temp.move_if_public('tb_history', 'crm');
SELECT pg_temp.move_if_public('friend', 'crm');
SELECT pg_temp.move_if_public('follow', 'crm');
SELECT pg_temp.move_if_public('block', 'crm');
SELECT pg_temp.move_if_public('blocks', 'crm');
SELECT pg_temp.move_if_public('sharing_access', 'crm');

-- -----------------------------------------------------------------------------
-- seo
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('seo_domain', 'seo');
SELECT pg_temp.move_if_public('seo_category', 'seo');
SELECT pg_temp.move_if_public('seo_links', 'seo');
SELECT pg_temp.move_if_public('seo_relative', 'seo');
SELECT pg_temp.move_if_public('seo_refs', 'seo');
SELECT pg_temp.move_if_public('seo_count', 'seo');
SELECT pg_temp.move_if_public('seo_count_log', 'seo');
SELECT pg_temp.move_if_public('seo_generation_log', 'seo');

-- -----------------------------------------------------------------------------
-- notification
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('notification', 'notification');
SELECT pg_temp.move_if_public('notification_history', 'notification');
SELECT pg_temp.move_if_public('noti_histories', 'notification');

-- -----------------------------------------------------------------------------
-- catalog: commercial plan / subscription (QHPro catalog)
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('catalog_plans', 'catalog');
SELECT pg_temp.move_if_public('catalog_plan_versions', 'catalog');
SELECT pg_temp.move_if_public('catalog_plan_entitlements', 'catalog');
SELECT pg_temp.move_if_public('catalog_plan_operation_policies', 'catalog');
SELECT pg_temp.move_if_public('catalog_products', 'catalog');
SELECT pg_temp.move_if_public('catalog_price_items', 'catalog');
SELECT pg_temp.move_if_public('catalog_subscriptions', 'catalog');
SELECT pg_temp.move_if_public('catalog_subscription_events', 'catalog');
SELECT pg_temp.move_if_public('subscription_checkout_snapshots', 'catalog');
SELECT pg_temp.move_if_public('subscription_settlement_receipts', 'catalog');

-- -----------------------------------------------------------------------------
-- tqd: quy hoạch / layer / planning (phần đã dump vào db_v4)
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('qh_layers', 'tqd');
SELECT pg_temp.move_if_public('qh_layer_families', 'tqd');
SELECT pg_temp.move_if_public('qh_layer_land_use', 'tqd');
SELECT pg_temp.move_if_public('qh_layer_legal_docs', 'tqd');
SELECT pg_temp.move_if_public('qh_layer_legals', 'tqd');
SELECT pg_temp.move_if_public('qh_layer_legends', 'tqd');
SELECT pg_temp.move_if_public('qh_layer_resolver_configs', 'tqd');
SELECT pg_temp.move_if_public('qh_land_use', 'tqd');
SELECT pg_temp.move_if_public('qh_legends', 'tqd');
SELECT pg_temp.move_if_public('qh_legal_documents', 'tqd');
SELECT pg_temp.move_if_public('qh_parcel_direction', 'tqd');
SELECT pg_temp.move_if_public('qh_planning_documents', 'tqd');
SELECT pg_temp.move_if_public('qh_planning_document_events', 'tqd');
SELECT pg_temp.move_if_public('qh_planning_events', 'tqd');
SELECT pg_temp.move_if_public('qh_planning_land_use', 'tqd');
SELECT pg_temp.move_if_public('qh_planning_project_tags', 'tqd');

-- -----------------------------------------------------------------------------
-- tx: giao dịch / contract
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('transactions', 'tx');
SELECT pg_temp.move_if_public('tx_transaction', 'tx');
SELECT pg_temp.move_if_public('tx_transaction_action', 'tx');
SELECT pg_temp.move_if_public('tx_contract_deals', 'tx');
SELECT pg_temp.move_if_public('tx_cost_types', 'tx');
SELECT pg_temp.move_if_public('tx_deal_costs', 'tx');

-- -----------------------------------------------------------------------------
-- legacy / backup / chưa phân loại rõ
-- -----------------------------------------------------------------------------
SELECT pg_temp.move_if_public('province_v2_backup', 'legacy');
SELECT pg_temp.move_if_public('wards_v2_backup', 'legacy');

-- schema_migrations: GIỮ public

COMMIT;

-- Kiểm tra còn sót trong public (trừ schema_migrations):
-- SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY 1;

-- Gợi ý search_path cho role app (ví dụ):
-- ALTER ROLE hieubap IN DATABASE db_v4
--   SET search_path = auth, "user", org, hub, bdspro, crm, seo, notification, catalog, tqd, tx, public;
