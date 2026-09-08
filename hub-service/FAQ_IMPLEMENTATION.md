# FAQ Implementation Guide

## ✅ Đã hoàn thành

### 1. Proto Definitions
- ✅ `shared/protobuf/schema/hub/user_guide.proto` (đổi tên thành `hub_services.proto` sau)
- ✅ Định nghĩa FAQService với đầy đủ CRUD methods
- ✅ Messages: FAQ, FAQDetail, FAQSimple
- ✅ Request/Response cho tất cả operations

### 2. Domain Layer
- ✅ `hub-service/internal/domain/faq_entity.go`
  - FAQEntity với Question, Answer, GroupKey
  - Embed BaseEntity (ID, CreatedAt, UpdatedAt, DeletedAt, AuditBase)

### 3. Repository Layer
- ✅ `hub-service/internal/repo/faq_repo.go` - Interface
- ✅ `hub-service/infra/postgre/faq_postgres.go` - Implementation
  - Embed CrudRepo[FAQEntity]
  - GetListWithFilter method với search và filter
  - GetSimpleListWithText method với text search (lowercase, space → %)

### 4. Usecase Layer
- ✅ `hub-service/internal/usecase/faq_usecase.go`
  - Embed BaseUsecase
  - GetListWithFilter method (Admin)
  - GetSimpleList method cho user API với text search support

### 5. Handler Layer
- ✅ `hub-service/infra/handler/faq_handler.go`
  - CreateFAQ, UpdateFAQ, DeleteFAQ
  - GetFAQ, GetFAQList (Admin)
  - GetFAQSimple (User) - truncate answer to 50 chars

### 6. Mapper Layer
- ✅ `hub-service/infra/mapper/faq_mapper.go`
  - CreateRequestToEntity
  - UpdateRequestToEntity
  - EntityToDetailProto
  - EntitiesToProto

### 7. Database Migration
- ✅ `hub-service/migrate/004_create_faqs_table.sql`
  - Bảng faqs với full-text search index
  - Trigger auto update updated_at

### 8. Documentation
- ✅ `hub-service/FAQ_API.md` - Chi tiết API usage
- ✅ `hub-service/FAQ_IMPLEMENTATION.md` - Hướng dẫn này

## 🔧 Cần làm tiếp

### 1. Generate Proto Code

```bash
cd shared/protobuf/schema/hub

# Generate Go protobuf + gRPC code
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  user_guide.proto

# Generate gRPC Gateway
protoc --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  user_guide.proto
```

### 2. Run Database Migration

```bash
cd hub-service

# Run migration
psql -U your_user -d your_db -f migrate/004_create_faqs_table.sql
```

Hoặc dùng migration tool:
```bash
migrate -path migrate -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up
```

### 3. Update Wire Dependency Injection

Thêm vào `hub-service/wire/wire.go`:

```go
// FAQ dependencies
postgre.NewFAQRepo,
usecase.NewFAQUsecase,
mapper.NewFAQMapper,
handler.NewFAQHandler,
```

Generate wire code:
```bash
cd hub-service/wire
wire
```

### 4. Register gRPC Service

Thêm vào `hub-service/cmd/grpc/main.go`:

```go
// Register FAQ service
hubpb.RegisterFAQServiceServer(grpcServer, wireApp.FAQHandler)
```

### 5. Register Gateway Routes

Thêm vào gateway registration:

```go
// FAQ routes
err = hubpb.RegisterFAQServiceHandlerFromEndpoint(ctx, gwmux, grpcServerEndpoint, opts)
if err != nil {
    log.Fatal("Cannot register FAQ gateway")
}
```

## 📋 API Endpoints Mới

### User API (Public)
```
GET /v2/hub/faqs
```

### Admin APIs
```
POST   /v2/hub/faqs
PUT    /v2/hub/faqs/{id}
DELETE /v2/hub/faqs/{id}
GET    /v2/hub/faqs/{id}
GET    /v2/hub/admin/faqs
```

## 🧪 Testing

### 1. Test User API (Simple)
```bash
# Lấy tất cả FAQs
curl http://localhost:8080/v2/hub/faqs

# Lọc theo groupKey
curl http://localhost:8080/v2/hub/faqs?groupKey=post

# Search theo text
curl "http://localhost:8080/v2/hub/faqs?text=tạo%20bài%20đăng"

# Search + filter
curl "http://localhost:8080/v2/hub/faqs?text=upload&groupKey=post"

# Phân trang
curl http://localhost:8080/v2/hub/faqs?page=1&size=10
```

