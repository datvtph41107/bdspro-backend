# ⚡ Quick Start - Training System

**TL;DR**: Hệ thống training đã sẵn sàng. Dữ liệu **CHỈ** nằm trong `data/train/`.

---

## 🎯 ĐIỂM QUAN TRỌNG

✅ **Dữ liệu training**: CHỈ từ `data/train/qa_dataset.json`  
✅ **Không cần** `knowledge_base/` để train  
✅ **291 Q&A pairs** đã có sẵn  
✅ **API endpoint** `/v2/ai/train` sẵn sàng

---

## 📂 CẤU TRÚC DATA

```
data/train/
├── README.md           ← Hướng dẫn thêm Q&A
└── qa_dataset.json     ← 291 Q&A pairs (READY!)
```

---

## 🚀 CÁCH DÙNG

### Option 1: Via API (Recommended)

```bash
# Train model
curl -X POST http://localhost:8040/v2/ai/train

# Check model info
curl http://localhost:8040/v2/ai/model-info
```

### Option 2: Via Script Helper

```bash
cd /Users/baphieu/Workspace/ME/go-microservices/ai-service

# Full flow
./train.sh full

# Or step by step:
./train.sh train     # Train model
./train.sh enable    # Enable fine-tuned model
./train.sh restart   # Restart service
./train.sh test      # Test query
```

### Option 3: Manual

```bash
# 1. Train
python scripts/train_model.py

# 2. Enable
echo "USE_FINE_TUNED_MODEL=true" >> .env

# 3. Restart
lsof -ti:8040 | xargs kill -9
HTTP_PORT=8040 python main.py
```

---

## ➕ THÊM Q&A MỚI

### Cách 1: Edit trực tiếp

```bash
nano data/train/qa_dataset.json

# Thêm vào array:
{
  "question": "Câu hỏi mới?",
  "answer": "Câu trả lời chi tiết...",
  "source": "manual"
}
```

### Cách 2: Python

```python
import json

with open('data/train/qa_dataset.json', 'r') as f:
    data = json.load(f)

data.append({
    "question": "New question?",
    "answer": "Detailed answer...",
    "source": "manual"
})

with open('data/train/qa_dataset.json', 'w') as f:
    json.dump(data, f, ensure_ascii=False, indent=2)
```

---

## ⏱️ THỜI GIAN

- **Generate data**: Đã có sẵn (không cần)
- **Training**: ~5-10 phút (CPU), ~2-3 phút (GPU)
- **Re-embed**: ~30 giây
- **Total**: ~10 phút

---

## 🔄 WORKFLOW

```
data/train/qa_dataset.json (291 pairs)
         ↓
   Train model
         ↓
models/fine_tuned/
         ↓
Enable fine-tuned model
         ↓
Restart service
         ↓
Better search results!
```

---

## 📋 DEPENDENCIES

Đã có trong `requirements.txt`:

```
sentence-transformers>=3.3.0
datasets>=4.2.0
accelerate>=0.26.0
```

Install:
```bash
pip install -r requirements.txt
```

---

## 🎨 HELPER SCRIPTS

### `./train.sh` Commands

```bash
./train.sh train      # Train model
./train.sh enable     # Enable fine-tuned
./train.sh disable    # Disable (use base)
./train.sh restart    # Restart service
./train.sh reembed    # Re-embed knowledge base
./train.sh info       # Show model info
./train.sh test       # Test query
./train.sh full       # Run full pipeline
```

---

## 🔍 CHECK STATUS

```bash
# Model info
curl http://localhost:8040/v2/ai/model-info

# Output:
{
  "model_type": "base" | "fine_tuned",
  "base_model": "all-MiniLM-L6-v2",
  "fine_tuned_available": true/false,
  "vector_db_count": 192
}
```

---

## 💡 NOTES

### Training đang chạy trong background?
```bash
# Check process
ps aux | grep train_model

# Kill if needed
pkill -f train_model
```

### Training xong làm gì?
```bash
# 1. Enable fine-tuned model
./train.sh enable

# 2. Restart service
./train.sh restart

# 3. Re-embed để dùng fine-tuned embeddings
curl -X POST http://localhost:8040/v2/ai/embed

# 4. Test
./train.sh test "Làm sao tạo API mới?"
```

### Backup model
```bash
cp -r models/fine_tuned models/backup_$(date +%Y%m%d)
```

---

## 🚨 TROUBLESHOOTING

### "Training data not found"
```bash
# Check file exists
ls -la data/train/qa_dataset.json

# If not, generate:
python scripts/generate_training_data.py
```

### "Out of memory"
```bash
# Edit scripts/train_model.py
batch_size = 8  # Reduce from 16
```

### Training quá lâu
```bash
# Reduce epochs in scripts/train_model.py
epochs = 1  # Instead of 3
```

---

## 📚 CHI TIẾT HƠN

- **Full guide**: `TRAINING_GUIDE.md`
- **Data README**: `data/train/README.md`
- **Local mode**: `LOCAL_MODE_README.md`

---

**Ready to train! 🚀**

```bash
./train.sh full
```

