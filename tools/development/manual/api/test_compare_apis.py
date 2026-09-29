#!/usr/bin/env python3
"""
Script test và so sánh 2 API: Gemini vs DeepSeek
Gọi cả 2 API với cùng dữ liệu và so sánh kết quả
"""

import requests
import csv
import time
from datetime import datetime
import json
import sys

# Cấu hình API
API_GEMINI_URL = "http://localhost:8000/v2/assistant/gemini/analyze/product"
API_DEEPSEEK_URL = "http://localhost:8000/v2/assistant/deepseek/analyze/product"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpdGllcyI6bnVsbCwiZXhwIjoxNzYyNDA0MjM0LCJpYXQiOjE3NjIzOTUyMzQsIm9yZ2FuaXphdGlvbiI6bnVsbCwicGxhbiI6bnVsbCwicGxhbkF0IjpudWxsLCJwcm9maWxlIjoxNTQsInJvbGUiOiIiLCJzZXNzaW9uIjo4ODksInN1YiI6MjY1LCJ0eXBlIjoiQUNDRVNTIn0.Zkv0vLGtMQpMgZUCFxjrnCGdSr7ZJ4LRnyLGl1-KAMU"

HEADERS = {
    "accept": "application/json",
    "Authorization": AUTH_TOKEN,
    "Content-Type": "application/json"
}

# Dữ liệu mẫu để test
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

def call_api(api_url, api_name, content, context):
    """
    Gọi API analyze product
    
    Args:
        api_url: URL của API
        api_name: Tên API (để log)
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
            api_url,
            headers=HEADERS,
            json=payload,
            timeout=60  # Tăng timeout lên 60s cho API AI
        )
        
        response_time = round((time.time() - start_time) * 1000, 2)  # Convert to ms
        
        result = {
            "api_name": api_name,
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
            "api_name": api_name,
            "status_code": 0,
            "response_time_ms": round((time.time() - start_time) * 1000, 2),
            "success": False,
            "response_data": None,
            "error": "Request timeout (>60s)"
        }
    except requests.exceptions.ConnectionError:
        return {
            "api_name": api_name,
            "status_code": 0,
            "response_time_ms": 0,
            "success": False,
            "response_data": None,
            "error": "Connection error - API server không khả dụng"
        }
    except Exception as e:
        return {
            "api_name": api_name,
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
        # Thử các key phổ biến
        for key in ["result", "message", "content", "response", "answer", "text", "data"]:
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
    Lưu kết quả so sánh vào file CSV
    
    Args:
        results: List các kết quả test
        filename: Tên file CSV
    """
    with open(filename, 'w', newline='', encoding='utf-8-sig') as csvfile:
        fieldnames = [
            'STT',
            'Content (50 ký tự đầu)',
            'Context',
            
            # Gemini
            'Gemini Status',
            'Gemini Success',
            'Gemini Time (ms)',
            'Gemini Response Length',
            'Gemini Response',
            'Gemini Error',
            
            # DeepSeek
            'DeepSeek Status',
            'DeepSeek Success',
            'DeepSeek Time (ms)',
            'DeepSeek Response Length',
            'DeepSeek Response',
            'DeepSeek Error',
            
            # So sánh
            'Faster API',
            'Time Diff (ms)',
            'Both Success',
            
            'Timestamp'
        ]
        
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        
        for i, result in enumerate(results, 1):
            gemini = result['gemini']
            deepseek = result['deepseek']
            
            # Tính API nào nhanh hơn
            faster_api = ""
            time_diff = 0
            if gemini['response_time_ms'] > 0 and deepseek['response_time_ms'] > 0:
                if gemini['response_time_ms'] < deepseek['response_time_ms']:
                    faster_api = "Gemini"
                    time_diff = deepseek['response_time_ms'] - gemini['response_time_ms']
                else:
                    faster_api = "DeepSeek"
                    time_diff = gemini['response_time_ms'] - deepseek['response_time_ms']
            
            both_success = "Yes" if (gemini['success'] and deepseek['success']) else "No"
            
            writer.writerow({
                'STT': i,
                'Content (50 ký tự đầu)': result['content'][:50] + "...",
                'Context': result['context'],
                
                # Gemini
                'Gemini Status': gemini['status_code'],
                'Gemini Success': 'Yes' if gemini['success'] else 'No',
                'Gemini Time (ms)': gemini['response_time_ms'],
                'Gemini Response Length': len(gemini['ai_response']),
                'Gemini Response': gemini['ai_response'][:1000],  # Giới hạn 1000 ký tự
                'Gemini Error': gemini['error'] or '',
                
                # DeepSeek
                'DeepSeek Status': deepseek['status_code'],
                'DeepSeek Success': 'Yes' if deepseek['success'] else 'No',
                'DeepSeek Time (ms)': deepseek['response_time_ms'],
                'DeepSeek Response Length': len(deepseek['ai_response']),
                'DeepSeek Response': deepseek['ai_response'][:1000],
                'DeepSeek Error': deepseek['error'] or '',
                
                # So sánh
                'Faster API': faster_api,
                'Time Diff (ms)': round(time_diff, 2),
                'Both Success': both_success,
                
                'Timestamp': result['timestamp']
            })

