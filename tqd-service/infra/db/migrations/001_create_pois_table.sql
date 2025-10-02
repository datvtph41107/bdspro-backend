-- Create POIs table
CREATE TABLE IF NOT EXISTS pois (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    category_id BIGINT REFERENCES poi_categories(id),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    lat DECIMAL(10,8),
    lng DECIMAL(11,8),
    geo_json TEXT,
    phone VARCHAR(20),
    website VARCHAR(255),
    rating_point DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    source VARCHAR(100),
    confidence DECIMAL(3,2) DEFAULT 0,
    last_verified_at TIMESTAMP,
    next_review_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by BIGINT,
    updated_by BIGINT
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_pois_code ON pois(code);
CREATE INDEX IF NOT EXISTS idx_pois_category_id ON pois(category_id);
CREATE INDEX IF NOT EXISTS idx_pois_name ON pois(name);
CREATE INDEX IF NOT EXISTS idx_pois_location ON pois(lat, lng);
CREATE INDEX IF NOT EXISTS idx_pois_is_active ON pois(is_active);
CREATE INDEX IF NOT EXISTS idx_pois_deleted_at ON pois(deleted_at);
CREATE INDEX IF NOT EXISTS idx_pois_rating ON pois(rating_point DESC, rating_count DESC);

-- Create spatial index for geographic queries (if using PostGIS)
-- CREATE INDEX IF NOT EXISTS idx_pois_geom ON pois USING GIST(ST_Point(lng, lat));
