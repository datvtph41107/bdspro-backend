# TQD Service - Usecase Implementation Analysis

**Analysis Date:** March 16, 2026  
**Service:** TQD (Travel/Tourism/Directory Service)  
**Analysis Scope:** All files in `/tqd-service/internal/usecase/`

---

## 1. BASE_USECASE.GO - Authorization Usecase

**File:** `base_usecase.go`

### Summary

- **Status:** ✅ COMPLETE (Base functionality)
- **Type:** Authorization/Base utilities for all other usecases
- **Pattern:** Composition pattern for other usecases to inherit authorization logic

### Implemented Methods

| Method                        | Signature                                                      | Status            |
| ----------------------------- | -------------------------------------------------------------- | ----------------- |
| `NewBaseUsecase`              | `func(userProvider, permissionProvider) *AuthorizationUsecase` | ✅ Full           |
| `GetCurrentUserID`            | `func(ctx) uint64`                                             | ✅ Full           |
| `GetCurrentOrganizationID`    | `func(ctx) uint64`                                             | ✅ Full           |
| `CheckPermission`             | `func(ctx, permission) (bool, error)`                          | ✅ Full           |
| `CheckPermissionWithUserID`   | `func(ctx, userID, permission) (bool, error)`                  | ✅ Full           |
| `CheckMultiplePermissions`    | `func(ctx, permissions[]) (map[Permission]bool, error)`        | ✅ Full           |
| `UserCanAccessTarget`         | `func(ctx, targetID, targetType) (bool, error)`                | ✅ Full           |
| `UserInOwner`                 | `func(ctx, ownerID, ownerType) (bool, error)`                  | ✅ Full           |
| `HasRoleWithOwner`            | `func(ctx, ownerID, ownerType, role) (bool, error)`            | ✅ Full           |
| `IsAdmin`                     | `func(ctx) (bool, error)`                                      | ✅ Full           |
| `IsSuperAdmin`                | `func(ctx) (bool, error)`                                      | ✅ Full           |
| `RequirePermission`           | `func(ctx, permission) error`                                  | ⚠️ **INCOMPLETE** |
| `RequirePermissionWithUserID` | `func(ctx, userID, permission) error`                          | ✅ Full           |

### Issues Found

#### 1. **RequirePermission() - COMMENTED LOGIC**

```go
// Lines 117-131: The actual permission check is commented out!
// The function currently always returns nil
func (b *AuthorizationUsecase) RequirePermission(ctx context.Context, permission enums.Permission) error {
    // hasPermission, err := b.CheckPermission(ctx, permission)
    // if err != nil {
    //     return err
    // }
    // if !hasPermission {
    //     return &PermissionError{
    //         Permission: permission.String(),
    //         Message:    "Insufficient permissions",
    //     }
    // }
    return nil  // ← ALWAYS RETURNS NIL
}
```

**Problem:** Authorization bypass! All permission checks pass.

### Dependencies

- `context`
- `tqd/internal/enums`
- `tqd/internal/interface/provider` (UserProvider, PermissionProvider)

### Recommendations

- Uncomment the permission check logic in `RequirePermission()`

---

## 2. AMENITY_USECASE.GO - Amenity Usecase

**File:** `amenity_usecase.go`

### Summary

- **Status:** ✅ COMPLETE (CRUD operations)
- **Type:** Generic CRUD wrapper
- **Pattern:** Inherits from generic `BaseUsecase[Amenity, AmenityRepo]`

### Implemented Methods

| Method                     | Signature                           | Status  |
| -------------------------- | ----------------------------------- | ------- |
| `NewAmenityUsecase`        | `func(amenityRepo) *AmenityUsecase` | ✅ Full |
| **Inherited CRUD methods** | From `BaseUsecase`                  | ✅ Full |

### CRUD Operations Inherited

- `Create(ctx, entity) error`
- `GetByID(ctx, id) (*Amenity, error)`
- `Update(ctx, id, entity) error`
- `Delete(ctx, id) error`
- `List(ctx, filter) ([]Amenity, int64, error)`
- `GetAll(ctx) ([]Amenity, error)`

### Dependencies

- `common/domain/crud` (BaseUsecase)
- `tqd/internal/domain` (Amenity)
- `tqd/internal/interface/repo` (AmenityRepo)

### Status: ✅ NO ISSUES

---

## 3. CONTACT_LABEL_USECASE.GO - Contact Label Usecase

**File:** `contact_label_usecase.go`

### Summary

