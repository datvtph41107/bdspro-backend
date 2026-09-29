# QH.11 Simplified Domain Model

## Entity 1: QHPlanningProject

Đại diện cho một Đồ án Quy hoạch.

### Fields

| Field          | Type     | Required |
| -------------- | -------- | -------- |
| id             | UUID     | Yes      |
| code           | String   | Yes      |
| name           | String   | Yes      |
| planningType   | Enum     | Yes      |
| planningLevel  | Enum     | Yes      |
| authority      | String   | No       |
| jurisdictionId | UUID     | No       |
| summary        | Text     | No       |
| approvalDate   | Date     | No       |
| effectiveDate  | Date     | No       |
| expiryDate     | Date     | No       |
| validityStatus | Enum     | Yes      |
| currentVersion | String   | No       |
| metadata       | JSONB    | No       |
| createdAt      | DateTime | Yes      |
| updatedAt      | DateTime | Yes      |

---

### PlanningType

```text
MASTER_PLAN
ZONING_PLAN
DETAILED_PLAN
LAND_USE_PLAN
SECTOR_PLAN
INFRASTRUCTURE_PLAN
TRANSPORT_PLAN
OTHER
```

---

### PlanningLevel

```text
NATIONAL
REGIONAL
PROVINCIAL
DISTRICT
COMMUNE
SPECIAL_ZONE
```

---

### ValidityStatus

```text
DRAFT
UNDER_REVIEW
APPROVED
ACTIVE
PARTIALLY_ADJUSTED
REPLACED
EXPIRED
UNKNOWN
```

---

## Entity 2: QHPlanningDocument

Đại diện cho mọi thành phần thuộc Hồ sơ Quy hoạch.

### Fields

| Field             | Type     | Required |
| ----------------- | -------- | -------- |
| id                | UUID     | Yes      |
| planningProjectId | UUID     | Yes      |
| documentType      | Enum     | Yes      |
| code              | String   | No       |
| title             | String   | Yes      |
| description       | Text     | No       |
| fileId            | UUID     | No       |
| versionNo         | String   | No       |
| validityStatus    | Enum     | No       |
| issueDate         | Date     | No       |
| effectiveDate     | Date     | No       |
| metadata          | JSONB    | No       |
| createdAt         | DateTime | Yes      |
| updatedAt         | DateTime | Yes      |

---

### DocumentType

```text
LEGAL_DOCUMENT

DECISION
RESOLUTION
NOTICE

EXPLANATION
REPORT
APPENDIX

MAP
DRAWING

CAD_FILE
PDF_FILE
IMAGE_FILE

GIS_SOURCE_DATA

ATTACHMENT
```

---

### DocumentSubType (Optional)

Dùng JSON hoặc enum phụ nếu cần chi tiết.

Ví dụ:

GIS_SOURCE_DATA

```text
SHP
GEOJSON
GEOPACKAGE
GDB
WMS
WFS
RASTER
```

MAP

```text
LAND_USE_MAP
TRANSPORT_MAP
INFRASTRUCTURE_MAP
ZONING_MAP
GENERAL_PLAN_MAP
```

---

## Quan hệ

QHPlanningProject

```text
1 Project
    ↓
N Documents
```

QHPlanningDocument

```text
N Documents
    ↓
1 Project
```

---

## Không tạo bảng riêng

Các khái niệm sau chỉ là metadata hoặc enum:

* Hồ sơ Quy hoạch
* Văn bản pháp lý
* Quyết định
* Nghị quyết
* Thông báo
* Thuyết minh
* Báo cáo
* Phụ lục
* Bản đồ
* Bản vẽ
* CAD
* PDF
* Ảnh
* GIS Source Data

---

## Các bảng phụ độc lập nếu cần

QHPlanningTimeline

QHPlanningVersion

QHPlanningRelation

QHPlanningNewsLink

Không lưu chung vào QHPlanningDocument.
