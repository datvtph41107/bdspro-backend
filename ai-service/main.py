"""
BDSPro AI Service - RAG-based Knowledge Assistant (LOCAL MODE)
Port: 8040 (HTTP) - No external API calls
"""

import uvicorn
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import time
import os
from pathlib import Path

# Import core components
from app.core.embeddings import EmbeddingManager
from app.core.retrieval import RAGRetriever
from app.core.llm import LLMClient

# ============================================================================
# Models
# ============================================================================

class QueryRequest(BaseModel):
    question: str
    top_k: int = 5
    temperature: float = 0.3

class QueryData(BaseModel):
    answer: str
    sources: List[str]
    confidence: float
    processing_time_ms: int
    context_chunks: int

class QueryResponse(BaseModel):
    success: bool
    message: str
    data: Optional[QueryData] = None

class HealthResponse(BaseModel):
    status: str
    version: str
    vector_db_count: int
    uptime_seconds: int

class StatsResponse(BaseModel):
    total_queries: int
    avg_processing_time_ms: float
    total_embeddings: int
    cache_hit_rate: float

# ============================================================================
# FastAPI App
# ============================================================================

app = FastAPI(
    title="BDSPro AI Service",
    description="RAG-based AI assistant for BDSPro codebase",
    version="1.0.0"
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ============================================================================
# Global State
# ============================================================================

embedding_manager: Optional[EmbeddingManager] = None
rag_retriever: Optional[RAGRetriever] = None
llm_client: Optional[LLMClient] = None
start_time = time.time()
query_count = 0
total_processing_time = 0

# ============================================================================
# Startup
# ============================================================================

@app.on_event("startup")
async def startup_event():
    """Initialize components on startup"""
    global embedding_manager, rag_retriever, llm_client
    
    print("🚀 Starting BDSPro AI Service (LOCAL MODE)...")
    
    # Check if should use fine-tuned model
    use_fine_tuned = os.getenv("USE_FINE_TUNED_MODEL", "false").lower() == "true"
    fine_tuned_path = os.getenv("FINE_TUNED_MODEL_PATH", "./models/fine_tuned")
    
    # Initialize components (no API keys needed)
    embedding_manager = EmbeddingManager(
        chroma_path=os.getenv("CHROMA_DB_PATH", "./data/chroma_db"),
        use_fine_tuned=use_fine_tuned,
        fine_tuned_path=fine_tuned_path
    )
    
    rag_retriever = RAGRetriever(embedding_manager)
    
    llm_client = LLMClient()  # No API key needed for local mode
    
    # Check vector DB
    count = embedding_manager.get_count()
    if count == 0:
        print("⚠️  Warning: Vector DB is empty!")
        print("   Run: python scripts/embed_context.py")
    else:
        print(f"✅ Loaded {count} document chunks")
    
    print(f"✅ AI Service ready on port {os.getenv('HTTP_PORT', 8040)}")

# ============================================================================
# Endpoints
# ============================================================================

@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "service": "BDSPro AI Service",
        "version": "1.0.0",
        "status": "running",
        "docs": "/docs"
    }

@app.post("/v2/ai/query", response_model=QueryResponse)
async def query(request: QueryRequest):
    """
    Query the knowledge base
    
    Request:
    {
        "question": "Làm sao tạo API mới?",
        "top_k": 5,
        "temperature": 0.3
    }
    
    Response:
    {
        "success": true,
        "message": "Query successful",
        "data": {
            "answer": "...",
            "sources": ["file1.md", "file2.md"],
            "confidence": 0.85,
            "processing_time_ms": 250,
            "context_chunks": 5
        }
    }
    """
    global query_count, total_processing_time
    
    start = time.time()
    
    try:
        # Validate input
        if not request.question or len(request.question.strip()) == 0:
            return QueryResponse(
                success=False,
                message="Question is required and cannot be empty",
                data=None
            )
        
        # 1. Retrieve relevant context
        context_results = rag_retriever.retrieve(
            question=request.question,
            top_k=request.top_k
        )
        
        if not context_results:
            return QueryResponse(
                success=False,
                message="No relevant context found in knowledge base",
                data=None
            )
        
        # 2. Generate answer with LLM
        answer = await llm_client.generate_answer(
            question=request.question,
            context=context_results['context'],
            temperature=request.temperature
        )
        
        # 3. Calculate processing time
        processing_time = int((time.time() - start) * 1000)
        
        # 4. Update stats
        query_count += 1
        total_processing_time += processing_time
        
        # 5. Return structured response
        return QueryResponse(
            success=True,
            message="Query successful",
            data=QueryData(
                answer=answer,
                sources=context_results['sources'],
                confidence=context_results['confidence'],
                processing_time_ms=processing_time,
                context_chunks=len(context_results['chunks'])
            )
        )
    
    except Exception as e:
        # Log error
        import traceback
        traceback.print_exc()
        
        return QueryResponse(
            success=False,
            message=f"Query failed: {str(e)}",
            data=None
        )

@app.get("/v2/ai/health", response_model=HealthResponse)
async def health():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        version="1.0.0",
        vector_db_count=embedding_manager.get_count() if embedding_manager else 0,
        uptime_seconds=int(time.time() - start_time)
    )

