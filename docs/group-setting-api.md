# Group Setting API

## Get Group Settings

```http
GET /v2/group-setting
```

Lấy danh sách cài đặt của một nhóm.

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
  "settings": [
    {
      "id": "uint32",
      "configKey": "string",
      "configValue": "string"
    }
  ],
  "total": "uint32"  // Tổng số cài đặt
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - User is not a member of the group
- 404: Group not found
- 500: Internal server error

## Get Group Setting

```http
GET /v2/group-setting/{id}
```

Lấy thông tin chi tiết của một cài đặt.

**Request Parameters:**
```json
{
  "id": "uint32"  // required
}
```

**Response:**
```json
{
  "setting": {
    "id": "uint32",
    "configKey": "string",
    "configValue": "string"
  }
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - User is not a member of the group
- 404: Setting not found
- 500: Internal server error

## Create Group Setting

```http
POST /v2/group-setting
```

Tạo mới một cài đặt cho nhóm.

**Request Body:**
```json
{
  "groupId": "uint32",    // required
  "configKey": "string",  // required
  "configValue": "string" // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of created setting
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - User is not a leader of the group
- 404: Group not found
- 500: Internal server error

## Update Group Setting

```http
PUT /v2/group-setting/{id}
```

Cập nhật thông tin của một cài đặt.

**Request Body:**
```json
{
  "id": "uint32",        // required
  "configKey": "string", // required
  "configValue": "string" // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of updated setting
}
```

**Error Codes:**
- 400: Invalid request body
- 403: Forbidden - User is not a leader of the group
- 404: Setting not found
- 500: Internal server error

## Delete Group Setting

```http
DELETE /v2/group-setting/{id}
```

Xóa một cài đặt.

**Request Parameters:**
```json
{
  "id": "uint32"  // required
}
```

**Response:**
```json
{
  "id": "uint32"  // ID of deleted setting
}
```

**Error Codes:**
- 400: Invalid request parameters
- 403: Forbidden - User is not a leader of the group
- 404: Setting not found
- 500: Internal server error 