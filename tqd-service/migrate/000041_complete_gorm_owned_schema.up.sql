BEGIN;

-- This migration closes the historical gap between the ordered SQL history
-- and tables that the old runtime expected GORM AutoMigrate to create.  Every
-- statement is additive so it is safe for both a fresh canonical database and
-- an adopted database that already contains part of the legacy GORM schema.

CREATE TABLE IF NOT EXISTS amenities (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    category VARCHAR(100),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS contact_labels (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description VARCHAR(500),
    color VARCHAR(7),
    icon VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    contact_count INTEGER DEFAULT 0,
    is_system BOOLEAN DEFAULT FALSE,
    status BIGINT DEFAULT 0,
    type BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS directory_categories (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100),
    description TEXT,
    icon VARCHAR(255),
    color VARCHAR(7),
    parent_id BIGINT,
    sort_order BIGINT DEFAULT 0,
    status BIGINT DEFAULT 0,
    type BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    level BIGINT DEFAULT 1,
    path VARCHAR(255),
    directory_count BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS directory_sources (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    description TEXT,
    category_id BIGINT REFERENCES directory_categories(id) ON DELETE SET NULL,
    type TEXT NOT NULL,
    icon TEXT,
    color TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order BIGINT DEFAULT 0,
    expected_amount BIGINT DEFAULT 0,
    actual_amount BIGINT DEFAULT 0,
    frequency TEXT,
    start_date TIMESTAMPTZ,
    is_recurring BOOLEAN DEFAULT FALSE,
    payment_method TEXT,
    notes TEXT,
    tags TEXT
);

CREATE TABLE IF NOT EXISTS directory_suppliers (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    description TEXT,
    contact_person TEXT,
    phone TEXT,
    email TEXT,
    address TEXT,
    website TEXT,
    tax_code TEXT,
    business_license TEXT,
    rating DOUBLE PRECISION DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    service_count BIGINT DEFAULT 0,
    contract_count BIGINT DEFAULT 0,
    total_value BIGINT DEFAULT 0,
    join_date TIMESTAMPTZ,
    last_contact_date TIMESTAMPTZ,
    notes TEXT,
    tags TEXT,
    categories TEXT
);

CREATE TABLE IF NOT EXISTS open_hours (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255),
    poi_id BIGINT,
    day_of_week BIGINT NOT NULL,
    open_time TIME,
    close_time TIME,
    note TEXT,
    is_open BOOLEAN DEFAULT TRUE,
    is_recurring BOOLEAN DEFAULT TRUE,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    status BIGINT DEFAULT 0,
    type BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS poi_media (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    poi_id BIGINT NOT NULL,
    media_type VARCHAR(50) NOT NULL,
    url VARCHAR(500) NOT NULL,
    thumbnail VARCHAR(500),
    description TEXT,
    alt_text VARCHAR(255),
    file_size BIGINT DEFAULT 0,
    duration BIGINT DEFAULT 0,
    width BIGINT DEFAULT 0,
    height BIGINT DEFAULT 0,
    sort_order BIGINT DEFAULT 0,
    status BIGINT DEFAULT 0,
    type BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS poi_amenities (
    poi_id BIGINT NOT NULL REFERENCES pois(id) ON DELETE CASCADE,
    amenity_id BIGINT NOT NULL REFERENCES amenities(id) ON DELETE CASCADE,
    PRIMARY KEY (poi_id, amenity_id)
);

CREATE TABLE IF NOT EXISTS qh_direction (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name TEXT
);

CREATE TABLE IF NOT EXISTS qh_jurisdictions (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS qh_layer_legals (
    id BIGSERIAL PRIMARY KEY,
    layer_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    file_type VARCHAR(50),
    created_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS qh_layer_resolver_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL,
    config_value JSONB NOT NULL,
    description TEXT,
    updated_by VARCHAR(100),
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS qh_region_extends (
    id BIGSERIAL PRIMARY KEY,
    layer_id BIGINT NOT NULL,
    name VARCHAR(255),
    display_name VARCHAR(255),
    description TEXT,
    label_id BIGINT,
    geom_type VARCHAR(50),
    geometry geometry(Geometry, 4326) NOT NULL,
    legal_doc TEXT,
    planning_name VARCHAR(255),
    original_properties JSONB,
    source_file VARCHAR(255),
    import_batch_id VARCHAR(100),
    status BIGINT DEFAULT 10,
    version BIGINT DEFAULT 1,
    is_latest BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS qh_region_import_error_logs (
    id BIGSERIAL PRIMARY KEY,
    import_batch_id VARCHAR(100),
    layer_id BIGINT,
    error_message TEXT,
    regions_payload JSONB,
    batch_size BIGINT DEFAULT 0,
    retry_count BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS qh_shape (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name TEXT
);

CREATE TABLE IF NOT EXISTS qh_parcel_direction (
    qh_parcel_info_id BIGINT NOT NULL REFERENCES qh_parcel_info(id) ON DELETE CASCADE,
    qh_direction_id BIGINT NOT NULL REFERENCES qh_direction(id) ON DELETE CASCADE,
    PRIMARY KEY (qh_parcel_info_id, qh_direction_id)
);

CREATE TABLE IF NOT EXISTS search_index (
    title TEXT,
    address TEXT,
    geom geometry(Point, 4326),
    search_text TEXT,
    poi_id BIGINT DEFAULT 0,
    parcel_id BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_followed_parcels (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    parcel_id BIGINT NOT NULL,
    note TEXT DEFAULT '',
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_reported (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    report_type BIGINT NOT NULL,
    status BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) DEFAULT '',
    address TEXT DEFAULT '',
    province VARCHAR(100) DEFAULT '',
    province_code VARCHAR(20) DEFAULT '',
    ward_code VARCHAR(20) DEFAULT '',
    parcel_id BIGINT,
    region_id BIGINT,
    min_lon DOUBLE PRECISION DEFAULT 0,
    min_lat DOUBLE PRECISION DEFAULT 0,
    max_lon DOUBLE PRECISION DEFAULT 0,
    max_lat DOUBLE PRECISION DEFAULT 0,
    center_lat DOUBLE PRECISION DEFAULT 0,
    center_lon DOUBLE PRECISION DEFAULT 0,
    thumbnail_url TEXT DEFAULT '',
    image_url TEXT DEFAULT '',
    pdf_url TEXT DEFAULT '',
    share_url TEXT DEFAULT '',
    file_size BIGINT DEFAULT 0,
    format VARCHAR(20) DEFAULT '',
    comparison JSONB DEFAULT '{}'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    command_key VARCHAR(128) DEFAULT '',
    request_hash VARCHAR(64) DEFAULT '',
    job_id VARCHAR(100) DEFAULT '',
    error_message TEXT DEFAULT '',
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_view_events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    client_event_id VARCHAR(100) DEFAULT '',
    session_id VARCHAR(100) DEFAULT '',
    device_id VARCHAR(100) DEFAULT '',
    entity_type BIGINT NOT NULL,
    entity_id BIGINT NOT NULL,
    parcel_id BIGINT,
    region_id BIGINT,
    source BIGINT NOT NULL DEFAULT 99,
    source_ref VARCHAR(255) DEFAULT '',
    route_name VARCHAR(100) DEFAULT '',
    zoom DOUBLE PRECISION DEFAULT 0,
    center_lat DOUBLE PRECISION DEFAULT 0,
    center_lon DOUBLE PRECISION DEFAULT 0,
    viewport_min_lon DOUBLE PRECISION DEFAULT 0,
    viewport_min_lat DOUBLE PRECISION DEFAULT 0,
    viewport_max_lon DOUBLE PRECISION DEFAULT 0,
    viewport_max_lat DOUBLE PRECISION DEFAULT 0,
    visible_ms BIGINT NOT NULL DEFAULT 0,
    count_intent BOOLEAN NOT NULL DEFAULT FALSE,
    counted BOOLEAN NOT NULL DEFAULT FALSE,
    dedupe_key VARCHAR(255) DEFAULT '',
    metadata JSONB DEFAULT '{}'::jsonb,
    viewed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_view_history (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    entity_type BIGINT NOT NULL,
    entity_id BIGINT NOT NULL,
    parcel_id BIGINT,
    region_id BIGINT,
    source BIGINT NOT NULL DEFAULT 99,
    viewed_at TIMESTAMPTZ NOT NULL,
    first_viewed_at TIMESTAMPTZ,
    view_count BIGINT NOT NULL DEFAULT 0,
    counted_view_count BIGINT NOT NULL DEFAULT 0,
    last_zoom DOUBLE PRECISION DEFAULT 0,
    last_center_lat DOUBLE PRECISION DEFAULT 0,
    last_center_lon DOUBLE PRECISION DEFAULT 0,
    viewport_min_lon DOUBLE PRECISION DEFAULT 0,
    viewport_min_lat DOUBLE PRECISION DEFAULT 0,
    viewport_max_lon DOUBLE PRECISION DEFAULT 0,
    viewport_max_lat DOUBLE PRECISION DEFAULT 0,
    last_client_event_id VARCHAR(100) DEFAULT '',
    last_visible_ms BIGINT NOT NULL DEFAULT 0,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- Bring adopted legacy tables to the current model contract. CREATE TABLE IF
-- NOT EXISTS alone cannot add columns to schemas previously made by GORM or
-- the old loose SQL files.
ALTER TABLE parcels
    ADD COLUMN IF NOT EXISTS created_by BIGINT,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon, 4326),
    ADD COLUMN IF NOT EXISTS lat DOUBLE PRECISION DEFAULT 0,
    ADD COLUMN IF NOT EXISTS lon DOUBLE PRECISION DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ref_id BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ref_type BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS seo_id BIGINT,
    ADD COLUMN IF NOT EXISTS centroid_geom geometry(Point, 4326),
    ADD COLUMN IF NOT EXISTS is_seo BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE qh_land_use
    ADD COLUMN IF NOT EXISTS created_by BIGINT,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS note TEXT,
    ADD COLUMN IF NOT EXISTS name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS code VARCHAR(10),
    ADD COLUMN IF NOT EXISTS display_order BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS is_visible BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS priority BIGINT DEFAULT 50,
    ADD COLUMN IF NOT EXISTS build_condition TEXT,
    ADD COLUMN IF NOT EXISTS color VARCHAR(10),
    ADD COLUMN IF NOT EXISTS can_build BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS warn_level BIGINT DEFAULT 0;

ALTER TABLE pois
    ADD COLUMN IF NOT EXISTS created_by BIGINT,
    ADD COLUMN IF NOT EXISTS updated_by BIGINT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS code VARCHAR(50),
    ADD COLUMN IF NOT EXISTS name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS address TEXT,
    ADD COLUMN IF NOT EXISTS phone VARCHAR(50),
    ADD COLUMN IF NOT EXISTS email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS website VARCHAR(255),
    ADD COLUMN IF NOT EXISTS latitude NUMERIC(10,8) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS longitude NUMERIC(11,8) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS category_id BIGINT,
    ADD COLUMN IF NOT EXISTS rating NUMERIC(3,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS review_count BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS is_featured BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS cover_image VARCHAR(500),
    ADD COLUMN IF NOT EXISTS images TEXT,
    ADD COLUMN IF NOT EXISTS tags TEXT,
    ADD COLUMN IF NOT EXISTS amenities TEXT,
    ADD COLUMN IF NOT EXISTS view_count BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS like_count BIGINT DEFAULT 0;

ALTER TABLE province_v2
    ADD COLUMN IF NOT EXISTS full_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS short_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS lat NUMERIC(10,8),
    ADD COLUMN IF NOT EXISTS lng NUMERIC(11,8),
    ADD COLUMN IF NOT EXISTS code VARCHAR(20),
    ADD COLUMN IF NOT EXISTS ward_count BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE ward_v2
    ADD COLUMN IF NOT EXISTS full_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS short_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS lat NUMERIC(10,8),
    ADD COLUMN IF NOT EXISTS lng NUMERIC(11,8),
    ADD COLUMN IF NOT EXISTS code VARCHAR(20),
    ADD COLUMN IF NOT EXISTS province_id VARCHAR(50),
    ADD COLUMN IF NOT EXISTS province_code VARCHAR(20),
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE qh_parcel_info
    ADD COLUMN IF NOT EXISTS parcel_id BIGINT,
    ADD COLUMN IF NOT EXISTS property_code VARCHAR(100),
    ADD COLUMN IF NOT EXISTS property_uuid UUID,
    ADD COLUMN IF NOT EXISTS map_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS land_number VARCHAR(50),
    ADD COLUMN IF NOT EXISTS total_area_sqm DOUBLE PRECISION DEFAULT 0,
    ADD COLUMN IF NOT EXISTS lat DOUBLE PRECISION DEFAULT 0,
    ADD COLUMN IF NOT EXISTS lon DOUBLE PRECISION DEFAULT 0,
    ADD COLUMN IF NOT EXISTS province_code VARCHAR(20),
    ADD COLUMN IF NOT EXISTS district_code VARCHAR(20),
    ADD COLUMN IF NOT EXISTS ward_code VARCHAR(20),
    ADD COLUMN IF NOT EXISTS address_text TEXT,
    ADD COLUMN IF NOT EXISTS adr_search TEXT,
    ADD COLUMN IF NOT EXISTS is_seo BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS province_v1_id BIGINT,
    ADD COLUMN IF NOT EXISTS district_v1_id BIGINT,
    ADD COLUMN IF NOT EXISTS ward_v1_id BIGINT,
    ADD COLUMN IF NOT EXISTS province_v2_id BIGINT,
    ADD COLUMN IF NOT EXISTS ward_v2_id BIGINT,
    ADD COLUMN IF NOT EXISTS planning_infos JSONB,
    ADD COLUMN IF NOT EXISTS planning_use_id BIGINT,
    ADD COLUMN IF NOT EXISTS shape_id BIGINT,
    ADD COLUMN IF NOT EXISTS facade BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS land_type_id BIGINT,
    ADD COLUMN IF NOT EXISTS geometry geometry(MultiPolygon, 4326),
    ADD COLUMN IF NOT EXISTS source_type BIGINT DEFAULT 2,
    ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_amenities_deleted_at ON amenities (deleted_at);
CREATE INDEX IF NOT EXISTS idx_amenities_category ON amenities (category);
CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_labels_code ON contact_labels (code);
CREATE INDEX IF NOT EXISTS idx_contact_labels_deleted_at ON contact_labels (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_directory_categories_code ON directory_categories (code);
CREATE INDEX IF NOT EXISTS idx_directory_categories_parent_id ON directory_categories (parent_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_directory_sources_code ON directory_sources (code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_directory_suppliers_code ON directory_suppliers (code);
CREATE INDEX IF NOT EXISTS idx_open_hours_poi_id ON open_hours (poi_id);
CREATE INDEX IF NOT EXISTS idx_poi_media_poi_id ON poi_media (poi_id);
CREATE INDEX IF NOT EXISTS idx_parcels_deleted_at ON parcels (deleted_at);
CREATE INDEX IF NOT EXISTS idx_parcels_geometry ON parcels USING GIST (geometry);
CREATE INDEX IF NOT EXISTS idx_provinces_name ON province_v2 (full_name);
CREATE INDEX IF NOT EXISTS idx_provinces_short_name ON province_v2 (short_name);
CREATE INDEX IF NOT EXISTS idx_provinces_code ON province_v2 (code);
CREATE INDEX IF NOT EXISTS idx_wards_name ON ward_v2 (full_name);
CREATE INDEX IF NOT EXISTS idx_wards_short_name ON ward_v2 (short_name);
CREATE INDEX IF NOT EXISTS idx_wards_code ON ward_v2 (code);
CREATE INDEX IF NOT EXISTS idx_wards_province_id ON ward_v2 (province_id);
CREATE INDEX IF NOT EXISTS idx_qh_direction_deleted_at ON qh_direction (deleted_at);
CREATE INDEX IF NOT EXISTS idx_qh_jurisdictions_deleted_at ON qh_jurisdictions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_layer_legal ON qh_layer_legals (layer_id);
CREATE INDEX IF NOT EXISTS idx_qh_layer_legals_deleted_at ON qh_layer_legals (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_qh_layer_resolver_configs_config_key ON qh_layer_resolver_configs (config_key);
CREATE INDEX IF NOT EXISTS idx_qh_parcel_info_parcel_id ON qh_parcel_info (parcel_id);
CREATE INDEX IF NOT EXISTS idx_qh_parcel_info_is_seo ON qh_parcel_info (is_seo);
CREATE INDEX IF NOT EXISTS idx_qh_region_extends_layer_id ON qh_region_extends (layer_id);
CREATE INDEX IF NOT EXISTS idx_qh_region_extends_geometry ON qh_region_extends USING GIST (geometry);
CREATE INDEX IF NOT EXISTS idx_qh_region_import_error_logs_batch ON qh_region_import_error_logs (import_batch_id);
CREATE INDEX IF NOT EXISTS idx_qh_region_import_error_logs_status ON qh_region_import_error_logs (status);
CREATE INDEX IF NOT EXISTS idx_qh_shape_deleted_at ON qh_shape (deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_followed_parcels_user_id ON user_followed_parcels (user_id);
CREATE INDEX IF NOT EXISTS idx_user_followed_parcels_parcel_id ON user_followed_parcels (parcel_id);
CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_followed_parcels_active_user_parcel
    ON user_followed_parcels (user_id, parcel_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_generated_reports_user_id ON user_reported (user_id);
CREATE INDEX IF NOT EXISTS idx_generated_reports_report_type ON user_reported (report_type);
CREATE INDEX IF NOT EXISTS idx_generated_reports_status ON user_reported (status);
CREATE INDEX IF NOT EXISTS idx_generated_reports_job_id ON user_reported (job_id);
CREATE INDEX IF NOT EXISTS idx_generated_reports_created_at ON user_reported (created_at);
CREATE INDEX IF NOT EXISTS idx_generated_reports_expires_at ON user_reported (expires_at);
CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_report_command_key
    ON user_reported (user_id, command_key)
    WHERE command_key <> '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_view_events_user_id ON user_view_events (user_id);
CREATE INDEX IF NOT EXISTS idx_user_view_events_entity ON user_view_events (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_user_view_events_viewed_at ON user_view_events (viewed_at);
CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_view_event_client_event
    ON user_view_events (user_id, client_event_id)
    WHERE client_event_id <> '' AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_view_history_user_id ON user_view_history (user_id);
CREATE INDEX IF NOT EXISTS idx_user_view_history_entity ON user_view_history (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_user_view_history_viewed_at ON user_view_history (viewed_at);
CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_view_history_user_entity
    ON user_view_history (user_id, entity_type, entity_id)
    WHERE deleted_at IS NULL;

COMMIT;
