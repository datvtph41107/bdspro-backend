#!/usr/bin/env python3
"""
Script test API analyze product - Đọc dữ liệu từ file JSON
Sử dụng: python3 test_api_from_json.py [path_to_json_file]
Nếu không truyền path, sẽ dùng file test_data_template.json
"""

import requests
import csv
import time
from datetime import datetime
import json
import sys
import os

# Cấu hình API
API_URL = "http://localhost:8000/v2/assistant/deepseek/analyze/product"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpdGllcyI6bnVsbCwiZXhwIjoxNzYyNDA0MjM0LCJpYXQiOjE3NjIzOTUyMzQsIm9yZ2FuaXphdGlvbiI6bnVsbCwicGxhbiI6bnVsbCwicGxhbkF0IjpudWxsLCJwcm9maWxlIjoxNTQsInJvbGUiOiIiLCJzZXNzaW9uIjo4ODksInN1YiI6MjY1LCJ0eXBlIjoiQUNDRVNTIn0.Zkv0vLGtMQpMgZUCFxjrnCGdSr7ZJ4LRnyLGl1-KAMU"

HEADERS = {
    "accept": "application/json",
    "Authorization": AUTH_TOKEN,
    "Content-Type": "application/json"
}

def load_test_data(json_file):
    """
    Load dữ liệu test từ file JSON
    
    Args:
        json_file: Path đến file JSON
        
    Returns:
        list: Danh sách test cases
    """
    try:
        with open(json_file, 'r', encoding='utf-8') as f:
            data = json.load(f)
        
        if not isinstance(data, list):
            print(f"❌ Error: File JSON phải chứa một mảng (array)")
            sys.exit(1)
        
        # Validate format
        for i, item in enumerate(data):
            if not isinstance(item, dict):
                print(f"❌ Error: Item {i+1} không phải là object")
                sys.exit(1)
            if 'content' not in item:
                print(f"❌ Error: Item {i+1} thiếu field 'content'")
                sys.exit(1)
            if 'context' not in item:
                item['context'] = ""  # Default context rỗng
        
        return data
    
    except FileNotFoundError:
        print(f"❌ Error: Không tìm thấy file {json_file}")
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"❌ Error: File JSON không hợp lệ - {e}")
        sys.exit(1)

def call_api(content, context):
    """Gọi API analyze product"""
    payload = {
        "content": content,
        "context": context
    }
    
    start_time = time.time()
    
    try:
        response = requests.post(
            API_URL,
            headers=HEADERS,
            json=payload,
            timeout=30
        )
        
        response_time = round((time.time() - start_time) * 1000, 2)
        
        result = {
            "status_code": response.status_code,
            "response_time_ms": response_time,
            "success": response.status_code == 200,
            "response_data": None,
            "error": None
        }
        
        try:
            result["response_data"] = response.json()
        except:
            result["response_data"] = response.text
            
        if not result["success"]:
            result["error"] = f"HTTP {response.status_code}: {response.text[:200]}"
            
        return result
        
    except requests.exceptions.Timeout:
        return {
            "status_code": 0,
            "response_time_ms": round((time.time() - start_time) * 1000, 2),
            "success": False,
            "response_data": None,
            "error": "Request timeout (>30s)"
        }
    except requests.exceptions.ConnectionError:
        return {
            "status_code": 0,
            "response_time_ms": 0,
            "success": False,
            "response_data": None,
            "error": "Connection error - API server không khả dụng"
        }
    except Exception as e:
        return {
            "status_code": 0,
            "response_time_ms": round((time.time() - start_time) * 1000, 2),
            "success": False,
            "response_data": None,
            "error": str(e)
        }

def extract_ai_response(response_data):
    """Trích xuất response text từ AI"""
    if not response_data:
        return ""
    
    if isinstance(response_data, dict):
        for key in ["result", "message", "content", "response", "answer", "text"]:
            if key in response_data:
                value = response_data[key]
                if isinstance(value, str):
                    return value
                elif isinstance(value, dict):
                    return json.dumps(value, ensure_ascii=False)
        
        return json.dumps(response_data, ensure_ascii=False)
    
    return str(response_data)

