# Assistant Service - API Endpoints Summary

## 🎯 Tổng quan

Assistant Service cung cấp **12 API endpoints** với **3 AI providers**:

| Provider | Analyze | Chat | Generate | Default |
|----------|---------|------|----------|---------|
| **DeepSeek** | ✅ | ✅ | ✅ | ✅ (Mặc định) |
| **OpenAI** | ✅ | ✅ | ✅ | ❌ |
| **Gemini** | ✅ | ✅ | ✅ | ❌ |

---

## 📋 Danh sách API Endpoints

### 1️⃣ Analyze Product (Phân tích sản phẩm BĐS)

| Provider | Endpoint | Model |
|----------|----------|-------|
| DeepSeek | `POST /v2/assistant/deepseek/analyze/product` | deepseek-chat |
| OpenAI | `POST /v2/assistant/openai/analyze/product` | gpt-4o |
| Gemini | `POST /v2/assistant/gemini/analyze/product` | gemini-1.5-flash |
| Default | `POST /v2/assistant/analyze/product` | deepseek-chat |

**Request:**
```json
{
  "content": "Bán nhà 3 tầng, 4 phòng ngủ, 80m2, giá 5 tỷ, Quận 1",
  "context": "Optional context"
}
```

**Response:**
```json
{
  "product": {
    "name": "Nhà 3 tầng Quận 1",
    "area": 80.0,
    "bedroom": 4,
    "floor": 3,
    "salePrice": 5000000000,
    ...
  },
  "metadata": {
    "processing_time_ms": 2500,
    "model": "deepseek-chat"
  }
}
```

---

### 2️⃣ Chat (Trò chuyện với AI)

| Provider | Endpoint | Model |
|----------|----------|-------|
| DeepSeek | `POST /v2/assistant/deepseek/chat` | deepseek-chat |
| OpenAI | `POST /v2/assistant/openai/chat` | gpt-4o |
| Gemini | `POST /v2/assistant/gemini/chat` | gemini-1.5-flash |
| Default | `POST /v2/assistant/chat` | deepseek-chat |

**Request:**
```json
{
  "message": "Tư vấn cho tôi về thị trường BĐS Hà Nội",
  "history": [...],
  "session_id": "session_123"
}
```

**Response:**
```json
{
  "reply": "Thị trường BĐS Hà Nội hiện đang...",
  "session_id": "session_123",
  "metadata": {
    "processing_time_ms": 1800,
    "model": "deepseek-chat"
  }
}
```

---

### 3️⃣ Generate Content (Tạo nội dung)

| Provider | Endpoint | Model |
|----------|----------|-------|
| DeepSeek | `POST /v2/assistant/deepseek/generate` | deepseek-chat |
| OpenAI | `POST /v2/assistant/openai/generate` | gpt-4o |
| Gemini | `POST /v2/assistant/gemini/generate` | gemini-1.5-flash |
| Default | `POST /v2/assistant/generate` | deepseek-chat |

**Request:**
```json
{
  "prompt": "Viết mô tả hấp dẫn",
  "template": "product_description",
  "variables": {
    "name": "Căn hộ Landmark 81",
    "area": "120",
    "price": "8 tỷ"
  }
}
```

**Response:**
```json
{
  "content": "🏢 CĂNH HỘ CAO CẤP LANDMARK 81...",
  "metadata": {
    "processing_time_ms": 3200,
    "model": "deepseek-chat"
  }
}
```

---

### 4️⃣ Health Check

| Endpoint | Method |
|----------|--------|
| `/v2/assistant/health` | GET |

**Response:**
```json
{
  "status": "healthy",
  "uptime_seconds": 3600,
  "version": "1.0.0",
  "checks": {
    "deepseek_client": "healthy",
    "openai_client": "healthy",
    "gemini_client": "healthy",
    "service": "healthy"
  }
}
```

---

## 🚀 Quick Start

### Test nhanh với curl:

```bash
# DeepSeek (nhanh, rẻ)
curl -X POST http://localhost:8061/v2/assistant/deepseek/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 80m2, giá 5 tỷ"}'

# OpenAI (chất lượng cao)
curl -X POST http://localhost:8061/v2/assistant/openai/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 80m2, giá 5 tỷ"}'

# Gemini (cân bằng)
curl -X POST http://localhost:8061/v2/assistant/gemini/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 80m2, giá 5 tỷ"}'
```

---

## 📊 So sánh Providers

| Tiêu chí | DeepSeek | OpenAI | Gemini |
|----------|----------|--------|--------|
| **Tốc độ** | ⚡⚡⚡ Rất nhanh | ⚡⚡ Nhanh | ⚡⚡⚡ Rất nhanh |
| **Chi phí** | 💰 Rẻ nhất | 💰💰💰 Đắt nhất | 💰💰 Trung bình |
| **Chất lượng** | ⭐⭐⭐ Tốt | ⭐⭐⭐⭐⭐ Xuất sắc | ⭐⭐⭐⭐ Rất tốt |
| **Tiếng Việt** | ✅ Tốt | ✅ Xuất sắc | ✅ Rất tốt |
| **Use case** | Bulk processing | High-value content | Balanced usage |

---

## 💡 Khi nào dùng provider nào?

### DeepSeek - Dùng khi:
- ✅ Cần xử lý nhiều request (bulk)
- ✅ Ưu tiên tốc độ và tiết kiệm chi phí
- ✅ Nội dung đơn giản, rõ ràng
- ✅ Phân tích tin đăng thông thường

### OpenAI - Dùng khi:
- ✅ Cần chất lượng cao nhất
- ✅ Nội dung phức tạp, yêu cầu hiểu ngữ cảnh sâu
- ✅ Generate marketing content quan trọng
- ✅ Chat support VIP customers

### Gemini - Dùng khi:
- ✅ Cần cân bằng giữa chất lượng và chi phí
- ✅ Multimodal tasks (text + image trong tương lai)
- ✅ Fallback cho DeepSeek nếu fail
- ✅ A/B testing so với các provider khác

---

## 🔧 Configuration

Xem file `config/app.yml`:

```yaml
deepseek:
  api_key: "your-deepseek-key"
  model: "deepseek-chat"

openai:
  api_key: "your-openai-key"
  model: "gpt-4o"

gemini:
  api_key: "your-gemini-key"
  model: "gemini-1.5-flash"
```

---

## 📚 Documentation

- **Full API Guide**: [MULTI_PROVIDER_API_GUIDE.md](./MULTI_PROVIDER_API_GUIDE.md)
- **Swagger UI**: http://localhost:8061/swagger/index.html
- **Health Check**: http://localhost:8061/v2/assistant/health

---

## 🎯 Implementation Summary

**Total APIs**: 12 endpoints
- Analyze Product: 4 endpoints (3 providers + 1 default)
- Chat: 4 endpoints (3 providers + 1 default)
- Generate: 4 endpoints (3 providers + 1 default)

**Files Created/Modified**:
- ✅ `infra/client/gemini_client.go` (NEW)
- ✅ `internal/interface/provider/gemini_provider.go` (NEW)
- ✅ `infra/handler/assistant_handler.go` (UPDATED - added 9 handlers)
- ✅ `internal/usecases/assistant_usecase.go` (UPDATED - added 3 methods)
- ✅ `shared/protobuf/schema/assistant/assistant.proto` (UPDATED - added 9 RPCs)
- ✅ `wire/wire.go` (UPDATED - added Gemini DI)
- ✅ `MULTI_PROVIDER_API_GUIDE.md` (NEW)

**Build Status**: ✅ Successful

---

**Last Updated**: 2025-01-06  
**Version**: 2.0.0

