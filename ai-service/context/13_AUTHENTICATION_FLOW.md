# Authentication & Authorization Flow - Deep Dive

## 📋 Overview

**Purpose**: Detailed explanation of authentication and authorization mechanisms
**Components**: JWT, Middleware, Context Propagation, RBAC
**Status**: ✅ Production Implementation

---

## 🔐 Complete Authentication Flow

### End-to-End Request Flow

```
┌─────────────┐
│   Client    │
│  (Browser)  │
└──────┬──────┘
       │
       │ 1. HTTP Request + JWT Token
       │ GET /v2/user/profile/me
       │ Authorization: Bearer eyJhbGc...
       │
       ▼
┌─────────────────────────────────────┐
│      Gateway Service (Port 8080)     │
│                                      │
│  ┌────────────────────────────────┐ │
│  │  JWTAuthMiddleware             │ │
│  │  - Extract token from header   │ │
│  │  - Parse & validate JWT        │ │
│  │  - Extract claims (profileId)  │ │
│  │  - Check public route          │ │
│  └────────────┬───────────────────┘ │
│               │                      │
│               │ 2. Valid token       │
│               ▼                      │
│  ┌────────────────────────────────┐ │
│  │  InjectMetadata                │ │
│  │  - Create gRPC metadata        │ │
│  │  - Add profileId               │ │
│  │  - Add organizationId          │ │
│  │  - Add role, session           │ │
│  └────────────┬───────────────────┘ │
│               │                      │
└───────────────┼──────────────────────┘
                │
                │ 3. gRPC Call + Metadata
                │ metadata: {
                │   profileId: "123"
                │   organizationId: "456"
                │   role: "admin"
                │ }
                │
                ▼
┌─────────────────────────────────────┐
│    User Service (Port 50051)        │
│                                      │
│  ┌────────────────────────────────┐ │
│  │  gRPC Interceptor              │ │
│  │  - Extract metadata            │ │
│  │  - Inject into context         │ │
│  │  - ctx.WithValue(ProfileIDKey) │ │
│  └────────────┬───────────────────┘ │
│               │                      │
│               │ 4. Context with user info
│               ▼                      │
│  ┌────────────────────────────────┐ │
│  │  Handler                       │ │
│  │  - Receive context             │ │
│  │  - Call usecase                │ │
│  └────────────┬───────────────────┘ │
│               │                      │
│               │ 5. Pass context      │
│               ▼                      │
│  ┌────────────────────────────────┐ │
│  │  Usecase                       │ │
│  │  - profileId := _utils.        │ │
│  │      GetProfileIdWithContext() │ │
│  │  - Business logic              │ │
│  └────────────┬───────────────────┘ │
│               │                      │
│               │ 6. Query with context│
│               ▼                      │
│  ┌────────────────────────────────┐ │
│  │  Repository                    │ │
│  │  - db.WithContext(ctx)         │ │
│  │  - Execute query               │ │
│  └────────────┬───────────────────┘ │
│               │                      │
│               │ 7. GORM Callback     │
│               ▼                      │
│  ┌────────────────────────────────┐ │
│  │  AuditBase.BeforeCreate        │ │
│  │  - Extract profileId from ctx  │ │
│  │  - Set CreatedBy automatically │ │
│  └────────────────────────────────┘ │
│                                      │
└──────────────────────────────────────┘
```

---

## 🎫 JWT Token Detailed Structure

### Access Token

**Algorithm**: HS256 (HMAC with SHA-256)
**Expiration**: 15 minutes
**Purpose**: Access protected resources

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "authId": 123,              // Auth method ID
    "profileId": 456,           // User profile ID
    "organizationId": 789,      // Current organization
    "role": "admin",            // User role
    "type": "access",           // Token type
    "session": 999,             // Session ID
    "planId": 1,                // Subscription plan
    "planFrom": "2025-01-01",   // Plan start date
    "exp": 1729123456,          // Expiration (15 min from issue)
    "iat": 1729122556           // Issued at
  },
  "signature": "..."
}
```

### Refresh Token

**Algorithm**: HS256
**Expiration**: 90 days (129,600 minutes)
**Purpose**: Obtain new access token

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "authId": 123,
    "sessionId": 999,
    "type": "refresh",
    "exp": 1737000000,          // Expiration (90 days from issue)
    "iat": 1729123456
  },
  "signature": "..."
}
```

