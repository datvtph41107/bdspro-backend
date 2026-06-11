---
title: "SEO Domain API - Cơ Chế Hoạt Động & Xử Lý Nghiệp Vụ"
author: "BDSPro CRM Service"
date: "2026-06-11"
output:
  html_document:
    toc: true
    toc_depth: 4
    number_sections: true
    theme: readable
  github_document:
    toc: true
    toc_depth: 4
params:
  service_name: "SeoDomainService"
  endpoint_base: "/v2/crm/admin"
  generated_from: "seo_domain_impl_v2.zip"
---

```{r setup, include=FALSE}
knitr::opts_chunk$set(echo = TRUE, warning = FALSE, message = FALSE)
```

# 1. Mục đích tài liệu

Tài liệu này giải thích **cơ chế hoạt động xử lý nghiệp vụ** của nhóm API `SeoDomainService` trong CRM service.

Tài liệu không chỉ liệt kê endpoint, mà mô tả đầy đủ:

- API nhận request như thế nào.
- Handler map protobuf request sang DTO như thế nào.
- Usecase kiểm tra nghiệp vụ gì.
- Repository thao tác DB ra sao.
- Response trả dữ liệu gì.
- API đó phục vụ cơ chế SEO nào: canonical, sitemap, index, robot, metadata, internal link, relative graph.
- Dẫn chứng từ code implementation và tài liệu nền SEO Domain.

Phạm vi API trong file hiện tại:

```text
POST   /v2/crm/admin/seo-domains
GET    /v2/crm/admin/seo-domains/{id}
GET    /v2/crm/admin/seo-domains
PUT    /v2/crm/admin/seo-domains/{id}
DELETE /v2/crm/admin/seo-domains/{id}

POST   /v2/crm/admin/seo-domains/internal-links
GET    /v2/crm/admin/seo-domains/internal-links
PUT    /v2/crm/admin/seo-domains/internal-links/{id}
DELETE /v2/crm/admin/seo-domains/internal-links/{id}

POST   /v2/crm/admin/seo-domains/relatives
GET    /v2/crm/admin/seo-domains/relatives
DELETE /v2/crm/admin/seo-domains/relatives/{id}
```

# 2. Dẫn chứng nền tảng

## 2.1 Dẫn chứng từ tài liệu chiến lược SEO Domain

| Dẫn chứng | Ý nghĩa nghiệp vụ | Ứng dụng trong API hiện tại |
|---|---|---|
| Hệ thống SEO Domain Key Template quản lý các loại trang SEO-A, SEO-B, SEO-C; có automatic deduplication cho SEO-C, authority distribution qua link hierarchy cho SEO-B, rich schema generation cho SEO-A và scalable template system. | SEO Domain không chỉ là bảng URL; nó là lớp quản lý output SEO đã sinh ra hoặc cần publish. | `seo_domain` lưu URL/canonical/metadata/index/sitemap/robot/ref; `seo_links` tạo internal authority graph; `seo_relative` tạo quan hệ liên quan. |
| Domain Key có thể hiểu là định danh duy nhất cho mỗi trang SEO: loại template + ref ID + param hash. | Mỗi page SEO phải có danh tính rõ ràng để tránh duplicate. | Trong bản API hiện tại, `canonicalUrl` là khóa nghiệp vụ quan trọng để chống duplicate. |
| Parameter normalization giải quyết bài toán nhiều URL khác nhau cùng content rồi gom về một canonical. | Search/faceted URL phải được consolidate để Google không coi là duplicate content. | `CreateSeoDomain` và `UpdateSeoDomain` check unique `canonicalUrl`; metadata có thể chứa normalized params/canonical hash nếu tích hợp template engine sau. |
| Pattern guide/detail/search mô tả SEO-B làm guide, SEO-A làm detail, SEO-C làm search/faceted. | Phân loại SEO page cần đi kèm luồng intent của người dùng. | `classify`: A=10, B=20, C=30, D=40, N=50. |
| Pattern internal link mô tả guide page link xuống detail page để truyền authority. | SEO không chỉ là metadata; internal link graph ảnh hưởng khả năng crawl, topical relevance, page authority. | `SeoInternalLink` là cạnh có anchor `title` và URL `link` từ `parentSeoId` sang `childSeoId`. |

## 2.2 Dẫn chứng từ code implementation

| File | Dẫn chứng kỹ thuật | Ý nghĩa |
|---|---|---|
| `shared/protobuf/schema/crm/seo_domain.proto` | Proto dùng `camelCase`, timestamp trả `string`, enum truyền `uint32`. | Khớp style `ReportService` trong codebase. |
| `crm/internal/enums/seo.go` | `ESeoClassifyA=10`, `B=20`, `C=30`, `D=40`, `N=50`. | Enum Go nằm trong package `enums`, API truyền nhận `uint32`. |
| `crm/internal/domain/seo/seo_domain.go` | `SeoDomain`, `SeoInternalLink`, `SeoRelative` có `DeletedAt gorm.DeletedAt`. | Xóa mềm, không mất lịch sử ngay. |
| `crm/internal/usecase/seo_domain_usecase.go` | Validate `originUrl`, validate `classify`, validate JSON metadata, parse `siteMapLastedAt`, check unique canonical URL. | Business rules nằm ở usecase, không đẩy xuống handler. |
| `crm/infra/repository/postgre/seo_domain_postgres.go` | Query đều lọc `deleted_at IS NULL`, preload links/relatives, count total trước pagination. | Data layer đảm bảo soft-delete và list có total. |
| `crm/infra/mapper/seo_domain_mapper.go` | `CreatedAt`, `UpdatedAt`, `SiteMapLastedAt` format RFC3339 string. | Đúng yêu cầu timestamp là string trong API response. |

# 3. Tư tưởng nghiệp vụ tổng thể

## 3.1 SeoDomain là gì?

