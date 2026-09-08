BEGIN;

-- SEO chỉ còn là hệ thống kỹ thuật để quản lý và phát hành trang.
-- Các bảng chiến lược, đối tượng, bằng chứng, hành vi và chấm điểm được loại bỏ.
-- Mô hình còn lại phục vụ CRUD, xuất bản, render, sitemap và metadata.

DROP TABLE IF EXISTS seo_query_cluster_member CASCADE;
DROP TABLE IF EXISTS seo_query_cluster CASCADE;
DROP TABLE IF EXISTS seo_intent_evidence_link CASCADE;
DROP TABLE IF EXISTS seo_intent_hypothesis CASCADE;
DROP TABLE IF EXISTS seo_query_evidence CASCADE;
DROP TABLE IF EXISTS seo_user_need_evidence CASCADE;
DROP TABLE IF EXISTS seo_user_need CASCADE;
DROP TABLE IF EXISTS seo_audience_segment CASCADE;
DROP TABLE IF EXISTS seo_objective_metric CASCADE;
DROP TABLE IF EXISTS seo_objective_evidence CASCADE;
DROP TABLE IF EXISTS seo_objective CASCADE;

DROP TABLE IF EXISTS seo_content_placement CASCADE;
DROP TABLE IF EXISTS seo_content_category_assignment CASCADE;
DROP TABLE IF EXISTS seo_category CASCADE;
DROP TABLE IF EXISTS seo_content_category CASCADE;
DROP TABLE IF EXISTS seo_content_domain CASCADE;

DROP TABLE IF EXISTS seo_runtime_event CASCADE;
DROP TABLE IF EXISTS seo_page_score_snapshot CASCADE;
DROP TABLE IF EXISTS seo_measurement_source_state CASCADE;
DROP TABLE IF EXISTS seo_measurement_import_run CASCADE;
DROP TABLE IF EXISTS seo_ranking_snapshot CASCADE;
DROP TABLE IF EXISTS seo_keyword_target CASCADE;
DROP TABLE IF EXISTS seo_traffic_daily CASCADE;
DROP TABLE IF EXISTS seo_search_performance_daily CASCADE;
DROP TABLE IF EXISTS seo_url_inspection_snapshot CASCADE;
DROP TABLE IF EXISTS seo_crawl_event CASCADE;
DROP TABLE IF EXISTS seo_issue CASCADE;
DROP TABLE IF EXISTS seo_trend CASCADE;

-- Loại bỏ các cột chiến lược và chấm điểm khỏi trang SEO kỹ thuật.
DROP INDEX IF EXISTS idx_seo_domain_classify_active;
DROP INDEX IF EXISTS idx_seo_domain_quality_score_active;
DROP INDEX IF EXISTS idx_seo_domain_visible_on_public;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_classify;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_quality_score;
ALTER TABLE IF EXISTS seo_domain DROP CONSTRAINT IF EXISTS chk_seo_domain_visible_on;
ALTER TABLE IF EXISTS seo_domain
    DROP COLUMN IF EXISTS classify,
    DROP COLUMN IF EXISTS quality_score,
    DROP COLUMN IF EXISTS visible_on;

