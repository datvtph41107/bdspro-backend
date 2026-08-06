BEGIN;

DROP TABLE IF EXISTS qh_planning_project_provinces;

ALTER TABLE province_v2 DROP CONSTRAINT IF EXISTS chk_province_v2_zoom;
DROP INDEX IF EXISTS idx_provinces_qh_order;
ALTER TABLE province_v2
    DROP COLUMN IF EXISTS qh_order,
    DROP COLUMN IF EXISTS zoom;

COMMIT;
