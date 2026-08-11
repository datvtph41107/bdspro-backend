BEGIN;

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- One real PostGIS parcel is enough to exercise SourceStore instead of trusting
-- location/spatial values supplied by the client.
CREATE TABLE province_v2 (
    code VARCHAR(20) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE parcels (
    id BIGSERIAL PRIMARY KEY,
    lat DOUBLE PRECISION NOT NULL DEFAULT 0,
    lon DOUBLE PRECISION NOT NULL DEFAULT 0,
    total_area_sqm DOUBLE PRECISION NOT NULL DEFAULT 0,
    address_text TEXT NOT NULL DEFAULT '',
    map_number VARCHAR(100) NOT NULL DEFAULT '',
    land_number VARCHAR(100) NOT NULL DEFAULT '',
    geometry geometry(MultiPolygon, 4326) NOT NULL,
    province_code VARCHAR(20) NOT NULL DEFAULT '',
    ward_code VARCHAR(20) NOT NULL DEFAULT '',
    deleted_at TIMESTAMPTZ
);

INSERT INTO province_v2(code, full_name)
VALUES ('79', 'Ho Chi Minh City');

INSERT INTO parcels(
    lat,
    lon,
    total_area_sqm,
    address_text,
    map_number,
    land_number,
    geometry,
    province_code,
    ward_code
) VALUES (
    10.7705,
    106.7005,
    12100,
    'QHPro local parcel',
    'LAB-01',
    '42',
    ST_GeomFromText(
        'MULTIPOLYGON(((106.7 10.77,106.701 10.77,106.701 10.771,106.7 10.771,106.7 10.77)))',
        4326
    ),
    '79',
    'LAB-WARD'
);

-- Minimal pre-000009 table required by the generated-report target.
-- This is a local lab fixture, not a replacement for the production migration chain.
CREATE TABLE user_reported (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    report_type INTEGER NOT NULL,
    status INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    province VARCHAR(100) NOT NULL DEFAULT '',
    province_code VARCHAR(20) NOT NULL DEFAULT '',
    ward_code VARCHAR(20) NOT NULL DEFAULT '',
    parcel_id BIGINT,
    region_id BIGINT,
    min_lon DOUBLE PRECISION NOT NULL DEFAULT 0,
    min_lat DOUBLE PRECISION NOT NULL DEFAULT 0,
    max_lon DOUBLE PRECISION NOT NULL DEFAULT 0,
    max_lat DOUBLE PRECISION NOT NULL DEFAULT 0,
    center_lat DOUBLE PRECISION NOT NULL DEFAULT 0,
    center_lon DOUBLE PRECISION NOT NULL DEFAULT 0,
    thumbnail_url TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    pdf_url TEXT NOT NULL DEFAULT '',
    share_url TEXT NOT NULL DEFAULT '',
    file_size BIGINT NOT NULL DEFAULT 0,
    format VARCHAR(20) NOT NULL DEFAULT '',
    comparison JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    client_request_id VARCHAR(100) NOT NULL DEFAULT '',
    job_id VARCHAR(100) NOT NULL DEFAULT '',
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_generated_reports_user_id ON user_reported(user_id);
CREATE INDEX idx_generated_reports_status ON user_reported(status);
CREATE INDEX idx_generated_reports_report_type ON user_reported(report_type);
CREATE INDEX idx_generated_reports_created_at ON user_reported(created_at);
CREATE INDEX idx_generated_reports_job_id ON user_reported(job_id);

COMMIT;
