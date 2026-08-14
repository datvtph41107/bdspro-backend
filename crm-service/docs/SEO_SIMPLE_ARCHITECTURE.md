# Kiến trúc SEO kỹ thuật tối giản

SEO chỉ là hệ thống quản lý và phát hành trang. Nó không lập chiến lược, chấm điểm, suy đoán ý định, đo hành vi hay tự quyết định trang nào “xứng đáng” được index.

## Luồng duy nhất

```text
Admin CRUD
    ↓
seo_domain (draft / published / archived)
    ↓ publish
Render HTML + public document
    ↓
Gateway ── sitemap.xml
    ├──── Next.js: SSR, metadata, canonical, robots, JSON-LD
    └──── React Native: page bootstrap và giao diện tương tác
```

Mọi client dùng cùng canonical, cùng cờ index và cùng phiên bản public document do CRM phát hành. Next.js và React Native không tính readiness, quality score hoặc index decision riêng.

## Mô hình dữ liệu còn lại

| Nhóm | Trường chính | Vai trò |
| --- | --- | --- |
| Định danh | `id`, `slug`, `canonical_url`, `origin_url` | Xác định duy nhất một trang |
| Nội dung | `title`, `description`, `summary`, `content`, `metadata` | Dữ liệu để render metadata, HTML và JSON-LD |
| Vòng đời | `page_status`, `published`, `published_at` | Thêm, sửa, xuất bản, gỡ xuất bản, lưu trữ, xóa |
| Search engine | `is_index`, `is_site_map`, `is_robot`, `sitemap_*` | Cấu hình kỹ thuật trực tiếp, không qua scoring |
| Render | `need_generate`, `render_status`, `rendered_html`, `static_html_*`, `generated_at` | Trạng thái triển khai trang |
| Nguồn tùy chọn | `ref_type`, `ref_id`, `ref_source`, snapshot và thời điểm đồng bộ | Chỉ hydrate nội dung; không được thay đổi cờ index/publish/sitemap |
| Liên kết | `seo_links`, `seo_relative` | Quan hệ kỹ thuật giữa các trang |
| Nhật ký | `seo_generation_log` | Retry và chẩn đoán render |

`classify`, `quality_score` và `visible_on` đã bị loại khỏi aggregate và API. Các bảng objective, audience, user need, evidence, intent, query cluster, taxonomy, category, placement, measurement, ranking, traffic, crawl, issue, trend và runtime-event bị xóa bởi migration `000024_simplify_seo_to_pages`.

## Quy tắc kỹ thuật

- Tạo/sửa/xóa chỉ thao tác trên page và các liên kết của page.
- Publish yêu cầu `slug`, `canonical_url` và `title`; sau publish trang được đưa vào hàng đợi render.
- Public chỉ phục vụ trang `published`, scope `public`, render `success` và có HTML.
- Sitemap chỉ lấy trang đã publish, bật `is_index`, bật `is_site_map`, render thành công và có canonical.
- `robots` và metadata lấy trực tiếp từ trường kỹ thuật của page.
- Source adapter chỉ cung cấp snapshot hoặc gợi ý nội dung. Không có authority, evidence, readiness hay quality gate.
- Telemetry hành vi SEO không nằm trong luồng phát hành.
- Danh sách tin public chỉ lấy trang đã publish và sắp xếp mới nhất trước; không có feed nổi bật, đọc nhiều, xu hướng, category hay placement biên tập.

## Hợp đồng public

Public document `qhpro-public-seo-document/v1` chỉ chứa:

- identity, heading, breadcrumbs;
- facts, sections, media, FAQ, related links/entities, map target;
- title, description, canonical, indexable, follow và image;
- tên/URL nguồn cùng thời điểm cập nhật/phát hành.

Thay đổi schema phải bắt đầu từ protobuf CRM, sinh lại code, rồi cập nhật validator Next.js và React Native. Không thêm lại các trường mang tính đánh giá hoặc quyết định hành vi.

## Triển khai migration

Migration `000024` là phá hủy có chủ đích và không thể khôi phục dữ liệu chiến lược cũ. Trước khi chạy ở môi trường thật cần backup database và xác nhận không còn consumer dùng các endpoint đã loại bỏ. Migration chưa được tự động chạy trong lần refactor này.
