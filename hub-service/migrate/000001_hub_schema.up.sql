-- Hub sở hữu metadata/location/deep-link và dữ liệu hỗ trợ ứng dụng.
-- File này là fresh-schema authority; dữ liệu location/FAQ là seed/import riêng.

CREATE TABLE applink (
    id BIGSERIAL PRIMARY KEY,
    code BIGINT NOT NULL,
    ref_id BIGINT NOT NULL,
    action INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_applink_code ON applink (code);
CREATE INDEX idx_applink_ref_action ON applink (ref_id, action);

CREATE TABLE tb_event_queue (
    id BIGSERIAL PRIMARY KEY,
    event_type TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_data TEXT,
    metadata TEXT,
    status INTEGER NOT NULL DEFAULT 10,
    profile_id BIGINT,
    organization_id BIGINT,
    scheduled_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    retry_count BIGINT NOT NULL DEFAULT 0,
    max_retries BIGINT NOT NULL DEFAULT 3,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX idx_tb_event_queue_deleted_at ON tb_event_queue (deleted_at);
CREATE INDEX idx_tb_event_queue_status ON tb_event_queue (status);
CREATE INDEX idx_tb_event_queue_profile_id ON tb_event_queue (profile_id);
CREATE INDEX idx_tb_event_queue_organization_id ON tb_event_queue (organization_id);

-- Location v1 dùng external code dạng chuỗi ở parent reference. Primary key
-- vẫn là BIGINT theo model production hiện hành.
CREATE TABLE provinces (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type INTEGER NOT NULL,
    type_text VARCHAR(50) NOT NULL,
    codename VARCHAR(100),
    short_codename VARCHAR(100),
    slug VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE UNIQUE INDEX idx_provinces_slug ON provinces (slug);
CREATE INDEX idx_provinces_deleted_at ON provinces (deleted_at);

CREATE TABLE districts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    province_id VARCHAR(10) NOT NULL,
    type INTEGER NOT NULL,
    type_text VARCHAR(50) NOT NULL,
    codename VARCHAR(100),
    short_codename VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX idx_districts_province_id ON districts (province_id);
CREATE INDEX idx_districts_deleted_at ON districts (deleted_at);

CREATE TABLE wards (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    district_id VARCHAR(10) NOT NULL,
    type INTEGER NOT NULL,
    type_text VARCHAR(50) NOT NULL,
    codename VARCHAR(100),
    short_codename VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX idx_wards_district_id ON wards (district_id);
CREATE INDEX idx_wards_deleted_at ON wards (deleted_at);

CREATE TABLE province_v2 (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code INTEGER NOT NULL,
    codename VARCHAR(100),
    division_type VARCHAR(100),
    phone_code INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE UNIQUE INDEX idx_province_v2_code ON province_v2 (code);
CREATE INDEX idx_province_v2_deleted_at ON province_v2 (deleted_at);

CREATE TABLE ward_v2 (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code INTEGER NOT NULL,
    codename VARCHAR(100),
    division_type VARCHAR(100),
    short_codename VARCHAR(100),
    province_code INTEGER NOT NULL,
    province_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT fk_ward_v2_province_code FOREIGN KEY (province_code)
        REFERENCES province_v2 (code) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE UNIQUE INDEX idx_ward_v2_code ON ward_v2 (code);
CREATE INDEX idx_ward_v2_province_code ON ward_v2 (province_code);
CREATE INDEX idx_ward_v2_province_id ON ward_v2 (province_id);
CREATE INDEX idx_ward_v2_deleted_at ON ward_v2 (deleted_at);

CREATE TABLE tb_user_guide (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    group_key TEXT NOT NULL,
    key TEXT,
    mode INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE UNIQUE INDEX idx_tb_user_guide_key ON tb_user_guide (key);
CREATE INDEX idx_tb_user_guide_group_key ON tb_user_guide (group_key);
CREATE INDEX idx_tb_user_guide_deleted_at ON tb_user_guide (deleted_at);

CREATE TABLE system_config (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    group_key INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE UNIQUE INDEX idx_system_config_key ON system_config (key);
CREATE INDEX idx_system_config_deleted_at ON system_config (deleted_at);

CREATE TABLE faqs (
    id BIGSERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    group_key VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX idx_faqs_group_key ON faqs (group_key);
CREATE INDEX idx_faqs_deleted_at ON faqs (deleted_at);

CREATE TABLE versions (
    id BIGSERIAL PRIMARY KEY,
    app_name VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    version_name VARCHAR(50) NOT NULL,
    build_number BIGINT NOT NULL DEFAULT 0,
    force_update BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    download_url VARCHAR(255),
    release_notes TEXT,
    bundle_name VARCHAR(255),
    bundle_size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX idx_versions_app_platform ON versions (app_name, platform);
CREATE INDEX idx_versions_deleted_at ON versions (deleted_at);

CREATE TABLE api_keys (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    app_name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    api_key VARCHAR(255) NOT NULL,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT uq_api_keys_name UNIQUE (name),
    CONSTRAINT uq_api_keys_value UNIQUE (api_key)
);
CREATE INDEX idx_api_keys_app_name ON api_keys (app_name);
CREATE INDEX idx_api_keys_expired_at ON api_keys (expired_at);
CREATE INDEX idx_api_keys_deleted_at ON api_keys (deleted_at);

CREATE TABLE user_guide_steps (
    id BIGSERIAL PRIMARY KEY,
    user_guide_id BIGINT NOT NULL,
    step_order INTEGER NOT NULL,
    image TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    CONSTRAINT fk_user_guide_steps_guide FOREIGN KEY (user_guide_id)
        REFERENCES tb_user_guide (id) ON DELETE CASCADE
);
CREATE INDEX idx_user_guide_steps_user_guide_id ON user_guide_steps (user_guide_id);
CREATE INDEX idx_user_guide_steps_deleted_at ON user_guide_steps (deleted_at);

-- Dòng tương tác dùng tổ hợp thuộc tính nghiệp vụ, không tự thêm ID ngoài model.
CREATE TABLE interactive_events (
    start_time BIGINT NOT NULL,
    end_time BIGINT,
    event VARCHAR(30) NOT NULL,
    ref_id BIGINT,
    screen VARCHAR(40),
    duration BIGINT,
    profile_id BIGINT
);
CREATE INDEX idx_interactive_events_start_time ON interactive_events (start_time);
CREATE INDEX idx_interactive_events_end_time ON interactive_events (end_time);
CREATE INDEX idx_interactive_events_event ON interactive_events (event);
CREATE INDEX idx_interactive_events_ref_id ON interactive_events (ref_id);
CREATE INDEX idx_interactive_events_profile_id ON interactive_events (profile_id);

CREATE TABLE error_logs (
    id BIGSERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    stack TEXT,
    screen VARCHAR(200),
    user_id VARCHAR(100),
    device JSONB,
    extra JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_error_logs_user_id ON error_logs (user_id);
CREATE INDEX idx_error_logs_created_at ON error_logs (created_at);

CREATE TABLE update_data (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    resource_type INTEGER NOT NULL,
    resource_id BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
CREATE INDEX idx_update_data_owner_resource ON update_data (owner_id, resource_type);
CREATE INDEX idx_update_data_updated_at ON update_data (updated_at);
