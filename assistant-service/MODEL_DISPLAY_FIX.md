# ✅ Fix hiển thị Model Name - Summary

## 🐛 Vấn đề ban đầu

Handler đang **HARDCODE** model name:
```go
// Line 134 - handler
Model: "deepseek-chat",  // ❌ SAI! Đang dùng OpenAI mà hiển thị deepseek
```

Trong khi thực tế:
- Line 44: Gọi `AnalyzeProductTextWithOpenAI` → dùng **gpt-3.5-turbo**
- Line 164: Gọi `Chat` → dùng **deepseek-chat**
- Line 206: Gọi `GenerateContent` → dùng **deepseek-chat**

## ✅ Giải pháp đã thực hiện

### 1. Cập nhật Usecase - Trả về Model Name

**Tất cả methods trong usecase giờ trả về `modelName`**:

#### Methods DeepSeek:
```go
// Before: (result, rawResponse, processingTime, error)
// After:  (result, rawResponse, modelName, processingTime, error)

func AnalyzeProductText(...) (*dto.ProductAnalysisResult, string, string, int64, error) {
    // ...
    modelName := "deepseek-chat"
    return result, rawResponse, modelName, processingTime, nil
}

func Chat(...) (string, string, string, int64, error) {
    // ...
    modelName := "deepseek-chat"
    return reply, sessionID, modelName, processingTime, nil
}

func GenerateContent(...) (string, string, int64, error) {
    // ...
    modelName := "deepseek-chat"
    return content, modelName, processingTime, nil
}
```

#### Methods OpenAI:
```go
func AnalyzeProductTextWithOpenAI(...) (*dto.ProductAnalysisResult, string, string, int64, error) {
    // ...
    modelName := "gpt-3.5-turbo"  // ✅ Hiển thị đúng OpenAI model
    return result, rawResponse, modelName, processingTime, nil
}

func ChatWithOpenAI(...) (string, string, string, int64, error) {
    // ...
    modelName := "gpt-3.5-turbo"
    return reply, sessionID, modelName, processingTime, nil
}

func GenerateContentWithOpenAI(...) (string, string, int64, error) {
    // ...
    modelName := "gpt-3.5-turbo"
    return content, modelName, processingTime, nil
}
```

#### Fallback Method:
```go
func AnalyzeProductTextWithFallback(...) (*dto.ProductAnalysisResult, string, string, int64, error) {
    // Try DeepSeek first
    result, rawResponse, modelName, processingTime, err := u.AnalyzeProductText(...)
    if err == nil {
        return result, rawResponse, modelName, processingTime, nil  // modelName = "deepseek-chat"
    }
    
    // Fallback to OpenAI
    result, rawResponse, modelName, processingTime, err = u.AnalyzeProductTextWithOpenAI(...)
    // modelName = "gpt-3.5-turbo"
    return result, rawResponse, modelName, processingTime, err
}
```

### 2. Cập nhật Handler - Sử dụng Model Name động

#### AnalyzeProductText Handler:
```go
// Before
result, rawResponse, processingTime, err := h.Usecase.AnalyzeProductTextWithOpenAI(...)
Metadata: &assistantpb.AnalysisMetadata{
    Model: "deepseek-chat",  // ❌ Hardcode sai
}

// After
result, rawResponse, modelName, processingTime, err := h.Usecase.AnalyzeProductTextWithOpenAI(...)
Metadata: &assistantpb.AnalysisMetadata{
    Model: modelName,  // ✅ Hiển thị đúng "gpt-3.5-turbo"
}
```

#### Chat Handler:
```go
// Before
reply, sessionID, processingTime, err := h.Usecase.Chat(...)
Metadata: &assistantpb.ChatMetadata{
    Model: "deepseek-chat",  // ❌ Hardcode
}

// After
reply, sessionID, modelName, processingTime, err := h.Usecase.Chat(...)
Metadata: &assistantpb.ChatMetadata{
    Model: modelName,  // ✅ Hiển thị động "deepseek-chat"
}
```

