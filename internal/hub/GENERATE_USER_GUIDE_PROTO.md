# Generate Proto Code for User Guide Simple API

## ✅ Đã hoàn thành

1. ✅ Proto definitions (`shared/protobuf/schema/hub/user_guide.proto`)
2. ✅ Usecase method (`hub-service/internal/usecase/user_guide_usecase.go`)
3. ✅ Handler method (`hub-service/infra/handler/user_guide_handler.go`)
4. ✅ Documentation (`hub-service/USER_GUIDE_SIMPLE_API.md`)

## 🔧 Cần generate proto code

### 1. Generate Go code từ proto

```bash
cd shared/protobuf/schema/hub

# Generate Go protobuf + gRPC code
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  user_guide.proto
```

### 2. Generate gRPC Gateway (nếu cần REST API)

```bash
cd shared/protobuf/schema/hub

protoc --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  user_guide.proto
```

### 3. Generate Swagger documentation (optional)

```bash
cd shared/protobuf/schema/hub

protoc --openapiv2_out=../../docs/hub \
  --openapiv2_opt=logtostderr=true \
  user_guide.proto
```

### 4. Compile và test

```bash
cd hub-service

# Build service
go build -o hub-service cmd/grpc/main.go

# Run service
./hub-service
```

## 📋 API Endpoint mới

```
GET /v2/hub/user-guides/simple
```

### Request Examples
```bash
# Lấy tất cả user guides
curl http://localhost:8080/v2/hub/user-guides/simple

# Lọc theo groupKey
curl http://localhost:8080/v2/hub/user-guides/simple?groupKey=post

# Phân trang
curl http://localhost:8080/v2/hub/user-guides/simple?page=1&size=10
```

### Response Example
```json
{
  "data": [
    {
      "id": 1,
      "title": "Hướng dẫn tạo bài đăng",
      "description": "Chi tiết cách tạo bài đăng mới với đầy đủ thông tin về tiêu đề, nội dung..."
    },
    {
      "id": 2,
      "title": "Cách upload hình ảnh",
      "description": "Hướng dẫn chi tiết cách upload và quản lý hình ảnh cho bài đăng của bạn..."
    }
  ],
  "total": 2
}
```

## 🎯 Use Cases

### Frontend Usage
```javascript
// Lấy user guides đơn giản
const response = await fetch('/v2/hub/user-guides/simple?groupKey=post');
const { data } = await response.json();

// Hiển thị danh sách
data.forEach(guide => {
  console.log(`${guide.title}: ${guide.description}`);
});
```

### Mobile App
```kotlin
// Android/Kotlin
suspend fun getSimpleUserGuides(groupKey: String? = null): List<UserGuideSimple> {
    val request = GetUserGuideSimpleRequest.newBuilder()
        .apply { groupKey?.let { setGroupKey(it) } }
        .setPage(1)
        .setSize(20)
        .build()
    
    val response = client.getUserGuideSimple(request)
    return response.dataList
}
```

### React Component
```jsx
const UserGuideList = () => {
  const [guides, setGuides] = useState([]);

  useEffect(() => {
    fetch('/v2/hub/user-guides/simple')
      .then(res => res.json())
      .then(data => setGuides(data.data));
  }, []);

  return (
    <div>
      {guides.map(guide => (
        <div key={guide.id} className="guide-card">
          <h3>{guide.title}</h3>
          <p>{guide.description}</p>
        </div>
      ))}
    </div>
  );
};
```

## 📁 Files Changed

1. **shared/protobuf/schema/hub/user_guide.proto**
   - Added `rpc GetUserGuideSimple`
   - Added `GetUserGuideSimpleRequest`
   - Added `GetUserGuideSimpleResponse`
   - Added `UserGuideSimple` message

2. **hub-service/internal/usecase/user_guide_usecase.go**
   - Added `GetSimpleList()` method to interface
   - Added `GetSimpleList()` implementation

3. **hub-service/infra/handler/user_guide_handler.go**
   - Added `GetUserGuideSimple()` handler
   - Added description truncation logic (max 80 chars)

## 🔍 Testing

### Test với grpcurl
```bash
grpcurl -plaintext \
  -d '{"groupKey": "post", "page": 1, "size": 10}' \
  localhost:50051 \
  hubpb.UserGuideService/GetUserGuideSimple
```

### Test với REST
```bash
# Lấy tất cả
curl http://localhost:8080/v2/hub/user-guides/simple

# Lọc theo group
curl http://localhost:8080/v2/hub/user-guides/simple?groupKey=post

# Phân trang
curl http://localhost:8080/v2/hub/user-guides/simple?page=2&size=5
```

## ⚠️ Notes

- Proto code **PHẢI được generate** trước khi build service
- Nếu lỗi `undefined: hubpb.GetUserGuideSimpleRequest` → chạy generate proto
- API endpoint: `/v2/hub/user-guides/simple` (khác với admin API `/v2/hub/user-guides`)
- Description tự động truncate xuống 80 ký tự + "..."
- Response chỉ có id, title, description (không có createdAt, updatedAt, groupKey)

## 🚀 Features

- ✅ **Lightweight Response**: Chỉ trả về data cần thiết
- ✅ **Description Truncation**: Tự động cắt description xuống 80 ký tự
- ✅ **Pagination Support**: Hỗ trợ phân trang
- ✅ **Group Filter**: Filter theo groupKey
- ✅ **Mobile Optimized**: Tối ưu cho mobile app
- ✅ **Easy Integration**: Dễ tích hợp vào frontend/mobile

