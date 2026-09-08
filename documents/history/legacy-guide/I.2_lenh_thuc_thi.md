# I.2 Lệnh Thực Thi - BDSPro Microservices

## 🚀 Các Lệnh Chính

### 1. **make buf-<service>** - Generate Protobuf Files
```bash
# Generate protobuf files cho tất cả services
make buf-all

# Generate protobuf cho từng service riêng lẻ
make buf-social          # Generate protobuf cho social service
make buf-user            # Generate protobuf cho user service
make buf-organization    # Generate protobuf cho organization service
make buf-bdspro          # Generate protobuf cho bdspro service
make buf-crm             # Generate protobuf cho crm service
make buf-notification    # Generate protobuf cho notification service
make buf-appointment     # Generate protobuf cho appointment service
make buf-payment         # Generate protobuf cho payment service
make buf-chat            # Generate protobuf cho chat service
make buf-task            # Generate protobuf cho task service
make buf-marketing       # Generate protobuf cho marketing service
make buf-transaction     # Generate protobuf cho transaction service
make buf-auth            # Generate protobuf cho auth service
make buf-shared          # Generate protobuf cho shared components
```

**Mục đích**: Tạo ra các file Go code từ định nghĩa protobuf (.proto files) để giao tiếp gRPC giữa các services.

**Khi nào sử dụng**: 
- Sau khi thay đổi file .proto
- Khi thêm/sửa/xóa RPC methods
- Khi thay đổi message structures
- Trước khi build hoặc chạy services

---

### 2. **make swag <service>** - Tạo Swagger Documentation
```bash
# Tạo Swagger docs cho từng service
make swag user           # Tạo Swagger cho user service
make swag organization   # Tạo Swagger cho organization service
make swag social         # Tạo Swagger cho social service
make swag bdspro         # Tạo Swagger cho bdspro service
make swag crm            # Tạo Swagger cho crm service
make swag chat           # Tạo Swagger cho chat service
make swag payment        # Tạo Swagger cho payment service
make swag notification   # Tạo Swagger cho notification service
make swag appointment    # Tạo Swagger cho appointment service
make swag marketing      # Tạo Swagger cho marketing service
make swag transaction    # Tạo Swagger cho transaction service
make swag auth           # Tạo Swagger cho auth service

# Merge tất cả Swagger docs
make merge-swagger

# Tạo Swagger cho file service (lệnh đặc biệt)
make swag-file
```

**Mục đích**: Tạo ra tài liệu API Swagger từ các comment trong code Go, giúp developers hiểu và test APIs.

**Khi nào sử dụng**:
- Sau khi thêm/sửa API endpoints
- Khi thay đổi request/response structures
- Trước khi deploy để cập nhật API documentation
- Khi cần merge tất cả Swagger docs

---

### 3. **make wire <service>** - Dependency Injection với Wire
```bash
# Generate wire dependencies cho từng service
make wire user           # Generate wire cho user service
make wire organization   # Generate wire cho organization service
make wire social         # Generate wire cho social service
make wire bdspro         # Generate wire cho bdspro service
make wire crm            # Generate wire cho crm service
make wire chat           # Generate wire cho chat service
make wire payment        # Generate wire cho payment service
make wire notification   # Generate wire cho notification service
make wire appointment    # Generate wire cho appointment service
make wire marketing      # Generate wire cho marketing service
make wire transaction    # Generate wire cho transaction service
make wire auth           # Generate wire cho auth service
```

**Mục đích**: Tự động tạo ra dependency injection code, giúp quản lý và inject các dependencies (repositories, usecases, handlers) một cách tự động.

**Khi nào sử dụng**:
- Sau khi thêm/sửa repositories
- Khi thêm/sửa usecases
- Khi thêm/sửa handlers
- Khi thay đổi dependency structure
- Trước khi build hoặc chạy services

---

## 🔄 Workflow Thực Thi

### **Khi thêm API mới:**
```bash
# 1. Thay đổi file .proto
# 2. Generate protobuf
make buf-<service>

# 3. Implement code (handler, usecase, repository)
# 4. Generate wire dependencies
make wire <service>

# 5. Generate Swagger docs
make swag <service>

# 6. Test service
make run-dev <service>
```

### **Khi thay đổi database schema:**
```bash
# 1. Thay đổi entity structures
# 2. Generate wire dependencies
make wire <service>

# 3. Test service
make run-dev <service>
```

### **Khi deploy:**
```bash
# 1. Generate tất cả protobuf
make buf-all

# 2. Generate tất cả wire dependencies
make wire <service1>
make wire <service2>
# ... cho tất cả services

# 3. Generate tất cả Swagger
make swag <service1>
make swag <service2>
# ... cho tất cả services

# 4. Merge Swagger docs
make merge-swagger

# 5. Build và deploy
make build-all
```

---

## 📋 Lệnh Kết Hợp Thường Dùng

### **Setup môi trường mới:**
```bash
# 1. Install dependencies
make init

# 2. Generate tất cả protobuf
make buf-all

# 3. Generate wire cho tất cả services
make wire user
make wire organization
make wire social
make wire bdspro
make wire crm
make wire chat
make wire payment
make wire notification
make wire appointment
make wire marketing
make wire transaction
make wire auth

# 4. Generate Swagger cho tất cả services
make swag user
make swag organization
make swag social
make swag bdspro
make swag crm
make swag chat
make swag payment
make swag notification
make swag appointment
make swag marketing
make swag transaction
make swag auth

# 5. Merge Swagger docs
make merge-swagger

# 6. Start tất cả services
make start-all
```

### **Development workflow:**
```bash
# 1. Start service với hot reload
make run-dev <service>

# 2. Trong terminal khác, monitor services
make dev

# 3. Khi thay đổi code, Air sẽ tự động reload
# 4. Khi thay đổi .proto, chạy:
make buf-<service>

# 5. Khi thay đổi dependencies, chạy:
make wire <service>

# 6. Khi thay đổi API, chạy:
make swag <service>
```

---

## ⚠️ Lưu Ý Quan Trọng

### **Thứ tự thực thi:**
1. **buf** → **wire** → **swag** → **run**
2. Luôn chạy `make buf-<service>` trước khi chạy `make wire <service>`
3. Luôn chạy `make wire <service>` trước khi chạy `make swag <service>`

### **Khi gặp lỗi:**
```bash
# Nếu lỗi protobuf
make buf-all

# Nếu lỗi dependency injection
make wire <service>

# Nếu lỗi Swagger
make swag <service>

# Nếu lỗi build
make init
make buf-all
make wire <service>
make build-all
```

### **Best Practices:**
- Luôn chạy `make buf-all` sau khi thay đổi .proto files
- Luôn chạy `make wire <service>` sau khi thay đổi dependencies
- Luôn chạy `make swag <service>` sau khi thay đổi API
- Sử dụng `make run-dev <service>` cho development
- Sử dụng `make start-all` cho testing toàn bộ hệ thống

---

## 🎯 Tóm Tắt

| Lệnh | Mục đích | Khi nào sử dụng |
|------|----------|-----------------|
| `make buf-<service>` | Generate protobuf files | Sau khi thay đổi .proto |
| `make swag <service>` | Tạo Swagger documentation | Sau khi thay đổi API |
| `make wire <service>` | Dependency injection | Sau khi thay đổi dependencies |

**Workflow chuẩn**: `buf` → `wire` → `swag` → `run`