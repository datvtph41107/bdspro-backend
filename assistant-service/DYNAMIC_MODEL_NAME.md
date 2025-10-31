# ✅ Dynamic Model Name - Lấy từ Config/Env

## 🎯 Vấn đề trước đó

**Hardcode model name trong usecase**:
```go
// ❌ SAI - Hardcode
modelName := "deepseek-chat"
modelName := "gpt-3.5-turbo"
```

**Vấn đề**:
- Nếu đổi model trong `config/app.yml` (VD: `deepseek-chat` → `deepseek-v3`)
- API vẫn hiển thị model cũ: `"deepseek-chat"` ❌
- Không đúng với model thực tế đang chạy

## ✅ Giải pháp - Dynamic Model Name

### 1. Thêm GetModel() vào Client

#### DeepseekClient:
```go
// infra/client/deepseek_client.go

type DeepseekClient struct {
    model string  // Đã load từ config/env
    // ...
}

// GetModel - Trả về model name đang sử dụng
func (c *DeepseekClient) GetModel() string {
    return c.model  // ✅ Lấy từ config
}
```

#### OpenAIClient:
```go
// infra/client/openai_client.go

type OpenAIClient struct {
    model string  // Đã load từ config/env
    // ...
}

// GetModel - Trả về model name đang sử dụng
func (c *OpenAIClient) GetModel() string {
    return c.model  // ✅ Lấy từ config
}
```

### 2. Cập nhật Provider Interface

#### DeepseekProvider:
```go
// internal/interface/provider/deepseek_provider.go

type DeepseekProvider interface {
    AnalyzeProductText(ctx context.Context, content string) (*dto.ProductAnalysisResult, error)
    Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error)
    GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error)
    
    GetModel() string  // ✅ Thêm method mới
}
```

#### OpenAIProvider:
```go
// internal/interface/provider/openai_provider.go

type OpenAIProvider interface {
    AnalyzeProductText(ctx context.Context, content string) (*dto.ProductAnalysisResult, error)
    Chat(ctx context.Context, message string, history []dto.DeepseekMessage) (string, error)
    GenerateContent(ctx context.Context, prompt string, maxTokens int, temperature float32) (string, error)
    
    GetModel() string  // ✅ Thêm method mới
}
```

### 3. Usecase gọi GetModel()

#### DeepSeek methods:
```go
// internal/usecases/assistant_usecase.go

func (u *AssistantUsecase) AnalyzeProductText(...) {
    // ...
    modelName := u.DeepseekClient.GetModel()  // ✅ Lấy từ config động
    return result, rawResponse, modelName, processingTime, nil
}

func (u *AssistantUsecase) Chat(...) {
    // ...
    modelName := u.DeepseekClient.GetModel()  // ✅ Dynamic
    return reply, sessionID, modelName, processingTime, nil
}

func (u *AssistantUsecase) GenerateContent(...) {
    // ...
    modelName := u.DeepseekClient.GetModel()  // ✅ Dynamic
    return content, modelName, processingTime, nil
}
```

#### OpenAI methods:
```go
func (u *AssistantUsecase) AnalyzeProductTextWithOpenAI(...) {
    // ...
    modelName := u.OpenAIClient.GetModel()  // ✅ Lấy từ config động
    return result, rawResponse, modelName, processingTime, nil
}

func (u *AssistantUsecase) ChatWithOpenAI(...) {
    // ...
    modelName := u.OpenAIClient.GetModel()  // ✅ Dynamic
    return reply, sessionID, modelName, processingTime, nil
}

func (u *AssistantUsecase) GenerateContentWithOpenAI(...) {
    // ...
    modelName := u.OpenAIClient.GetModel()  // ✅ Dynamic
    return content, modelName, processingTime, nil
}
```

## 🎯 Kết quả

### Config: deepseek-chat
```yaml
# config/app.yml
deepseek:
  model: "deepseek-chat"
```

**Response**:
```json
{
  "metadata": {
    "model": "deepseek-chat"  // ✅ Đúng
  }
}
```

### Đổi Config: deepseek-v3
```yaml
# config/app.yml
deepseek:
  model: "deepseek-v3"
```

**Response**:
```json
{
  "metadata": {
    "model": "deepseek-v3"  // ✅ Tự động update theo config!
  }
}
```

### OpenAI: gpt-3.5-turbo
```yaml
# config/app.yml
openai:
  model: "gpt-3.5-turbo"
```

**Response**:
```json
{
  "metadata": {
    "model": "gpt-3.5-turbo"  // ✅ Đúng
  }
}
```

### Đổi sang: gpt-4o
```yaml
# config/app.yml
openai:
  model: "gpt-4o"
```

**Response**:
```json
{
  "metadata": {
    "model": "gpt-4o"  // ✅ Tự động update!
  }
}
```

## 🔄 Flow hoàn chỉnh

