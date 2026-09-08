# 🤖 Model Training Guide

**AI Service** - Hướng dẫn train/fine-tune model để cải thiện độ chính xác

---

## 📋 TỔNG QUAN

Hệ thống hỗ trợ **fine-tune** sentence-transformers model trên Q&A dataset của BDSPro:

- ✅ **291 Q&A pairs** được generate tự động từ knowledge base
- ✅ **Fine-tuning** với CosineSimilarity loss
- ✅ **API endpoint** để trigger training
- ✅ **Hot-swap** giữa base model và fine-tuned model

---

## 🎯 KHI NÀO CẦN TRAIN?

Train lại model khi:
1. ✅ **Thêm documents mới** vào knowledge base
2. ✅ **Câu trả lời không chính xác** cho một số câu hỏi
3. ✅ **Cập nhật nội dung** trong knowledge base
4. ✅ **Cải thiện độ chính xác** semantic search

---

## 🚀 CÁCH SỬ DỤNG

### Option 1: Training qua API (Recommended) ⭐

```bash
# 1. Trigger training (TỰ ĐỘNG LÀM TẤT CẢ!)
curl -X POST http://localhost:8040/v2/ai/train

# Response (sau 10-20 phút):
{
  "status": "success",
  "message": "Model trained and knowledge base re-embedded successfully",
  "results": {
    "model_path": "./models/fine_tuned",
    "training_examples": 291,
    "total_steps": 54,
    "epochs": 3
  },
  "embedded_chunks": 192,
  "model_activated": true,
  "info": [
    "✅ Model trained successfully",
    "✅ Embedding manager reloaded with fine-tuned model",
    "✅ Knowledge base re-embedded with new model",
    "✅ Service now using fine-tuned model automatically"
  ]
}

# ✅ XONG! Service tự động:
# - Train model
# - Reload model mới
# - Re-embed toàn bộ knowledge base
# - Kích hoạt model ngay lập tức

# 2. Verify (optional)
curl http://localhost:8040/v2/ai/model-info
```

**🎉 Không cần restart service hay chỉnh .env nữa!**

### Option 2: Training qua Script

```bash
cd /Users/baphieu/Workspace/ME/go-microservices/ai-service

# 1. Generate training data (nếu chưa có)
source venv/bin/activate
python scripts/generate_training_data.py

# 2. Train model
python scripts/train_model.py

# 3. Update .env
echo "USE_FINE_TUNED_MODEL=true" >> .env

# 4. Restart service
lsof -ti:8040 | xargs kill -9
HTTP_PORT=8040 python main.py
```

---

## 📊 TRAINING DATA

### Structure

```json
{
  "question": "Làm sao tạo API mới trong BDSPro?",
  "answer": "Để tạo API mới...",
  "source": "BDSPro_Project_Overview.md",
  "section": "API Development",
  "metadata": {
    "level": 2,
    "type": "knowledge_base"
  }
}
```

### Files

- **Input**: `knowledge_base/*.md` (3 files)
- **Training Data**: `data/train/*.json` (TẤT CẢ file JSON trong folder)
  - `person_dataset.json` - 203 Q&A pairs
  - Có thể thêm nhiều file JSON khác
- **Model**: `models/fine_tuned/` (fine-tuned weights)

**💡 Tip**: Training endpoint tự động đọc **TẤT CẢ** file `.json` trong `data/train/` và combine chúng lại!

### Dataset Breakdown

| Source | Sections | Q&A Pairs |
|--------|----------|-----------|
| BDSPro_Project_Overview.md | 77 | 115 |
| Code_Patterns_and_Best_Practices.md | 39 | 60 |
| Service_Architecture_Deep_Dive.md | 73 | 111 |
| Code-specific (manual) | - | 5 |
| **TOTAL** | **189** | **291** |

---

## ⚙️ TRAINING PARAMETERS

```python
{
  "base_model": "all-MiniLM-L6-v2",
  "epochs": 3,
  "batch_size": 16,
  "loss_function": "CosineSimilarityLoss",
  "warmup_steps": 100,
  "output_path": "./models/fine_tuned"
}
```

### Time Estimates

- **Generate Q&A data**: ~5 seconds
- **Training**: ~5-10 minutes (CPU), ~2-3 minutes (GPU)
- **Total**: ~10-15 minutes

---

## 🔄 WORKFLOW (AUTOMATED)

```
┌─────────────────────────────────────────────────────────┐
│  POST /v2/ai/train                                      │
│  (1 API call làm TẤT CẢ!)                              │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│  STEP 1: Load Training Data                            │
│  ├─ Scan data/train/*.json                             │
│  ├─ Load ALL JSON files                                │
│  └─ Combine Q&A pairs                                  │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│  STEP 2: Fine-tune Model                               │
│  ├─ Base: all-MiniLM-L6-v2                            │
│  ├─ Loss: CosineSimilarity                            │
│  ├─ Epochs: 3                                          │
│  └─ Save: models/fine_tuned/                          │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│  STEP 3: Reload Model                                  │
│  ├─ Create new EmbeddingManager                       │
│  ├─ Load fine-tuned model                             │
│  └─ Update global state                               │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│  STEP 4: Re-embed Knowledge Base                       │
│  ├─ Clear old embeddings                              │
│  ├─ Re-embed với model mới                            │
│  ├─ Save to ChromaDB                                  │
│  └─ Update RAG retriever                              │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│  ✅ DONE!                                               │
│  Service sử dụng fine-tuned model ngay lập tức        │
└─────────────────────────────────────────────────────────┘
```

