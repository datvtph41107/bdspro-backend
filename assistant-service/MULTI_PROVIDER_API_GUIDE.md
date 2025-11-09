# Assistant Service - Multi-Provider AI API Guide

Hướng dẫn sử dụng các API tích hợp 3 AI providers trong Assistant Service.

## 📋 Tổng quan

Assistant Service hỗ trợ 3 AI providers với endpoints riêng biệt:
- **DeepSeek** - Nhanh và rẻ, tốt cho tiếng Việt
- **OpenAI** - GPT-4, chất lượng cao nhất
- **Gemini** - Google AI, cân bằng giá/hiệu suất

**Tất cả 3 providers đều có đầy đủ API cho:**
- ✅ Analyze Product (Phân tích sản phẩm BĐS)
- ✅ Chat (Trò chuyện với AI)
- ✅ Generate Content (Tạo nội dung)

## 🔑 Cấu hình

### File config (`config/app.yml`)

```yaml
gemini:
  api_key: "YOUR_GEMINI_API_KEY"
  base_url: "https://generativelanguage.googleapis.com/v1beta"  # Optional, default
  model: "gemini-1.5-flash"  # Optional, default
  timeout_seconds: 60  # Optional, default
```

### Environment Variables

```bash
export GEMINI_API_KEY="your-api-key-here"
```

## 🚀 API Endpoints

Mỗi provider có 3 endpoints giống nhau:

| Provider | Analyze Product | Chat | Generate |
|----------|----------------|------|----------|
| **DeepSeek** | `/v2/assistant/deepseek/analyze/product` | `/v2/assistant/deepseek/chat` | `/v2/assistant/deepseek/generate` |
| **OpenAI** | `/v2/assistant/openai/analyze/product` | `/v2/assistant/openai/chat` | `/v2/assistant/openai/generate` |
| **Gemini** | `/v2/assistant/gemini/analyze/product` | `/v2/assistant/gemini/chat` | `/v2/assistant/gemini/generate` |
| **Default** | `/v2/assistant/analyze/product` | `/v2/assistant/chat` | `/v2/assistant/generate` |

> **Note:** Default endpoints sử dụng DeepSeek

---

## 📝 API Examples

### 1. Phân tích sản phẩm BĐS (Analyze Product)

Tất cả 3 providers có cùng request/response format.

**Endpoints:**
- `POST /v2/assistant/deepseek/analyze/product` (DeepSeek)
- `POST /v2/assistant/openai/analyze/product` (OpenAI)
- `POST /v2/assistant/gemini/analyze/product` (Gemini)

**Request:**
```json
{
  "content": "Bán nhà 3 tầng, 4 phòng ngủ, 80m2, giá 5 tỷ, Quận 1, TP.HCM",
  "context": "Gần trường học và chợ"
}
```

**Response:**
```json
{
  "product": {
    "name": "Nhà 3 tầng Quận 1",
    "area": 80.0,
    "transaction_type": 10,
    "address": "Quận 1, TP.HCM",
    "province": "TP.HCM",
    "district": "Quận 1",
    "bedroom": 4,
    "floor": 3,
    "salePrice": 5000000000,
    "propertyType": "Nhà riêng"
  },
  "metadata": {
    "processing_time_ms": 2500,
    "model": "gemini-1.5-flash"
  }
}
```

**cURL Examples:**

```bash
# DeepSeek (nhanh nhất, rẻ nhất)
curl -X POST http://localhost:8061/v2/assistant/deepseek/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Bán nhà 3 tầng, 80m2, giá 5 tỷ, Quận 1"}'

# OpenAI (chất lượng cao nhất)
curl -X POST http://localhost:8061/v2/assistant/openai/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Bán nhà 3 tầng, 80m2, giá 5 tỷ, Quận 1"}'

# Gemini (cân bằng)
curl -X POST http://localhost:8061/v2/assistant/gemini/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Bán nhà 3 tầng, 80m2, giá 5 tỷ, Quận 1"}'
```

---

### 2. Chat với AI

**Endpoints:**
- `POST /v2/assistant/deepseek/chat` (DeepSeek)
- `POST /v2/assistant/openai/chat` (OpenAI)
- `POST /v2/assistant/gemini/chat` (Gemini)

**Request:**
```json
{
  "message": "Tư vấn cho tôi về thị trường BĐS Hà Nội",
  "history": [
    {
      "role": "user",
      "content": "Xin chào",
      "timestamp": 1699000000
    },
    {
      "role": "assistant",
      "content": "Xin chào! Tôi có thể giúp gì cho bạn?",
      "timestamp": 1699000001
    }
  ],
  "context": "User quan tâm khu vực Cầu Giấy",
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
    "model": "gemini-1.5-flash"
  }
}
```

**cURL Examples:**

