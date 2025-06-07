# Organization Role API Documentation

## Overview

Organization Role API cung cấp các endpoint để quản lý vai trò trong tổ chức, bao gồm tạo mới, cập nhật, xóa và truy vấn thông tin vai trò.

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

### Create Role

```http
POST /v2/organization-role
```

Tạo mới một vai trò trong tổ chức. Chỉ admin mới có quyền tạo vai trò.

**Request Body:**
```json
{
  "name": "string",           // required, max 255 chars
  "key": "string",           // required, max 50 chars, unique
  "organizationId": "uint32"  // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created role
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Organization not found
- 409: Role key already exists
- 500: Internal server error

### Update Role

```http
PUT /v2/organization-role/{id}
```

Cập nhật thông tin của một vai trò. Chỉ admin mới có quyền cập nhật.

**Request Body:**
```json
{
  "id": "uint32",            // required
  "name": "string",          // required, max 255 chars
  "key": "string",           // required, max 50 chars
  "permissionIds": ["uint32"] // optional, list of permission IDs
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated role
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Role not found
- 409: Role key already exists
- 500: Internal server error

### Delete Role

```http
DELETE /v2/organization-role/{id}
```

Xóa một vai trò. Chỉ admin mới có quyền xóa vai trò.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, role ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of deleted role
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - not admin
- 404: Role not found
- 500: Internal server error

### Get Roles

```http
GET /v2/organization-role
```

Lấy danh sách vai trò của một tổ chức.

**Request Parameters:**
```json
{
  "organizationId": "uint32",  // required
  "page": "uint32",           // optional, default: 0
  "size": "uint32"            // optional, default: 10
}
```

**Response:**
```json
{
  "roles": [
    {
      "id": "uint32",
      "name": "string",
      "key": "string"
    }
  ],
  "total": "uint32"  // Tổng số vai trò
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

## Data Models

### Role

```json
{
  "id": "uint32",
  "name": "string",    // max 255 chars
  "key": "string"      // max 50 chars, unique
}
```

## Validation Rules

1. **Role Name**
   - Required
   - Maximum length: 255 characters
   - Cannot contain special characters except spaces and hyphens

2. **Role Key**
   - Required
   - Maximum length: 50 characters
   - Must be unique within the organization
   - Format: lowercase with underscores (e.g., "admin", "member", "viewer")

3. **Organization ID**
   - Required
   - Must be a valid organization ID

4. **Permission IDs**
   - Optional
   - Must be valid permission IDs
   - Maximum 20 permissions per role

## Error Codes

- 400: Bad Request - Invalid input data
- 401: Unauthorized - Missing or invalid authentication
- 403: Forbidden - Insufficient permissions
- 404: Not Found - Resource not found
- 409: Conflict - Resource already exists
- 500: Internal Server Error - Server-side error

## Rate Limiting

- Maximum 100 requests per minute per IP
- Maximum 1000 requests per hour per IP

## Logging

Service logs các sự kiện sau:

1. **Role Events**
   - ROLE_CREATE: Tạo vai trò mới
   - ROLE_UPDATE: Cập nhật vai trò
   - ROLE_DELETE: Xóa vai trò
   - ROLE_PERMISSION_CHANGE: Thay đổi quyền của vai trò 