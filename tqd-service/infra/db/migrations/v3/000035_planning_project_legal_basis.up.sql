ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS legal_basis_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_legal_basis_id
    ON qh_planning_projects(legal_basis_id)
    WHERE legal_basis_id IS NOT NULL AND legal_basis_id > 0;

COMMENT ON COLUMN qh_planning_projects.legal_basis_id IS
    'Căn cứ pháp lý — id tài liệu quy hoạch (qh_planning_documents / plan_doc)';
