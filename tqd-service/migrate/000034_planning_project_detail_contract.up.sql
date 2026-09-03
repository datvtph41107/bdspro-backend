BEGIN;

COMMENT ON COLUMN qh_planning_projects.research_scope IS
    'Mô tả phạm vi, ranh giới hoặc khu vực nghiên cứu được công bố trong hồ sơ đồ án quy hoạch.';

COMMENT ON COLUMN qh_planning_projects.indicators IS
    'Các chỉ tiêu quy hoạch chính của đồ án, lưu dưới dạng JSON có cấu trúc để hiển thị và truy vấn theo từng chỉ tiêu.';

COMMIT;