`SeoDomain` là bản ghi đại diện cho **một trang hoặc URL SEO có thể được quản lý, publish, index hoặc loại khỏi index**.

Nó lưu các nhóm dữ liệu:

| Nhóm | Field | Vai trò |
|---|---|---|
| URL | `originUrl`, `canonicalUrl` | Xác định URL gốc và URL chuẩn để chống duplicate. |
| Phân loại SEO | `classify` | Xác định page là detail, guide, search, share/viral hay non-SEO. |
| Indexing | `isSiteMap`, `siteMapLastedAt`, `isIndex`, `isRobot` | Điều khiển sitemap, index, robot/crawler behavior. |
| Metadata | `metadata` | Lưu JSON metadata: title, description, keywords, schema, OG, breadcrumbs, canonical hash, normalized params. |
| Mapping entity | `refType`, `refId` | Gắn SEO page với entity thật như parcel, project, report, map, document. |
| Link graph | `internalLinks`, `relatives` | Tạo cấu trúc liên kết nội bộ và quan hệ liên quan. |

## 3.2 SeoInternalLink là gì?

`SeoInternalLink` là cạnh liên kết có URL/anchor từ một SEO page cha sang một SEO page con.

Ví dụ:

```text
Guide Hà Nội /real-estate/hanoi
    └── link title: "Đất Ba Đình 100m²"
        link: /real-estate/hanoi/ba-dinh/a001
        childSeoId: detail page A001
```

Ý nghĩa nghiệp vụ:

- Dùng SEO-B guide page để dẫn crawler xuống SEO-A detail page.
- Tạo breadcrumb/navigation thực tế.
- Làm rõ topical relationship giữa page cha và page con.
- Tăng khả năng discover/crawl của detail pages.

## 3.3 SeoRelative là gì?

`SeoRelative` là quan hệ liên quan giữa hai SEO page, nhưng **không nhất thiết là một hyperlink hiển thị trực tiếp**.

Khác biệt:

| Loại | Có link URL? | Có anchor title? | Mục đích |
|---|---:|---:|---|
| `SeoInternalLink` | Có | Có | Link hiển thị/crawlable từ parent sang child. |
| `SeoRelative` | Không | Không | Quan hệ semantic/related/cùng chủ đề/cùng entity cluster. |

Ví dụ:

```text
Page "Đất bán Ba Đình"
    relative -> Page "Quy hoạch Ba Đình"
    relative -> Page "Dự án quanh Ba Đình"
    relative -> Page "Bản đồ quy hoạch Ba Đình"
```

# 4. Luồng kiến trúc request → response

```text
Client / Admin Panel
    ↓ HTTP transcoding / gRPC Gateway
Proto request: crmpb.*Request
    ↓
Handler: crm/internal/handler/seo_domain_handler.go
    - map proto → dto
    - gọi usecase
    - map domain → proto response
    ↓
DTO: crm/internal/dto/seo_domain_dto.go
    - chuẩn hóa shape request trong app layer
    ↓
Usecase: crm/internal/usecase/seo_domain_usecase.go
    - validate nghiệp vụ
    - check duplicate
    - parse time
    - verify parent/child tồn tại
    - build entity hoặc update map
    ↓
Repository interface: crm/internal/repo/seo_domain_repo.go
    ↓
Postgres implementation: crm/infra/repository/postgre/seo_domain_postgres.go
    - GORM Create/First/Find/Updates
    - soft delete bằng deleted_at
    - preload relation
    - pagination/count
    ↓
Mapper: crm/infra/mapper/seo_domain_mapper.go
    - domain → protobuf
    - time.Time → RFC3339 string
    ↓
Response
```

# 5. Enum classify và ý nghĩa SEO

API truyền `classify` bằng `uint32`; domain lưu `enums.ESeoClassify`.

| Giá trị | Enum Go | Tên nghiệp vụ | Ý nghĩa |
|---:|---|---|---|
| 10 | `ESeoClassifyA` | SEO-A / Detail | Trang chi tiết, giàu content, phục vụ long-tail intent. |
| 20 | `ESeoClassifyB` | SEO-B / Guide | Trang điều hướng, category, guide, hub page. |
| 30 | `ESeoClassifyC` | SEO-C / Search | Trang kết quả tìm kiếm/faceted, cần canonical/dedup mạnh. |
| 40 | `ESeoClassifyD` | Share/Viral/Preview | Link chia sẻ, preview, không phải SEO chính. |
| 50 | `ESeoClassifyN` | Non-SEO | Login, admin, 404, 500 hoặc page không phục vụ SEO. |

Business rule:

```go
if req.Classify != 0 && !enums.IsValidSeoClassify(req.Classify) {
    return nil, _errors.ReturnError(400, "classify không hợp lệ")
}
```

Nếu tạo mới mà `classify = 0`, usecase mặc định `ESeoClassifyN = 50` để tránh page chưa phân loại bị index nhầm.

# 6. API: CreateSeoDomain

## 6.1 Contract

```text
POST /v2/crm/admin/seo-domains
```

```protobuf
message CreateSeoDomainRequest {
  string originUrl = 1;
  optional string canonicalUrl = 2;
  uint32 classify = 3;
  optional string siteMapLastedAt = 4;
  bool isSiteMap = 5;
  bool isIndex = 6;
  optional bool isRobot = 7;
  optional string metadata = 8;
  string note = 9;
  uint32 refType = 10;
  optional uint64 refId = 11;
}
```

## 6.2 Mục đích nghiệp vụ

Tạo một SEO page record mới để hệ thống có thể quản lý URL, canonical, metadata, sitemap/index/robot và mapping tới entity thật.

Ví dụ dùng cho:

- Tạo page detail cho thửa đất: SEO-A.
- Tạo guide page cho quận/huyện: SEO-B.
- Tạo search/faceted canonical page: SEO-C.
- Tạo page non-SEO và đánh dấu không index: SEO-N.