---

## 🔑 JWT Middleware Implementation

### File: `shared/common/jwt/middleware.go`

### 1. JWT Validation Process

```go
func JWTAuthMiddleware(pubRoutes []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Step 1: Check if OPTIONS request (CORS preflight)
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        // Step 2: Get token from Authorization header
        tokenString := c.GetHeader("Authorization")
        
        // Step 3: Check if route is public
        ignoreToken := false
        if isPublicRoute(pubRoutes, c.Request.URL.Path) {
            ignoreToken = true
        }
        
        // Step 4: Validate token presence
        if tokenString == "" && !ignoreToken {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code":    401,
                "message": "Authorization header missing",
            })
            c.Abort()
            return
        }
        
        // Step 5: Remove "Bearer " prefix
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        
        // Step 6: Parse and validate JWT
        claims, err := ParseJWT(tokenString)
        if err != nil && !ignoreToken {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code":    401,
                "message": "Invalid or expired token",
            })
            c.Abort()
            return
        }
        
        // Step 7: Inject claims into gin.Context
        if claims != nil {
            c.Set("claims", claims)
            c.Set("organizationId", claims.OrganizationId)
            c.Set("profileId", claims.ProfileId)
            c.Set("planId", claims.PlanId)
        }
        
        // Step 8: Continue to next handler
        c.Next()
    }
}
```

### 2. Public Route Matching

**Pattern Matching Support**:
- `*` - Wildcard (matches any characters)
- `:` - Parameter placeholder
- Prefix matching

```go
func isPublicRoute(pubRoutes []string, path string) bool {
    for _, route := range pubRoutes {
        if matchRoute(route, path) {
            return true
        }
    }
    return false
}

// Examples:
// Pattern: "/v2/auth/otp" → Matches: "/v2/auth/otp", "/v2/auth/otp/request"
// Pattern: "/v2/user/profile/info/*" → Matches: "/v2/user/profile/info/123"
// Pattern: "/v2/auth/admin/login" → Matches: "/v2/auth/admin/login"
```

**Public Routes List** (126 routes total):
```go
var publicRoutes = []string{
    // Auth routes
    "/v2/auth/otp",              // OTP authentication
    "/v2/auth/token/refresh",    // Token refresh
    "/v2/auth/oauth",            // OAuth login
    "/v2/auth/admin/login",      // Admin login
    "/v2/auth/qr",               // QR authentication
    
    // Public content
    "/v2/user/profile/info/",    // Public profiles
    "/v2/user/profile/search/public",
    "/v2/bdspro/v2/post/global", // Public property listings
    "/v2/bdspro/v2/product/market",
    "/v2/social/news-feed/global",
    
    // File uploads
    "/v1/file/upload",
    "/v1/file/load",
    
    // Payment webhooks
    "/v2/payment/sepay/webhook",
    
    // Feedback
    "/v2/feedback/rate/list",
    "/v2/feedback/report/reasons",
    
    // ... 100+ more public routes
}
```

---

## 🔄 Metadata Injection for gRPC

### File: `shared/common/jwt/middleware.go`

