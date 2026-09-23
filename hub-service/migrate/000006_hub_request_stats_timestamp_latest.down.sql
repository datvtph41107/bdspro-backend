DROP INDEX IF EXISTS uq_user_interaction_log_hour_bucket;
DROP INDEX IF EXISTS uq_request_stats_log_hour_bucket;

ALTER TABLE user_interaction_log DROP COLUMN IF EXISTS bucket_hour;
ALTER TABLE request_stats_log DROP COLUMN IF EXISTS bucket_hour;

-- Khôi phục unique 000005 (timestamp truncated = bucket).
CREATE UNIQUE INDEX IF NOT EXISTS uq_request_stats_log_hourly
    ON request_stats_log (method, path, service, status_code, "timestamp");

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_interaction_log_hourly
    ON user_interaction_log (profile_id, organization_id, "timestamp");
