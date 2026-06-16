# SEO UI/UX Redesign Map

Ngay cap nhat: 2026-06-16

## Ket luan thiet ke

Ban frontend SEO moi duoc dung nhu mot **SEO Operating System** thay vi mot bo CRUD page roi rac.

Tinh than tham khao tu usrnk1/Dribbble la: minimal SaaS admin, table density cao, detail drawer, popover/action compact, status chip, tab/section ro rang, bo goc 6-8px, nen sang, border mong, typography gon, khong hero trang tri.

Luu y: day la ban thiet ke nguyen ban lay cam hung ve language/pattern, khong copy pixel hay asset ban quyen tu Dribbble.

## Huong nghiep vu moi

SEO page khong con la `refType + refId` lam trung tam. UI moi xem moi trang la mot asset:

```txt
Source contract -> Source slot -> Snapshot/context -> Template contract -> Studio -> Render -> Publish -> Index -> Monitor
```

Ket qua:

- Admin thay ngay page nao bi ket va ket o dau.
- Template noi ro required/optional source slots.
- Source catalog noi ro adapter nao ready/stale/missing.
- Studio cho thay source context, editor, SERP preview, checklist va publish bar cung luc.
- Render/index/internal links/events/settings duoc tach thanh cac man van hanh ro rang.

## File da trien khai

```txt
frontend/src/pages/SEO/SeoShell.tsx
frontend/src/pages/SEO/seoNavigation.ts
frontend/src/pages/SEO/styles/seo-console.css
frontend/src/pages/SEO/design/seoWorkspaceMock.ts
frontend/src/pages/SEO/design/SeoDesignPrimitives.tsx
frontend/src/pages/SEO/overview/SeoOverviewDashboard.tsx
frontend/src/pages/SEO/pages/SeoPagesRoute.tsx
frontend/src/pages/SEO/create/SeoCreatePage.tsx
frontend/src/pages/SEO/studio/SeoStudioPage.tsx
frontend/src/pages/SEO/studio/SeoStudioHeader.tsx
frontend/src/pages/SEO/studio/SeoSourcePanel.tsx
frontend/src/pages/SEO/studio/SeoInsightPanel.tsx
frontend/src/pages/SEO/studio/SeoPublishBar.tsx
frontend/src/pages/SEO/templates/SeoTemplatesPage.tsx
frontend/src/pages/SEO/generation/SeoGenerationPage.tsx
frontend/src/pages/SEO/indexing/SeoIndexingPage.tsx
frontend/src/pages/SEO/performance/SeoPerformancePage.tsx
frontend/src/pages/SEO/internal-links/SeoInternalLinksPage.tsx
frontend/src/pages/SEO/system-events/SeoSystemEventsPage.tsx
frontend/src/pages/SEO/settings/SeoSettingsPage.tsx
frontend/src/components/ContentEditor/surfaces/StaticToolbar.tsx
frontend/src/components/ContentEditor/styles/content-editor.css
```

## Mapping 40 Dribbble references sang SEO Console

