# API User Guide - Danh sách Newsfeed và Product của User

## Tổng quan
Đã bổ sung đầy đủ API để lấy danh sách newsfeed (bài viết) và product (sản phẩm) của một user cụ thể.

---

## 1. API Newsfeed (Bài viết)

### 1.1. API Newsfeed của User - Cần Authentication
**Endpoint**: `GET /v2/social/news-feed/user/{userId}`

**Service**: `social-service`

**Description**: Lấy danh sách bài viết của một user cụ thể (yêu cầu authentication)

**Authentication**: Bearer Token required

**Parameters**:
- **Path Parameters**:
  - `userId` (uint64, required): ID của user cần lấy bài viết

- **Query Parameters**:
  - `page` (int32, optional): Số trang (mặc định: 1)
  - `size` (int32, optional): Số bài viết mỗi trang (mặc định: 10)
  - `text` (string, optional): Tìm kiếm theo nội dung
  - `isReel` (bool, optional): Lọc theo reel

**Response**:
```json
{
  "data": [
    {
      "id": 123,
      "title": "Tiêu đề bài viết",
      "content": "Nội dung bài viết",
      "image": "url_image",
      "visibility": 10,
      "numLike": 100,
      "numComment": 50,
      "numShare": 20,
      "numView": 1000,
      "updatedAt": "2025-10-21T10:00:00Z",
      "author": {
        "id": 456,
        "fullName": "Tên user",
        "avatar": "url_avatar"
      },
      "mediaFeeds": [...],
      "friendTags": [...],
      "post": {...}
    }
  ],
  "total": 100
}
```

**Example**:
```bash
# Lấy danh sách bài viết của user 123 (cần token)
curl -X GET "http://api.example.com/v2/social/news-feed/user/123?page=1&size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Lấy danh sách reel của user 123
curl -X GET "http://api.example.com/v2/social/news-feed/user/123?page=1&size=10&isReel=true" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Tìm kiếm bài viết của user 123 theo nội dung
curl -X GET "http://api.example.com/v2/social/news-feed/user/123?page=1&size=10&text=bất động sản" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 1.2. API Newsfeed của User - Public (No Authentication) 🆕
**Endpoint**: `GET /v2/social/public/news-feed/user/{userId}`

**Service**: `social-service`

**Description**: Lấy danh sách bài viết công khai của một user cụ thể - API public không cần authentication

**Authentication**: ❌ Không yêu cầu

**Parameters**:
- **Path Parameters**:
  - `userId` (uint64, required): ID của user cần lấy bài viết

- **Query Parameters**:
  - `page` (int32, optional): Số trang (mặc định: 1)
  - `size` (int32, optional): Số bài viết mỗi trang (mặc định: 10)
  - `text` (string, optional): Tìm kiếm theo nội dung
  - `isReel` (bool, optional): Lọc theo reel
  - `isPost` (bool, optional): Lọc theo bài viết thường

**Response**: Giống như API 1.1

**Example**:
```bash
# Lấy bài viết công khai của user 123 (không cần token)
curl -X GET "http://api.example.com/v2/social/public/news-feed/user/123?page=1&size=10"

# Lấy danh sách reel công khai
curl -X GET "http://api.example.com/v2/social/public/news-feed/user/123?page=1&size=10&isReel=true"

# Tìm kiếm bài viết theo nội dung
curl -X GET "http://api.example.com/v2/social/public/news-feed/user/123?page=1&size=10&text=nhà đẹp"
```

---

## 2. API Product (Sản phẩm) 🆕

### 2.1. API Product của User - Cần Authentication
**Endpoint**: `GET /v2/bdspro/v2/product/user/{userId}`

**Service**: `bdspro-service`

**Description**: Lấy danh sách sản phẩm của một user cụ thể (yêu cầu authentication)

**Authentication**: Bearer Token required

**Parameters**:
- **Path Parameters**:
  - `userId` (uint64, required): ID của user cần lấy sản phẩm

- **Query Parameters**:
  - `page` (uint32, optional): Số trang (mặc định: 1)
  - `size` (uint32, optional): Số sản phẩm mỗi trang (mặc định: 10, max: 100)
  - `text` (string, optional): Tìm kiếm theo tên/mô tả
  - `parentId` (uint64, optional): ID sản phẩm cha (null: all, 0: parent, >0: children)
  - `propertyTypeId` (uint64, optional): Loại tài sản
  - `docTypeId` (uint64, optional): Loại giấy tờ pháp lý
  - `projectId` (uint64, optional): Dự án
  - `provinceId` (uint64, optional): Tỉnh/Thành phố
  - `districtId` (uint64, optional): Quận/Huyện
  - `wardId` (uint64, optional): Phường/Xã
  - `transactionType` (int32, optional): Loại giao dịch (10: mua bán, 20: cho thuê)
  - `saleStatus` (int32, optional): Trạng thái bán
  - `rentStatus` (int32, optional): Trạng thái cho thuê
  - `saleVisibility` (int32, optional): Hiển thị bán
  - `rentVisibility` (int32, optional): Hiển thị cho thuê
  - `sourceType` (int32, optional): Nguồn sản phẩm
  - `archived` (bool, optional): Đã lưu trữ

**Response**:
```json
{
  "data": [
    {
      "id": 789,
      "name": "Nhà đẹp 3 tầng",
      "code": "SP000123",
      "area": 120.5,
      "description": "Mô tả sản phẩm...",
      "address": "123 Đường ABC",
      "transactionType": 10,
      "saleStatus": 10,
      "rentStatus": 0,
      "propertyType": {...},
      "docType": {...},
      "province": {...},
      "district": {...},
      "ward": {...},
      "price": {...},
      "mediaList": [...],
      "houseInfo": {...}
    }
  ],
  "totalElements": 50
}
```

**Example**:
```bash
# Lấy tất cả sản phẩm của user 456
curl -X GET "http://api.example.com/v2/bdspro/v2/product/user/456?page=1&size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Lọc sản phẩm bán tại Hà Nội
curl -X GET "http://api.example.com/v2/bdspro/v2/product/user/456?transactionType=10&provinceId=1" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Tìm kiếm sản phẩm theo từ khóa
curl -X GET "http://api.example.com/v2/bdspro/v2/product/user/456?text=nhà đẹp&page=1&size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 2.2. API Product của User - Public (No Authentication) 🆕
**Endpoint**: `GET /v2/bdspro/public/product/user/{profileId}`

