# 🎨 Test UI - Hướng dẫn sử dụng

## 📝 Mô tả

File `test_ui.html` là giao diện web đơn giản để test AI Service một cách trực quan, không cần dùng curl hay Postman.

## 🚀 Cách sử dụng

### Bước 1: Chạy AI Service

```bash
cd ai-service
source venv/bin/activate
python main.py
```

Service sẽ chạy tại: `http://localhost:8040`

### Bước 2: Mở Test UI

Có 2 cách:

**Cách 1: Mở trực tiếp file**
```bash
open test_ui.html
```

**Cách 2: Double-click file trong Finder**
- Tìm file `test_ui.html` trong thư mục `ai-service`
- Double-click để mở trong trình duyệt

### Bước 3: Sử dụng

1. **Kiểm tra kết nối**: Thanh status ở trên cùng sẽ hiển thị:
   - ✅ **Màu xanh**: API online
   - ❌ **Màu đỏ**: API offline (cần start service)

2. **Nhập câu hỏi**: 
   - Ví dụ: "Làm sao để tạo một API mới trong BDSPro?"
   - Ví dụ: "Kiến trúc của auth-service như thế nào?"

3. **Tùy chỉnh tham số**:
   - `top_k`: Số chunks context (mặc định: 5)
   - `temperature`: Độ sáng tạo của AI (0-2, mặc định: 0.3)

4. **Submit**:
   - Click "Gửi câu hỏi"
   - Hoặc nhấn `Ctrl/Cmd + Enter`

5. **Xem kết quả**:
   - Câu trả lời chi tiết
   - Thời gian xử lý
   - Số chunks được sử dụng
   - Độ tin cậy (confidence)
   - Danh sách nguồn tham khảo

## ✨ Tính năng

- ✅ Giao diện đẹp, hiện đại
- ✅ Responsive (mobile-friendly)
- ✅ Kiểm tra kết nối API tự động
- ✅ Loading indicator khi xử lý
- ✅ Hiển thị meta info (time, chunks, confidence)
- ✅ Keyboard shortcut (Ctrl/Cmd + Enter)
- ✅ Error handling với hướng dẫn khắc phục
- ✅ CORS-ready

## 🔧 Thay đổi API URL

Nếu service chạy ở port khác, sửa dòng này trong file HTML:

```javascript
const API_URL = 'http://localhost:8040';  // Đổi thành port của bạn
```

## 🎯 Các câu hỏi mẫu để test

```
1. Làm sao để tạo một API mới theo chuẩn v2?
2. Giải thích cấu trúc Clean Architecture của project
3. Auth service hoạt động như thế nào?
4. Cách sử dụng wire để dependency injection?
5. Làm sao để kết nối với database PostgreSQL?
6. Quy tắc đặt tên biến và hàm trong project?
7. Làm sao để thêm một entity mới?
8. Cách triển khai repository pattern?
9. Giải thích về infra layer
10. Làm sao để gọi API từ service khác?
```

## 🐛 Troubleshooting

### Lỗi: "Không thể kết nối đến API"

**Nguyên nhân**: AI service chưa chạy

**Giải pháp**:
```bash
cd ai-service
source venv/bin/activate
python main.py
```

### Lỗi: "No relevant context found"

**Nguyên nhân**: Vector DB chưa có dữ liệu

**Giải pháp**:
```bash
python scripts/embed_context.py
```

### Lỗi CORS

**Nguyên nhân**: CORS chưa được cấu hình đúng

**Giải pháp**: Đã được cấu hình sẵn trong `main.py`:
```python
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
```

## 📸 Screenshot

Giao diện bao gồm:
- Header gradient màu tím đẹp mắt
- Form nhập câu hỏi với textarea lớn
- Tùy chỉnh top_k và temperature
- Hiển thị kết quả với:
  - Box câu trả lời (màu xanh nhạt)
  - Meta cards hiển thị số liệu
  - Danh sách nguồn tham khảo (màu vàng nhạt)
- Footer thông tin

## 💡 Tips

1. **Câu hỏi cụ thể**: Hỏi cụ thể sẽ cho kết quả tốt hơn
2. **Tăng top_k**: Nếu câu trả lời thiếu context, tăng top_k lên 8-10
3. **Temperature**: 
   - Thấp (0.1-0.3): Câu trả lời chính xác, ít sáng tạo
   - Cao (0.7-1.0): Câu trả lời sáng tạo hơn nhưng có thể sai
4. **Refresh**: Nếu API status không đổi, refresh trang

## 🔗 Liên quan

- `main.py`: API endpoint implementation
- `LOCAL_MODE_README.md`: Hướng dẫn chạy AI service
- `scripts/embed_context.py`: Script tạo vector embeddings
- `knowledge_base/`: Dữ liệu training

---

**Enjoy testing! 🎉**

