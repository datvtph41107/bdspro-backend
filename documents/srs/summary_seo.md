A.1 Mục tiêu tài liệu
A.1.1 Mục tiêu SEO của hệ thống

QH Pro là nền tảng dữ liệu quy hoạch, đất đai và BĐS. Tài liệu SEO nhằm:

Tăng khả năng hiển thị trên Google, Bing và các công cụ tìm kiếm.
Tăng khả năng được AI Search và LLM hiểu, trích xuất và sử dụng.
Xây dựng Semantic SEO, Knowledge Graph và Entity Search.
Tối ưu khám phá dữ liệu và liên kết nội bộ.
Đảm bảo tuân thủ pháp lý, bảo mật và phân quyền dữ liệu.

👉 Tài liệu này là yêu cầu kỹ thuật SEO, không phải SEO marketing.

A.1.2 Phạm vi áp dụng

Áp dụng cho các phần có thể công khai của hệ thống QH Pro:

Trang tra cứu và hiển thị dữ liệu quy hoạch.
Trang chi tiết: thửa đất, khu quy hoạch, đồ án, văn bản pháp lý, thực thể liên quan.
Kho dữ liệu và kho tri thức quy hoạch.
URL, metadata, structured data, Open Graph, sitemap, AI summary.

Không áp dụng cho:

Hệ thống quản trị nội bộ.
API nội bộ.
Chức năng vận hành.
Dữ liệu không công khai.
UI tạm thời, popup kỹ thuật.
A.1.3 Đối tượng sử dụng tài liệu
BA: xác định phạm vi SEO trong SRS.
UI/UX: thiết kế breadcrumb, internal link, SEO content block, structured content.
Frontend: metadata, canonical, structured data, SSR/hybrid rendering.
Backend: URL, metadata, sitemap, AI summary, SEO output.
QA: kiểm thử và nghiệm thu SEO.
Product: kiểm soát định hướng SEO tổng thể.
A.1.4 Quan hệ với bộ SRS QH Pro
Đây là phụ lục kỹ thuật, không thay thế SRS.
SRS là nguồn chính cho nghiệp vụ, API, luồng xử lý.
SEO chỉ bổ sung lớp kỹ thuật:
URL, canonical
metadata, structured data
Open Graph, sitemap
AI summary, internal link
security & visibility SEO
QA checklist

A.1.5 Quan hệ Backend / Frontend / Mobile
Backend (SEO Engine)

Backend chịu trách nhiệm tạo toàn bộ SEO output:

URL chuẩn + Canonical URL
Metadata (title, description,…)
Structured Data (JSON-LD)
Sitemap
AI Summary
Trạng thái Index / NoIndex
Cung cấp SEO payload cho client
Frontend Web (Render SEO)

Frontend chịu trách nhiệm render SEO từ backend:

Render <head> tags (metadata, canonical,…)
Render Structured Data
Render Open Graph
Render Internal Link
Đảm bảo SSR/SEO-friendly rendering để index được
Mobile App

Không trực tiếp SEO nhưng phải hỗ trợ:

