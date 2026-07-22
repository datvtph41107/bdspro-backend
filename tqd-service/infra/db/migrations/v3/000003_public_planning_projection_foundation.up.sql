BEGIN;

-- Stable public identity and structured answer fields used by Search, map
-- quick-view and Next.js server rendering. The migration is additive so it can
-- run against databases that were previously maintained by GORM AutoMigrate.
ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS slug VARCHAR(180),
    ADD COLUMN IF NOT EXISTS total_area DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS legal_status INTEGER NOT NULL DEFAULT 800,
    ADD COLUMN IF NOT EXISTS research_scope TEXT,
    ADD COLUMN IF NOT EXISTS indicators JSONB NOT NULL DEFAULT '[]'::jsonb;

UPDATE qh_planning_projects
SET slug = COALESCE(
    NULLIF(BTRIM(slug), ''),
    NULLIF(BTRIM(metadata->>'slug'), ''),
    NULLIF(BTRIM(metadata->>'seoSlug'), ''),
    NULLIF(BTRIM(code), ''),
    id::text
)
WHERE slug IS NULL OR BTRIM(slug) = '';

CREATE UNIQUE INDEX IF NOT EXISTS uq_qh_planning_projects_slug_active
    ON qh_planning_projects(LOWER(slug))
    WHERE deleted_at IS NULL AND slug IS NOT NULL AND BTRIM(slug) <> '';

CREATE INDEX IF NOT EXISTS idx_qh_planning_projects_public_lookup
    ON qh_planning_projects(legal_status, validity_status, updated_at DESC)
    WHERE deleted_at IS NULL;

ALTER TABLE qh_planning_documents
    ADD COLUMN IF NOT EXISTS thumbnail VARCHAR(500);

CREATE TABLE IF NOT EXISTS qh_planning_events (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    event_type VARCHAR(80),
    event_name VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMPTZ,
    document_id BIGINT NULL REFERENCES qh_planning_documents(id) ON DELETE SET NULL,
    source_url TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    created_by BIGINT,
    updated_by BIGINT
);

ALTER TABLE qh_planning_events
    ADD COLUMN IF NOT EXISTS event_type VARCHAR(80),
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS document_id BIGINT NULL,
    ADD COLUMN IF NOT EXISTS source_url TEXT,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_qh_planning_events_project_timeline
    ON qh_planning_events(planning_project_id, event_date DESC, id DESC)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS qh_planning_relations (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    related_planning_project_id BIGINT NOT NULL REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL DEFAULT 'related',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT chk_qh_planning_relation_not_self
        CHECK (planning_project_id <> related_planning_project_id)
);

ALTER TABLE qh_planning_relations
    ADD COLUMN IF NOT EXISTS relation_type VARCHAR(50) NOT NULL DEFAULT 'related',
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

CREATE UNIQUE INDEX IF NOT EXISTS uq_qh_planning_relations_active
    ON qh_planning_relations(planning_project_id, related_planning_project_id, relation_type)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_relations_related
    ON qh_planning_relations(related_planning_project_id, planning_project_id)
    WHERE deleted_at IS NULL;

COMMIT;
