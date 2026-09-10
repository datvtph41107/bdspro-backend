#!/usr/bin/env python3
"""
Script test API analyze product của assistant-service
Gọi API với dữ liệu mẫu và lưu kết quả vào CSV
"""

import requests
import csv
import time
from datetime import datetime
import json

# Cấu hình API
API_URL = "http://localhost:8000/v2/assistant/deepseek/analyze/product"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpdGllcyI6bnVsbCwiZXhwIjoxNzYyNDA0MjM0LCJpYXQiOjE3NjIzOTUyMzQsIm9yZ2FuaXphdGlvbiI6bnVsbCwicGxhbiI6bnVsbCwicGxhbkF0IjpudWxsLCJwcm9maWxlIjoxNTQsInJvbGUiOiIiLCJzZXNzaW9uIjo4ODksInN1YiI6MjY1LCJ0eXBlIjoiQUNDRVNTIn0.Zkv0vLGtMQpMgZUCFxjrnCGdSr7ZJ4LRnyLGl1-KAMU"

HEADERS = {
    "accept": "application/json",
    "Authorization": AUTH_TOKEN,
    "Content-Type": "application/json"
}

# Dữ liệu mẫu để test
# Bạn có thể thêm/sửa dữ liệu ở đây
TEST_DATA = [
    {
        "content": """⭐️ ÔI TRỜI E CẠN LỜI ĐỂ NÓI VỀ LÔ ĐẤT NÀY 
💫💫 Lô đất tuyến 2 đẹp như này giờ hiếm lắm.
🎗Nằm cách Quốc Lộ 21A chỉ có 30m, đường đã  được phê duyệt mở rộng 2 bên 80m.
🎗 Mặt 5m, hậu 5m, dài 18m, có diện tích 89,0 m2, có 60m2 đất ở?
🎗 Cách khu công nghệ cao có 4km, cách ngã tư lục quân 800m, cách sân golf Đồng Mô 500m?
🎗Đường to Thông suốt oto vào tận đất, vị trí quá lý tưởng để ở và làm vc?
🎗 Cách chùa Khai Nguyên có 300m?
🎗Xung quanh toàn nhà cao tầng, dân ở đông vui, an ninh cực tốt?
🎗Sổ hồng sang tên nhanh gọn?
💵 Anh chị quan tâm LH 
———————————————————   
VPBĐS Gia Hoà ☎️0973954457. Nhanh còn chậm hết ạ?""",
        "context": ""
    },
    {
        "content": """
        - Bằng Y , Sơn Đà , Huyện Ba Vì ( cũ ) . 
Chủ nhờ bán lô đất sẵn nhà , sẵn truồng trại . 
- S : 946m2 có 60m thổ cư còn lại là đất trồng cây lâu năm .
- Lô góc 2 mặt tiền , mặt chính 14,5m . Đường thông thoáng ô tô đánh võng . Thế đất cao , thoáng . View hồ sen . 
- Đất cách đường TL413 chỉ vài trăm mét . 
- Giá chưa tới 2 tỷ 
- Lh zalo 0969722996
        """,
        "context": ""
    },
    {
        "content": """
        ✈ Bán lô góc 2 mặt tiền diện tích chỉ 88m2 
🚀 Khu đấu giá X1, Lục Xuân, xã Phúc Lộc, (Phúc Thọ) Hà Nội
🍗 Mặt tiền 6m, bên cạnh cụm công nghiệp Võng Xuyên 
🍺 80m ra đường tỉnh 418, một phút ra ĐL Tây Thăng Long 
🍾 Đối diện đang làm hạ tầng đấu giá, đất trung tâm xã 
🍽 Giá chỉ 4 tỷ 850 triệu liên hệ 0968.785.586 / 0372.556.875 để làm việc
Hoa hồng 30 ạ
        """,
        "context": ""
    },
    {
        "content": """
        🤙🤙 ALO QUÝ KHÁCH HÀNG, E XIN NHẮC LẠI LẦN CUỐI CHO QUÝ NHÀ ĐẦU TƯ NHỚ Ạ.
🤙 E BÁN LÔ ĐẤT HIẾM : Mặt 5,5m, dài 16m, diện tích 85m2 có 60m2 
đất ở.
♥️ ĐẤT Ở XÃ SƠN ĐÔNG, THỊ XÃ SƠN TÂY, TP HÀ NỘI( nay là xã ĐOÀI PHƯƠNG, TP HÀ NÔI).
♥️ ĐƯỜNG OTO VÀO TẬN ĐẤT TÌM CẢ XÃ CÒN ĐÚNG LÔ.
♥️ ĐẶC BIỆT CÁC BÁC THÍCH ĐI LỄ CHÙA CẦU MAY THÌ CẠNH ĐẤY LÀ CHÙA KHAI NGUYÊN LỚN NHẤT VIỆT NAM TIỆN QUÁ TRỜI LUÔN
💵 GIÁ CHỈ NHỈNH 1,3 tỷ tí RẺ NHẤT KHU VỰC Ạ.
✌️ Nhanh tay Alo e để biết thêm thông tin nhé.
🤙 Số e: ☎️ 0973954457
        """,
        "context": ""
    },
    {
        "content": """
        • 61,3m  Thôn 7, xã Tân Xã, huyện Thạch Thất, Hà Nội
• Lô đất nằm gần trục đường liên xã DH10, giao thông thuận tiện
• Gần khu vực đấu giá đất, quy hoạch đồng bộ, dân cư hiện hữu
• Chỉ mất vài phút di chuyển tới khu công nghệ cao Hòa Lạc và trường FPT
Giá bán: Chỉ nhỉnh 2 tỷ – cơ hội đầu tư sinh lời cao hoặc an cư lý tưởng
Liên hệ ngay: 0374522137
        """,
        "context": ""
    },
]

