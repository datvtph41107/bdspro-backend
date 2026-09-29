# BDSPro AI Service

**Port**: 50070 (gRPC), 8090 (HTTP)  
**Purpose**: RAG-based AI assistant for BDSPro codebase knowledge

---

## 🎯 Features

- ✅ Query codebase knowledge với natural language
- ✅ Tự động embed context files
- ✅ REST API + gRPC support
- ✅ Real-time response streaming
- ✅ Source tracking & accuracy scoring

---

## 🚀 Quick Start

```bash
# 1. Setup
cd ai-service
python -m venv venv
source venv/bin/activate  # Linux/Mac

# 2. Install
pip install -r requirements.txt

# 3. Configure service owner
# Edit ai-service/.env: Add OPENAI_API_KEY khi cần chạy AI thật

# 4. Embed context (first time only)
python scripts/embed_context.py

# 5. Run service
python main.py

# Service running at:
# - HTTP: http://localhost:8090
# - gRPC: localhost:50070
```

---

## 📡 API Usage

### HTTP REST API

```bash
# Query
curl -X POST http://localhost:8090/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Làm sao tạo API mới?",
    "top_k": 5
  }'

# Response
{
  "answer": "Để tạo API mới, follow 7 bước:\n1. Proto Definition...",
  "sources": ["23_QUICK_REFERENCE_GUIDE.md", "07_DEVELOPMENT_WORKFLOW.md"],
  "confidence": 0.92,
  "processing_time_ms": 1234
}

# Health check
curl http://localhost:8090/v2/ai/health

# Stats
curl http://localhost:8090/v2/ai/stats
```

### gRPC API

```protobuf
service AIService {
  rpc Query(QueryRequest) returns (QueryResponse);
  rpc Health(HealthRequest) returns (HealthResponse);
}
```

---

## 🛠️ Development

### Project Structure

```
ai-service/
├── main.py              # FastAPI app entry
├── config.py            # Configuration
├── requirements.txt     # Dependencies
├── .env                 # Runtime environment của AI service
├── app/
│   ├── api/
│   │   ├── http/        # REST endpoints
│   │   └── grpc/        # gRPC handlers
│   ├── core/
│   │   ├── embeddings.py   # Embedding logic
│   │   ├── retrieval.py    # RAG retrieval
│   │   └── llm.py          # LLM interface
│   └── models/
│       └── schemas.py   # Request/Response models
├── scripts/
│   └── embed_context.py # Embed context files
├── data/
│   └── chroma_db/       # Vector database
└── tests/
    └── test_api.py      # Tests
```

---

## 🐳 Docker Deployment

```bash
# Build
docker build -t bdspro-ai-service:latest .

# Run
docker run -d \
  -p 8090:8090 \
  -p 50070:50070 \
  -e OPENAI_API_KEY=sk-xxx \
  -v $(pwd)/data:/app/data \
  --name bdspro-ai \
  bdspro-ai-service:latest
```

---

## 🔧 Configuration

### Environment Variables

```bash
# ai-service/.env
OPENAI_API_KEY=sk-your-key-here
EMBEDDING_MODEL=text-embedding-3-small
LLM_MODEL=gpt-4-turbo-preview
CHROMA_DB_PATH=./data/chroma_db
CONTEXT_PATH=../context
HTTP_PORT=8090
GRPC_PORT=50070
```

---

## 📊 Performance

```
Query latency:      1-2 seconds
Embedding (init):   30 seconds (one-time)
Concurrent users:   100+
Memory usage:       ~500MB
CPU usage:          Low (unless embedding)
```

---

## 🧪 Testing

```bash
# Unit tests
pytest tests/

# Integration test
python scripts/test_queries.py

# Load test
locust -f tests/load_test.py
```

---

## 📈 Monitoring

```bash
# Prometheus metrics at /metrics
# - ai_query_total
# - ai_query_duration_seconds
# - ai_embedding_total
# - ai_errors_total

# Logs
tail -f logs/ai-service.log
```

---

## 🔐 Security

- API key validation
- Rate limiting (100 req/min per IP)
- Input sanitization
- No sensitive data in embeddings

---

## 💰 Cost Estimation

```
OpenAI Costs:
- Embedding (one-time): ~$0.002
- Query (GPT-4): ~$0.01/query
- Monthly (100 queries/day): ~$30
```
