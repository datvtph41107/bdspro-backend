-- Cơ quan ban hành (master) cho layer quy hoạch
CREATE TABLE IF NOT EXISTS qh_authority_issuring (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_qh_authority_issuring_code_active
    ON qh_authority_issuring (code)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_qh_authority_issuring_name_active
    ON qh_authority_issuring (name)
    WHERE deleted_at IS NULL;

COMMENT ON TABLE qh_authority_issuring IS 'Danh mục cơ quan ban hành văn bản quy hoạch';