Deep Link
Share Link
Open Graph compatibility
Sync URL + Entity với Web
Hỗ trợ điều hướng phục vụ share/discover/AI
Nguyên tắc chung
SEO là xuyên suốt hệ thống (Backend + Frontend + Mobile)
Phải đảm bảo đồng bộ URL + Entity + Share behavior
A.2 Mục tiêu SEO tổng thể
Mô hình hệ thống
QH Pro là Data Platform, không phải website content/blog
SEO tập trung vào entity + dữ liệu có cấu trúc + truy vấn dữ liệu
A.2.1 Google Search (SEO truyền thống)
Đối tượng cần SEO mạnh (core entities)
Đơn vị hành chính (hiện tại + lịch sử)
Thửa đất
Khu quy hoạch
Đồ án / dự án quy hoạch
Văn bản pháp lý
Bản đồ quy hoạch
Báo cáo / phân tích
Yêu cầu kỹ thuật chính
URL rõ ràng, có nghĩa
Mỗi entity có trang index độc lập
Giảm duplicate content
Hỗ trợ lịch sử địa danh / đơn vị hành chính
Có khả năng mở rộng schema trong tương lai
A.2.2 AI Search (LLM / AI Assistant)
Target systems
ChatGPT, Gemini, Copilot, Perplexity, LLM khác
Yêu cầu dữ liệu
Structured data rõ ràng
Entity được định danh nhất quán
Quan hệ giữa entity phải rõ
Dữ liệu dễ parse / extract / summarize
Mục tiêu
AI hiểu đúng entity + quan hệ địa chính
AI có thể trích dẫn QH Pro
AI trả lời dựa trên dữ liệu hệ thống
A.2.3 Knowledge Graph
Entity core
Đơn vị hành chính (và lịch sử)
Thửa đất
Khu quy hoạch
Đồ án / dự án
Văn bản pháp lý
Bản đồ / báo cáo
Loại quan hệ cần có (quan trọng để code graph)
Thuộc về / nằm trong
Kế thừa / thay thế / sáp nhập
Áp dụng cho / liên quan
Được điều chỉnh bởi
Được phê duyệt bởi
Mục tiêu kỹ thuật
Tạo graph thay vì page rời rạc
Hỗ trợ semantic search + AI retrieval
Tăng rich result / contextual understanding
A.2.4 Internal Linking
Các liên kết bắt buộc
Hierarchy hành chính: quốc gia → tỉnh → huyện → xã
Địa bàn ↔ thửa đất
Địa bàn ↔ khu quy hoạch
Khu quy hoạch ↔ đồ án
Đồ án ↔ văn bản pháp lý
Văn bản ↔ bản đồ
Báo cáo ↔ dữ liệu nguồn
Mục tiêu kỹ thuật
Tăng crawlability
Truyền ngữ cảnh giữa entity
Xây knowledge graph nội bộ
Hỗ trợ navigation UI + SEO bot
A.2.5 Discover & Share
Use cases
Share: thửa đất, khu quy hoạch, đồ án, văn bản, bản đồ, báo cáo, snapshot
Yêu cầu kỹ thuật
URL stable, permanent
Có metadata share (title, description, preview)
Open Graph chuẩn cho social + chat apps
Không leak data restricted
Sync Web ↔ Mobile ↔ Share system
A.3 Nguyên tắc triển khai SEO (CORE RULES)
Áp dụng cho toàn hệ thống SEO

Dùng để quyết định:

URL design
Metadata design
Structured data
Sitemap
Internal link
AI summary
Index / NoIndex
Nguyên tắc cốt lõi (quan trọng cho dev)
SEO phải thống nhất toàn hệ thống
Ưu tiên cấu trúc entity-based
URL + entity + metadata phải đồng bộ
Dữ liệu phải crawlable + machine-readable
Internal link phải có ngữ nghĩa, không chỉ navigation
Nội dung phải hỗ trợ index và AI extraction

A.3 Nguyên tắc triển khai SEO (CORE RULES)
A.3.1 Entity First (quan trọng nhất)
SEO xoay quanh Entity, không xoay quanh UI hay keyword.
Core entities:
Đơn vị hành chính (và lịch sử)
Thửa đất
Khu quy hoạch
Đồ án / dự án quy hoạch
Văn bản pháp lý
Bản đồ, báo cáo, snapshot
Quy tắc kỹ thuật:
Mỗi Entity phải có:
ID định danh rõ ràng
URL ổn định
Metadata / Structured Data / AI Summary phải sinh từ Entity
Internal link dựa trên quan hệ Entity
Không SEO theo filter, UI state, search result

👉 Entity = lõi toàn bộ SEO system

A.3.2 Canonical First
Mỗi Entity chỉ có 1 URL canonical
Mục tiêu:
tránh duplicate content
tập trung SEO value
Quy tắc:
URL filter/search/share/snapshot:
phải xác định rõ:
Index hay không
Canonical về đâu
Canonical luôn trỏ về:
Entity gốc
A.3.3 Public / Permission First
Chỉ SEO dữ liệu được phép công khai
Rule:
Public → có thể SEO
Restricted → kiểm soát theo policy
Private → không index
Login-required → không SEO mặc định

👉 SEO không được phá security model

A.3.4 NoIndex First
Default: KHÔNG INDEX
Chỉ index khi:
Dữ liệu hoàn chỉnh
Có Entity rõ ràng
Có giá trị độc lập
Luôn NoIndex:
search result
filter state
UI state
draft / temp data
nội dung nội bộ

👉 Index phải là quyết định có kiểm soát

A.3.5 Legal Safe
SEO phải tuân thủ pháp lý dữ liệu địa chính
Rule:
Chỉ SEO dữ liệu công bố hợp lệ
Không SEO dữ liệu nhạy cảm / hạn chế
Không gây hiểu nhầm pháp lý trong metadata / AI summary

