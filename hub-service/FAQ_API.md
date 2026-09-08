# FAQ API - Frequently Asked Questions

## Tổng quan

API quản lý câu hỏi thường gặp (FAQs) với đầy đủ chức năng CRUD và API đơn giản cho user.

## API Endpoints

### 1. User API (Đơn giản)

#### GET /v2/hub/faqs

Lấy danh sách FAQs đơn giản cho user, chỉ trả về id, question và answer (tối đa 50 ký tự). Hỗ trợ search theo text (lowercase, khoảng trống thay bằng %).

**Parameters:**
- `text` (query, optional): Search theo question hoặc answer (lowercase, khoảng trống → %)
- `groupKey` (query, optional): Lọc theo groupKey
- `page` (query, optional): Số trang (default: 1)
- `size` (query, optional): Số bản ghi trên trang (default: 20, max: 100)

**Search Logic:**
- Input: "tạo bài đăng"
- Lowercase: "tạo bài đăng"
- Replace spaces: "tạo%bài%đăng"
- SQL LIKE: `LOWER(question) LIKE '%tạo%bài%đăng%' OR LOWER(answer) LIKE '%tạo%bài%đăng%'`

**Response Example:**
```json
{
  "data": [
    {
      "id": 1,
      "question": "Làm thế nào để tạo bài đăng mới?",
      "answer": "Để tạo bài đăng mới, bạn vào menu Bài đăng > Tạo mới, điền thông tin và nhấn..."
    },
    {
      "id": 2,
      "question": "Cách upload hình ảnh?",
      "answer": "Bạn có thể upload hình ảnh bằng cách click vào nút Upload hoặc kéo thả file vào..."
    }
  ],
  "total": 2
}
```

### 2. Admin APIs

#### POST /v2/hub/faqs
Tạo FAQ mới.

