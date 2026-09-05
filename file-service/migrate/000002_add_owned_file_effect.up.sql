ALTER TABLE file
    ADD COLUMN IF NOT EXISTS owner_namespace VARCHAR(64),
    ADD COLUMN IF NOT EXISTS owner_key VARCHAR(128),
    ADD COLUMN IF NOT EXISTS owner_request_hash VARCHAR(64);

CREATE UNIQUE INDEX IF NOT EXISTS uq_file_owned_effect
ON file (owner_namespace, owner_key)
WHERE owner_namespace IS NOT NULL
  AND owner_key IS NOT NULL;