```go
func ExtractMetadataFromRequest(ctx context.Context, r *http.Request) metadata.MD {
    md := metadata.MD{}
    
    // Step 1: Get token from header
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        return md
    }
    
    // Step 2: Parse token
    tokenStr := strings.TrimPrefix(tokenString, "Bearer ")
    claims, err := ParseJWT(tokenStr)
    if err != nil || claims == nil {
        return md
    }
    
    // Step 3: Inject into gRPC metadata
    md.Append("profileId", fmt.Sprintf("%d", claims.ProfileId))
    
    if claims.OrganizationId != nil {
        md.Append("organizationId", fmt.Sprintf("%d", *claims.OrganizationId))
    }
    
    if claims.Session != 0 {
        md.Append("session", fmt.Sprintf("%d", claims.Session))
    }
    
    if claims.Role != "" {
        md.Append("role", claims.Role)
    }
    
    if claims.Type != "" {
        md.Append("type", claims.Type)
    }
    
    if claims.AuthID != 0 {
        md.Append("authId", fmt.Sprintf("%d", claims.AuthID))
    }
    
    if claims.PlanId != nil {
        md.Append("planId", fmt.Sprintf("%d", *claims.PlanId))
    }
    
    md.Append("method", r.Method)
    
    return md
}
```

**Result**: gRPC metadata contains all user context
```
metadata: {
  "profileId": ["123"],
  "organizationId": ["456"],
  "role": ["admin"],
  "session": ["999"],
  "authId": ["111"],
  "planId": ["1"],
  "method": ["GET"]
}
```

---

## 📡 gRPC Context Extraction

### File: `shared/common/middleware/grpc_metadata_interceptor.go`

```go
func InjectGrpcMetadataContextMiddleware(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler) context.Context {
    // Extract gRPC metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return ctx
    }
    
    // Inject into context
    if profileIds := md.Get("profileId"); len(profileIds) > 0 {
        profileId, _ := strconv.ParseUint(profileIds[0], 10, 64)
        ctx = context.WithValue(ctx, _enum.ProfileIDKey, profileId)
    }
    
    if orgIds := md.Get("organizationId"); len(orgIds) > 0 {
        orgId, _ := strconv.ParseUint(orgIds[0], 10, 64)
        ctx = context.WithValue(ctx, _enum.OrganizationIDKey, orgId)
    }
    
    if roles := md.Get("role"); len(roles) > 0 {
        ctx = context.WithValue(ctx, _enum.RoleKey, roles[0])
    }
    
    // ... other fields
    
    return ctx
}
```

---

## 🎯 Context Usage in Services

### Usecase Layer

```go
// File: user-service/internal/usecase/profile_usecase.go
package usecase

import _utils "common/utils"

func (u *ProfileUsecase) GetMyProfile(ctx context.Context) (*Profile, *_err.ErrorDTO) {
    // Extract profileId from context
    profileId := _utils.GetProfileIdWithContext(ctx)
    if profileId == 0 {
        return nil, &_err.ErrorDTO{
            Code:    401,
            Message: "Unauthorized",
        }
    }
    
    // Get profile
    profile, err := u.repo.GetByID(ctx, profileId)
    if err != nil {
        return nil, &_err.ErrorDTO{
            Code:    404,
            Message: "Profile not found",
        }
    }
    
    return profile, nil
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *Product) (*Product, *_err.ErrorDTO) {
    // Get current user and organization
    profileId := _utils.GetProfileIdWithContext(ctx)
    organizationId := _utils.GetOrganizationIdFromContext(ctx)
    
    // Set ownership
    product.OwnerId = profileId
    product.OrganizationId = organizationId
    
    // Create product
    if err := u.repo.Create(ctx, product); err != nil {
        return nil, &_err.ErrorDTO{Code: 500, Message: "Cannot create product"}
    }
    
    return product, nil
}
```

### Repository Layer

```go
// File: bdspro-service/infra/postgre/product_postgres.go
func (r *ProductRepo) Create(ctx context.Context, product *Product) error {
    // Context is passed to GORM
    // AuditBase callbacks will use it
    return r.db.WithTx(ctx, func(tx *gorm.DB) error {
        return tx.WithContext(ctx).Create(product).Error
    })
}
```

### Automatic Audit Trail

