# TQD Service HTTP Handlers - Complete Analysis

## Overview

This document provides a comprehensive analysis of all HTTP handlers in the TQD service that need to be converted to gRPC handlers. The analysis is based on the files in `/tqd-service/infra/handler/http/`.

## Handler Structure Pattern

### Common Structure

Each handler follows this pattern:

```go
type XxxHandler struct {
    xxxUsecase *usecase.XxxUsecase
    mapper     *mapper.XxxMapper
    validator  *validator.XxxValidator  // Optional
}

func NewXxxHandler(
    xxxUsecase *usecase.XxxUsecase,
    mapper *mapper.XxxMapper,
    validator *validator.XxxValidator,
) *XxxHandler {
    return &XxxHandler{
        xxxUsecase: xxxUsecase,
        mapper: mapper,
        validator: validator,
    }
}
```

### Common HTTP Method Signatures

```go
func (h *XxxHandler) XxxMethodHTTP(c *gin.Context) {
    // Parse parameters/body
    // Validate
    // Call usecase
    // Return response
}
```

---

## 1. AmenityHandler

**File**: `amenity_http_handler.go`

### Handler struct

```go
type AmenityHandler struct {
    amenityUsecase *usecase.AmenityUsecase
    amenityMapper  *mapper.AmenityMapper
}
```

### HTTP Methods Implemented

| Method            | HTTP Verb | Route                  | Request DTO             | Response                 | Description          |
| ----------------- | --------- | ---------------------- | ----------------------- | ------------------------ | -------------------- |
| CreateAmenityHTTP | POST      | /v2/tqd/amenities      | CreateAmenityRequestDTO | AmenityDTO               | Create new amenity   |
| GetAmenityHTTP    | GET       | /v2/tqd/amenities/{id} | Path: id (uint64)       | AmenityDTO               | Get amenity by ID    |
| UpdateAmenityHTTP | PUT       | /v2/tqd/amenities/{id} | UpdateAmenityRequestDTO | AmenityDTO               | Update amenity by ID |
| DeleteAmenityHTTP | DELETE    | /v2/tqd/amenities/{id} | Path: id (uint64)       | SubmitResponse           | Delete amenity       |
| ListAmenitiesHTTP | GET       | /v2/tqd/amenities      | Query: page, size       | ListAmenitiesResponseDTO | List with pagination |

**gRPC Conversion Notes**:

- Request IDs from path params should become request message fields
- Query pagination params → request message
- Response codes standardized via gRPC status codes

---

## 2. ContactLabelHandler

**File**: `contact_label_handler.go`

**Note**: This handler is already hybrid (has both gRPC and HTTP methods)

### Handler struct

```go
type ContactLabelHandler struct {
    tqdpb.UnimplementedContactLabelServiceServer
    contactLabelUsecase usecase.ContactLabelUsecase
    mapper              *mapper.ContactLabelMapper
    validator           *validator.ContactLabelValidator
}
```

### gRPC Methods Implemented (Already Defined)

| Method             | Request                         | Response                        | Description                  |
| ------------------ | ------------------------------- | ------------------------------- | ---------------------------- |
| CreateContactLabel | tqdpb.ContactLabel              | tqdpb.ContactLabel              | Create new contact label     |
| GetContactLabel    | tqdpb.GetContactLabelRequest    | tqdpb.ContactLabel              | Get by ID                    |
| UpdateContactLabel | tqdpb.ContactLabel              | tqdpb.ContactLabel              | Update contact label         |
| DeleteContactLabel | tqdpb.DeleteContactLabelRequest | sharepb.SubmitResponse          | Delete contact label         |
| ListContactLabels  | tqdpb.ListContactLabelsRequest  | tqdpb.ListContactLabelsResponse | List with pagination/filters |

### HTTP Methods Also Implemented

