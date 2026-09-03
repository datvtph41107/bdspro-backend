# 💬 Hướng Dẫn Cho Conversations Mới

**📍 File này dành cho AI ở conversations sau**

---

## 🎯 CONTEXT NHANH

### Project này là gì?

**BDSPro AI Service** - RAG-based AI assistant đã được implement đầy đủ:
- ✅ FastAPI server (Python)
- ✅ Vector database (ChromaDB)
- ✅ OpenAI GPT-4 integration
- ✅ REST API endpoints
- ✅ Docker deployment
- ✅ Documentation hoàn chỉnh

### Files quan trọng:

```
ai-service/
├── SUMMARY.md              ← ⭐ ĐỌC FILE NÀY TRƯỚC!
├── README.md               ← API docs & full guide
├── QUICKSTART.md           ← Setup trong 5 phút
├── main.py                 ← FastAPI server
├── app/core/               ← Core components
└── scripts/embed_context.py ← Embedding script
```

---

## 📖 BƯỚC 1: ĐỌC SUMMARY.MD

**QUAN TRỌNG**: Đọc `/ai-service/SUMMARY.md` để hiểu:

✅ Architecture & how RAG works  
✅ What has been built (complete file list)  
✅ Technical details of each component  
✅ How to use & integrate  
✅ What's done vs not done  
✅ How to extend  

**→ SUMMARY.md có TẤT CẢ context cần thiết!**

---

## 🚀 BƯỚC 2: HIỂU TRẠNG THÁI

### ✅ Đã Hoàn Thành (Ready to Use)

- [x] FastAPI server với 5 endpoints
- [x] Embedding manager (OpenAI)
- [x] RAG retriever (semantic search)
- [x] LLM client (GPT-4)
- [x] Context embedding script
- [x] Docker deployment files
- [x] Complete documentation

### ❌ Chưa Implement (Future Enhancements)

- [ ] Streaming responses
- [ ] Caching layer (Redis)
- [ ] gRPC server
- [ ] Auto re-embedding
- [ ] Authentication
- [ ] Rate limiting
- [ ] Analytics dashboard

---

## 💡 BƯỚC 3: COMMON TASKS

### Task 1: "Chạy service"

```bash
cd ai-service
python3 -m venv venv && source venv/bin/activate
pip install -r requirements.txt
# Dùng ai-service/.env đã có sẵn; chỉ inject OPENAI_API_KEY khi cần gọi provider thật.
python scripts/embed_context.py  # First time only
python main.py  # Service at :8090
```

### Task 2: "Test service"

```bash
curl http://localhost:8090/v2/ai/health
curl -X POST http://localhost:8090/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{"question": "Làm sao tạo API mới?"}'
```

### Task 3: "Add new feature"

**Example**: Add authentication

```python
# In main.py
from fastapi import Depends, HTTPException
from fastapi.security import HTTPBearer

security = HTTPBearer()

@app.post("/v2/ai/query")
async def query(
    request: QueryRequest,
    token: str = Depends(security)
):
    if not validate_token(token):
        raise HTTPException(401)
    # ... existing logic
```

### Task 4: "Deploy with Docker"

```bash
cd ai-service
docker-compose build
docker-compose up -d
docker-compose logs -f ai-service
```

### Task 5: "Update context files"

```bash
# When ../context/*.md changes
python scripts/embed_context.py  # Re-embed
# Or via API:
curl -X POST http://localhost:8090/v2/ai/embed
```

---

## 🎯 BƯỚC 4: XÁC ĐỊNH TASK

### User muốn gì?

**Scenario A**: "Chạy service"
→ Follow setup guide (QUICKSTART.md hoặc SUMMARY.md section "HOW TO USE")

**Scenario B**: "Thêm feature X"
→ Check SUMMARY.md section "FUTURE ENHANCEMENTS" xem đã plan chưa
→ Nếu chưa plan → Design new feature
→ Implement theo pattern hiện có