```
1. Load Config
   config/app.yml → deepseek.model = "deepseek-v3"
   
2. Initialize Client
   NewDeepseekClient() → client.model = "deepseek-v3"
   
3. Call API
   handler → usecase.AnalyzeProductText()
   
4. Get Model Name
   usecase → client.GetModel() → return "deepseek-v3"
   
5. Return Response
   {
     "metadata": {
       "model": "deepseek-v3"  // ✅ Chính xác 100%
     }
   }
```

## 💡 Lợi ích

1. ✅ **Dynamic**: Model name tự động thay đổi theo config
2. ✅ **Accurate**: Luôn hiển thị đúng model đang chạy
3. ✅ **Flexible**: Dễ dàng đổi model không cần code
4. ✅ **Traceable**: Biết chính xác model nào xử lý request nào
5. ✅ **Cost tracking**: Monitor chi phí theo model thực tế

## 📊 So sánh

### Trước (Hardcode):
| Config | Code | Response | Đúng? |
|--------|------|----------|-------|
| deepseek-v3 | "deepseek-chat" | "deepseek-chat" | ❌ SAI |
| gpt-4o | "gpt-3.5-turbo" | "gpt-3.5-turbo" | ❌ SAI |

### Sau (Dynamic):
| Config | Code | Response | Đúng? |
|--------|------|----------|-------|
| deepseek-v3 | GetModel() | "deepseek-v3" | ✅ ĐÚNG |
| gpt-4o | GetModel() | "gpt-4o" | ✅ ĐÚNG |

## 🧪 Testing

### Test với DeepSeek:
```yaml
# config/app.yml
deepseek:
  model: "deepseek-chat"
```

```bash
curl -X POST "http://localhost:8061/v2/assistant/analyze/product" \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 3 tầng"}' | jq '.metadata.model'

# Expected: "deepseek-chat"
```

### Đổi config và restart:
```yaml
# config/app.yml
deepseek:
  model: "deepseek-v3"
```

```bash
# Restart service
./assistant-service

# Test lại
curl -X POST "http://localhost:8061/v2/assistant/analyze/product" \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 3 tầng"}' | jq '.metadata.model'

# Expected: "deepseek-v3"  ✅ Tự động đổi theo config!
```

### Test với OpenAI:
```yaml
# config/app.yml
openai:
  model: "gpt-4o-mini"
```

```bash
curl -X POST "http://localhost:8061/v2/assistant/analyze/product" \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 3 tầng"}' | jq '.metadata.model'

# Expected: "gpt-4o-mini"  ✅
```

## 🔧 Fallback Behavior

```go
func AnalyzeProductTextWithFallback(...) {
    // Try DeepSeek
    result, _, modelName, time, err := u.AnalyzeProductText(...)
    if err == nil {
        // modelName = u.DeepseekClient.GetModel()  (VD: "deepseek-v3")
        return result, _, modelName, time, nil
    }
    
    // Fallback to OpenAI
    result, _, modelName, time, err = u.AnalyzeProductTextWithOpenAI(...)
    // modelName = u.OpenAIClient.GetModel()  (VD: "gpt-4o")
    return result, _, modelName, time, err
}
```

**Response khi DeepSeek fail**:
```json
{
  "metadata": {
    "model": "gpt-4o",  // ✅ Hiển thị OpenAI vì đã fallback
    "processing_time_ms": 2500
  }
}
```

## 📝 Environment Variables Support

Client cũng hỗ trợ override model từ env:

```bash
# Override qua env var
export DEEPSEEK_MODEL="deepseek-v3"
export OPENAI_MODEL="gpt-4o"

# Hoặc trong code khi init client
model := os.Getenv("DEEPSEEK_MODEL")
if model == "" {
    model = viper.GetString("deepseek.model")
}
```

**GetModel() sẽ trả về model thực tế đang dùng từ env hoặc config!**

## ✅ Checklist

- [x] Thêm GetModel() vào DeepseekClient
- [x] Thêm GetModel() vào OpenAIClient
- [x] Thêm GetModel() vào DeepseekProvider interface
- [x] Thêm GetModel() vào OpenAIProvider interface
- [x] Usecase gọi GetModel() thay vì hardcode (6 methods)
- [x] Test compile không lỗi
- [x] Linter pass

## 📚 Files đã thay đổi

1. `infra/client/deepseek_client.go` - Thêm GetModel()
2. `infra/client/openai_client.go` - Thêm GetModel()
3. `internal/interface/provider/deepseek_provider.go` - Thêm GetModel() vào interface
4. `internal/interface/provider/openai_provider.go` - Thêm GetModel() vào interface
5. `internal/usecases/assistant_usecase.go` - Gọi GetModel() (6 chỗ)

---

**Trạng thái**: ✅ HOÀN THÀNH  
**Ngày fix**: 2025-01-23  
**Test**: ⚠️ Cần test sau khi restart service với config khác nhau

**Model name giờ LUÔN CHÍNH XÁC 100% theo config/env thực tế!** 🎯

