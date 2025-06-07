# Organization Permission API Documentation

## Overview

Organization Permission API cung cấp các endpoint để quản lý quyền trong tổ chức, bao gồm tạo mới, cập nhật, xóa và phân quyền cho các vai trò.

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

### Create Permission

```http
POST /v2/organization-permission
```

Tạo mới một quyền trong tổ chức. Chỉ admin mới có quyền tạo quyền.

**Request Body:**
```json
{
  "name": "string",           // required, max 255 chars
  "key": "string",           // required, max 50 chars, unique
  "description": "string",    // optional
  "organizationId": "uint32"  // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created permission
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Organization not found
- 409: Permission key already exists
- 500: Internal server error

### Update Permission

```http
PUT /v2/organization-permission/{id}
```

Cập nhật thông tin của một quyền. Chỉ admin mới có quyền cập nhật.

**Request Body:**
```json
{
  "id": "uint32",           // required
  "name": "string",         // required, max 255 chars
  "key": "string",          // required, max 50 chars
  "description": "string"   // optional
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated permission
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Permission not found
- 409: Permission key already exists
- 500: Internal server error

### Delete Permission

```http
DELETE /v2/organization-permission/{id}
```

Xóa một quyền. Chỉ admin mới có quyền xóa quyền.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, permission ID
}
```

**Response:**
```json
{
  "success": "boolean"
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - not admin
- 404: Permission not found
- 500: Internal server error

### Get Permissions

```http
GET /v2/organization-permission
```

Lấy danh sách quyền của một tổ chức.

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
  "permissions": [
    {
      "id": "uint32",
      "name": "string",
      "key": "string"
    }
  ],
  "total": "uint32"  // Tổng số quyền
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

### Assign Permission to Role

```http
POST /v2/organization-permission/assign
```

Gán quyền cho một vai trò. Chỉ admin mới có quyền gán quyền.

**Request Body:**
```json
{
  "roleId": "uint32",        // required
  "permissionIds": ["uint32"] // required, array of permission IDs
}
```

**Response:**
```json
{
  "success": "boolean"
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Role or permission not found
- 500: Internal server error

### Remove Permission from Role

```http
DELETE /v2/organization-permission/assign
```

Xóa quyền khỏi một vai trò. Chỉ admin mới có quyền xóa quyền.

**Request Body:**
```json
{
  "roleId": "uint32",        // required
  "permissionIds": ["uint32"] // required, array of permission IDs
}
```

**Response:**
```json
{
  "success": "boolean"
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Role or permission not found
- 500: Internal server error

## Data Models

### Permission

```json
{
  "id": "uint32",
  "name": "string",     // max 255 chars
  "key": "string",      // max 50 chars, unique
  "description": "string"
}
```

### Role Permission

```json
{
  "roleId": "uint32",
  "permissionId": "uint32"
}
```

## Validation Rules

1. **Permission Name**
   - Required
   - Maximum length: 255 characters

2. **Permission Key**
   - Required
   - Maximum length: 50 characters
   - Must be unique within the organization
   - Only alphanumeric characters and underscores allowed
   - Example: manage_users, view_reports

3. **Description**
   - Optional
   - Maximum length: 500 characters

4. **Organization ID**
   - Required
   - Must be a valid organization ID

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

1. **Permission Events**
   - PERMISSION_CREATE: Tạo quyền mới
   - PERMISSION_UPDATE: Cập nhật thông tin quyền
   - PERMISSION_DELETE: Xóa quyền

2. **Role Permission Events**
   - PERMISSION_ASSIGN: Gán quyền cho vai trò
   - PERMISSION_REMOVE: Xóa quyền khỏi vai trò 