👉 Legal > SEO luôn ưu tiên

A.3.6 AI Search Ready

SEO phải tối ưu cho AI / LLM

Yêu cầu kỹ thuật:
Data structured rõ ràng
Entity consistent (định danh thống nhất)
Quan hệ entity rõ ràng
Có khả năng extract + summarize
Mục tiêu:
AI hiểu đúng dữ liệu QH Pro
AI có thể trích dẫn QH Pro
AI dùng QH Pro làm nguồn trả lời
A.3.7 Internal Link First
Internal link là backbone SEO
Link graph chính:
Hành chính hierarchy (tỉnh → huyện → xã)
Địa bàn ↔ thửa đất
Địa bàn ↔ khu quy hoạch
Khu quy hoạch ↔ đồ án
Đồ án ↔ văn bản pháp lý
Báo cáo ↔ data source
Mục tiêu:
Tăng crawlability
Xây knowledge graph nội bộ
Tăng semantic SEO value
A.3.8 Không SEO theo UI state

KHÔNG SEO:

popup / tooltip / drawer
filter / search result
map zoom / layer state
selection state
session/user state
Chỉ SEO khi:
Là Entity độc lập
Có URL public ổn định

👉 Tránh tạo URL rác + duplicate + instability

A.4 SEO Class chuẩn (RẤT QUAN TRỌNG CHO IMPLEMENTATION)
Mục đích

Phân loại mức độ SEO để điều khiển:

Index / NoIndex
Sitemap inclusion
Metadata level
Structured data level
AI summary level
Internal link priority
A.4.1 SEO-A (Core Entity Pages)
Vai trò:
Quan trọng nhất hệ thống SEO
Bao gồm:
Trang chi tiết Entity:
đất
khu quy hoạch
đồ án
văn bản pháp lý
bản đồ
báo cáo
hành chính
Rule:
ALWAYS:
Index
Sitemap
Canonical
Metadata full
Structured Data
OpenGraph
AI Summary
Internal link đầy đủ

👉 Đây là main SEO revenue pages

A.4.2 SEO-B (Navigation / Listing pages)
Vai trò:
Trang trung gian điều hướng
Bao gồm:
danh sách:
đồ án
văn bản
khu quy hoạch
bản đồ
hành chính
Rule:
Có thể Index
Có thể vào Sitemap
Có Metadata
Có Internal Link mạnh
Dẫn về SEO-A

👉 Vai trò: hỗ trợ crawl + discovery

A.4.3 SEO-C (Conditional SEO pages)
Vai trò:
SEO có điều kiện (dynamic / uncertain)
Bao gồm:
search result
GPS lookup
tờ/thửa lookup
vùng phân tích
so sánh dữ liệu
trang tổng hợp động
Rule:
Default: NoIndex
Chỉ Index khi:
đủ chất lượng dữ liệu
có Entity rõ ràng
có giá trị ổn định
Risk:
dễ tạo nhiều URL → phải kiểm soát canonical + index chặt

👉 Đây là nhóm dễ gây SEO noise nếu không kiểm soát

TÓM TẮT CHO DEV (RẤT NGẮN GỌN)

Nếu dùng để implement hệ thống:

SEO = Entity-driven system
Canonical = 1 entity = 1 URL
Default = NoIndex
Index chỉ khi đủ điều kiện + entity rõ
Internal link = graph giữa entities
UI state = KHÔNG BAO GIỜ SEO
SEO Class:
A = entity page (core index)
B = listing/navigation
C = dynamic/conditional

A.4.4 SEO-D – Share / Discover (Share Layer)
Vai trò
Nhóm URL phục vụ:
Share
Snapshot
Preview
Discover
Không phải landing page SEO chính
Rule kỹ thuật
Có URL riêng
BẮT BUỘC:
Open Graph
Metadata tối thiểu (title, description, image)
Có thể:
NoIndex = true
Canonical → Entity gốc (nếu tồn tại)
Không đưa vào SEO ranking system
Mapping logic
if page.type == SHARE || SNAPSHOT:
    seoClass = SEO-D
    index = false (default)
    allowOpenGraph = true
    canonical = entityUrl (if exists)
