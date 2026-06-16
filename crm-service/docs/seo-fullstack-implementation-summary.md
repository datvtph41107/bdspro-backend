# Tong ket trien khai SEO Fullstack

Ngay cap nhat: 2026-06-15

## 1. Ket luan chung

Da ghep noi module SEO frontend voi CRM SEO backend theo huong source-aware:

- Frontend SEO Pages, Create SEO va Studio khong con phu thuoc mock chinh cho CRUD/operations.
- CRM backend da co `ref_source`, `source_status`, source snapshot, hash va thoi diem sync tren `seo_domain`.
- `GET /v2/crm/admin/seo/ref-types` tiep tuc la source of truth cho frontend.
- `GET /v2/crm/admin/seo/ref-values` va `GET /v2/crm/admin/seo/ref-values/{refType}/{refId}` da noi resolver registry thay vi chi tra 501.
- Create/update/touch SEO domain da dung chung resolver registry de resolve source, luu snapshot/hash/status va check duplicate theo `ref_type + ref_source + ref_id`.

## 2. Backend da trien khai

### 2.1 Data model va migration

Da them migration:

```txt
crm-service/infra/db/migrate_v2/000006_seo_source_snapshot.up.sql
crm-service/infra/db/migrate_v2/000006_seo_source_snapshot.down.sql
```

Migration bo sung cac cot:

```txt
ref_source
ref_label
ref_url
source_status
ref_missing
ref_snapshot_json
ref_hash
ref_last_synced_at
quality_score
```

Da bo sung index cho:

```txt
ref_type + ref_source + ref_id
source_status
ref_hash
ref_last_synced_at
quality_score
```

### 2.2 Domain/DTO/Proto/Mapper

Da mo rong:

```txt
crm-service/internal/domain/seo/seo_domain.go
crm-service/internal/dto/seo_domain_dto.go
crm-service/internal/dto/seo_ref_dto.go
shared/protobuf/schema/crm/seo_domain.proto
crm-service/infra/mapper/seo_domain_mapper.go
crm-service/infra/handler/seo_domain_handler.go
```

Response `SeoDomainResponse`, create/update/list request, ref-value response va touch-by-ref request da co source fields moi.

### 2.3 Resolver registry

Da them resolver layer:

```txt
crm-service/internal/usecase/seo_ref_resolver.go
```

Hien da noi resolver that cho:

| Resolver | Trang thai | Ghi chu |
|---|---|---|
| `bdspro.project` | Da search that | Resolve theo id se tim trong list vi BDSPro AdminProjectService chua co detail API |
| `bdspro.region` | Da search/resolve that | Goi AdminRegionService |
| `bdspro.area_region` | Da list/filter/resolve that | Goi AdminAreaRegionService |
| `tqd.parcel` | Chua co adapter | Can them TQD provider vao CRM |
| `tqd.qh_layer` | Chua co adapter | Can them TQD provider vao CRM |
| `tqd.legal_document` | Chua co adapter | Can chot source API tai TQD |
| `crm.report` | Chua co adapter | Can chot report resolver noi bo CRM |

### 2.4 Source-aware operations

Da nang cac flow:

- `SearchSeoRefValues`: chon resolver theo registry, normalize ket qua, check `alreadyHasSeoPage`.
- `ResolveSeoRefValue`: tra snapshot gom label/url/suggested slug/title/meta/hash/dataJson.
- `CreateSeoDomain`: neu co source thi resolve snapshot, luu `ref_source/ref_label/ref_url/source_status/ref_snapshot_json/ref_hash/ref_last_synced_at`.
- `UpdateSeoDomain`: khi doi `refType/refId/refSource` thi resolve lai source binding.
- `TouchByRef`: resolve source, so sanh hash, cap nhat `linked/stale/missing`, set `need_generate` khi source thay doi.

## 3. Frontend da trien khai

Da bo sung/ket noi:

```txt
frontend/src/interface/seo.ts
frontend/src/service/seo-service.ts
frontend/src/service/seo-pages-service.ts
frontend/src/service/seo-studio-service.ts
frontend/src/pages/SEO/create/SeoCreatePage.tsx
frontend/src/pages/SEO/create/store/seoCreateStore.ts
frontend/src/pages/SEO/pages/seoPagesConstants.ts
frontend/src/pages/SEO/studio/SeoSourcePanel.tsx
```

### 3.1 SEO API layer

`seo-service.ts` da co:

```txt
seoApi.list/getById/create/update/delete
seoApi.touchByRef
seoApi.getRefTypes
seoApi.searchRefValues
seoApi.resolveRefValue
```

Mapper frontend da doc du:

```txt
refSource
refLabel
refUrl
sourceStatus
refMissing
refSnapshotJson
refHash
refLastSyncedAt
qualityScore
```

### 3.2 SEO Pages operations

`seo-pages-service.ts` da chuyen tu mock memory sang backend-backed service:

- List doc tu `GET /v2/crm/admin/seo-domains`.
- Map `SeoDomain` thanh `SeoPage` de UI hien tai khong can doi lon.
- Publish/archive/update metadata/touch-by-ref di qua backend.
- KPI va issue summary tinh lai tu du lieu backend.

### 3.3 Create SEO

Create flow da:

- Doc ref types tu CRM.
- Search ref values tu CRM.
- Luu `refType/refId/refSource/sourceService`.
- Tao draft bang `POST /v2/crm/admin/seo-domains`.
- Mo Studio theo id record backend tra ve.

### 3.4 SEO Studio

Studio service da:

- Load draft tu `GET /v2/crm/admin/seo-domains/{id}`.
- Save content/title/summary qua `PUT /v2/crm/admin/seo-domains/{id}`.
- Publish qua backend update.
- Generate queue/touch source qua `touch-by-ref`.
- Source panel hien `refSource/sourceStatus/refHash`.

## 4. Ket qua verify

Da chay thanh cong:

```txt
buf generate --template ./buf.gen.yaml --path crm
go test ./infra/handler ./infra/client ./infra/mapper ./infra/postgre ./internal/domain/seo ./internal/dto ./internal/repo ./internal/interface/provider
GOCACHE=/tmp/codex-go-build go test -vet=off ./internal/usecase
```

Ket qua:

```txt
Backend SEO-related packages: OK
internal/usecase compile voi -vet=off: OK
```

`go test ./...` trong `crm-service` hien bi dung boi loi cu khong thuoc SEO:

```txt
crm-service/internal/usecase/follow_usecase.go
fmt.Errorf call has arguments but no formatting directives
```

`npm run build` trong frontend hien bi dung boi dependency editor chua co trong `package.json/node_modules`:

```txt
@tiptap/react
@tiptap/core
@tiptap/starter-kit
@tiptap/suggestion
@tiptap/extension-placeholder
@tiptap/extension-link
```

Can bo sung cac dependency Tiptap nay truoc khi build production frontend.

## 5. Viec con lai de production-ready tuyet doi

1. Them TQD provider vao CRM service de implement that:
   - `tqd.parcel`
   - `tqd.qh_layer`
   - `tqd.legal_document`
   - `tqd.report`
2. Them CRM report resolver cho `crm.report`.
3. Neu business yeu cau mot source chi co mot SEO page, enable unique index:

```sql
CREATE UNIQUE INDEX IF NOT EXISTS uniq_seo_domain_active_ref_source
ON seo_domain(ref_type, ref_source, ref_id)
WHERE deleted_at IS NULL AND ref_id IS NOT NULL;
```

4. Sua loi vet cu trong `follow_usecase.go` de `go test ./...` sach toan service.
5. Bo sung dependency Tiptap de `npm run build` frontend hoan tat.

## 6. Cau chot

He thong da duong hoang buoc qua giai doan scaffold: frontend da goi CRM SEO APIs, backend da luu source-aware state, BDSPro resolver da chay qua client that, create/update/touch da cung dung snapshot/hash/status. Phan con lai khong phai UI nua ma la bo sung adapter source service TQD/CRM report va don sach dependency/build pipeline.
