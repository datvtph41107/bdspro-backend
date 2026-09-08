-- List admin order by created_at; soft-delete filter.
CREATE INDEX IF NOT EXISTS idx_reports_created_at_alive
  ON reports (created_at DESC)
  WHERE deleted_at IS NULL;
