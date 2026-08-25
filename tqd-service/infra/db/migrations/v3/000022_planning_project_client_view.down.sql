COMMENT ON COLUMN qh_planning_projects.client_view IS NULL;

ALTER TABLE qh_planning_projects
    DROP COLUMN IF EXISTS client_view;
