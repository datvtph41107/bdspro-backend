BEGIN;

ALTER TABLE province_v2
    ADD COLUMN IF NOT EXISTS qh_order INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS zoom INTEGER NOT NULL DEFAULT 10;

CREATE INDEX IF NOT EXISTS idx_provinces_qh_order
    ON province_v2(qh_order DESC, full_name ASC);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_province_v2_zoom'
          AND conrelid = 'province_v2'::regclass
    ) THEN
        ALTER TABLE province_v2
            ADD CONSTRAINT chk_province_v2_zoom
            CHECK (zoom >= 5 AND zoom <= 22);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS qh_planning_project_provinces (
    id BIGSERIAL PRIMARY KEY,
    planning_project_id BIGINT NOT NULL
        REFERENCES qh_planning_projects(id) ON DELETE CASCADE,
    province_id VARCHAR(50) NOT NULL
        REFERENCES province_v2(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_qh_planning_project_provinces
        UNIQUE (planning_project_id, province_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_planning_project_provinces_project
    ON qh_planning_project_provinces(planning_project_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_planning_project_provinces_province
    ON qh_planning_project_provinces(province_id)
    WHERE deleted_at IS NULL;

COMMIT;
