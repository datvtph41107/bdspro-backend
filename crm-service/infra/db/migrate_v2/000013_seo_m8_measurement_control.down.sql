DROP INDEX IF EXISTS uq_seo_measurement_import_run_active_checksum;
DROP INDEX IF EXISTS idx_seo_crawl_source_run;
DROP INDEX IF EXISTS idx_seo_ranking_source_run;
DROP INDEX IF EXISTS idx_seo_traffic_source_run;
DROP INDEX IF EXISTS idx_seo_inspection_source_run;
DROP INDEX IF EXISTS idx_seo_spd_source_run;

ALTER TABLE IF EXISTS seo_crawl_event DROP COLUMN IF EXISTS data_source;
ALTER TABLE IF EXISTS seo_crawl_event DROP COLUMN IF EXISTS source_run_id;
ALTER TABLE IF EXISTS seo_ranking_snapshot DROP COLUMN IF EXISTS source_run_id;
ALTER TABLE IF EXISTS seo_traffic_daily DROP COLUMN IF EXISTS source_run_id;
ALTER TABLE IF EXISTS seo_url_inspection_snapshot DROP COLUMN IF EXISTS source_run_id;
ALTER TABLE IF EXISTS seo_search_performance_daily DROP COLUMN IF EXISTS source_run_id;

DROP TABLE IF EXISTS seo_page_score_snapshot;
DROP TABLE IF EXISTS seo_measurement_source_state;
DROP TABLE IF EXISTS seo_measurement_import_run;
