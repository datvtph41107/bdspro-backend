# II.3 Swagger - API Documentation

## 🎯 Tổng Quan

Swagger là công cụ tạo tài liệu API tự động từ các comment trong code Go. Trong dự án BDSPro, Swagger được sử dụng để tạo ra tài liệu API cho tất cả các endpoints, giúp developers hiểu và test APIs một cách dễ dàng.

---

## 🔧 Cách Sử Dụng Swagger

### **1. Generate Swagger Documentation**

```bash
# Generate Swagger cho service cụ thể
make swag <service-name>
# Ví dụ: make swag user

# Hoặc sử dụng script trực tiếp
./script.sh swag <service-name>

# Generate Swagger cho file service (lệnh đặc biệt)
make swag-file

# Merge tất cả Swagger docs
make merge-swagger
```

### **2. Quy Tắc Viết Swagger Comments**

Swagger comments được viết **TRỰC TIẾP TẠI CÁC HANDLER** theo format chuẩn:

```go
// @Summary     Lấy thông tin user theo ID
// @Description Lấy thông tin chi tiết của user dựa trên ID
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id   path     int  true  "User ID"
// @Success     200  {object} userpb.GetUserResponse
// @Failure     400  {object} common.ErrorResponse
// @Failure     404  {object} common.ErrorResponse
// @Failure     500  {object} common.ErrorResponse
// @Router      /v2/user/{id} [get]
// @Security    BearerAuth
func (h *userHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
    // Implementation
}
```

---

## 📝 Các Thẻ Swagger Quan Trọng

### **1. Thông Tin Cơ Bản**
```go
// @Summary     Tóm tắt ngắn gọn về API
// @Description Mô tả chi tiết về chức năng của API
// @Tags        Nhóm API (User, Organization, etc.)
```

### **2. Request/Response**
```go
// @Accept      json                    // Loại dữ liệu nhận vào
// @Produce     json                    // Loại dữ liệu trả về
// @Param       id   path     int  true  "User ID"           // Path parameter
// @Param       body body     userpb.CreateUserRequest true "User data"  // Body parameter
// @Param       page query    int  false "Page number"       // Query parameter
// @Success     200  {object} userpb.GetUserResponse         // Response thành công
// @Failure     400  {object} common.ErrorResponse           // Response lỗi
```

### **3. Routing & Security**
```go
// @Router      /v2/user/{id} [get]     // Route và HTTP method
// @Security    BearerAuth              // Yêu cầu authentication
```

---

## 🔍 Ví Dụ Thực Tế

### **1. GET API - Lấy dữ liệu**
```go
// @Summary     Lấy danh sách users
// @Description Lấy danh sách users với phân trang
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       page query int false "Page number" default(1)
// @Param       size query int false "Page size" default(10)
// @Success     200 {object} userpb.GetUsersResponse
// @Failure     400 {object} common.ErrorResponse
// @Failure     500 {object} common.ErrorResponse
// @Router      /v2/user/users [get]
// @Security    BearerAuth
func (h *userHandler) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
    // Implementation
}
```

### **2. POST API - Tạo mới**
```go
// @Summary     Tạo user mới
// @Description Tạo một user mới trong hệ thống
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       body body userpb.CreateUserRequest true "User data"
// @Success     201 {object} userpb.CreateUserResponse
// @Failure     400 {object} common.ErrorResponse
// @Failure     409 {object} common.ErrorResponse
// @Failure     500 {object} common.ErrorResponse
// @Router      /v2/user/users [post]
// @Security    BearerAuth
func (h *userHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
    // Implementation
}
```

### **3. PUT API - Cập nhật**
```go
// @Summary     Cập nhật thông tin user
// @Description Cập nhật thông tin của user theo ID
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id   path     int  true "User ID"
// @Param       body body     userpb.UpdateUserRequest true "Updated user data"
// @Success     200 {object} userpb.UpdateUserResponse
// @Failure     400 {object} common.ErrorResponse
// @Failure     404 {object} common.ErrorResponse
// @Failure     500 {object} common.ErrorResponse
// @Router      /v2/user/{id} [put]
// @Security    BearerAuth
func (h *userHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
    // Implementation
}
```

### **4. DELETE API - Xóa**
```go
// @Summary     Xóa user
// @Description Xóa user theo ID
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} common.SuccessResponse
// @Failure     400 {object} common.ErrorResponse
// @Failure     404 {object} common.ErrorResponse
// @Failure     500 {object} common.ErrorResponse
// @Router      /v2/user/{id} [delete]
// @Security    BearerAuth
func (h *userHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
    // Implementation
}
```

---

## 🏷️ Các Loại Parameters

### **1. Path Parameters**
```go
// @Param       id   path     int  true  "User ID"
// @Param       orgId path    int  true  "Organization ID"
```

### **2. Query Parameters**
```go
// @Param       page query    int  false "Page number" default(1)
// @Param       size query    int  false "Page size" default(10)
// @Param       search query  string false "Search keyword"
// @Param       status query  string false "Filter by status"
```

### **3. Body Parameters**
```go
// @Param       body body     userpb.CreateUserRequest true "User data"
// @Param       body body     userpb.UpdateUserRequest true "Updated user data"
```

### **4. Header Parameters**
```go
// @Param       Authorization header string true "Bearer token"
// @Param       Content-Type header string true "application/json"
```

---

## 📊 Response Types

### **1. Success Responses**
```go
// @Success     200 {object} userpb.GetUserResponse
// @Success     201 {object} userpb.CreateUserResponse
// @Success     204 "No Content"
```