A.4.5 SEO-N – Non SEO (System / Internal only)
Vai trò
Toàn bộ URL không phục vụ search engine
Rule bắt buộc
index = false
sitemap = false
structuredData = false
aiSummary = false
internalLinkSEO = false
Danh sách nhóm loại trừ
Auth:
login, register
user profile
System:
admin, config
internal API
UI state:
popup, drawer, tooltip
map state (zoom, layer, selection)
Runtime:
error pages
temporary states
Mapping logic
if page.isSystem || page.isUIState || page.isPrivate:
    seoClass = SEO-N
    index = false
=========================
PHẦN B – SEO MATRIX SYSTEM
=========================
B.1 Mục tiêu hệ thống
Purpose

Rà soát toàn bộ system để:

xác định URL nào được SEO
xác định Entity chính
gán SEO Class
quyết định index/sitemap/canonical
B.1 Rule chuẩn cho mọi screen

Mỗi screen phải có:

screen:
  seoEligible: boolean
  primaryEntity: Entity
  secondaryEntity: [Entity]
  seoClass: SEO-A | SEO-B | SEO-C | SEO-D | SEO-N
  url:
    exists: boolean
    indexable: boolean
    inSitemap: boolean
  notes: string
=========================
QH MODULE SEO CLASSIFICATION
=========================
B.1.1 QH.1 – Map / Tra cứu quy hoạch
QH.1.1 Map View
seoEligible: true
primaryEntity: PlanningRegion
seoClass: SEO-B
index: true
sitemap: true
QH.1.2 Search / Locate UI
seoEligible: false
seoClass: SEO-N
index: false
QH.1.3 Layer Panel
seoEligible: false
seoClass: SEO-N
index: false
QH.1.4 Popup Info
seoEligible: false
seoClass: SEO-N
index: false
QH.1.5 Detail Page (IMPORTANT)
seoEligible: true
primaryEntity:
  - Parcel
  - PlanningRegion
  - PlanningProject

seoClass: SEO-A
index: true
sitemap: true
=========================
B.1.2 QH.2 – Search / Location Engine
QH.2.1 Address Search (HIGH VALUE SEO)
seoEligible: true
primaryEntity: AdministrativeUnit
seoClass: SEO-C
index: conditional
QH.2.2 GPS Lookup
seoEligible: false
seoClass: SEO-N
index: false
QH.2.3 Parcel Search (VERY IMPORTANT)
seoEligible: true
primaryEntity: Parcel
seoClass: SEO-A
index: true
sitemap: true
QH.2.4 Polygon / Area Search
seoEligible: conditional
primaryEntity: PlanningRegion
seoClass: SEO-C
index: conditional
=========================
B.1.3 QH.3 – GIS Layer System
Rule
Mostly NOT SEO
Only public semantic layers are SEO-B / SEO-C
seoEligible: conditional
primaryEntity: GISLayer
seoClass: SEO-B | SEO-C
index: conditional
=========================
B.1.4 QH.4 – Analysis / Explanation Engine (HIGH VALUE)
seoEligible: true
primaryEntity:
  - PlanningRegion
  - Parcel

seoClass: SEO-A
index: true
sitemap: true
Note (implementation important)
This module generates:
AI Summary pages
semantic content pages
landing pages
=========================
B.1.5 QH.5 – Comparison Engine
seoEligible: conditional
primaryEntity: PlanningRegion
seoClass: SEO-C
index: conditional
Rule
Only persistent comparisons are indexable
Temporary user comparisons = SEO-N behavior
=========================
B.1.6 QH.6 – Alerts / Tracking
seoEligible: false
seoClass: SEO-N
index: false
Rule
User-specific system → never SEO
=========================
B.1.7 QH.7 – Snapshot / Report
Split behavior
Public Report
seoClass: SEO-C
index: conditional
Share Snapshot
seoClass: SEO-D
index: false
Private Snapshot
seoClass: SEO-N
=========================
B.1.8 QH.8 – Account / Billing
seoClass: SEO-N
index: false
B.1.9 QH.9 – API / Integration
seoClass: SEO-N
index: false
Exception rule
API docs portal (future):
→ có thể SEO-B (separate system)
=========================
FINAL SYSTEM RULE (IMPORTANT FOR CODE)
=========================
SEO Decision Engine
function resolveSEO(page):

if page.isSystem || page.isAuth || page.isPrivate:
   return SEO-N
if page.isShareOrSnapshot:
   return SEO-D
if page.isDynamicSearchResult:
   return SEO-C (default noindex)
