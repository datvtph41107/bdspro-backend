CREATE TABLE IF NOT EXISTS file_access (
    file_id   BIGINT NOT NULL,
    access_id BIGINT NOT NULL
);

-- Existing installations may already have this table from the retired
-- AutoMigrate runtime. Index creation is idempotent so the migration adopts
-- that durable fact without requiring runtime schema authority.
CREATE UNIQUE INDEX IF NOT EXISTS uq_file_access_file_access
ON file_access (file_id, access_id);

CREATE INDEX IF NOT EXISTS idx_file_access_access_id
ON file_access (access_id);
