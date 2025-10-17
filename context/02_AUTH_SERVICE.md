# Auth Service - Complete Context

## 📋 Service Overview

**Service Name**: Auth Service
**Port**: 50062 (gRPC)
**Primary Role**: Authentication, Authorization & Access Control
**Status**: ✅ Production Ready

---

## 🎯 Core Responsibilities

### 1. User Authentication
- **OTP Authentication**: SMS/Phone-based OTP login
- **OAuth Integration**: Google, Facebook, Zalo
- **Admin Authentication**: Username/Password for admin users
- **QR Code Authentication**: QR-based login for mobile apps

### 2. Authorization & Access Control
- **Role Management**: Create, update, delete roles
- **Permission Management**: Fine-grained permission control
- **Role Groups**: Group multiple roles together
- **RBAC**: Role-Based Access Control enforcement

### 3. Token Management
- **JWT Generation**: Access & Refresh tokens
- **Token Refresh**: Automatic token renewal
- **Session Management**: User session tracking
- **Token Revocation**: Logout & token invalidation

### 4. Account Management
- **Account Creation**: Multi-provider account creation
- **Account Deletion**: Soft delete with restore capability
- **Account Restoration**: Restore deleted accounts
- **Account Locking**: Temporary & permanent locks

### 5. Admin Access Control
- **IP Whitelisting**: IP-based access control
- **Device Management**: Device registration & tracking
- **2FA Enforcement**: Two-factor authentication
- **Access Logging**: Audit trail for admin access

---

## 📁 Service Architecture

### Domain Layer (`internal/domain/`)
```
auth-service/internal/domain/
├── auth_method.go           # Main authentication entity
├── session.go               # User session entity
├── otp.go                   # OTP verification entity
├── role.go                  # Role entity
├── permission.go            # Permission entity
├── role_group.go            # Role group entity
├── admin_access.go          # Admin access control entity
├── user_info.go             # Extended user information
├── module.go                # Module/feature entity
├── color.go                 # UI color preferences
└── status.go                # Account status entity
```

### Interface Layer (`internal/interface/`)

#### Repository Interfaces (`repo/`)
- `AuthMethodRepository`: Auth method CRUD operations
- `SessionRepository`: Session management
- `OTPRepository`: OTP storage & validation
- `RoleRepository`: Role management
- `PermissionRepository`: Permission management
- `RoleGroupRepository`: Role group management
- `AdminAccessRepository`: Admin access control
- `UserInfoRepository`: User info management

#### Provider Interfaces (`provider/`)
- `CacheProvider`: Redis cache operations
- `CookieProvider`: Cookie management
- `NotificationProvider`: Notification service client
- `OrganizationProvider`: Organization service client
- `PaymentProvider`: Payment service client
- `ProfileProvider`: User profile service client
- `AuthValidator`: Authentication validation
- `ZNSProvider`: ZNS notification (Zalo Notification Service)

### Usecase Layer (`internal/usecase/`)
```
auth-service/internal/usecase/
├── auth_user_usecase.go     # User authentication logic
├── auth_admin_usecase.go    # Admin authentication logic
├── otp_usecase.go           # OTP generation & validation
├── role_usecase.go          # Role management
├── permission_usecase.go    # Permission management
├── role_group_usecase.go    # Role group management
├── admin_access_usecase.go  # Admin access control
├── user_info_usecase.go     # User info management
└── module_usecase.go        # Module management
```

### Infrastructure Layer (`infra/`)

#### Handlers (`handler/`)
- `AuthHandler`: Main authentication endpoints
- `OAuthHandler`: OAuth provider endpoints
- `RoleHandler`: Role management endpoints
- `PermissionHandler`: Permission management endpoints
- `RoleGroupHandler`: Role group endpoints
- `InternalHandler`: Internal service-to-service calls
- `UserInfoHandler`: User info endpoints

#### PostgreSQL Repositories (`postgre/`)
- `AuthMethodPostgres`: Auth method database operations
- `SessionPostgres`: Session storage
- `OTPPostgres`: OTP storage
- `RolePostgres`: Role database operations
- `PermissionPostgres`: Permission storage
- `RoleGroupPostgres`: Role group storage
- `AdminAccessPostgres`: Admin access control storage

#### OAuth Providers (`oauth/`)
- `GoogleOAuth`: Google authentication
- `FacebookOAuth`: Facebook authentication
- `ZaloOAuth`: Zalo authentication