def print_statistics(results):
    """In thống kê so sánh"""
    total = len(results)
    
    # Thống kê Gemini
    gemini_success = sum(1 for r in results if r['gemini']['success'])
    gemini_failed = total - gemini_success
    gemini_times = [r['gemini']['response_time_ms'] for r in results if r['gemini']['response_time_ms'] > 0]
    gemini_avg_time = sum(gemini_times) / len(gemini_times) if gemini_times else 0
    gemini_avg_length = sum(len(r['gemini']['ai_response']) for r in results) / total if total > 0 else 0
    
    # Thống kê DeepSeek
    deepseek_success = sum(1 for r in results if r['deepseek']['success'])
    deepseek_failed = total - deepseek_success
    deepseek_times = [r['deepseek']['response_time_ms'] for r in results if r['deepseek']['response_time_ms'] > 0]
    deepseek_avg_time = sum(deepseek_times) / len(deepseek_times) if deepseek_times else 0
    deepseek_avg_length = sum(len(r['deepseek']['ai_response']) for r in results) / total if total > 0 else 0
    
    # Thống kê so sánh
    both_success = sum(1 for r in results if r['gemini']['success'] and r['deepseek']['success'])
    both_failed = sum(1 for r in results if not r['gemini']['success'] and not r['deepseek']['success'])
    
    print("\n" + "="*80)
    print("THỐNG KÊ SO SÁNH GEMINI vs DEEPSEEK")
    print("="*80)
    
    print(f"\nTổng số request: {total}")
    print(f"Cả 2 API đều thành công: {both_success} ({both_success/total*100:.1f}%)")
    print(f"Cả 2 API đều thất bại: {both_failed} ({both_failed/total*100:.1f}%)")
    
    print("\n" + "-"*80)
    print("GEMINI API")
    print("-"*80)
    print(f"Thành công:             {gemini_success}/{total} ({gemini_success/total*100:.1f}%)")
    print(f"Thất bại:               {gemini_failed}/{total} ({gemini_failed/total*100:.1f}%)")
    if gemini_times:
        print(f"Thời gian phản hồi TB:  {gemini_avg_time:.2f} ms")
        print(f"Thời gian phản hồi min: {min(gemini_times):.2f} ms")
        print(f"Thời gian phản hồi max: {max(gemini_times):.2f} ms")
    print(f"Độ dài response TB:     {gemini_avg_length:.0f} ký tự")
    
    print("\n" + "-"*80)
    print("DEEPSEEK API")
    print("-"*80)
    print(f"Thành công:             {deepseek_success}/{total} ({deepseek_success/total*100:.1f}%)")
    print(f"Thất bại:               {deepseek_failed}/{total} ({deepseek_failed/total*100:.1f}%)")
    if deepseek_times:
        print(f"Thời gian phản hồi TB:  {deepseek_avg_time:.2f} ms")
        print(f"Thời gian phản hồi min: {min(deepseek_times):.2f} ms")
        print(f"Thời gian phản hồi max: {max(deepseek_times):.2f} ms")
    print(f"Độ dài response TB:     {deepseek_avg_length:.0f} ký tự")
    
    print("\n" + "-"*80)
    print("SO SÁNH TỐC ĐỘ")
    print("-"*80)
    if gemini_times and deepseek_times:
        if gemini_avg_time < deepseek_avg_time:
            faster = "Gemini"
            diff = deepseek_avg_time - gemini_avg_time
            percent = (diff / deepseek_avg_time) * 100
        else:
            faster = "DeepSeek"
            diff = gemini_avg_time - deepseek_avg_time
            percent = (diff / gemini_avg_time) * 100
        
        print(f"API nhanh hơn:          {faster}")
        print(f"Nhanh hơn trung bình:   {diff:.2f} ms ({percent:.1f}%)")
    
    print("\n" + "-"*80)
    print("SO SÁNH ĐỘ DÀI RESPONSE")
    print("-"*80)
    if gemini_avg_length > deepseek_avg_length:
        print(f"Gemini dài hơn:         {gemini_avg_length - deepseek_avg_length:.0f} ký tự")
    else:
        print(f"DeepSeek dài hơn:       {deepseek_avg_length - gemini_avg_length:.0f} ký tự")
    
    print("\n" + "="*80)
    
    # Chi tiết từng request
    print("\nCHI TIẾT TỪNG REQUEST:")
    print("-"*80)
    for i, result in enumerate(results, 1):
        print(f"\n[{i}] {result['content'][:70]}...")
        
        gemini_status = "✓" if result['gemini']['success'] else "✗"
        deepseek_status = "✓" if result['deepseek']['success'] else "✗"
        
        print(f"  {gemini_status} Gemini:   {result['gemini']['status_code']} - {result['gemini']['response_time_ms']}ms - {len(result['gemini']['ai_response'])} chars")
        if result['gemini']['error']:
            print(f"    Error: {result['gemini']['error']}")
        
        print(f"  {deepseek_status} DeepSeek: {result['deepseek']['status_code']} - {result['deepseek']['response_time_ms']}ms - {len(result['deepseek']['ai_response'])} chars")
        if result['deepseek']['error']:
            print(f"    Error: {result['deepseek']['error']}")

