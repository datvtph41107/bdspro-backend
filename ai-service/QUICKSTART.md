# BDSPro AI Service - Quick Start Guide

**🎯 Mục tiêu**: Chạy AI service trong 5 phút

---

## ⚡ Quick Start (5 Minutes)

```bash
# 1. Navigate
cd ai-service

# 2. Setup Python environment
python3 -m venv venv
source venv/bin/activate  # Mac/Linux
# venv\Scripts\activate   # Windows

# 3. Install dependencies
pip install -r requirements.txt

# 4. Configure
# ai-service/.env đã là file runtime duy nhất; điền OPENAI_API_KEY nếu cần.
# Edit .env: Thêm OPENAI_API_KEY

# 5. Embed context (first time only - ~30 seconds)
python scripts/embed_context.py

# 6. Run service
python main.py

# ✅ Service running at http://localhost:8090
```

---

## 🧪 Test the Service

### 1. Health Check

```bash
curl http://localhost:8090/v2/ai/health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "vector_db_count": 120,
  "uptime_seconds": 45
}
```

### 2. Query Knowledge Base

```bash
curl -X POST http://localhost:8090/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Làm sao tạo API mới trong BDSPro?",
    "top_k": 5
  }'
```

**Response:**
```json
{
  "answer": "Để tạo API mới trong BDSPro, bạn follow 7 bước:\n\n1. **Proto Definition**: Định nghĩa trong `shared/protobuf/schema/<service>/<name>.proto`\n2. **Generate**: Chạy `make buf-<service>`\n3. **Handler**: Implement trong `infra/handler/<name>_handler.go` với Swagger comments và `@bind`\n4. **Usecase**: Viết business logic trong `internal/usecase/<name>_usecase.go`\n5. **Interface**: Định nghĩa interface trong `internal/interface/repo/<name>_repo.go`\n6. **Repository**: Implement trong `infra/postgre/<name>_postgres.go` (embed `_provider.CrudRepo[T]`)\n7. **Mapper**: Tạo mapper trong `infra/mapper/<name>_mapper.go`\n\n**Quan trọng:**\n- ❌ KHÔNG dùng `gin.Context` → ✅ Dùng `context.Context`\n- ✅ Repository phải gọi `repo.Init(repo, db)` trong constructor\n- ❌ KHÔNG cần viết test/wire ngay (sẽ gen sau)",
  "sources": [
    "23_QUICK_REFERENCE_GUIDE.md",
    "07_DEVELOPMENT_WORKFLOW.md"
  ],
  "confidence": 0.92,
  "processing_time_ms": 1234,
  "context_chunks": 5
}
```

### 3. More Examples

```bash
# Repository pattern
curl -X POST http://localhost:8090/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{"question": "Repository pattern viết như thế nào?"}'

# Payment service
curl -X POST http://localhost:8090/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{"question": "Payment service hoạt động ra sao?"}'

# Get user info
curl -X POST http://localhost:8090/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{"question": "Làm sao lấy profileId từ context?"}'
```

---

## 🐳 Docker Quick Start

```bash
# 1. Build
docker-compose build

# 2. Run
docker-compose up -d

# 3. Check logs
docker-compose logs -f ai-service

# 4. Test
curl http://localhost:8090/v2/ai/health
```

---

## 🔧 Integration với Code của Bạn

### Python Client

```python
import requests

def ask_ai(question: str):
    response = requests.post(
        "http://localhost:8090/v2/ai/query",
        json={"question": question, "top_k": 5}
    )
    return response.json()

# Usage
result = ask_ai("Làm sao tạo API mới?")
print(result['answer'])
print(f"Sources: {', '.join(result['sources'])}")
```

### Go Client

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type QueryRequest struct {
    Question string `json:"question"`
    TopK     int    `json:"top_k"`
}

type QueryResponse struct {
    Answer           string   `json:"answer"`
    Sources          []string `json:"sources"`
    Confidence       float64  `json:"confidence"`
    ProcessingTimeMS int      `json:"processing_time_ms"`
}

func AskAI(question string) (*QueryResponse, error) {
    req := QueryRequest{
        Question: question,
        TopK:     5,
    }
    
    body, _ := json.Marshal(req)
    resp, err := http.Post(
        "http://localhost:8090/v2/ai/query",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result QueryResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

// Usage
func main() {
    result, _ := AskAI("Làm sao tạo API mới?")
    fmt.Println(result.Answer)
}
```

### cURL in Terminal

```bash
# Add to ~/.bashrc or ~/.zshrc
askbdspro() {
    curl -s -X POST http://localhost:8090/v2/ai/query \
      -H "Content-Type: application/json" \
      -d "{\"question\": \"$1\"}" | jq -r '.answer'
}

# Usage
askbdspro "Làm sao tạo API mới?"
```

---

## 📊 Performance

```
First query:    ~2 seconds (cold start)
Subsequent:     ~1 second
Accuracy:       90%+
Cost:           ~$0.01 per query
Memory:         ~500MB
CPU:            Low
```

---

## 🐛 Troubleshooting

### Error: "OPENAI_API_KEY not found"

```bash
# Check .env file
cat .env | grep OPENAI_API_KEY

# Or set in environment
export OPENAI_API_KEY=sk-your-key-here
```

### Error: "Vector DB is empty"

```bash
# Re-embed context
python scripts/embed_context.py
```

### Error: "Port 8090 already in use"

```bash
# Change port in .env
echo "HTTP_PORT=8091" >> .env

# Or kill existing process
lsof -ti:8090 | xargs kill -9
```

---

## 🎯 Next Steps

1. ✅ Service running
2. ✅ Test với 5-10 câu hỏi
3. ✅ Integrate vào workflow
4. 📈 Monitor performance
5. 🔄 Update context khi cần

---

## 💡 Pro Tips

### 1. Keep Service Running

```bash
# Add to startup (Mac)
# Create ~/Library/LaunchAgents/bdspro-ai.plist

# Or use tmux/screen
tmux new -s bdspro-ai
cd ai-service && python main.py
# Ctrl+B, D to detach
```

### 2. Quick Re-embed

```bash
# When context files change
curl -X POST http://localhost:8090/v2/ai/embed
```

### 3. Alias for Quick Query

```bash
# Add to .bashrc/.zshrc
alias ask='f(){ curl -s -X POST http://localhost:8090/v2/ai/query -H "Content-Type: application/json" -d "{\"question\": \"$1\"}" | jq -r .answer; }; f'

# Usage
ask "Repository pattern là gì?"
```

---

**🎉 Done! Giờ bạn có AI assistant riêng cho BDSPro!**