## 6.3 Cơ chế xử lý

```text
Bước 1: Handler nhận crmpb.CreateSeoDomainRequest
Bước 2: Handler map sang dto.SeoDomainCreateRequest
Bước 3: Usecase validate originUrl không rỗng
Bước 4: Usecase validate classify nếu classify != 0
Bước 5: Usecase validate metadata nếu metadata không rỗng thì phải là JSON hợp lệ
Bước 6: Nếu canonicalUrl không truyền, canonicalUrl = originUrl
Bước 7: Check canonicalUrl đã tồn tại chưa
Bước 8: Parse siteMapLastedAt nếu có
Bước 9: Build domain.SeoDomain
Bước 10: Repository Create
Bước 11: Repository GetByID để preload relation và trả bản ghi mới nhất
Bước 12: Mapper format time thành string RFC3339 và trả SeoDomainResponse
```

## 6.4 Business validations

| Rule | Lý do | Error |
|---|---|---|
| `originUrl` không được rỗng | Không thể quản lý SEO page nếu không có URL gốc. | `400 originUrl không được để trống` |
| `classify` phải thuộc 10/20/30/40/50 | Tránh lưu page sai loại SEO. | `400 classify không hợp lệ` |
| `metadata` phải là JSON hợp lệ | Metadata lưu `jsonb`, sai JSON sẽ gây lỗi DB/render. | `400 metadata phải là JSON hợp lệ` |
| `canonicalUrl` unique | Chống duplicate canonical, tránh nhiều bản ghi cạnh tranh cùng URL chuẩn. | `400 canonicalUrl đã tồn tại` |
| `siteMapLastedAt` parse được | Đồng bộ sitemap cần timestamp chuẩn. | `400 siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD` |

## 6.5 Ví dụ request

```json
{
  "originUrl": "/real-estate/hanoi/ba-dinh/a001",
  "canonicalUrl": "/real-estate/hanoi/ba-dinh/a001",
  "classify": 10,
  "siteMapLastedAt": "2026-06-11T00:00:00Z",
  "isSiteMap": true,
  "isIndex": true,
  "isRobot": true,
  "metadata": "{\"title\":\"Thửa đất A001 - 100m² ở Ba Đình\",\"description\":\"Mua thửa đất A001 tại Ba Đình, Hà Nội\",\"schemaType\":\"RealEstateListing\"}",
  "note": "SEO-A detail page generated from parcel",
  "refType": 3,
  "refId": 1001
}
```

## 6.6 Ví dụ response

```json
{
  "id": 1,
  "originUrl": "/real-estate/hanoi/ba-dinh/a001",
  "canonicalUrl": "/real-estate/hanoi/ba-dinh/a001",
  "classify": 10,
  "siteMapLastedAt": "2026-06-11T00:00:00Z",
  "isSiteMap": true,
  "isIndex": true,
  "isRobot": true,
  "metadata": "{\"title\":\"Thửa đất A001 - 100m² ở Ba Đình\"}",
  "note": "SEO-A detail page generated from parcel",
  "refType": 3,
  "refId": 1001,
  "internalLinks": [],
  "relatives": [],
  "createdAt": "2026-06-11T10:00:00Z",
  "updatedAt": "2026-06-11T10:00:00Z"
}
```

## 6.7 Dẫn chứng nghiệp vụ

- SEO-A detail page cần rich metadata/schema để phục vụ long-tail search.
- SEO-C/faceted page cần canonical unique để gom nhiều biến thể URL về một URL chuẩn.
- API enforce unique `canonicalUrl`, vì nếu cho phép trùng canonical thì cùng một canonical có thể bị nhiều record quản lý, gây lỗi sitemap, internal link, update và thống kê.

## 6.8 Acceptance criteria

- Tạo page thiếu `originUrl` phải lỗi 400.
- Tạo page có `metadata = "abc"` phải lỗi 400.
- Tạo page có `classify = 99` phải lỗi 400.
- Tạo page không truyền `canonicalUrl` thì response `canonicalUrl = originUrl`.
- Tạo page trùng `canonicalUrl` với page chưa xóa phải lỗi 400.
- Response trả `createdAt`, `updatedAt`, `siteMapLastedAt` dạng string.

# 7. API: GetSeoDomain

## 7.1 Contract

```text
GET /v2/crm/admin/seo-domains/{id}
```

Input dùng `sharepb.IdRequest`.

## 7.2 Mục đích nghiệp vụ

Lấy chi tiết một SEO page, kèm internal links và relatives để admin/debug thấy:

- Page này đang canonical về đâu.
- Page có được index/sitemap/robot hay không.
- Page liên kết tới các page nào.
- Page có quan hệ semantic với page nào.

## 7.3 Cơ chế xử lý

```text
Handler nhận id
    ↓
Usecase gọi GetSeoDomain(id)
    ↓
Repo query seo_domain WHERE id = ? AND deleted_at IS NULL
    ↓
Repo preload InternalLinks, InternalLinks.ChildSeo, Relatives, Relatives.ChildSeo
    ↓
Nếu không có: trả 404 seo domain không tồn tại
    ↓
Mapper domain → SeoDomainResponse
```

## 7.4 Dẫn chứng kỹ thuật

Repository `GetByID` preload:

```go
Preload("InternalLinks", "deleted_at IS NULL").
Preload("InternalLinks.ChildSeo").
Preload("Relatives", "deleted_at IS NULL").
Preload("Relatives.ChildSeo")
```

Điều này giúp response không chỉ có bản ghi SEO page, mà có luôn graph cấp 1 phục vụ quản trị SEO.

## 7.5 Acceptance criteria

- ID không tồn tại hoặc đã soft-delete trả 404.
- Response có `internalLinks` và `relatives` nếu đã được tạo.
- Không trả relation đã `deleted_at`.