def call_api(content, context):
    """
    Gọi API analyze product
    
    Args:
        content: Nội dung cần phân tích
        context: Ngữ cảnh phân tích
        
    Returns:
        dict: Kết quả bao gồm status, response_time, response_data, error
    """
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
        
        response_time = round((time.time() - start_time) * 1000, 2)  # Convert to ms
        
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
    """
    Trích xuất response text từ AI
    """
    if not response_data:
        return ""
    
    if isinstance(response_data, dict):
        # Thử các key phổ biến
        for key in ["result", "message", "content", "response", "answer", "text"]:
            if key in response_data:
                value = response_data[key]
                if isinstance(value, str):
                    return value
                elif isinstance(value, dict):
                    return json.dumps(value, ensure_ascii=False)
        
        # Nếu không tìm thấy, trả về toàn bộ dict
        return json.dumps(response_data, ensure_ascii=False)
    
    return str(response_data)

def save_to_csv(results, filename):
    """
    Lưu kết quả vào file CSV
    
    Args:
        results: List các kết quả test
        filename: Tên file CSV
    """
    with open(filename, 'w', newline='', encoding='utf-8-sig') as csvfile:
        fieldnames = [
            'STT',
            'Content',
            'Context', 
            'Status Code',
            'Success',
            'Response Time (ms)',
            'AI Response',
            'Error',
            'Timestamp'
        ]
        
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        
        for i, result in enumerate(results, 1):
            writer.writerow({
                'STT': i,
                'Content': result['content'],
                'Context': result['context'],
                'Status Code': result['status_code'],
                'Success': 'Thành công' if result['success'] else 'Thất bại',
                'Response Time (ms)': result['response_time_ms'],
                'AI Response': result['ai_response'][:500] if result['ai_response'] else '',  # Giới hạn 500 ký tự
                'Error': result['error'] or '',
                'Timestamp': result['timestamp']
            })

def print_statistics(results):
    """
    In thống kê kết quả
    """
    total = len(results)
    success = sum(1 for r in results if r['success'])
    failed = total - success
    
    avg_response_time = sum(r['response_time_ms'] for r in results) / total if total > 0 else 0
    min_response_time = min((r['response_time_ms'] for r in results if r['response_time_ms'] > 0), default=0)
    max_response_time = max((r['response_time_ms'] for r in results), default=0)
    
    print("\n" + "="*60)
    print("THỐNG KÊ KẾT QUẢ TEST API")
    print("="*60)
    print(f"Tổng số request:        {total}")
    print(f"Thành công:             {success} ({success/total*100:.1f}%)")
    print(f"Thất bại:               {failed} ({failed/total*100:.1f}%)")
    print(f"Thời gian phản hồi TB:  {avg_response_time:.2f} ms")
    print(f"Thời gian phản hồi min: {min_response_time:.2f} ms")
    print(f"Thời gian phản hồi max: {max_response_time:.2f} ms")
    print("="*60)
    
    # Chi tiết từng request
    print("\nCHI TIẾT TỪNG REQUEST:")
    print("-"*60)
    for i, result in enumerate(results, 1):
        status = "✓" if result['success'] else "✗"
        print(f"{status} Request {i}: {result['status_code']} - {result['response_time_ms']}ms")
        print(f"  Content: {result['content'][:60]}...")
        if result['error']:
            print(f"  Error: {result['error']}")
        print()

def main():
    """
    Hàm chính
    """
    print("BẮT ĐẦU TEST API ANALYZE PRODUCT")
    print(f"API URL: {API_URL}")
    print(f"Số lượng test cases: {len(TEST_DATA)}\n")
    
    results = []
    
    for i, test_case in enumerate(TEST_DATA, 1):
        print(f"[{i}/{len(TEST_DATA)}] Đang gọi API...")
        print(f"  Content: {test_case['content'][:60]}...")
        
        api_result = call_api(
            content=test_case['content'],
            context=test_case['context']
        )
        
        result = {
            'content': test_case['content'],
            'context': test_case['context'],
            'status_code': api_result['status_code'],
            'success': api_result['success'],
            'response_time_ms': api_result['response_time_ms'],
            'ai_response': extract_ai_response(api_result['response_data']),
            'error': api_result['error'],
            'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        }
        
        results.append(result)
        
        status_icon = "✓" if result['success'] else "✗"
        print(f"  {status_icon} Status: {result['status_code']} - {result['response_time_ms']}ms")
        
        if result['error']:
            print(f"  Error: {result['error']}")
        
        print()
        
        # Delay giữa các request để tránh quá tải
        if i < len(TEST_DATA):
            time.sleep(0.5)
    
    # Lưu kết quả vào CSV
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    csv_filename = f"api_test_results_{timestamp}.csv"
    save_to_csv(results, csv_filename)
    print(f"\n✓ Đã lưu kết quả vào file: {csv_filename}")
    
    # In thống kê
    print_statistics(results)
    
    # Lưu thêm file JSON chi tiết (optional)
    json_filename = f"api_test_results_{timestamp}.json"
    with open(json_filename, 'w', encoding='utf-8') as f:
        json.dump(results, f, ensure_ascii=False, indent=2)
    print(f"✓ Đã lưu kết quả chi tiết vào file: {json_filename}")

if __name__ == "__main__":
    main()

