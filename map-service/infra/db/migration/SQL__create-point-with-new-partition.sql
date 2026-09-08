-- 2) Tạo hoặc thay thế hàm create_pole_map_points
CREATE OR REPLACE FUNCTION create_point_with_new_partition(
    v_name    TEXT,
    v_lat     NUMERIC,
    v_lng     NUMERIC,
    v_from_id BIGINT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    v_pole_point_id INT;
    v_lat1 NUMERIC := v_lat - 0.5;
    v_lat2 NUMERIC := v_lat + 0.5;
    v_lng1 NUMERIC := v_lng - 0.5;
    v_lng2 NUMERIC := v_lng + 0.5;
BEGIN
    -------------------------------------------------------------------------
    -- 1) Thêm bản ghi vào pole_points (Polygon có 4 đỉnh quanh lat-lng)
    -------------------------------------------------------------------------
    INSERT INTO pole_points(geom, from_id, current_id)
    VALUES (
        ST_GeogFromText(
            FORMAT(
                'SRID=4326;POLYGON((%s %s, %s %s, %s %s, %s %s, %s %s))',
                v_lng1, v_lat1,
                v_lng1, v_lat2,
                v_lng2, v_lat2,
                v_lng2, v_lat1,
                v_lng1, v_lat1
            )
        ),
		v_from_id,
		v_from_id + 1
    )
    RETURNING id INTO v_pole_point_id;

    -------------------------------------------------------------------------
    -- 2) Tạo partition mới cho map_points
    -------------------------------------------------------------------------
    EXECUTE FORMAT($f$
        CREATE TABLE map_points_p_%s
        PARTITION OF map_points
        FOR VALUES FROM (%s) TO (%s)
    $f$, v_from_id, v_from_id, v_from_id + 100000);

    -------------------------------------------------------------------------
    -- 3) Thêm bản ghi vào map_points (id = v_from_id)
    -------------------------------------------------------------------------
    INSERT INTO map_points(id, name, geom)
    VALUES (
        v_from_id,
        v_name,
        ST_SetSRID(ST_MakePoint(v_lng, v_lat), 4326)::GEOGRAPHY
    );

    -------------------------------------------------------------------------
    -- 4) Thêm bản ghi vào pole_point_tag
    -------------------------------------------------------------------------
--    INSERT INTO pole_point_tag (pole_point_id, from_id, current_id)
--    VALUES (v_pole_point_id, v_from_id, v_from_id + 1);

END;
$$;
