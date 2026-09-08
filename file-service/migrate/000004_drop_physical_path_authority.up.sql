-- Physical deployment paths are runtime storage details, not durable File
-- identity. All reads resolve the stable File path/reference through the
-- configured StorageService root, so the database must not retain a mount-
-- specific absolute path.
ALTER TABLE file
    DROP COLUMN IF EXISTS absolute_path;
