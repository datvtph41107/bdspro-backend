-- Historical databases received qh_labels and qh_regions.label_id from GORM
-- AutoMigrate. Canonical SQL establishes these dependencies before their first
-- use. Every statement is additive and remains compatible with adopted data.
CREATE TABLE IF NOT EXISTS qh_labels (
    id          BIGSERIAL PRIMARY KEY,
    layer_id    BIGINT       NOT NULL REFERENCES qh_layers(id),
    name        VARCHAR(255) NOT NULL,
    color       VARCHAR(20)  NOT NULL DEFAULT '#1890ff',
    description TEXT,
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

ALTER TABLE qh_regions
    ADD COLUMN IF NOT EXISTS label_id BIGINT REFERENCES qh_labels(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_qh_regions_label_id
    ON qh_regions(label_id)
    WHERE deleted_at IS NULL;

-- Denormalized region count per label (tự cập nhật qua trigger trên qh_regions)
ALTER TABLE qh_labels ADD COLUMN IF NOT EXISTS region_count BIGINT NOT NULL DEFAULT 0;

-- Backfill một lần
UPDATE qh_labels l
SET region_count = (
    SELECT COUNT(*)::bigint
    FROM qh_regions r
    WHERE r.label_id = l.id
      AND r.deleted_at IS NULL
);

CREATE OR REPLACE FUNCTION qh_label_recount_region_count(p_label_id bigint)
RETURNS void AS $$
BEGIN
    IF p_label_id IS NULL THEN
        RETURN;
    END IF;
    UPDATE qh_labels l
    SET region_count = (
        SELECT COUNT(*)::bigint
        FROM qh_regions r
        WHERE r.label_id = l.id
          AND r.deleted_at IS NULL
    )
    WHERE l.id = p_label_id
      AND l.deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_qh_regions_maintain_label_region_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM qh_label_recount_region_count(NEW.label_id);
    ELSIF TG_OP = 'UPDATE' THEN
        IF OLD.label_id IS DISTINCT FROM NEW.label_id THEN
            PERFORM qh_label_recount_region_count(OLD.label_id);
            PERFORM qh_label_recount_region_count(NEW.label_id);
        ELSIF OLD.deleted_at IS DISTINCT FROM NEW.deleted_at THEN
            PERFORM qh_label_recount_region_count(NEW.label_id);
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM qh_label_recount_region_count(OLD.label_id);
    END IF;
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_qh_regions_label_region_count ON qh_regions;
CREATE TRIGGER trg_qh_regions_label_region_count
    AFTER INSERT OR UPDATE OR DELETE ON qh_regions
    FOR EACH ROW
    EXECUTE FUNCTION trg_qh_regions_maintain_label_region_count();

COMMENT ON COLUMN qh_labels.region_count IS 'Số qh_regions (deleted_at IS NULL) có label_id = id; trigger trg_qh_regions_label_region_count';
