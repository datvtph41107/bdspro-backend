CREATE OR REPLACE FUNCTION create_point(
    v_lat  NUMERIC,
    v_lng  NUMERIC,
    v_name TEXT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    pole_ids            INT[];      -- Mảng ID của pole_points thỏa điều kiện ST_Contains
    v_pt_tag_id         INT;        -- ID của bản ghi trong pole_point_tag
    partition_target_id BIGINT;     -- current_id cần gán cho map_points.id
BEGIN
    -----------------------------------------------------------------------------
    -- 1) Tìm danh sách pole_points có polygon chứa điểm (lng, lat) đầu vào
    -----------------------------------------------------------------------------
--    SELECT array_agg(id)
--      INTO pole_ids
--      FROM pole_points
--     WHERE ST_Contains(
--             geom::geometry,
--             ST_SetSRID(ST_MakePoint(v_lng, v_lat), 4326)::geometry
--           );

    -----------------------------------------------------------------------------
    -- 2) Từ pole_point_tag, lấy 1 bản ghi có pole_point_id ∈ pole_ids
    --    và (from_id + 999) > current_id
    -----------------------------------------------------------------------------
    SELECT id, current_id
      INTO v_pt_tag_id, partition_target_id
      FROM pole_points
     WHERE ST_Contains(
             geom::geometry,
             ST_SetSRID(ST_MakePoint(v_lng, v_lat), 4326)::geometry
           )
--		pole_point_id = ANY (pole_ids)
       AND (from_id + 100000) > current_id
     LIMIT 1;

--	RAISE NOTICE 'RESULT: %', partition_target_id;
    -----------------------------------------------------------------------------
    -- 3) Nếu partition_target_id khác NULL => tạo map_points + cập nhật current_id
    -----------------------------------------------------------------------------
    IF partition_target_id IS NOT NULL THEN
        -- Thêm bản ghi vào map_points
        INSERT INTO map_points(id, name, geom)
        VALUES (
            partition_target_id,
            v_name,
            ST_SetSRID(ST_MakePoint(v_lng, v_lat), 4326)::GEOGRAPHY
        );

        -- Tăng current_id của bản ghi pole_point_tag thêm 1
        UPDATE pole_points
           SET current_id = current_id + 1
         WHERE id = v_pt_tag_id;

    -----------------------------------------------------------------------------
    -- 4) Nếu partition_target_id = NULL => gọi hàm create_pole_map_points
    -----------------------------------------------------------------------------
    ELSE
		partition_target_id := (SELECT COUNT(*) * 100000 + 1 FROM pole_points);
        -- Ở đây ta giả sử from_id là 1000 hoặc giá trị phù hợp khác
        PERFORM create_point_with_new_partition(v_name, v_lat, v_lng, partition_target_id);
    END IF;
END;
$$;