| Method                                | HTTP Verb | Request            | Response           |
| ------------------------------------- | --------- | ------------------ | ------------------ |
| CreateContactLabelHTTP                | POST      | tqdpb.ContactLabel | tqdpb.ContactLabel |
| (More HTTP methods partially visible) |           |                    |                    |

**gRPC Conversion Notes**:

- Already uses proto messages
- Status handling with `status.Errorf(codes.XXX, ...)`
- Proto message field conversion pattern established

---

## 3. DirectoryCategoryHandler

**File**: `directory_category_http_handler.go`

### Handler struct

```go
type DirectoryCategoryHandler struct {
    directoryCategoryUsecase *usecase.DirectoryCategoryUsecase
    mapper                   *mapper.DirectoryCategoryMapper
    validator                *validator.DirectoryCategoryValidator
}
```

### HTTP Methods Implemented

| Method                      | HTTP Verb | Route                             | Request                             | Response                                      | Description          |
| --------------------------- | --------- | --------------------------------- | ----------------------------------- | --------------------------------------------- | -------------------- |
| CreateDirectoryCategoryHTTP | POST      | /v2/tqd/directory/categories      | domain.DirectoryCategory            | domain.DirectoryCategory                      | Create new category  |
| GetDirectoryCategoryHTTP    | GET       | /v2/tqd/directory/categories/{id} | Path: id (uint64)                   | domain.DirectoryCategory (via mapper.ToProto) | Get by ID            |
| UpdateDirectoryCategoryHTTP | PUT       | /v2/tqd/directory/categories/{id} | domain.DirectoryCategory            | domain.DirectoryCategory                      | Update by ID         |
| DeleteDirectoryCategoryHTTP | DELETE    | /v2/tqd/directory/categories/{id} | Path: id (uint64)                   | map[string]interface{}                        | Delete               |
| ListDirectoryCategoriesHTTP | GET       | /v2/tqd/directory/categories      | Query: page, size, search, isActive | map[string]interface{} with data array        | List with pagination |

**Filter Parameters**:

- page (default: 0)
- size (default: 10)
- search (string)
- isActive (optional bool)

---

## 4. DirectorySourceHandler

**File**: `directory_source_http_handler.go`

### Handler struct

```go
type DirectorySourceHandler struct {
    directorySourceUsecase *usecase.DirectorySourceUsecase
    mapper                 *mapper.DirectorySourceMapper
    validator              *validator.DirectorySourceValidator
}
```

### HTTP Methods Implemented

| Method                    | HTTP Verb | Route                   | Request                | Response               | Description       |
| ------------------------- | --------- | ----------------------- | ---------------------- | ---------------------- | ----------------- |
| CreateDirectorySourceHTTP | POST      | /directory/sources      | DirectorySourceDTO     | Proto via mapper       | Create new source |
| GetDirectorySourceHTTP    | GET       | /directory/sources/{id} | Path: id (uint64)      | Proto via mapper       | Get by ID         |
| UpdateDirectorySourceHTTP | PUT       | /directory/sources/{id} | domain.DirectorySource | Proto via mapper       | Update by ID      |
| DeleteDirectorySourceHTTP | DELETE    | /directory/sources/{id} | Path: id (uint64)      | sharepb.SubmitResponse | Delete            |
| ListDirectorySourcesHTTP  | GET       | /directory/sources      | Query filters          | map[string]interface{} | List with filters |

**Filter Parameters**:

- page, size (pagination)
- search (string)
- category (string)
- type (string)
- isActive (optional bool)
- isRecurring (optional bool)

---

## 5. DirectorySupplierHandler

**File**: `directory_supplier_http_handler.go`

### Handler struct

```go
type DirectorySupplierHandler struct {
    directorySupplierUsecase *usecase.DirectorySupplierUsecase
    mapper                   *mapper.DirectorySupplierMapper
}
```

### HTTP Methods Implemented

