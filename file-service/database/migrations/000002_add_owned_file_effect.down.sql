DROP INDEX IF EXISTS uq_file_owned_effect;

ALTER TABLE file
    DROP COLUMN IF EXISTS owner_request_hash,
    DROP COLUMN IF EXISTS owner_key,
    DROP COLUMN IF EXISTS owner_namespace;