# 8. API: GetSeoDomainList

## 8.1 Contract

```text
GET /v2/crm/admin/seo-domains
```

Query:

```protobuf
message GetSeoDomainListRequest {
  optional string keyword = 1;
  optional uint32 classify = 2;
  optional bool isSiteMap = 3;
  optional bool isIndex = 4;
  optional bool isRobot = 5;
  optional uint32 refType = 6;
  optional uint64 refId = 7;
  uint32 page = 8;
  uint32 size = 9;
}
```

## 8.2 Mục đích nghiệp vụ

Cho admin SEO/devops lọc danh sách SEO pages theo trạng thái publish/index/entity để kiểm soát coverage.

Ví dụ câu hỏi nghiệp vụ:

- Có bao nhiêu page SEO-A đang index?
- Các URL nào được đưa vào sitemap?
- Parcel `refId=1001` đang có SEO page nào?
- URL chứa `ba-dinh` có bao nhiêu record?
- Có page nào `isRobot=false` nhưng `isIndex=true` không?

## 8.3 Cơ chế xử lý

```text
Handler nhận query params
    ↓
Map vào dto.SeoDomainListRequest + Pagable
    ↓
Usecase gọi repo.GetList
    ↓
Repo build query base: seo_domain WHERE deleted_at IS NULL
    ↓
Apply filters nếu có
    ↓
Count total
    ↓
Order updated_at DESC
    ↓
Offset/Limit theo page/size
    ↓
Preload InternalLinks và Relatives
    ↓
Return data + total
```

## 8.4 Filter behavior

| Filter | SQL behavior | Ý nghĩa |
|---|---|---|
| `keyword` | `origin_url ILIKE ? OR canonical_url ILIKE ? OR note ILIKE ?` | Search URL/note nhanh. |
| `classify` | `classify = ?` | Lọc SEO-A/B/C/D/N. |
| `isSiteMap` | `is_site_map = ?` | Lọc page có/không có trong sitemap. |
| `isIndex` | `is_index = ?` | Lọc page có/không cho index. |
| `isRobot` | `is_robot = ?` | Lọc trạng thái robot/crawler. |
| `refType` | `ref_type = ?` | Lọc theo loại entity. |
| `refId` | `ref_id = ?` | Lọc theo entity cụ thể. |
| `page`, `size` | offset/limit | Pagination. |

## 8.5 Acceptance criteria

- Không truyền filter: trả danh sách page chưa xóa, sort `updated_at DESC`.
- Truyền `keyword=ba-dinh`: trả page có origin/canonical/note chứa `ba-dinh`.
- Truyền `classify=10&isIndex=true`: trả SEO-A detail đang index.
- `total` là tổng số record match filter trước pagination.

# 9. API: UpdateSeoDomain

## 9.1 Contract

```text
PUT /v2/crm/admin/seo-domains/{id}
```

```protobuf
message UpdateSeoDomainRequest {
  uint64 id = 1;
  optional string originUrl = 2;
  optional string canonicalUrl = 3;
  optional uint32 classify = 4;
  optional string siteMapLastedAt = 5;
  optional bool isSiteMap = 6;
  optional bool isIndex = 7;
  optional bool isRobot = 8;
  optional string metadata = 9;
  optional string note = 10;
  optional uint32 refType = 11;
  optional uint64 refId = 12;
}
```

## 9.2 Mục đích nghiệp vụ

Cập nhật SEO page sau khi có thay đổi chiến lược SEO hoặc entity gốc thay đổi.

Ví dụ:

- Bật/tắt `isIndex` khi page không còn đủ chất lượng.
- Bật/tắt `isSiteMap` khi muốn submit/remove khỏi sitemap.
- Thay metadata title/description/schema sau khi template được tối ưu.
- Đổi canonical khi faceted search được normalize lại.
- Gắn lại `refType/refId` nếu mapping entity bị sai.

## 9.3 Cơ chế xử lý

```text
Bước 1: Usecase GetSeoDomain(id) để đảm bảo record tồn tại
Bước 2: Tạo map updates rỗng
Bước 3: Với từng optional field, chỉ update nếu field được truyền
Bước 4: Validate originUrl nếu truyền thì không rỗng
Bước 5: Validate canonicalUrl nếu truyền:
        - trim space
        - nếu rỗng thì fallback existing.OriginURL
        - check unique canonical, bỏ qua chính id hiện tại
Bước 6: Validate classify nếu truyền
Bước 7: Parse siteMapLastedAt nếu truyền:
        - string rỗng → set DB null
        - RFC3339/YYYY-MM-DD/YYYY-MM-DD HH:mm:ss → time
Bước 8: Validate metadata JSON nếu truyền
Bước 9: Nếu updates rỗng → trả existing, không gọi DB update
Bước 10: Repo Updates(updates)
Bước 11: GetSeoDomain lại để trả dữ liệu mới nhất
```

## 9.4 Lý do dùng optional field

Update request dùng `optional` để phân biệt:

| Tình huống | Ý nghĩa |
|---|---|
| Field không gửi | Không thay đổi DB. |
| Field gửi giá trị mới | Update DB. |
| `siteMapLastedAt = ""` | Clear sitemap time về null. |
| `canonicalUrl = ""` | Fallback về `existing.OriginURL`. |
| `metadata = ""` | Cho phép clear metadata. |

## 9.5 Acceptance criteria

- Update ID không tồn tại trả 404.
- Update chỉ `note` không ảnh hưởng các field khác.
- Update `canonicalUrl` trùng record khác trả 400.
- Update `metadata` invalid JSON trả 400.
- Update không truyền field nào trả existing, không lỗi.

# 10. API: DeleteSeoDomain

## 10.1 Contract

```text
DELETE /v2/crm/admin/seo-domains/{id}
```

