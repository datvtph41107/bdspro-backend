# Organization Log API Documentation

## Overview

Organization Log API cung cấp các endpoint để xem và quản lý lịch sử hoạt động trong tổ chức, bao gồm các sự kiện liên quan đến tổ chức, thành viên, vai trò, quyền và chi nhánh.

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

### Get Organization Logs

```http
GET /v2/organization/log/{id}
```

Lấy thông tin chi tiết của một log.

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
    "organizationId": "uint32",
    "userId": "uint32",
    "eventType": "string",
    "eventData": "string",     // JSON string containing event details
    "createdAt": "string"      // ISO 8601 format
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Log not found
- 500: Internal server error

### Get Organization Logs By Organization ID

```http
GET /v2/organization/{organization_id}/logs
```

Lấy danh sách log của một tổ chức. Có thể lọc theo loại sự kiện, thời gian và người thực hiện.

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
      "organizationId": "uint32",
      "userId": "uint32",
      "eventType": "string",
      "eventData": "string",     // JSON string containing event details
      "createdAt": "string"      // ISO 8601 format
    }
  ],
  "total": "uint32"  // Tổng số bản ghi
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

### Get Organization Logs By Actor ID

```http
GET /v2/organization/logs/actor/{actor_id}
```

Lấy danh sách log của một người dùng trong tổ chức.

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
      "organizationId": "uint32",
      "userId": "uint32",
      "eventType": "string",
      "eventData": "string",     // JSON string containing event details
      "createdAt": "string"      // ISO 8601 format
    }
  ],
  "total": "uint32"  // Tổng số bản ghi
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Actor not found
- 500: Internal server error

### Get Organization Member Logs

```http
GET /v2/organization-log/member
```

Lấy danh sách log liên quan đến thành viên của tổ chức.

**Request Parameters:**
```json
{
  "organizationId": "uint32",    // required
  "memberId": "uint32",          // optional, filter by member ID
  "startTime": "string",         // optional, ISO 8601 format
  "endTime": "string",           // optional, ISO 8601 format
  "page": "uint32",              // optional, default: 0
  "size": "uint32"               // optional, default: 10
}
```

