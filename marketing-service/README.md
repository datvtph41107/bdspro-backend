# Marketing Service

Marketing Service - Quản lý chiến dịch marketing theo kiến trúc microservices.

## Cấu trúc thư mục

```
marketing-service/
├── cmd/
│   ├── runtime.go              # Root command
│   └── grpc/
│       └── runtime.go          # gRPC server command
├── config/
│   ├── config.yaml             # Configuration file
│   └── runtime.go              # Config loader
├── env/
│   └── constants.go            # Environment constants
├── infrastructure/
│   ├── client/
│   │   └── payment_client.go   # Payment service client
│   ├── repository/
│   │   └── postgres/
│   │       └── campaign_postgres_repository.go
│   └── server/
│       └── grpc/
│           ├── runtime.go      # gRPC server runtime
│           └── interceptors/
│               └── interceptors.go
├── internal/
│   ├── domain/
│   │   ├── model/
│   │   │   ├── campaign.go     # Campaign domain model
│   │   │   └── package.go      # Package domain model
│   │   └── repository/
│   │       └── campaign_repository.go
│   ├── dto/
│   │   └── campaign.go         # Data Transfer Objects
│   ├── handler/
│   │   └── grpc_handler.go     # gRPC handler
│   └── usecase/
│       └── campaign_usecase.go # Business logic
├── wire/
│   ├── wire.go                 # Wire configuration
│   └── wire_gen.go             # Generated wire code
├── main.go                     # Application entry point
└── go.mod                      # Go module
```

## Tính năng đã implement

### 1. **Domain Models**
- **Campaign**: Model chiến dịch quảng cáo với đầy đủ các field
- **Package**: Model gói quảng cáo

### 2. **Repository Layer**
- **CampaignRepository**: Interface cho campaign operations
- **CampaignPostgresRepository**: Implementation sử dụng GORM

### 3. **Use Case Layer**
- **CampaignUsecase**: Business logic cho campaign management
- Tích hợp với payment service để tạo ví tự động

### 4. **gRPC Handler**
- **MarketingGrpcHandler**: Implement tất cả gRPC methods
- Các methods: Create, Update, Delete, Get, Search, GetByProduct, UpdateStatus

### 5. **Configuration**
- **config.yaml**: File cấu hình
- **env/constants.go**: Environment constants
- **config/runtime.go**: Config loader với viper

### 6. **Infrastructure**
- **PaymentClient**: Client gọi payment service qua gRPC
- **GRPCServer**: gRPC server với interceptors
- **Interceptors**: Logger, Recovery, Auth (TODO)

### 7. **Dependency Injection**
- **Wire**: Google Wire cho dependency injection
- **wire_gen.go**: Generated code

## Cách chạy

### 1. Build
```bash
go build -o marketing-service main.go
```

### 2. Chạy gRPC server
```bash
./marketing-service grpc-server
```

### 3. Generate wire code (khi cần)
```bash
cd wire && wire
```

## Configuration

### Environment Variables
- `GRPC_PORT`: Port cho gRPC server (default: 8209)
- `PAYMENT_SERVICE_HOST`: Payment service host (default: localhost)
- `PAYMENT_SERVICE_PORT`: Payment service port (default: 8208)
- `DB_DNS`: Database connection string
- `CONFIG_FILE`: Path to config file

### Config File (config.yaml)
```yaml
GRPC:
  PORT: 8209

Database:
  PRIMARY_DB_TYPE: postgres
  DB_DNS: host=localhost user=postgres password=postgres dbname=marketing port=5432 sslmode=disable

Payment:
  SERVICE_HOST: localhost
  SERVICE_PORT: 8208

ENV_RUNTIME: localhost
```

## API Endpoints

### Campaign Management
- `CreateCampaign`: Tạo chiến dịch mới
- `UpdateCampaign`: Cập nhật chiến dịch
- `DeleteCampaign`: Xóa chiến dịch
- `GetCampaign`: Lấy chi tiết chiến dịch
- `SearchCampaigns`: Tìm kiếm chiến dịch
- `GetCampaignsByProduct`: Lấy chiến dịch theo sản phẩm
- `UpdateCampaignStatus`: Cập nhật trạng thái

## Dependencies

- **gRPC**: Google gRPC
- **GORM**: Database ORM
- **Viper**: Configuration management
- **Wire**: Dependency injection
- **Cobra**: CLI framework
- **Fabric Logging**: Logging

## TODO

- [ ] Implement database connection
- [ ] Add authentication interceptor
- [ ] Add database migrations
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Add API documentation
- [ ] Add monitoring and metrics 