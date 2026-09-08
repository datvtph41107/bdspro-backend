# II.4 JWT Authentication - User Information Extraction

## 🎯 Tổng Quan

JWT (JSON Web Token) là phương thức xác thực chính trong dự án BDSPro Microservices. Hệ thống sử dụng JWT để lưu trữ thông tin người dùng và truyền tải qua các request, cho phép các service xác định danh tính và quyền hạn của người dùng.

---

## 📋 Thông Tin Trong JWT Token

### **JWT Claims chứa các thông tin sau:**
- **`sub` (AuthID)**: ID xác thực của user
- **`profile` (ProfileId)**: ID profile của user  
- **`organization` (OrganizationId)**: ID tổ chức hiện tại
- **`plan` (PlanId)**: ID gói subscription
- **`type`**: Loại token (ACCESS/REFRESH)
- **`role`**: Vai trò của user
- **`session`**: ID phiên làm việc
- **`start` (PlanFrom)**: Ngày bắt đầu gói
- **`iat`**: Thời gian tạo token
- **`exp`**: Thời gian hết hạn token

---

## 🔧 Các Function Chính

### **1. Lấy Thông Tin User từ Context**
```go
import _utils "common/utils"

// Lấy Profile ID của user hiện tại
profileID := _utils.GetProfileIdWithContext(ctx)

// Lấy Organization ID của tổ chức hiện tại  
organizationID := _utils.GetOrganizationIdFromContext(ctx)

// Lấy Session ID
sessionID := _utils.GetSessionIdFromContext(ctx)

// Lấy Plan ID
planID := _utils.GetPlanIdFromContext(ctx)
```

### **2. Context Management**
```go
// Tạo context mới với thông tin user hiện tại
newCtx := _utils.CloneContext(ctx)

// Tạo context từ Gin context (HTTP)
ctx := _utils.CreateContextFromGin(c)
```

### **3. JWT Token Operations**
```go
import _jwt "common/jwt"

// Tạo JWT token
token, err := _jwt.GenerateToken(jwtTokenProperties)

// Parse JWT token
claims, err := _jwt.ParseJWT(tokenString)

// Lấy Principal từ context
principal := _jwt.GetPrincipalContext(ctx)
```

---

## 🔄 JWT Flow trong Hệ Thống

### **Authentication Flow:**
1. User Login → Auth Service
2. Auth Service validates credentials  
3. Generate JWT Token with user info
4. Return token to client
5. Client stores token and sends in subsequent requests

### **Request Processing Flow:**
1. Client sends request with JWT token in Authorization header
2. Gateway/Service extracts token
3. Parse and validate JWT token
4. Extract user information from claims
5. Store user info in context
6. Pass context to business logic
7. Business logic uses user info via utils functions

---

## 🛠️ JWT Middleware & Context Processing

### **gRPC Metadata Interceptor:**
- Parse JWT token từ gRPC metadata
- Extract user info và store vào context
- Sử dụng constants keys: ProfileIDKey, OrganizationIDKey, AuthIDKey, etc.

### **HTTP JWT Middleware:**
- Extract token từ Authorization header
- Parse JWT payload
- Store claims vào context với CONTEXT_USER_KEY

---

## 📋 Context Keys

### **Constants Keys:**
- **ProfileIDKey**: "profileId" - User profile ID
- **OrganizationIDKey**: "organizationId" - Organization ID  
- **AuthIDKey**: "authId" - Authentication ID
- **SessionKey**: "session" - Session ID
- **RoleKey**: "role" - User role
- **TypeKey**: "type" - Token type
- **PlanIDKey**: "planId" - Plan ID
- **PlanFromKey**: "planFrom" - Plan start date

---

## 🚀 Cách Sử Dụng trong Business Logic

### **Usecase Layer:**
```go
import _utils "common/utils"

func (u *userUsecase) GetUserProfile(ctx context.Context, req *dto.GetUserProfileRequest) (*dto.GetUserProfileResponse, error) {
    // Lấy Profile ID của user hiện tại
    profileID := _utils.GetProfileIdWithContext(ctx)
    if profileID == 0 {
        return nil, errors.New("user not authenticated")
    }

    // Lấy Organization ID của tổ chức hiện tại
    organizationID := _utils.GetOrganizationIdFromContext(ctx)
    
    // Sử dụng thông tin user trong business logic
    user, err := u.userRepo.GetByID(ctx, profileID)
    if err != nil {
        return nil, err
    }

    return &dto.GetUserProfileResponse{
        User: user,
        OrganizationID: organizationID,
    }, nil
}
```

### **Handler Layer:**
- Context đã chứa thông tin user từ middleware
- Chỉ cần gọi usecase với context
- Không cần parse JWT manually

### **Repository Layer:**
- Sử dụng user info để audit trail
- Filter data theo organization context
- Log user actions

---

## ⚠️ Lưu Ý Quan Trọng

### **Quy Tắc Sử Dụng Utils:**
- **Luôn import utils với alias `_utils`**
- **Chỉ sử dụng trong layer usecase, handler, hoặc service**
- **Không tự parse JWT hay context ngoài utils**

### **Error Handling:**
- Kiểm tra user đã authenticated chưa
- Kiểm tra organization context
- Validate permissions

### **Security Best Practices:**
- Luôn validate user permissions
- Kiểm tra organization context
- Log audit trail với user info
- Không expose sensitive information

---

## 🎯 Tóm Tắt

| Function | Mục đích | Khi nào sử dụng |
|----------|----------|-----------------|
| `GetProfileIdWithContext(ctx)` | Lấy Profile ID của user | Trong usecase, handler, service |
| `GetOrganizationIdFromContext(ctx)` | Lấy Organization ID | Khi cần context tổ chức |
| `GetSessionIdFromContext(ctx)` | Lấy Session ID | Để tracking session |
| `GetPlanIdFromContext(ctx)` | Lấy Plan ID | Kiểm tra subscription |
| `CloneContext(ctx)` | Tạo context mới | Khi cần context riêng biệt |

**Quy tắc vàng**: **Luôn sử dụng `_utils.GetProfileIdWithContext(ctx)` và `_utils.GetOrganizationIdFromContext(ctx)` để lấy thông tin user thay vì tự parse JWT hay context!**
