# Organization API Documentation

## Overview

Organization API cung cấp các endpoint để quản lý thông tin và thành viên của các tổ chức. API này cho phép tạo mới, cập nhật, và truy vấn thông tin tổ chức, cũng như quản lý thành viên và phân quyền.

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

### Organization

#### Create Organization

```http
POST /v2/organization
```

Tạo mới một tổ chức và thêm người tạo làm admin.

**Request Body:**
```json
{
  "name": "string",           // required, max 255 chars
  "taxCode": "string",        // required, max 50 chars, unique
  "businessLicenseUrl": "string", // optional
  "address": "string",        // optional
  "phone": "string",          // optional, valid phone format
  "email": "string",          // optional, valid email format
  "website": "string",        // optional, must start with http:// or https://
  "logoUrl": "string"         // optional
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created organization
}
```

**Error Codes:**
- 400: Invalid request body
- 409: Tax code already exists
- 500: Internal server error

#### Update Organization

```http
PUT /v2/organization/{id}
```

Cập nhật thông tin của một tổ chức. Chỉ admin mới có quyền cập nhật.

**Request Body:**
```json
{
  "id": "uint32",            // required
  "name": "string",          // required, max 255 chars
  "taxCode": "string",       // required, max 50 chars
  "businessLicenseUrl": "string", // optional
  "address": "string",       // optional
  "phone": "string",         // optional, valid phone format
  "email": "string",         // optional, valid email format
  "website": "string",       // optional, must start with http:// or https://
  "logoUrl": "string",       // optional
  "status": "string"         // optional, one of: active, inactive, pending
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated organization
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Organization not found
- 409: Tax code already exists
- 500: Internal server error

#### Get Organizations

```http
GET /v2/organization
```

Lấy danh sách các tổ chức mà user là creator hoặc member. Kết quả được sắp xếp theo thời gian tạo mới nhất.

**Request Parameters:**
```json
{
  "id": "uint32"  // optional, filter by organization ID
}
```

**Response:**
```json
{
  "organizations": [
    {
      "id": "uint32",
      "name": "string",
      "taxCode": "string",
      "businessLicenseUrl": "string",
      "address": "string",
      "phone": "string",
      "email": "string",
      "website": "string",
      "logoUrl": "string",
      "status": "string"
    }
  ]
}
```

**Error Codes:**
- 500: Internal server error

### Organization Member

#### Create Member

```http
POST /v2/organization/member
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
- 403: Forbidden - not admin or max members reached
- 404: Organization not found
- 409: Member already exists
- 500: Internal server error

#### Update Member

```http
PUT /v2/organization/member
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

#### Delete Member

```http
DELETE /v2/organization/member
```

Xóa một thành viên khỏi tổ chức. Chỉ admin mới có quyền xóa thành viên.

**Request Body:**
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
- 400: Invalid request body
- 403: Forbidden - not admin
- 404: Member not found
- 500: Internal server error

## Data Models

### Organization

```json
{
  "id": "uint32",
  "name": "string",           // required, max 255 chars
  "taxCode": "string",        // required, max 50 chars, unique
  "businessLicenseUrl": "string", // optional
  "address": "string",        // optional
  "phone": "string",          // optional, valid phone format
  "email": "string",          // optional, valid email format
  "website": "string",        // optional, must start with http:// or https://
  "logoUrl": "string",        // optional
  "status": "string",         // active, inactive, pending
  "createdAt": "string",
  "updatedAt": "string",
  "createdBy": "uint32",
  "updatedBy": "uint32"
}
```

### Organization Member

```json
{
  "id": "uint32",
  "organizationId": "uint32",
  "userId": "uint32",
  "roleId": "uint32",
  "status": "uint32",         // 1: active, 2: inactive, 3: pending
  "joinedAt": "string",
  "removedAt": "string"
}
```

## Validation Rules

1. **Organization Name**
   - Required
   - Maximum length: 255 characters

2. **Tax Code**
   - Required
   - Maximum length: 50 characters
   - Must be unique across all organizations

3. **Email**
   - Optional
   - Must be a valid email format
   - Example: user@example.com

4. **Phone**
   - Optional
   - Must be a valid phone number format
   - Can include country code
   - Example: +84123456789

5. **Website**
   - Optional
   - Must start with http:// or https://
   - Must be a valid URL format
   - Example: https://example.com

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

1. **Organization Events**
   - ORGANIZATION_CREATE: Tạo mới tổ chức
   - ORGANIZATION_UPDATE: Cập nhật tổ chức

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