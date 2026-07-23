# QHPro Public SEO Document v1

## Mục đích

`qhpro-public-seo-document/v1` là semantic document contract nằm giữa dữ liệu nguồn và HTML renderer.

```text
TQD planning-project projection
        │
        ▼
PublicSeoDocument v1
        │
        ├── Answer Readiness
        ├── index eligibility
        ├── visible rich HTML
        ├── JSON-LD
        └── telemetry hooks
```

Contract này không thay thế `seo_domain`:

- `seo_domain` vẫn sở hữu lifecycle, publication, canonical URL, title, description, index/sitemap policy và rendered artifact.
- Public document sở hữu cấu trúc nội dung: breadcrumbs, facts, sections, FAQ, media, relations, map target và source disclosure.

## Vị trí lưu trong Gate 2.2A

Gate 2.2A đọc document từ `seo_domain.metadata`, hỗ trợ hai dạng:

```json
{
    "schemaVersion": "qhpro-public-seo-document/v1",
    "moduleId": "planning-project"
}
```

hoặc:

```json
{
    "publicDocument": {
        "schemaVersion": "qhpro-public-seo-document/v1",
        "moduleId": "planning-project"
    }
}
```

Metadata legacy có `@context` được bỏ qua để tránh hai nguồn JSON-LD.

## Required planning-project answer specification

Facts bắt buộc:

- `planning-type`
- `status`
- `administrative-unit`
- `updated-at`

Sections bắt buộc:

- `overview`
- `research-scope`
- `legal-timeline`
- `source`

Ngoài ra planning-project cần:

- map target hợp lệ;
- source name hoặc source URL;
- source authority/quality đủ tin cậy;
- environment bằng `production`.

Nếu thiếu, renderer vẫn tạo HTML nhưng tự hạ xuống:

```html
<meta name="robots" content="noindex,follow" />
```

## Supported blocks

```text
paragraph
list
notice
 timeline
```

`items` của `list` là chuỗi. `items` của `timeline` là object gồm `id`, `title`, `date`, `description`.

## Security

- Go `html/template` escape visible HTML.
- Rich blocks là controlled structure, không render raw arbitrary HTML.
- Legacy `content` được chuyển thành plain text.
- JSON-LD sinh từ semantic document bằng `json.Marshal` và không dùng metadata JSON-LD cũ.
- Internal links và map paths phải là public relative path hợp lệ.

## Gate 2.2B

Gate 2.2A dùng script operator để map TQD projection thật vào metadata nhằm kiểm chứng renderer. Đây không phải runtime dependency.

Gate 2.2B sẽ thêm TQD planning-project projection vào internal gRPC contract để CRM tự compose document khi sync/generate, loại bỏ bước operator sync.
