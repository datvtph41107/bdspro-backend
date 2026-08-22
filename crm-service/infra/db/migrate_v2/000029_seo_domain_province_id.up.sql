BEGIN;

ALTER TABLE seo_domain
    ADD COLUMN IF NOT EXISTS province_id BIGINT NULL;

CREATE INDEX IF NOT EXISTS idx_seo_domain_province_id
    ON seo_domain(province_id)
    WHERE deleted_at IS NULL AND province_id IS NOT NULL;

COMMENT ON COLUMN seo_domain.province_id IS
    'Province (province_v2.id) denormalized from source: adm unit / planning project / document / layer.';

COMMIT;