```bash
# DeepSeek Chat
curl -X POST http://localhost:8061/v2/assistant/deepseek/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Tư vấn cho tôi về thị trường BĐS Hà Nội"}'

# OpenAI Chat
curl -X POST http://localhost:8061/v2/assistant/openai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Tư vấn cho tôi về thị trường BĐS Hà Nội"}'

# Gemini Chat
curl -X POST http://localhost:8061/v2/assistant/gemini/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Tư vấn cho tôi về thị trường BĐS Hà Nội"}'
```

---

### 3. Generate nội dung

**Endpoints:**
- `POST /v2/assistant/deepseek/generate` (DeepSeek)
- `POST /v2/assistant/openai/generate` (OpenAI)
- `POST /v2/assistant/gemini/generate` (Gemini)

**Request:**
```json
{
  "prompt": "Viết mô tả hấp dẫn",
  "template": "product_description",
  "variables": {
    "name": "Căn hộ cao cấp Landmark 81",
    "area": "120",
    "location": "Quận Bình Thạnh",
    "price": "8 tỷ"
  },
  "max_tokens": 2000,
  "temperature": 0.7
}
```

**Supported Templates:**
- `product_description` - Mô tả sản phẩm BĐS
- `email` - Viết email chuyên nghiệp
- `social_post` - Bài đăng mạng xã hội

**Response:**
```json
{
  "content": "🏢 CĂNH HỘ CAO CẤP LANDMARK 81...",
  "metadata": {
    "processing_time_ms": 3200,
    "model": "gemini-1.5-flash"
  }
}
```

**cURL Examples:**

```bash
# DeepSeek Generate
curl -X POST http://localhost:8061/v2/assistant/deepseek/generate \
  -H "Content-Type: application/json" \
  -d '{
    "template": "product_description",
    "variables": {"name": "Căn hộ Landmark 81", "area": "120", "price": "8 tỷ"}
  }'

# OpenAI Generate
curl -X POST http://localhost:8061/v2/assistant/openai/generate \
  -H "Content-Type: application/json" \
  -d '{
    "template": "product_description",
    "variables": {"name": "Căn hộ Landmark 81", "area": "120", "price": "8 tỷ"}
  }'

# Gemini Generate
curl -X POST http://localhost:8061/v2/assistant/gemini/generate \
  -H "Content-Type: application/json" \
  -d '{
    "template": "product_description",
    "variables": {"name": "Căn hộ Landmark 81", "area": "120", "price": "8 tỷ"}
  }'
```

---

## 🔄 So sánh Providers

| Feature | DeepSeek | OpenAI | Gemini |
|---------|----------|--------|--------|
| **Model** | deepseek-chat | gpt-4o | gemini-1.5-flash |
| **Speed** | ⚡⚡⚡ Rất nhanh | ⚡⚡ Nhanh | ⚡⚡⚡ Rất nhanh |
| **Cost** | 💰 Rẻ nhất | 💰💰💰 Đắt nhất | 💰💰 Trung bình |
| **Quality** | ⭐⭐⭐ Tốt | ⭐⭐⭐⭐⭐ Xuất sắc | ⭐⭐⭐⭐ Rất tốt |
| **Vietnamese** | ✅ Tốt | ✅ Xuất sắc | ✅ Rất tốt |
| **Endpoint** | `/v2/assistant/deepseek/*` | `/v2/assistant/openai/*` | `/v2/assistant/gemini/*` |

## 🧪 Testing

### Quick Test với curl

```bash
# Test tất cả 3 providers cho Analyze
curl -X POST http://localhost:8061/v2/assistant/deepseek/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 80m2, 3 tầng, giá 5 tỷ"}'

curl -X POST http://localhost:8061/v2/assistant/openai/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 80m2, 3 tầng, giá 5 tỷ"}'

curl -X POST http://localhost:8061/v2/assistant/gemini/analyze/product \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 80m2, 3 tầng, giá 5 tỷ"}'

# Test Chat
curl -X POST http://localhost:8061/v2/assistant/deepseek/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Xin chào"}'

# Test Generate
curl -X POST http://localhost:8061/v2/assistant/openai/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Viết mô tả nhà đẹp"}'
```

### Test với Postman

Import collection từ Swagger docs:
```
http://localhost:8061/swagger/index.html
```

## 📊 Performance Tips

1. **Chọn model phù hợp:**
   - `gemini-1.5-flash` - Nhanh, giá rẻ (recommended)
   - `gemini-1.5-pro` - Chất lượng cao hơn, chậm hơn

2. **Tối ưu timeout:**
   ```yaml
   gemini:
     timeout_seconds: 60  # Tăng nếu request lớn
   ```

3. **Sử dụng context hiệu quả:**
   - Giữ context ngắn gọn
   - Chỉ truyền thông tin cần thiết

## 🐛 Troubleshooting

### Lỗi: "gemini API error: status=401"

**Nguyên nhân:** API key không hợp lệ