if page.isEntityDetail:
   return SEO-A
if page.isListing:
   return SEO-B

return SEO-N

=========================
B.1.10 QH.10 – Admin System
=========================
SEO Rule
seoClass: SEO-N
index: false
sitemap: false
structuredData: false
aiSummary: false
internalLinkSEO: false
Scope
Admin / internal system only
Mapping logic
if module == QH.10:
    return SEO-N
Global Rule (Admin)
Always:
noindex
nofollow
no sitemap
no AI content generation
no structured data
=========================
B.1 SUMMARY – MODULE SEO MATRIX
=========================
Module classification (normalized)
QH.3: SEO-C (conditional)
QH.4: SEO-A (core content engine)
QH.5: SEO-C (conditional)
QH.6: SEO-N
QH.7: SEO-C / SEO-D
QH.8: SEO-N
QH.9: SEO-N
QH.10: SEO-N
QH.11: SEO-A / SEO-B (core SEO hub)
System rule
if module in NON_SEO_MODULES:
    seoClass = SEO-N
=========================
REQUIRED SEO MODULES (PHẦN D)
=========================
Must implement deep SEO spec for:
QH.3 (GIS Layer)
QH.4 (Analysis Engine)
QH.5 (Comparison Engine)
QH.7 (Snapshot/Report)
QH.11 (Library Core)
=========================
B.1.11 QH.11 – PLANNING LIBRARY (CORE SEO HUB)
=========================
Role definition
Primary SEO content hub
Long-term stable entities
High authority pages
Knowledge Graph backbone
Entity types
coreEntities:
  - PlanningProject
  - LegalDocument
  - PlanningMap
  - PlanningLibrary
  - AdministrativeUnit
=========================
B.1.11.1 Workspace Library
=========================
seoClass: SEO-B
entity: PlanningLibrary
index: true
sitemap: true
type: landing_page
=========================
B.1.11.2 Project List
=========================
seoClass: SEO-B
entity: PlanningProject
index: true
sitemap: true
type: listing_page
=========================
B.1.11.3 Project Detail (CRITICAL)
=========================
seoClass: SEO-A
entity: PlanningProject
index: true
sitemap: true
type: entity_page
Required outputs:
metadata
structured data
ai summary
internal links
KG relations
=========================
B.1.11.4 Legal Document List
=========================
seoClass: SEO-B
entity: LegalDocument
index: true
sitemap: true
=========================
B.1.11.5 Legal Document Detail
=========================
seoClass: SEO-A
entity: LegalDocument
index: true
sitemap: true
=========================
B.1.11.6 Map List
=========================
seoClass: SEO-B
entity: PlanningMap
index: true
sitemap: true
=========================
B.1.11.7 Map Detail (HIGH SEO)
=========================
seoClass: SEO-A
entity: PlanningMap
index: true
sitemap: true
=========================
B.1.11.8 Timeline Legal
=========================
seoClass: SEO-C
index: conditional
canonical: true (to Project or LegalDocument)
=========================
B.1.11.9 Version Relations
=========================
seoClass: SEO-C
index: conditional
purpose:
  - historical seo
  - entity evolution graph
  - knowledge graph linking
=========================
QH.11 PRIORITY RULE
=========================
High SEO priority pages:
1. PlanningProject detail (SEO-A)
2. LegalDocument detail (SEO-A)
3. PlanningMap detail (SEO-A)
4. AdministrativeUnit pages (SEO-A)
5. Analysis Engine pages (QH.4)
=========================
CORE SEO CLUSTER (SYSTEM LEVEL)
=========================
Top SEO value clusters:
QH.11 Library
QH.1.5 Entity Detail Pages
QH.2.1 Address Search
QH.2.3 Parcel Search
QH.4 Analysis Engine
=========================
B.2 ENTITY SEO MATRIX
=========================
Purpose (normalized)

Entity-level SEO configuration for system-wide generation:

URL generation
metadata generation
structured data
AI summary
knowledge graph
internal linking
Entity SEO contract (STANDARD MODEL)
EntitySEO:
  id: string
  type: string
  seoEligible: boolean
  seoClass: SEO-A | SEO-B | SEO-C | SEO-D | SEO-N
  url:
    enabled: boolean
    indexable: boolean
    canonical: string
  features:
    metadata: boolean
    structuredData: boolean
    aiSummary: boolean
    internalLinks: boolean
    sitemap: boolean
