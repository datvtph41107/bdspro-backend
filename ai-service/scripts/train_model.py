"""
Train Model - Fine-tune embedding model on BDSPro Q&A data
"""

import sys
from pathlib import Path
import time

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from app.core.trainer import train_model

def main():
    """Run model training"""
    
    print("🤖 BDSPro Model Training")
    print("="*60)
    
    # Paths
    data_path = Path(__file__).parent.parent / "data" / "train" / "qa_dataset.json"
    output_path = Path(__file__).parent.parent / "models" / "fine_tuned"
    
    if not data_path.exists():
        print(f"❌ Training data not found: {data_path}")
        print("   Run: python scripts/generate_training_data.py")
        sys.exit(1)
    
    print(f"\n📂 Training data: {data_path}")
    print(f"📂 Output path: {output_path}")
    
    # Training parameters
    epochs = 3
    batch_size = 16
    
    print(f"\n⚙️  Training parameters:")
    print(f"   Epochs: {epochs}")
    print(f"   Batch size: {batch_size}")
    print(f"   Base model: all-MiniLM-L6-v2")
    
    print(f"\n🚀 Starting training...")
    print("   This may take 5-15 minutes depending on your hardware...")
    
    start_time = time.time()
    
    try:
        results = train_model(
            data_path=str(data_path),
            output_path=str(output_path),
            epochs=epochs,
            batch_size=batch_size
        )
        
        elapsed_time = time.time() - start_time
        
        print(f"\n{'='*60}")
        print(f"✅ TRAINING COMPLETE!")
        print(f"   Status: {results['status']}")
        print(f"   Model saved: {results['model_path']}")
        print(f"   Training examples: {results['training_examples']}")
        print(f"   Total steps: {results['total_steps']}")
        print(f"   Time elapsed: {elapsed_time:.1f}s")
        print(f"{'='*60}\n")
        
        print("💡 Next steps:")
        print("   1. Update .env: USE_FINE_TUNED_MODEL=true")
        print("   2. Restart service to use fine-tuned model")
        print("   3. Test with: curl -X POST http://localhost:8040/v2/ai/query")
        
    except Exception as e:
        print(f"\n❌ Training failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

if __name__ == "__main__":
    main()

