# Script Export Permissions & Role Groups

Script này dùng để export dữ liệu từ bảng `organization_permissions` và `role_permissions` ra file JSON.

## Cấu trúc dữ liệu

### Permissions
- **Format**: Hierarchy tree structure
- **Cấu trúc**: Các permission root sẽ chứa các child permissions trong mảng `childs`

### Role Groups
- **Format**: Array of role-permission mappings
- **Cấu trúc**: Danh sách các mapping giữa role và permission

## Cách sử dụng

### 1. Sử dụng script shell (Khuyến nghị)

```bash
cd organization-service/scripts
./run_export.sh
```

### 2. Chạy trực tiếp với Go

```bash
cd organization-service/scripts
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=organization_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
go run export_permissions_roles.go
```

## Output

File JSON sẽ được tạo với tên `permissions_roles_export.json`.

Xem file `example_output.json` để biết cấu trúc output.
