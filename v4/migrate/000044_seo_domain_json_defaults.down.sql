ALTER TABLE seo_domain
    ALTER COLUMN ref_snapshot_json DROP DEFAULT,
    ALTER COLUMN metadata DROP DEFAULT;
