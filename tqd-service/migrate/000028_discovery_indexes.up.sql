-- Discovery V1 indexes.
-- Source domain tables remain the source of truth; these indexes accelerate spatial identify and federated search.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- The old source expected these tables to have been created by GORM before
-- this migration ran.  The canonical SQL owner must be self-sufficient, so
-- establish the minimum current repository contract before creating indexes.
CREATE TABLE IF NOT EXISTS parcels (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    geometry geometry(MultiPolygon, 4326),
    lat DOUBLE PRECISION DEFAULT 0,
    lon DOUBLE PRECISION DEFAULT 0,
    ref_id BIGINT DEFAULT 0,
    ref_type BIGINT DEFAULT 0,
    seo_id BIGINT,
    centroid_geom geometry(Point, 4326),
    is_seo BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS qh_planning_land_use (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    planning_type_code TEXT,
    planning_type_name TEXT
);

CREATE TABLE IF NOT EXISTS qh_parcel_info (
    id BIGSERIAL PRIMARY KEY,
    parcel_id BIGINT NOT NULL,
    property_code VARCHAR(100),
    property_uuid UUID,
    map_number VARCHAR(50),
    land_number VARCHAR(50),
    total_area_sqm DOUBLE PRECISION DEFAULT 0,
    lat DOUBLE PRECISION DEFAULT 0,
    lon DOUBLE PRECISION DEFAULT 0,
    province_code VARCHAR(20),
    district_code VARCHAR(20),
    ward_code VARCHAR(20),
    address_text TEXT,
    adr_search TEXT,
    is_seo BOOLEAN DEFAULT FALSE,
    province_v1_id BIGINT,
    district_v1_id BIGINT,
    ward_v1_id BIGINT,
    province_v2_id BIGINT,
    ward_v2_id BIGINT,
    planning_infos JSONB,
    planning_use_id BIGINT REFERENCES qh_planning_land_use(id) ON DELETE SET NULL,
    shape_id BIGINT,
    facade BIGINT DEFAULT 0,
    land_type_id BIGINT,
    geometry geometry(MultiPolygon, 4326),
    source_type BIGINT DEFAULT 2,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS poi_categories (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    color VARCHAR(7),
    parent_id BIGINT,
    level BIGINT DEFAULT 1,
    path VARCHAR(500),
    sort_order BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    poi_count BIGINT DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_poi_categories_code
    ON poi_categories (code);

CREATE TABLE IF NOT EXISTS pois (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    code VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    website VARCHAR(255),
    latitude NUMERIC(10,8) DEFAULT 0,
    longitude NUMERIC(11,8) DEFAULT 0,
    category_id BIGINT NOT NULL REFERENCES poi_categories(id) ON DELETE RESTRICT,
    rating NUMERIC(3,2) DEFAULT 0,
    review_count BIGINT DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    cover_image VARCHAR(500),
    images TEXT,
    tags TEXT,
    amenities TEXT,
    view_count BIGINT DEFAULT 0,
    like_count BIGINT DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_pois_code ON pois (code);

CREATE TABLE IF NOT EXISTS province_v2 (
    id VARCHAR(50) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    short_name VARCHAR(100),
    lat NUMERIC(10,8),
    lng NUMERIC(11,8),
    code VARCHAR(20),
    ward_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_province_v2_code ON province_v2 (code);

CREATE TABLE IF NOT EXISTS ward_v2 (
    id VARCHAR(50) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    short_name VARCHAR(100),
    lat NUMERIC(10,8),
    lng NUMERIC(11,8),
    code VARCHAR(20),
    province_id VARCHAR(50) NOT NULL REFERENCES province_v2(id) ON DELETE CASCADE,
    province_code VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ward_v2_code ON ward_v2 (code);

CREATE INDEX IF NOT EXISTS idx_discovery_parcels_geometry_active
    ON parcels USING GIST (geometry)
    WHERE deleted_at IS NULL AND geometry IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_parcel_id
    ON qh_parcel_info (parcel_id);

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_property_code_lower
    ON qh_parcel_info (LOWER(property_code));

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_map_land
    ON qh_parcel_info (map_number, land_number);

CREATE INDEX IF NOT EXISTS idx_discovery_parcel_info_adr_trgm
    ON qh_parcel_info USING GIN (LOWER(adr_search) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_discovery_regions_geometry_active
    ON qh_regions USING GIST (geometry)
    WHERE deleted_at IS NULL AND status = 10 AND is_latest = TRUE;

CREATE INDEX IF NOT EXISTS idx_discovery_regions_display_name_trgm
    ON qh_regions USING GIN (LOWER(COALESCE(display_name, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_regions_name_trgm
    ON qh_regions USING GIN (LOWER(COALESCE(name, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_pois_name_trgm
    ON pois USING GIN (LOWER(name) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_pois_address_trgm
    ON pois USING GIN (LOWER(COALESCE(address, '')) gin_trgm_ops)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_discovery_ward_name_trgm
    ON ward_v2 USING GIN (LOWER(COALESCE(full_name, '')) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_discovery_province_name_trgm
    ON province_v2 USING GIN (LOWER(COALESCE(full_name, '')) gin_trgm_ops);
