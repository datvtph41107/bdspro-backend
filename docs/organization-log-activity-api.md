# Organization Log Activity API Documentation

## Overview

Organization Log Activity API cung cấp các endpoint để truy vấn lịch sử hoạt động của tổ chức, bao gồm các sự kiện như tạo tổ chức, cập nhật thông tin, thêm/xóa thành viên, và các hoạt động khác.

## Base URL

```
http://localhost:8080/v2
```

## Authentication

Tất cả các request cần phải có header `Authorization` với JWT token:

```
Authorization: Bearer <token>
```

## API Endpoints

### Get Organization Log Activity By ID

```http
GET /v2/organization/log/{id}
```

Lấy thông tin chi tiết về một hoạt động trong tổ chức.

**Request Parameters:**
```json
{
  "id": "uint32"  // required
}
```

**Response:**
```json
{
  "log": {
    "id": "uint32",
    "created_at": "int64",
    "updated_at": "int64",
    "created_by": "uint32",
    "updated_by": "uint32",
    "organization_id": "uint32",
    "actor_id": "uint32",
    "log_type": "string",
    "log_data": "string"
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Log not found
- 500: Internal server error

### Get Organization Log Activities By Organization ID

```http
GET /v2/organization/{organization_id}/logs
```

Lấy danh sách các hoạt động của một tổ chức.

**Request Parameters:**
```json
{
  "organization_id": "uint32",  // required
  "page": "uint32",            // optional, default: 0
  "size": "uint32"             // optional, default: 10
}
```

**Response:**
```json
{
  "logs": [
    {
      "id": "uint32",
      "created_at": "int64",
      "updated_at": "int64",
      "created_by": "uint32",
      "updated_by": "uint32",
      "organization_id": "uint32",
      "actor_id": "uint32",
      "log_type": "string",
      "log_data": "string"
    }
  ],
  "total": "uint32"  // Tổng số bản ghi
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

### Get Organization Log Activities By Actor ID

```http
GET /v2/organization/logs/actor/{actor_id}
```

Lấy danh sách các hoạt động của một người dùng trong tổ chức.

**Request Parameters:**
```json
{
  "actor_id": "uint32",  // required
  "page": "uint32",      // optional, default: 0
  "size": "uint32"       // optional, default: 10
}
```

**Response:**
```json
{
  "logs": [
    {
      "id": "uint32",
      "created_at": "int64",
      "updated_at": "int64",
      "created_by": "uint32",
      "updated_by": "uint32",
      "organization_id": "uint32",
      "actor_id": "uint32",
      "log_type": "string",
      "log_data": "string"
    }
  ],
  "total": "uint32"  // Tổng số bản ghi
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Actor not found
- 500: Internal server error

## Data Models

### Organization Log Activity

```json
{
  "id": "uint32",
  "created_at": "int64",
  "updated_at": "int64",
  "created_by": "uint32",
  "updated_by": "uint32",
  "organization_id": "uint32",
  "actor_id": "uint32",
  "log_type": "string",
  "log_data": "string"
}
```

## Log Types

1. **Organization Events**
   - ORGANIZATION_CREATE: Tạo tổ chức mới
   - ORGANIZATION_UPDATE: Cập nhật thông tin tổ chức
   - ORGANIZATION_DELETE: Xóa tổ chức

2. **Member Events**
   - MEMBER_ADD: Thêm thành viên
   - MEMBER_REMOVE: Xóa thành viên
   - MEMBER_UPDATE: Cập nhật vai trò thành viên
   - MEMBER_SWITCH: Chuyển đổi tổ chức

3. **Role Events**
   - ROLE_CREATE: Tạo vai trò mới
   - ROLE_UPDATE: Cập nhật vai trò
   - ROLE_DELETE: Xóa vai trò

4. **Permission Events**
   - PERMISSION_CREATE: Tạo quyền mới
   - PERMISSION_UPDATE: Cập nhật quyền
   - PERMISSION_DELETE: Xóa quyền
   - PERMISSION_ASSIGN: Gán quyền cho vai trò

5. **Branch Events**
   - BRANCH_CREATE: Tạo chi nhánh mới
   - BRANCH_UPDATE: Cập nhật chi nhánh
   - BRANCH_DELETE: Xóa chi nhánh
   - BRANCH_MEMBER_ADD: Thêm thành viên vào chi nhánh
   - BRANCH_MEMBER_REMOVE: Xóa thành viên khỏi chi nhánh

## Validation Rules

1. **ID Fields**
   - Must be positive integers
   - Must exist in the database

2. **Timestamps**
   - Unix timestamp in seconds
   - Must be valid timestamps

3. **Log Type**
   - Must be one of the predefined log types
   - Case sensitive

4. **Log Data**
   - Must be valid JSON string
   - Maximum length: 1000 characters

## Error Codes

- 400: Bad Request - Invalid input data
- 401: Unauthorized - Missing or invalid authentication
- 403: Forbidden - Insufficient permissions
- 404: Not Found - Resource not found
- 500: Internal Server Error - Server-side error

## Rate Limiting

- Maximum 100 requests per minute per IP
- Maximum 1000 requests per hour per IP

## Log Retention

- Logs được lưu trữ trong 30 ngày
- Có thể xuất logs trong khoảng thời gian này
- Sau 30 ngày, logs sẽ được tự động xóa 