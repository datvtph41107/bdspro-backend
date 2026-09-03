BEGIN;

-- Additive fields required by the planning-project client read model. Some of
-- these columns already exist in pilot databases through AutoMigrate or the
-- public projection foundation, therefore every change is idempotent.
ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS decision_number VARCHAR(255),
    ADD COLUMN IF NOT EXISTS research_scope TEXT,
    ADD COLUMN IF NOT EXISTS indicators JSONB NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE qh_planning_documents
    ADD COLUMN IF NOT EXISTS status INTEGER NOT NULL DEFAULT 10;

-- Layer linkage and render metadata required by the project-related-layer task.
-- Existing environments may already have these columns through earlier v2
-- migrations or AutoMigrate, therefore this remains additive and idempotent.
ALTER TABLE qh_layers
    ADD COLUMN IF NOT EXISTS planning_project_id BIGINT NULL REFERENCES qh_planning_projects(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS display_order INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS default_visible BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS source_code VARCHAR(100),
    ADD COLUMN IF NOT EXISTS source_type VARCHAR(20) NOT NULL DEFAULT 'raster',
    ADD COLUMN IF NOT EXISTS layer_url TEXT,
    ADD COLUMN IF NOT EXISTS style_config JSONB NOT NULL DEFAULT '[]'::jsonb;

UPDATE qh_layers
SET source_code = COALESCE(NULLIF(BTRIM(source_code), ''), NULLIF(BTRIM(name), ''), id::text)
WHERE source_code IS NULL OR BTRIM(source_code) = '';

CREATE INDEX IF NOT EXISTS idx_qh_layers_planning_project
    ON qh_layers(planning_project_id, default_visible DESC, display_order, id)
    WHERE deleted_at IS NULL AND planning_project_id IS NOT NULL;

UPDATE qh_planning_documents
SET status = 10
WHERE status IS NULL OR status NOT IN (10, 20);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_qh_planning_documents_status'
          AND conrelid = 'qh_planning_documents'::regclass
    ) THEN
        ALTER TABLE qh_planning_documents
            ADD CONSTRAINT chk_qh_planning_documents_status
            CHECK (status IN (10, 20));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_decision_number
    ON qh_planning_projects(LOWER(decision_number))
    WHERE deleted_at IS NULL AND decision_number IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_updated
    ON qh_planning_projects(updated_at DESC, id DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_documents_project_status
    ON qh_planning_documents(planning_project_id, status, issue_date DESC, id DESC)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS user_followed_planning_projects (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    planning_project_id BIGINT NOT NULL
        REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_user_followed_planning_projects
        UNIQUE (user_id, planning_project_id)
);

CREATE INDEX IF NOT EXISTS idx_user_followed_planning_projects_user
    ON user_followed_planning_projects(user_id, updated_at DESC, id DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_followed_planning_projects_project
    ON user_followed_planning_projects(planning_project_id, user_id)
    WHERE deleted_at IS NULL;

CREATE OR REPLACE FUNCTION update_user_followed_planning_projects_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_user_followed_planning_projects_updated_at
    ON user_followed_planning_projects;
CREATE TRIGGER trg_user_followed_planning_projects_updated_at
BEFORE UPDATE ON user_followed_planning_projects
FOR EACH ROW
EXECUTE FUNCTION update_user_followed_planning_projects_updated_at();

COMMENT ON COLUMN qh_planning_documents.status IS
    'Planning document display status: 10=normal, 20=important.';
COMMENT ON TABLE user_followed_planning_projects IS
    'User-owned planning-project follows. TQD owns project identity and follow state.';

COMMIT;
