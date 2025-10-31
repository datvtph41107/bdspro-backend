# User Guide API Documentation

API để quản lý user guides (hướng dẫn sử dụng) trong hệ thống.

## Endpoints

Base URL: `/v2/hub/user-guides`

### 1. Tạo User Guide Mới
**POST** `/v2/hub/user-guides`

**Request Body:**
```json
{
  "title": "Hướng dẫn sử dụng tính năng X",
  "description": "Đây là mô tả chi tiết về cách sử dụng tính năng X...",
  "groupKey": "feature-x"
}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Hướng dẫn sử dụng tính năng X",
    "description": "Đây là mô tả chi tiết về cách sử dụng tính năng X...",
    "groupKey": "feature-x",
    "createdAt": "2025-10-28T10:00:00Z",
    "updatedAt": "2025-10-28T10:00:00Z"
  }
}
```

---

### 2. Cập Nhật User Guide
**PUT** `/v2/hub/user-guides/{id}`

**Request Body:**
```json
{
  "id": 1,
  "title": "Hướng dẫn sử dụng tính năng X (đã cập nhật)",
  "description": "Mô tả đã được cập nhật...",
  "groupKey": "feature-x"
}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Hướng dẫn sử dụng tính năng X (đã cập nhật)",
    "description": "Mô tả đã được cập nhật...",
    "groupKey": "feature-x",
    "createdAt": "2025-10-28T10:00:00Z",
    "updatedAt": "2025-10-28T10:15:00Z"
  }
}
```

---

### 3. Xóa User Guide
**DELETE** `/v2/hub/user-guides/{id}`

**Response:**
```json
{
  "success": true
}
```

---

### 4. Lấy Chi Tiết User Guide
**GET** `/v2/hub/user-guides/{id}`

**Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Hướng dẫn sử dụng tính năng X",
    "description": "Đây là mô tả chi tiết về cách sử dụng tính năng X... (FULL DESCRIPTION - Không bị cắt)",
    "groupKey": "feature-x",
    "createdAt": "2025-10-28T10:00:00Z",
    "updatedAt": "2025-10-28T10:00:00Z"
  }
}
```

**Lưu ý:** API chi tiết trả về **FULL description** (toàn bộ nội dung).

---

### 5. Lấy Danh Sách User Guides (Có Phân Trang & Tìm Kiếm)
**GET** `/v2/hub/user-guides`

**Query Parameters:**
- `title` (string, optional): Tìm kiếm theo title (ILIKE)
- `groupKey` (string, optional): Lọc theo groupKey
- `page` (int, optional, default: 1): Số trang
- `size` (int, optional, default: 20, max: 100): Số bản ghi trên trang

**Example Requests:**

1. Lấy tất cả (trang 1, 20 bản ghi):
   ```
   GET /v2/hub/user-guides
   ```

2. Tìm kiếm theo title:
   ```
   GET /v2/hub/user-guides?title=tính năng
   ```

3. Lọc theo groupKey:
   ```
   GET /v2/hub/user-guides?groupKey=feature-x
   ```

4. Kết hợp tìm kiếm và phân trang:
   ```
   GET /v2/hub/user-guides?title=hướng dẫn&groupKey=feature-x&page=2&size=10
   ```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Hướng dẫn sử dụng tính năng X",
      "description": "Đây là mô tả chi tiết về cách sử dụng tính ...",
      "groupKey": "feature-x",
      "createdAt": "2025-10-28T10:00:00Z",
      "updatedAt": "2025-10-28T10:00:00Z"
    },
    {
      "id": 2,
      "title": "Hướng dẫn sử dụng tính năng Y",
      "description": "Mô tả khác về tính năng Y, giúp người dùng...",
      "groupKey": "feature-y",
      "createdAt": "2025-10-28T10:05:00Z",
      "updatedAt": "2025-10-28T10:05:00Z"
    }
  ],
  "total": 2
}
```

**Lưu ý:** Description trong danh sách **chỉ hiển thị 50 ký tự đầu tiên** và thêm "..." nếu vượt quá. Để xem full description, sử dụng API chi tiết (GET by ID).

