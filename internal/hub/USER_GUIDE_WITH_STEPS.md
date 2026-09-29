# User Guide với Steps - Implementation Guide

## Tổng quan

User Guide đã được cập nhật để hỗ trợ **Steps** (các bước hướng dẫn). Mỗi step bao gồm:
- **image**: URL hình ảnh minh họa
- **content**: Nội dung mô tả bước
- **stepOrder**: Thứ tự bước (1, 2, 3...)

## Database Schema

### Bảng `user_guides` (existing)
```sql
CREATE TABLE user_guides (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    group_key VARCHAR(100),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT
);
```

### Bảng `user_guide_steps` (NEW)
```sql
CREATE TABLE user_guide_steps (
    id BIGSERIAL PRIMARY KEY,
    user_guide_id BIGINT NOT NULL,  -- Foreign key to user_guides
    step_order INT NOT NULL,         -- Thứ tự bước (1, 2, 3...)
    image TEXT,                      -- URL hình ảnh
    content TEXT NOT NULL,           -- Nội dung bước
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_by BIGINT,
    updated_by BIGINT,
    
    FOREIGN KEY(user_guide_id) REFERENCES user_guides(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_user_guide_steps_user_guide_id ON user_guide_steps (user_guide_id);
CREATE INDEX idx_user_guide_steps_step_order ON user_guide_steps (step_order);
```

## Proto Definitions

```protobuf
message UserGuideStep {
  uint64 id = 1;
  uint64 userGuideId = 2;
  int32 stepOrder = 3;
  string image = 4;
  string content = 5;
}

message UserGuideStepInput {
  int32 stepOrder = 1;
  string image = 2;
  string content = 3;
}

message UserGuide {
  uint64 id = 1;
  string title = 2;
  string description = 3;
  string groupKey = 4;
  string createdAt = 5;
  string updatedAt = 6;
  repeated UserGuideStep steps = 7;  // ✅ Thêm steps
}

message CreateUserGuideRequest {
  string title = 1;
  string description = 2;
  string groupKey = 3;
  repeated UserGuideStepInput steps = 4;  // ✅ Thêm steps
}

message UpdateUserGuideRequest {
  uint64 id = 1;
  string title = 2;
  string description = 3;
  string groupKey = 4;
  repeated UserGuideStepInput steps = 5;  // ✅ Thêm steps
}
```

## API Examples

### 1. Tạo User Guide với Steps

**Request:**
```json
POST /v2/hub/user-guides
{
  "title": "Hướng dẫn tạo bài đăng",
  "description": "Cách tạo bài đăng mới trên hệ thống",
  "groupKey": "post",
  "steps": [
    {
      "stepOrder": 1,
      "image": "https://example.com/step1.png",
      "content": "Bước 1: Click vào nút 'Tạo bài đăng' ở góc trên bên phải"
    },
    {
      "stepOrder": 2,
      "image": "https://example.com/step2.png",
      "content": "Bước 2: Điền tiêu đề và nội dung cho bài đăng của bạn"
    },
    {
      "stepOrder": 3,
      "image": "https://example.com/step3.png",
      "content": "Bước 3: Upload hình ảnh minh họa (tùy chọn)"
    },
    {
      "stepOrder": 4,
      "image": "",
      "content": "Bước 4: Click nút 'Đăng' để hoàn tất"
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Hướng dẫn tạo bài đăng",
    "description": "Cách tạo bài đăng mới trên hệ thống",
    "groupKey": "post",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "steps": [
      {
        "id": 1,
        "userGuideId": 1,
        "stepOrder": 1,
        "image": "https://example.com/step1.png",
        "content": "Bước 1: Click vào nút 'Tạo bài đăng' ở góc trên bên phải"
      },
      {
        "id": 2,
        "userGuideId": 1,
        "stepOrder": 2,
        "image": "https://example.com/step2.png",
        "content": "Bước 2: Điền tiêu đề và nội dung cho bài đăng của bạn"
      },
      {
        "id": 3,
        "userGuideId": 1,
        "stepOrder": 3,
        "image": "https://example.com/step3.png",
        "content": "Bước 3: Upload hình ảnh minh họa (tùy chọn)"
      },
      {
        "id": 4,
        "userGuideId": 1,
        "stepOrder": 4,
        "image": "",
        "content": "Bước 4: Click nút 'Đăng' để hoàn tất"
      }
    ]
  }
}
```

### 2. Cập nhật User Guide và Steps

**Request:**
```json
PUT /v2/hub/user-guides/1
{
  "id": 1,
  "title": "Hướng dẫn tạo bài đăng (Updated)",
  "description": "Cách tạo bài đăng mới - phiên bản cập nhật",
  "groupKey": "post",
  "steps": [
    {
      "stepOrder": 1,
      "image": "https://example.com/new-step1.png",
      "content": "Bước 1: Vào menu Bài đăng > Tạo mới"
    },
    {
      "stepOrder": 2,
      "image": "https://example.com/new-step2.png",
      "content": "Bước 2: Nhập thông tin bài đăng"
    }
  ]
}
```

**Note:** Khi update, tất cả steps cũ sẽ bị xóa và thay thế bằng steps mới.

### 3. Lấy chi tiết User Guide với Steps

