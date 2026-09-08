#!/bin/bash
# AI Service Training Helper Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}=====================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}=====================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if venv exists
if [ ! -d "venv" ]; then
    print_error "Virtual environment not found!"
    print_info "Run: python3 -m venv venv && source venv/bin/activate && pip install -r requirements.txt"
    exit 1
fi

# Activate venv
source venv/bin/activate

# Parse command
case "$1" in
    "data")
        print_header "Generating Training Data"
        python scripts/generate_training_data.py
        print_success "Training data generated!"
        ;;
        
    "train")
        print_header "Training Model"
        print_info "This will take 5-15 minutes..."
        python scripts/train_model.py
        print_success "Model trained!"
        print_info "Next: ./train.sh enable"
        ;;
        
    "enable")
        print_header "Enabling Fine-tuned Model"
        
        # Check if fine-tuned model exists
        if [ ! -d "models/fine_tuned" ]; then
            print_error "Fine-tuned model not found!"
            print_info "Run: ./train.sh train first"
            exit 1
        fi
        
        # Update .env
        if grep -q "USE_FINE_TUNED_MODEL" .env 2>/dev/null; then
            sed -i '' 's/USE_FINE_TUNED_MODEL=.*/USE_FINE_TUNED_MODEL=true/' .env
        else
            echo "USE_FINE_TUNED_MODEL=true" >> .env
        fi
        
        print_success "Fine-tuned model enabled in .env"
        print_info "Restart service: ./train.sh restart"
        ;;
        
    "disable")
        print_header "Disabling Fine-tuned Model"
        
        if grep -q "USE_FINE_TUNED_MODEL" .env 2>/dev/null; then
            sed -i '' 's/USE_FINE_TUNED_MODEL=.*/USE_FINE_TUNED_MODEL=false/' .env
        else
            echo "USE_FINE_TUNED_MODEL=false" >> .env
        fi
        
        print_success "Fine-tuned model disabled"
        print_info "Restart service: ./train.sh restart"
        ;;
        
    "restart")
        print_header "Restarting Service"
        
        # Stop service
        lsof -ti:8040 | xargs kill -9 2>/dev/null && print_success "Service stopped" || print_info "Service not running"
        
        sleep 2
        
        # Start service in background
        HTTP_PORT=8040 nohup python main.py > logs/service.log 2>&1 &
        
        sleep 3
        
        # Check health
        if curl -s http://localhost:8040/v2/ai/health > /dev/null; then
            print_success "Service started on port 8040"
            curl -s http://localhost:8040/v2/ai/model-info | python3 -m json.tool
        else
            print_error "Service failed to start"
            print_info "Check logs: tail -f logs/service.log"
            exit 1
        fi
        ;;
        
    "reembed")
        print_header "Re-embedding Knowledge Base"
        
        curl -X POST http://localhost:8040/v2/ai/embed | python3 -m json.tool
        
        print_success "Knowledge base re-embedded"
        ;;
        
    "info")
        print_header "Model Information"
        curl -s http://localhost:8040/v2/ai/model-info | python3 -m json.tool
        ;;
        
    "test")
        print_header "Testing Query"
        
        QUESTION="${2:-Làm sao tạo API mới trong BDSPro?}"
        
        print_info "Question: $QUESTION"
        echo ""
        
        curl -s -X POST http://localhost:8040/v2/ai/query \
          -H "Content-Type: application/json" \
          -d "{\"question\": \"$QUESTION\", \"top_k\": 3}" | \
          python3 -c "import sys, json; d=json.load(sys.stdin); print(f'Sources: {d[\"sources\"]}\nConfidence: {d[\"confidence\"]:.2f}\nTime: {d[\"processing_time_ms\"]}ms\n\nAnswer:\n{d[\"answer\"][:500]}...')"
        ;;
        
    "full")
        print_header "Full Training Pipeline"
        
        # 1. Generate data
        print_info "Step 1/5: Generating training data..."
        python scripts/generate_training_data.py
        
        # 2. Train
        print_info "Step 2/5: Training model (this takes time)..."
        python scripts/train_model.py
        
        # 3. Enable
        print_info "Step 3/5: Enabling fine-tuned model..."
        if grep -q "USE_FINE_TUNED_MODEL" .env 2>/dev/null; then
            sed -i '' 's/USE_FINE_TUNED_MODEL=.*/USE_FINE_TUNED_MODEL=true/' .env
        else
            echo "USE_FINE_TUNED_MODEL=true" >> .env
        fi
        
        # 4. Restart
        print_info "Step 4/5: Restarting service..."
        lsof -ti:8040 | xargs kill -9 2>/dev/null || true
        sleep 2
        HTTP_PORT=8040 nohup python main.py > logs/service.log 2>&1 &
        sleep 3
        
        # 5. Re-embed
        print_info "Step 5/5: Re-embedding knowledge base..."
        curl -s -X POST http://localhost:8040/v2/ai/embed > /dev/null
        
        print_success "Full training pipeline complete!"
        print_info "Test with: ./train.sh test"
        ;;
        
    *)
        echo "AI Service Training Helper"
        echo ""
        echo "Usage: ./train.sh <command>"
        echo ""
        echo "Commands:"
        echo "  data      - Generate training data (Q&A pairs)"
        echo "  train     - Train/fine-tune model"
        echo "  enable    - Enable fine-tuned model"
        echo "  disable   - Disable fine-tuned model (use base)"
        echo "  restart   - Restart service"
        echo "  reembed   - Re-embed knowledge base"
        echo "  info      - Show model info"
        echo "  test      - Test query (optional: ./train.sh test 'your question')"
        echo "  full      - Run full training pipeline"
        echo ""
        echo "Examples:"
        echo "  ./train.sh data           # Generate Q&A dataset"
        echo "  ./train.sh train          # Train model"
        echo "  ./train.sh enable         # Enable fine-tuned model"
        echo "  ./train.sh restart        # Restart service"
        echo "  ./train.sh test 'How to create API?'"
        echo "  ./train.sh full           # Do everything"
        ;;
esac

