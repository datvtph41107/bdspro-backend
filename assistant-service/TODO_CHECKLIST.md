# ✅ Checklist - Hoàn thiện OpenAI Integration

## Bước 1: Generate Wire ✅ (BẮT BUỘC)

```bash
cd assistant-service
wire ./wire
```

**Kiểm tra**: File `wire/wire_gen.go` đã được cập nhật với OpenAI client

## Bước 2: Set API Keys ⚠️

```bash
# Tạo file .env hoặc set environment variables
export OPENAI_API_KEY="sk-your-openai-api-key"
export DEEPSEEK_API_KEY="sk-your-deepseek-api-key"
```

**Lưu ý**: Cần ít nhất 1 API key để service hoạt động

## Bước 3: Build & Test

```bash
# Build service
go build -o assistant-service .

# Hoặc run trực tiếp
go run main.go
```

**Kiểm tra**: Service start thành công, không có error về missing dependencies

## Bước 4: Test API (Optional)

### Test với DeepSeek
```bash
curl -X POST "http://localhost:8061/v2/assistant/analyze/product" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Nhà 3 tầng, 80m2, Cầu Giấy, Hà Nội, giá 5 tỷ"
  }'
```

### Test Health Check
```bash
curl "http://localhost:8061/health"
```

## Bước 5: Update Handler (Optional - cho OpenAI methods)

Nếu muốn expose OpenAI methods qua gRPC/HTTP API:

1. **Update handler**: `infra/handler/assistant_handler.go`
   - Thêm methods gọi `usecase.AnalyzeProductTextWithOpenAI`
   - Thêm methods gọi `usecase.ChatWithOpenAI`
   - Thêm methods gọi `usecase.GenerateContentWithOpenAI`

2. **Update protobuf** (nếu cần): `shared/protobuf/schema/assistant/`
   - Thêm field `provider` (enum: DEEPSEEK, OPENAI, AUTO_FALLBACK)
   - Generate lại proto

## Bước 6: Documentation

- [ ] Cập nhật API documentation
- [ ] Thêm examples vào Postman/Swagger
- [ ] Update team wiki/docs nếu có

## Bước 7: Monitoring & Logging (Production)

```go
// Thêm logging cho provider selection
log.Printf("Using OpenAI provider for request_id=%s", requestID)

// Track metrics
metrics.RecordProviderUsage("openai")
metrics.RecordLatency("openai", processingTime)
```

## Bước 8: Cost Optimization (Production)

- [ ] Implement caching cho repeated queries
- [ ] Set request rate limits
- [ ] Monitor API usage và cost
- [ ] Setup alerts cho high usage

## 🎯 Các tính năng mở rộng (Future)

### 1. Smart Provider Selection
```go
// Tự động chọn provider dựa trên task complexity
func (u *AssistantUsecase) SelectProvider(task Task) Provider {
    if task.IsComplex() {
        return OpenAI
    }
    return DeepSeek
}
```

### 2. Load Balancing
- Round-robin giữa providers
- Weighted selection dựa trên success rate

### 3. A/B Testing
- So sánh kết quả của 2 providers
- Chọn provider tốt nhất cho từng domain

### 4. Response Caching
```go
// Cache kết quả để tránh duplicate API calls
cache.Set(hash(content), result, 1*time.Hour)
```

### 5. Request Queuing
- Queue requests khi API rate limit
- Batch processing cho efficiency

## 📋 Pre-Production Checklist

- [ ] Wire đã được generate
- [ ] API keys đã được set
- [ ] Service build thành công
- [ ] Unit tests passed
- [ ] Integration tests passed
- [ ] Load testing completed
- [ ] Error handling tested
- [ ] Logging configured
- [ ] Metrics/monitoring setup
- [ ] Documentation completed

## 🚨 Troubleshooting

### Lỗi "cannot find package"
```bash
go mod tidy
go mod download
```

### Lỗi "wire: no provider found"
- Kiểm tra lại `wire/wire.go`
- Đảm bảo đã import đúng packages
- Chạy `wire ./wire` lại

### Lỗi "API key not found"
- Kiểm tra environment variables
- Kiểm tra config/app.yml
- Restart service sau khi set env vars

### OpenAI timeout
- Tăng `timeout_seconds` trong config
- Kiểm tra network connection
- Thử với model nhẹ hơn (gpt-3.5-turbo)

## 📞 Support

- Xem [USAGE_EXAMPLE.md](./USAGE_EXAMPLE.md) cho ví dụ chi tiết
- Xem [README.md](./README.md) cho tổng quan
- Xem [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) cho chi tiết kỹ thuật

---

**Trạng thái hiện tại**: 
- ✅ Code implementation hoàn thành
- ⚠️ Cần generate wire
- ⚠️ Cần test

**Ưu tiên**: 
1. Generate wire (BẮT BUỘC)
2. Set API keys
3. Test basic functionality

