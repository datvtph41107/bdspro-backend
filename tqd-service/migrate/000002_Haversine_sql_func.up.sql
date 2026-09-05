-- ===================================================
-- FUNCTION: Haversine Distance cho PostgreSQL
-- MÔ TẢ: Tính khoảng cách giữa 2 điểm theo tọa độ (km)
-- ===================================================

CREATE OR REPLACE FUNCTION haversine_distance(
    lat1 DECIMAL(10,8), 
    lng1 DECIMAL(11,8), 
    lat2 DECIMAL(10,8), 
    lng2 DECIMAL(11,8)
) RETURNS DECIMAL(10,2) AS $$
DECLARE
    R CONSTANT DECIMAL(5,2) := 6371; -- Bán kính trái đất (km)
    dlat DECIMAL(10,8);
    dlng DECIMAL(11,8);
    a DECIMAL(10,8);
    c DECIMAL(10,8);
BEGIN
    -- Convert độ sang radian
    dlat := RADIANS(lat2 - lat1);
    dlng := RADIANS(lng2 - lng1);
    
    -- Công thức Haversine
    a := SIN(dlat/2) * SIN(dlat/2) +
         COS(RADIANS(lat1)) * COS(RADIANS(lat2)) *
         SIN(dlng/2) * SIN(dlng/2);
    
    c := 2 * ATAN2(SQRT(a), SQRT(1-a));
    
    -- Kết quả (km)
    RETURN R * c;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- ===================================================
-- TEST FUNCTION
-- ===================================================
-- SELECT haversine_distance(21.0285, 105.8542, 10.8231, 106.6297); -- Hà Nội -> HCM