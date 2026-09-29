# User Guide Simple API - Simplified User Guide List

## Tổng quan

API đơn giản cho user để lấy danh sách user guides, chỉ trả về id, title và description (tối đa 80 ký tự).

## Khác biệt với UserGuide API đầy đủ

### UserGuide API đầy đủ (Admin):
```
GET /v2/hub/user-guides

Response:
{
  "data": [
    {
      "id": 1,
      "title": "Hướng dẫn tạo bài đăng",
      "description": "Chi tiết cách tạo bài đăng mới với đầy đủ thông tin...",
      "groupKey": "post",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

### UserGuide Simple API (User):
```
GET /v2/hub/user-guides/simple

Response:
{
  "data": [
    {
      "id": 1,
      "title": "Hướng dẫn tạo bài đăng",
      "description": "Chi tiết cách tạo bài đăng mới với đầy đủ thông tin..."
    }
  ],
  "total": 1
}
```

## API Endpoint

### GET /v2/hub/user-guides/simple

Lấy danh sách user guides đơn giản với phân trang và filter theo groupKey.

**Parameters:**
- `groupKey` (query, optional): Lọc theo groupKey
- `page` (query, optional): Số trang (default: 1)
- `size` (query, optional): Số bản ghi trên trang (default: 20, max: 100)

**Các groupKey hỗ trợ:**
- `post` - Hướng dẫn bài đăng
- `contact` - Hướng dẫn liên hệ
- `product` - Hướng dẫn sản phẩm
- `asset` - Hướng dẫn tài sản
- `pipeline` - Hướng dẫn pipeline
- `campaign` - Hướng dẫn chiến dịch
- `general` - Hướng dẫn chung

## Request Examples

### cURL
```bash
# Lấy tất cả user guides
curl -X GET "http://localhost:8080/v2/hub/user-guides/simple" \
  -H "Accept: application/json"

# Lọc theo groupKey
curl -X GET "http://localhost:8080/v2/hub/user-guides/simple?groupKey=post" \
  -H "Accept: application/json"

# Phân trang
curl -X GET "http://localhost:8080/v2/hub/user-guides/simple?page=1&size=10" \
  -H "Accept: application/json"
```

### JavaScript (Fetch)
```javascript
// Lấy user guides đơn giản
async function getUserGuidesSimple(groupKey = null, page = 1, size = 20) {
  const params = new URLSearchParams();
  if (groupKey) params.append('groupKey', groupKey);
  params.append('page', page.toString());
  params.append('size', size.toString());
  
  const response = await fetch(`/v2/hub/user-guides/simple?${params}`);
  const data = await response.json();
  
  return data;
}

// Usage
const guides = await getUserGuidesSimple('post', 1, 10);
console.log(guides.data); // Array of simple user guides
```

### Go Client
```go
import hubpb "pb/types/hub"

req := &hubpb.GetUserGuideSimpleRequest{
    GroupKey: "post",
    Page:     1,
    Size:     20,
}

