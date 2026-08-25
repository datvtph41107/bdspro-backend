BEGIN;

-- Create table to store neighboring wards
CREATE TABLE IF NOT EXISTS ward_neighbors (
    id BIGSERIAL PRIMARY KEY,
    ward_id BIGINT NOT NULL
        REFERENCES ward_v2(id) ON DELETE CASCADE,
    neighbor_ward_id BIGINT NOT NULL
        REFERENCES ward_v2(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_ward_neighbors UNIQUE (ward_id, neighbor_ward_id)
);

CREATE INDEX IF NOT EXISTS idx_ward_neighbors_ward
    ON ward_neighbors(ward_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_ward_neighbors_neighbor
    ON ward_neighbors(neighbor_ward_id)
    WHERE deleted_at IS NULL;

COMMIT;