@app.get("/v2/ai/stats", response_model=StatsResponse)
async def stats():
    """Service statistics"""
    avg_time = (total_processing_time / query_count) if query_count > 0 else 0
    
    return StatsResponse(
        total_queries=query_count,
        avg_processing_time_ms=avg_time,
        total_embeddings=embedding_manager.get_count() if embedding_manager else 0,
        cache_hit_rate=0.0  # TODO: Implement caching
    )

@app.post("/v2/ai/embed")
async def embed():
    """
    Re-embed context files
    (Admin endpoint - should be protected in production)
    """
    try:
        from scripts.embed_context import embed_all_context
        
        count = await embed_all_context(embedding_manager)
        
        return {
            "status": "success",
            "embedded_chunks": count,
            "message": "Context files embedded successfully"
        }
    
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Embedding failed: {str(e)}"
        )

@app.post("/v2/ai/train")
async def train():
    """
    Train/fine-tune the embedding model
    
    Trains from: all JSON files in data/train/
    
    Steps:
    1. Load all Q&A pairs from JSON files in data/train/
    2. Fine-tune sentence-transformers model
    3. Save to ./models/fine_tuned
    4. Reload embedding manager with new model
    5. Re-embed knowledge base automatically
    
    Note: This is a long-running operation (10-20 minutes)
    """
    global embedding_manager, rag_retriever
    
    try:
        from pathlib import Path
        
        # Check if training directory exists
        train_dir_path = Path("./data/train")
        
        if not train_dir_path.exists():
            raise HTTPException(
                status_code=400,
                detail=f"Training directory not found. Please ensure {str(train_dir_path)} exists."
            )
        
        # Count JSON files
        json_files = list(train_dir_path.glob("*.json"))
        if not json_files:
            raise HTTPException(
                status_code=400,
                detail=f"No JSON training files found in {str(train_dir_path)}"
            )
        
        print(f"\n{'='*60}")
        print(f"🚀 STEP 1/3: Training model from {len(json_files)} file(s)...")
        print(f"{'='*60}")
        for json_file in json_files:
            print(f"  - {json_file.name}")
        
        from app.core.trainer import train_model
        
        # Run training (pass directory path, it will read all JSON files)
        results = train_model(
            data_path=str(train_dir_path),
            output_path="./models/fine_tuned",
            epochs=16,
            batch_size=16
        )
        
        print(f"\n{'='*60}")
        print(f"✅ Training completed!")
        print(f"{'='*60}\n")
        
        # Step 2: Reload embedding manager with fine-tuned model
        print(f"{'='*60}")
        print(f"🔄 STEP 2/3: Reloading embedding manager with fine-tuned model...")
        print(f"{'='*60}\n")
        
        # Create new embedding manager with fine-tuned model
        new_embedding_manager = EmbeddingManager(
            chroma_path=os.getenv("CHROMA_DB_PATH", "./data/chroma_db"),
            use_fine_tuned=True,
            fine_tuned_path="./models/fine_tuned"
        )
        
        print("✅ New embedding manager loaded with fine-tuned model\n")
        
        # Step 3: Re-embed knowledge base
        print(f"{'='*60}")
        print(f"📚 STEP 3/3: Re-embedding knowledge base...")
        print(f"{'='*60}\n")
        
        from scripts.embed_context import embed_all_context
        
        embedded_chunks = embed_all_context(new_embedding_manager)
        
        # Update global state to use new model
        embedding_manager = new_embedding_manager
        rag_retriever = RAGRetriever(embedding_manager)
        
        print(f"\n{'='*60}")
        print(f"🎉 ALL STEPS COMPLETED!")
        print(f"{'='*60}\n")
        
        return {
            "status": "success",
            "message": "Model trained and knowledge base re-embedded successfully",
            "results": results,
            "training_files": [f.name for f in json_files],
            "training_data_dir": str(train_dir_path),
            "embedded_chunks": embedded_chunks,
            "model_activated": True,
            "info": [
                "✅ Model trained successfully",
                "✅ Embedding manager reloaded with fine-tuned model",
                "✅ Knowledge base re-embedded with new model",
                "✅ Service now using fine-tuned model automatically"
            ]
        }
    
    except HTTPException:
        raise
    except Exception as e:
        import traceback
        traceback.print_exc()
        raise HTTPException(
            status_code=500,
            detail=f"Training failed: {str(e)}"
        )

@app.get("/v2/ai/model-info")
async def model_info():
    """Get information about current model"""
    return {
        "model_type": embedding_manager.model_type if embedding_manager else "unknown",
        "base_model": "all-MiniLM-L6-v2",
        "fine_tuned_available": Path("./models/fine_tuned").exists(),
        "vector_db_count": embedding_manager.get_count() if embedding_manager else 0
    }

# ============================================================================
# Main
# ============================================================================

if __name__ == "__main__":
    port = int(os.getenv("HTTP_PORT", 8040))
    
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=port,
        reload=True,  # Auto-reload during development
        log_level="info"
    )

