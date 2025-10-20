"""
Embedding Manager - Handles document embeddings and vector storage (LOCAL)
"""

import chromadb
from chromadb.config import Settings
from sentence_transformers import SentenceTransformer
from typing import List, Dict, Optional
import os

class EmbeddingManager:
    """Manages document embeddings and ChromaDB operations"""
    
    def __init__(
        self, 
        chroma_path: str = "./data/chroma_db",
        use_fine_tuned: bool = False,
        fine_tuned_path: str = "./models/fine_tuned"
    ):
        """
        Initialize Embedding Manager
        
        Args:
            chroma_path: Path to ChromaDB storage
            use_fine_tuned: Whether to use fine-tuned model
            fine_tuned_path: Path to fine-tuned model
        """
        self.chroma_path = chroma_path
        
        # Check if fine-tuned model exists and should be used
        from pathlib import Path
        fine_tuned_exists = Path(fine_tuned_path).exists()
        
        if use_fine_tuned and fine_tuned_exists:
            print(f"📦 Loading FINE-TUNED model from: {fine_tuned_path}")
            self.embedding_model = SentenceTransformer(fine_tuned_path)
            self.model_type = "fine_tuned"
            print("✅ Fine-tuned model loaded!")
        else:
            print("📦 Loading base model (sentence-transformers)...")
            self.embedding_model = SentenceTransformer('all-MiniLM-L6-v2')
            self.model_type = "base"
            print("✅ Base embedding model loaded!")
        
        # Initialize ChromaDB
        self.chroma_client = chromadb.PersistentClient(
            path=chroma_path,
            settings=Settings(
                anonymized_telemetry=False
            )
        )
        
        # Get or create collection
        self.collection = self.chroma_client.get_or_create_collection(
            name="bdspro_context",
            metadata={"hnsw:space": "cosine"}
        )
    
    def embed_text(self, text: str) -> List[float]:
        """
        Create embedding for text using local model
        
        Args:
            text: Input text
            
        Returns:
            Embedding vector
        """
        # Use sentence-transformers (local)
        embedding = self.embedding_model.encode(text, convert_to_numpy=True)
        return embedding.tolist()
    
    def add_documents(
        self,
        texts: List[str],
        metadatas: List[Dict],
        ids: List[str]
    ):
        """
        Add documents to vector store
        
        Args:
            texts: List of text chunks
            metadatas: List of metadata dicts
            ids: List of unique IDs
        """
        # Create embeddings
        embeddings = []
        for text in texts:
            embedding = self.embed_text(text)
            embeddings.append(embedding)
        
        # Add to collection
        self.collection.add(
            embeddings=embeddings,
            documents=texts,
            metadatas=metadatas,
            ids=ids
        )
    
    def search(
        self,
        query: str,
        top_k: int = 5
    ) -> Dict:
        """
        Search for relevant documents
        
        Args:
            query: Search query
            top_k: Number of results to return
            
        Returns:
            Search results with documents, metadatas, distances
        """
        # Embed query
        query_embedding = self.embed_text(query)
        
        # Search
        results = self.collection.query(
            query_embeddings=[query_embedding],
            n_results=top_k
        )
        
        return results
    
    def get_count(self) -> int:
        """Get total number of documents in collection"""
        return self.collection.count()
    
    def clear(self):
        """Clear all documents from collection"""
        # Delete and recreate collection
        self.chroma_client.delete_collection(name="bdspro_context")
        self.collection = self.chroma_client.get_or_create_collection(
            name="bdspro_context",
            metadata={"hnsw:space": "cosine"}
        )

