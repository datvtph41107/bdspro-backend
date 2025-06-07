# Group Notification API Documentation

## Overview

Group Notification API cung cấp các endpoint để quản lý thông báo trong nhóm, bao gồm xem danh sách thông báo và chi tiết thông báo.

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

### Get Group Notifications

```http
GET /v2/group-notification
```

Lấy danh sách thông báo của một nhóm. Kết quả được phân trang.

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
  "notifications": [
    {
      "id": "uint32",
      "groupId": "uint32",
      "type": "string",      // notification type
      "content": "string",   // notification content
      "createdAt": "string", // ISO 8601 format
      "updatedAt": "string"  // ISO 8601 format
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

### Group Notification

```json
{
  "id": "uint32",
  "groupId": "uint32",
  "type": "string",      // notification type
  "content": "string",   // notification content
  "createdAt": "string", // ISO 8601 format
  "updatedAt": "string"  // ISO 8601 format
}
```

## Notification Types

1. **Member Events**
   - MEMBER_JOINED: Thành viên mới tham gia nhóm
   - MEMBER_LEFT: Thành viên rời khỏi nhóm
   - MEMBER_ROLE_CHANGED: Thay đổi vai trò thành viên

2. **Group Events**
   - GROUP_UPDATED: Cập nhật thông tin nhóm
   - GROUP_SETTING_CHANGED: Thay đổi cài đặt nhóm

3. **Document Events**
   - DOCUMENT_ADDED: Thêm tài liệu mới
   - DOCUMENT_UPDATED: Cập nhật tài liệu
   - DOCUMENT_DELETED: Xóa tài liệu

4. **Deal Events**
   - DEAL_CREATED: Tạo giao dịch mới
   - DEAL_UPDATED: Cập nhật giao dịch
   - DEAL_STATUS_CHANGED: Thay đổi trạng thái giao dịch

## Validation Rules

1. **Group ID**
   - Required
   - Must be a valid group ID

2. **Pagination**
   - page: số nguyên dương, mặc định: 1
   - size: số nguyên dương, mặc định: 20, tối đa: 100

## Error Codes

- 400: Bad Request - Invalid input data
- 401: Unauthorized - Missing or invalid authentication
- 403: Forbidden - Insufficient permissions
- 404: Not Found - Resource not found
- 500: Internal Server Error - Server-side error

## Rate Limiting

- Maximum 100 requests per minute per IP
- Maximum 1000 requests per hour per IP

## Notification Retention

- Thông báo được lưu trữ trong 30 ngày
- Thông báo cũ hơn 30 ngày sẽ tự động bị xóa
- Có thể export thông báo trong khoảng thời gian 30 ngày 