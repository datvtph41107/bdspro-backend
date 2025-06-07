# Group API Documentation

## Overview

Group API cung cấp các endpoint để quản lý thông tin cơ bản của các nhóm, bao gồm tạo mới, cập nhật, xóa và truy vấn thông tin nhóm.

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

### Create Group

```http
POST /v2/group
```

Tạo mới một nhóm.

**Request Body:**
```json
{
  "name": "string",           // required, max 255 chars
  "description": "string",    // optional
  "avatarUrl": "string"       // optional
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created group
}
```

**Error Codes:**
- 400: Invalid request body
- 500: Internal server error

### Update Group

```http
PUT /v2/group/{id}
```

Cập nhật thông tin của một nhóm.

**Request Body:**
```json
{
  "id": "uint32",           // required
  "name": "string",         // required, max 255 chars
  "description": "string"   // optional
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated group
}
```

**Error Codes:**
- 400: Invalid request body
- 404: Group not found
- 500: Internal server error

### Delete Group

```http
DELETE /v2/group/{id}
```

Xóa một nhóm.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, group ID
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of deleted group
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Group not found
- 500: Internal server error

### Get Group

```http
GET /v2/group/{id}
```

Lấy thông tin chi tiết của một nhóm.

**Request Parameters:**
```json
{
  "id": "uint32"  // required, group ID
}
```