def save_to_csv(results, filename):
    """Lưu kết quả vào file CSV"""
    with open(filename, 'w', newline='', encoding='utf-8-sig') as csvfile:
        fieldnames = [
            'STT',
            'Content',
            'Context', 
            'Status Code',
            'Success',
            'Response Time (ms)',
            'AI Response',
            'Full Response',
            'Error',
            'Timestamp'
        ]
        
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        
        for i, result in enumerate(results, 1):
            # Format full response
            full_response = ""
            if result.get('full_response_data'):
                if isinstance(result['full_response_data'], dict):
                    full_response = json.dumps(result['full_response_data'], ensure_ascii=False)
                else:
                    full_response = str(result['full_response_data'])
            
            writer.writerow({
                'STT': i,
                'Content': result['content'],
                'Context': result['context'],
                'Status Code': result['status_code'],
                'Success': 'Thành công' if result['success'] else 'Thất bại',
                'Response Time (ms)': result['response_time_ms'],
                'AI Response': result['ai_response'][:1000] if result['ai_response'] else '',
                'Full Response': full_response[:2000] if full_response else '',
                'Error': result['error'] or '',
                'Timestamp': result['timestamp']
            })

def print_statistics(results):
    """In thống kê kết quả"""
    total = len(results)
    success = sum(1 for r in results if r['success'])
    failed = total - success
    
    response_times = [r['response_time_ms'] for r in results if r['response_time_ms'] > 0]
    avg_response_time = sum(response_times) / len(response_times) if response_times else 0
    min_response_time = min(response_times) if response_times else 0
    max_response_time = max(response_times) if response_times else 0
    
    print("\n" + "="*70)
    print("THỐNG KÊ KẾT QUẢ TEST API")
    print("="*70)
    print(f"Tổng số request:        {total}")
    print(f"Thành công:             {success} ({success/total*100:.1f}%)")
    print(f"Thất bại:               {failed} ({failed/total*100:.1f}%)")
    if response_times:
        print(f"Thời gian phản hồi TB:  {avg_response_time:.2f} ms")
        print(f"Thời gian phản hồi min: {min_response_time:.2f} ms")
        print(f"Thời gian phản hồi max: {max_response_time:.2f} ms")
    print("="*70)
    
    # Chi tiết từng request
    print("\nCHI TIẾT TỪNG REQUEST:")
    print("-"*70)
    for i, result in enumerate(results, 1):
        status = "✓" if result['success'] else "✗"
        print(f"{status} Request {i}: {result['status_code']} - {result['response_time_ms']}ms")
        print(f"  Content: {result['content'][:70]}...")
        if result['ai_response'] and len(result['ai_response']) > 0:
            preview = result['ai_response'][:100].replace('\n', ' ')
            print(f"  Response: {preview}...")
        if result['error']:
            print(f"  Error: {result['error']}")
        print()

def main():
    """Hàm chính"""
    # Xác định file JSON để load
    if len(sys.argv) > 1:
        json_file = sys.argv[1]
    else:
        json_file = "test_data_template.json"
    
    print("="*70)
    print("TEST API ANALYZE PRODUCT - Load từ JSON")
    print("="*70)
    print(f"API URL: {API_URL}")
    print(f"Data file: {json_file}")
    
    # Load dữ liệu
    test_data = load_test_data(json_file)
    print(f"Số lượng test cases: {len(test_data)}")
    print("="*70)
    print()
    
    results = []
    
    for i, test_case in enumerate(test_data, 1):
        print(f"[{i}/{len(test_data)}] Đang gọi API...")
        content_preview = test_case['content'][:60] + "..." if len(test_case['content']) > 60 else test_case['content']
        print(f"  Content: {content_preview}")
        
        api_result = call_api(
            content=test_case['content'],
            context=test_case.get('context', '')
        )
        
        result = {
            'content': test_case['content'],
            'context': test_case.get('context', ''),
            'status_code': api_result['status_code'],
            'success': api_result['success'],
            'response_time_ms': api_result['response_time_ms'],
            'ai_response': extract_ai_response(api_result['response_data']),
            'full_response_data': api_result['response_data'],
            'error': api_result['error'],
            'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        }
        
        results.append(result)
        
        status_icon = "✓" if result['success'] else "✗"
        print(f"  {status_icon} Status: {result['status_code']} - {result['response_time_ms']}ms")
        
        if result['error']:
            print(f"  ❌ Error: {result['error']}")
        
        print()
        
        # Delay giữa các request
        if i < len(test_data):
            time.sleep(0.5)
    
    # Lưu kết quả vào CSV
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    csv_filename = f"api_test_results_{timestamp}.csv"
    save_to_csv(results, csv_filename)
    print(f"✓ Đã lưu kết quả vào file CSV: {csv_filename}")
    
    # Lưu file JSON chi tiết
    json_filename = f"api_test_results_{timestamp}.json"
    with open(json_filename, 'w', encoding='utf-8') as f:
        json.dump(results, f, ensure_ascii=False, indent=2)
    print(f"✓ Đã lưu kết quả chi tiết vào file JSON: {json_filename}")
    
    # In thống kê
    print_statistics(results)
    
    print(f"\n{'='*70}")
    print("✓ HOÀN THÀNH!")
    print(f"{'='*70}")

if __name__ == "__main__":
    main()