Input dùng `sharepb.IdRequest`.

## 10.2 Mục đích nghiệp vụ

Xóa mềm một SEO page khi page không còn hợp lệ hoặc không muốn quản lý nữa.

Không hard-delete vì:

- Cần giữ audit/history.
- Tránh mất dữ liệu phân tích ngay lập tức.
- Có thể cần biết URL từng tồn tại để redirect/410/sitemap cleanup.

## 10.3 Cơ chế xử lý

```text
Usecase GetSeoDomain(id)
    ↓
Nếu không tồn tại: 404
    ↓
Repo Delete
    ↓
GORM update deleted_at = NOW()
    ↓
Return SubmitResponse { id, message }
```

## 10.4 Side effect cần lưu ý

Bản hiện tại chỉ soft-delete `seo_domain`, chưa cascade mềm `seo_links` và `seo_relative`. Vì các API list relation đều lọc `deleted_at IS NULL` ở relation table, relation cũ vẫn tồn tại nếu chưa bị xóa riêng. Khi mở rộng production nên cân nhắc:

- Transaction soft-delete domain + links + relatives.
- Hoặc background cleanup orphaned links.
- Hoặc DB constraint/foreign key `ON DELETE` phù hợp nếu hard-delete.

## 10.5 Acceptance criteria

- Delete ID không tồn tại trả 404.
- Delete thành công trả message `Seo domain deleted successfully`.
- Get lại ID đã xóa trả 404.
- List không trả record đã xóa.

# 11. API: CreateSeoInternalLink

## 11.1 Contract

```text
POST /v2/crm/admin/seo-domains/internal-links
```

```protobuf
message CreateSeoInternalLinkRequest {
  string link = 1;
  string title = 2;
  uint64 parentSeoId = 3;
  uint64 childSeoId = 4;
}
```

## 11.2 Mục đích nghiệp vụ

Tạo liên kết nội bộ từ page cha sang page con để xây dựng SEO hierarchy.

Ví dụ:

```text
SEO-B guide: /real-estate/hanoi/ba-dinh
    link → SEO-A detail: /real-estate/hanoi/ba-dinh/a001
```

Trong logic SEO, đây là cơ chế để:

- Page guide/category dẫn crawler tới detail pages.
- Tăng khả năng discover của detail page.
- Tạo anchor text có ngữ nghĩa.
- Hình thành cấu trúc Homepage → Guide → Detail.

## 11.3 Cơ chế xử lý

```text
Handler map request → SeoInternalLinkCreateRequest
    ↓
Usecase validate parentSeoId và childSeoId bắt buộc
    ↓
Validate parentSeoId != childSeoId
    ↓
GetSeoDomain(parentSeoId) để đảm bảo parent tồn tại
    ↓
GetSeoDomain(childSeoId) để đảm bảo child tồn tại
    ↓
Repo Exist(parentSeoId, childSeoId, nil) để chống duplicate edge
    ↓
Build SeoInternalLink { link, title, parentSeoId, childSeoId }
    ↓
Repo Create
    ↓
Repo GetByID preload ParentSeo và ChildSeo
    ↓
Mapper trả SeoInternalLinkResponse
```

## 11.4 Business validations

| Rule | Lý do | Error |
|---|---|---|
| `parentSeoId`, `childSeoId` bắt buộc | Link graph cần đủ 2 node. | `400 parentSeoId và childSeoId là bắt buộc` |
| Parent không được trùng child | Tránh self-loop vô nghĩa. | `400 parentSeoId không được trùng childSeoId` |
| Parent phải tồn tại | Không tạo edge mồ côi. | 404 từ `GetSeoDomain` |
| Child phải tồn tại | Không tạo edge mồ côi. | 404 từ `GetSeoDomain` |
| Không trùng cặp parent-child | Tránh duplicate internal link edge. | `400 internal link đã tồn tại` |

## 11.5 Ví dụ request

```json
{
  "link": "/real-estate/hanoi/ba-dinh/a001",
  "title": "Thửa đất A001 - 100m² tại Ba Đình",
  "parentSeoId": 10,
  "childSeoId": 1
}
```

## 11.6 Ví dụ response

```json
{
  "id": 100,
  "link": "/real-estate/hanoi/ba-dinh/a001",
  "title": "Thửa đất A001 - 100m² tại Ba Đình",
  "parentSeoId": 10,
  "childSeoId": 1,
  "parentSeo": {
    "id": 10,
    "originUrl": "/real-estate/hanoi/ba-dinh",
    "canonicalUrl": "/real-estate/hanoi/ba-dinh",
    "classify": 20
  },
  "childSeo": {
    "id": 1,
    "originUrl": "/real-estate/hanoi/ba-dinh/a001",
    "canonicalUrl": "/real-estate/hanoi/ba-dinh/a001",
    "classify": 10
  },
  "createdAt": "2026-06-11T10:00:00Z",
  "updatedAt": "2026-06-11T10:00:00Z"
}
```

## 11.7 Acceptance criteria

- Không cho tạo parent = child.
- Không cho tạo link nếu parent hoặc child đã bị xóa.
- Không cho tạo duplicate `parentSeoId-childSeoId`.
- Response có `parentSeo` và `childSeo` brief.

# 12. API: GetSeoInternalLinkList

## 12.1 Contract

```text
GET /v2/crm/admin/seo-domains/internal-links
```

Query:

```protobuf
message GetSeoInternalLinkListRequest {
  optional uint64 parentSeoId = 1;
  optional uint64 childSeoId = 2;
  optional string keyword = 3;
  uint32 page = 4;
  uint32 size = 5;
}
```

## 12.2 Mục đích nghiệp vụ

Audit và debug internal linking graph.

Ví dụ câu hỏi nghiệp vụ:

- Guide page `10` đang link xuống các detail page nào?
- Detail page `1` đang được page nào link tới?
- Anchor text nào chứa `Ba Đình`?
- Có bao nhiêu link đã tạo trong cluster này?

## 12.3 Cơ chế xử lý

```text
Repo query seo_links WHERE deleted_at IS NULL
    ↓
Nếu parentSeoId có: filter parent_seo_id
Nếu childSeoId có: filter child_seo_id
Nếu keyword có: filter link ILIKE hoặc title ILIKE
    ↓
Count total
    ↓
Preload ParentSeo + ChildSeo
    ↓
Order updated_at DESC
    ↓
Offset/Limit
```

## 12.4 Acceptance criteria

- `parentSeoId=10` trả outbound links của page 10.
- `childSeoId=1` trả inbound links tới page 1.
- `keyword=ba-dinh` search được trong URL hoặc anchor title.
- Không trả link đã soft-delete.

# 13. API: UpdateSeoInternalLink

## 13.1 Contract

```text
PUT /v2/crm/admin/seo-domains/internal-links/{id}
```

```protobuf
message UpdateSeoInternalLinkRequest {
  uint64 id = 1;
  optional string link = 2;
  optional string title = 3;
  optional uint64 parentSeoId = 4;
  optional uint64 childSeoId = 5;
}
```

## 13.2 Mục đích nghiệp vụ

Cập nhật URL/anchor hoặc đổi cạnh link khi cấu trúc SEO hierarchy thay đổi.

Ví dụ:

- Đổi anchor title để tăng relevance.
- Đổi link path sau khi canonical URL thay đổi.
- Chuyển child page từ guide cũ sang guide mới.

## 13.3 Cơ chế xử lý

```text
Usecase get internal link by id
    ↓
Nếu không tồn tại: 404
    ↓
Lấy parentID/childID hiện tại
    ↓
Nếu request truyền parentSeoId hoặc childSeoId thì thay tạm để validate
    ↓
Validate parentID/childID bắt buộc và không trùng nhau
    ↓
Validate parent/child SeoDomain tồn tại
    ↓
Check duplicate edge, bỏ qua chính id hiện tại
    ↓
Build update map cho field được truyền
    ↓
Nếu update map rỗng: trả existing
    ↓
Repo Updates
    ↓
Repo GetByID trả bản mới nhất
```

## 13.4 Acceptance criteria

- Update ID không tồn tại trả 404.
- Update đổi parent/child sang cặp đã tồn tại trả 400.
- Update parent = child trả 400.
- Update chỉ title không ảnh hưởng parent/child.

# 14. API: DeleteSeoInternalLink

## 14.1 Contract

```text
DELETE /v2/crm/admin/seo-domains/internal-links/{id}
```

## 14.2 Mục đích nghiệp vụ

Xóa mềm một internal link khi:

- Anchor/link không còn đúng canonical.
- Muốn bỏ đường dẫn crawl tới page con.
- Page con không còn nên được truyền authority.
- Cần tái cấu trúc hierarchy.

## 14.3 Cơ chế xử lý

```text
Usecase GetByID(id)
    ↓
Nếu không tồn tại: 404
    ↓
Repo Delete: update deleted_at = NOW()
    ↓
Return SubmitResponse
```

## 14.4 Acceptance criteria

- Delete link không tồn tại trả 404.
- Delete thành công không hiện trong list.
- Parent/child SEO pages không bị xóa theo.

# 15. API: CreateSeoRelative

## 15.1 Contract

```text
POST /v2/crm/admin/seo-domains/relatives
```

```protobuf
message CreateSeoRelativeRequest {
  uint64 parentSeoId = 1;
  uint64 childSeoId = 2;
}
```

## 15.2 Mục đích nghiệp vụ

Tạo quan hệ liên quan giữa hai SEO pages để dùng cho:

- Module “bài/trang liên quan”.
- Gợi ý internal link tự động về sau.
- Semantic cluster cho cùng khu vực/chủ đề/entity.
- Phân tích graph giữa các page không nhất thiết có link trực tiếp.

Ví dụ:

```text
SEO-B: /real-estate/hanoi/ba-dinh
    relative → SEO-B: /planning/hanoi/ba-dinh
    relative → SEO-C: /real-estate/search/price-1-3b/location-ba-dinh
```

## 15.3 Cơ chế xử lý

```text
Validate parentSeoId và childSeoId bắt buộc
    ↓
Validate parentSeoId != childSeoId
    ↓
Validate parent SEO domain tồn tại
    ↓
Validate child SEO domain tồn tại
    ↓
Check duplicate relation parent-child
    ↓
Create SeoRelative
    ↓
Repo GetByID preload ParentSeo và ChildSeo
```

## 15.4 Khác biệt với internal link

| Câu hỏi | SeoInternalLink | SeoRelative |
|---|---|---|
| Có URL link cụ thể không? | Có `link`. | Không. |
| Có anchor title không? | Có `title`. | Không. |
| Dùng để render HTML anchor? | Có. | Không bắt buộc. |
| Dùng để semantic clustering? | Có thể. | Chính yếu. |
| Duplicate rule | parent-child unique, ignore id khi update. | parent-child unique. |

## 15.5 Acceptance criteria

- Không cho tạo parent = child.
- Không cho tạo duplicate parent-child.
- Không cho tạo nếu parent/child không tồn tại.
- Response có parentSeo và childSeo brief.

# 16. API: GetSeoRelativeList

## 16.1 Contract

```text
GET /v2/crm/admin/seo-domains/relatives
```

Query:

```protobuf
message GetSeoRelativeListRequest {
  optional uint64 parentSeoId = 1;
  optional uint64 childSeoId = 2;
  uint32 page = 3;
  uint32 size = 4;
}
```

## 16.2 Mục đích nghiệp vụ

Xem graph related pages để phục vụ SEO cluster và recommendation.