**Service**: `bdspro-service`

**Description**: Lấy danh sách sản phẩm của một user cụ thể - API public không cần authentication

**Authentication**: ❌ Không yêu cầu

**Parameters**:
- **Path Parameters**:
  - `profileId` (uint64, required): ID của profile user cần lấy sản phẩm

- **Query Parameters**: (Giống như API 2.1)
  - `page`, `size`, `text`, `parentId`, `propertyTypeId`, `docTypeId`, `projectId`
  - `provinceId`, `districtId`, `wardId`, `transactionType`, `saleStatus`, `rentStatus`
  - `saleVisibility`, `rentVisibility`, `sourceType`, `archived`

**Response**: Giống như API 2.1

**Example**:
```bash
# Lấy sản phẩm công khai của user 456 (không cần token)
curl -X GET "http://api.example.com/v2/bdspro/public/product/user/456?page=1&size=20"

# Lọc sản phẩm cho thuê
curl -X GET "http://api.example.com/v2/bdspro/public/product/user/456?transactionType=20&page=1"

# Tìm kiếm sản phẩm theo vị trí
curl -X GET "http://api.example.com/v2/bdspro/public/product/user/456?provinceId=1&districtId=5"
```

---

## 3. So sánh API Private vs Public

### 3.1. Newsfeed APIs

| Feature | API Private (`/v2/social/news-feed/user/{userId}`) | API Public (`/v2/social/public/news-feed/user/{userId}`) |
|---------|---------------------------------------------------|----------------------------------------------------------|
| **Authentication** | ✅ Required (Bearer Token) | ❌ Not required |
| **Use Case** | Dùng trong app khi user đã đăng nhập | Dùng cho public profile, chia sẻ link |
| **Parameters** | Full filter options | Full filter options |
| **Data Access** | Có thể xem thêm thông tin riêng tư | Chỉ xem bài viết public |
| **Performance** | Standard | Standard |

### 3.2. Product APIs

| Feature | API Private (`/v2/bdspro/v2/product/user/{userId}`) | API Public (`/v2/bdspro/public/product/user/{profileId}`) |
|---------|-----------------------------------------------------|-----------------------------------------------------------|
| **Authentication** | ✅ Required (Bearer Token) | ❌ Not required |
| **Use Case** | Dùng trong app khi user đã đăng nhập | Dùng cho public profile, chia sẻ link |
| **Parameters** | Đầy đủ filter options | Đầy đủ filter options |
| **Data Access** | Có thể xem thêm thông tin riêng tư (nếu có permission) | Chỉ xem thông tin public |
| **Performance** | Standard | Standard |

---

## 4. Kết hợp API trong ứng dụng

### 4.1. Use Case: Xem profile của user khác (Public Profile)
```javascript
// 1. Lấy thông tin newsfeed công khai (không cần token)
const newsfeeds = await fetch(`/v2/social/public/news-feed/user/${userId}?page=1&size=10`);

// 2. Lấy danh sách sản phẩm công khai (không cần token)
const products = await fetch(`/v2/bdspro/public/product/user/${userId}?page=1&size=10`);

// 3. Hiển thị trên profile page
displayUserProfile(newsfeeds.data, products.data);
```

### 4.2. Use Case: Admin xem chi tiết user
```javascript
// Sử dụng API private với token admin
const newsfeeds = await fetch(`/v2/social/news-feed/user/${userId}?page=1&size=20`, {
  headers: {
    'Authorization': `Bearer ${adminToken}`
  }
});

const products = await fetch(`/v2/bdspro/v2/product/user/${userId}?page=1&size=20`, {
  headers: {
    'Authorization': `Bearer ${adminToken}`
  }
});
```