-- Tài liệu ngay trong PostgreSQL cho mô hình SEO kỹ thuật còn lại.
COMMENT ON TABLE seo_domain IS 'Mỗi bản ghi là một trang SEO từ bản nháp đến HTML đã phát hành.';
COMMENT ON COLUMN seo_domain.id IS 'Khóa chính của trang SEO.';
COMMENT ON COLUMN seo_domain.slug IS 'Phần định danh cuối đường dẫn.';
COMMENT ON COLUMN seo_domain.origin_url IS 'Đường dẫn gốc trước khi chuẩn hóa.';
COMMENT ON COLUMN seo_domain.canonical_url IS 'URL chính thức dành cho máy tìm kiếm.';
COMMENT ON COLUMN seo_domain.ref_type IS 'Loại dữ liệu nguồn được liên kết; 0 là nhập thủ công.';
COMMENT ON COLUMN seo_domain.ref_id IS 'Mã bản ghi tại hệ thống nguồn.';
COMMENT ON COLUMN seo_domain.ref_source IS 'Khóa bộ kết nối dùng để đọc dữ liệu nguồn.';
COMMENT ON COLUMN seo_domain.ref_label IS 'Tên hiển thị của bản ghi nguồn.';
COMMENT ON COLUMN seo_domain.ref_url IS 'Đường dẫn tham chiếu đến dữ liệu nguồn.';
COMMENT ON COLUMN seo_domain.source_status IS 'Trạng thái đồng bộ hiện tại của dữ liệu nguồn.';
COMMENT ON COLUMN seo_domain.ref_missing IS 'Đánh dấu bản ghi nguồn không còn tồn tại.';
COMMENT ON COLUMN seo_domain.ref_snapshot_json IS 'Bản chụp dữ liệu nguồn dùng để render ổn định.';
COMMENT ON COLUMN seo_domain.ref_hash IS 'Mã nhận biết dữ liệu nguồn đã thay đổi.';
COMMENT ON COLUMN seo_domain.ref_last_synced_at IS 'Thời điểm đồng bộ nguồn gần nhất.';
COMMENT ON COLUMN seo_domain.source_updated_at IS 'Thời điểm dữ liệu nguồn thay đổi gần nhất.';
COMMENT ON COLUMN seo_domain.scope IS 'Phạm vi sử dụng của trang: public, internal, app hoặc admin.';
COMMENT ON COLUMN seo_domain.page_status IS 'Vòng đời trang: draft, published hoặc archived.';
COMMENT ON COLUMN seo_domain.title IS 'Tiêu đề chính và tiêu đề SEO mặc định.';
COMMENT ON COLUMN seo_domain.description IS 'Nội dung cho thẻ meta description.';
COMMENT ON COLUMN seo_domain.summary IS 'Phần tóm tắt hiển thị đầu trang.';
COMMENT ON COLUMN seo_domain.content IS 'Nội dung đầy đủ dùng để render trang.';
COMMENT ON COLUMN seo_domain.published IS 'Đánh dấu trang đã được phát hành.';
COMMENT ON COLUMN seo_domain.published_at IS 'Thời điểm phát hành gần nhất.';
COMMENT ON COLUMN seo_domain.is_index IS 'Cho phép máy tìm kiếm lập chỉ mục trang.';
COMMENT ON COLUMN seo_domain.is_site_map IS 'Cho phép đưa trang vào sitemap.';
COMMENT ON COLUMN seo_domain.is_robot IS 'Cho phép bot thu thập trang.';
COMMENT ON COLUMN seo_domain.site_map_lasted_at IS 'Thời điểm dữ liệu sitemap được cập nhật.';
COMMENT ON COLUMN seo_domain.sitemap_priority IS 'Độ ưu tiên kỹ thuật ghi trong sitemap.';
COMMENT ON COLUMN seo_domain.sitemap_change_freq IS 'Tần suất thay đổi gợi ý trong sitemap.';
COMMENT ON COLUMN seo_domain.need_generate IS 'Đánh dấu trang đang cần tạo lại HTML.';
COMMENT ON COLUMN seo_domain.render_status IS 'Trạng thái hiện tại của quá trình tạo HTML.';
COMMENT ON COLUMN seo_domain.rendered_html IS 'HTML hoàn chỉnh được phục vụ công khai.';
COMMENT ON COLUMN seo_domain.last_render_error IS 'Lỗi tạo HTML gần nhất.';
COMMENT ON COLUMN seo_domain.generated_at IS 'Thời điểm tạo HTML thành công gần nhất.';
COMMENT ON COLUMN seo_domain.static_html_path IS 'Vị trí tệp HTML đã tạo.';
COMMENT ON COLUMN seo_domain.static_html_hash IS 'Mã nhận biết nội dung HTML thay đổi.';
COMMENT ON COLUMN seo_domain.template_key IS 'Khóa mẫu HTML dùng để render.';
COMMENT ON COLUMN seo_domain.template_version IS 'Phiên bản mẫu HTML đã sử dụng.';
COMMENT ON COLUMN seo_domain.deep_link IS 'Đường dẫn mở đúng màn hình trong ứng dụng.';
COMMENT ON COLUMN seo_domain.metadata IS 'Dữ liệu mở rộng như JSON-LD, breadcrumb và Open Graph.';
COMMENT ON COLUMN seo_domain.note IS 'Ghi chú vận hành dành cho quản trị viên.';
COMMENT ON COLUMN seo_domain.created_at IS 'Thời điểm tạo bản ghi.';
COMMENT ON COLUMN seo_domain.updated_at IS 'Thời điểm cập nhật bản ghi gần nhất.';
COMMENT ON COLUMN seo_domain.deleted_at IS 'Thời điểm xóa mềm bản ghi.';