| Method                      | HTTP Verb | Route                            | Request                  | Response                      | Description         |
| --------------------------- | --------- | -------------------------------- | ------------------------ | ----------------------------- | ------------------- |
| CreateDirectorySupplierHTTP | POST      | /v2/tqd/directory/suppliers      | domain.DirectorySupplier | JSON response object          | Create new supplier |
| GetDirectorySupplierHTTP    | GET       | /v2/tqd/directory/suppliers/{id} | Path: id (uint64)        | JSON response object          | Get by ID           |
| UpdateDirectorySupplierHTTP | PUT       | /directory/suppliers/{id}        | domain.DirectorySupplier | JSON response object          | Update by ID        |
| DeleteDirectorySupplierHTTP | DELETE    | /v2/tqd/directory/suppliers/{id} | Path: id (uint64)        | JSON response object          | Delete              |
| ListDirectorySuppliersHTTP  | GET       | /v2/tqd/directory/suppliers      | Query filters            | JSON response with pagination | List with filters   |

**Filter Parameters**:

- page, size (default: 1, 10)
- search (string)
- categories (comma-separated string)
- isActive (optional bool)
- minRating, maxRating (optional float64)

---

## 6. OpenHourHandler

**File**: `open_hour_http_handler.go`

### Handler struct

```go
type OpenHourHandler struct {
    openHourUsecase usecase.OpenHourUsecase
}
```

### HTTP Methods Implemented

| Method                    | HTTP Verb | Route                           | Request                   | Response                          | Description          |
| ------------------------- | --------- | ------------------------------- | ------------------------- | --------------------------------- | -------------------- |
| CreateOpenHourHTTP        | POST      | /v2/tqd/open-hours              | dto.OpenHourCreateRequest | JSON response with created data   | Create new open hour |
| GetOpenHourHTTP           | GET       | /v2/tqd/open-hours/{id}         | Path: id (uint64)         | JSON response with open hour data | Get by ID            |
| UpdateOpenHourHTTP        | PUT       | /v2/tqd/open-hours/{id}         | dto.OpenHourUpdateRequest | JSON response with updated data   | Update by ID         |
| DeleteOpenHourHTTP        | DELETE    | /v2/tqd/open-hours/{id}         | Path: id (uint64)         | JSON response                     | Delete (soft)        |
| ListOpenHoursHTTP         | GET       | /v2/tqd/open-hours              | Query filters             | JSON array with total             | List with pagination |
| GetOpenHourTimesByDayHTTP | GET       | /v2/tqd/open-hours/times-by-day | Query: poiId              | JSON response                     | Get grouped by day   |

**Filter Parameters for List**:

- page, size (pagination)
- poiId (uint64)
- dayOfWeek (2-8 range)
- isOpen (optional bool)
- status (int)
- type (int)

**Request DTOs**:

- `OpenHourCreateRequest`: Name, POIID, DayOfWeek, OpenTime, CloseTime, Note, IsOpen, Status, Type
- `OpenHourUpdateRequest`: All fields optional (using pointers)

---

## 7. PoiCategoryHandler

**File**: `poi_category_http_handler.go`

### Handler struct

```go
type PoiCategoryHandler struct {
    poiCategoryUsecase usecase.PoiCategoryUsecase
}
```

### HTTP Methods Implemented

| Method                 | HTTP Verb | Route                       | Request            | Response             | Description             |
| ---------------------- | --------- | --------------------------- | ------------------ | -------------------- | ----------------------- |
| CreatePoiCategoryHTTP  | POST      | /v2/tqd/poi-categories      | domain.PoiCategory | JSON response        | Create new POI category |
| GetPoiCategoryHTTP     | GET       | /v2/tqd/poi-categories/{id} | Path: id (uint64)  | JSON response        | Get by ID               |
| UpdatePoiCategoryHTTP  | PUT       | /v2/tqd/poi-categories/{id} | domain.PoiCategory | JSON response        | Update by ID            |
| DeletePoiCategoryHTTP  | DELETE    | /v2/tqd/poi-categories/{id} | Path: id (uint64)  | JSON response        | Delete                  |
| ListPoiCategoriesHTTP  | GET       | /v2/tqd/poi-categories      | Query filters      | JSON with data array | List with pagination    |
| GetPoiCategoryTreeHTTP | GET       | /v2/tqd/poi-categories/tree | -                  | JSON response        | Get hierarchical tree   |

