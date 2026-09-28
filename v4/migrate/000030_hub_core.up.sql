-- Hub domain tables for modules/hub (local v4 DSN).
-- Idempotent: CREATE IF NOT EXISTS. Indexes are non-unique so existing
-- db_v4 hub dumps with duplicate slug/code rows can still migrate cleanly.
CREATE TABLE IF NOT EXISTS provinces (
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
CREATE INDEX IF NOT EXISTS idx_provinces_slug ON provinces (slug);
CREATE INDEX IF NOT EXISTS idx_provinces_deleted_at ON provinces (deleted_at);

CREATE TABLE IF NOT EXISTS districts (
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
CREATE INDEX IF NOT EXISTS idx_districts_province_id ON districts (province_id);
CREATE INDEX IF NOT EXISTS idx_districts_deleted_at ON districts (deleted_at);

CREATE TABLE IF NOT EXISTS wards (
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
CREATE INDEX IF NOT EXISTS idx_wards_district_id ON wards (district_id);
CREATE INDEX IF NOT EXISTS idx_wards_deleted_at ON wards (deleted_at);

CREATE TABLE IF NOT EXISTS province_v2 (
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
CREATE INDEX IF NOT EXISTS idx_province_v2_code ON province_v2 (code);
CREATE INDEX IF NOT EXISTS idx_province_v2_deleted_at ON province_v2 (deleted_at);

CREATE TABLE IF NOT EXISTS ward_v2 (
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
    updated_by BIGINT
);
CREATE INDEX IF NOT EXISTS idx_ward_v2_code ON ward_v2 (code);
CREATE INDEX IF NOT EXISTS idx_ward_v2_province_code ON ward_v2 (province_code);
CREATE INDEX IF NOT EXISTS idx_ward_v2_province_id ON ward_v2 (province_id);
CREATE INDEX IF NOT EXISTS idx_ward_v2_deleted_at ON ward_v2 (deleted_at);

CREATE TABLE IF NOT EXISTS tb_user_guide (
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
CREATE INDEX IF NOT EXISTS idx_tb_user_guide_key ON tb_user_guide (key);
CREATE INDEX IF NOT EXISTS idx_tb_user_guide_group_key ON tb_user_guide (group_key);
CREATE INDEX IF NOT EXISTS idx_tb_user_guide_deleted_at ON tb_user_guide (deleted_at);

CREATE TABLE IF NOT EXISTS faqs (
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
CREATE INDEX IF NOT EXISTS idx_faqs_group_key ON faqs (group_key);
CREATE INDEX IF NOT EXISTS idx_faqs_deleted_at ON faqs (deleted_at);

CREATE TABLE IF NOT EXISTS versions (
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
CREATE INDEX IF NOT EXISTS idx_versions_app_platform ON versions (app_name, platform);
CREATE INDEX IF NOT EXISTS idx_versions_deleted_at ON versions (deleted_at);

CREATE TABLE IF NOT EXISTS api_keys (
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
    updated_by BIGINT
);
CREATE INDEX IF NOT EXISTS idx_api_keys_name ON api_keys (name);
CREATE INDEX IF NOT EXISTS idx_api_keys_value ON api_keys (api_key);
CREATE INDEX IF NOT EXISTS idx_api_keys_app_name ON api_keys (app_name);
CREATE INDEX IF NOT EXISTS idx_api_keys_expired_at ON api_keys (expired_at);
CREATE INDEX IF NOT EXISTS idx_api_keys_deleted_at ON api_keys (deleted_at);

CREATE TABLE IF NOT EXISTS user_guide_steps (
    id BIGSERIAL PRIMARY KEY,
    user_guide_id BIGINT NOT NULL,
    step_order INTEGER NOT NULL,
    image TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
CREATE INDEX IF NOT EXISTS idx_user_guide_steps_user_guide_id ON user_guide_steps (user_guide_id);
CREATE INDEX IF NOT EXISTS idx_user_guide_steps_deleted_at ON user_guide_steps (deleted_at);
