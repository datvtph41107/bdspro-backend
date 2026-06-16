# CRM SEO Enums va RefType

Ngay cap nhat: 2026-06-15

Tai lieu nay tong hop hien trang xu ly `enums`, `seo refType`, `resolverKey`, `sourceService` va muc do hoan thien cua backend SEO ref-aware trong CRM service.

## Ket luan nhanh

He thong da co nen tang dung huong cho SEO ref-aware:

- CRM service da co API enum pho thong: `GET /v2/enums/crm/{name}`.
- CRM SEO da co schema API rieng cho ref registry:
  - `GET /v2/crm/admin/seo/ref-types`
  - `GET /v2/crm/admin/seo/ref-values`
  - `GET /v2/crm/admin/seo/ref-values/{refType}/{refId}`
- Backend da co `SeoRefTypeConfig` lam registry metadata cho `refType`.
- Frontend co the doc danh sach `refType` SEO tu CRM thay vi hard-code service nguon.

Nhung he thong CHUA day du production cho search/resolve source:

- `GetSeoRefTypes` da tra du metadata.
- `SearchSeoRefValues` hien moi validate registry va tra `501 resolver SEO chua duoc trien khai`.
- `ResolveSeoRefValue` hien moi validate registry/refId va tra `501 resolver SEO chua duoc trien khai`.
- Chua co resolver thuc te goi BDSPro/TQD/CRM source service.
- Chua co `ref_source` trong model/migration de phan biet cung `refType/refId` nhung khac source.
- Chua co snapshot/hash/source status luu vao `seo_domain`.

Trang thai: API contract va registry da san sang; resolver/source snapshot/sync flow van can lam tiep.

## 1. Enum API hien tai

### Endpoint

```txt
GET /v2/enums/crm/{name}
```

Proto:

```txt
shared/protobuf/schema/crm/enum.proto
```

Handler:

```txt
crm-service/infra/handler/enum_handler.go
```

Usecase:

```txt
crm-service/internal/usecase/enum_usecase.go
```

### Muc dich

Enum API hien tai dung cho enum pho thong cua CRM, dang response theo format:

```json
{
  "data": [
    {
      "id": 1,
      "name": "..."
    }
  ],
  "total": 1
}
```

### Cac enum name dang ho tro

| name | Nguon map trong code | Ghi chu |
|---|---|---|
| `contact_tag` | `enums.TagContactMap` | Tag contact |
| `statusContact` | `enums.StatusContactMap` | Trang thai contact |
| `visibility` | `enums.VisibilityMap` | Visibility |
| `priority` | `enums.PriorityMap` | Priority |
| `ownerOf` | `enums.OwnerOfNames` | Owner type |
| `step` | `enums.StepMap` | Step |
| `sourceLead` | `enums.SourceLeadMap` | Lead source |
| `friendStatus` | `enums.FriendStatusMap` | Friend status |
| `condition` | `enums.ConditionMap` | Rule condition |
| `ruleThen` | `enums.RuleThenMap` | Rule action |
| `targetType` | `enums.TargetTypeMap` | Target type |
| `reportStatus` | `enums.ReportStatusMap` | Report status |

### Gioi han cua Enum API

Enum API chi phu hop voi enum don gian `id/name`.

SEO refType khong nen dua vao Enum API nay vi can nhieu metadata hon:

- `key`
- `label`
- `family`
- `sourceService`
- `resolverKey`
- `searchable`
- `uniquePerRef`
- `supportsTemplate`
- `defaultPageType`
- `adminUrlPattern`
- `publicUrlPattern`

Vi vay SEO refType duoc tach thanh API rieng: `GET /v2/crm/admin/seo/ref-types`.

## 2. SEO refType enum hien tai

File:

```txt
crm-service/internal/enums/seo_ref_type.go
crm-service/internal/enums/seo.go
```

Gia tri enum hien dang dung trong repo:

| Value | Go enum | Key | Y nghia |
|---:|---|---|---|
| `0` | manual fallback | `manual` | Trang thu cong, khong gan source |
| `100` | `ESEORefTypeAdmUnit` | `adm_unit` | Don vi hanh chinh |
| `200` | `ESEORefTypeParcel` | `parcel` | Thua dat / parcel |
| `300` | `ESEORefTypeRegion` | `region` | Khu vuc / area region |
| `400` | `ESEORefTypeProject` | `project` | Du an bat dong san |
| `500` | `ESEORefTypeDocument` | `document` | Tai lieu phap ly/quy hoach/ho so |
| `600` | `ESEORefTypeMap` | `map` | Ban do / layer / workspace |
| `700` | `ESEORefTypeReport` | `report` | Bao cao |

