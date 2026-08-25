-- 000015_layer_plan_doc_relation.up.sql
-- Quan hệ layer ↔ tài liệu quy hoạch (có thứ tự sắp xếp).

CREATE TABLE IF NOT EXISTS layer_plan_doc_relation (
    id                   BIGSERIAL PRIMARY KEY,
    layer_id             BIGINT NOT NULL
        REFERENCES qh_layers(id) ON DELETE CASCADE,
    planning_document_id BIGINT NOT NULL
        REFERENCES qh_planning_documents(id) ON DELETE CASCADE,
    sort_order           INT NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ NULL,
    created_by           BIGINT,
    updated_by           BIGINT,
    CONSTRAINT uq_layer_plan_doc_relation
        UNIQUE (layer_id, planning_document_id)
);

CREATE INDEX IF NOT EXISTS idx_layer_plan_doc_relation_layer_order
    ON layer_plan_doc_relation (layer_id, sort_order ASC, id ASC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_layer_plan_doc_relation_document
    ON layer_plan_doc_relation (planning_document_id)
    WHERE deleted_at IS NULL;
