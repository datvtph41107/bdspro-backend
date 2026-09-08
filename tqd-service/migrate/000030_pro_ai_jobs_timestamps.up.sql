-- 000004_pro_ai_jobs_timestamps.up.sql
-- Thời gian hoàn thành job + thời gian phê duyệt (created_at đã có)

ALTER TABLE pro_ai_jobs
    ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_completed_at ON pro_ai_jobs(completed_at);
CREATE INDEX IF NOT EXISTS idx_pro_ai_jobs_approved_at ON pro_ai_jobs(approved_at);

-- Backfill: job đã phân loại/lỗi → completed_at; đã duyệt → approved_at (+ completed nếu thiếu)
UPDATE pro_ai_jobs
SET completed_at = COALESCE(completed_at, updated_at, created_at)
WHERE deleted_at IS NULL
  AND process_status IN (30, 40, 50)
  AND completed_at IS NULL;

UPDATE pro_ai_jobs
SET approved_at = COALESCE(approved_at, updated_at, created_at)
WHERE deleted_at IS NULL
  AND process_status = 40
  AND approved_at IS NULL;
