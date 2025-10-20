"""
LLM Client - LOCAL version (no external API calls)
"""

from typing import Optional
import os

class LLMClient:
    """Handles response generation (LOCAL mode - no LLM API)"""
    
    def __init__(
        self,
        api_key: str = None,
        model: str = "local"
    ):
        """
        Initialize LLM Client (LOCAL mode)
        
        Args:
            api_key: Not used in local mode
            model: Not used in local mode
        """
        self.model = "local-rag-only"
        print("✅ LLM Client initialized in LOCAL mode (no external API)")
    
    async def generate_answer(
        self,
        question: str,
        context: str,
        temperature: float = 0.3
    ) -> str:
        """
        Generate answer using retrieved context only (no LLM)
        
        Args:
            question: User question
            context: Retrieved context
            temperature: Not used in local mode
            
        Returns:
            Formatted context as answer
        """
        # Simple response with context only
        answer = f"""**🔍 Thông tin tìm được trong tài liệu:**

{context}

---

**💡 Ghi chú:** Service đang chạy ở chế độ LOCAL (không dùng AI model). 
Kết quả trên là context tìm được từ vector database dựa trên câu hỏi: "{question}"

Để có câu trả lời tự nhiên hơn, cần integrate với LLM (GPT/Claude/Local LLM)."""

        return answer
    
    async def generate_with_streaming(
        self,
        question: str,
        context: str,
        temperature: float = 0.3
    ):
        """
        Generate answer with streaming response
        
        Args:
            question: User question
            context: Retrieved context
            temperature: Sampling temperature
            
        Yields:
            Response chunks
        """
        user_prompt = f"""Context từ documentation:

{context}

---

Câu hỏi: {question}

Trả lời:"""

        stream = await self.client.chat.completions.create(
            model=self.model,
            messages=[
                {"role": "system", "content": self.system_prompt},
                {"role": "user", "content": user_prompt}
            ],
            temperature=temperature,
            max_tokens=2000,
            stream=True
        )
        
        async for chunk in stream:
            if chunk.choices[0].delta.content:
                yield chunk.choices[0].delta.content

