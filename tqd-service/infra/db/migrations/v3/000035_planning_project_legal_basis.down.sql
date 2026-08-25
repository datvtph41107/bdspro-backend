DROP INDEX IF EXISTS idx_qh_planning_projects_legal_basis_id;

ALTER TABLE qh_planning_projects
    DROP COLUMN IF EXISTS legal_basis_id;
