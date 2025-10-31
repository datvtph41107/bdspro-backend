# Hướng dẫn sử dụng OpenAI và DeepSeek trong Assistant Service

## Tổng quan

Assistant Service hỗ trợ 2 AI providers:
- **DeepSeek**: Model tốc độ cao, chi phí thấp, tốt cho tiếng Việt
- **OpenAI**: GPT-4o, GPT-3.5 - chất lượng cao, đa dạng use cases

## Cấu hình

### 1. Cấu hình trong `config/app.yml`

```yaml
# DeepSeek Configuration
deepseek:
  api_key: "sk-e55042ea5c144ac1ada1829ed11595f9"
  base_url: "https://api.deepseek.com/v1"
  model: "deepseek-v3"
  timeout_seconds: 30
  max_tokens: 4000

# OpenAI Configuration  
openai:
  api_key: ""  # Set via OPENAI_API_KEY env variable
  base_url: "https://api.openai.com/v1"
  model: "gpt-4o-mini"
  timeout_seconds: 60
  max_tokens: 4000
```

### 2. Environment Variables

```bash
# Required: Set at least one API key
export DEEPSEEK_API_KEY="sk-your-deepseek-key"
export OPENAI_API_KEY="sk-your-openai-key"
```

## Sử dụng trong Code

### 1. DeepSeek Client

```go
package main

import (
    "context"
    "assistant/infra/client"
    "assistant/internal/dto"
)

func main() {
    // Initialize DeepSeek client
    deepseekClient := client.NewDeepseekClient()
    
    ctx := context.Background()
    
    // Example 1: Analyze Product Text
    productText := `Cần bán nhà 3 tầng, 4 phòng ngủ, 3WC, diện tích 80m2, 
    mặt tiền 5m, hướng Đông, giá 5 tỷ, sổ đỏ chính chủ, 
    tại Cầu Giấy, Hà Nội`
    
    result, err := deepseekClient.AnalyzeProductText(ctx, productText)
    if err != nil {
        panic(err)
    }
    
    // result contains extracted product information
    
    // Example 2: Chat
    message := "Giá nhà ở Cầu Giấy như thế nào?"
    history := []dto.DeepseekMessage{
        {Role: "user", Content: "Xin chào"},
        {Role: "assistant", Content: "Chào bạn! Tôi có thể giúp gì?"},
    }
    
    reply, err := deepseekClient.Chat(ctx, message, history)
    if err != nil {
        panic(err)
    }
    
    // Example 3: Generate Content
    prompt := "Viết mô tả sản phẩm cho nhà 3 tầng ở Cầu Giấy"
    content, err := deepseekClient.GenerateContent(ctx, prompt, 2000, 0.7)
    if err != nil {
        panic(err)
    }
}
```

### 2. OpenAI Client

```go
package main

import (
    "context"
    "assistant/infra/client"
    "assistant/internal/dto"
)

func main() {
    // Initialize OpenAI client
    openaiClient := client.NewOpenAIClient()
    
    ctx := context.Background()
    
    // Example 1: Analyze Product Text
    productText := `Cần bán nhà 3 tầng, 4 phòng ngủ, 3WC, diện tích 80m2, 
    mặt tiền 5m, hướng Đông, giá 5 tỷ, sổ đỏ chính chủ, 
    tại Cầu Giấy, Hà Nội`
    
    result, err := openaiClient.AnalyzeProductText(ctx, productText)
    if err != nil {
        panic(err)
    }
    
    // result contains extracted product information
    
    // Example 2: Chat
    message := "Giá nhà ở Cầu Giấy như thế nào?"
    history := []dto.DeepseekMessage{
        {Role: "user", Content: "Xin chào"},
        {Role: "assistant", Content: "Chào bạn! Tôi có thể giúp gì?"},
    }
    
    reply, err := openaiClient.Chat(ctx, message, history)
    if err != nil {
        panic(err)
    }
    
    // Example 3: Generate Content with custom parameters
    prompt := "Viết mô tả sản phẩm chuyên nghiệp cho nhà cao cấp"
    content, err := openaiClient.GenerateContent(ctx, prompt, 3000, 0.8)
    if err != nil {
        panic(err)
    }
}
```

## Dependency Injection với Wire

### wire.go