```go
// File: shared/common/domain/entity/audit.go
func (base *AuditBase) BeforeCreate(tx *gorm.DB) error {
    // GORM callback automatically called before INSERT
    
    // Extract profileId from context
    userID := GetCurrentUserID(tx)
    
    // Set CreatedBy and UpdatedBy
    base.CreatedBy = &userID
    base.UpdatedBy = &userID
    
    return nil
}

func (base *AuditBase) BeforeUpdate(tx *gorm.DB) error {
    // GORM callback automatically called before UPDATE
    
    userID := GetCurrentUserID(tx)
    base.UpdatedBy = &userID
    
    return nil
}

func GetCurrentUserID(tx *gorm.DB) uint64 {
    // Try gin context first
    ginProfileId := tx.Statement.Context.Value("profileId")
    if ginProfileId != nil {
        return ginProfileId.(uint64)
    }
    
    // Try standard context key
    profileId := tx.Statement.Context.Value(_enums.ProfileIDKey)
    if profileId == nil {
        return 0
    }
    
    return profileId.(uint64)
}
```

**Result**: Every database record automatically tracks who created/updated it!

---

## 🔐 Authorization (RBAC)

### Permission Checking

```go
// Define permissions
const (
    PermissionContactView   = "contact.view"
    PermissionContactCreate = "contact.create"
    PermissionContactUpdate = "contact.update"
    PermissionContactDelete = "contact.delete"
)

// Check in usecase
func (u *ContactUsecase) Create(ctx context.Context, contact *Contact) (*Contact, *_err.ErrorDTO) {
    // Get user role from context
    role := ctx.Value(_enum.RoleKey).(string)
    
    // Check permission
    hasPermission, err := u.permissionRepo.CheckPermission(ctx, role, PermissionContactCreate)
    if err != nil || !hasPermission {
        return nil, &_err.ErrorDTO{
            Code:    403,
            Message: "Bạn không có quyền tạo liên hệ",
        }
    }
    
    // Continue with creation
    // ...
}
```

### Permission Repository

```go
func (r *PermissionRepo) CheckPermission(ctx context.Context, role string, permission string) (bool, error) {
    // Get role ID from role key
    var roleEntity Role
    if err := r.db.Where("key = ?", role).First(&roleEntity).Error; err != nil {
        return false, err
    }
    
    // Parse permission (e.g., "contact.create" → module="contact", action="create")
    parts := strings.Split(permission, ".")
    if len(parts) != 2 {
        return false, errors.New("invalid permission format")
    }
    module := parts[0]
    action := parts[1]
    
    // Check if permission exists for this role
    var perm Permission
    err := r.db.Where("role_id = ? AND module = ? AND action = ? AND is_granted = ?", 
        roleEntity.ID, module, action, true).First(&perm).Error
    
    if err == gorm.ErrRecordNotFound {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    
    return true, nil
}
```

---

## 🛡️ Security Features

### 1. Token Validation

```go
func ParseJWT(tokenString string) (*JWTClaims, error) {
    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        
        // Return secret key
        return []byte(config.JWTSecret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Extract claims
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        // Check expiration
        if claims.ExpiresAt.Before(time.Now()) {
            return nil, errors.New("token expired")
        }
        
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### 2. Admin Access Control

**File**: `auth-service/internal/usecase/admin_access_usecase.go`

```go
func (u *AdminAccessUsecase) ValidateAccess(
    ctx context.Context,
    userId uint64,
    ipAddress string,
    deviceId string,
) (bool, error) {
    // Get admin access rules for user
    accessRules, err := u.repo.GetByUserId(ctx, userId)
    if err != nil {
        return false, err
    }
    
    // Check IP whitelist
    if !isIPAllowed(accessRules, ipAddress) {
        u.logAccess(ctx, userId, ipAddress, "IP_BLOCKED")
        return false, errors.New("IP address not allowed")
    }
    
    // Check device restrictions
    if !isDeviceAllowed(accessRules, deviceId) {
        u.logAccess(ctx, userId, deviceId, "DEVICE_BLOCKED")
        return false, errors.New("Device not allowed")
    }
    
    // Check time restrictions
    if !isTimeAllowed(accessRules) {
        u.logAccess(ctx, userId, ipAddress, "TIME_RESTRICTED")
        return false, errors.New("Access outside allowed hours")
    }
    
    // Log successful access
    u.logAccess(ctx, userId, ipAddress, "SUCCESS")
    
    return true, nil
}
```

### 3. Session Management

```go
func (u *SessionUsecase) CreateSession(
    ctx context.Context,
    authId uint64,
    platform string,
    version string,
    os string,
) (string, error) {
    // Generate session key
    sessionKey := generateSessionKey()
    
    // Create session record
    session := &Session{
        AuthId:     authId,
        Platform:   platform,
        Version:    version,
        OS:         os,
        SessionKey: sessionKey,
        Activate:   true,
        IPRequest:  getIPFromContext(ctx),
        UserAgent:  getUserAgentFromContext(ctx),
        CreatedAt:  time.Now(),
    }
    
    if err := u.repo.Create(ctx, session); err != nil {
        return "", err
    }
    
    // Store in Redis for fast lookup
    u.cacheProvider.Set(ctx, "session:"+sessionKey, session.ID, 90*24*time.Hour)
    
    return sessionKey, nil
}