- **Status:** ✅ COMPLETE (Full CRUD + custom operations)
- **Type:** Domain-specific CRUD with custom business logic
- **Pattern:** Interface + implementation struct with permission checks

### Implemented Methods

| Method                   | Signature                                                                            | Status  |
| ------------------------ | ------------------------------------------------------------------------------------ | ------- |
| `NewContactLabelUsecase` | `func(repo, userProvider, permissionProvider) ContactLabelUsecase`                   | ✅ Full |
| `Create`                 | `func(ctx, req *CreateContactLabelRequestDTO) (*ContactLabelDTO, error)`             | ✅ Full |
| `GetByID`                | `func(ctx, id) (*ContactLabelDTO, error)`                                            | ✅ Full |
| `Update`                 | `func(ctx, id, req *UpdateContactLabelRequestDTO) (*ContactLabelDTO, error)`         | ✅ Full |
| `Delete`                 | `func(ctx, id) error`                                                                | ✅ Full |
| `List`                   | `func(ctx, req *ListContactLabelsRequestDTO) (*ListContactLabelsResponseDTO, error)` | ✅ Full |
| `UpdateContactCount`     | `func(ctx, id, count) error`                                                         | ✅ Full |

### Implementation Details

**Create:**

- Gets profile ID from context
- Checks CRUD permission
- Validates unique code constraint
- Creates domain object with defaults
- Converts to DTO

**GetByID:**

- Permission check (READ)
- Retrieves and converts to DTO

**Update:**

- Permission check (UPDATE)
- Gets existing record
- Checks unique code constraint (excluding current ID)
- Updates fields
- Returns updated entity

**Delete:**

- Permission check (DELETE)
- Deletes from repository

**List:**

- Permission check (READ)
- Supports filtering by search, isActive, isSystem
- Returns paginated results with total count

**UpdateContactCount:**

- Updates contact count (no permission check!)

### Issues Found

#### 1. **UpdateContactCount() - Missing Permission Check**

```go
// Line 221-227: No permission validation!
func (u *contactLabelUsecase) UpdateContactCount(ctx context.Context, id uint64, count int) error {
    err := u.contactLabelRepo.UpdateContactCount(ctx, id, count)
    if err != nil {
        return fmt.Errorf("failed to update contact count: %v", err)
    }
    return nil
}
```

**Problem:** Can be called without authorization

### Dependencies

- `common/utils`
- `context`
- `fmt`
- `tqd/internal/domain`
- `tqd/internal/dto`
- `tqd/internal/enums`
- `tqd/internal/interface/provider`
- `tqd/internal/interface/repo`

### Status: ⚠️ MINOR ISSUES (1 missing permission check)

---

## 4. DIRECTORY_CATEGORY_USECASE.GO - Directory Category Usecase

**File:** `directory_category_usecase.go`

### Summary

- **Status:** ✅ COMPLETE (Full CRUD operations)
- **Type:** Domain-specific CRUD with authorization
- **Pattern:** Inherits from AuthorizationUsecase + custom implementation

### Implemented Methods

| Method                        | Signature                                                                | Status  |
| ----------------------------- | ------------------------------------------------------------------------ | ------- |
| `NewDirectoryCategoryUsecase` | `func(repo, userProvider, permissionProvider) *DirectoryCategoryUsecase` | ✅ Full |
| `CreateDirectoryCategory`     | `func(ctx, category *DirectoryCategory) error`                           | ✅ Full |
| `GetDirectoryCategoryByID`    | `func(ctx, id) (*DirectoryCategory, error)`                              | ✅ Full |
| `UpdateDirectoryCategory`     | `func(ctx, id, category *DirectoryCategory) (*DirectoryCategory, error)` | ✅ Full |
| `DeleteDirectoryCategory`     | `func(ctx, id) error`                                                    | ✅ Full |
| `ListDirectoryCategories`     | `func(ctx, request *ListRequest) ([]DirectoryCategory, int64, error)`    | ✅ Full |

### Implementation Details

**Create:**

- Checks permission via inherited method `RequirePermission()`
- Creates domain object
- Returns error or success

**GetByID:**

- Checks permission (READ)
- Returns domain object

**Update:**

- Checks permission (UPDATE)
- Sets ID on entity
- Updates and retrieves updated version

**Delete:**

- Checks permission (DELETE)
- Deletes from repository

**List:**

- Checks permission (READ)
- Returns paginated results

### Issues Found

**NOTE:** This usecase relies on `RequirePermission()` from base_usecase.go which has the commented logic bug. All permission checks here will pass incorrectly!

