# Organization Branch API Documentation

## Overview

Organization Branch API cung cấp các endpoint để quản lý chi nhánh trong tổ chức, bao gồm tạo mới, cập nhật, xóa và quản lý thành viên của chi nhánh.

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

### Create Organization Branch

```http
POST /v2/organization-branch
```

Tạo mới một chi nhánh trong tổ chức. Chỉ admin mới có quyền tạo chi nhánh.

**Request Body:**
```json
{
  "name": "string",           // required, max 255 chars
  "address": "string",        // required
  "phone": "string",          // required, valid phone format
  "email": "string",          // required, valid email format
  "managerId": "uint32",      // required, user ID of branch manager
  "organizationId": "uint32"  // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created branch
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Organization not found
- 409: Branch name already exists
- 500: Internal server error

### Update Organization Branch

```http
PUT /v2/organization-branch/{id}
```

Cập nhật thông tin của một chi nhánh. Chỉ admin mới có quyền cập nhật.

**Request Body:**
```json
{
  "id": "uint32",            // required
  "name": "string",          // required, max 255 chars
  "address": "string",       // required
  "phone": "string",         // required, valid phone format
  "email": "string"          // required, valid email format
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated branch
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Branch not found
- 409: Branch name already exists
- 500: Internal server error

### Delete Organization Branch

```http
DELETE /v2/organization-branch/{id}
```

Xóa một chi nhánh. Chỉ admin mới có quyền xóa chi nhánh.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, branch ID
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
- 404: Branch not found
- 500: Internal server error

### Add Member to Organization Branch

```http
POST /v2/organization-branch/{id}/add-member
```

Thêm một thành viên vào chi nhánh. Chỉ admin mới có quyền thêm thành viên.

**Request Body:**
```json
{
  "id": "uint32",      // required, branch ID
  "userId": "uint32"   // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of added member
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Branch not found
- 409: Member already exists
- 500: Internal server error

### Remove Member from Organization Branch

```http
DELETE /v2/organization-branch/{id}/remove-member/{userId}
```

Xóa một thành viên khỏi chi nhánh. Chỉ admin mới có quyền xóa thành viên.

**Request Parameters:**
```json
{
  "id": "uint32",      // required, branch ID
  "userId": "uint32"   // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of removed member
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - not admin
- 404: Branch or member not found
- 500: Internal server error

### Get Organization Branch Members

```http
GET /v2/organization-branch/{id}/members
```

Lấy danh sách thành viên của một chi nhánh.

**Request Parameters:**
```json
{
  "id": "uint32",        // required, branch ID
  "page": "uint32",      // optional, default: 0
  "size": "uint32"       // optional, default: 10
}
```

**Response:**
```json
{
  "members": [
    {
      "id": "uint32",
      "user_id": "uint32"
    }
  ],
  "total": "uint32"  // Tổng số thành viên
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Branch not found
- 500: Internal server error

### Get Organization Branches

```http
GET /v2/organization-branch
```

Lấy danh sách chi nhánh của một tổ chức.

**Request Parameters:**
```json
{
  "organizationId": "uint32"  // required
}
```

**Response:**
```json
{
  "branches": [
    {
      "id": "uint32",
      "name": "string",
      "address": "string",
      "phone": "string",
      "email": "string"
    }
  ]
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Organization not found
- 500: Internal server error

## Data Models

### Organization Branch

```json
{
  "id": "uint32",
  "name": "string",     // max 255 chars
  "address": "string",
  "phone": "string",    // valid phone format
  "email": "string"     // valid email format
}
```

### Organization Branch Member

```json
{
  "id": "uint32",
  "user_id": "uint32"
}
```

## Validation Rules

1. **Branch Name**
   - Required
   - Maximum length: 255 characters
   - Must be unique within the organization

2. **Address**
   - Required
   - Maximum length: 500 characters

3. **Phone**
   - Required
   - Must be a valid phone number format
   - Can include country code
   - Example: +84123456789

4. **Email**
   - Required
   - Must be a valid email format
   - Example: branch@example.com

5. **Manager ID**
   - Required
   - Must be a valid user ID
   - User must be a member of the organization

6. **Organization ID**
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

1. **Branch Events**
   - BRANCH_CREATE: Tạo chi nhánh mới
   - BRANCH_UPDATE: Cập nhật thông tin chi nhánh
   - BRANCH_DELETE: Xóa chi nhánh

2. **Member Events**
   - BRANCH_MEMBER_ADD: Thêm thành viên vào chi nhánh
   - BRANCH_MEMBER_REMOVE: Xóa thành viên khỏi chi nhánh
   - BRANCH_MANAGER_CHANGE: Thay đổi quản lý chi nhánh 