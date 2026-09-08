BEGIN;

DO $$
DECLARE
    has_client_request_id BOOLEAN;
    has_command_key BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1
          FROM information_schema.columns
         WHERE table_schema = current_schema()
           AND table_name = 'user_reported'
           AND column_name = 'client_request_id'
    ) INTO has_client_request_id;

    SELECT EXISTS (
        SELECT 1
          FROM information_schema.columns
         WHERE table_schema = current_schema()
           AND table_name = 'user_reported'
           AND column_name = 'command_key'
    ) INTO has_command_key;

    IF has_command_key AND NOT has_client_request_id THEN
        ALTER TABLE user_reported
            RENAME COLUMN command_key TO client_request_id;
    ELSIF has_client_request_id AND has_command_key THEN
        RAISE EXCEPTION
            'user_reported has both command_key and client_request_id; refusing ambiguous rollback';
    ELSIF NOT has_client_request_id AND NOT has_command_key THEN
        RAISE EXCEPTION
            'user_reported has neither command_key nor client_request_id';
    END IF;
END $$;

DROP INDEX IF EXISTS uidx_user_report_command_key;

CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_report_client_request
    ON user_reported(user_id, client_request_id)
    WHERE client_request_id <> '' AND deleted_at IS NULL;

COMMIT;