---

## 5. Error Handling

### Common Error Codes:

#### 400 Bad Request
```json
{
  "code": 400,
  "message": "Invalid request parameters"
}
```

#### 401 Unauthorized (Private API only)
```json
{
  "code": 401,
  "message": "Unauthorized - Missing or invalid token"
}
```

#### 404 Not Found
```json
{
  "code": 404,
  "message": "User not found"
}
```

#### 500 Internal Server Error
```json
{
  "code": 500,
  "message": "Internal server error"
}
```

---

## 6. Pagination Best Practices

### Default Values:
- **Default Page**: 1
- **Default Size**: 10
- **Max Size**: 100 (newsfeed), 100 (product)

### Example Response với Pagination:
```json
{
  "data": [...],
  "total": 250,
  "page": 1,
  "size": 20
}
```

### Tính tổng số trang:
```javascript
const totalPages = Math.ceil(response.total / response.size);
```

---

## 7. Workflow đầy đủ

### Frontend Flow:
```
1. User truy cập profile page: /profile/{userId}
   ↓
2. Gọi API public lấy newsfeed (không cần token)
   GET /v2/social/public/news-feed/user/{userId}
   ↓
3. Gọi API public lấy products (không cần token)
   GET /v2/bdspro/public/product/user/{userId}
   ↓
4. Render UI với data
```

### Backend Implementation:
```
1. Proto định nghĩa:
   ✅ social/news_feed.proto → GetNewsFeedByUserID, GetUserNewsFeedPublic
   ✅ bdspro/product.proto → GetUserProducts
   ✅ bdspro/public.proto → GetUserProductsPublic

2. Handler layer:
   ✅ social-service/infra/service/news_feed.go → GetNewsFeedByUserID, GetUserNewsFeedPublic
   ✅ bdspro-service/infra/handler/product_handler.go → GetUserProducts
   ✅ bdspro-service/infra/handler/public_handler.go → GetUserProductsPublic

3. Usecase layer:
   ✅ social-service/internal/usecase/news_feed.go
   ✅ bdspro-service/internal/usecases/shared/product_usecase.go

4. Repository layer:
   ✅ social-service/infra/postgre/news_feed_postgre.go
   ✅ bdspro-service/infra/postgres/product_postgre.go
```

---

## 8. Testing

### Test Newsfeed Private API:
```bash
# Cần token
curl -X GET "http://localhost:8003/v2/social/news-feed/user/123?page=1&size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test Newsfeed Public API:
```bash
# Không cần token
curl -X GET "http://localhost:8003/v2/social/public/news-feed/user/123?page=1&size=10"
```

### Test Product Private API:
```bash
# Cần token
curl -X GET "http://localhost:8002/v2/bdspro/v2/product/user/123?page=1&size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test Product Public API:
```bash
# Không cần token
curl -X GET "http://localhost:8002/v2/bdspro/public/product/user/123?page=1&size=10"
```

---

## 9. Tóm tắt

| Feature | Newsfeed (Private) | Newsfeed (Public) | Product (Private) | Product (Public) |
|---------|-------------------|-------------------|-------------------|------------------|
| **Endpoint** | `/v2/social/news-feed/user/{userId}` | `/v2/social/public/news-feed/user/{userId}` | `/v2/bdspro/v2/product/user/{userId}` | `/v2/bdspro/public/product/user/{profileId}` |
| **Service** | social-service | social-service | bdspro-service | bdspro-service |
| **Authentication** | Required | Not required | Required | Not required |
| **Pagination** | ✅ | ✅ | ✅ | ✅ |
| **Filter** | Full | Full | Full | Full |
| **Use Case** | Admin/Internal | Public profile | Admin/Internal | Public profile |

---

## 10. Notes

### APIs đã có:
- ✅ **Newsfeed Private API**: `/v2/social/news-feed/user/{userId}` - Cần token
- 🆕 **Newsfeed Public API**: `/v2/social/public/news-feed/user/{userId}` - **Không cần token** ⭐
- 🆕 **Product Private API**: `/v2/bdspro/v2/product/user/{userId}` - Cần token
- 🆕 **Product Public API**: `/v2/bdspro/public/product/user/{profileId}` - **Không cần token** ⭐

### Đặc điểm chung:
Tất cả API đều hỗ trợ:
- ✅ Phân trang (page, size)
- ✅ Tìm kiếm (text search)
- ✅ Filter đa dạng
- ✅ Response format chuẩn
- ✅ Error handling

### Khuyến nghị sử dụng:
- **Public APIs** (`/v2/.../public/...`): Dùng cho public profile, share link, không cần authentication
- **Private APIs**: Dùng cho trang admin, internal tool, yêu cầu authentication

---

**Ngày cập nhật**: 2025-10-21  
**Version**: 2.0.0