### Dependencies

- `context`
- `tqd/internal/domain`
- `tqd/internal/dto`
- `tqd/internal/enums`
- `tqd/internal/interface/provider`
- `tqd/internal/interface/repo`

### Status: ⚠️ INHERITS PERMISSION BUG (from base_usecase)

---

## 5. DIRECTORY_SOURCE_USECASE.GO - Directory Source Usecase

**File:** `directory_source_usecase.go`

### Summary

- **Status:** ⚠️ PARTIALLY COMPLETE (some features incomplete)
- **Type:** Domain-specific operations with custom logic
- **Pattern:** Standard implementation with manual permission checks

### Implemented Methods

| Method                      | Signature                                                                          | Status            |
| --------------------------- | ---------------------------------------------------------------------------------- | ----------------- |
| `NewDirectorySourceUsecase` | `func(repo, userProvider, permissionProvider) *DirectorySourceUsecase`             | ✅ Full           |
| `Create`                    | `func(ctx, req *DirectorySourceDTO) (*DirectorySourceDTO, error)`                  | ✅ Full           |
| `GetByID`                   | `func(ctx, id) (*DirectorySourceDTO, error)`                                       | ✅ Full           |
| `Update`                    | `func(ctx, req *DirectorySource) (*DirectorySourceDTO, error)`                     | ⚠️ **INCOMPLETE** |
| `Delete`                    | `func(ctx, id) error`                                                              | ✅ Full           |
| `List`                      | `func(ctx, req *ListDirectorySourcesRequestDTO) ([]DirectorySource, int64, error)` | ✅ Full           |

### Issues Found

#### 1. **Update() - COMMENTED VALIDATION LOGIC**

```go
// Lines 96-109: Critical validation is commented out!
func (u *DirectorySourceUsecase) Update(ctx context.Context, req *domain.DirectorySource) (*dto.DirectorySourceDTO, error) {
    // ...
    // 3. Convert DTO to Domain
    // domainSource := req.ToDomain()  ← COMMENTED

    // 4. Validate domain
    // if err := domainSource.Validate(); err != nil {  ← COMMENTED
    //     return nil, fmt.Errorf("validation failed: %v", err)
    // }

    // 5. Call repository (no validation!)
    err = u.directorySourceRepo.Update(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to update directory source: %v", err)
    }
    // ...
}
```

**Problem:**

- No validation on update
- No DTO to domain conversion
- Accepts domain object directly instead of DTO
- Inconsistent with Create() which validates

#### 2. **Inconsistent API Design**

- `Create()` accepts `*dto.DirectorySourceDTO`
- `Update()` accepts `*domain.DirectorySource` (inconsistent!)

### Dependencies

- `context`
- `encoding/json`
- `fmt`
- `tqd/internal/domain`
- `tqd/internal/dto`
- `tqd/internal/enums`
- `tqd/internal/interface/provider`
- `tqd/internal/interface/repo`
- `common/utils`

### Status: ⚠️ INCOMPLETE (Update function has commented logic)

---

## 6. DIRECTORY_SUPPLIER_USECASE.GO - Directory Supplier Usecase

**File:** `directory_supplier_usecase.go`

### Summary

- **Status:** ⚠️ MOSTLY COMPLETE (missing permission checks)
- **Type:** Domain-specific CRUD operations
- **Pattern:** Standard implementation (but incomplete authorization)

### Implemented Methods

| Method                        | Signature                                                                                               | Status            |
| ----------------------------- | ------------------------------------------------------------------------------------------------------- | ----------------- |
| `NewDirectorySupplierUsecase` | `func(repo, userProvider, permissionProvider) *DirectorySupplierUsecase`                                | ✅ Full           |
| `Create`                      | `func(ctx, supplier *DirectorySupplier) (*DirectorySupplier, error)`                                    | ⚠️ **INCOMPLETE** |
| `GetByID`                     | `func(ctx, id) (*DirectorySupplier, error)`                                                             | ⚠️ **INCOMPLETE** |
| `Update`                      | `func(ctx, supplier *DirectorySupplier) (*DirectorySupplier, error)`                                    | ⚠️ **INCOMPLETE** |
| `Delete`                      | `func(ctx, id) error`                                                                                   | ⚠️ **INCOMPLETE** |
| `List`                        | `func(ctx, page, size, search, categories[], minRating, maxRating) ([]DirectorySupplier, int64, error)` | ⚠️ **INCOMPLETE** |

### Issues Found

