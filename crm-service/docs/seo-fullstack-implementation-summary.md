# Tong ket trien khai SEO Domain

Ngay cap nhat: 2026-06-17

## 1. Ket luan chung

Da ghep noi module SEO frontend voi CRM SEO backend theo huong source-aware:

- Frontend SEO Pages, Create SEO va Studio khong con phu thuoc mock chinh cho CRUD/operations.
- CRM backend da co `ref_source`, `source_status`, source snapshot, hash va thoi diem sync tren `seo_domain`.
- `GET /v2/crm/admin/seo/ref-types` tiep tuc la source of truth cho frontend.
- `GET /v2/crm/admin/seo/ref-values` va `GET /v2/crm/admin/seo/ref-values/{refType}/{refId}` da noi resolver registry thay vi chi tra 501.
- Create/update/touch SEO domain da dung chung resolver registry de resolve source, luu snapshot/hash/status va check duplicate theo `ref_type + ref_source + ref_id`.
- Da bo sung `GET /v2/crm/admin/seo/operations/dashboard` lam API tong hop cho admin SEO: KPI, trend, opportunity, top query, event, checklist, render queue, indexing rows va source health.
- Frontend SEO da chuyen sang backend-first/no-mock-runtime: API co du lieu thi hien thi, API chua co thi hien empty state/notice "chua co du lieu", khong tu dung du lieu mau.

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

| Resolver             | Trang thai                  | Ghi chu                                                                                                             |
| -------------------- | --------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `bdspro.project`     | Da search that              | Resolve theo id se tim trong list vi BDSPro AdminProjectService chua co detail API                                  |
| `bdspro.region`      | Da search/resolve that      | Goi AdminRegionService                                                                                              |
| `bdspro.area_region` | Da list/filter/resolve that | Goi AdminAreaRegionService                                                                                          |
| `tqd.parcel`         | Da search/resolve that      | Goi TQD provider; search text numeric se resolve truc tiep theo parcel id, text thuong se loc tu danh sach generate |
| `tqd.qh_layer`       | Chua co adapter             | Can them TQD provider vao CRM                                                                                       |
| `tqd.legal_document` | Chua co adapter             | Can chot source API tai TQD                                                                                         |
| `crm.report`         | Chua co adapter             | Can chot report resolver noi bo CRM                                                                                 |

### 2.4 Source-aware operations

Da nang cac flow:

- `SearchSeoRefValues`: chon resolver theo registry, normalize ket qua, check `alreadyHasSeoPage`.
- `ResolveSeoRefValue`: tra snapshot gom label/url/suggested slug/title/meta/hash/dataJson.
- `CreateSeoDomain`: neu co source thi resolve snapshot, luu `ref_source/ref_label/ref_url/source_status/ref_snapshot_json/ref_hash/ref_last_synced_at`.
- `UpdateSeoDomain`: khi doi `refType/refId/refSource` thi resolve lai source binding.
- `TouchByRef`: resolve source, so sanh hash, cap nhat `linked/stale/missing`, set `need_generate` khi source thay doi.
- `GetSeoOperationsDashboard`: gom du lieu tu `seo_domain` de tao KPI, checklist van hanh, danh sach viec can xu ly, render queue, indexing rows va source health cho UI admin.
- `GetSeoOperationsDashboard`: trend/clicks/impressions/topQueries khong con suy dien tu title/qualityScore. Backend da them `SeoTrendRepo` doc bang `seo_trend`; neu chua import GSC/analytics thi tra mang rong.

### 2.5 Tu duy nghiep vu moi

UI admin khong yeu cau nguoi dung hieu `refType/refId` nhu mot khai niem ky thuat. Admin nhin thay:

- Nguon du lieu: khu vuc, du an, thua dat, tai lieu, bao cao.
- Trang SEO: title, slug, template, trang thai publish, render, indexing.
- Trang thai van hanh: du lieu moi/cu/thieu, can render, co duoc index khong.

Backend van giu `refType/refId/refSource` nhu "ma van don noi bo" de tim lai nguon su that. Frontend hien thi thanh ngon ngu nghiep vu de admin thao tac nhanh va dung.

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
frontend/src/pages/SEO/insights/SeoInsightPage.tsx
frontend/src/pages/SEO/overview/SeoOverviewDashboard.tsx
frontend/src/pages/SEO/generation/SeoGenerationPage.tsx
frontend/src/pages/SEO/indexing/SeoIndexingPage.tsx
frontend/src/pages/SEO/design/seoBackendMappers.ts
frontend/src/pages/SEO/design/seoDesignTypes.ts
frontend/src/pages/SEO/rendering/seoPreviewRenderer.ts
```

### 3.1 SEO API layer

`seo-service.ts` da co:

```txt
seoApi.list/getById/create/update/delete
seoApi.touchByRef
seoApi.getRefTypes
seoApi.searchRefValues
seoApi.resolveRefValue
seoOperationsApi.dashboard
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

Operations mapper da doc:

```txt
kpis
trend
opportunities
topQueries
systemEvents
requirements
renderQueue
indexingRows
sourceHealth
```

### 3.2 SEO Pages operations

`seo-pages-service.ts` da chuyen tu mock memory sang backend-backed service:

- List doc tu `GET /v2/crm/admin/seo-domains`.
- Map `SeoDomain` thanh `SeoPage` de UI hien tai khong can doi lon.
- Publish/archive/update metadata/touch-by-ref di qua backend.
- KPI va issue summary tinh lai tu du lieu backend.
- Giao dien SEO Pages moi doc `GET /v2/crm/admin/seo-domains` voi `page=0`, map `SeoDomain` thanh `SeoAsset`, neu API loi thi hien notice loi va empty state, khong fallback mock.
- Toan bo table trong module SEO moi da chuyen sang custom directory rows de kiem soat UI/UX tot hon.

### 3.3 Create SEO

Create flow da:

- Doc ref types tu CRM.
- Search ref values tu CRM.
- Luu `refType/refId/refSource/sourceService`.
- Tao draft bang `POST /v2/crm/admin/seo-domains`.
- Mo Studio theo id record backend tra ve.
- RefType khac manual bat buoc phai co refValue/snapshot resolve that truoc khi tao draft.
- HTML preview duoc render tu snapshot/source data that bang renderer frontend mac dinh; day la preview de admin xem bo cuc, chua thay the template engine backend.

### 3.4 SEO Studio

Studio service da:

- Load draft tu `GET /v2/crm/admin/seo-domains/{id}`.
- Save content/title/summary qua `PUT /v2/crm/admin/seo-domains/{id}`.
- Publish qua backend update.
- Generate queue/touch source qua `touch-by-ref`.
- Source panel hien `refSource/sourceStatus/refHash`.
- Preview modal hien HTML hien tai cua SEO domain/editor, canh bao khi `needGenerate=true`.

### 3.5 Insight, Overview, Render, Indexing

Da noi cac man hinh van hanh moi voi `GET /v2/crm/admin/seo/operations/dashboard`:

- Insight: KPI, top query, opportunity.
- Tong quan: KPI, source health, event gan day, danh sach trang can theo doi.
- Render: render queue tu backend.
- Indexing: robots/canonical/sitemap/indexing rows tu backend.
- Nhat ky: systemEvents tu backend.
- Internal Links: doc `GET /v2/crm/admin/seo-domains/internal-links`.
- Performance: doc `topQueries` tu `seo_trend` thong qua operations dashboard.

Khi backend/API chua co du lieu, UI hien empty state va noi ro API nao chua co. File mock SEO cu da duoc xoa khoi frontend SEO.

### 3.6 Template, Render HTML, XML

Da co:

- Sitemap XML backend: `GET /v2/crm/sitemap.xml` va `GET /v2/seo/sitemap.xml`.
- Render data noi bo: `/v2/internal/seo/render-data/...` va `/v2/internal/seo/render-data/by-slug/...`.
- Preview HTML trong admin: dung `seoPreviewRenderer.ts` de xem bo cuc public tu `title/description/content/refSnapshotJson/internalLinks`.

Chua co:

- Template API rieng.
- Render-preview API backend cho admin.
- Job render HTML that ghi `staticHtmlPath/staticHtmlHash`.
- Job ping sitemap/indexing inspection that.

Vi vay UI Templates va Settings khong hien du lieu mau; cac man nay hien "chua co API" va danh sach endpoint can bo sung.

## 4. Ket qua verify

Da chay thanh cong:

```txt
buf generate --template ./buf.gen.yaml --path crm
GOCACHE=/tmp/codex-go-build go test -vet=off ./internal/usecase ./infra/postgre ./infra/handler ./internal/repo
NODE_OPTIONS=--max-old-space-size=4096 ./node_modules/.bin/vite build --ssr src/pages/SEO/index.tsx --minify=false --outDir dist-seo-verify
```

Ket qua:

```txt
Backend SEO-related packages: OK
Frontend SEO SSR/module verify: OK
```

Ghi chu frontend:

```txt
Full production build hien transform xong 8443 modules nhung bi kill 137 o buoc rendering chunks, ke ca khi tang NODE_OPTIONS=8192.
Day la gioi han tai nguyen/bundle cua frontend hien tai, khong phai loi truc tiep tu module SEO trong output.
Canh bao cu cua react-leaflet-draw ve leaflet.draw default export van con.
```

## 5. Viec con lai de production-ready tuyet doi

1. Hoan thien cac resolver source con thieu:
    - `tqd.qh_layer`
    - `tqd.legal_document`
    - `crm.report`
2. Them CRM report resolver cho `crm.report`.
3. Neu business yeu cau mot source chi co mot SEO page, enable unique index:

```sql
CREATE UNIQUE INDEX IF NOT EXISTS uniq_seo_domain_active_ref_source
ON seo_domain(ref_type, ref_source, ref_id)
WHERE deleted_at IS NULL AND ref_id IS NOT NULL;
```

4. Noi Google Search Console/Search Analytics/importer that de ghi vao `seo_trend`; neu bang nay rong thi Insight/Performance se hien chua co du lieu.
5. Bo sung Template API, render-preview API backend va renderer/job tao HTML public that.
6. Tach chunk frontend/giam bundle de `npm run build` full khong bi kill 137.
7. Bo sung job scheduler that cho render queue, sitemap ping, source refresh va indexing inspection neu can van hanh tu dong.

## 6. Cau chot

He thong da buoc qua giai doan scaffold: backend da co source-aware state, resolver registry, ref value API, operations dashboard API, sitemap XML, render-data noi bo, TQD parcel resolver va trend repo doc `seo_trend`; frontend da co UI admin no-mock-runtime cho Pages/Insight/Overview/Render/Indexing/Create/Studio/Internal Links/Nhat ky. Phan con lai de production-ready tuyet doi la noi tiep cac source service chua co contract, nap SEO analytics that vao `seo_trend`, them Template/Render API that va dua render/indexing thanh job van hanh tu dong.
