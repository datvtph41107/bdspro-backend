# Group Deal API Documentation

## Overview

Group Deal API cung cấp các endpoint để quản lý giao dịch trong nhóm, bao gồm tạo mới, cập nhật, xóa và truy vấn thông tin giao dịch.

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

### Create Group Deal

```http
POST /v2/group-deal
```

Tạo mới một giao dịch trong nhóm.

**Request Body:**
```json
{
  "name": "string",     // required, max 255 chars
  "amount": "double",   // required, positive number
  "status": "string",   // required, one of: pending, approved, rejected
  "groupId": "uint32"   // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created deal
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Group not found
- 500: Internal server error

### Update Group Deal

```http
PUT /v2/group-deal/{id}
```

Cập nhật thông tin của một giao dịch.

**Request Body:**
```json
{
  "id": "uint32",       // required
  "name": "string",     // required, max 255 chars
  "amount": "double",   // required, positive number
  "status": "string"    // required, one of: pending, approved, rejected
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated deal
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Deal not found
- 500: Internal server error

### Delete Group Deal

```http
DELETE /v2/group-deal/{id}
```

Xóa một giao dịch.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, deal ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of deleted deal
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Deal not found
- 500: Internal server error

### Get Group Deal

```http
GET /v2/group-deal/{id}
```

Lấy thông tin chi tiết của một giao dịch.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, deal ID
}
```

**Response:**
```json
{
  "deal": {
    "id": "uint32",
    "name": "string",
    "amount": "double",
    "status": "string"
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Deal not found
- 500: Internal server error

### Get Group Deals

```http
GET /v2/group-deal
```

Lấy danh sách giao dịch của một nhóm. Kết quả được phân trang.

**Request Parameters:**
```json
{
  "groupId": "uint32",    // required
  "page": "uint32",       // optional, default: 1
  "size": "uint32"        // optional, default: 20, max: 100
}
```

**Response:**
```json
{
  "deals": [
    {
      "id": "uint32",
      "name": "string",
      "amount": "double",
      "status": "string"
    }
  ],
  "total": "uint32"
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Group not found
- 500: Internal server error

### Update Group Deal Status

```http
PUT /v2/group-deal/{id}/status
```

Cập nhật trạng thái của một giao dịch.

**Request Body:**
```json
{
  "id": "uint32",     // required
  "status": "string"  // required, one of: pending, approved, rejected
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated deal
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Deal not found
- 500: Internal server error

## Data Models

### Group Deal

```json
{
  "id": "uint32",
  "name": "string",     // max 255 chars
  "amount": "double",   // positive number
  "status": "string"    // pending, approved, rejected
}
```

## Validation Rules

1. **Deal Name**
   - Required
   - Maximum length: 255 characters

2. **Amount**
   - Required
   - Must be a positive number
   - Maximum value: 999999999.99

3. **Status**
   - Required
   - Must be one of: pending, approved, rejected
   - Only admin can change status from pending to approved/rejected

4. **Group ID**
   - Required
   - Must be a valid group ID

## Error Codes

- 400: Bad Request - Invalid input data
- 401: Unauthorized - Missing or invalid authentication
- 403: Forbidden - Insufficient permissions
- 404: Not Found - Resource not found
- 500: Internal Server Error - Server-side error

## Rate Limiting

- Maximum 100 requests per minute per IP
- Maximum 1000 requests per hour per IP

## Logging

Service logs các sự kiện sau:

1. **Deal Events**
   - DEAL_CREATE: Tạo giao dịch mới
   - DEAL_UPDATE: Cập nhật thông tin giao dịch
   - DEAL_DELETE: Xóa giao dịch
   - DEAL_STATUS_CHANGE: Thay đổi trạng thái giao dịch 