### 2. Test Admin API
```bash
# Tạo FAQ mới
curl -X POST http://localhost:8080/v2/hub/faqs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "question": "Làm thế nào để tạo bài đăng?",
    "answer": "Vào menu Bài đăng > Tạo mới...",
    "groupKey": "post"
  }'

# Cập nhật FAQ
curl -X PUT http://localhost:8080/v2/hub/faqs/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "id": 1,
    "question": "Updated question",
    "answer": "Updated answer",
    "groupKey": "post"
  }'

# Xóa FAQ
curl -X DELETE http://localhost:8080/v2/hub/faqs/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# Lấy danh sách (Admin)
curl http://localhost:8080/v2/hub/admin/faqs?question=tạo \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Test với grpcurl
```bash
# Create FAQ
grpcurl -plaintext -d '{
  "question": "Test question?",
  "answer": "Test answer",
  "groupKey": "general"
}' localhost:50051 hubpb.FAQService/CreateFAQ

# Get Simple FAQs
grpcurl -plaintext -d '{
  "groupKey": "post",
  "page": 1,
  "size": 10
}' localhost:50051 hubpb.FAQService/GetFAQSimple
```

## 📁 File Structure

```
hub-service/
├── internal/
│   ├── domain/
│   │   └── faq_entity.go          ✅ Entity definition
│   ├── repo/
│   │   └── faq_repo.go            ✅ Repository interface
│   └── usecase/
│       └── faq_usecase.go         ✅ Business logic
├── infra/
│   ├── handler/
│   │   └── faq_handler.go         ✅ gRPC handlers
│   ├── mapper/
│   │   └── faq_mapper.go          ✅ DTO conversions
│   └── postgre/
│       └── faq_postgres.go        ✅ DB implementation
├── migrate/
│   └── 004_create_faqs_table.sql  ✅ Database migration
└── FAQ_API.md                     ✅ API documentation

shared/protobuf/schema/hub/
└── user_guide.proto               ✅ Proto definitions (FAQ + UserGuide)
```

## 🎯 Features Implemented

- ✅ **Full CRUD**: Create, Read, Update, Delete FAQs
- ✅ **User API**: Simple API chỉ trả id, question, answer (max 50 chars)
- ✅ **Admin API**: Full API với search, filter, pagination
- ✅ **Answer Truncation**: Tự động cắt answer xuống 50 ký tự cho user API
- ✅ **Text Search**: Search theo text với lowercase và space → % (LIKE pattern)
- ✅ **Group Filter**: Lọc FAQs theo groupKey
- ✅ **Search**: Full-text search theo question (Admin)
- ✅ **Pagination**: Hỗ trợ phân trang
- ✅ **Soft Delete**: Xóa mềm với DeletedAt
- ✅ **Audit Trail**: CreatedBy, UpdatedBy tracking
- ✅ **Clean Architecture**: Tách rõ layers (domain, repo, usecase, handler)

## 🚀 Next Steps

1. **Generate Proto Code** ✅ Quan trọng nhất
2. **Run Migration** ✅ Tạo bảng trong DB
3. **Update Wire** ✅ Dependency injection
4. **Register Services** ✅ gRPC và Gateway
5. **Test APIs** ✅ Verify functionality
6. **Write Unit Tests** (Optional)
7. **Update Swagger Docs** (Optional)

## ⚠️ Important Notes

- Proto code **PHẢI được generate** trước khi build
- Migration **PHẢI chạy** trước khi start service
- Wire code **PHẢI regenerate** sau khi thêm dependencies
- User API không yêu cầu auth, Admin APIs cần auth
- Answer tự động truncate to 50 chars cho user API
- Full-text search index cho performance tốt hơn

## 📝 Example Usage

### Frontend Integration
```javascript
// Lấy FAQs cho trang help
async function loadHelpFAQs() {
  const response = await fetch('/v2/hub/faqs?groupKey=general');
  const { data } = await response.json();
  
  return data.map(faq => ({
    id: faq.id,
    question: faq.question,
    answer: faq.answer
  }));
}

// Hiển thị FAQ accordion
const FAQSection = () => {
  const [faqs, setFaqs] = useState([]);
  
  useEffect(() => {
    loadHelpFAQs().then(setFaqs);
  }, []);
  
  return (
    <div>
      {faqs.map(faq => (
        <FAQItem key={faq.id} {...faq} />
      ))}
    </div>
  );
};
```

## 🎉 Ready to Use!

Sau khi complete các bước trong "Cần làm tiếp", FAQ API sẽ hoàn toàn sẵn sàng!