**Giải pháp:**
1. Kiểm tra `GEMINI_API_KEY` trong env
2. Verify API key tại: https://makersuite.google.com/app/apikey
3. Đảm bảo API key có quyền truy cập Gemini API

### Lỗi: "no content in gemini response"

**Nguyên nhân:** Response từ Gemini bị lỗi hoặc content bị filter

**Giải pháp:**
1. Kiểm tra content request có vi phạm policy không
2. Xem log để debug response structure
3. Thử lại với content khác

### Lỗi: "failed to parse product analysis"

**Nguyên nhân:** AI trả về JSON không đúng format

**Giải pháp:**
1. Kiểm tra prompt có rõ ràng không
2. Thêm context để AI hiểu rõ hơn
3. Fallback sang OpenAI nếu Gemini không ổn định

## 🔐 Security Notes

⚠️ **QUAN TRỌNG:**
- KHÔNG commit API key vào git
- Sử dụng environment variables trong production
- Rotate API key định kỳ
- Monitor usage để tránh vượt quota

## 📖 Related Documentation

- [Gemini API Official Docs](https://ai.google.dev/docs)
- [Assistant Service README](./README.md)
- [OpenAI Integration](./USAGE_EXAMPLE.md)

## 💡 Use Cases

### 1. Phân tích tin đăng tự động với provider lựa chọn

```javascript
// Frontend example - chọn provider theo nhu cầu
const analyzeProduct = async (text, provider = 'deepseek') => {
  const endpoints = {
    deepseek: '/v2/assistant/deepseek/analyze/product',
    openai: '/v2/assistant/openai/analyze/product', 
    gemini: '/v2/assistant/gemini/analyze/product'
  };
  
  const response = await fetch(endpoints[provider], {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content: text })
  });
  return await response.json();
};

// Use case: Analyze nhanh với DeepSeek
const result = await analyzeProduct(text, 'deepseek');

// Use case: Analyze chính xác với OpenAI
const result = await analyzeProduct(text, 'openai');
```

### 2. Chat support BĐS với fallback

```javascript
// Chatbot với fallback từ DeepSeek → Gemini → OpenAI
const chatWithAI = async (message, sessionId) => {
  const providers = ['deepseek', 'gemini', 'openai'];
  
  for (const provider of providers) {
    try {
      const response = await fetch(`/v2/assistant/${provider}/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message, session_id: sessionId })
      });
      
      if (response.ok) {
        return await response.json();
      }
    } catch (error) {
      console.log(`${provider} failed, trying next...`);
    }
  }
  
  throw new Error('All AI providers failed');
};
```

### 3. A/B Testing providers

```javascript
// Marketing automation với A/B testing
const generateDescription = async (product, testGroup) => {
  // Group A: DeepSeek (nhanh, rẻ)
  // Group B: OpenAI (chất lượng cao)
  const provider = testGroup === 'A' ? 'deepseek' : 'openai';
  
  const response = await fetch(`/v2/assistant/${provider}/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      template: 'product_description',
      variables: {
        name: product.name,
        area: product.area.toString(),
        location: product.location,
        price: product.price.toString()
      }
    })
  });
  
  const result = await response.json();
  
  // Track performance metrics
  trackMetrics({
    provider,
    testGroup,
    processingTime: result.metadata.processing_time_ms,
    model: result.metadata.model
  });
  
  return result;
};
```

### 4. Cost optimization strategy

```javascript
// Optimize cost dựa trên độ phức tạp
const smartAnalyze = async (text) => {
  const wordCount = text.split(' ').length;
  
  let provider;
  if (wordCount < 50) {
    // Simple text → DeepSeek (rẻ nhất)
    provider = 'deepseek';
  } else if (wordCount < 200) {
    // Medium text → Gemini (cân bằng)
    provider = 'gemini';
  } else {
    // Complex text → OpenAI (chính xác nhất)
    provider = 'openai';
  }
  
  console.log(`Using ${provider} for ${wordCount} words`);
  
  return await fetch(`/v2/assistant/${provider}/analyze/product`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content: text })
  }).then(r => r.json());
};
```

## 📝 Changelog

### v2.0.0 (2025-01-06)
- ✅ **Multi-Provider Support**: DeepSeek, OpenAI, Gemini
- ✅ **9 API Endpoints Total**: 3 providers × 3 functions each
- ✅ Provider-specific endpoints:
  - DeepSeek: `/v2/assistant/deepseek/*`
  - OpenAI: `/v2/assistant/openai/*`
  - Gemini: `/v2/assistant/gemini/*`
- ✅ Gemini client integration
- ✅ OpenAI dedicated endpoints
- ✅ DeepSeek dedicated endpoints
- ✅ Wire dependency injection cho tất cả providers
- ✅ Protobuf schema hoàn chỉnh
- ✅ Swagger documentation đầy đủ
- ✅ Comprehensive API Guide với use cases

---

**Maintained by:** Assistant Service Team  
**Last Updated:** 2025-01-06

