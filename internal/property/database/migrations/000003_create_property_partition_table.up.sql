CREATE TABLE property_system
PARTITION OF property
FOR VALUES IN (10)
PARTITION BY RANGE (created_at);

CREATE UNIQUE INDEX ux_property_system_identifier
ON property_system (identifier, created_at);

CREATE TABLE property_user
PARTITION OF property
FOR VALUES IN (20)
PARTITION BY RANGE (created_at);

-- DEFAULT PARTITIONS: insert dữ liệu khi chưa có partition tháng
CREATE TABLE property_system_default
PARTITION OF property_system DEFAULT;

CREATE TABLE property_user_default
PARTITION OF property_user DEFAULT;

CREATE TABLE property_system_202602
PARTITION OF property_system
FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE property_user_202602
PARTITION OF property_user
FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE property_system_202603
PARTITION OF property_system
FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE TABLE property_user_202603
PARTITION OF property_user
FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');