#### External Clients (`client/`)
- `NotificationClient`: Notification service gRPC client
- `OrganizationClient`: Organization service gRPC client
- `PaymentClient`: Payment service gRPC client
- `UserClient`: User service gRPC client

#### Other Infrastructure
- `RedisProvider`: Redis cache implementation
- `CookieProvider`: Cookie management implementation
- `ZNSProvider`: Zalo Notification Service integration

---

## 🗄️ Database Schema

### auth_method Table
```sql
CREATE TABLE auth_method (
    id BIGSERIAL PRIMARY KEY,
    provider provider_enum NOT NULL,          -- PHONE, GOOGLE, FACEBOOK, ZALO, ADMIN
    auth_name VARCHAR(255),                   -- username or identifier
    password VARCHAR(255),                    -- hashed password
    is_sensitive BOOLEAN DEFAULT FALSE,
    avatar VARCHAR(255),
    full_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(255),
    role_key INT,
    status SMALLINT DEFAULT 1,                -- 0:inactive, 1:active, 2:temp_locked, 3:perm_locked
    locked_at TIMESTAMP,
    locked_until TIMESTAMP,
    lock_reason VARCHAR(500),
    locked_by BIGINT,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_auth_method_phone ON auth_method(phone);
CREATE INDEX idx_auth_method_email ON auth_method(email);
CREATE INDEX idx_auth_method_provider ON auth_method(provider);
CREATE INDEX idx_auth_method_user_id ON auth_method(user_id);
```

### session Table
```sql
CREATE TABLE session (
    id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,
    platform VARCHAR(50),                     -- web, ios, android
    version VARCHAR(20),
    os VARCHAR(50),
    session_key VARCHAR(255) UNIQUE,
    client_id VARCHAR(255),
    activate BOOLEAN DEFAULT TRUE,
    ip_request VARCHAR(50),
    total_request BIGINT DEFAULT 0,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP,
    FOREIGN KEY (auth_id) REFERENCES auth_method(id)
);

CREATE INDEX idx_session_auth_id ON session(auth_id);
CREATE INDEX idx_session_key ON session(session_key);
```

### otp Table
```sql
CREATE TABLE otp (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    otp_code VARCHAR(10) NOT NULL,
    auth_id BIGINT,
    expired_at TIMESTAMP NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (auth_id) REFERENCES auth_method(id)
);

CREATE INDEX idx_otp_phone ON otp(phone);
CREATE INDEX idx_otp_auth_id ON otp(auth_id);
```

### role Table
```sql
CREATE TABLE role (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### permission Table
```sql
CREATE TABLE permission (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    module VARCHAR(100) NOT NULL,              -- contact, post, product, etc.
    action VARCHAR(50) NOT NULL,               -- view, create, update, delete
    resource VARCHAR(100),
    is_granted BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (role_id) REFERENCES role(id)
);