**Response:**
```json
{
  "logs": [
    {
      "id": "uint32",
      "organizationId": "uint32",
      "userId": "uint32",
      "memberId": "uint32",
      "eventType": "string",     // MEMBER_CREATE, MEMBER_UPDATE, MEMBER_DELETE
      "eventData": "string",     // JSON string containing event details
      "createdAt": "string"      // ISO 8601 format
    }
  ],
  "total": "uint32"  // Tổng số bản ghi
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

### Get Organization Role Logs

```http
GET /v2/organization-log/role
```

Lấy danh sách log liên quan đến vai trò của tổ chức.

**Request Parameters:**
```json
{
  "organizationId": "uint32",    // required
  "roleId": "uint32",           // optional, filter by role ID
  "startTime": "string",         // optional, ISO 8601 format
  "endTime": "string",           // optional, ISO 8601 format
  "page": "uint32",              // optional, default: 1
  "pageSize": "uint32"           // optional, default: 20, max: 100
}
```

**Response:**
```json
{
  "logs": [
    {
      "id": "uint32",
      "organizationId": "uint32",
      "userId": "uint32",
      "roleId": "uint32",
      "eventType": "string",     // ROLE_CREATE, ROLE_UPDATE, ROLE_DELETE
      "eventData": "string",     // JSON string containing event details
      "createdAt": "string"      // ISO 8601 format
    }
  ],
  "total": "uint32",
  "page": "uint32",
  "pageSize": "uint32"
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

### Get Organization Permission Logs

```http
GET /v2/organization-log/permission
```

Lấy danh sách log liên quan đến quyền của tổ chức.

**Request Parameters:**
```json
{
  "organizationId": "uint32",    // required
  "permissionId": "uint32",      // optional, filter by permission ID
  "startTime": "string",         // optional, ISO 8601 format
  "endTime": "string",           // optional, ISO 8601 format
  "page": "uint32",              // optional, default: 1
  "pageSize": "uint32"           // optional, default: 20, max: 100
}
```

**Response:**
```json
{
  "logs": [
    {
      "id": "uint32",
      "organizationId": "uint32",
      "userId": "uint32",
      "permissionId": "uint32",
      "eventType": "string",     // PERMISSION_CREATE, PERMISSION_UPDATE, PERMISSION_DELETE
      "eventData": "string",     // JSON string containing event details
      "createdAt": "string"      // ISO 8601 format
    }
  ],
  "total": "uint32",
  "page": "uint32",
  "pageSize": "uint32"
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

### Get Organization Branch Logs

```http
GET /v2/organization-log/branch
```

Lấy danh sách log liên quan đến chi nhánh của tổ chức.

**Request Parameters:**
```json
{
  "organizationId": "uint32",    // required
  "branchId": "uint32",         // optional, filter by branch ID
  "startTime": "string",         // optional, ISO 8601 format
  "endTime": "string",           // optional, ISO 8601 format
  "page": "uint32",              // optional, default: 1
  "pageSize": "uint32"           // optional, default: 20, max: 100
}
```

**Response:**
```json
{
  "logs": [
    {
      "id": "uint32",
      "organizationId": "uint32",
      "userId": "uint32",
      "branchId": "uint32",
      "eventType": "string",     // BRANCH_CREATE, BRANCH_UPDATE, BRANCH_DELETE
      "eventData": "string",     // JSON string containing event details
      "createdAt": "string"      // ISO 8601 format
    }
  ],
  "total": "uint32",
  "page": "uint32",
  "pageSize": "uint32"
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

## Data Models

### Organization Log

```json
{
  "id": "uint32",
  "organizationId": "uint32",
  "userId": "uint32",
  "eventType": "string",
  "eventData": "string",     // JSON string containing event details
  "createdAt": "string"      // ISO 8601 format
}
```

## Event Types

1. **Organization Events**
   - ORGANIZATION_CREATE: Tạo mới tổ chức
   - ORGANIZATION_UPDATE: Cập nhật tổ chức
   - ORGANIZATION_DELETE: Xóa tổ chức

2. **Member Events**
   - MEMBER_CREATE: Thêm thành viên
   - MEMBER_UPDATE: Cập nhật thành viên
   - MEMBER_DELETE: Xóa thành viên

3. **Role Events**
   - ROLE_CREATE: Tạo vai trò
   - ROLE_UPDATE: Cập nhật vai trò
   - ROLE_DELETE: Xóa vai trò

4. **Permission Events**
   - PERMISSION_CREATE: Tạo quyền
   - PERMISSION_UPDATE: Cập nhật quyền
   - PERMISSION_DELETE: Xóa quyền
   - PERMISSION_ASSIGN: Gán quyền cho vai trò
   - PERMISSION_REMOVE: Xóa quyền khỏi vai trò

5. **Branch Events**
   - BRANCH_CREATE: Tạo chi nhánh
   - BRANCH_UPDATE: Cập nhật chi nhánh
   - BRANCH_DELETE: Xóa chi nhánh
   - BRANCH_MEMBER_ADD: Thêm thành viên vào chi nhánh
   - BRANCH_MEMBER_REMOVE: Xóa thành viên khỏi chi nhánh
   - BRANCH_MANAGER_CHANGE: Thay đổi quản lý chi nhánh

## Validation Rules

1. **Time Range**
   - startTime và endTime phải theo định dạng ISO 8601
   - endTime phải lớn hơn startTime
   - Khoảng thời gian tối đa: 30 ngày

2. **Pagination**
   - page: số nguyên dương, mặc định: 1
   - pageSize: số nguyên dương, mặc định: 20, tối đa: 100

3. **Organization ID**
   - Required
   - Must be a valid organization ID

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

- Logs được lưu trữ trong 90 ngày
- Logs cũ hơn 90 ngày sẽ tự động bị xóa
- Có thể export logs trong khoảng thời gian 30 ngày 