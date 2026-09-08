-- =====================================================
-- MIGRATION UP: qh_label_layers — many-to-many label ↔ layer
-- =====================================================

-- Lịch sử cũ tạo bảng nối trước bảng nhãn ở version 000025. Canonical SQL
-- phải chạy được trên database mới, vì vậy bảng cha được khởi tạo tại đúng
-- dependency đầu tiên; version 000025 chỉ hoàn thiện index/trigger.
CREATE TABLE IF NOT EXISTS qh_labels (
    id          BIGSERIAL PRIMARY KEY,
    layer_id    BIGINT       NOT NULL REFERENCES qh_layers(id),
    name        VARCHAR(255) NOT NULL,
    color       VARCHAR(20)  NOT NULL DEFAULT '#1890ff',
    description TEXT,
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS qh_label_layers (
    label_id     BIGINT NOT NULL REFERENCES qh_labels(id) ON DELETE CASCADE,
    layer_id  BIGINT NOT NULL REFERENCES qh_layers(id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (label_id, layer_id)
);

CREATE INDEX IF NOT EXISTS idx_qh_label_layers_layer_id
    ON qh_label_layers(layer_id);

COMMENT ON TABLE qh_label_layers IS 'Liên kết nhiều-nhiều: một nhãn có thể thuộc nhiều layer';

-- Đồng bộ dữ liệu cũ (qh_labels.layer_id) vào bảng nối
INSERT INTO qh_label_layers (label_id, layer_id)
SELECT l.id, l.layer_id
FROM qh_labels l
WHERE l.deleted_at IS NULL
  AND l.layer_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM qh_label_layers x
    WHERE x.label_id = l.id AND x.layer_id = l.layer_id
  );
