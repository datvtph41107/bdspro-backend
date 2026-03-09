-- ─────────────────────────────────────────────────────────────
-- Table: applink
-- Deep link resolver table
-- code → ref_id + action
-- ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS applink (
    -- Primary key
    id BIGSERIAL PRIMARY KEY,

    -- Random numeric code [10^6, 10^10]
    -- Client chỉ thấy base64(code)
    code BIGINT NOT NULL,

    -- ID của entity tương ứng
    -- (user_id, property_id, room_id, ...)
    ref_id BIGINT NOT NULL,

    -- ApplinkAction enum (XXYY convention)
    --   01xx = user service
    --   02xx = bdspro service
    --   08xx = chat service
    action INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Unique index cho code (query chính)
CREATE UNIQUE INDEX IF NOT EXISTS idx_applink_code
ON applink (code);

-- Query theo entity
CREATE INDEX IF NOT EXISTS idx_applink_ref_action
ON applink (ref_id, action);