Luu y quan trong:

- Tai lieu thiet ke ban dau co de xuat value `10/20/30/...`.
- Code hien tai dang dung `100/200/300/...`.
- Nen tiep tuc giu `100/200/300/...` de khong pha du lieu SEO da ton tai.

Validation hien co:

```txt
IsValidSEORefType(value)
```

Key mapping hien co:

```txt
SEORefTypeKey(value)
```

## 3. SEO Ref Registry

File:

```txt
crm-service/internal/domain/seo/seo_ref_registry.go
```

Registry la source of truth cho SEO refType. Frontend khong can biet `project` lay tu BDSPro hay `parcel` lay tu TQD. Frontend chi goi CRM.

Struct:

```go
type SeoRefTypeConfig struct {
    Value            uint32
    Key              string
    Label            string
    Description      string
    Family           string
    SourceService    string
    ResolverKey      string
    Searchable       bool
    UniquePerRef     bool
    SupportsTemplate bool
    DefaultPageType  string
    AdminURLPattern  string
    PublicURLPattern string
}
```

### Registry hien tai

| Value | Key | Label | Family | Source service | Resolver key | Searchable | Unique per ref |
|---:|---|---|---|---|---|---|---|
| `0` | `manual` | Trang thu cong | `manual` | `crm` | `manual` | `false` | `false` |
| `100` | `adm_unit` | Don vi hanh chinh | `location` | `bdspro` | `bdspro.region` | `true` | `true` |
| `300` | `region` | Khu vuc | `location` | `bdspro` | `bdspro.area_region` | `true` | `true` |
| `200` | `parcel` | Thua dat | `planning` | `tqd` | `tqd.parcel` | `true` | `true` |
| `400` | `project` | Du an | `real_estate` | `bdspro` | `bdspro.project` | `true` | `true` |
| `500` | `document` | Tai lieu | `document` | `tqd` | `tqd.legal_document` | `true` | `true` |
| `600` | `map` | Ban do / lop quy hoach | `map` | `tqd` | `tqd.qh_layer` | `true` | `true` |
| `700` | `report` | Bao cao | `report` | `crm` | `crm.report` | `true` | `true` |

Ham lay config:

```txt
SeoRefTypeConfigs()
GetSeoRefTypeConfig(value, resolverKey, sourceService)
```

## 4. SEO Ref APIs

Proto:

```txt
shared/protobuf/schema/crm/seo_domain.proto
```

Handler:

```txt
crm-service/infra/handler/seo_domain_handler.go
```

Usecase:

```txt
crm-service/internal/usecase/seo_domain_usecase.go
```

Mapper:

```txt
crm-service/infra/mapper/seo_domain_mapper.go
```

DTO:

```txt
crm-service/internal/dto/seo_ref_dto.go
```

### 4.1 Get SEO ref types

```txt
GET /v2/crm/admin/seo/ref-types
```

Response:

```json
{
  "data": [
    {
      "value": 400,
      "key": "project",
      "label": "Du an",
      "description": "Du an bat dong san",
      "family": "real_estate",
      "sourceService": "bdspro",
      "resolverKey": "bdspro.project",
      "searchable": true,
      "uniquePerRef": true,
      "supportsTemplate": true,
      "defaultPageType": "entity",
      "adminUrlPattern": "",
      "publicUrlPattern": ""
    }
  ]
}
```

Trang thai: da implement.

### 4.2 Search SEO ref values

```txt
GET /v2/crm/admin/seo/ref-values?refType=400&text=abc&page=1&size=20
GET /v2/crm/admin/seo/ref-values?refType=500&resolverKey=tqd.legal_document&text=abc
```

Request fields:

| Field | Bat buoc | Ghi chu |
|---|---|---|
| `refType` | Co | Gia tri refType |
| `text` | Khong | Tu khoa search |
| `page` | Khong | Trang |
| `size` | Khong | So item/trang |
| `resolverKey` | Khong | Resolver cu the neu refType co nhieu source |
| `sourceService` | Khong | Filter source service |

Response mong muon:

```json
{
  "data": [
    {
      "refType": 400,
      "refTypeKey": "project",
      "refId": 123,
      "label": "Project ABC",
      "description": "",
      "imageUrl": "",
      "adminUrl": "",
      "publicUrl": "",
      "alreadyHasSeoPage": false,
      "seoPageId": 0,
      "sourceService": "bdspro",
      "resolverKey": "bdspro.project"
    }
  ],
  "total": 1
}
```

