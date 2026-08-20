BEGIN;

DROP INDEX IF EXISTS idx_qh_planning_document_events_doc_timeline;
DROP TABLE IF EXISTS qh_planning_document_events;

COMMIT;
