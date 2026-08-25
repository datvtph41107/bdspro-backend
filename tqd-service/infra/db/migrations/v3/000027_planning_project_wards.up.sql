BEGIN;

-- Create join table between planning projects and wards
CREATE TABLE IF NOT EXISTS qh_planning_project_wards (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL
        REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    ward_id VARCHAR(50) NOT NULL
        REFERENCES ward_v2(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_qh_planning_project_wards
        UNIQUE (planning_project_id, ward_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_planning_project_wards_project
    ON qh_planning_project_wards(planning_project_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_project_wards_ward
    ON qh_planning_project_wards(ward_id)
    WHERE deleted_at IS NULL;

COMMIT;

