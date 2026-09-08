BEGIN;

DROP TRIGGER IF EXISTS trg_user_followed_planning_projects_updated_at
    ON user_followed_planning_projects;
DROP FUNCTION IF EXISTS update_user_followed_planning_projects_updated_at();

DROP TABLE IF EXISTS user_followed_planning_projects;

DROP INDEX IF EXISTS idx_qh_planning_documents_project_status;
DROP INDEX IF EXISTS idx_qh_layers_planning_project;
DROP INDEX IF EXISTS idx_qh_planning_projects_updated;
DROP INDEX IF EXISTS idx_qh_planning_projects_decision_number;

ALTER TABLE qh_planning_documents
    DROP CONSTRAINT IF EXISTS chk_qh_planning_documents_status,
    DROP COLUMN IF EXISTS status;

-- research_scope and indicators are owned by migration
-- 000003_public_planning_projection_foundation and are intentionally retained.
ALTER TABLE qh_planning_projects
    DROP COLUMN IF EXISTS decision_number;

COMMIT;
