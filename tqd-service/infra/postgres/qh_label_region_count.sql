-- Đồng bộ qh_labels.region_count theo một layer (chạy tay khi cần repair).
-- Thường không cần: migration 000008 gắn trigger trên qh_regions tự cập nhật.
-- Tham số: $1 = layer_id (prepared statement / nhiều client SQL).

UPDATE qh_labels
SET region_count = 0
WHERE layer_id = $1
  AND deleted_at IS NULL;

WITH counts AS (
    SELECT label_id, COUNT(*)::bigint AS cnt
    FROM qh_regions
    WHERE layer_id = $1
      AND deleted_at IS NULL
    GROUP BY label_id
)
UPDATE qh_labels l
SET region_count = counts.cnt
FROM counts
WHERE l.id = counts.label_id
  AND l.layer_id = $1
  AND l.deleted_at IS NULL;