| # | Reference | Pattern rut ra | Ap dung vao code |
|---|---|---|---|
| 1 | Employee Profile Preview Cards | Preview cards co thong tin can quyet dinh, khong marketing | Drawer asset preview va source fields |
| 2 | Vendor Directory Table | Dense table + filter + status chip | `SeoPagesRoute` assets table |
| 3 | Employee Profile Accounts Table | Entity detail + related accounts | Source slots trong drawer va Studio left panel |
| 4 | Download Agent View | Step/accordion onboarding ro rang | `SeoCreatePage` workflow steps |
| 5 | Automatically Detected Accounts | Auto-detected list + action | Source catalog/results in Compose |
| 6 | Popover Requirements Details | Requirement popover/details | Template required/optional slots |
| 7 | Trust Center Subprocessors Listing | Listing tab + trust table | `SeoTemplatesPage` source catalog table |
| 8 | Request Changes Dialog | Review/approval dialog mindset | Publish readiness va blocked states |
| 9 | Security Measures Tab | Category tabs + measures | Settings governance groups |
| 10 | Download Documents Modal | Modal compact, checkbox/action clarity | Drawer footer action grouping |
| 11 | Request Access Dialog | Minimal modal decision | Drawer open/close/action design |
| 12 | Compliance Documents Tab | High-density document table | Indexing URL policy table |
| 13 | Device Issue Slide-Over | Slide-over detail view | `SeoDrawer` asset detail |
| 14 | App Versions Detail View | Version/detail table | Template status/readiness table |
| 15 | Table UI Upsell Callout | Inline table callout | Metrics + blocked/render cards |
| 16 | Trust Center Overview | System overview cards | `SeoOverviewDashboard` metrics/source health |
| 17 | MDM Apps Directory | Directory listing | Source catalog va assets registry |
| 18 | Device Table Actions Popover | Compact row actions | Icon buttons in assets table |
| 19 | Employee Portal Overview | Empty/overview structure | Overview summary without hero noise |
| 20 | Compliance Card Video | Status/progress storytelling | Lifecycle dots and progress bars |
| 21 | Security Training Initial State | Guided first-run state | Compose flow first step |
| 22 | Policy Detail View | Detail page structure | Studio/editor layout |
| 23 | Device Details Header | Detail header with actions | `SeoStudioHeader` |
| 24 | Monitor Bulk Actions Toolbar | Bulk/monitor toolbar | Assets command bar |
| 25 | Control Detail Collapsed Sidenav | Strong left nav | `SeoShell` sidebar |
| 26 | Compliance Monitoring Checks | Checklist thinking | `SeoQualityChecklist` panel |
| 27 | Activity History Control | Audit trail | `SeoSystemEventsPage` |
| 28 | Evidence Dropzone | Evidence/source attachment mindset | Source slot contract display |
| 29 | Overlay Drawer Assets Table | Table + overlay drawer | `SeoPagesRoute` + `SeoDrawer` |
| 30 | Control Detail View | Control detail and remediation | Asset blocked/source detail |
| 31 | Category Status Popover | Category status chip | `SeoPill`, `SeoProgress` |
| 32 | Compliance Program Overview | Program-level dashboard | Overview operating system |
| 33 | Compact Charts | Compact bars, no heavy dashboard | `SeoPerformancePage` chart bars |
| 34 | Faceted Filters Sidebar | Faceted filtering | Assets filter/select/search |
| 35 | Listing API Clients Detail View | Listing + detail view style | Core pattern cho Assets registry |
| 36 | Dashboard Compliance Overview | KPI + operational summary | Metrics grid |
| 37 | Category Grid | Template/category grid | `SeoTemplatesPage` cards |
| 38 | Schema Usage Insights | Schema/table insight | Template contract + source readiness |
| 39 | Config Schema Polling | Settings/schema polling | `SeoSettingsPage` governance cards |
| 40 | Query Error Location Highlight | Error trace clarity | Render failed/index blocked/event logs |

## Design rules da ap dung

- Bo goc: 6-8px cho button, panel, drawer, cards.
- Layout: sidebar nav + content, khong hero gradient.
- Table: cot uu tien theo nghiep vu `Asset / Sources / Template / Lifecycle / Generate / Index / Quality`.
- Drawer: mo detail tai cho, khong day admin sang trang khac khi chi can inspect.
- Editor: left source context, center writing paper, right SERP/checklist, bottom publish state.
- Color: neutral gray/white/black, blue cho primary/info, green/warn/red cho health state.
- Interaction: icon buttons, compact toolbar, clear hover/focus, select/search/filter dat gan table.

## Huong backend tuong lai nen khop voi UI

UI moi can backend cung cap cac concept:

```txt
SeoAsset
SeoSourceContract
SeoSourceSlot
SeoTemplateContract
SeoRenderContext
SeoRenderJob
SeoPublishPolicy
SeoIndexStatus
SeoAuditEvent
```

Khi backend doi, frontend chi can thay `seoWorkspaceMock.ts` bang service/API moi, con layout va interaction model co the giu.