func (u *SessionUsecase) ValidateSession(ctx context.Context, sessionKey string) (bool, error) {
    // Check Redis first
    sessionId, err := u.cacheProvider.Get(ctx, "session:"+sessionKey)
    if err == nil && sessionId != "" {
        return true, nil
    }
    
    // Fallback to database
    session, err := u.repo.GetBySessionKey(ctx, sessionKey)
    if err != nil {
        return false, err
    }
    
    if !session.Activate {
        return false, errors.New("session deactivated")
    }
    
    return true, nil
}
```

---

## 🔒 Account Lock Mechanism

### File: `auth-service/internal/domain/auth_method.go`

```go
const (
    StatusInactive          = 0  // Account disabled
    StatusActive            = 1  // Normal active
    StatusTemporarilyLocked = 2  // Locked with expiry
    StatusPermanentlyLocked = 3  // Permanently banned
)

type AuthMethod struct {
    ID          uint64
    Status      uint8
    LockedAt    *time.Time
    LockedUntil *time.Time
    LockReason  string
    LockedBy    *uint64
    // ... other fields
}

// Business logic methods
func (a *AuthMethod) IsLocked() bool {
    return a.Status == StatusTemporarilyLocked || a.Status == StatusPermanentlyLocked
}

func (a *AuthMethod) IsTemporarilyLocked() bool {
    return a.Status == StatusTemporarilyLocked
}

func (a *AuthMethod) IsPermanentlyLocked() bool {
    return a.Status == StatusPermanentlyLocked
}

func (a *AuthMethod) IsLockExpired() bool {
    if !a.IsTemporarilyLocked() || a.LockedUntil == nil {
        return false
    }
    return time.Now().After(*a.LockedUntil)
}

func (a *AuthMethod) CanLogin() bool {
    // Check inactive
    if a.Status == StatusInactive {
        return false
    }
    
    // Check permanent lock
    if a.IsPermanentlyLocked() {
        return false
    }
    
    // Check temporary lock (with expiry check)
    if a.IsTemporarilyLocked() && !a.IsLockExpired() {
        return false
    }
    
    return true
}
```

**Usage in Login Flow**:
```go
func (u *AuthUsecase) VerifyOTP(ctx context.Context, phone string, otp string) (*AuthResponse, *_err.ErrorDTO) {
    // Find auth method
    authMethod, _ := u.authMethodRepo.FindByPhone(ctx, phone)
    
    // Check if account can login
    if !authMethod.CanLogin() {
        return nil, &_err.ErrorDTO{
            Code:    403,
            Message: fmt.Sprintf("Tài khoản bị khóa: %s", authMethod.LockReason),
        }
    }
    
    // Verify OTP
    // ...
    
    // Generate tokens
    // ...
}
```

---

## 🔄 Complete Authentication Lifecycle

### 1. Login (OTP)
```
User → Request OTP
    ↓
