BEGIN;
ALTER TABLE qh_planning_documents DROP CONSTRAINT IF EXISTS chk_qh_planning_document_type;
ALTER TABLE qh_planning_projects
    DROP CONSTRAINT IF EXISTS chk_qh_planning_project_type,
    DROP CONSTRAINT IF EXISTS chk_qh_planning_project_level;
-- Do not convert bigint/int identities back to the incompatible legacy types.
COMMIT;