CREATE INDEX idx_permission_role_id ON permission(role_id);
CREATE INDEX idx_permission_module ON permission(module);
```

### admin_access Table
```sql
CREATE TABLE admin_access (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    ip_address VARCHAR(50),
    ip_range VARCHAR(100),
    device_id VARCHAR(255),
    device_name VARCHAR(255),
    device_type VARCHAR(50),                   -- desktop, mobile, tablet
    auth_type VARCHAR(50),                     -- password, 2fa, biometric
    status VARCHAR(20),                        -- active, inactive, blocked
    effective_from TIMESTAMP,
    effective_to TIMESTAMP,
    max_devices INT DEFAULT 5,
    require_2fa BOOLEAN DEFAULT FALSE,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_admin_access_user_id ON admin_access(user_id);
CREATE INDEX idx_admin_access_ip ON admin_access(ip_address);
```

---

## 🔐 Authentication Flows

### 1. OTP Authentication Flow
```
1. User requests OTP
   POST /v2/auth/otp/request
   Body: { phone, fullname }
   
2. Generate & send OTP
   - Generate 6-digit OTP
   - Store in database with expiry (5 minutes)
   - Send via SMS/ZNS
   - Return authId
   
3. User submits OTP
   POST /v2/auth/otp/verify
   Body: { phone, otp, authId }
   
4. Verify OTP
   - Check OTP validity & expiry
   - Validate against stored OTP
   - Mark as verified
   
5. Create/Update AuthMethod
   - Find or create auth_method by phone
   - Link to user profile
   
6. Generate tokens
   - Generate access token (15 min)
   - Generate refresh token (90 days)
   - Create session record
   
7. Return tokens
   Response: { 
     accessToken, 
     refreshToken, 
     profileId, 
     organizationId 
   }
```

### 2. OAuth Flow (Google/Facebook/Zalo)
```
1. Client initiates OAuth
   POST /v2/auth/oauth/{provider}/login
   
2. Redirect to provider
   - Generate state token
   - Redirect to OAuth provider
   
3. Provider callback
   - Receive authorization code
   - Exchange for access token
   
4. Get user profile
   - Fetch user info from provider
   - Extract email, name, avatar
   
5. Find or create AuthMethod
   - Search by oauth_id
   - Create if not exists
   - Link to user profile
   
6. Generate JWT tokens
   - Create session
   - Generate access & refresh tokens
   
7. Return tokens
```

### 3. Admin Login Flow
```
1. Admin submits credentials
   POST /v2/auth/admin/login
   Body: { username, password }
   
2. Validate credentials
   - Find auth_method by username
   - Verify password hash
   - Check account status
   
3. Check admin access control
   - Validate IP address
   - Check device restrictions
   - Verify 2FA if required
   
4. Log access attempt
   - Create admin_access_log
   - Record IP, device, timestamp
   
5. Generate tokens
   - Create admin session
   - Generate JWT with admin role
   
6. Return tokens
```

### 4. QR Authentication Flow
```
1. Web client initializes QR
   POST /v2/auth/qr/init
   - Generate session key
   - Create QR code data
   - Store in Redis (5 min expiry)
   
2. Display QR code
   - Show QR on web browser
   - Poll for status updates
   
3. Mobile app scans QR
   - Parse session key
   - Validate session exists
   
4. User confirms on mobile
   POST /v2/auth/qr/confirm
   Body: { sessionId, sessionKey }
   - Update session status to 'confirmed'
   
5. Web polls status
   POST /v2/auth/qr/status
   - Check session status
   - Return 'pending'/'confirmed'/'expired'
   
6. On confirmed
   - Generate tokens for web session
   - Link to mobile user
   - Return tokens
```

---

## 🔑 JWT Token Structure

### Access Token Payload
```json
{
  "authId": 123,
  "profileId": 456,
  "organizationId": 789,
  "role": "admin",
  "permissions": ["contact.view", "post.create"],
  "exp": 1234567890,
  "iat": 1234567000
}
```

### Refresh Token Payload
```json
{
  "authId": 123,
  "sessionId": 999,
  "type": "refresh",
  "exp": 1234567890,
  "iat": 1234567000
}
```

### Token Expiry
- **Access Token**: 15 minutes
- **Refresh Token**: 90 days (129,600 minutes)

---

## 🛡️ Role-Based Access Control (RBAC)

### Permission Format
```
{module}.{action}

Examples:
- contact.view
- contact.create
- contact.update
- contact.delete
- post.view
- post.create
- product.view
- product.update
```

### Default Roles
```go
const (
    RoleSuperAdmin  = "SUPER_ADMIN"     // Full system access
    RoleAdmin       = "ADMIN"           // Organization admin
    RoleMember      = "MEMBER"          // Regular member
    RoleViewer      = "VIEWER"          // Read-only access
)
```

### Permission Checking
```go
// In usecase/middleware
func CheckPermission(ctx context.Context, module, action string) error {
    profileId := GetProfileIdFromContext(ctx)
    role := GetRoleFromContext(ctx)
    
    // Check if role has permission
    hasPermission := permissionRepo.CheckPermission(role, module, action)
    if !hasPermission {
        return ErrForbidden
    }
    
    return nil
}
```

---

## 🔒 Account Locking Mechanism

### Lock Types
1. **Inactive (status=0)**: Account disabled, can be reactivated
2. **Active (status=1)**: Normal active account
3. **Temporarily Locked (status=2)**: Locked with expiry time
4. **Permanently Locked (status=3)**: Permanently banned

### Lock Logic
```go
func (a *AuthMethod) CanLogin() bool {
    if a.Status == StatusInactive {
        return false
    }
    if a.IsPermanentlyLocked() {
        return false
    }
    if a.IsTemporarilyLocked() && !a.IsLockExpired() {
        return false
    }
    return true
}

func (a *AuthMethod) IsLockExpired() bool {
    if !a.IsTemporarilyLocked() || a.LockedUntil == nil {
        return false
    }
    return time.Now().After(*a.LockedUntil)
}
```

---

## 📡 External Service Integrations

### 1. Notification Service
**Purpose**: Send OTP, notifications
**Methods**:
- `SendOTP(phone, otp)`: Send OTP via SMS
- `SendEmail(email, subject, body)`: Send email notification

### 2. Organization Service
**Purpose**: Get organization info, validate membership
**Methods**:
- `GetOrganizationByUser(userId)`: Get user's organizations
- `ValidateMembership(userId, orgId)`: Check if user is member

### 3. Payment Service
**Purpose**: Check subscription status, features access
**Methods**:
- `GetSubscriptionStatus(orgId)`: Get org subscription
- `ValidateFeatureAccess(orgId, feature)`: Check feature access

### 4. User Service (Profile)
**Purpose**: Get/update user profile
**Methods**:
- `GetProfileById(profileId)`: Get user profile
- `UpdateProfile(profileId, data)`: Update profile

### 5. ZNS (Zalo Notification Service)
**Purpose**: Send OTP via Zalo
**Methods**:
- `SendZaloOTP(phone, otp)`: Send OTP via Zalo

---

## 🔧 Configuration

### Environment Variables
```yaml
# config/develop.yml
server:
  port: 8216
  tcp_port: 50062

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: auth_db

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  access_exp_minutes: 15
  refresh_exp_minutes: 129600

oauth:
  google:
    client_id: "google-client-id"
    client_secret: "google-client-secret"
    redirect_url: "http://localhost:8080/auth/google/callback"
  
  facebook:
    client_id: "facebook-app-id"
    client_secret: "facebook-app-secret"
    redirect_url: "http://localhost:8080/auth/facebook/callback"
  
  zalo:
    app_id: "zalo-app-id"
    app_secret: "zalo-app-secret"
    redirect_url: "http://localhost:8080/auth/zalo/callback"

zns:
  api_key: "zns-api-key"
  template_id: "otp-template-id"
```

---

## 🚀 API Endpoints Summary

### Public Endpoints (No Auth)
- `POST /v2/auth/otp/request` - Request OTP
- `POST /v2/auth/otp/verify` - Verify OTP
- `POST /v2/auth/otp/resend` - Resend OTP
- `POST /v2/auth/token/refresh` - Refresh token
- `PUT /v2/auth/restore` - Restore account
- `POST /v2/auth/qr/init` - Initialize QR auth
- `POST /v2/auth/admin/login` - Admin login
- `POST /v2/auth/oauth/{provider}/login` - OAuth login

### Protected Endpoints (Auth Required)
- `DELETE /v2/auth/logout` - Logout
- `DELETE /v2/auth/delete` - Delete account
- `POST /v2/auth/switch/organization` - Switch org
- `GET /v2/auth/admin-access` - List admin access
- `POST /v2/auth/admin-access` - Create admin access
- `PUT /v2/auth/admin-access/{id}` - Update admin access
- `DELETE /v2/auth/admin-access/{id}` - Delete admin access

---

## 🧪 Testing Scenarios

### 1. OTP Login Test
```bash
# Request OTP
curl -X POST http://localhost:8080/v2/auth/otp/request \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567", "fullname": "Test User"}'

# Verify OTP
curl -X POST http://localhost:8080/v2/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone": "0901234567", "otp": "123456", "authId": 1}'
```

### 2. Token Refresh Test
```bash
curl -X POST http://localhost:8080/v2/auth/token/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken": "your-refresh-token"}'
```

### 3. Admin Login Test
```bash
curl -X POST http://localhost:8080/v2/auth/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

---

## 📊 Key Metrics

### Database Tables: 10
- auth_method
- session
- otp
- role
- permission
- role_group
- admin_access
- admin_access_log
- user_info
- module

### Handlers: 7
### Usecases: 9
### Repositories: 11
### External Clients: 4
### OAuth Providers: 3

---

## 🔄 Service Dependencies

### Depends On:
- **User Service**: Profile management
- **Organization Service**: Organization validation
- **Payment Service**: Subscription checks
- **Notification Service**: OTP & email sending

### Depended By:
- **Gateway Service**: JWT validation
- **All Services**: Authentication & authorization

---

**Last Updated**: October 15, 2025
**Service Version**: v2
**Status**: ✅ Production Ready