#### 1. **ALL METHODS - MISSING PERMISSION CHECKS**

```go
// Example: Create() - No permission check at all!
func (u *DirectorySupplierUsecase) Create(ctx context.Context, supplier *domain.DirectorySupplier) (*domain.DirectorySupplier, error) {
    // NO PERMISSION CHECK!
    // 3. Validate domain
    if err := supplier.Validate(); err != nil {
        // ...
    }
}

// GetByID() - No permission check
func (u *DirectorySupplierUsecase) GetByID(ctx context.Context, id uint64) (*domain.DirectorySupplier, error) {
    // NO PERMISSION CHECK!
    supplier, err := u.directorySupplierRepo.GetByID(ctx, id)
}

// Update() - No permission check
func (u *DirectorySupplierUsecase) Update(ctx context.Context, supplier *domain.DirectorySupplier) (*domain.DirectorySupplier, error) {
    // NO PERMISSION CHECK!
}

// Delete() - No permission check
func (u *DirectorySupplierUsecase) Delete(ctx context.Context, id uint64) error {
    // NO PERMISSION CHECK! Only checks existence
}

// List() - No permission check
func (u *DirectorySupplierUsecase) List(ctx context.Context, page, size int, ...) ([]domain.DirectorySupplier, int64, error) {
    // NO PERMISSION CHECK!
}
```

#### 2. **Missing Authorization Logic**

The struct has `userProvider` and `permissionProvider` fields but never uses them!

#### 3. **Fields Definition Inconsistency**

Constructor comments suggest authorization should be checked but methods don't implement it.

### Dependencies

- `context`
- `fmt`
- `tqd/internal/domain`
- `tqd/internal/interface/provider`
- `tqd/internal/interface/repo`

### Status: ❌ CRITICAL - All methods missing permission checks

---

## 7. OPEN_HOUR_USECASE.GO - Open Hour Usecase

**File:** `open_hour_usecase.go`

### Summary

- **Status:** ✅ COMPLETE (CRUD + custom operations)
- **Type:** Generic CRUD + specialized operations
- **Pattern:** Interface + implementation with inherited CRUD + custom methods

### Implemented Methods

| Method                | Signature                                                                    | Status              |
| --------------------- | ---------------------------------------------------------------------------- | ------------------- |
| `NewOpenHourUsecase`  | `func(openHourRepo) OpenHourUsecase`                                         | ✅ Full             |
| **CRUD Methods**      | From `BaseUsecase`                                                           | ✅ Full (inherited) |
| `GetTimesByDayOfWeek` | `func(ctx, poiID *uint64) ([]OpenHourTimesByDayResponse, error)`             | ✅ Full             |
| `BulkSave`            | `func(ctx, req *OpenHourBulkSaveRequest) (*OpenHourBulkSaveResponse, error)` | ✅ Full             |

### Implementation Details

**GetTimesByDayOfWeek:**

- Calls repo to get grouped by day
- Maps Vietnamese day names
- Sorts by day (2-8, where 8 = Sunday)
- Returns structured response with time slots per day

**BulkSave:**

- Processes create/update items
- Handles delete IDs
- Returns detailed response with:
    - Total operations
    - Created count
    - Updated count
    - Deleted count
    - Failed count
    - Individual result items with messages

### Dependencies

- `common/domain/crud`
- `context`
- `fmt`
- `tqd/internal/domain`
- `tqd/internal/dto`
- `tqd/internal/enums`
- `tqd/internal/interface/repo`

### Status: ✅ NO ISSUES

---

## 8. POI_CATEGORY_USECASE.GO - POI Category Usecase

**File:** `poi_category_usecase.go`

### Summary

- **Status:** ✅ COMPLETE (CRUD + tree operations)
- **Type:** Generic CRUD + specialized tree operations
- **Pattern:** Interface + implementation with inherited CRUD + custom methods

### Implemented Methods

| Method                  | Signature                                          | Status              |
| ----------------------- | -------------------------------------------------- | ------------------- |
| `NewPoiCategoryUsecase` | `func(poiCategoryRepo) PoiCategoryUsecase`         | ✅ Full             |
| **CRUD Methods**        | From `BaseUsecase`                                 | ✅ Full (inherited) |
| `GetTree`               | `func(ctx) ([]PoiCategoryTreeResponse, *ErrorDTO)` | ✅ Full             |

### Implementation Details

**GetTree:**

- Retrieves all categories in tree structure
- Recursively maps parent-child relationships
- Returns hierarchical response
- Handles error with custom ErrorDTO
- Maps domain to DTO with children

