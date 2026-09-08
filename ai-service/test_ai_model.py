#!/usr/bin/env python3
"""
Test script for BDSPro AI Service
Tests the AI model with various questions about the project
"""

import requests
import json
import time

# AI Service endpoint
AI_SERVICE_URL = "http://localhost:8090"

def test_ai_service():
    """Test the AI service with various questions"""
    
    print("🤖 Testing BDSPro AI Service")
    print("=" * 50)
    
    # Test questions about the project
    test_questions = [
        "What is the overall architecture of BDSPro microservices?",
        "How does authentication work in the auth service?",
        "What are the main entities in BDSPro service?",
        "How does payment processing work?",
        "What is Clean Architecture pattern used in this project?",
        "How to create a new API endpoint?",
        "What are the common code patterns used?",
        "How does inter-service communication work?",
        "What is the database schema for products?",
        "How to implement repository pattern?"
    ]
    
    for i, question in enumerate(test_questions, 1):
        print(f"\n📝 Test {i}: {question}")
        print("-" * 40)
        
        try:
            # Make request to AI service
            response = requests.post(
                f"{AI_SERVICE_URL}/v2/ai/query",
                json={"question": question, "top_k": 3},
                timeout=30
            )
            
            if response.status_code == 200:
                result = response.json()
                print(f"✅ Answer: {result.get('answer', 'No answer')[:200]}...")
                print(f"📚 Sources: {len(result.get('sources', []))} documents")
                print(f"⏱️  Time: {result.get('processing_time_ms', 0)}ms")
            else:
                print(f"❌ Error: {response.status_code} - {response.text}")
                
        except requests.exceptions.RequestException as e:
            print(f"❌ Connection error: {e}")
        
        time.sleep(1)  # Wait between requests
    
    print("\n🎉 AI Service testing completed!")

def test_health_check():
    """Test health check endpoint"""
    print("\n🏥 Testing Health Check")
    print("-" * 30)
    
    try:
        response = requests.get(f"{AI_SERVICE_URL}/v2/ai/health", timeout=10)
        if response.status_code == 200:
            print("✅ AI Service is healthy")
            print(f"📊 Status: {response.json()}")
        else:
            print(f"❌ Health check failed: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"❌ Health check error: {e}")

if __name__ == "__main__":
    print("🚀 Starting BDSPro AI Service Tests")
    print("=" * 50)
    
    # Wait for service to start
    print("⏳ Waiting for AI service to start...")
    time.sleep(5)
    
    # Test health check first
    test_health_check()
    
    # Test AI queries
    test_ai_service()
    
    print("\n✨ All tests completed!")