```go
//go:build wireinject
// +build wireinject

package wire

import (
    "assistant/infra/client"
    "assistant/internal/interface/provider"
    "github.com/google/wire"
)

func InitializeDeepseekProvider() provider.DeepseekProvider {
    wire.Build(
        client.NewDeepseekClient,
        wire.Bind(new(provider.DeepseekProvider), new(*client.DeepseekClient)),
    )
    return nil
}

func InitializeOpenAIProvider() provider.OpenAIProvider {
    wire.Build(
        client.NewOpenAIClient,
        wire.Bind(new(provider.OpenAIProvider), new(*client.OpenAIClient)),
    )
    return nil
}
```

## So sánh DeepSeek vs OpenAI

| Feature | DeepSeek | OpenAI |
|---------|----------|--------|
| **Model** | deepseek-v3 | gpt-4o-mini, gpt-4o |
| **Chi phí** | Rất thấp (~$0.001/request) | Cao hơn (~$0.002-0.01/request) |
| **Tốc độ** | Rất nhanh (1-2s) | Trung bình (2-4s) |
| **Tiếng Việt** | Tốt | Rất tốt |
| **Độ chính xác** | Tốt | Rất tốt |
| **Use case** | Chat, phân tích cơ bản | Phân tích phức tạp, creative writing |

## Lựa chọn Provider

### Khi nào dùng DeepSeek:
- ✅ Phân tích văn bản sản phẩm BĐS cơ bản
- ✅ Chatbot hỗ trợ khách hàng
- ✅ Generate mô tả sản phẩm đơn giản
- ✅ Cần tốc độ cao, chi phí thấp
- ✅ Xử lý tiếng Việt

### Khi nào dùng OpenAI:
- ✅ Phân tích phức tạp, multi-step reasoning
- ✅ Creative content generation (marketing copy, blog posts)
- ✅ Summarization dài
- ✅ Cần độ chính xác cao nhất
- ✅ Multi-language support

## Xử lý lỗi

```go
// DeepSeek Error Handling
result, err := deepseekClient.AnalyzeProductText(ctx, text)
if err != nil {
    // Check error types
    if strings.Contains(err.Error(), "timeout") {
        // Handle timeout - retry or fallback to OpenAI
    } else if strings.Contains(err.Error(), "API error") {
        // Handle API error - check API key, quota
    }
    return err
}

// OpenAI Error Handling
result, err := openaiClient.AnalyzeProductText(ctx, text)
if err != nil {
    // Similar error handling
    if strings.Contains(err.Error(), "rate limit") {
        // Handle rate limit - wait or fallback to DeepSeek
    }
    return err
}
```

## Fallback Strategy

```go
func AnalyzeWithFallback(ctx context.Context, text string) (*dto.ProductAnalysisResult, error) {
    // Try DeepSeek first (faster, cheaper)
    deepseekClient := client.NewDeepseekClient()
    result, err := deepseekClient.AnalyzeProductText(ctx, text)
    if err == nil {
        return result, nil
    }
    
    // Fallback to OpenAI if DeepSeek fails
    log.Printf("DeepSeek failed: %v, falling back to OpenAI", err)
    openaiClient := client.NewOpenAIClient()
    result, err = openaiClient.AnalyzeProductText(ctx, text)
    if err != nil {
        return nil, fmt.Errorf("both providers failed: deepseek=%v, openai=%v", err, err)
    }
    
    return result, nil
}
```

## Best Practices

1. **API Key Management**
   - Luôn dùng environment variables
   - Không commit API keys vào git
   - Rotate keys định kỳ

2. **Timeout Configuration**
   - DeepSeek: 30s (nhanh)
   - OpenAI: 60s (chậm hơn)
   - Tăng timeout cho requests phức tạp

3. **Cost Optimization**
   - Dùng DeepSeek cho use cases đơn giản
   - Dùng OpenAI khi cần chất lượng cao nhất
   - Cache kết quả khi có thể

4. **Error Handling**
   - Implement retry logic với exponential backoff
   - Log errors để debug
   - Có fallback provider

5. **Testing**
   - Mock providers trong unit tests
   - Test với cả 2 providers trong integration tests
   - Monitor response time và quality

## Troubleshooting

### DeepSeek

```bash
# Check API key
curl https://api.deepseek.com/v1/chat/completions \
  -H "Authorization: Bearer $DEEPSEEK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"deepseek-v3","messages":[{"role":"user","content":"test"}]}'
```

### OpenAI

```bash
# Check API key
curl https://api.openai.com/v1/chat/completions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"test"}]}'
```

## References

- [DeepSeek API Docs](https://platform.deepseek.com/docs)
- [OpenAI API Docs](https://platform.openai.com/docs)
- [Assistant Service README](./README.md)