**Response:**
```json
{
  "group": {
    "id": "uint32",
    "name": "string",
    "description": "string",
    "avatarUrl": "string"
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 404: Group not found
- 500: Internal server error

### Get Groups

```http
GET /v2/group
```

Lấy danh sách các nhóm. Kết quả được phân trang.

**Request Parameters:**
```json
{
  "page": "uint32",    // optional, default: 1
  "size": "uint32"     // optional, default: 20, max: 100
}
```

**Response:**
```json
{
  "groups": [
    {
      "id": "uint32",
      "name": "string",
      "description": "string",
      "avatarUrl": "string"
    }
  ],
  "total": "uint32"
}
```

**Error Codes:**
- 500: Internal server error

## Data Models

### Group

```json
{
  "id": "uint32",
  "name": "string",           // required, max 255 chars
  "description": "string",    // optional
  "avatarUrl": "string"       // optional
}
```

## Validation Rules

1. **Group Name**
   - Required
   - Maximum length: 255 characters

2. **Description**
   - Optional
   - Maximum length: 1000 characters

3. **Avatar URL**
   - Optional
   - Must be a valid URL format
   - Must be an image file (jpg, png, gif)

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

1. **Group Events**
   - GROUP_CREATE: Tạo mới nhóm
   - GROUP_UPDATE: Cập nhật thông tin nhóm
   - GROUP_DELETE: Xóa nhóm

## Group Member Management

### Create Member

```http
POST /v2/group/member
```

Thêm một thành viên mới vào nhóm.

**Request Body:**
```json
{
  "groupId": "uint32",  // required
  "userId": "uint32",   // required
  "role": "uint32"      // required, role ID
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
- 403: Forbidden - not authorized
- 404: Group not found
- 409: Member already exists
- 500: Internal server error

### Update Member

```http
PUT /v2/group/member
```

Cập nhật vai trò của thành viên trong nhóm.

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
- 403: Forbidden - not authorized
- 404: Member not found
- 500: Internal server error

### Delete Member

```http
DELETE /v2/group/member
```

Xóa một thành viên khỏi nhóm.

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
- 403: Forbidden - not authorized
- 404: Member not found
- 500: Internal server error

## Group Document Management

### 1. Create Group Document
**Endpoint:** `POST /v2/group-document`

#### Request
```json
{
    "groupId": "number",       // Required
    "name": "string",          // Required
    "description": "string",   // Optional
    "fileUrl": "string",       // Required
    "type": "string",          // Required
    "status": "string"         // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 2. Update Group Document
**Endpoint:** `PUT /v2/group-document/{id}`

#### Request
```json
{
    "id": "number",           // Required
    "name": "string",         // Required
    "description": "string",  // Optional
    "fileUrl": "string",      // Required
    "type": "string",         // Required
    "status": "string"        // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 3. Delete Group Document
**Endpoint:** `DELETE /v2/group-document/{id}`

#### Request
```json
{
    "id": "number"  // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 4. Get Group Documents
**Endpoint:** `GET /v2/group-document/{groupId}`

#### Response
```json
{
    "documents": [
        {
            "id": "number",
            "groupId": "number",
            "name": "string",
            "description": "string",
            "fileUrl": "string",
            "type": "string",
            "status": "string",
            "createdAt": "string",
            "updatedAt": "string"
        }
    ]
}
```

## Group Deal Management

### 1. Create Group Deal
**Endpoint:** `POST /v2/group-deal`

#### Request
```json
{
    "groupId": "number",       // Required
    "name": "string",          // Required
    "description": "string",   // Optional
    "value": "number",         // Required
    "status": "string",        // Required
    "startDate": "string",     // Required
    "endDate": "string"        // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 2. Update Group Deal
**Endpoint:** `PUT /v2/group-deal/{id}`

#### Request
```json
{
    "id": "number",           // Required
    "name": "string",         // Required
    "description": "string",  // Optional
    "value": "number",        // Required
    "status": "string",       // Required
    "startDate": "string",    // Required
    "endDate": "string"       // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 3. Delete Group Deal
**Endpoint:** `DELETE /v2/group-deal/{id}`

#### Request
```json
{
    "id": "number"  // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 4. Get Group Deals
**Endpoint:** `GET /v2/group-deal/{groupId}`

#### Response
```json
{
    "deals": [
        {
            "id": "number",
            "groupId": "number",
            "name": "string",
            "description": "string",
            "value": "number",
            "status": "string",
            "startDate": "string",
            "endDate": "string",
            "createdAt": "string",
            "updatedAt": "string"
        }
    ]
}
```

## Group Settings Management

### 1. Update Group Settings
**Endpoint:** `PUT /v2/group-settings/{groupId}`

#### Request
```json
{
    "groupId": "number",           // Required
    "notificationEnabled": "boolean", // Required
    "autoApproveMembers": "boolean",  // Required
    "allowMemberInvite": "boolean",   // Required
    "privacyLevel": "string"          // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 2. Get Group Settings
**Endpoint:** `GET /v2/group-settings/{groupId}`

#### Response
```json
{
    "groupId": "number",
    "notificationEnabled": "boolean",
    "autoApproveMembers": "boolean",
    "allowMemberInvite": "boolean",
    "privacyLevel": "string",
    "updatedAt": "string"
}
```

## Group Notification Management

### 1. Create Group Notification
**Endpoint:** `POST /v2/group-notification`

#### Request
```json
{
    "groupId": "number",       // Required
    "title": "string",         // Required
    "content": "string",       // Required
    "type": "string",          // Required
    "priority": "string"       // Required
}
```

#### Response
```json
{
    "id": "number"
}
```

#### Authorization
- Requires valid JWT token
- User must be a member of the group
- User must have appropriate permissions

### 2. Get Group Notifications
**Endpoint:** `GET /v2/group-notification/{groupId}`

#### Response
```json
{
    "notifications": [
        {
            "id": "number",
            "groupId": "number",
            "title": "string",
            "content": "string",
            "type": "string",
            "priority": "string",
            "createdAt": "string"
        }
    ]
}
```

## Rate Limiting

- Maximum 100 requests per minute per IP
- Maximum 1000 requests per hour per IP
- Maximum 50 concurrent requests per user
- Maximum 10 requests per second per endpoint

## Validation Rules

1. **Group Name**
   - Required
   - Maximum length: 255 characters
   - Cannot contain special characters except spaces and hyphens

2. **Group Description**
   - Optional
   - Maximum length: 1000 characters
   - Can contain any characters

3. **Group Type**
   - Required
   - Must be one of: "department", "project", "team", "other"
   - Case sensitive

4. **Group Status**
   - Required
   - Must be one of: "active", "inactive", "archived"
   - Default: "active"

5. **Member Limit**
   - Maximum 100 members per group
   - Cannot exceed organization's member limit

6. **Role Name**
   - Required
   - Maximum length: 50 characters
   - Must be unique within the group
   - Cannot contain special characters

7. **Permission Key**
   - Required
   - Maximum length: 50 characters
   - Must be unique within the group
   - Format: "resource:action" (e.g., "group:create", "member:delete")

## Error Codes

- 400: Bad Request - Invalid input data
  - Invalid group name format
  - Invalid group type
  - Invalid group status
  - Member limit exceeded
  - Invalid role name format
  - Invalid permission key format

- 401: Unauthorized - Missing or invalid authentication
  - Missing JWT token
  - Invalid JWT token
  - Expired JWT token

- 403: Forbidden - Insufficient permissions
  - Not a group admin
  - Not a group member
  - Insufficient role permissions
  - Cannot modify admin role
  - Cannot remove last admin

- 404: Not Found - Resource not found
  - Group not found
  - Member not found
  - Role not found
  - Permission not found
  - Branch not found

- 409: Conflict - Resource already exists
  - Group name already exists
  - Role name already exists
  - Permission key already exists
  - Member already exists
  - Branch name already exists

- 429: Too Many Requests - Rate limit exceeded
  - Too many requests per minute
  - Too many requests per hour
  - Too many concurrent requests

- 500: Internal Server Error - Server-side error
  - Database error
  - Cache error
  - External service error
  - Unexpected error

## Common Error Response Format

```json
{
  "code": "number",
  "message": "string",
  "details": [
    {
      "field": "string",
      "description": "string"
    }
  ]
}
```

## Notes

1. All timestamps are in ISO 8601 format (UTC)
2. All IDs are uint32
3. All status fields are case sensitive
4. All role and permission keys are case sensitive
5. Maximum file size for group avatar: 5MB
6. Supported image formats: JPG, PNG, GIF
7. Maximum number of roles per group: 10
8. Maximum number of permissions per role: 20
9. Maximum number of branches per group: 5
10. Group deletion is soft delete (archived)

## Logging

Service logs các sự kiện sau:

1. **Group Events**
   - GROUP_CREATE: Tạo mới nhóm
   - GROUP_UPDATE: Cập nhật thông tin nhóm
   - GROUP_DELETE: Xóa nhóm

2. **Member Events**
   - MEMBER_ADD: Thêm thành viên vào nhóm
   - MEMBER_REMOVE: Xóa thành viên khỏi nhóm
   - MEMBER_ROLE_CHANGE: Thay đổi vai trò thành viên
   - MEMBER_STATUS_CHANGE: Thay đổi trạng thái thành viên

3. **Role Events**
   - ROLE_CREATE: Tạo vai trò mới
   - ROLE_UPDATE: Cập nhật vai trò
   - ROLE_DELETE: Xóa vai trò
   - ROLE_PERMISSION_CHANGE: Thay đổi quyền của vai trò

4. **Permission Events**
   - PERMISSION_CREATE: Tạo quyền mới
   - PERMISSION_UPDATE: Cập nhật quyền
   - PERMISSION_DELETE: Xóa quyền

5. **Branch Events**
   - BRANCH_CREATE: Tạo chi nhánh mới
   - BRANCH_UPDATE: Cập nhật thông tin chi nhánh
   - BRANCH_DELETE: Xóa chi nhánh
   - BRANCH_MANAGER_CHANGE: Thay đổi quản lý chi nhánh

Mỗi log event bao gồm:
- Timestamp
- Event type
- User ID thực hiện
- Resource ID bị ảnh hưởng
- Thông tin chi tiết về thay đổi
- Organization ID (nếu có)
- Group ID (nếu có)
- Branch ID (nếu có) 