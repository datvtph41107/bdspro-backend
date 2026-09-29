BEGIN;

DROP TABLE IF EXISTS report_jobs;
DROP TABLE IF EXISTS quota_usage_events;

ALTER TABLE user_reported
    DROP COLUMN IF EXISTS request_hash;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
          FROM user_reported
         WHERE LENGTH(client_request_id) > 100
    ) THEN
        RAISE EXCEPTION
            'cannot narrow user_reported.client_request_id to varchar(100): data exceeds 100 characters';
    END IF;
END $$;

ALTER TABLE user_reported
    ALTER COLUMN client_request_id TYPE VARCHAR(100);

COMMIT;
