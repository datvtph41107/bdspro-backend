# Generate Proto Code for User Settings API

## ✅ Đã hoàn thành

1. ✅ Proto definitions (`shared/protobuf/schema/hub/system_config.proto`)
2. ✅ Usecase method (`hub-service/internal/usecase/system_config_usecase.go`)
3. ✅ Handler method (`hub-service/infra/handler/system_config_handler.go`)
4. ✅ Documentation (`hub-service/USER_SETTINGS_API.md`)

## 🔧 Cần generate proto code

### 1. Generate Go code từ proto

```bash
cd shared/protobuf/schema/hub

# Generate Go protobuf + gRPC code
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  system_config.proto
```

### 2. Generate gRPC Gateway (nếu cần REST API)

```bash
cd shared/protobuf/schema/hub

protoc --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  system_config.proto
```

### 3. Generate Swagger documentation (optional)

```bash
cd shared/protobuf/schema/hub

protoc --openapiv2_out=../../docs/hub \
  --openapiv2_opt=logtostderr=true \
  system_config.proto
```

### 4. Compile và test

```bash
cd hub-service

# Build service
go build -o hub-service cmd/grpc/main.go

# Run service
./hub-service
```

## 📋 API Endpoint mới

```
GET /v2/hub/settings/{groupKey}
```

### Request Example
```bash
curl http://localhost:8080/v2/hub/settings/post
```

### Response Example
```json
{
  "settings": {
    "display_name": "true",
    "auto_save": "false",
    "max_upload": "10"
  },
  "groupKey": "post"
}
```

## 🎯 Use Cases

### Frontend Usage
```javascript
// Lấy settings cho post
const response = await fetch('/v2/hub/settings/post');
const { settings } = await response.json();

// Sử dụng settings
const maxUpload = parseInt(settings.max_upload || '5');
const allowComments = settings.allow_comments === 'true';
```

### Mobile App
```kotlin
// Android/Kotlin
suspend fun getPostSettings(): Map<String, String> {
    val response = client.getUserSettingsByGroup(
        GetUserSettingsByGroupRequest.newBuilder()
            .setGroupKey("post")
            .build()
    )
    return response.settingsMap
}
```

## 📁 Files Changed

1. **shared/protobuf/schema/hub/system_config.proto**
   - Added `rpc GetUserSettingsByGroup`
   - Added `GetUserSettingsByGroupRequest`
   - Added `GetUserSettingsByGroupResponse`

2. **hub-service/internal/usecase/system_config_usecase.go**
   - Added `GetUserSettingsByGroup()` method

3. **hub-service/infra/handler/system_config_handler.go**
   - Added `GetUserSettingsByGroup()` handler

## 🔍 Testing

### Test với grpcurl
```bash
grpcurl -plaintext \
  -d '{"groupKey": "post"}' \
  localhost:50051 \
  hubpb.SystemConfigService/GetUserSettingsByGroup
```

### Test với REST
```bash
curl http://localhost:8080/v2/hub/settings/post
```

## ⚠️ Notes

- Proto code **PHẢI được generate** trước khi build service
- Nếu lỗi `undefined: hubpb.GetUserSettingsByGroupRequest` → chạy generate proto
- API endpoint: `/v2/hub/settings/{groupKey}` (khác với admin API `/v2/hub/system-config/...`)
- Response chỉ có key-value, không có metadata

