-- Repair request_stats / user_interaction upsert keys.
-- Tables may exist from AutoMigrate without PK/unique indexes required by ON CONFLICT.

ALTER TABLE request_stats_log
    ALTER COLUMN method SET NOT NULL,
    ALTER COLUMN path SET NOT NULL,
    ALTER COLUMN service SET NOT NULL,
    ALTER COLUMN status_code SET NOT NULL,
    ALTER COLUMN count SET NOT NULL,
    ALTER COLUMN timestamp SET NOT NULL,
    ALTER COLUMN bucket_hour SET NOT NULL;

CREATE TABLE request_stats_log_dedup AS
SELECT
    method,
    path,
    service,
    status_code,
    SUM(count) AS count,
    MAX(timestamp) AS timestamp,
    bucket_hour
FROM request_stats_log
GROUP BY method, path, service, status_code, bucket_hour;

TRUNCATE request_stats_log RESTART IDENTITY;

INSERT INTO request_stats_log (method, path, service, status_code, count, timestamp, bucket_hour)
SELECT method, path, service, status_code, count, timestamp, bucket_hour
FROM request_stats_log_dedup;

DROP TABLE request_stats_log_dedup;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'request_stats_log'::regclass AND contype = 'p'
    ) THEN
        ALTER TABLE request_stats_log ADD CONSTRAINT request_stats_log_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_request_stats_log_method ON request_stats_log (method);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_path ON request_stats_log (path);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_service ON request_stats_log (service);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_status_code ON request_stats_log (status_code);
CREATE INDEX IF NOT EXISTS idx_request_stats_log_timestamp ON request_stats_log (timestamp);
DROP INDEX IF EXISTS uq_request_stats_log_hour_bucket;
CREATE UNIQUE INDEX uq_request_stats_log_hour_bucket
    ON request_stats_log (method, path, service, status_code, bucket_hour);

UPDATE user_interaction_log SET organization_id = 0 WHERE organization_id IS NULL;

ALTER TABLE user_interaction_log
    ALTER COLUMN profile_id SET NOT NULL,
    ALTER COLUMN organization_id SET NOT NULL,
    ALTER COLUMN count SET NOT NULL,
    ALTER COLUMN timestamp SET NOT NULL,
    ALTER COLUMN bucket_hour SET NOT NULL;

CREATE TABLE user_interaction_log_dedup AS
SELECT
    profile_id,
    organization_id,
    SUM(count) AS count,
    MAX(timestamp) AS timestamp,
    bucket_hour
FROM user_interaction_log
GROUP BY profile_id, organization_id, bucket_hour;

TRUNCATE user_interaction_log RESTART IDENTITY;

INSERT INTO user_interaction_log (profile_id, organization_id, count, timestamp, bucket_hour)
SELECT profile_id, organization_id, count, timestamp, bucket_hour
FROM user_interaction_log_dedup;

DROP TABLE user_interaction_log_dedup;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'user_interaction_log'::regclass AND contype = 'p'
    ) THEN
        ALTER TABLE user_interaction_log ADD CONSTRAINT user_interaction_log_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_user_interaction_log_profile_id ON user_interaction_log (profile_id);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_organization_id ON user_interaction_log (organization_id);
CREATE INDEX IF NOT EXISTS idx_user_interaction_log_timestamp ON user_interaction_log (timestamp);
DROP INDEX IF EXISTS uq_user_interaction_log_hour_bucket;
CREATE UNIQUE INDEX uq_user_interaction_log_hour_bucket
    ON user_interaction_log (profile_id, organization_id, bucket_hour);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'request_audit_log'
    ) AND NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'request_audit_log'::regclass AND contype = 'p'
    ) THEN
        ALTER TABLE request_audit_log ADD CONSTRAINT request_audit_log_pkey PRIMARY KEY (id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_request_audit_log_method ON request_audit_log (method);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_path ON request_audit_log (path);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_service ON request_audit_log (service);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_status_code ON request_audit_log (status_code);
CREATE INDEX IF NOT EXISTS idx_request_audit_log_timestamp ON request_audit_log (timestamp);
