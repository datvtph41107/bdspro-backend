# Assistant Service

AI-powered assistant service sử dụng DeepSeek và OpenAI để phân tích và xử lý thông tin bất động sản.

## 📋 Tổng quan

Assistant Service cung cấp các tính năng AI:
- **Phân tích văn bản**: Trích xuất thông tin sản phẩm BĐS từ văn bản tiếng Việt
- **Chatbot**: Conversation với AI assistant
- **Generate nội dung**: Tạo mô tả sản phẩm, email, social posts
- **Multi-provider support**: Hỗ trợ cả DeepSeek và OpenAI

## 🚀 Cài đặt

### Prerequisites

- Go 1.21+
- Protocol Buffers compiler (protoc)
- Wire (dependency injection)
- DeepSeek API Key (optional)
- OpenAI API Key (optional)

### 1. Generate Protobuf

```bash
# Từ root project
./generate_assistant_proto.sh
```

### 2. Generate Wire

```bash
cd assistant-service
wire ./wire
```

### 3. Cấu hình runtime

`assistant-service/.env` là file cấu hình duy nhất của process này. Các biến
`*_API_KEY` để trống là hợp lệ: provider sẽ trả lỗi capability rõ ràng khi
được gọi, còn dev có thể chạy health/contract test mà không cần secret. Khi
cần gọi provider thật, inject key từ secret manager hoặc shell local; không
commit key vào Git.

### 4. Build & Run

```bash
cd assistant-service
make test
make build
./bin/assistant-service
```

Service sẽ chạy trên:
- **gRPC**: Port `50061`
- **HTTP** (qua gateway): Port `8080`

## 📡 API Endpoints

### 1. Analyze Product Text

**Endpoint**: `POST /v2/assistant/analyze/product`

Phân tích văn bản mô tả BĐS và trích xuất thông tin.

**Request**:
```json
{
  "content": "Cần bán nhà 3 tầng, 4 phòng ngủ, 3WC, diện tích 80m2, mặt tiền 5m, hướng Đông, giá 5 tỷ, sổ đỏ chính chủ, tại Cầu Giấy, Hà Nội",
  "context": "Thông tin bổ sung (optional)"
}
```

**Response**:
```json
{
  "product": {
    "name": "Nhà 3 tầng Cầu Giấy",
    "area": 80,
    "address": "Cầu Giấy, Hà Nội",
    "province": "Hà Nội",
    "district": "Cầu Giấy",
    "property_type": "Nhà riêng",
    "transaction_type": 1,
    "sale_price": 5000000000,
    "bedroom": 4,
    "bathroom": 3,
    "floor": 3,
    "frontage": 5,
    "orientation": "Đông",
    "legal_doc": "Sổ đỏ"
  },
  "metadata": {
    "processing_time_ms": 1500,
    "tokens_used": 600,
    "model": "deepseek-chat"
  }
}
```

### 2. Chat

**Endpoint**: `POST /v2/assistant/chat`

Conversation với AI assistant.

**Request**:
```json
{
  "message": "Cho tôi biết giá nhà ở Cầu Giấy như thế nào?",
  "history": [
    {
      "role": "user",
      "content": "Xin chào",
      "timestamp": 1704067200
    },
    {
      "role": "assistant",
      "content": "Chào bạn! Tôi có thể giúp gì cho bạn?",
      "timestamp": 1704067201
    }
  ],
  "context": "Context bổ sung",
  "session_id": "session_123"
}
```

**Response**:
```json
{
  "reply": "Giá nhà ở Cầu Giấy hiện tại dao động từ 100-200 triệu/m² tùy vào vị trí và loại nhà...",
  "session_id": "session_123",
  "metadata": {
    "processing_time_ms": 2000,
    "tokens_used": 450,
    "model": "deepseek-chat"
  }
}
```

### 3. Generate Content

**Endpoint**: `POST /v2/assistant/generate`

Tạo nội dung theo template.

**Request**:
```json
{
  "prompt": "Tạo mô tả sản phẩm",
  "template": "product_description",
  "variables": {
    "name": "Nhà 3 tầng Cầu Giấy",
    "area": "80",
    "location": "Cầu Giấy, Hà Nội",
    "price": "5 tỷ"
  },
  "max_tokens": 2000,
  "temperature": 0.7
}
```

