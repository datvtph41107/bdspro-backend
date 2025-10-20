# 🚀 AI Service - LOCAL MODE

**Status**: ✅ Running on port **8040**  
**Mode**: 100% Local (No external API calls)

---

## 📋 THÔNG TIN NHANH

### Service đã được cấu hình:
- ✅ **Embedding**: sentence-transformers (model: all-MiniLM-L6-v2)
- ✅ **Vector DB**: ChromaDB với 192 document chunks
- ✅ **Response**: Direct context retrieval (no LLM)
- ✅ **Port**: 8040
- ✅ **Knowledge Base**: 3 markdown files trong `knowledge_base/`

---

## 🎯 CÁCH SỬ DỤNG

### 1️⃣ Khởi động Service

```bash
cd /Users/baphieu/Workspace/ME/go-microservices/ai-service
source venv/bin/activate
HTTP_PORT=8040 python main.py
```

### 2️⃣ Test API

#### Health Check:
```bash
curl http://localhost:8040/v2/ai/health
```

#### Query (Hỏi đáp):
```bash
curl -X POST http://localhost:8040/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Làm sao tạo API mới trong BDSPro?",
    "top_k": 3
  }'
```

#### Stats:
```bash
curl http://localhost:8040/v2/ai/stats
```

### 3️⃣ Update Knowledge Base

Khi thêm/sửa file trong `knowledge_base/`:

```bash
cd /Users/baphieu/Workspace/ME/go-microservices/ai-service
source venv/bin/activate
python scripts/embed_context.py

# Sau đó restart service
lsof -ti:8040 | xargs kill -9
HTTP_PORT=8040 python main.py
```

---

## 📊 PERFORMANCE

- **Embedding model**: 80MB (local, fast)
- **Vector search**: ~50-100ms
- **Total response time**: ~200-400ms
- **Memory**: ~500MB
- **No API costs** ✨

---

## 🔧 COMPONENTS

### Modified Files:
1. `app/core/embeddings.py` - Dùng sentence-transformers
2. `app/core/llm.py` - Simple context return (no GPT)
3. `main.py` - Không cần OPENAI_API_KEY
4. `scripts/embed_context.py` - Local embedding

### Data:
- Knowledge base: `knowledge_base/*.md`
- Vector DB: `data/chroma_db/`
- Chunks: 192 embeddings

---

## 🎨 EXAMPLE RESPONSE

```json
{
  "answer": "**🔍 Thông tin tìm được trong tài liệu:**\n\n[Context từ docs]...",
  "sources": ["Service_Architecture_Deep_Dive.md", "BDSPro_Project_Overview.md"],
  "confidence": 0.34,
  "processing_time_ms": 237,
  "context_chunks": 5
}
```

---

## 💡 NOTES

### Local Mode:
- ✅ Không cần OPENAI_API_KEY
- ✅ Không có API costs
- ✅ Chạy offline hoàn toàn
- ⚠️ Response là context gốc (chưa qua LLM)

### Để có AI response tự nhiên hơn:
Có thể integrate với:
- Local LLM: Ollama, LLaMA, etc.
- Cloud LLM: GPT-4, Claude (cần API key)
- Custom fine-tuned models

---

## 🚀 QUICK COMMANDS

```bash
# Start service
cd ai-service && source venv/bin/activate && HTTP_PORT=8040 python main.py

# Stop service
lsof -ti:8040 | xargs kill -9

# Re-embed knowledge base
python scripts/embed_context.py

# Test query
curl -X POST http://localhost:8040/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{"question": "Làm sao tạo API mới?"}'
```

---

**Ready to use! 🎉**