**Request Body:**
```json
{
  "question": "Làm thế nào để tạo bài đăng mới?",
  "answer": "Để tạo bài đăng mới, bạn vào menu Bài đăng > Tạo mới, điền đầy đủ thông tin như tiêu đề, nội dung, hình ảnh và nhấn nút Lưu.",
  "groupKey": "post"
}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "question": "Làm thế nào để tạo bài đăng mới?",
    "answer": "Để tạo bài đăng mới, bạn vào menu Bài đăng > Tạo mới...",
    "groupKey": "post",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### PUT /v2/hub/faqs/{id}
Cập nhật FAQ.

**Request Body:**
```json
{
  "id": 1,
  "question": "Làm thế nào để tạo bài đăng mới? (Updated)",
  "answer": "Câu trả lời đã được cập nhật...",
  "groupKey": "post"
}
```

#### DELETE /v2/hub/faqs/{id}
Xóa FAQ (soft delete).

**Response:**
```json
{
  "success": true
}
```

#### GET /v2/hub/faqs/{id}
Lấy chi tiết một FAQ.

**Response:**
```json
{
  "data": {
    "id": 1,
    "question": "Làm thế nào để tạo bài đăng mới?",
    "answer": "Để tạo bài đăng mới, bạn vào menu Bài đăng > Tạo mới...",
    "groupKey": "post",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /v2/hub/admin/faqs
Lấy danh sách FAQs đầy đủ (Admin).

**Parameters:**
- `question` (query, optional): Tìm kiếm theo câu hỏi
- `groupKey` (query, optional): Lọc theo groupKey
- `page` (query, optional): Số trang
- `size` (query, optional): Số bản ghi trên trang

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "question": "Làm thế nào để tạo bài đăng mới?",
      "answer": "Để tạo bài đăng mới, bạn vào menu Bài đăng > Tạo mới...",
      "groupKey": "post",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

## Group Keys

Các groupKey hỗ trợ:
- `post` - FAQ về bài đăng
- `contact` - FAQ về liên hệ
- `product` - FAQ về sản phẩm
- `asset` - FAQ về tài sản
- `pipeline` - FAQ về pipeline
- `campaign` - FAQ về chiến dịch
- `general` - FAQ chung

## cURL Examples

### Lấy danh sách FAQs (User)
```bash
# Lấy tất cả
curl http://localhost:8080/v2/hub/faqs

# Lọc theo groupKey
curl http://localhost:8080/v2/hub/faqs?groupKey=post

# Search theo text
curl "http://localhost:8080/v2/hub/faqs?text=tạo%20bài%20đăng"

# Search + filter groupKey
curl "http://localhost:8080/v2/hub/faqs?text=upload&groupKey=post"

# Phân trang
curl http://localhost:8080/v2/hub/faqs?page=1&size=10
```

### Tạo FAQ mới (Admin)
```bash
curl -X POST http://localhost:8080/v2/hub/faqs \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Làm thế nào để tạo bài đăng mới?",
    "answer": "Để tạo bài đăng mới, bạn vào menu Bài đăng > Tạo mới...",
    "groupKey": "post"
  }'
```

### Cập nhật FAQ (Admin)
```bash
curl -X PUT http://localhost:8080/v2/hub/faqs/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "question": "Updated question",
    "answer": "Updated answer",
    "groupKey": "post"
  }'
```

### Xóa FAQ (Admin)
```bash
curl -X DELETE http://localhost:8080/v2/hub/faqs/1
```

## JavaScript Examples

### Frontend - Hiển thị FAQs
```javascript
// Lấy FAQs theo category và search text
async function loadFAQs(groupKey = null, searchText = null) {
  const params = new URLSearchParams();
  if (groupKey) params.append('groupKey', groupKey);
  if (searchText) params.append('text', searchText);
  
  const url = `/v2/hub/faqs?${params.toString()}`;
  const response = await fetch(url);
  const data = await response.json();
  
  return data.data;
}

// Hiển thị FAQs
const FAQList = () => {
  const [faqs, setFaqs] = useState([]);

  useEffect(() => {
    loadFAQs('post').then(setFaqs);
  }, []);

  return (
    <div className="faq-container">
      {faqs.map(faq => (
        <div key={faq.id} className="faq-item">
          <h3 className="question">{faq.question}</h3>
          <p className="answer">{faq.answer}</p>
        </div>
      ))}
    </div>
  );
};
```

### Accordion Component
```javascript
const FAQAccordion = ({ groupKey }) => {
  const [faqs, setFaqs] = useState([]);
  const [expandedId, setExpandedId] = useState(null);

  useEffect(() => {
    fetch(`/v2/hub/faqs?groupKey=${groupKey}`)
      .then(res => res.json())
      .then(data => setFaqs(data.data));
  }, [groupKey]);

  const toggleFAQ = (id) => {
    setExpandedId(expandedId === id ? null : id);
  };

  return (
    <div className="faq-accordion">
      {faqs.map(faq => (
        <div key={faq.id} className="faq-accordion-item">
          <button 
            className="faq-question"
            onClick={() => toggleFAQ(faq.id)}
          >
            {faq.question}
            <span>{expandedId === faq.id ? '−' : '+'}</span>
          </button>
          {expandedId === faq.id && (
            <div className="faq-answer">
              {faq.answer}
            </div>
          )}
        </div>
      ))}
    </div>
  );
};
```

### Search FAQs (User - Simple API)
```javascript
const FAQSearch = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [groupKey, setGroupKey] = useState('');

  const searchFAQs = async (searchText, group) => {
    const params = new URLSearchParams();
    if (searchText) params.append('text', searchText);
    if (group) params.append('groupKey', group);
    
    const response = await fetch(`/v2/hub/faqs?${params.toString()}`);
    const data = await response.json();
    setResults(data.data);
  };

  useEffect(() => {
    if (query.length >= 2) {
      const timer = setTimeout(() => searchFAQs(query, groupKey), 300);
      return () => clearTimeout(timer);
    } else if (query.length === 0) {
      searchFAQs('', groupKey); // Load all if empty
    }
  }, [query, groupKey]);

  return (
    <div>
      <select value={groupKey} onChange={(e) => setGroupKey(e.target.value)}>
        <option value="">Tất cả</option>
        <option value="post">Bài đăng</option>
        <option value="contact">Liên hệ</option>
        <option value="product">Sản phẩm</option>
      </select>
      
      <input
        type="text"
        placeholder="Tìm kiếm câu hỏi hoặc câu trả lời..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />
      
      <div className="results">
        {results.map(faq => (
          <div key={faq.id}>
            <strong>{faq.question}</strong>
            <p>{faq.answer}</p>
          </div>
        ))}
      </div>
    </div>
  );
};
```

## Mobile App Examples

### React Native
```jsx
import React, { useState, useEffect } from 'react';
import { View, Text, TouchableOpacity, FlatList } from 'react-native';

const FAQScreen = ({ groupKey }) => {
  const [faqs, setFaqs] = useState([]);
  const [expandedId, setExpandedId] = useState(null);

  useEffect(() => {
    fetch(`/v2/hub/faqs?groupKey=${groupKey}`)
      .then(res => res.json())
      .then(data => setFaqs(data.data));
  }, [groupKey]);

  const renderFAQItem = ({ item }) => (
    <View style={styles.faqItem}>
      <TouchableOpacity onPress={() => setExpandedId(item.id)}>
        <Text style={styles.question}>{item.question}</Text>
      </TouchableOpacity>
      {expandedId === item.id && (
        <Text style={styles.answer}>{item.answer}</Text>
      )}
    </View>
  );

  return (
    <FlatList
      data={faqs}
      renderItem={renderFAQItem}
      keyExtractor={item => item.id.toString()}
    />
  );
};
```

## Database Schema

```sql
CREATE TABLE faqs (
    id BIGSERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    group_key VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);

CREATE INDEX idx_faqs_group_key ON faqs (group_key);
CREATE INDEX idx_faqs_deleted_at ON faqs (deleted_at);
CREATE INDEX idx_faqs_question ON faqs USING gin(to_tsvector('english', question));
```

## Proto Definition

```protobuf
service FAQService {
  rpc CreateFAQ (CreateFAQRequest) returns (CreateFAQResponse);
  rpc UpdateFAQ (UpdateFAQRequest) returns (UpdateFAQResponse);
  rpc DeleteFAQ (DeleteFAQRequest) returns (DeleteFAQResponse);
  rpc GetFAQ (GetFAQRequest) returns (GetFAQResponse);
  rpc GetFAQList (GetFAQListRequest) returns (GetFAQListResponse);
  rpc GetFAQSimple (GetFAQSimpleRequest) returns (GetFAQSimpleResponse);
}

message FAQSimple {
  uint64 id = 1;
  string question = 2;
  string answer = 3;  // Tối đa 50 ký tự
}
```

## Features

- ✅ **CRUD Operations**: Tạo, đọc, cập nhật, xóa FAQs
- ✅ **Simple User API**: API đơn giản chỉ trả id, question, answer
- ✅ **Answer Truncation**: Tự động cắt answer xuống 50 ký tự cho API user
- ✅ **Text Search**: Search theo text (lowercase, khoảng trống → %) trong question và answer
- ✅ **Group Filter**: Lọc FAQs theo groupKey
- ✅ **Search**: Tìm kiếm FAQs theo question (Admin)
- ✅ **Pagination**: Hỗ trợ phân trang
- ✅ **Soft Delete**: Xóa mềm để có thể khôi phục
- ✅ **Full-text Search**: Index cho tìm kiếm nhanh

## Notes

- User API (`/v2/hub/faqs`) không yêu cầu authentication
- Admin APIs yêu cầu authentication và authorization
- Answer được tự động truncate xuống 50 ký tự + "..." cho user API
- Full-text search sử dụng PostgreSQL GIN index