**Filter Parameters for List**:

- page, size (pagination)
- code (string)
- name (string)
- parentId (uint64)
- isActive (optional bool)

**Request/Response DTO**:

- `PoiCategoryListRequest`: query binding structure

---

## 8. POIHTTPHandler

**File**: `poi_http_handler.go`

### Handler struct

```go
type POIHTTPHandler struct {
    poiUsecase usecase.IPOIUsecase
}
```

### HTTP Methods Implemented

| Method            | HTTP Verb | Route                              | Request                        | Response                | Description          |
| ----------------- | --------- | ---------------------------------- | ------------------------------ | ----------------------- | -------------------- |
| CreatePOI         | POST      | /api/v1/pois                       | dto.CreatePOIRequestDTO        | dto.POIDTO              | Create new POI       |
| UpdatePOI         | PUT       | /api/v1/pois/{id}                  | dto.UpdatePOIRequestDTO        | dto.POIDTO              | Update POI           |
| DeletePOI         | DELETE    | /api/v1/pois/{id}                  | Path: id (uint64)              | 204 No Content          | Delete POI           |
| GetPOIByID        | GET       | /api/v1/pois/{id}                  | Path: id (uint64)              | dto.POIDTO              | Get by ID            |
| GetPOIByCode      | GET       | /api/v1/pois/code/{code}           | Path: code (string)            | dto.POIDTO              | Get by code          |
| GetPOIList        | GET       | /api/v1/pois                       | Query filters                  | dto.ListPOIsResponseDTO | List with pagination |
| GetPOIsByCategory | GET       | /api/v1/pois/category/{categoryId} | Path: categoryId, Query: limit | array of dto.POIDTO     | Get by category      |
| GetNearbyPOIs     | GET       | /api/v1/pois/nearby                | Query: lat, lng, radius, limit | array of dto.POIDTO     | Get nearby           |
| UpdatePOIRating   | PUT       | /api/v1/pois/{id}/rating           | JSON: ratingPoint, ratingCount | JSON response           | Update rating        |

**Filter Parameters for GetPOIList**:

- page, size (pagination)
- text (search by name/address/code)
- categoryId (uint64)
- isActive (optional bool)
- minLat, maxLat, minLng, maxLng (geographic bounds)

**Special Methods**:

- GetNearbyPOIs: uses coordinates + radius + limit
- UpdatePOIRating: updates ratingPoint (0-5) and ratingCount
- GetPoiCategoryTreeHTTP: returns hierarchical structure

---

## Common Patterns & DTOs Used

### Pagination Pattern

```go
type Pagable struct {
    Page uint32  // 0-indexed
    Size uint32  // Default: 10-20
}
```

### Common Request DTOs

- `CreateXxxRequestDTO`: Creation payload
- `UpdateXxxRequestDTO`: Update payload (often with optional fields using pointers)
- `ListXxxRequestDTO`: Pagination + filters

### Common Response DTOs

- `XxxDTO`: Single entity response
- `ListXxxResponseDTO`: Array response with total count
- `SubmitResponse`: Generic success response from sharepb
- `map[string]interface{}`: Generic JSON responses

### Error Handling Pattern

```go
// Consistent error codes
- http.StatusBadRequest (400): Invalid input/validation
- http.StatusNotFound (404): Resource not found
- http.StatusInternalServerError (500): Processing errors
```

---

## gRPC Conversion Mapping

### Method Naming Convention

