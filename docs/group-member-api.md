# Group Member API Documentation

## Overview

Group Member API cung cấp các endpoint để quản lý thành viên trong nhóm, bao gồm thêm, cập nhật, xóa và truy vấn thông tin thành viên.

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

### Create Group Member

```http
POST /v2/group-member
```

Thêm một thành viên mới vào nhóm.

**Request Body:**
```json
{
  "groupId": "uint32",  // required
  "userId": "uint32",   // required
  "role": "string"      // required, one of: admin, member
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
- 404: Group not found
- 409: Member already exists
- 500: Internal server error

### Update Group Member

```http
PUT /v2/group-member/{id}
```

Cập nhật vai trò của thành viên trong nhóm.

**Request Body:**
```json
{
  "id": "uint32",     // required
  "role": "string"    // required, one of: admin, member
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
- 404: Member not found
- 500: Internal server error

### Delete Group Member

```http
DELETE /v2/group-member/{id}
```

Xóa một thành viên khỏi nhóm.

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
- 404: Member not found
- 500: Internal server error

### Get Group Member

```http
GET /v2/group-member/{id}
```

Lấy thông tin chi tiết của một thành viên trong nhóm.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, member ID
}
```

**Response:**
```json
{
  "member": {
    "id": "uint32",
    "groupId": "uint32",
    "userId": "uint32",
    "role": "string"
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Member not found
- 500: Internal server error

### Get Group Members

```http
GET /v2/group-member/{groupId}
```

Lấy danh sách thành viên của một nhóm.

**Request Parameters:**
```json
{
  "groupId": "uint32",  // required
  "page": "uint32",     // optional, default: 0
  "size": "uint32"      // optional, default: 10
}
```

**Response:**
```json
{
  "members": [
    {
      "id": "uint32",
      "userId": "uint32",
      "role": "string",
      "status": "string"
    }
  ],
  "total": "uint32"  // Tổng số thành viên
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Group not found
- 500: Internal server error

## Data Models

### Group Member

```json
{
  "id": "uint32",
  "groupId": "uint32",
  "userId": "uint32",
  "role": "string"  // admin, member
}
```

## Validation Rules

1. **Group ID**
   - Required
   - Must be a valid group ID

2. **User ID**
   - Required
   - Must be a valid user ID
   - User must not already be a member of the group

3. **Role**
   - Required
   - Must be one of: admin, member
   - Only one admin per group

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
   - MEMBER_CREATE: Thêm thành viên mới
   - MEMBER_UPDATE: Cập nhật vai trò thành viên
   - MEMBER_DELETE: Xóa thành viên 