### **2. Error Responses**
```go
// @Failure     400 {object} common.ErrorResponse "Bad Request"
// @Failure     401 {object} common.ErrorResponse "Unauthorized"
// @Failure     403 {object} common.ErrorResponse "Forbidden"
// @Failure     404 {object} common.ErrorResponse "Not Found"
// @Failure     409 {object} common.ErrorResponse "Conflict"
// @Failure     422 {object} common.ErrorResponse "Validation Error"
// @Failure     500 {object} common.ErrorResponse "Internal Server Error"
```

---

## 🔐 Security & Authentication

### **1. Bearer Token Authentication**
```go
// @Security    BearerAuth
```

### **2. API Key Authentication**
```go
// @Security    ApiKeyAuth
```

### **3. No Authentication**
```go
// Không cần thêm @Security tag
```

---

## 🚀 Workflow Sử Dụng Swagger

### **1. Khi Thêm API Mới**
```go
// 1. Viết handler function
func (h *userHandler) NewAPI(ctx context.Context, req *userpb.NewAPIRequest) (*userpb.NewAPIResponse, error) {
    // Implementation
}

// 2. Thêm Swagger comments
// @Summary     Mô tả API mới
// @Description Chi tiết về API mới
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       body body userpb.NewAPIRequest true "Request data"
// @Success     200 {object} userpb.NewAPIResponse
// @Failure     400 {object} common.ErrorResponse
// @Router      /v2/user/new-api [post]
// @Security    BearerAuth
func (h *userHandler) NewAPI(ctx context.Context, req *userpb.NewAPIRequest) (*userpb.NewAPIResponse, error) {
    // Implementation
}

// 3. Generate Swagger
make swag user
```

### **2. Khi Sửa API Hiện Tại**
```go
// 1. Sửa Swagger comments
// @Summary     Mô tả đã được cập nhật
// @Description Chi tiết mới về API
// ...

// 2. Generate Swagger
make swag user
```

### **3. Khi Merge Tất Cả Swagger**
```bash
# Generate Swagger cho tất cả services
make swag user
make swag organization
make swag social
make swag bdspro
# ... các services khác

# Merge tất cả Swagger docs
make merge-swagger
```

---

## 📁 Cấu Trúc Swagger Files

### **1. Generated Files**
```
<service>-service/
├── docs/
│   ├── docs.go          # Generated Swagger docs
│   ├── swagger.json     # Swagger JSON format
│   └── swagger.yaml     # Swagger YAML format
```

### **2. Gateway Service**
```
gateway-service/
├── docs/
│   ├── user/            # User service Swagger
│   ├── organization/    # Organization service Swagger
│   ├── social/          # Social service Swagger
│   └── merged/          # Merged Swagger docs
```

---

## 🌐 Truy Cập Swagger UI

### **1. Individual Service Swagger**
```
http://localhost:8080/swagger/user/index.html
http://localhost:8080/swagger/organization/index.html
http://localhost:8080/swagger/social/index.html
```

### **2. Merged Swagger (Tất cả services)**
```
http://localhost:8080/swagger/merged/index.html
```

### **3. File Service Swagger**
```
http://localhost:8081/swagger/file/swagger.json
http://localhost:8081/swagger/index.html
```

---

## ⚠️ Lưu Ý Quan Trọng

### **1. Quy Tắc Viết Comments**
- **Luôn viết comments TRỰC TIẾP TẠI HANDLER**
- **Sử dụng format chuẩn Swagger**
- **Mô tả rõ ràng và chi tiết**
- **Định nghĩa đầy đủ parameters và responses**

### **2. Best Practices**
- **Tóm tắt ngắn gọn trong @Summary**
- **Mô tả chi tiết trong @Description**
- **Sử dụng Tags để nhóm APIs**
- **Định nghĩa đầy đủ error responses**
- **Thêm @Security cho APIs cần authentication**

### **3. Common Mistakes**
```go
// ❌ SAI - Thiếu thông tin
// @Summary Get user
// @Router /user [get]

// ✅ ĐÚNG - Đầy đủ thông tin
// @Summary     Lấy thông tin user theo ID
// @Description Lấy thông tin chi tiết của user dựa trên ID
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} userpb.GetUserResponse
// @Failure     400 {object} common.ErrorResponse
// @Failure     404 {object} common.ErrorResponse
// @Router      /v2/user/{id} [get]
// @Security    BearerAuth
```

---

## 🔍 Ví Dụ Hoàn Chỉnh

### **Organization Service - Get Branch Deals**
```go
// @Summary     Lấy danh sách thương vụ theo chi nhánh
// @Description Lấy danh sách thương vụ thuộc về một chi nhánh cụ thể với phân trang
// @Tags        Organization
// @Accept      json
// @Produce     json
// @Param       branchId path int true "Branch ID"
// @Param       page query int false "Page number" default(1)
// @Param       size query int false "Page size" default(10)
// @Success     200 {object} organizationpb.GetGroupDealsV2Response
// @Failure     400 {object} common.ErrorResponse
// @Failure     404 {object} common.ErrorResponse
// @Failure     500 {object} common.ErrorResponse
// @Router      /v2/org/organization/branch/{branchId}/deals [get]
// @Security    BearerAuth
func (h *GroupDealHandler) GetBranchDeals(ctx context.Context, req *organizationpb.GetBranchDealsRequest) (*organizationpb.GetGroupDealsV2Response, error) {
    // Implementation
}
```

---

## 🎯 Tóm Tắt

| Bước | Hành Động | Lệnh |
|------|-----------|------|
| 1 | Viết handler function | - |
| 2 | Thêm Swagger comments | `// @Summary, @Description, @Tags, etc.` |
| 3 | Generate Swagger | `make swag <service>` |
| 4 | Merge Swagger (nếu cần) | `make merge-swagger` |

**Quy tắc vàng**: **Swagger comments phải được viết TRỰC TIẾP TẠI HANDLER** để đảm bảo tài liệu API chính xác và đầy đủ!