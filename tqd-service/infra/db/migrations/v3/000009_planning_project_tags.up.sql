BEGIN;

CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100),
    type INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_tags_type
    ON tags(type)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tags_code
    ON tags(LOWER(code))
    WHERE deleted_at IS NULL AND code IS NOT NULL AND BTRIM(code) <> '';

CREATE TABLE IF NOT EXISTS qh_planning_project_tags (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL
        REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL
        REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_qh_planning_project_tags
        UNIQUE (planning_project_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_planning_project_tags_project
    ON qh_planning_project_tags(planning_project_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_project_tags_tag
    ON qh_planning_project_tags(tag_id)
    WHERE deleted_at IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_tags_type'
          AND conrelid = 'tags'::regclass
    ) THEN
        ALTER TABLE tags
            ADD CONSTRAINT chk_tags_type
            CHECK (type IN (10));
    END IF;
END $$;

COMMIT;