---

## 📈 EXPECTED IMPROVEMENTS

Fine-tuned model vs Base model:

| Metric | Base Model | Fine-tuned |
|--------|------------|------------|
| Domain relevance | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Code understanding | ⭐⭐ | ⭐⭐⭐⭐ |
| Architecture queries | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Terminology matching | ⭐⭐ | ⭐⭐⭐⭐ |
| Context retrieval | ⭐⭐⭐ | ⭐⭐⭐⭐ |

**Estimated improvement**: 20-40% better retrieval accuracy

---

## 🔧 API ENDPOINTS

### 1. Train Model
```bash
POST /v2/ai/train

Response:
{
  "status": "success",
  "message": "Model trained successfully",
  "results": {...},
  "next_steps": [...]
}
```

### 2. Model Info
```bash
GET /v2/ai/model-info

Response:
{
  "model_type": "base" | "fine_tuned",
  "base_model": "all-MiniLM-L6-v2",
  "fine_tuned_available": true/false,
  "vector_db_count": 192
}
```

### 3. Re-embed with Fine-tuned Model
```bash
POST /v2/ai/embed

# Re-embeds knowledge base with current model
# Should run after training to use fine-tuned embeddings
```

---

## 💡 BEST PRACTICES

### 1. **Training Flow (AUTOMATED)** ✨

Training giờ tự động làm TẤT CẢ:
```bash
# Chỉ cần 1 lệnh!
curl -X POST http://localhost:8040/v2/ai/train

# ✅ Tự động train
# ✅ Tự động reload model
# ✅ Tự động re-embed knowledge base
# ✅ Tự động activate ngay lập tức

# KHÔNG CẦN restart service!
# KHÔNG CẦN chỉnh .env!
# KHÔNG CẦN gọi /v2/ai/embed!
```

### 2. **Update Training Data**

Khi thêm documents mới:
```bash
# 1. Add .md files to knowledge_base/

# 2. (Optional) Generate thêm Q&A dataset
python scripts/generate_training_data.py
# Output: data/train/new_dataset.json

# 3. Re-train (đọc TẤT CẢ JSON trong data/train/)
curl -X POST http://localhost:8040/v2/ai/train

# ✅ XONG! Tự động làm hết
```

### 3. **Backup Models**

```bash
# Backup fine-tuned model
cp -r models/fine_tuned models/fine_tuned_backup_$(date +%Y%m%d)

# Restore if needed
cp -r models/fine_tuned_backup_20251018 models/fine_tuned
```

### 4. **A/B Testing**

```bash
# Test with base model
export USE_FINE_TUNED_MODEL=false
python main.py

# Test with fine-tuned
export USE_FINE_TUNED_MODEL=true
python main.py

# Compare results
```

---

## 🐛 TROUBLESHOOTING

### Issue: "Training data not found"

```bash
# Generate training data first
python scripts/generate_training_data.py
```

### Issue: "Out of memory during training"

```python
# Reduce batch size in scripts/train_model.py
batch_size = 8  # Instead of 16
```

### Issue: "Fine-tuned model worse than base"

```bash
# 1. Check training data quality
cat data/train/qa_dataset.json | jq '.[:5]'

# 2. Try more epochs
# Edit scripts/train_model.py: epochs = 5

# 3. Revert to base model
export USE_FINE_TUNED_MODEL=false
```

---

## 📚 FILES STRUCTURE

```
ai-service/
├── data/
│   └── train/
│       ├── person_dataset.json      ← 203 Q&A pairs
│       ├── qa_dataset.json          ← (optional) thêm dataset khác
│       └── *.json                   ← Training đọc TẤT CẢ file JSON
├── models/
│   └── fine_tuned/                  ← Fine-tuned model weights
│       ├── config.json
│       ├── pytorch_model.bin
│       └── ...
├── scripts/
│   ├── generate_training_data.py    ← Generate Q&A dataset
│   ├── train_model.py               ← Train model (manual)
│   └── embed_context.py             ← Embed knowledge base
├── app/core/
│   ├── embeddings.py                ← Support fine-tuned model
│   └── trainer.py                   ← Training logic (auto-load multiple files)
└── main.py
    └── POST /v2/ai/train            ← Auto train + embed + activate
```

---

## 🚀 QUICK COMMANDS

### Simple Way (Recommended) ⭐

```bash
# 1. Chỉ cần 1 lệnh để train + embed + activate!
curl -X POST http://localhost:8040/v2/ai/train

# 2. Test ngay
curl -X POST http://localhost:8040/v2/ai/query \
  -H "Content-Type: application/json" \
  -d '{"question": "Làm sao tạo API mới?"}'

# ✅ DONE! Đơn giản vậy thôi!
```

### Advanced Way (Manual)

```bash
cd /Users/baphieu/Workspace/ME/go-microservices/ai-service
source venv/bin/activate

# Generate training data (nếu cần thêm dataset mới)
python scripts/generate_training_data.py

# Train manually
python scripts/train_model.py

# Enable fine-tuned
echo "USE_FINE_TUNED_MODEL=true" >> .env

# Restart
lsof -ti:8040 | xargs kill -9
HTTP_PORT=8040 python main.py

# Re-embed manually
curl -X POST http://localhost:8040/v2/ai/embed
```

---

## 📊 MONITORING

```bash
# Check model type
curl http://localhost:8040/v2/ai/model-info

# Check stats
curl http://localhost:8040/v2/ai/stats

# Check health
curl http://localhost:8040/v2/ai/health
```

---

**Happy Training! 🎉**

