# Organization Member API Documentation

## Overview

Organization Member API cung cấp các endpoint để quản lý thành viên trong tổ chức, bao gồm thêm, cập nhật, xóa và truy vấn thông tin thành viên.

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

### Switch Organization

```http
POST /v2/organization-member/switch
```

Chuyển đổi tổ chức hiện tại của user.

**Request Body:**
```json
{
  "organizationId": "uint32",  // required
  "userId": "uint32"          // required
}
```

**Response:**
```json
{
  "id": "uint32",
  "userId": "uint32",
  "isSuccess": "boolean",
  "role": "string",
  "permissions": ["string"]
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Organization not found
- 500: Internal server error

### Create Organization Member

```http
POST /v2/organization-member
```

Thêm một thành viên mới vào tổ chức. Chỉ admin mới có quyền thêm thành viên.

**Request Body:**
```json
{
  "organizationId": "uint32",  // required
  "userId": "uint32",         // required
  "role": "uint32"            // required, role ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created member
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Organization not found
- 409: Member already exists
- 500: Internal server error

### Update Organization Member

```http
PUT /v2/organization-member/{id}
```

Cập nhật vai trò của thành viên. Chỉ admin mới có quyền cập nhật.

**Request Body:**
```json
{
  "id": "uint32",      // required
  "role": "uint32"     // required, role ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated member
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Member not found
- 500: Internal server error

### Delete Organization Member

```http
DELETE /v2/organization-member/{id}
```

Xóa một thành viên khỏi tổ chức. Chỉ admin mới có quyền xóa thành viên.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, member ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of deleted member
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - not admin
- 404: Member not found
- 500: Internal server error

### Get Organization Members

```http
GET /v2/organization-member/{organizationId}
```

Lấy danh sách thành viên của một tổ chức.

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
  "organizationMembers": [
    {
      "id": "uint32",
      "userId": "uint32",
      "role": "uint32",
      "status": "string"
    }
  ],
  "total": "uint32"  // Tổng số thành viên
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

## Data Models

### Organization Member

```json
{
  "id": "uint32",
  "userId": "uint32",
  "role": "uint32",
  "status": "string"  // active, inactive, pending
}
```

## Validation Rules

1. **User ID**
   - Required
   - Must be a valid user ID

2. **Role ID**
   - Required
   - Must be a valid role ID in the organization

3. **Organization ID**
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

1. **Member Events**
   - MEMBER_CREATE: Thêm thành viên
   - MEMBER_UPDATE: Cập nhật thành viên
   - MEMBER_DELETE: Xóa thành viên
   - MEMBER_SWITCH: Chuyển đổi tổ chức 