#### GenerateContent Handler:
```go
// Before
content, processingTime, err := h.Usecase.GenerateContent(...)
Metadata: &assistantpb.GenerateMetadata{
    Model: "deepseek-chat",  // ❌ Hardcode
}

// After
content, modelName, processingTime, err := h.Usecase.GenerateContent(...)
Metadata: &assistantpb.GenerateMetadata{
    Model: modelName,  // ✅ Hiển thị động
}
```

## 📊 Kết quả

### Request sử dụng DeepSeek:
```json
{
  "metadata": {
    "model": "deepseek-chat",          // ✅ ĐÚNG
    "processing_time_ms": 1200
  }
}
```

### Request sử dụng OpenAI:
```json
{
  "metadata": {
    "model": "gpt-3.5-turbo",          // ✅ ĐÚNG
    "processing_time_ms": 1800
  }
}
```

### Request với Fallback (DeepSeek fail → OpenAI):
```json
{
  "metadata": {
    "model": "gpt-3.5-turbo",          // ✅ ĐÚNG - Hiển thị model đã dùng thực tế
    "processing_time_ms": 2500
  }
}
```

## 🎯 Lợi ích

1. ✅ **Hiển thị chính xác** model đang được sử dụng
2. ✅ **Dễ debug**: Biết request nào dùng model gì
3. ✅ **Dễ monitor**: Track usage của từng model
4. ✅ **Transparent**: User biết mình đang dùng AI nào
5. ✅ **Cost tracking**: Biết chi phí từ model nào

## 📝 Breaking Changes

**Có** - Signature của usecase methods đã thay đổi:

### Before:
```go
AnalyzeProductText(...) (*Result, string, int64, error)
Chat(...) (string, string, int64, error)
GenerateContent(...) (string, int64, error)
```

### After:
```go
AnalyzeProductText(...) (*Result, string, string, int64, error)  // +modelName
Chat(...) (string, string, string, int64, error)                 // +modelName
GenerateContent(...) (string, string, int64, error)              // +modelName
```

⚠️ **Nếu có code khác gọi trực tiếp usecase methods, cần update!**

## 🧪 Testing

### Test hiển thị đúng model:

```bash
# Test với API
curl -X POST "http://localhost:8061/v2/assistant/analyze/product" \
  -H "Content-Type: application/json" \
  -d '{"content": "Nhà 3 tầng Cầu Giấy"}' | jq '.metadata.model'

# Expected output: "gpt-3.5-turbo"  (vì đang dùng OpenAI)
```

### Test với DeepSeek:
```go
// Đổi handler line 44 thành:
result, rawResponse, modelName, processingTime, err := h.Usecase.AnalyzeProductText(...)

// Response sẽ hiển thị:
// "model": "deepseek-chat"
```

## 🔄 Tương lai - Dynamic Model từ Config

Hiện tại model name vẫn hardcode trong usecase. Để dynamic 100%:

```go
type DeepseekClient struct {
    model string  // Lấy từ config
}

func (u *AssistantUsecase) AnalyzeProductText(...) {
    // ...
    modelName := u.DeepseekClient.GetModelName()  // Từ config
    return result, rawResponse, modelName, processingTime, nil
}
```

## ✅ Checklist

- [x] Cập nhật usecase return modelName
- [x] Cập nhật handler nhận modelName
- [x] Fix tất cả methods (6 methods)
- [x] Test compile không lỗi
- [x] Linter pass
- [x] Update documentation

## 📚 Files đã thay đổi

1. `internal/usecases/assistant_usecase.go` - 7 methods updated
2. `infra/handler/assistant_handler.go` - 3 handlers updated

---

**Trạng thái**: ✅ HOÀN THÀNH  
**Ngày fix**: 2025-01-23  
**Test**: ⚠️ Cần test sau khi deploy

