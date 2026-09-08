"""
RAG Retriever - Retrieves relevant context for queries
"""

from typing import Dict, List
from .embeddings import EmbeddingManager

class RAGRetriever:
    """Handles retrieval of relevant context for RAG"""
    
    def __init__(self, embedding_manager: EmbeddingManager):
        """
        Initialize RAG Retriever
        
        Args:
            embedding_manager: Embedding manager instance
        """
        self.embedding_manager = embedding_manager
    
    def retrieve(
        self,
        question: str,
        top_k: int = 5
    ) -> Dict:
        """
        Retrieve relevant context for question
        
        Args:
            question: User question
            top_k: Number of chunks to retrieve
            
        Returns:
            Dict with context, sources, chunks, confidence
        """
        # Search vector DB
        results = self.embedding_manager.search(
            query=question,
            top_k=top_k
        )
        
        if not results['documents'] or not results['documents'][0]:
            return {
                'context': '',
                'sources': [],
                'chunks': [],
                'confidence': 0.0
            }
        
        # Extract results
        documents = results['documents'][0]
        metadatas = results['metadatas'][0]
        distances = results['distances'][0] if 'distances' in results else [0] * len(documents)
        
        # Build context
        context_parts = []
        sources = set()
        
        for doc, metadata in zip(documents, metadatas):
            context_parts.append(doc)
            if 'source' in metadata:
                sources.add(metadata['source'])
        
        # Join context with separators
        context = "\n\n---\n\n".join(context_parts)
        
        # Calculate confidence (inverse of average distance)
        avg_distance = sum(distances) / len(distances) if distances else 1.0
        confidence = max(0.0, min(1.0, 1.0 - avg_distance))
        
        return {
            'context': context,
            'sources': list(sources),
            'chunks': documents,
            'confidence': confidence
        }
    
    def retrieve_with_scores(
        self,
        question: str,
        top_k: int = 5,
        threshold: float = 0.7
    ) -> List[Dict]:
        """
        Retrieve with relevance scores, filtered by threshold
        
        Args:
            question: User question
            top_k: Number of chunks to retrieve
            threshold: Minimum relevance score (0-1)
            
        Returns:
            List of dicts with chunk, source, score
        """
        results = self.embedding_manager.search(
            query=question,
            top_k=top_k
        )
        
        if not results['documents'] or not results['documents'][0]:
            return []
        
        documents = results['documents'][0]
        metadatas = results['metadatas'][0]
        distances = results['distances'][0] if 'distances' in results else [0] * len(documents)
        
        # Build results with scores
        scored_results = []
        for doc, metadata, distance in zip(documents, metadatas, distances):
            score = 1.0 - distance
            
            # Filter by threshold
            if score >= threshold:
                scored_results.append({
                    'chunk': doc,
                    'source': metadata.get('source', 'unknown'),
                    'score': score
                })
        
        # Sort by score descending
        scored_results.sort(key=lambda x: x['score'], reverse=True)
        
        return scored_results

