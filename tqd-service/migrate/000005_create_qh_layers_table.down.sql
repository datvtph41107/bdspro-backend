-- Drop trigger
DROP TRIGGER IF EXISTS trg_qh_layers_updated_at ON qh_layers;

-- Drop function
DROP FUNCTION IF EXISTS update_qh_layers_updated_at;

-- Drop indexes
-- DROP INDEX IF EXISTS idx_qh_layers_name;
DROP INDEX IF EXISTS idx_qh_layers_client;
DROP INDEX IF EXISTS idx_qh_layers_admin_status;
DROP INDEX IF EXISTS idx_qh_layers_dates;
DROP INDEX IF EXISTS idx_qh_layers_search;

-- Drop table
DROP TABLE IF EXISTS qh_layers;