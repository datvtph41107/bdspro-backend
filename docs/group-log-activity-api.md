# Group Log Activity API Documentation

## Overview

Group Log Activity API cung cấp các endpoint để truy vấn lịch sử hoạt động của nhóm, bao gồm các sự kiện như tạo nhóm, cập nhật thông tin, thêm/xóa thành viên, và các hoạt động khác.

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

### Get Group Log Activities

```http
GET /v2/group/log-activity
```

Lấy danh sách các hoạt động của một nhóm.

**Request Parameters:**
```json
{
  "groupId": "uint32",    // required
  "page": "uint32",       // optional, default: 0
  "size": "uint32"        // optional, default: 10
}
```

**Response:**
```json
{
  "logActivities": [
    {
      "id": "uint32",
      "groupId": "uint32",
      "actorId": "uint32",
      "logType": "string",
      "logData": "string",
      "createdAt": "string",  // ISO 8601 format
      "updatedAt": "string"   // ISO 8601 format
    }
  ],
  "total": "uint32"  // Tổng số bản ghi
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Group not found
- 500: Internal server error

## Data Models

### Group Log Activity

```json
{
  "id": "uint32",
  "groupId": "uint32",
  "actorId": "uint32",    // ID của người thực hiện hành động
  "logType": "string",    // Loại hoạt động
  "logData": "string",    // Dữ liệu chi tiết của hoạt động
  "createdAt": "string",  // ISO 8601 format
  "updatedAt": "string"   // ISO 8601 format
}
```

## Log Types

1. **Group Events**
   - GROUP_CREATE: Tạo nhóm mới
   - GROUP_UPDATE: Cập nhật thông tin nhóm
   - GROUP_DELETE: Xóa nhóm

2. **Member Events**
   - MEMBER_ADD: Thêm thành viên
   - MEMBER_REMOVE: Xóa thành viên
   - MEMBER_UPDATE: Cập nhật vai trò thành viên

3. **Document Events**
   - DOCUMENT_CREATE: Tạo tài liệu mới
   - DOCUMENT_UPDATE: Cập nhật tài liệu
   - DOCUMENT_DELETE: Xóa tài liệu
   - DOCUMENT_DOWNLOAD: Tải xuống tài liệu

4. **Deal Events**
   - DEAL_CREATE: Tạo giao dịch mới
   - DEAL_UPDATE: Cập nhật giao dịch
   - DEAL_DELETE: Xóa giao dịch
   - DEAL_STATUS_CHANGE: Thay đổi trạng thái giao dịch

5. **Setting Events**
   - SETTING_CREATE: Tạo cài đặt mới
   - SETTING_UPDATE: Cập nhật cài đặt
   - SETTING_DELETE: Xóa cài đặt

## Validation Rules

1. **Group ID**
   - Required
   - Must be a valid group ID

2. **Pagination**
   - Page: Optional, default: 1
   - Size: Optional, default: 20, max: 100

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

- Logs được lưu trữ trong 30 ngày
- Có thể xuất logs trong khoảng thời gian này
- Sau 30 ngày, logs sẽ được tự động xóa 