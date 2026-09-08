-- =====================================================
-- DOWN MIGRATION: DROP QH_REGIONS
-- =====================================================

-- =====================================================
-- DROP TRIGGERS
-- =====================================================

DROP TRIGGER IF EXISTS trg_qh_regions_updated_at ON qh_regions;
DROP TRIGGER IF EXISTS trg_qh_regions_calculate_metrics ON qh_regions;

-- =====================================================
-- DROP FUNCTIONS (trigger functions + import function)
-- =====================================================

DROP FUNCTION IF EXISTS update_qh_regions_updated_at();
DROP FUNCTION IF EXISTS calculate_qh_region_metrics();
DROP FUNCTION IF EXISTS import_region_from_geojson(jsonb, bigint);

-- =====================================================
-- DROP INDEXES
-- =====================================================

DROP INDEX IF EXISTS idx_qh_regions_geometry;
DROP INDEX IF EXISTS idx_qh_regions_layer_id;
DROP INDEX IF EXISTS idx_qh_regions_land_use;
DROP INDEX IF EXISTS idx_qh_regions_client;

-- =====================================================
-- DROP TABLE
-- =====================================================

DROP TABLE IF EXISTS qh_regions CASCADE;