Trang thai: API contract da co, usecase da validate registry, nhung resolver thuc te chua implement. Hien tra `501 resolver SEO chua duoc trien khai: {resolverKey}`.

### 4.3 Resolve SEO ref value

```txt
GET /v2/crm/admin/seo/ref-values/{refType}/{refId}
GET /v2/crm/admin/seo/ref-values/400/123
GET /v2/crm/admin/seo/ref-values/500/123?resolverKey=tqd.legal_document
```

Response mong muon:

```json
{
  "refType": 400,
  "refTypeKey": "project",
  "refId": 123,
  "label": "Project ABC",
  "description": "",
  "imageUrl": "",
  "adminUrl": "",
  "publicUrl": "",
  "suggestedSlug": "project-abc",
  "suggestedTitle": "Project ABC",
  "suggestedMetaTitle": "Project ABC",
  "suggestedMetaDescription": "",
  "sourceService": "bdspro",
  "resolverKey": "bdspro.project",
  "hash": "",
  "dataJson": "{}"
}
```

Trang thai: API contract da co, usecase da validate registry/refId, nhung resolver thuc te chua implement. Hien tra `501 resolver SEO chua duoc trien khai: {resolverKey}`.

## 5. Gateway

File:

```txt
gateway-service/service/crm.go
```

CRM gateway da register:

```txt
RegisterSeoDomainServiceHandlerFromEndpoint
```

Vi vay cac route SEO Domain va SEO Ref API co the di qua gateway sau khi generated protobuf duoc cap nhat o moi moi truong build.

## 6. Luong xu ly mong muon

### 6.1 Frontend lay refType

```txt
Frontend
-> GET /v2/crm/admin/seo/ref-types
-> Hien thi dropdown/group theo family/sourceService/resolverKey
```

Frontend khong hard-code service nguon.

### 6.2 Frontend search ref value

```txt
Frontend
-> GET /v2/crm/admin/seo/ref-values?refType=400&text=abc
-> CRM doc registry
-> CRM chon resolver bdspro.project
-> Resolver goi BDSPro source API
-> Normalize ve SeoRefValueResponse
-> Check duplicate SEO page
-> Tra ve frontend
```

### 6.3 Frontend resolve ref value

```txt
Frontend
-> GET /v2/crm/admin/seo/ref-values/400/123
-> CRM doc registry
-> CRM chon resolver bdspro.project
-> Resolver goi source detail
-> Normalize snapshot
-> Tinh hash
-> Tra ve SeoRefSnapshotResponse
```

### 6.4 Tao SEO page tu source

```txt
POST /v2/crm/admin/seo-domains
-> Validate refType/refId/refSource
-> Resolve source
-> Check duplicate
-> Save snapshot/hash/label/url
-> source_status = linked
-> need_generate = true
```

Trang thai hien tai: create SEO domain van luu `refType/refId` co ban, chua resolve source snapshot/hash.

## 7. Cac source service du kien

### BDSPro

Da co public list API nen phu hop lam source search ban dau:

| Resolver key | API nguon goi y |
|---|---|
| `bdspro.region` | `/v2/bdspro/v2/list/region` |
| `bdspro.area_region` | `/v2/bdspro/v2/list/area-region` |
| `bdspro.project` | `/v2/bdspro/v2/list/project` |
| `bdspro.doc_type` | `/v2/bdspro/v2/list/doc-type` |

### TQD

TQD phu hop voi quy hoach, parcel, layer, map, legal document, report:

| Resolver key | Proto/API nguon goi y |
|---|---|
| `tqd.parcel` | `tqd/parcel.proto` |
| `tqd.qh_layer` | `tqd/qh_layer.proto`, `tqd/map_workspace.proto` |
| `tqd.legal_document` | `tqd/layer_legal.proto` hoac document source |
| `tqd.report` | `tqd/report.proto` |

### CRM

CRM source noi bo:

| Resolver key | Nguon goi y |
|---|---|
| `crm.report` | `crm/report.proto`, report repo/usecase |

## 8. Van de refSource

Hien model `SeoDomain` moi co:

```txt
ref_type
ref_id
```

Voi cac loai nhu `document`, `report`, `map`, cung mot `refType/refId` co the nam o nhieu service khac nhau.

Vi du:

```txt
refType = document, refId = 123, source = tqd.legal_document
refType = document, refId = 123, source = crm.document
```

Neu khong co `ref_source`, hai record nay co the bi trung logic.

Khuyen nghi them:

```sql
ALTER TABLE seo_domain
ADD COLUMN IF NOT EXISTS ref_source VARCHAR(80) NOT NULL DEFAULT '';
```

Unique index nen tinh ca `ref_source`:

