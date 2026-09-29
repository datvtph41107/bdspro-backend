# Hướng dẫn sử dụng Prompt hiệu quả trên Cursor

## Mục lục
- [Giới thiệu](#giới-thiệu)
- [Quy trình làm việc với Cursor](#quy-trình-làm-việc-với-cursor)
- [Các nguyên tắc cơ bản](#các-nguyên-tắc-cơ-bản)
- [Kỹ thuật prompt nâng cao](#kỹ-thuật-prompt-nâng-cao)
- [Ví dụ thực tế](#ví-dụ-thực-tế)
- [Lưu ý quan trọng](#lưu-ý-quan-trọng)

## Giới thiệu

Cursor là một IDE thông minh được tích hợp AI, cho phép bạn tương tác với AI assistant để hỗ trợ viết code, debug, và giải thích code. Để tận dụng tối đa khả năng của Cursor, bạn cần biết cách viết prompt hiệu quả và có quy trình làm việc rõ ràng.

## Quy trình làm việc với Cursor

### Bước 1: Set Context cho Project
Trước khi bắt đầu làm việc, luôn set context cho Cursor để AI hiểu rõ về project của bạn:

```
Hãy đọc và hiểu toàn bộ project này. Tôi muốn bạn nắm rõ:
- Cấu trúc thư mục và architecture
- Các services chính và chức năng
- Technology stack được sử dụng
- Database schema và models
- API patterns và conventions
```

**Ví dụ cụ thể cho từng service:**
```
Hãy đọc module organization trong organization-service và cho tôi biết:
1. Cấu trúc thư mục internal/domain/organization
2. Các models chính trong organization
3. Business logic xử lý organization
4. API endpoints có sẵn
5. Database tables liên quan
```

### Bước 2: Xác định loại task và set vai trò

#### 2.1: Thêm module hoàn toàn mới
```
Giả sử bạn là một senior Go developer với 10 năm kinh nghiệm trong microservices architecture. Tôi muốn tạo một module mới hoàn toàn cho [tên module] với các yêu cầu sau:

**Context:**
- Module này sẽ xử lý [mô tả chức năng chính]
- Tuân thủ clean architecture pattern
- Sử dụng Go 1.21+ và các best practices

**Yêu cầu cụ thể:**
- Tạo cấu trúc thư mục internal/domain/[module_name]
- Implement models, repositories, services, handlers
- Thêm validation rules
- Viết unit tests với coverage >80%
- Tích hợp với database và API gateway

Hãy implement từng phần một cách chi tiết.
```

#### 2.2: Thêm chức năng mới vào module hiện có
```
Giả sử bạn là một senior Go developer chuyên về clean architecture. Tôi muốn thêm chức năng [tên chức năng] vào module [tên module] hiện có.

**Context hiện tại:**
- Module đã có cấu trúc: [mô tả cấu trúc hiện tại]
- Database schema: [mô tả schema]
- API patterns: [mô tả patterns]

**Chức năng mới cần thêm:**
- [Mô tả chi tiết chức năng]
- Validation rules: [liệt kê rules]
- Business logic: [mô tả logic]
- API endpoint: [mô tả endpoint]

Hãy implement theo cấu trúc hiện có và đảm bảo consistency.
```

#### 2.3: Chỉnh sửa chức năng cũ của module
```
Giả sử bạn là một senior Go developer với kinh nghiệm refactoring và optimization. Tôi muốn chỉnh sửa chức năng [tên chức năng] trong module [tên module].

**Chức năng hiện tại:**
[Paste code hiện tại hoặc mô tả]

**Vấn đề cần giải quyết:**
- [Liệt kê các vấn đề]
- [Yêu cầu cải thiện]

**Yêu cầu:**
- Giữ nguyên API contract
- Cải thiện performance/security/maintainability
- Thêm proper error handling
- Update unit tests

Hãy phân tích và đưa ra giải pháp tối ưu.
```


## Các nguyên tắc cơ bản

### 1. Rõ ràng và cụ thể
```
❌ Không tốt: "Sửa lỗi này"
✅ Tốt: "Có lỗi syntax error ở dòng 15 trong file user_service.go, hãy kiểm tra và sửa lỗi"
```

### 2. Cung cấp context đầy đủ
- Mô tả vấn đề bạn đang gặp phải
- Chia sẻ code liên quan
- Giải thích mục tiêu bạn muốn đạt được

### 3. Sử dụng ngôn ngữ lập trình cụ thể
```
❌ Không tốt: "Tạo function tính tổng"
✅ Tốt: "Tạo function Go tính tổng hai số nguyên với error handling"
```

### 4. Chia nhỏ vấn đề phức tạp
```
❌ Không tốt: "Tạo toàn bộ user management system"
✅ Tốt: "Tạo model User với các fields: id, name, email, created_at"
```

## Kỹ thuật prompt nâng cao

### 1. Prompt theo vai trò (Role-based prompting)
```
Bạn là một senior Go developer với 10 năm kinh nghiệm. Hãy review code sau và đưa ra gợi ý cải thiện về performance và security.
```

### 2. Prompt theo bước (Step-by-step prompting)
```
Hãy giải thích từng bước:
1. Đầu tiên, phân tích cấu trúc code hiện tại
2. Sau đó, xác định các vấn đề tiềm ẩn
3. Cuối cùng, đề xuất giải pháp tối ưu
```

### 3. Prompt với ví dụ (Example-based prompting)
```
Tạo một API endpoint tương tự như endpoint này:
[Paste code example here]

Nhưng thay đổi để xử lý user authentication thay vì product data.
```

### 4. Prompt với ràng buộc (Constraint-based prompting)
```
Tạo function với các yêu cầu:
- Sử dụng Go 1.21+
- Tuân thủ clean architecture
- Có unit test coverage >80%
- Sử dụng dependency injection
```

## Ví dụ thực tế

### Ví dụ 1: Debug code
```
Tôi có lỗi sau trong Go service:
[Paste error message]

Code liên quan:
[Paste relevant code]

Hãy phân tích nguyên nhân và đưa ra giải pháp.
```

### Ví dụ 2: Code review
```
Hãy review code sau theo các tiêu chí:
- Code quality và readability
- Performance considerations
- Security vulnerabilities
- Best practices

[Paste code to review]
```

### Ví dụ 3: Tối ưu hóa
```
Tôi có function này chạy chậm:
[Paste function code]

Hãy phân tích performance và đề xuất cách tối ưu hóa.
```

### Ví dụ 4: Thêm tính năng mới
```
Tôi muốn thêm API endpoint để tạo organization mới. Dựa trên cấu trúc hiện tại:
- Tạo handler mới trong internal/interface/handler
- Thêm service method trong internal/usecase
- Cập nhật repository nếu cần
- Thêm validation rules
- Viết unit test

Hãy implement theo clean architecture pattern.
```

### Ví dụ 5: Refactor code
```
Tôi muốn refactor function này để dễ test hơn:
[Paste function code]

Hãy áp dụng dependency injection và tách business logic ra khỏi data access.
```

## Lưu ý quan trọng

### 1. Bảo mật
- Không chia sẻ sensitive data (passwords, API keys, etc.)
- Kiểm tra code được generate trước khi sử dụng
- Không tin tưởng hoàn toàn vào AI suggestions

### 2. Hiệu quả
- Sử dụng @ để reference files cụ thể
- Chia nhỏ vấn đề phức tạp thành nhiều prompt đơn giản
- Sử dụng chat history để duy trì context

### 3. Tương tác
- Hãy hỏi lại nếu câu trả lời không rõ ràng
- Yêu cầu giải thích thêm khi cần thiết
- Sử dụng follow-up questions để đi sâu vào chi tiết

### 4. Best Practices
- Luôn test code được generate
- Review và hiểu code trước khi sử dụng
- Kết hợp AI assistance với kiến thức của bạn
- Học từ các suggestions của AI để cải thiện kỹ năng

## Các phím tắt hữu ích

- `Cmd/Ctrl + K`: Mở chat với AI
- `Cmd/Ctrl + L`: Focus vào chat input
- `Cmd/Ctrl + I`: Inline chat (chat trong editor)
- `@`: Reference files hoặc functions

## Kết luận

Việc sử dụng prompt hiệu quả trên Cursor sẽ giúp bạn:
- Tăng productivity trong development
- Học hỏi từ AI suggestions
- Giảm thời gian debug và problem-solving
- Cải thiện code quality

Hãy thực hành thường xuyên và điều chỉnh cách viết prompt theo nhu cầu cụ thể của bạn.