### Dependencies

- `common/domain/crud`
- `common/domain/err`
- `context`
- `tqd/internal/domain`
- `tqd/internal/dto`
- `tqd/internal/interface/repo`

### Status: ✅ NO ISSUES

---

## 9. POI_USECASE.GO - Point of Interest (POI) Usecase

**File:** `poi_usecase.go`

### Summary

- **Status:** ⚠️ MOSTLY COMPLETE (some implementations incomplete/buggy)
- **Type:** Specialized business operations
- **Pattern:** Interface + implementation with custom business logic

### Implemented Methods

| Method              | Signature                                                              | Status            |
| ------------------- | ---------------------------------------------------------------------- | ----------------- |
| `NewPOIUsecase`     | `func(poiRepo, mapManager) IPOIUsecase`                                | ✅ Full           |
| `CreatePOI`         | `func(ctx, req *CreatePOIRequestDTO) (*POIDTO, *ErrorDTO)`             | ✅ Full           |
| `UpdatePOI`         | `func(ctx, id, req *UpdatePOIRequestDTO) (*POIDTO, *ErrorDTO)`         | ✅ Full           |
| `DeletePOI`         | `func(ctx, id) *ErrorDTO`                                              | ✅ Full           |
| `GetPOIByID`        | `func(ctx, id) (*POIDTO, *ErrorDTO)`                                   | ✅ Full           |
| `GetPOIByCode`      | `func(ctx, code) (*POIDTO, *ErrorDTO)`                                 | ✅ Full           |
| `GetPOIList`        | `func(ctx, req *ListPOIsRequestDTO) (*ListPOIsResponseDTO, *ErrorDTO)` | ✅ Full           |
| `GetPOIsByCategory` | `func(ctx, categoryID, limit) ([]POIDTO, *ErrorDTO)`                   | ✅ Full           |
| `GetNearbyPOIs`     | `func(ctx, lat, lng, radius, limit) ([]POIDTO, *ErrorDTO)`             | ⚠️ **INCOMPLETE** |
| `UpdatePOIRating`   | `func(ctx, id, ratingPoint, ratingCount) *ErrorDTO`                    | ✅ Full           |

### Issues Found

#### 1. **GetNearbyPOIs() - INCOMPLETE MAPPER**

```go
// Lines 173-183: Mapper logic incomplete!
func (u *POIUsecase) GetNearbyPOIs(...) ([]dto.POIDTO, *_err.ErrorDTO) {
    // ...
    pois, err := u.mm.Nearby(ctx, lat, lng, radius)

    // Old code commented out
    // pois, err := u.poiRepo.GetNearby(ctx, lat, lng, radius, limit)
    if err != nil {
        // ...
    }

    // INCOMPLETE MAPPING - Only sets Lat and Lng!
    poiDTOs := make([]dto.POIDTO, len(pois))
    for i, poi := range pois {
        poiDTO := &dto.POIDTO{
            Lat: float64(poi.Lat),  // Only these two fields!
            Lng: float64(poi.Lon),
        }
        if poiDTO != nil {  // Always true, useless check
            poiDTOs[i] = *poiDTO
        }
    }

    return poiDTOs, nil
}
```

**Problems:**

- Only maps `Lat` and `Lng` fields from response
- All other POI fields are missing (ID, Name, Code, Description, etc.)
- The `if poiDTO != nil` check is always true (useless)
- Need to call full mapper like in other methods

**Expected Fix:**

```go
// Should use complete mapper
poiDTOs := mapper.POIListToDTO(pois)
return poiDTOs, nil
```

#### 2. **CreatePOI() - Category Validation Logic Issue**

```go
// Lines 43-47: Category check validates wrong entity
if req.CategoryID != nil {
    if _, err := u.poiRepo.GetByID(ctx, *req.CategoryID); err != nil {  // ← Wrong!
        // Checking if POI exists, not category!
        // ...
    }
}
```

**Problem:** Uses `GetByID()` on POI repo instead of category repo. Should validate category exists!

#### 3. **UpdatePOI() - Same Category Validation Issue**

```go
// Lines 80-84
if req.CategoryID != nil {
    if _, err := u.poiRepo.GetByID(ctx, *req.CategoryID); err != nil {  // ← Wrong!
        return nil, &_err.ErrorDTO{
            Code:    400,
            Message: "Danh mục không tồn tại",
        }
    }
}
```

### Dependencies

