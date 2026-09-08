CREATE TABLE file (
    id             BIGSERIAL PRIMARY KEY,
    created_by     BIGINT,
    updated_by     BIGINT,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,
    name           VARCHAR(255),
    path           VARCHAR(255),
    thumbnail_path VARCHAR(255),
    absolute_path  VARCHAR(255),
    extension      VARCHAR(20),
    size           BIGINT,
    hash           VARCHAR(255),
    description    TEXT,
    mine           VARCHAR(50)
);

CREATE INDEX idx_file_deleted_at
ON file (deleted_at);
