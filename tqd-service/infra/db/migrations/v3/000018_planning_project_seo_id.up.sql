ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS seo_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_seo_id
    ON qh_planning_projects(seo_id)
    WHERE seo_id IS NOT NULL AND seo_id > 0;
