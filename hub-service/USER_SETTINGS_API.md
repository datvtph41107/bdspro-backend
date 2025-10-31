# User Settings API - Simplified Key-Value Settings

## Tổng quan

API đơn giản cho user để lấy settings theo group, chỉ trả về key-value pairs (không có metadata như id, createdAt, etc).

## Khác biệt với SystemConfig API

### SystemConfig API (Admin):
```
GET /v2/hub/system-config/by-group/{groupKey}

Response:
{
  "data": [
    {
      "id": 1,
      "name": "Tên hiển thị",
      "key": "display_name",
      "value": "true",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z",
      "groupKey": "post",
      "groupName": "Cấu hình bài đăng"
    }
  ],
  "total": 1,
  "groupName": "Cấu hình bài đăng"
}
```

### User Settings API (User):
```
GET /v2/hub/settings/{groupKey}

Response:
{
  "settings": {
    "display_name": "true",
    "auto_save": "false",
    "max_upload": "10"
  },
  "groupKey": "post"
}
```

## API Endpoint

### GET /v2/hub/settings/{groupKey}

Lấy settings theo group dưới dạng key-value map đơn giản.

**Parameters:**
- `groupKey` (path, required): Group key của settings

**Các groupKey hỗ trợ:**
- `post` - Cấu hình bài đăng
- `contact` - Cấu hình liên hệ
- `product` - Cấu hình sản phẩm
- `asset` - Cấu hình tài sản
- `notification` - Cấu hình thông báo
- `pipeline` - Cấu hình pipeline
- `campaign` - Cấu hình chiến dịch
- `general` - Cấu hình chung

## Request Examples

### cURL
```bash
curl -X GET "http://localhost:8080/v2/hub/settings/post" \
  -H "Accept: application/json"
```

### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/v2/hub/settings/post')
  .then(response => response.json())
  .then(data => {
    console.log(data.settings);
    // Output: { "display_name": "true", "auto_save": "false", ... }
  });
```

### Go Client
```go
import hubpb "pb/types/hub"

req := &hubpb.GetUserSettingsByGroupRequest{
    GroupKey: "post",
}

resp, err := client.GetUserSettingsByGroup(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Access settings
displayName := resp.Settings["display_name"]
autoSave := resp.Settings["auto_save"]
```

## Response Examples

### Success Response (200)
```json
{
  "settings": {
    "display_name": "true",
    "auto_save": "false",
    "max_upload": "10",
    "allow_comments": "true",
    "default_status": "draft"
  },
  "groupKey": "post"
}
```

### Error Response (400) - Invalid Group Key
```json
{
  "error": "failed to get user settings: Group key không hợp lệ: invalid_group"
}
```

### Error Response (500) - Internal Error
```json
{
  "error": "failed to get user settings: Lỗi khi lấy settings"
}
```

## Usage Scenarios

### 1. Lấy settings cho màn hình Post
```javascript
// Frontend code
async function loadPostSettings() {
  const response = await fetch('/v2/hub/settings/post');
  const { settings } = await response.json();
  
  // Apply settings
  const maxUpload = parseInt(settings.max_upload || '5');
  const allowComments = settings.allow_comments === 'true';
  const defaultStatus = settings.default_status || 'draft';
  
  return { maxUpload, allowComments, defaultStatus };
}
```

### 2. Cache settings trong client
```javascript
class SettingsManager {
  constructor() {
    this.cache = {};
  }
  
  async getSettings(groupKey) {
    if (!this.cache[groupKey]) {
      const response = await fetch(`/v2/hub/settings/${groupKey}`);
      const data = await response.json();
      this.cache[groupKey] = data.settings;
    }
    return this.cache[groupKey];
  }
  
  getSetting(groupKey, key, defaultValue = null) {
    const settings = this.cache[groupKey] || {};
    return settings[key] || defaultValue;
  }
}

// Usage
const manager = new SettingsManager();
await manager.getSettings('post');
const maxUpload = manager.getSetting('post', 'max_upload', '5');
```

### 3. Load multiple groups
```javascript
async function loadAllSettings() {
  const groups = ['post', 'contact', 'product', 'asset'];
  
  const responses = await Promise.all(
    groups.map(group => 
      fetch(`/v2/hub/settings/${group}`).then(r => r.json())
    )
  );
  
  const allSettings = {};
  responses.forEach(({ settings, groupKey }) => {
    allSettings[groupKey] = settings;
  });
  
  return allSettings;
}
```

## Proto Definition

```protobuf
service SystemConfigService {
    // API cho user - Lấy settings theo group (chỉ trả key-value)
    rpc GetUserSettingsByGroup(GetUserSettingsByGroupRequest) returns (GetUserSettingsByGroupResponse) {
        option (google.api.http) = {
            get: "/v2/hub/settings/{groupKey}"
        };
    }
}

message GetUserSettingsByGroupRequest {
    string groupKey = 1;
}

message GetUserSettingsByGroupResponse {
    map<string, string> settings = 1;  // key-value pairs
    string groupKey = 2;
}
```

## Implementation Details

### Usecase
```go
// GetUserSettingsByGroup lấy settings theo group dạng key-value map cho user
func (uc *SystemConfigUsecase) GetUserSettingsByGroup(ctx context.Context, groupKey string) (map[string]string, error) {
    // Parse groupKey string to enum
    configGroup, err := enums.GetSystemConfigGroup(groupKey)
    if err != nil {
        return nil, _errors.ReturnError(400, fmt.Sprintf("Group key không hợp lệ: %s", groupKey))
    }

    // Lấy configs từ DB
    configs, err := uc.repo.GetByGroup(ctx, configGroup)
    if err != nil {
        return nil, _errors.ReturnError(500, "Lỗi khi lấy settings")
    }

    // Convert sang map key-value
    settings := make(map[string]string)
    for _, config := range configs {
        settings[config.Key] = config.Value
    }

    return settings, nil
}
```

### Handler
```go
func (h *SystemConfigHandler) GetUserSettingsByGroup(ctx context.Context, req *hubpb.GetUserSettingsByGroupRequest) (*hubpb.GetUserSettingsByGroupResponse, error) {
    settings, err := h.usecase.GetUserSettingsByGroup(ctx, req.GroupKey)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get user settings: %v", err)
    }

    return &hubpb.GetUserSettingsByGroupResponse{
        Settings: settings,
        GroupKey: req.GroupKey,
    }, nil
}
```

## Benefits

1. **Simple Response**: Chỉ trả về data cần thiết (key-value)
2. **Easy to Use**: Dễ parse và sử dụng trong frontend
3. **Performance**: Giảm kích thước response
4. **Clean API**: Phân biệt rõ API cho admin vs user
5. **Type Safety**: Protobuf map type đảm bảo type safety

## Notes

- API này **không yêu cầu authentication** (public)
- Nếu cần authentication, thêm middleware check token
- Settings được cache ở usecase layer
- Response luôn trả về map, nếu group không có settings thì trả về empty map `{}`
- GroupKey không hợp lệ sẽ trả về error 400

