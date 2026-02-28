CREATE TABLE property (
    id BIGSERIAL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,

    property_type_id BIGINT,
    title VARCHAR(255),

    location_id BIGINT,
    address_detail VARCHAR(100),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    map_url VARCHAR(100),

    project_id BIGINT,
    block_id BIGINT,
    building_id BIGINT,
    level VARCHAR(100),

    unit_code VARCHAR(50),
    identifier VARCHAR(50),
    national_id VARCHAR(50),
    national_id_verified BOOLEAN DEFAULT FALSE,

    area_total DECIMAL(15,2),
    area_land DECIMAL(15,2),
    area_residential DECIMAL(15,2),

    legal_note TEXT,
    legal_status SMALLINT,

    record_status SMALLINT NOT NULL DEFAULT 10,
    scope SMALLINT NOT NULL DEFAULT 10,
    source_type SMALLINT NOT NULL, -- 10: system, 20: user

    avatar_id BIGINT,
    building_info_id BIGINT,

    province_id BIGINT,
    district_id BIGINT,
    ward_id BIGINT,
    region_id BIGINT,
    PRIMARY KEY (id, source_type, created_at)
) PARTITION BY LIST (source_type);

CREATE INDEX idx_property_identifier ON property (identifier, source_type);
CREATE INDEX idx_property_type_id ON property (property_type_id);
CREATE INDEX idx_property_province ON property (province_id);
CREATE INDEX idx_property_scope_created ON property (scope, created_at);
CREATE INDEX idx_property_filter ON property (source_type, scope, created_at DESC);
CREATE INDEX idx_property_not_deleted ON property (deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_property_region ON property(region_id);