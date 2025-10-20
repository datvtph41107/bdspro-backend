"""
Embed Context Files - Create embeddings for all context documentation
"""

import sys
from pathlib import Path
import os

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from app.core.embeddings import EmbeddingManager
from dotenv import load_dotenv

load_dotenv()

def chunk_text(text: str, chunk_size: int = 1000, overlap: int = 200) -> list:
    """
    Split text into overlapping chunks
    
    Args:
        text: Input text
        chunk_size: Size of each chunk
        overlap: Overlap between chunks
        
    Returns:
        List of text chunks
    """
    chunks = []
    start = 0
    
    while start < len(text):
        end = start + chunk_size
        chunk = text[start:end]
        chunks.append(chunk)
        start += chunk_size - overlap
    
    return chunks

def embed_all_context(embedding_manager: EmbeddingManager = None):
    """
    Embed all context files from ../context/
    
    Args:
        embedding_manager: Optional embedding manager instance
    """
    if embedding_manager is None:
        embedding_manager = EmbeddingManager()
    
    # Clear existing embeddings
    print("🗑️  Clearing existing embeddings...")
    embedding_manager.clear()
    
    # Get knowledge_base directory (inside ai-service)
    context_dir = Path(__file__).parent.parent / "knowledge_base"
    
    if not context_dir.exists():
        print(f"❌ Knowledge base directory not found: {context_dir}")
        return 0
    
    print(f"📂 Reading knowledge base files from: {context_dir}")
    
    total_chunks = 0
    
    # Process each markdown file
    for md_file in sorted(context_dir.glob("*.md")):
        print(f"\n📄 Processing: {md_file.name}")
        
        # Read file
        with open(md_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Skip empty files
        if not content.strip():
            print(f"   ⏭️  Skipped (empty)")
            continue
        
        # Split into chunks
        chunks = chunk_text(content)
        print(f"   ✂️  Split into {len(chunks)} chunks")
        
        # Prepare data
        texts = []
        metadatas = []
        ids = []
        
        for i, chunk in enumerate(chunks):
            texts.append(chunk)
            metadatas.append({
                "source": md_file.name,
                "chunk_id": i,
                "file_path": str(md_file)
            })
            ids.append(f"{md_file.stem}_{i}")
        
        # Add to vector store
        print(f"   🔢 Creating embeddings...")
        embedding_manager.add_documents(
            texts=texts,
            metadatas=metadatas,
            ids=ids
        )
        
        total_chunks += len(chunks)
        print(f"   ✅ Embedded {len(chunks)} chunks")
    
    # Summary
    print(f"\n{'='*60}")
    print(f"✅ DONE!")
    print(f"   Total files processed: {len(list(context_dir.glob('*.md')))}")
    print(f"   Total chunks embedded: {total_chunks}")
    print(f"   Vector DB count: {embedding_manager.get_count()}")
    print(f"{'='*60}\n")
    
    return total_chunks

if __name__ == "__main__":
    print("🚀 BDSPro Context Embedder (LOCAL MODE)")
    print("   Using sentence-transformers for embeddings")
    print("="*60)
    
    try:
        embed_all_context()
    except Exception as e:
        print(f"\n❌ Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