COMMENT ON TABLE seo_links IS 'Các liên kết HTML giữa những trang SEO.';
COMMENT ON COLUMN seo_links.id IS 'Khóa chính của liên kết.';
COMMENT ON COLUMN seo_links.parent_seo_id IS 'Trang chứa liên kết.';
COMMENT ON COLUMN seo_links.child_seo_id IS 'Trang được liên kết tới.';
COMMENT ON COLUMN seo_links.link IS 'URL được ghi vào HTML.';
COMMENT ON COLUMN seo_links.title IS 'Nội dung hiển thị của liên kết.';
COMMENT ON COLUMN seo_links.link_type IS 'Vị trí hoặc mục đích kỹ thuật của liên kết.';
COMMENT ON COLUMN seo_links.priority IS 'Thứ tự ưu tiên khi render liên kết.';
COMMENT ON COLUMN seo_links.created_at IS 'Thời điểm tạo liên kết.';
COMMENT ON COLUMN seo_links.updated_at IS 'Thời điểm cập nhật liên kết gần nhất.';
COMMENT ON COLUMN seo_links.deleted_at IS 'Thời điểm xóa mềm liên kết.';

COMMENT ON TABLE seo_relative IS 'Các quan hệ dữ liệu giữa những trang SEO.';
COMMENT ON COLUMN seo_relative.id IS 'Khóa chính của quan hệ.';
COMMENT ON COLUMN seo_relative.parent_seo_id IS 'Trang nguồn của quan hệ.';
COMMENT ON COLUMN seo_relative.child_seo_id IS 'Trang liên quan.';
COMMENT ON COLUMN seo_relative.relation_type IS 'Loại quan hệ giữa hai trang.';
COMMENT ON COLUMN seo_relative.priority IS 'Thứ tự ưu tiên khi đọc quan hệ.';
COMMENT ON COLUMN seo_relative.created_at IS 'Thời điểm tạo quan hệ.';
COMMENT ON COLUMN seo_relative.updated_at IS 'Thời điểm cập nhật quan hệ gần nhất.';
COMMENT ON COLUMN seo_relative.deleted_at IS 'Thời điểm xóa mềm quan hệ.';

COMMENT ON TABLE seo_generation_log IS 'Lịch sử từng lần tạo HTML của trang SEO.';
COMMENT ON COLUMN seo_generation_log.id IS 'Khóa chính của lần tạo HTML.';
COMMENT ON COLUMN seo_generation_log.seo_domain_id IS 'Trang SEO được tạo HTML.';
COMMENT ON COLUMN seo_generation_log.job_id IS 'Mã tác vụ bên ngoài nếu tiến trình chạy qua hàng đợi.';
COMMENT ON COLUMN seo_generation_log.action IS 'Tên thao tác kỹ thuật của tiến trình cũ.';
COMMENT ON COLUMN seo_generation_log.trigger_type IS 'Nguyên nhân bắt đầu tiến trình.';
COMMENT ON COLUMN seo_generation_log.status IS 'Trạng thái xử lý của tiến trình.';
COMMENT ON COLUMN seo_generation_log.static_html_path IS 'Vị trí HTML được tạo.';
COMMENT ON COLUMN seo_generation_log.static_html_hash IS 'Mã nhận diện nội dung HTML.';
COMMENT ON COLUMN seo_generation_log.error_message IS 'Lỗi khi tiến trình thất bại.';
COMMENT ON COLUMN seo_generation_log.metadata IS 'Thông tin kỹ thuật bổ sung của tiến trình.';
COMMENT ON COLUMN seo_generation_log.started_at IS 'Thời điểm bắt đầu xử lý.';
COMMENT ON COLUMN seo_generation_log.finished_at IS 'Thời điểm kết thúc xử lý.';
COMMENT ON COLUMN seo_generation_log.duration_ms IS 'Thời gian xử lý tính bằng mili giây.';
COMMENT ON COLUMN seo_generation_log.created_at IS 'Thời điểm tạo bản ghi tiến trình.';
COMMENT ON COLUMN seo_generation_log.updated_at IS 'Thời điểm cập nhật tiến trình gần nhất.';
COMMENT ON COLUMN seo_generation_log.deleted_at IS 'Thời điểm xóa mềm bản ghi tiến trình.';

COMMIT;
