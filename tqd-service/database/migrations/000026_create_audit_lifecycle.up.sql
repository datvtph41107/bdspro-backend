-- Migration 000019: Create audit and lifecycle tracking tables
-- Mô tả: Bảng ghi nhận audit log mọi thao tác CRUD và lịch sử chuyển trạng thái vòng đời layer

-- 1. Bảng qh_audit_entries — ghi nhận mọi thay đổi dữ liệu (audit trail)
CREATE TABLE IF NOT EXISTS qh_audit_entries (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(64) NOT NULL,            -- loại entity: layer, parcel, label, region
    entity_id BIGINT NOT NULL,                   -- ID của entity bị thay đổi
    action VARCHAR(32) NOT NULL,                 -- loại hành động: created, updated, deleted, replaced, state_changed
    actor VARCHAR(128),                          -- user ID thực hiện
    actor_name VARCHAR(256),                     -- tên hiển thị của user
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),-- thời điểm thực hiện
    changes JSONB,                               -- diff các field thay đổi [{field, oldValue, newValue}]
    reason TEXT,                                 -- lý do thay đổi
    request_id VARCHAR(64),                      -- request ID để trace
    ip_address VARCHAR(45),                      -- IP của client
    metadata JSONB,                              -- metadata bổ sung
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_entity ON qh_audit_entries(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON qh_audit_entries(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_audit_actor ON qh_audit_entries(actor);
CREATE INDEX IF NOT EXISTS idx_audit_action ON qh_audit_entries(action);

-- 2. Bảng qh_lifecycle_transitions — ghi nhận lịch sử chuyển trạng thái vòng đời layer
CREATE TABLE IF NOT EXISTS qh_lifecycle_transitions (
    id BIGSERIAL PRIMARY KEY,
    layer_id BIGINT NOT NULL,                    -- layer được chuyển trạng thái
    from_state VARCHAR(32) NOT NULL,             -- trạng thái trước: draft, approved, effective, expired, replaced, suspended
    to_state VARCHAR(32) NOT NULL,               -- trạng thái sau
    actor VARCHAR(128),                          -- user thực hiện
    reason TEXT,                                 -- lý do chuyển trạng thái
    legal_doc VARCHAR(512),                      -- văn bản pháp lý liên quan
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lifecycle_layer ON qh_lifecycle_transitions(layer_id);
CREATE INDEX IF NOT EXISTS idx_lifecycle_created ON qh_lifecycle_transitions(created_at DESC);
