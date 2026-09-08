-- 000020_create_qh_planning_tables.up.sql

CREATE TABLE IF NOT EXISTS qh_planning_projects (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    planning_type INT4 NOT NULL,
    planning_level INT4 NOT NULL,
    authority VARCHAR(255),
    jurisdiction_id BIGINT,
    summary TEXT,
    approval_date DATE,
    effective_date DATE,
    expiry_date DATE,
    validity_status VARCHAR(50) NOT NULL,
    current_version VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT,
    updated_by BIGINT
);

CREATE TABLE IF NOT EXISTS qh_planning_documents (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL REFERENCES qh_planning_projects(id),
    document_type INT4 NOT NULL,
    code VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    filepath VARCHAR(255),
    version_no VARCHAR(50),
    validity_status VARCHAR(50),
    issue_date DATE,
    effective_date DATE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT,
    updated_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_code ON qh_planning_projects(code);
CREATE INDEX IF NOT EXISTS idx_qh_planning_documents_project_id ON qh_planning_documents(planning_project_id);
