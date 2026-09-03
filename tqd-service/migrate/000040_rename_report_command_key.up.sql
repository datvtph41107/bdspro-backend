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

    IF has_client_request_id AND NOT has_command_key THEN
        ALTER TABLE user_reported
            RENAME COLUMN client_request_id TO command_key;
    ELSIF has_client_request_id AND has_command_key THEN
        RAISE EXCEPTION
            'user_reported has both client_request_id and command_key; refusing ambiguous durable identity migration';
    ELSIF NOT has_client_request_id AND NOT has_command_key THEN
        RAISE EXCEPTION
            'user_reported has neither client_request_id nor command_key';
    END IF;
END $$;

DROP INDEX IF EXISTS uidx_user_report_client_request;

CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_report_command_key
    ON user_reported(user_id, command_key)
    WHERE command_key <> '' AND deleted_at IS NULL;

COMMIT;
