BEGIN;

DROP INDEX IF EXISTS idx_qh_planning_relations_related;
DROP INDEX IF EXISTS uq_qh_planning_relations_active;
DROP INDEX IF EXISTS idx_qh_planning_events_project_timeline;
DROP INDEX IF EXISTS idx_qh_planning_projects_public_lookup;
DROP INDEX IF EXISTS uq_qh_planning_projects_slug_active;

-- qh_planning_events and qh_planning_relations may predate this migration in
-- environments previously maintained by GORM. A safe rollback therefore keeps
-- the tables and data, removing only additive columns owned by this migration.
ALTER TABLE IF EXISTS qh_planning_relations
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS relation_type;

ALTER TABLE IF EXISTS qh_planning_events
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS source_url,
    DROP COLUMN IF EXISTS document_id,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS event_type;

ALTER TABLE IF EXISTS qh_planning_documents
    DROP COLUMN IF EXISTS thumbnail;

ALTER TABLE IF EXISTS qh_planning_projects
    DROP COLUMN IF EXISTS indicators,
    DROP COLUMN IF EXISTS research_scope,
    DROP COLUMN IF EXISTS legal_status,
    DROP COLUMN IF EXISTS total_area,
    DROP COLUMN IF EXISTS slug;

COMMIT;
