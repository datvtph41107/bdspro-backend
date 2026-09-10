# Utility Service

Service tiện ích cung cấp các tính năng chung cho hệ thống microservices, bao gồm:
- Event Queue: Quản lý hàng đợi các thao tác người dùng cần thực hiện

## Cấu trúc dự án

```
utility-service/
├── cmd/                    # Entry points
│   ├── grpc/              # gRPC server
│   └── http/              # HTTP server (deprecated)
├── config/                # Configuration files
│   ├── config.go
│   ├── develop.yml
│   └── product.yml
├── env/                   # Environment runtime
├── infra/                 # Infrastructure layer
│   ├── db/               # Database migration
│   ├── handler/          # gRPC handlers
│   ├── mapper/           # DTO mappers
│   └── postgre/          # PostgreSQL repositories
├── initial/              # App initialization
├── internal/             # Internal domain logic
│   ├── domain/          # Domain entities
│   ├── dto/             # Data Transfer Objects
│   ├── enums/           # Enumerations
│   ├── repo/            # Repository interfaces
│   └── usecase/         # Business logic
├── wire/                # Dependency injection
├── Dockerfile
├── go.mod
├── main.go
└── README.md
```

## API Endpoints

### Event Queue API

#### 1. Lấy danh sách event queue
```
GET /v2/utility/event-queue/list
Query params:
  - page: int (Trang)
  - size: int (Kích thước trang)
  - text: string (Tìm kiếm)
  - profileId: int (Profile ID)
  - organizationId: int (Organization ID)
  - status: []int (Trạng thái: 10-Pending, 20-Processing, 30-Completed, 40-Failed, 50-Cancelled)
  - eventType: []string (Loại event)
  - fromDate: string (Từ ngày YYYY-MM-DD)
  - toDate: string (Đến ngày YYYY-MM-DD)
```

#### 2. Lấy chi tiết event queue
```
GET /v2/utility/event-queue/detail/{id}
Path params:
  - id: int (ID event queue)
```

#### 3. Tạo mới event queue
```
POST /v2/utility/event-queue
Body: {
  "eventType": "string",       // Bắt buộc
  "eventName": "string",       // Bắt buộc
  "eventData": "string",       // Optional (JSON string)
  "metadata": "string",        // Optional (JSON string)
  "status": int,               // Optional (default: 10-Pending)
  "profileId": int,            // Optional
  "organizationId": int,       // Optional
  "scheduledAt": "string",     // Optional (ISO 8601 format)
  "maxRetries": int            // Optional (default: 3)
}
```

#### 4. Cập nhật event queue
```
PUT /v2/utility/event-queue/edit/{id}
Path params:
  - id: int (ID event queue)
Body: (same as create)
```

#### 5. Xóa event queue
```
DELETE /v2/utility/event-queue/{id}
Path params:
  - id: int (ID event queue)
```

#### 6. Cập nhật trạng thái
```
PUT /v2/utility/event-queue/status/{id}
Path params:
  - id: int (ID event queue)
Body: {
  "status": int,              // Bắt buộc
  "errorMessage": "string"    // Optional
}
```

## Cấu hình

### Database
```yaml
database:
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  database: utility_db
  sslmode: disable
```

### Server
```yaml
server:
  name: utility-service
  http_port: 8030
  tcp_port: 50030
```

## Chạy service

### Development
```bash
# Makefile nạp ../.env rồi ./.env và chọn topology host.
make server

# Hoặc chạy cả process và dependency trong container.
make compose-up
```

### Production
```bash
# Build artifact local; môi trường deploy inject biến runtime/secret riêng.
make build

QHPRO_ENVIRONMENT=production QHPRO_EXECUTION_MODE=container ./bin/hub-service grpc
```

### Docker
```bash
# Build image
docker build -t utility-service:latest .

# Run container
docker run -d \
  -p 8030:8030 \
  -p 50030:50030 \
  -e QHPRO_ENVIRONMENT=production \
  -e QHPRO_EXECUTION_MODE=container \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USERNAME=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=utility_db \
  utility-service:latest
```

## Database Schema

### tb_event_queue
```sql
CREATE TABLE tb_event_queue (
  id BIGSERIAL PRIMARY KEY,
  event_type VARCHAR(255) NOT NULL,
  event_name VARCHAR(255) NOT NULL,
  event_data TEXT,
  metadata TEXT,
  status INT DEFAULT 10,
  profile_id BIGINT,
  organization_id BIGINT,
  scheduled_at TIMESTAMP,
  processed_at TIMESTAMP,
  retry_count INT DEFAULT 0,
  max_retries INT DEFAULT 3,
  error_message TEXT,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_event_queue_status ON tb_event_queue(status);
CREATE INDEX idx_event_queue_profile ON tb_event_queue(profile_id);
CREATE INDEX idx_event_queue_org ON tb_event_queue(organization_id);
CREATE INDEX idx_event_queue_created ON tb_event_queue(created_at);
```

## Dependencies

- Go 1.24+
- PostgreSQL 13+
- gRPC
- GORM
- Wire (dependency injection)
- Viper (configuration)

## Kiến trúc

Service tuân theo **Clean Architecture** với các layer:

1. **Domain Layer** (`internal/domain`, `internal/dto`, `internal/enums`): Domain entities và business objects
2. **Use Case Layer** (`internal/usecase`): Business logic
3. **Repository Layer** (`internal/repo`): Repository interfaces
4. **Infrastructure Layer** (`infra/`): Implementation details
   - `postgre/`: PostgreSQL implementation
   - `handler/`: gRPC handlers
   - `mapper/`: DTO converters
5. **Dependency Injection** (`wire/`): Wire configuration

## Lưu ý

- Service sử dụng gRPC làm protocol chính
- HTTP server được đánh dấu deprecated, ưu tiên dùng gRPC
- Sử dụng soft delete cho tất cả entities
- Auto audit tracking (created_by, updated_by) qua GORM callbacks
- Hỗ trợ phân trang với default size=20, max size=100