---

### 6. Lấy Danh Sách Nhóm (Groups)
**GET** `/v2/hub/user-guides/groups`

Trả về danh sách các `groupKey` và số lượng user guide trong mỗi nhóm.

**Response:**
```json
{
  "data": [
    {
      "groupKey": "feature-x",
      "count": 5
    },
    {
      "groupKey": "feature-y",
      "count": 3
    },
    {
      "groupKey": "getting-started",
      "count": 10
    }
  ]
}
```

---

## Database Migration

Để tạo bảng `tb_user_guide`, chạy migration:

```sql
-- File: migrate/004_create_user_guide_table.sql
CREATE TABLE IF NOT EXISTS tb_user_guide (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    group_key VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT NULL,
    updated_by BIGINT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_guide_deleted_at ON tb_user_guide(deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_guide_group_key ON tb_user_guide(group_key);
CREATE INDEX IF NOT EXISTS idx_user_guide_title ON tb_user_guide(title);
CREATE INDEX IF NOT EXISTS idx_user_guide_created_at ON tb_user_guide(created_at DESC);
```

---

## Các Tính Năng Chính

1. **CRUD đầy đủ**: Create, Read, Update, Delete
2. **Tìm kiếm**: Tìm kiếm theo title (case-insensitive)
3. **Lọc theo nhóm**: Lọc user guides theo groupKey
4. **Phân trang**: Hỗ trợ phân trang với page và size
5. **Soft Delete**: Xóa mềm, không xóa vật lý
6. **Audit Trail**: Tự động lưu created_by, updated_by, created_at, updated_at
7. **Nhóm thống kê**: API lấy danh sách nhóm với số lượng
8. **Smart Truncation**: 
   - API danh sách: Description tự động cắt còn 50 ký tự + "..." (tối ưu hiệu suất)
   - API chi tiết: Trả về full description (đầy đủ nội dung)

---

## Ví Dụ Sử Dụng với cURL

### Tạo user guide mới:
```bash
curl -X POST http://localhost:8080/v2/hub/user-guides \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Cách tạo sản phẩm mới",
    "description": "Bước 1: Vào menu Sản phẩm...",
    "groupKey": "product-management"
  }'
```

### Lấy danh sách với tìm kiếm:
```bash
curl -X GET "http://localhost:8080/v2/hub/user-guides?title=sản phẩm&page=1&size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Lấy danh sách nhóm:
```bash
curl -X GET http://localhost:8080/v2/hub/user-guides/groups \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Cập nhật user guide:
```bash
curl -X PUT http://localhost:8080/v2/hub/user-guides/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "id": 1,
    "title": "Cách tạo sản phẩm mới (Cập nhật)",
    "description": "Hướng dẫn chi tiết...",
    "groupKey": "product-management"
  }'
```

### Xóa user guide:
```bash
curl -X DELETE http://localhost:8080/v2/hub/user-guides/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Gợi Ý GroupKey

Các `groupKey` có thể sử dụng:

- `getting-started`: Hướng dẫn bắt đầu
- `product-management`: Quản lý sản phẩm
- `customer-management`: Quản lý khách hàng
- `order-processing`: Xử lý đơn hàng
- `reporting`: Báo cáo
- `settings`: Cài đặt hệ thống
- `troubleshooting`: Khắc phục sự cố
- `advanced-features`: Tính năng nâng cao

---

## Notes

- API yêu cầu JWT authentication (header `Authorization: Bearer <token>`)
- Tất cả response đều ở dạng JSON
- Xóa user guide là soft delete (deleted_at != null)
- Pagination mặc định: page=1, size=20, max size=100
- Tìm kiếm title không phân biệt hoa thường (case-insensitive)
- **Description truncation**:
  - **List API** (`GET /v2/hub/user-guides`): Description được cắt ở **50 ký tự** + "..." để tối ưu performance
  - **Detail API** (`GET /v2/hub/user-guides/{id}`): Trả về **full description** không bị cắt
  - Logic cắt chuỗi tính theo **rune** (hỗ trợ Unicode/tiếng Việt đúng)