**Response**:
```json
{
  "content": "🏠 **Nhà 3 tầng đẹp tại Cầu Giấy, Hà Nội**\n\n✨ Diện tích: 80m²\n📍 Vị trí: Cầu Giấy, Hà Nội...",
  "metadata": {
    "processing_time_ms": 1800,
    "tokens_used": 500,
    "model": "deepseek-chat"
  }
}
```

### 4. Health Check

**Endpoint**: `GET /v2/assistant/health`

Kiểm tra trạng thái service.

**Response**:
```json
{
  "status": "healthy",
  "uptime_seconds": 3600,
  "version": "1.0.0",
  "checks": {
    "deepseek_client": "healthy",
    "service": "healthy"
  }
}
```

## 🔧 Configuration

### config/app.yml

```yaml
app:
  name: "Assistant Service"
  version: "1.0.0"
  port:
    grpc: "50061"

deepseek:
  api_key: "${DEEPSEEK_API_KEY}"
  base_url: "https://api.deepseek.com/v1"
  model: "deepseek-v3"
  timeout_seconds: 30
  max_tokens: 4000

openai:
  api_key: "${OPENAI_API_KEY}"
  base_url: "https://api.openai.com/v1"
  model: "gpt-4o-mini"
  timeout_seconds: 60
  max_tokens: 4000
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DEEPSEEK_API_KEY` | DeepSeek API key | Optional |
| `OPENAI_API_KEY` | OpenAI API key | Optional |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `CONFIG_PATH` | Config file path | `config` |

## 🏗️ Architecture

```
assistant-service/
├── cmd/grpc/          # gRPC server entry point
├── config/            # Configuration files
├── env/               # Environment setup
├── infra/
│   ├── client/        # AI client implementations (DeepSeek, OpenAI)
│   └── handler/       # gRPC handlers
├── internal/
│   ├── dto/           # Data Transfer Objects
│   ├── interface/     # Interfaces
│   │   └── provider/  # Provider interfaces
│   └── usecases/      # Business logic
├── initial/           # Initialization
├── wire/              # Dependency injection
├── main.go            # Main entry point
└── Dockerfile         # Docker configuration
```

## 🧪 Testing

```bash
# Test với curl
curl -X POST "http://localhost:8080/v2/assistant/analyze/product" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Nhà 3 tầng, 100m2, Cầu Giấy Hà Nội, giá 5 tỷ"
  }'

# Health check
curl "http://localhost:8080/v2/assistant/health"
```

## 🐳 Docker

### Build

```bash
docker build -t assistant-service:latest .
```

### Run

```bash
docker run -d \
  -p 50061:50061 \
  -e DEEPSEEK_API_KEY="sk-your-key" \
  -e ENVIRONMENT="production" \
  --name assistant-service \
  assistant-service:latest
```

## 📊 Performance

| Metric | Value |
|--------|-------|
| Average response time | ~1.5s |
| Max concurrent requests | 100 |
| Token usage (avg) | 500-800 tokens |
| Cost per request | ~$0.002 |

## 🔐 Security

- API key được lưu trong environment variables
- Không hardcode credentials trong code
- Support JWT authentication qua gateway
- Rate limiting tại gateway layer

## 🚨 Troubleshooting

### Lỗi "API_KEY not found"

```bash
# Cho DeepSeek
export DEEPSEEK_API_KEY="sk-your-deepseek-key"

# Cho OpenAI
export OPENAI_API_KEY="sk-your-openai-key"
```

### Lỗi connection timeout

Tăng timeout trong config:

```yaml
deepseek:
  timeout_seconds: 60

openai:
  timeout_seconds: 90
```

### Lỗi "failed to generate wire"

```bash
cd assistant-service
go get github.com/google/wire/cmd/wire
wire ./wire
```

## 📝 Development

### Add new API

1. Update proto file: `shared/protobuf/schema/assistant/assistant.proto`
2. Generate proto: `./generate_assistant_proto.sh`
3. Add handler method in `infra/handler/assistant_handler.go`
4. Add usecase method in `internal/usecases/assistant_usecase.go`
5. Test & deploy

### Code Style

- Follow Clean Architecture
- Use dependency injection (Wire)
- Interface-based design
- Comprehensive error handling

## 📚 References

- [DeepSeek API Documentation](https://platform.deepseek.com/docs)
- [OpenAI API Documentation](https://platform.openai.com/docs)
- [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Wire](https://github.com/google/wire)

## 📞 Support

- **Issues**: Create issue in GitLab
- **Documentation**: See `/docs` folder
- **Contact**: Development Team

## 📄 License

Copyright © 2025 - All rights reserved