=========================
B.2.1 Administrative Unit (CORE ENTITY)
=========================
Role
Root node of all geo + planning data
Central hub for SEO + KG
SEO Configuration
entity: AdministrativeUnit
seoClass: SEO-A
priority: VERY_HIGH
index: true
sitemap: true
aiSearch: true
kgNode: true
Hierarchy model
Country
  → Province/City
    → District
      → Ward
        → Historical Units (pre-merger)
SEO Role mapping
Must support:
Landing pages per administrative unit
Entry point SEO pages
Knowledge Graph root nodes
Internal link hubs
URL behavior
/admin/{type}/{slug}
SEO rules
Always index
Always canonicalize
Always generate structured data
Always generate KG links
=========================
FINAL SYSTEM BEHAVIOR (IMPORTANT)
GLOBAL SEO ENGINE RULE
function resolveEntitySEO(entity):

  if entity.type in SYSTEM_ENTITIES:
      return SEO-N

  if entity.type == SHARE or SNAPSHOT:
      return SEO-D

  if entity.isDynamic:
      return SEO-C

  if entity.isCoreDetail:
      return SEO-A

  if entity.isListing:
      return SEO-B

  return SEO-N
=========================
OUTPUT SUMMARY (FOR IMPLEMENTATION)
System architecture rules:
SEO is Entity-driven
QH.11 = main SEO hub
QH.4 = AI content engine (high SEO value)
QH.1.5 + QH.2.3 = main entry SEO pages
Default rule = NoIndex unless explicitly allowed
Admin/API/System = always SEO-N
Share layer = SEO-D only (OG + preview)
Entity (AdministrativeUnit) = root of Knowledge Graph

B.2.2 Parcel (Thửa đất) — Entity SEO Specification
B.2.2.1 Entity Definition
entity:
  name: Parcel
  type: real_estate_land_parcel
  description: "Thửa đất là đơn vị dữ liệu đất đai nhỏ nhất trong hệ thống QH Pro, đại diện cho một polygon địa lý có thông tin pháp lý và quy hoạch."
B.2.2.2 SEO Attributes
seo:
  enabled: true
  priority: highest
  class: SEO-A

  url:
    enabled: true
    type: detail_url
    pattern: "/parcel/{parcel_id}"

  index:
    enabled: true

  sitemap:
    enabled: true

  ai_search:
    enabled: true

  knowledge_graph:
    enabled: true
B.2.2.3 Functional Roles (SEO Use Cases)

Parcel được sử dụng để sinh các loại trang sau:

use_cases:
  - parcel_detail_page
  - parcel_lookup_page
  - parcel_planning_info_page
  - parcel_analysis_page
B.2.2.4 Entity Relationships
relations:
  parent:
    - AdministrativeUnit

  linked_entities:
    - PlanningRegion
    - PlanningProject
    - LegalDocument
    - Report
    - Snapshot
B.2.2.5 SEO Content Model (Output Contract)
Input (Backend → SEO Engine)
{
  "parcel_id": "string",
  "geometry": "polygon",
  "admin_unit_id": "string",
  "planning_region_ids": ["string"],
  "legal_documents": ["string"],
  "metadata": {
    "address": "string",
    "map_code": "string",
    "area": "number"
  }
}
Output (SEO Page Model)
page:
  type: parcel_detail

  title: "Thửa đất {parcel_id} - {address}"
  h1: "Thông tin thửa đất {parcel_id}"

  meta:
    description: "Thông tin quy hoạch, pháp lý và bản đồ thửa đất {parcel_id}"
    robots: "index,follow"

  breadcrumb:
    - "Trang chủ"
    - "Quy hoạch"
    - "Thửa đất {parcel_id}"

  ai_summary:
    enabled: true
    sources:
      - Parcel
      - PlanningRegion
      - LegalDocument

  key_facts:
    - parcel_id
    - area
    - administrative_unit
    - zoning_status
    - legal_status

  internal_links:
    - administrative_unit_page
    - planning_region_page
    - planning_project_page
B.2.2.6 Indexing Rules
index_rules:
  allow_index:
    - parcel.is_public == true
    - parcel.has_valid_geometry == true

  deny_index:
    - parcel.is_private == true
    - parcel.status == "internal"
B.2.2.7 AI / Knowledge Graph Rules
ai_kg:
  node_type: Parcel

  properties:
    - parcel_id
    - geometry
    - area
    - legal_status

  edges:
    - BELONGS_TO -> AdministrativeUnit
    - LOCATED_IN -> PlanningRegion
    - AFFECTED_BY -> PlanningProject