```
HTTP Method        →  gRPC Method
CreateXxxHTTP      →  Create (request: CreateXxxRequest, response: Xxx)
GetXxxHTTP         →  Get (request: GetXxxRequest, response: Xxx)
UpdateXxxHTTP      →  Update (request: UpdateXxxRequest, response: Xxx)
DeleteXxxHTTP      →  Delete (request: DeleteXxxRequest, response: SubmitResponse)
ListXxxHTTP        →  List (request: ListXxxRequest, response: ListXxxResponse)
GetXxxByCodeHTTP   →  GetByCode (request: GetByCodeRequest, response: Xxx)
GetXxxsByCategory  →  GetsByCategory (request: GetsByCategoryRequest, response: GetsByCategoryResponse)
GetNearbyXxx       →  GetNearby (request: GetNearbyRequest, response: GetNearbyResponse)
```

### Method Signatures for gRPC

```go
// CRUD Operations
func (s *Server) Create(ctx context.Context, req *CreateXxxRequest) (*Xxx, error)
func (s *Server) GetByID(ctx context.Context, req *GetXxxRequest) (*Xxx, error)
func (s *Server) Update(ctx context.Context, req *UpdateXxxRequest) (*Xxx, error)
func (s *Server) Delete(ctx context.Context, req *DeleteXxxRequest) (*SubmitResponse, error)
func (s *Server) List(ctx context.Context, req *ListXxxRequest) (*ListXxxResponse, error)

// Advanced Operations
func (s *Server) GetByCode(ctx context.Context, req *GetByCodeRequest) (*Xxx, error)
func (s *Server) GetsByCategory(ctx context.Context, req *GetsByCategoryRequest) (*GetsByCategoryResponse, error)
func (s *Server) GetNearby(ctx context.Context, req *GetNearbyRequest) (*GetNearbyResponse, error)
```

---

## Summary of gRPC Services to Create

### 1. AmenityService

- Create, GetByID, Update, Delete, List

### 2. ContactLabelService (Already Exists)

- Create, GetByID, Update, Delete, List
- Convert HTTP methods to reuse gRPC implementation

### 3. DirectoryCategoryService

- Create, GetByID, Update, Delete, List

### 4. DirectorySourceService

- Create, GetByID, Update, Delete, List

### 5. DirectorySupplierService

- Create, GetByID, Update, Delete, List

### 6. OpenHourService

- Create, GetByID, Update, Delete, List
- GetTimesByDay (special method)

### 7. PoiCategoryService

- Create, GetByID, Update, Delete, List
- GetTree (special hierarchical method)

### 8. POIService

- Create, Update, Delete, GetByID, GetByCode
- GetList (with pagination & filters)
- GetsByCategory, GetNearby, UpdateRating (special methods)

---

## Proto Message Templates to Generate

Each service should have corresponding proto messages:

```protobuf
// Standard CRUD Request/Response messages
message CreateXxxRequest { /* fields */ }
message UpdateXxxRequest { /* fields */ }
message GetXxxRequest { uint64 id = 1; }
message DeleteXxxRequest { uint64 id = 1; }
message ListXxxRequest {
    common.Pagable pagable = 1;
    // filter fields
}
message ListXxxResponse {
    repeated Xxx data = 1;
    int64 total = 2;
}

// Special method messages
message GetByCodeRequest { string code = 1; }
message GetNearbyRequest {
    float lat = 1;
    float lng = 2;
    float radius = 3;
    int32 limit = 4;
}
```

---

## Note on Handler Initialization

All handlers use dependency injection pattern via constructor:

```go
func NewXxxHandler(
    usecase *usecase.XxxUsecase,
    mapper *mapper.XxxMapper,
    validator *validator.XxxValidator,  // Optional
) *XxxHandler
```

For gRPC conversion, ensure:

1. Same usecase injections maintained
2. Mappers updated to handle proto↔domain conversion
3. Validators integrated with gRPC validation pattern
4. Wire configuration updated to register gRPC handlers
