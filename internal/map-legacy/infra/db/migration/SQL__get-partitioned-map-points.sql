CREATE OR REPLACE FUNCTION get_partitioned_map_points(
	input_geo GEOGRAPHY(POLYGON, 4326)
)
RETURNS TABLE (name VARCHAR, geom GEOGRAPHY(POINT, 4326))
AS $$
DECLARE
    partition_name TEXT;
    sql_query TEXT := '';
BEGIN
    -- Lặp qua danh sách các from_id hợp lệ để lấy danh sách partitions
    FOR partition_name IN
        SELECT DISTINCT 'map_points_p_' || from_id
        FROM pole_points p
        WHERE ST_Intersects(p.geom, input_geo)
    LOOP
        -- Tạo truy vấn động UNION ALL
--        sql_query := sql_query || 'SELECT name, geom FROM ' || partition_name
--		|| 'where ST_Contains(cast(geom as geometry),' || cast(input_geo as varchar) ||') UNION ALL ';
		sql_query := sql_query
            || 'SELECT name, geom FROM ' || partition_name
--            || ' WHERE ST_Intersects('
--            || 'geom, '
--            || 'ST_GeomFromText(''' || ST_AsText(input_geo) || ''', 4326)) '
            || ' UNION ALL ';
    END LOOP;

    -- Loại bỏ UNION ALL cuối cùng nếu có
    IF sql_query <> '' THEN
        sql_query := LEFT(sql_query, LENGTH(sql_query) - 11); -- Xóa ' UNION ALL' cuối cùng

        -- Thực thi truy vấn động với đúng số cột
        RETURN QUERY EXECUTE sql_query;
    END IF;
END $$ LANGUAGE plpgsql;