resp, err := client.GetUserGuideSimple(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Access simple guides
for _, guide := range resp.Data {
    fmt.Printf("ID: %d, Title: %s, Description: %s\n", 
        guide.Id, guide.Title, guide.Description)
}
```

## Response Examples

### Success Response (200)
```json
{
  "data": [
    {
      "id": 1,
      "title": "Hướng dẫn tạo bài đăng",
      "description": "Chi tiết cách tạo bài đăng mới với đầy đủ thông tin về tiêu đề, nội dung..."
    },
    {
      "id": 2,
      "title": "Cách upload hình ảnh",
      "description": "Hướng dẫn chi tiết cách upload và quản lý hình ảnh cho bài đăng của bạn..."
    },
    {
      "id": 3,
      "title": "Thiết lập thông báo",
      "description": "Cách thiết lập và quản lý các loại thông báo trong hệ thống để không bỏ lỡ..."
    }
  ],
  "total": 3
}
```

### Success Response với phân trang
```json
{
  "data": [
    {
      "id": 1,
      "title": "Hướng dẫn tạo bài đăng",
      "description": "Chi tiết cách tạo bài đăng mới với đầy đủ thông tin về tiêu đề, nội dung..."
    }
  ],
  "total": 15
}
```

### Error Response (500)
```json
{
  "error": "failed to get simple user guides: Lỗi khi lấy danh sách user guide đơn giản: database connection failed"
}
```

## Usage Scenarios

### 1. Hiển thị danh sách hướng dẫn trong mobile app
```javascript
// React Native example
const UserGuideList = () => {
  const [guides, setGuides] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadGuides = async () => {
      try {
        const response = await fetch('/v2/hub/user-guides/simple?size=50');
        const data = await response.json();
        setGuides(data.data);
      } catch (error) {
        console.error('Error loading guides:', error);
      } finally {
        setLoading(false);
      }
    };

    loadGuides();
  }, []);

  if (loading) return <LoadingSpinner />;

  return (
    <FlatList
      data={guides}
      keyExtractor={(item) => item.id.toString()}
      renderItem={({ item }) => (
        <TouchableOpacity onPress={() => openGuideDetail(item.id)}>
          <View style={styles.guideItem}>
            <Text style={styles.title}>{item.title}</Text>
            <Text style={styles.description}>{item.description}</Text>
          </View>
        </TouchableOpacity>
      )}
    />
  );
};
```

### 2. Filter theo category
```javascript
const GuideCategory = ({ category }) => {
  const [guides, setGuides] = useState([]);

  useEffect(() => {
    const loadCategoryGuides = async () => {
      const response = await fetch(`/v2/hub/user-guides/simple?groupKey=${category}`);
      const data = await response.json();
      setGuides(data.data);
    };

    loadCategoryGuides();
  }, [category]);

  return (
    <div>
      <h2>{category} Guides</h2>
      {guides.map(guide => (
        <div key={guide.id} className="guide-card">
          <h3>{guide.title}</h3>
          <p>{guide.description}</p>
          <button onClick={() => viewFullGuide(guide.id)}>
            Xem chi tiết
          </button>
        </div>
      ))}
    </div>
  );
};
```

### 3. Infinite scroll với pagination
```javascript
const InfiniteGuideList = () => {
  const [guides, setGuides] = useState([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const loadMoreGuides = async () => {
    if (!hasMore) return;

    const response = await fetch(`/v2/hub/user-guides/simple?page=${page}&size=20`);
    const data = await response.json();
    
    if (data.data.length === 0) {
      setHasMore(false);
      return;
    }

    setGuides(prev => [...prev, ...data.data]);
    setPage(prev => prev + 1);
  };

  return (
    <InfiniteScroll
      dataLength={guides.length}
      next={loadMoreGuides}
      hasMore={hasMore}
      loader={<LoadingSpinner />}
    >
      {guides.map(guide => (
        <GuideCard key={guide.id} guide={guide} />
      ))}
    </InfiniteScroll>
  );
};
```

### 4. Search và filter
```javascript
const GuideSearch = () => {
  const [guides, setGuides] = useState([]);
  const [selectedGroup, setSelectedGroup] = useState('');

  const searchGuides = async (groupKey = '') => {
    const url = groupKey 
      ? `/v2/hub/user-guides/simple?groupKey=${groupKey}`
      : '/v2/hub/user-guides/simple';
    
    const response = await fetch(url);
    const data = await response.json();
    setGuides(data.data);
  };

  return (
    <div>
      <select 
        value={selectedGroup} 
        onChange={(e) => {
          setSelectedGroup(e.target.value);
          searchGuides(e.target.value);
        }}
      >
        <option value="">Tất cả</option>
        <option value="post">Bài đăng</option>
        <option value="contact">Liên hệ</option>
        <option value="product">Sản phẩm</option>
        <option value="asset">Tài sản</option>
      </select>

      <div className="guide-list">
        {guides.map(guide => (
          <div key={guide.id} className="guide-item">
            <h3>{guide.title}</h3>
            <p>{guide.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
};
```

## Proto Definition

```protobuf
service UserGuideService {
    // API cho user - Lấy danh sách user guides đơn giản (chỉ id, title, description tối đa 80 ký tự)
    rpc GetUserGuideSimple (GetUserGuideSimpleRequest) returns (GetUserGuideSimpleResponse) {
        option (google.api.http) = {
            get: "/v2/hub/user-guides/simple"
        };
    };
}

message GetUserGuideSimpleRequest {
    string groupKey = 1;   // Filter by groupKey (optional)
    int32 page = 2;        // Page number (default: 1)
    int32 size = 3;        // Page size (default: 20)
}

message UserGuideSimple {
    uint64 id = 1;
    string title = 2;
    string description = 3;  // Tối đa 80 ký tự
}

message GetUserGuideSimpleResponse {
    repeated UserGuideSimple data = 1;
    int64 total = 2;
}
```

## Implementation Details

### Usecase
```go
// GetSimpleList retrieves user guides for simple API (only id, title, description)
func (u *UserGuideUsecase) GetSimpleList(ctx context.Context, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, *_err.ErrorDTO) {
    data, total, err := u.Repo.GetListWithFilter(ctx, "", groupKey, pagable)
    if err != nil {
        return nil, 0, &_err.ErrorDTO{
            Code:    500,
            Message: "Lỗi khi lấy danh sách user guide đơn giản: " + err.Error(),
        }
    }

    return data, total, nil
}
```

### Handler
```go
func (h *UserGuideHandler) GetUserGuideSimple(ctx context.Context, req *hubpb.GetUserGuideSimpleRequest) (*hubpb.GetUserGuideSimpleResponse, error) {
    // Create pagable from request
    pagable := &_dto.Pagable{
        Page: uint32(req.GetPage()),
        Size: uint32(req.GetSize()),
    }

    // Get simple list
    data, total, err := h.userGuideUsecase.GetSimpleList(
        ctx,
        req.GetGroupKey(),
        pagable,
    )
    if err != nil {
        return nil, err
    }

    // Convert to simple proto format
    var simpleItems []*hubpb.UserGuideSimple
    for _, item := range data {
        description := item.Description
        // Truncate description to max 80 characters
        if len(description) > 80 {
            description = description[:80] + "..."
        }

        simpleItems = append(simpleItems, &hubpb.UserGuideSimple{
            Id:          item.ID,
            Title:       item.Title,
            Description: description,
        })
    }

    return &hubpb.GetUserGuideSimpleResponse{
        Data:  simpleItems,
        Total: total,
    }, nil
}
```

## Benefits

1. **Lightweight Response**: Chỉ trả về data cần thiết (id, title, description)
2. **Description Truncation**: Tự động cắt description xuống 80 ký tự
3. **Easy to Parse**: Dễ parse và hiển thị trong UI
4. **Performance**: Giảm kích thước response đáng kể
5. **Mobile Friendly**: Tối ưu cho mobile app với data nhẹ
6. **Pagination Support**: Hỗ trợ phân trang để load từng phần

## Notes

- API này **không yêu cầu authentication** (public)
- Description được tự động truncate xuống 80 ký tự + "..."
- Nếu description <= 80 ký tự thì giữ nguyên
- Pagination: default page=1, size=20, max size=100
- GroupKey filter là optional, nếu không có thì lấy tất cả
- Response luôn có `total` để biết tổng số records

