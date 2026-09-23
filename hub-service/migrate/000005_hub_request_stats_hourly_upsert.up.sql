-- Gộp request/user-interaction stats theo bucket giờ; flush chỉ tăng count
-- nếu (key + hour) đã tồn tại — không INSERT thêm row mỗi lần flush.

CREATE TABLE request_stats_log_hourly AS
SELECT
    method,
    path,
    service,
    status_code,
    SUM(count)::bigint AS count,
    date_trunc('hour', "timestamp") AS timestamp
FROM request_stats_log
GROUP BY method, path, service, status_code, date_trunc('hour', "timestamp");

TRUNCATE request_stats_log;

INSERT INTO request_stats_log (method, path, service, status_code, count, timestamp)
SELECT method, path, service, status_code, count, timestamp
FROM request_stats_log_hourly;

DROP TABLE request_stats_log_hourly;

CREATE UNIQUE INDEX IF NOT EXISTS uq_request_stats_log_hourly
    ON request_stats_log (method, path, service, status_code, "timestamp");

CREATE TABLE user_interaction_log_hourly AS
SELECT
    profile_id,
    COALESCE(organization_id, 0) AS organization_id,
    SUM(count)::bigint AS count,
    date_trunc('hour', "timestamp") AS timestamp
FROM user_interaction_log
GROUP BY profile_id, COALESCE(organization_id, 0), date_trunc('hour', "timestamp");

TRUNCATE user_interaction_log;

INSERT INTO user_interaction_log (profile_id, organization_id, count, timestamp)
SELECT profile_id, organization_id, count, timestamp
FROM user_interaction_log_hourly;

DROP TABLE user_interaction_log_hourly;

UPDATE user_interaction_log SET organization_id = 0 WHERE organization_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_interaction_log_hourly
    ON user_interaction_log (profile_id, organization_id, "timestamp");
