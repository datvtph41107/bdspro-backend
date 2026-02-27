-- tạo partition cron theo tháng
CREATE OR REPLACE FUNCTION create_property_partition(
    p_type TEXT,
    p_date DATE
)
RETURNS VOID AS $$
DECLARE
    start_date DATE;
    end_date DATE;
    partition_name TEXT;
BEGIN
    start_date := date_trunc('month', p_date);
    end_date := start_date + INTERVAL '1 month';

    partition_name := format('property_%s_%s',
        p_type,
        to_char(start_date, 'YYYYMM')
    );

    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF property_%s
         FOR VALUES FROM (%L) TO (%L)',
        partition_name,
        p_type,
        start_date,
        end_date
    );
END;
$$ LANGUAGE plpgsql;

-- FUNCTION CLEAN DEFAULT
CREATE OR REPLACE FUNCTION move_default_partition_data()
RETURNS VOID AS $$
BEGIN
    -- SYSTEM
    INSERT INTO property_system
    SELECT * FROM property_system_default
    ON CONFLICT DO NOTHING;

    DELETE FROM property_system_default;

    -- USER
    INSERT INTO property_user
    SELECT * FROM property_user_default
    ON CONFLICT DO NOTHING;

    DELETE FROM property_user_default;
END;
$$ LANGUAGE plpgsql;