B.2.3 Planning Region (Khu quy hoạch)
B.2.3.1 Entity Definition
entity:
  name: PlanningRegion
  type: spatial_planning_area
  description: "Vùng quy hoạch đại diện cho khu vực áp dụng quy hoạch trong hệ thống GIS."
B.2.3.2 SEO Attributes
seo:
  enabled: true
  priority: highest
  class: SEO-A

  url:
    type: detail_url
    pattern: "/planning-region/{region_id}"

  index: true
  sitemap: true
  ai_search: true
  knowledge_graph: true
B.2.3.3 Functional Roles
use_cases:
  - planning_region_landing_page
  - planning_region_detail_page
  - planning_analysis_page
B.2.3.4 Entity Relations
relations:
  parent:
    - AdministrativeUnit

  children:
    - Parcel

  linked_entities:
    - PlanningProject
    - GISLayer
    - LegalDocument
B.2.4 Planning Project (Đồ án quy hoạch)
B.2.4.1 Entity Definition
entity:
  name: PlanningProject
  type: planning_document_project
  description: "Đồ án quy hoạch đại diện cho một kế hoạch quy hoạch được phê duyệt hoặc điều chỉnh."
B.2.4.2 SEO Attributes
seo:
  enabled: true
  priority: highest
  class: SEO-A

  url:
    pattern: "/planning-project/{project_id}"

  index: true
  sitemap: true
  ai_search: true
  knowledge_graph: true
B.2.4.3 Functional Roles
use_cases:
  - planning_project_landing_page
  - planning_project_detail_page
  - planning_project_document_view
B.2.4.4 Relations
relations:
  parent:
    - AdministrativeUnit

  linked_entities:
    - PlanningRegion
    - LegalDocument
    - PlanningMap
B.2.5 Planning Library (Thư viện quy hoạch)
B.2.5.1 Entity Definition
entity:
  name: PlanningLibrary
  type: knowledge_repository
  description: "Kho tri thức quy hoạch tập hợp các đồ án, văn bản, bản đồ."
B.2.5.2 SEO Attributes
seo:
  enabled: true
  class: SEO-B
  index: true
  sitemap: true
B.2.5.3 Role
use_cases:
  - library_homepage
  - content_hub_page
  - navigation_page
B.2.6 Legal Document (Văn bản pháp lý)
B.2.6.1 Entity Definition
entity:
  name: LegalDocument
  type: legal_planning_document
B.2.6.2 SEO Attributes
seo:
  enabled: true
  class: SEO-A
  priority: highest
  index: true
  sitemap: true
  ai_search: true
  knowledge_graph: true
B.2.6.3 Use Cases
use_cases:
  - legal_document_detail_page
  - legal_document_list_page
  - legal_document_reference_page
B.2.6.4 Relations
relations:
  linked_entities:
    - PlanningProject
    - AdministrativeUnit
B.2.7 GIS Layer
entity:
  name: GISLayer
  type: spatial_data_layer
seo:
  conditional: true
  rules:
    include_if:
      - layer.is_public == true
      - layer.type in ["planning", "official", "published"]
    exclude_if:
      - layer.type in ["internal", "system", "temporary"]
B.2.8 Report
entity:
  name: Report
  type: analytical_report
seo:
  class: SEO-C
  conditional_indexing: true
B.2.9 Snapshot
entity:
  name: Snapshot
  type: temporal_state_record
seo:
  class: SEO-D
  default_index: false
  allow_index_if_public: true
B.2.10 AI Summary Content (Derived Entity)
entity:
  name: AISummaryContent
  type: derived_semantic_content
  source: [Parcel, PlanningRegion, PlanningProject, LegalDocument]
seo:
  index: false
  reason: "derived_content_not_primary_entity"
  ai_search: true
  knowledge_graph: true
KẾT LUẬN (Machine Rule Summary)
Entity Priority Map
SEO_PRIORITY:
  SEO-A:
    - Parcel
    - PlanningRegion
    - PlanningProject
    - LegalDocument

  SEO-B:
    - PlanningLibrary
    - GISLayer (public)
    - List pages

  SEO-C:
    - Report
    - Compare views
    - Search results

  SEO-D:
    - Snapshot
    - Share URL

  SEO-N:
    - API
    - Admin
    - Internal systems