Ví dụ:

- Page guide Ba Đình đang liên quan tới những page nào?
- Page detail A001 được page nào xem là relative?
- Có duplicate cluster không?

## 16.3 Cơ chế xử lý

```text
Repo query seo_relative WHERE deleted_at IS NULL
    ↓
Apply parentSeoId/childSeoId filter nếu có
    ↓
Count total
    ↓
Preload ParentSeo + ChildSeo
    ↓
Order updated_at DESC
    ↓
Offset/Limit
```

## 16.4 Acceptance criteria

- `parentSeoId=10` trả các relation đi ra từ page 10.
- `childSeoId=1` trả các relation trỏ tới page 1.
- Không trả relation đã soft-delete.

# 17. API: DeleteSeoRelative

## 17.1 Contract

```text
DELETE /v2/crm/admin/seo-domains/relatives/{id}
```

## 17.2 Mục đích nghiệp vụ

Xóa mềm một relation khi semantic cluster không còn đúng.

Ví dụ:

- Page đã đổi khu vực/chủ đề.
- Search canonical mới thay thế page cũ.
- Entity source bị xóa hoặc ngưng publish.

## 17.3 Cơ chế xử lý

```text
Usecase GetByID(id)
    ↓
Nếu không tồn tại: 404
    ↓
Repo Delete: update deleted_at = NOW()
    ↓
Return SubmitResponse
```

## 17.4 Acceptance criteria

- Delete relation không tồn tại trả 404.
- Delete thành công không hiện trong list.
- Parent/child SEO pages không bị xóa theo.

# 18. Metadata JSON - cấu trúc khuyến nghị

Field `metadata` là `*string` nhưng được validate `json.Valid` và lưu DB `jsonb`. Vì vậy frontend/admin panel nên gửi JSON string hợp lệ.

## 18.1 Metadata cho SEO-A detail

```json
{
  "title": "Thửa đất A001 - 100m² ở Ba Đình, Hà Nội",
  "description": "Mua thửa đất A001 diện tích 100m² tại Ba Đình, Hà Nội. Giá 2.5 tỷ.",
  "keywords": ["đất bán ba đình", "thửa đất A001", "bất động sản hà nội"],
  "schemaType": "RealEstateListing",
  "schema": {
    "@context": "https://schema.org",
    "@type": "RealEstateListing",
    "name": "Thửa đất A001",
    "address": "Ba Đình, Hà Nội"
  },
  "breadcrumb": [
    {"text": "Home", "url": "/"},
    {"text": "Hà Nội", "url": "/real-estate/hanoi"},
    {"text": "Ba Đình", "url": "/real-estate/hanoi/ba-dinh"}
  ]
}
```

## 18.2 Metadata cho SEO-B guide

```json
{
  "title": "Đất bán ở Ba Đình, Hà Nội | Hướng dẫn thị trường",
  "description": "Tổng hợp thông tin thị trường đất Ba Đình, giá, khu vực nổi bật và danh sách bất động sản liên quan.",
  "keywords": ["đất bán ba đình", "bất động sản ba đình", "nhà đất hà nội"],
  "schemaType": "Place",
  "internalLinkPolicy": {
    "linkToClassify": [10],
    "maxLinks": 50,
    "anchorPattern": "{name} - {area}m²"
  }
}
```

## 18.3 Metadata cho SEO-C faceted search

```json
{
  "title": "Tìm đất bán ở Hà Nội | Giá 1-3 tỷ",
  "description": "Danh sách đất bán Hà Nội trong khoảng giá 1-3 tỷ, đã chuẩn hóa canonical để tránh trùng lặp.",
  "keywords": ["đất bán hà nội 1-3 tỷ", "tìm đất hà nội"],
  "schemaType": "SearchResultsPage",
  "normalizedParams": {
    "keyword": "đất bán",
    "location": "hanoi",
    "price": "1-3b"
  },
  "canonicalHash": "md5-or-sha256-hash"
}
```

# 19. Sitemap / Index / Robot policy

## 19.1 Field meaning

| Field | Ý nghĩa | Rule khuyến nghị |
|---|---|---|
| `isSiteMap` | Có đưa vào sitemap không | Chỉ true với canonical page có chất lượng. |
| `siteMapLastedAt` | Lần cuối cập nhật sitemap | Set khi page được generate/update và đã sync sitemap. |
| `isIndex` | Có cho index không | True với SEO-A/B/C canonical; false với duplicate, admin, preview. |
| `isRobot` | Robot/crawler có được phép crawl không | Pointer để phân biệt unset và false. |

## 19.2 Policy matrix

| Page type | classify | isSiteMap | isIndex | isRobot | Ghi chú |
|---|---:|---:|---:|---:|---|
| Detail parcel canonical | 10 | true | true | true | SEO-A, rich metadata/schema. |
| District guide | 20 | true | true | true | SEO-B, link hub. |
| Search canonical | 30 | true | true | true | SEO-C, only canonical after dedup. |
| Share preview | 40 | false | false | true | Không index nhưng có thể crawl preview. |
| Admin/login/404 | 50 | false | false | false | Non-SEO. |

# 20. Error handling matrix

| Tình huống | API | Error |
|---|---|---|
| Tạo domain thiếu `originUrl` | CreateSeoDomain | 400 `originUrl không được để trống` |
| `classify` ngoài 10/20/30/40/50 | Create/UpdateSeoDomain | 400 `classify không hợp lệ` |
| `metadata` invalid JSON | Create/UpdateSeoDomain | 400 `metadata phải là JSON hợp lệ` |
| `canonicalUrl` trùng | Create/UpdateSeoDomain | 400 `canonicalUrl đã tồn tại` |
| `siteMapLastedAt` sai format | Create/UpdateSeoDomain | 400 `siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD` |
| Domain không tồn tại | Get/Update/DeleteSeoDomain | 404 `seo domain không tồn tại` |
| Parent/child thiếu | Create/UpdateInternalLink, CreateRelative | 400 `parentSeoId và childSeoId là bắt buộc` |
| Parent trùng child | Create/UpdateInternalLink, CreateRelative | 400 `parentSeoId không được trùng childSeoId` |
| Duplicate internal link | Create/UpdateInternalLink | 400 `internal link đã tồn tại` |
| Duplicate relative | CreateRelative | 400 `seo relative đã tồn tại` |
| Internal link không tồn tại | Update/DeleteInternalLink | 404 `seo internal link không tồn tại` |
| Relative không tồn tại | DeleteRelative | 404 `seo relative không tồn tại` |

