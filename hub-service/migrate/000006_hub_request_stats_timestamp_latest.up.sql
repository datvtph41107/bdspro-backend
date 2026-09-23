-- timestamp = lần update gần nhất; bucket_hour = khóa giờ (chỉ INSERT khi giờ mới).

ALTER TABLE request_stats_log
    ADD COLUMN IF NOT EXISTS bucket_hour TIMESTAMPTZ;

UPDATE request_stats_log
SET bucket_hour = date_trunc('hour', "timestamp")
WHERE bucket_hour IS NULL;

-- Gộp trùng (method/path/service/status/hour) trước khi tạo unique.
CREATE TABLE request_stats_log_dedup AS
SELECT
    method,
    path,
    service,
    status_code,
    SUM(count)::bigint AS count,
    MAX("timestamp") AS timestamp,
    bucket_hour
FROM request_stats_log
GROUP BY method, path, service, status_code, bucket_hour;

TRUNCATE request_stats_log;

INSERT INTO request_stats_log (method, path, service, status_code, count, timestamp, bucket_hour)
SELECT method, path, service, status_code, count, timestamp, bucket_hour
FROM request_stats_log_dedup;

DROP TABLE request_stats_log_dedup;

ALTER TABLE request_stats_log
    ALTER COLUMN bucket_hour SET NOT NULL;

DROP INDEX IF EXISTS uq_request_stats_log_hourly;
DROP INDEX IF EXISTS uq_request_stats_log_hour_bucket;

CREATE UNIQUE INDEX IF NOT EXISTS uq_request_stats_log_hour_bucket
    ON request_stats_log (method, path, service, status_code, bucket_hour);

ALTER TABLE user_interaction_log
    ADD COLUMN IF NOT EXISTS bucket_hour TIMESTAMPTZ;

UPDATE user_interaction_log
SET bucket_hour = date_trunc('hour', "timestamp")
WHERE bucket_hour IS NULL;

CREATE TABLE user_interaction_log_dedup AS
SELECT
    profile_id,
    COALESCE(organization_id, 0) AS organization_id,
    SUM(count)::bigint AS count,
    MAX("timestamp") AS timestamp,
    bucket_hour
FROM user_interaction_log
GROUP BY profile_id, COALESCE(organization_id, 0), bucket_hour;

TRUNCATE user_interaction_log;

INSERT INTO user_interaction_log (profile_id, organization_id, count, timestamp, bucket_hour)
SELECT profile_id, organization_id, count, timestamp, bucket_hour
FROM user_interaction_log_dedup;

DROP TABLE user_interaction_log_dedup;

UPDATE user_interaction_log SET organization_id = 0 WHERE organization_id IS NULL;

ALTER TABLE user_interaction_log
    ALTER COLUMN bucket_hour SET NOT NULL;

DROP INDEX IF EXISTS uq_user_interaction_log_hourly;
DROP INDEX IF EXISTS uq_user_interaction_log_hour_bucket;

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_interaction_log_hour_bucket
    ON user_interaction_log (profile_id, organization_id, bucket_hour);
