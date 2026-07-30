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
    ver_no VARCHAR(50),
    legal_status INT4 NOT NULL DEFAULT 600,
    doc_ids BIGINT[] NOT NULL DEFAULT '{}'::bigint[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    created_by BIGINT,
    updated_by BIGINT
);

-- Existing pilot databases may already contain this table with an older GORM
-- layout. CREATE TABLE IF NOT EXISTS does not add missing columns, therefore
-- every column used by the current repositories must also be added explicitly.
ALTER TABLE qh_planning_events
    ADD COLUMN IF NOT EXISTS planning_project_id BIGINT,
    ADD COLUMN IF NOT EXISTS event_type VARCHAR(80),
    ADD COLUMN IF NOT EXISTS event_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS event_date TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS document_id BIGINT NULL,
    ADD COLUMN IF NOT EXISTS source_url TEXT,
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS ver_no VARCHAR(50),
    ADD COLUMN IF NOT EXISTS legal_status INT4 NOT NULL DEFAULT 600,
    ADD COLUMN IF NOT EXISTS doc_ids BIGINT[] NOT NULL DEFAULT '{}'::bigint[],
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL,
    ADD COLUMN IF NOT EXISTS created_by BIGINT,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT;

-- Backfill known legacy project identity column names without assuming that
-- every environment used the same historic schema.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema='public'
          AND table_name='qh_planning_events'
          AND column_name='project_id'
    ) THEN
        EXECUTE '
            UPDATE qh_planning_events
            SET planning_project_id = project_id
            WHERE planning_project_id IS NULL
              AND project_id IS NOT NULL
        ';
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema='public'
          AND table_name='qh_planning_events'
          AND column_name='qh_planning_project_id'
    ) THEN
        EXECUTE '
            UPDATE qh_planning_events
            SET planning_project_id = qh_planning_project_id
            WHERE planning_project_id IS NULL
              AND qh_planning_project_id IS NOT NULL
        ';
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_qh_planning_events_project_timeline
    ON qh_planning_events(planning_project_id, event_date DESC, id DESC)
    WHERE deleted_at IS NULL AND planning_project_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS qh_planning_relations (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    related_planning_project_id BIGINT NOT NULL REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL DEFAULT 'related',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    from_id BIGINT,
    to_id BIGINT,
    is_peer BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT chk_qh_planning_relation_not_self
        CHECK (planning_project_id <> related_planning_project_id)
);

-- Support both the current public-projection columns and the legacy
-- QHPlanningRelation repository columns.
ALTER TABLE qh_planning_relations
    ADD COLUMN IF NOT EXISTS planning_project_id BIGINT,
    ADD COLUMN IF NOT EXISTS related_planning_project_id BIGINT,
    ADD COLUMN IF NOT EXISTS relation_type VARCHAR(50) NOT NULL DEFAULT 'related',
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS from_id BIGINT,
    ADD COLUMN IF NOT EXISTS to_id BIGINT,
    ADD COLUMN IF NOT EXISTS is_peer BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL,
    ADD COLUMN IF NOT EXISTS created_by BIGINT,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT;

UPDATE qh_planning_relations
SET planning_project_id = COALESCE(planning_project_id, NULLIF(from_id, 0)),
    related_planning_project_id = COALESCE(related_planning_project_id, NULLIF(to_id, 0)),
    from_id = COALESCE(NULLIF(from_id, 0), planning_project_id),
    to_id = COALESCE(NULLIF(to_id, 0), related_planning_project_id),
    relation_type = CASE
        WHEN NULLIF(BTRIM(relation_type), '') IS NOT NULL THEN relation_type
        WHEN COALESCE(is_peer, FALSE) THEN 'peer'
        ELSE 'related'
    END;

CREATE UNIQUE INDEX IF NOT EXISTS uq_qh_planning_relations_active
    ON qh_planning_relations(
        planning_project_id,
        related_planning_project_id,
        relation_type
    )
    WHERE deleted_at IS NULL
      AND planning_project_id IS NOT NULL
      AND related_planning_project_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_relations_related
    ON qh_planning_relations(
        related_planning_project_id,
        planning_project_id
    )
    WHERE deleted_at IS NULL
      AND planning_project_id IS NOT NULL
      AND related_planning_project_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_relations_legacy_from_to
    ON qh_planning_relations(from_id, to_id)
    WHERE deleted_at IS NULL
      AND from_id IS NOT NULL
      AND to_id IS NOT NULL;

COMMIT;