def main():
    """Hàm chính"""
    print("="*80)
    print("BẮT ĐẦU TEST VÀ SO SÁNH API: GEMINI vs DEEPSEEK")
    print("="*80)
    print(f"Gemini API:   {API_GEMINI_URL}")
    print(f"DeepSeek API: {API_DEEPSEEK_URL}")
    print(f"Số lượng test cases: {len(TEST_DATA)}")
    print("="*80)
    print()
    
    results = []
    
    for i, test_case in enumerate(TEST_DATA, 1):
        print(f"[{i}/{len(TEST_DATA)}] Testing với content:")
        content_preview = test_case['content'][:70].replace('\n', ' ') + "..."
        print(f"  {content_preview}")
        
        # Gọi Gemini API
        print(f"  → Calling Gemini API...")
        gemini_result = call_api(
            api_url=API_GEMINI_URL,
            api_name="Gemini",
            content=test_case['content'],
            context=test_case['context']
        )
        gemini_status = "✓" if gemini_result['success'] else "✗"
        print(f"    {gemini_status} {gemini_result['status_code']} - {gemini_result['response_time_ms']}ms")
        
        # Delay nhỏ giữa 2 API calls
        time.sleep(0.3)
        
        # Gọi DeepSeek API
        print(f"  → Calling DeepSeek API...")
        deepseek_result = call_api(
            api_url=API_DEEPSEEK_URL,
            api_name="DeepSeek",
            content=test_case['content'],
            context=test_case['context']
        )
        deepseek_status = "✓" if deepseek_result['success'] else "✗"
        print(f"    {deepseek_status} {deepseek_result['status_code']} - {deepseek_result['response_time_ms']}ms")
        
        # Lưu kết quả
        result = {
            'content': test_case['content'],
            'context': test_case['context'],
            'gemini': {
                'status_code': gemini_result['status_code'],
                'success': gemini_result['success'],
                'response_time_ms': gemini_result['response_time_ms'],
                'ai_response': extract_ai_response(gemini_result['response_data']),
                'error': gemini_result['error'],
                'full_response_data': gemini_result['response_data']
            },
            'deepseek': {
                'status_code': deepseek_result['status_code'],
                'success': deepseek_result['success'],
                'response_time_ms': deepseek_result['response_time_ms'],
                'ai_response': extract_ai_response(deepseek_result['response_data']),
                'error': deepseek_result['error'],
                'full_response_data': deepseek_result['response_data']
            },
            'timestamp': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        }
        
        results.append(result)
        print()
        
        # Delay giữa các test case
        if i < len(TEST_DATA):
            time.sleep(1)
    
    # Lưu kết quả vào CSV
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    csv_filename = f"api_comparison_{timestamp}.csv"
    save_to_csv(results, csv_filename)
    print(f"✓ Đã lưu kết quả so sánh vào file CSV: {csv_filename}")
    
    # Lưu file JSON chi tiết
    json_filename = f"api_comparison_{timestamp}.json"
    with open(json_filename, 'w', encoding='utf-8') as f:
        json.dump(results, f, ensure_ascii=False, indent=2)
    print(f"✓ Đã lưu kết quả chi tiết vào file JSON: {json_filename}")
    
    # In thống kê
    print_statistics(results)
    
    print(f"\n{'='*80}")
    print("✓ HOÀN THÀNH SO SÁNH!")
    print(f"{'='*80}")

if __name__ == "__main__":
    main()