Generate OTP → Save to DB → Send via SMS
    ↓
User → Submit OTP
    ↓
Validate OTP → Check account status → Create session
    ↓
Generate JWT (access + refresh) → Return to user
    ↓
User stores tokens → Use in subsequent requests
```

### 2. API Request with JWT
```
User → API Request + JWT
    ↓
Gateway → JWT Middleware → Validate token
    ↓
Extract claims → Inject into gRPC metadata
    ↓
Service → gRPC Interceptor → Extract metadata
    ↓
Inject into context → Pass to handler
    ↓
Handler → Usecase (with context)
    ↓
Usecase → Extract user info → Business logic
    ↓
Repository → Use context → Execute query
    ↓
GORM → Audit callback → Auto-set CreatedBy
    ↓
Response → Return to user
```

### 3. Token Refresh
```
User → Refresh token expired warning
    ↓
POST /v2/auth/token/refresh + refreshToken
    ↓
Validate refresh token → Extract authId, sessionId
    ↓
Check session still active
    ↓
Generate new access token (15 min)
    ↓
Generate new refresh token (90 days)
    ↓
Return new tokens → User updates stored tokens
```

### 4. Logout
```
User → DELETE /v2/auth/logout
    ↓
Extract session from JWT
    ↓
Deactivate session in DB
    ↓
Delete from Redis cache
    ↓
Return success → User deletes local tokens
```

---

## 🔑 Context Keys Reference

### File: `shared/common/domain/enum/context_key_enum.go`

```go
type ContextKey string

const (
    ProfileIDKey      ContextKey = "profileId"       // User profile ID
    OrganizationIDKey ContextKey = "organizationId"  // Current organization
    AuthIDKey         ContextKey = "authId"          // Auth method ID
    SessionKey        ContextKey = "sessionKey"      // Session key
    RoleKey           ContextKey = "role"            // User role
    TypeKey           ContextKey = "type"            // Auth type
    PlanIDKey         ContextKey = "planId"          // Subscription plan
    PlanFromKey       ContextKey = "planFrom"        // Plan start date
)
```

**Utility Functions**:
```go
// File: shared/common/utils/context.go

func GetProfileIdWithContext(ctx context.Context) uint64 {
    profileId := ctx.Value(_enums.ProfileIDKey)
    if profileId == nil {
        return 0
    }
    return profileId.(uint64)
}

func GetOrganizationIdFromContext(ctx context.Context) uint64 {
    organizationId := ctx.Value(_enums.OrganizationIDKey)
    if organizationId == nil {
        return 0
    }
    return organizationId.(uint64)
}

func CloneContext(ctx context.Context) context.Context {
    // Create new context with user info copied
    cloneCtx := context.Background()
    
    profileId := GetProfileIdWithContext(ctx)
    organizationId := GetOrganizationIdFromContext(ctx)
    
    cloneCtx = context.WithValue(cloneCtx, _enums.ProfileIDKey, profileId)
    cloneCtx = context.WithValue(cloneCtx, _enums.OrganizationIDKey, organizationId)
    
    return cloneCtx
}
```

---

## 📊 Security Metrics

### Token Management
- **Access Token TTL**: 15 minutes
- **Refresh Token TTL**: 90 days
- **Session Storage**: Redis + PostgreSQL
- **Token Algorithm**: HMAC-SHA256

### Access Control
- **Public Routes**: 126 routes (no authentication)
- **Protected Routes**: All others (JWT required)
- **Admin Routes**: Additional IP/device restrictions
- **Permission Granularity**: module.action format

### Audit Trail
- **Auto-tracking**: CreatedBy, UpdatedBy on all entities
- **GORM Callbacks**: Automatic population
- **History Service**: Detailed action logging
- **Admin Activity**: Special tracking for admin actions

---

**Last Updated**: October 15, 2025
**Version**: v1.0
**Status**: ✅ Production Implementation

