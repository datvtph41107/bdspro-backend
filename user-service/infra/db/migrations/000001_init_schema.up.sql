CREATE TABLE user_addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    -- Loại địa chỉ: home, work,...
    address_type VARCHAR(20) NOT NULL,
    full_address TEXT NOT NULL,
    -- Các thành phần địa chỉ (để dễ lọc, phân tích)
    street VARCHAR(255),
    ward VARCHAR(100),
    district VARCHAR(100),
    city VARCHAR(100) NOT NULL,
    region VARCHAR(50),               -- Bắc, Trung, Nam 
    country VARCHAR(50) DEFAULT 'Vietnam',
    postal_code VARCHAR(20),
    
    -- Tọa độ (lấy từ geo-service khi tạo/cập nhật)
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    
    -- (tìm kiếm vùng)
    geohash_5 VARCHAR(5) NOT NULL,    -- ~5km
    geohash_7 VARCHAR(7) NOT NULL,    -- ~150m
    geohash_9 VARCHAR(9),              -- ~5m 
    
    -- Thông geocoding (nguồn gốc)
    place_id VARCHAR(255),             -- ID từ Google Maps hoặc OSM
    formatted_address TEXT,            -- Địa chỉ
    location_type VARCHAR(50),         -- 'rooftop', 'range_interpolated', 'geometric_center', v.v.
    
    -- Thông tin xác thực địa chỉ
    is_verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMPTZ,
    verification_method VARCHAR(50),   -- 'register',...
    
    -- Thông tin thời gian
    moved_in_at DATE,                   -- Ngày bắt đầu vị trí
    moved_out_at DATE,                 
    is_current BOOLEAN DEFAULT false,
    
    -- Thông tin chi tiết 
    property_type VARCHAR(50),          -- 'apartment', 'house'
    living_area_sqm FLOAT,              -- Diện tích sử dụng (m²)
    number_of_residents INTEGER,        -- Số người ở
    
    -- Thông tin tương tác (cho Living Map)
    total_checkins INTEGER DEFAULT 0,   -- Tổng số lần check-in tại đây
    last_checkin_at TIMESTAMPTZ,        -- Lần check-in gần nhất
    
    -- Tags để phân loại linh hoạt (JSONB)
    tags JSONB DEFAULT '[]'::jsonb,     -- ["favorite", "investment"]
    is_default BOOLEAN DEFAULT false,   -- Địa chỉ mặc định (ưu tiên hiển thị)
    notes TEXT,                         -- Ghi chú cá nhân
    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index trên user_id để lấy địa chỉ của một user nhanh chóng
CREATE INDEX idx_user_addresses_user ON user_addresses(user_id);
-- Index trên address_type để phân loại
CREATE INDEX idx_user_addresses_type ON user_addresses(address_type);

-- Index trên geohash để tìm kiếm theo vùng (dùng geohash_5 trước, sau đó geohash_7)
CREATE INDEX idx_user_addresses_geohash5 ON user_addresses(geohash_5);
CREATE INDEX idx_user_addresses_geohash7 ON user_addresses(geohash_7);

-- Index trên city và district để thống kê theo khu vực
CREATE INDEX idx_user_addresses_city ON user_addresses(city);
CREATE INDEX idx_user_addresses_district ON user_addresses(district);

-- Index trên is_default để lấy địa chỉ mặc định nhanh
CREATE INDEX idx_user_addresses_default ON user_addresses(user_id, is_default) WHERE is_default = true;

-- Index trên is_current để lọc địa chỉ đang ở
CREATE INDEX idx_user_addresses_current ON user_addresses(user_id, is_current) WHERE is_current = true;

-- Index trên is_verified để lọc địa chỉ đã xác thực
CREATE INDEX idx_user_addresses_verified ON user_addresses(is_verified);

-- Index trên moved_in_at, moved_out_at cho phân tích thời gian
CREATE INDEX idx_user_addresses_moved_in ON user_addresses(moved_in_at);
CREATE INDEX idx_user_addresses_moved_out ON user_addresses(moved_out_at);

-- Index GIN cho tags (nếu cần tìm theo tag)
CREATE INDEX idx_user_addresses_tags ON user_addresses USING GIN (tags);