- `context`
- `common/domain/err`
- `tqd/infra/mapper`
- `tqd/internal/dto`
- `tqd/internal/interface/repo`
- `tqd/map`

### Status: ⚠️ INCOMPLETE (GetNearbyPOIs incomplete mapping + category validation bugs)

---

## SUMMARY TABLE

| File                          | Status             | CRUD | Auth | Custom Methods | Issues                                            |
| ----------------------------- | ------------------ | ---- | ---- | -------------- | ------------------------------------------------- |
| base_usecase.go               | ✅ Complete        | N/A  | ✅   | 11/12 methods  | 1 critical (commented check)                      |
| amenity_usecase.go            | ✅ Complete        | ✅   | ✅   | Inherited      | None                                              |
| contact_label_usecase.go      | ⚠️ Mostly Complete | ✅   | ✅   | 7/7            | 1 minor (missing perm check)                      |
| directory_category_usecase.go | ⚠️ Partial         | ✅   | ⚠️   | 5/5            | Inherits permission bug                           |
| directory_source_usecase.go   | ⚠️ Partial         | ⚠️   | ✅   | 6/6            | 2 issues (commented validation, inconsistent API) |
| directory_supplier_usecase.go | ❌ Stubbish        | ✅   | ❌   | 6/6            | 5+ critical (no auth checks)                      |
| open_hour_usecase.go          | ✅ Complete        | ✅   | ✅   | 2/2            | None                                              |
| poi_category_usecase.go       | ✅ Complete        | ✅   | ✅   | 1/1            | None                                              |
| poi_usecase.go                | ⚠️ Partial         | ✅   | ✅   | 8/10           | 3 issues (incomplete mapper, category validation) |

---

## CRITICAL ISSUES TO FIX (Priority Order)

### TIER 1 - SECURITY CRITICAL ⚠️

1. **base_usecase.go::RequirePermission()** - Line 117-131
    - Commented permission check always returns nil
    - **Impact:** Authorization bypass for all dependent usecases

2. **directory_supplier_usecase.go** - All methods
    - Zero permission checks in any CRUD method
    - **Impact:** Complete authorization bypass for supplier operations

### TIER 2 - HIGH PRIORITY 🔴

3. **poi_usecase.go::GetNearbyPOIs()** - Line 173-183
    - Incomplete DTO mapping (only Lat/Lng)
    - **Impact:** Missing data in API responses

4. **poi_usecase.go::CreatePOI/UpdatePOI** - Lines 43-47, 80-84
    - Category validation uses wrong repository
    - **Impact:** Category validation doesn't work

5. **directory_source_usecase.go::Update()** - Line 96-109
    - Commented validation logic
    - **Impact:** No validation on updates

### TIER 3 - MEDIUM PRIORITY 🟡

6. **contact_label_usecase.go::UpdateContactCount()** - Line 221-227
    - Missing permission check
    - **Impact:** Unauthorized count updates possible

7. **directory_source_usecase.go** - API inconsistency
    - Create accepts DTO, Update accepts Domain entity
    - **Impact:** Confusing API, potential bugs

---

## INCOMPLETE IMPLEMENTATIONS CHECKLIST

### Methods Needing Implementation/Fixes:

- [ ] Fix `base_usecase.go::RequirePermission()` - uncomment logic
- [ ] Fix `contact_label_usecase.go::UpdateContactCount()` - add permission check
- [ ] Fix `directory_source_usecase.go::Update()` - uncomment validation
- [ ] Fix `directory_source_usecase.go` - make API consistent (use DTO)
- [ ] Implement `directory_supplier_usecase.go` - add permission checks to all methods
- [ ] Fix `poi_usecase.go::GetNearbyPOIs()` - complete DTO mapping
- [ ] Fix `poi_usecase.go::CreatePOI()` - use correct category validation
- [ ] Fix `poi_usecase.go::UpdatePOI()` - use correct category validation

---

## DEPENDENCY ISSUES

### Current Patterns

1. Some usecases inherit from `AuthorizationUsecase`
2. Some usecases manually check permissions
3. Some usecases have fields for permission checking but don't use them

### Recommendation

Standardize on one authorization pattern across all usecases.

---

## NEXT STEPS

1. **Immediate (Critical):** Fix authorization bypass in `base_usecase.go`
2. **Short-term:** Add missing permission checks to `directory_supplier_usecase`
3. **Medium-term:** Fix incomplete implementations in `poi_usecase` and `directory_source_usecase`
4. **Long-term:** Audit all error handling and add comprehensive logging