**Request:**
```bash
GET /v2/hub/user-guides/1
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "title": "Hướng dẫn tạo bài đăng",
    "description": "Cách tạo bài đăng mới trên hệ thống",
    "groupKey": "post",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "steps": [
      {
        "id": 1,
        "userGuideId": 1,
        "stepOrder": 1,
        "image": "https://example.com/step1.png",
        "content": "Bước 1: Click vào nút 'Tạo bài đăng' ở góc trên bên phải"
      },
      // ... more steps
    ]
  }
}
```

## Frontend Integration

### React Component Example

```jsx
const UserGuideDetail = ({ guideId }) => {
  const [guide, setGuide] = useState(null);

  useEffect(() => {
    fetch(`/v2/hub/user-guides/${guideId}`)
      .then(res => res.json())
      .then(data => setGuide(data.data));
  }, [guideId]);

  if (!guide) return <Loading />;

  return (
    <div className="user-guide">
      <h1>{guide.title}</h1>
      <p>{guide.description}</p>
      
      <div className="steps">
        {guide.steps.map((step, index) => (
          <div key={step.id} className="step">
            <div className="step-number">Bước {step.stepOrder}</div>
            {step.image && (
              <img src={step.image} alt={`Step ${step.stepOrder}`} />
            )}
            <p>{step.content}</p>
          </div>
        ))}
      </div>
    </div>
  );
};
```

### Create/Update Form

```jsx
const UserGuideForm = ({ guideId = null }) => {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    groupKey: 'post',
    steps: [
      { stepOrder: 1, image: '', content: '' }
    ]
  });

  const addStep = () => {
    setFormData({
      ...formData,
      steps: [
        ...formData.steps,
        { stepOrder: formData.steps.length + 1, image: '', content: '' }
      ]
    });
  };

  const removeStep = (index) => {
    const newSteps = formData.steps.filter((_, i) => i !== index);
    // Reorder
    newSteps.forEach((step, i) => {
      step.stepOrder = i + 1;
    });
    setFormData({ ...formData, steps: newSteps });
  };

  const updateStep = (index, field, value) => {
    const newSteps = [...formData.steps];
    newSteps[index][field] = value;
    setFormData({ ...formData, steps: newSteps });
  };

  const handleSubmit = async () => {
    const url = guideId 
      ? `/v2/hub/user-guides/${guideId}`
      : '/v2/hub/user-guides';
    const method = guideId ? 'PUT' : 'POST';

    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData)
    });

    const data = await response.json();
    console.log('Saved:', data);
  };

  return (
    <form onSubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      <input
        value={formData.title}
        onChange={(e) => setFormData({ ...formData, title: e.target.value })}
        placeholder="Tiêu đề"
      />
      
      <textarea
        value={formData.description}
        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
        placeholder="Mô tả"
      />

      <h3>Các bước:</h3>
      {formData.steps.map((step, index) => (
        <div key={index} className="step-form">
          <h4>Bước {step.stepOrder}</h4>
          <input
            value={step.image}
            onChange={(e) => updateStep(index, 'image', e.target.value)}
            placeholder="URL hình ảnh"
          />
          <textarea
            value={step.content}
            onChange={(e) => updateStep(index, 'content', e.target.value)}
            placeholder="Nội dung bước"
          />
          <button type="button" onClick={() => removeStep(index)}>Xóa</button>
        </div>
      ))}

      <button type="button" onClick={addStep}>Thêm bước</button>
      <button type="submit">Lưu</button>
    </form>
  );
};
```

## Files Changed

### Domain Layer
- ✅ `hub-service/internal/domain/user_guide_step_entity.go` - Entity cho steps

### Repository Layer
- ✅ `hub-service/internal/repo/user_guide_step_repo.go` - Interface
- ✅ `hub-service/infra/postgre/user_guide_step_postgres.go` - Implementation

### Usecase Layer
- ✅ `hub-service/internal/usecase/user_guide_usecase.go` - Updated constructor

### Handler Layer
- ✅ `hub-service/infra/handler/user_guide_handler.go` - Updated Create/Update/GetDetail

### Mapper Layer
- ✅ `hub-service/infra/mapper/user_guide_mapper.go` - Step conversion methods

### Proto Layer
- ✅ `shared/protobuf/schema/hub/user_guide.proto` - Added Steps messages

### Migration
- ✅ `hub-service/migrate/003_create_user_guide_steps_table.sql` - DB migration

## Migration Steps

1. **Run SQL Migration:**
```bash
psql -U your_user -d your_db -f hub-service/migrate/003_create_user_guide_steps_table.sql
```

2. **Generate Proto Code:**
```bash
cd shared/protobuf/schema/hub
protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_opt=paths=source_relative \
  user_guide.proto
```

3. **Update Wire Dependencies:**
- Thêm `NewUserGuideStepRepo` vào wire.go
- Regenerate: `cd hub-service/wire && wire`

## Features

- ✅ **Ordered Steps**: Steps được sắp xếp theo stepOrder
- ✅ **Image Support**: Mỗi step có thể có hình ảnh minh họa
- ✅ **Cascade Delete**: Xóa user guide sẽ tự động xóa tất cả steps
- ✅ **Update Logic**: Update user guide sẽ xóa và tạo lại steps mới
- ✅ **Clean Architecture**: Tách biệt rõ ràng giữa các layers

## Notes

- Steps được lưu riêng trong bảng `user_guide_steps`
- Foreign key constraint đảm bảo data integrity
- ON DELETE CASCADE tự động xóa steps khi xóa user guide
- Steps được load cùng với user guide detail
- Simple list API không trả về steps (chỉ có title và description)