# 21. Test scenarios

## 21.1 CreateSeoDomain tests

```text
TC01: Create SEO-A valid
Input classify=10, valid metadata JSON, originUrl/canonicalUrl unique
Expected 200, response classify=10, timestamps string

TC02: Create missing originUrl
Input originUrl=""
Expected 400

TC03: Create invalid metadata
Input metadata="not-json"
Expected 400

TC04: Create duplicate canonical
Precondition existing canonicalUrl=/a
Input canonicalUrl=/a
Expected 400

TC05: Create without canonicalUrl
Input originUrl=/a, canonicalUrl unset
Expected response canonicalUrl=/a
```

## 21.2 UpdateSeoDomain tests

```text
TC06: Update note only
Expected only note changed

TC07: Update classify invalid
Input classify=999
Expected 400

TC08: Update siteMapLastedAt empty
Input siteMapLastedAt=""
Expected DB field null

TC09: Update canonical to duplicate
Expected 400

TC10: Update no fields
Expected return existing, no error
```

## 21.3 Internal link tests

```text
TC11: Create valid guide → detail link
Parent classify=20, child classify=10
Expected 200, response has parentSeo and childSeo

TC12: Create self-loop
parentSeoId = childSeoId
Expected 400

TC13: Create duplicate link
Same parent-child already exists
Expected 400

TC14: List outbound links
Query parentSeoId=guideId
Expected all child links

TC15: Delete link
Expected soft-delete, list no longer returns it
```

## 21.4 Relative tests

```text
TC16: Create valid relative
Expected 200

TC17: Create duplicate relative
Expected 400

TC18: List by childSeoId
Expected inbound related pages

TC19: Delete relative
Expected soft-delete
```

# 22. Operational checklist

## 22.1 Khi tạo SEO page mới

- [ ] `originUrl` đúng URL thật.
- [ ] `canonicalUrl` là URL chuẩn, không chứa tracking params.
- [ ] `classify` đúng A/B/C/D/N.
- [ ] `metadata` là JSON hợp lệ.
- [ ] `isSiteMap`, `isIndex`, `isRobot` phù hợp page type.
- [ ] `refType/refId` map đúng entity nguồn.
- [ ] Nếu là guide page, tạo internal links xuống detail pages.
- [ ] Nếu là cluster page, tạo relatives tới page liên quan.

## 22.2 Khi update SEO page

- [ ] Không đổi canonical tùy tiện nếu page đã index.
- [ ] Nếu đổi canonical, cập nhật internal links trỏ tới URL mới.
- [ ] Nếu tắt `isIndex`, cân nhắc remove khỏi sitemap.
- [ ] Nếu metadata thay đổi từ template, validate JSON và render preview ở frontend/admin.

## 22.3 Khi delete SEO page

- [ ] Kiểm tra inbound internal links tới page này.
- [ ] Kiểm tra relatives liên quan.
- [ ] Cập nhật sitemap nếu page trước đó `isSiteMap=true`.
- [ ] Cân nhắc redirect/410 ở web layer nếu URL public đã từng index.

# 23. Khuyến nghị mở rộng sau bản CRUD hiện tại

Bản CRUD hiện tại là nền móng quản lý SEO Domain. Để đạt cơ chế SEO automation đầy đủ như Domain Key Template System, nên bổ sung các module sau:

| Module | Mục đích | API gợi ý |
|---|---|---|
| Template management | Quản lý rule generate metadata/canonical/schema. | `/seo-templates` |
| Generate page | Sinh `SeoDomain` từ template + entity data. | `/seo-pages/generate` |
| Batch generation | Sinh hàng loạt SEO-A/B/C pages. | `/seo-pages/generate-batch` |
| Dedup check | Normalize params và check canonical tồn tại. | `/seo-pages/check-dedup` |
| Preview metadata | Render title/desc/schema trước khi lưu. | `/seo-templates/{id}/preview` |
| Validate metadata | Check length, JSON schema, canonical, sitemap policy. | `/seo-templates/{id}/validate` |
| Metrics | Đo index/rank/CTR/click/impression. | `/seo-metrics/*` |

# 24. Kết luận

Nhóm API `SeoDomainService` hiện tại đóng vai trò **SEO registry layer** trong hệ thống CRM:

1. `SeoDomain` quản lý từng SEO page/canonical/metadata/indexing state.
2. `SeoInternalLink` quản lý crawlable internal link từ page cha sang page con.
3. `SeoRelative` quản lý semantic relationship giữa các page.
4. Usecase chịu trách nhiệm validate nghiệp vụ: canonical unique, enum hợp lệ, metadata JSON hợp lệ, parent/child tồn tại, chống duplicate edge.
5. Repository chịu trách nhiệm query DB, soft-delete, preload relation và pagination.
6. Mapper đảm bảo API response theo convention: camelCase + enum uint32 + timestamp string.

Khi kết hợp với template generation/dedup engine, lớp API này sẽ trở thành nơi lưu output SEO chuẩn hóa, giúp hệ thống đi từ metadata thủ công sang quản lý SEO có cấu trúc, có canonical, có sitemap/index policy và có link graph rõ ràng.
