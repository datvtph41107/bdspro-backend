-- Đồ án chỉ hiện trên API client khi client_view = true. Mặc định true cho dữ liệu hiện có.
ALTER TABLE qh_planning_projects
    ADD COLUMN IF NOT EXISTS client_view BOOLEAN NOT NULL DEFAULT true;

COMMENT ON COLUMN qh_planning_projects.client_view IS
    'true = đồ án hiển thị trên API client; false = chỉ admin thấy';