**Scenario C**: "Deploy to production"
→ Use Docker deployment (docker-compose.yml)
→ Add monitoring, load balancer, etc.

**Scenario D**: "Fix bug/issue"
→ Check SUMMARY.md section "MAINTENANCE" → Troubleshooting
→ Check logs: `tail -f logs/ai-service.log`

**Scenario E**: "Hiểu code"
→ Read SUMMARY.md section "TECHNICAL DETAILS"
→ Read specific file in `app/core/`

---

## 📋 QUICK CHECKLIST

Trước khi bắt đầu, check:

- [ ] Đã đọc SUMMARY.md chưa? ⭐
- [ ] Service files tồn tại trong `/ai-service/`?
- [ ] Hiểu rõ what's done vs not done?
- [ ] Biết user muốn làm gì?
- [ ] Có cần setup môi trường không?
- [ ] Có cần test service không?

---

## 🎨 CODE EXAMPLES

### Read a component:

```python
# app/core/embeddings.py - Embedding manager
# app/core/retrieval.py - RAG retriever
# app/core/llm.py - LLM client
# main.py - FastAPI server
# scripts/embed_context.py - Embedding script
```

### Key patterns:

```python
# 1. Embedding
embedding = openai.embeddings.create(
    model="text-embedding-3-small",
    input=text
)

# 2. Vector search
results = collection.query(
    query_embeddings=[embedding],
    n_results=5
)

# 3. LLM generation
response = await openai.chat.completions.create(
    model="gpt-4-turbo-preview",
    messages=[
        {"role": "system", "content": system_prompt},
        {"role": "user", "content": user_prompt}
    ]
)

# 4. FastAPI endpoint
@app.post("/v2/ai/query")
async def query(request: QueryRequest) -> QueryResponse:
    # Retrieve context
    # Generate answer
    # Return response
```

---

## 🔧 TROUBLESHOOTING GUIDE

### Issue: "File not found"

**Fix**: Check working directory
```bash
pwd  # Should be in /ai-service/
ls -la  # Should see main.py, requirements.txt
```

### Issue: "Import error"

**Fix**: Activate venv & install deps
```bash
source venv/bin/activate
pip install -r requirements.txt
```

### Issue: "API key error"

**Fix**: Check .env file
```bash
cat .env | grep OPENAI_API_KEY
# Should show: OPENAI_API_KEY=sk-...
```

### Issue: "Empty vector DB"

**Fix**: Re-embed context
```bash
python scripts/embed_context.py
```

### Issue: "Port in use"

**Fix**: Kill process or change port
```bash
lsof -ti:8090 | xargs kill -9
# Or change in .env: HTTP_PORT=8091
```

---

## 💬 RESPONSE TEMPLATE

Khi user hỏi về AI service:

```
✅ Tôi đã đọc SUMMARY.md và hiểu project.

📋 Status: 
- Service đã implement đầy đủ với RAG
- FastAPI server ready
- Docker deployment included
- [List relevant features đã có]

🎯 Về task của bạn "[user's request]":
[Analysis & approach]

📝 Next steps:
1. [Step 1]
2. [Step 2]
...

💡 Notes:
[Any important considerations]
```

---

## 🎉 SUMMARY

**3 Điều Quan Trọng Nhất**:

1. **ĐỌC SUMMARY.MD TRƯỚC** ← ⭐ Most important!
2. **Service đã HOÀN THÀNH** - ready to use/extend
3. **Context đầy đủ** - không cần hỏi thêm về architecture

**Khi không chắc**:
→ Refer to SUMMARY.md section tương ứng
→ Read specific code file if needed
→ Check QUICKSTART.md for setup

---

**Good luck! 🚀**

*P/S: File này được tạo để giúp AI conversations sau có thể continue seamlessly. Nếu bạn là AI đang đọc này → Perfect! You have all context needed!*
