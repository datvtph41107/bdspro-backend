BEGIN;

CREATE TABLE IF NOT EXISTS qh_planning_document_events (
    id BIGSERIAL PRIMARY KEY,
    planning_document_id BIGINT NOT NULL REFERENCES qh_planning_documents(id) ON DELETE CASCADE,
    event_type VARCHAR(80),
    event_name VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    created_by BIGINT,
    updated_by BIGINT
);

CREATE INDEX IF NOT EXISTS idx_qh_planning_document_events_doc_timeline
    ON qh_planning_document_events(planning_document_id, event_date DESC, id DESC)
    WHERE deleted_at IS NULL AND planning_document_id IS NOT NULL;

COMMIT;
