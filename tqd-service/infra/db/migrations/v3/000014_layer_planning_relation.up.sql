-- 000014_layer_planning_relation.up.sql
-- Quan hệ layer ↔ đồ án quy hoạch (có thứ tự sắp xếp).
-- Flow API: parcel → regions (intersects) → layers → bảng này → planning-projects.

DROP TABLE IF EXISTS parcel_planning_relation;

CREATE TABLE IF NOT EXISTS layer_planning_relation (
    id                   BIGSERIAL PRIMARY KEY,
    layer_id             BIGINT NOT NULL
        REFERENCES qh_layers(id) ON DELETE CASCADE,
    planning_project_id  BIGINT NOT NULL
        REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    sort_order           INT NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ NULL,
    created_by           BIGINT,
    updated_by           BIGINT,
    CONSTRAINT uq_layer_planning_relation
        UNIQUE (layer_id, planning_project_id)
);

CREATE INDEX IF NOT EXISTS idx_layer_planning_relation_layer_order
    ON layer_planning_relation (layer_id, sort_order ASC, id ASC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_layer_planning_relation_project
    ON layer_planning_relation (planning_project_id)
    WHERE deleted_at IS NULL;
