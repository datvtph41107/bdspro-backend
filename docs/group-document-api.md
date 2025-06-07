# Group Document API Documentation

## Overview

Group Document API cung cấp các endpoint để quản lý tài liệu trong nhóm, bao gồm tạo mới, cập nhật, xóa và truy vấn thông tin tài liệu.

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

### Create Group Document

```http
POST /v2/group/{groupId}/documents
```

Tạo mới một tài liệu trong nhóm.

**Request Body:**
```json
{
  "groupId": "uint32",      // required
  "name": "string",         // required, max 255 chars
  "description": "string",  // optional
  "fileUrl": "string",      // required, valid URL
  "fileType": "string",     // required, valid file type
  "fileSize": "uint64"      // required, file size in bytes
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created document
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Group not found
- 500: Internal server error

### Update Group Document

```http
PUT /v2/documents/{id}
```

Cập nhật thông tin của một tài liệu.

**Request Body:**
```json
{
  "id": "uint32",          // required
  "name": "string",        // required, max 255 chars
  "description": "string", // optional
  "fileUrl": "string",     // required, valid URL
  "fileType": "string",    // required, valid file type
  "fileSize": "uint64"     // required, file size in bytes
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated document
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Document not found
- 500: Internal server error

### Delete Group Document

```http
DELETE /v2/documents/{id}
```

Xóa một tài liệu.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, document ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of deleted document
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Document not found
- 500: Internal server error

### Get Group Document

```http
GET /v2/documents/{id}
```

Lấy thông tin chi tiết của một tài liệu.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, document ID
}
```

**Response:**
```json
{
  "document": {
    "id": "uint32",
    "groupId": "uint32",
    "name": "string",
    "description": "string",
    "fileUrl": "string",
    "fileType": "string",
    "fileSize": "uint64",
    "createdAt": "string",  // ISO 8601 format
    "updatedAt": "string"   // ISO 8601 format
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Document not found
- 500: Internal server error

### Get Group Documents

```http
GET /v2/group/{groupId}/documents
```

Lấy danh sách tài liệu của một nhóm. Kết quả được phân trang.

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
  "documents": [
    {
      "id": "uint32",
      "groupId": "uint32",
      "name": "string",
      "description": "string",
      "fileUrl": "string",
      "fileType": "string",
      "fileSize": "uint64",
      "createdAt": "string",  // ISO 8601 format
      "updatedAt": "string"   // ISO 8601 format
    }
  ],
  "total": "uint32"
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Group not found
- 500: Internal server error

## Data Models

### Group Document

```json
{
  "id": "uint32",
  "groupId": "uint32",
  "name": "string",         // max 255 chars
  "description": "string",
  "fileUrl": "string",      // valid URL
  "fileType": "string",     // valid file type
  "fileSize": "uint64",     // file size in bytes
  "createdAt": "string",    // ISO 8601 format
  "updatedAt": "string"     // ISO 8601 format
}
```

## Validation Rules

1. **Document Name**
   - Required
   - Maximum length: 255 characters

2. **Description**
   - Optional
   - Maximum length: 1000 characters

3. **File URL**
   - Required
   - Must be a valid URL format
   - Must be accessible

4. **File Type**
   - Required
   - Must be one of: pdf, doc, docx, xls, xlsx, ppt, pptx, txt, jpg, jpeg, png, gif
   - Maximum file size: 50MB

5. **File Size**
   - Required
   - Must be a positive number
   - Maximum size: 50MB (52428800 bytes)

6. **Group ID**
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

1. **Document Events**
   - DOCUMENT_CREATE: Tạo tài liệu mới
   - DOCUMENT_UPDATE: Cập nhật thông tin tài liệu
   - DOCUMENT_DELETE: Xóa tài liệu
   - DOCUMENT_DOWNLOAD: Tải xuống tài liệu 