```sql
CREATE UNIQUE INDEX IF NOT EXISTS uniq_seo_domain_active_ref_source
ON seo_domain(page_type, ref_type, ref_source, ref_id)
WHERE deleted_at IS NULL
  AND ref_id IS NOT NULL
  AND ref_type <> 0;
```

Trang thai: chua implement trong model/migration hien tai.

## 9. Source snapshot va hash

Can bo sung vao `seo_domain`:

| Field | Muc dich |
|---|---|
| `page_type` | Phan loai manual/entity/hub/guide |
| `ref_source` | Source cu the/resolver key |
| `source_status` | `manual`, `linked`, `stale`, `missing`, `permission_denied` |
| `ref_label` | Ten hien thi source tai thoi diem sync |
| `ref_url` | URL source/public source |
| `ref_snapshot_json` | Snapshot normalized |
| `ref_hash` | Hash source snapshot |
| `ref_last_synced_at` | Lan sync source gan nhat |
| `ref_missing` | Source khong con ton tai |
| `template_id` | Template generation |
| `quality_score` | Diem quality SEO |
| `issues_json` | Loi/canh bao audit |

Trang thai: chua implement trong model/migration hien tai.

## 10. Checklist muc do hoan thien

| Hang muc | Trang thai | Ghi chu |
|---|---|---|
| Enum API pho thong `GET /v2/enums/crm/{name}` | Da co | Phu hop enum `id/name` |
| SEO refType enum values | Da co | Dang dung `100/200/300/...` |
| `SEORefTypeKey` | Da co | Map value sang key |
| SEO ref registry metadata | Da co | `SeoRefTypeConfigs()` |
| API `GET /v2/crm/admin/seo/ref-types` | Da co | Tra metadata tu registry |
| Proto cho `ref-values` search/resolve | Da co | Contract da chot |
| Handler cho `ref-values` search/resolve | Da co khung | Goi usecase |
| Usecase validate registry cho search/resolve | Da co | Tra `501` vi chua co resolver |
| BDSPro resolver | Chua co | Can goi public list/detail |
| TQD resolver | Chua co | Can goi parcel/layer/document/report |
| CRM report resolver | Chua co | Can goi report usecase/repo |
| Duplicate check theo ref source | Chua co | Can repo `GetByRefSource` |
| `ref_source` DB column | Chua co | Can migration/model/dto |
| Source snapshot/hash | Chua co | Can model/migration/usecase |
| Touch by ref source-aware | Chua co | Hien moi mark `need_generate` theo `refType/refId` |
| Sync/relink/detach source | Chua co | Can API/usecase/event |
| Event/audit SEO source | Chua co | Can repo/event table neu production |

## 11. De xuat thu tu lam tiep

1. Them `ref_source`, `source_status`, `ref_snapshot_json`, `ref_hash`, `ref_last_synced_at`, `ref_missing` vao model va migration.
2. Them repo methods:
   - `GetByRefSource(refType, refSource, refID)`
   - `GetListByRefSource(refType, refSource, refID)`
   - `MarkNeedGenerateByRefSource(refType, refSource, refID)`
3. Tao package resolver:
   - `crm-service/internal/resolver/seo`
4. Implement resolver registry:
   - `Registry.Get(key)`
   - `Registry.Search(req)`
   - `Registry.Resolve(req)`
5. Implement resolver uu tien:
   - `bdspro.project`
   - `bdspro.region`
   - `bdspro.area_region`
   - `tqd.parcel`
   - `tqd.qh_layer`
6. Nang create SEO domain:
   - resolve source truoc khi create
   - save snapshot/hash/label/url
   - check duplicate theo `ref_source`
7. Nang `TouchSeoDomainByRef` thanh source-aware:
   - nhan them `refSource`/`resolverKey`
   - resolve source
   - hash khac thi stale + need_generate
   - resolve fail thi missing
8. Them API operation:
   - validate
   - sync-source
   - relink-source
   - detach-source

## 12. Dinh nghia "day du"

He thong chi duoc xem la day du khi dat ca 5 dieu kien:

1. Frontend chi goi CRM SEO API, khong goi truc tiep BDSPro/TQD de chon SEO source.
2. Moi `refType` trong registry co resolver search/resolve hoat dong.
3. SEO domain luu du `ref_source`, snapshot, hash, source status.
4. Create/update/touch/sync/relink/detach deu dung chung resolver registry.
5. Duplicate SEO page duoc check theo `page_type + ref_type + ref_source + ref_id`.

Tinh den ban cap nhat nay: chua day du production, nhung API contract va registry da dung huong de tiep tuc hoan thien.
