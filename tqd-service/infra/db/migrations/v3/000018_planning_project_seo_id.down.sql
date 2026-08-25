DROP INDEX IF EXISTS idx_qh_planning_projects_seo_id;

ALTER TABLE qh_planning_projects
    DROP COLUMN IF EXISTS seo_id;
