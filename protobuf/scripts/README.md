# Swagger Merge Scripts

Scripts để merge tất cả file Swagger JSON trong thư mục `docs` và thêm Bearer authentication cho tất cả API.

## 📁 Files

- `merge_swagger.go` - Script Go để merge swagger files
- `merge_swagger.py` - Script Python để merge swagger files  
- `merge_swagger.sh` - Shell script để chạy Go script
- `README.md` - Hướng dẫn sử dụng

## 🚀 Cách sử dụng

### Option 1: Merge và Serve (Recommended)

```bash
# Merge swagger và hướng dẫn truy cập
cd protobuf
chmod +x scripts/merge_and_serve.sh
./scripts/merge_and_serve.sh
```

### Option 2: Sử dụng Go script

```bash
# Chạy trực tiếp
cd protobuf
go run scripts/merge_swagger.go

# Hoặc sử dụng shell script
chmod +x scripts/merge_swagger.sh
./scripts/merge_swagger.sh
```

### Option 3: Sử dụng Python script

```bash
cd protobuf
python3 scripts/merge_swagger.py
```

## 📋 Yêu cầu

### Cho Go script:
- Go 1.16+ installed

### Cho Python script:
- Python 3.6+ installed

## 📊 Output

Script sẽ tạo file `merged_swagger.json` trong thư mục `protobuf` với:

- ✅ Tất cả API endpoints từ tất cả services
- ✅ Bearer authentication cho mọi API
- ✅ Tất cả definitions/schemas
- ✅ Tất cả tags (không trùng lặp)
- ✅ Base path `/v2`
- ✅ Security definitions cho Bearer token

## 🔧 Tính năng

### Merge Logic:
1. **Paths**: Merge tất cả API endpoints từ các file swagger
2. **Definitions**: Merge tất cả schemas/definitions
3. **Tags**: Merge tags, loại bỏ duplicates
4. **Security**: Thêm Bearer authentication cho tất cả operations

### Bearer Authentication:
- Thêm `securityDefinitions.BearerAuth` với type `apiKey`
- Thêm `Authorization` header parameter cho mọi operation
- Set `security: [{"BearerAuth": []}]` cho mọi operation

## 📁 Cấu trúc thư mục

```
protobuf/
├── docs/
│   ├── organization/
│   │   ├── organization_branch.swagger.json
│   │   ├── organization_member.swagger.json
│   │   └── ...
│   ├── user/
│   │   └── ...
│   ├── bdspro/
│   │   └── ...
│   └── ...
├── scripts/
│   ├── merge_swagger.go
│   ├── merge_swagger.py
│   ├── merge_swagger.sh
│   └── README.md
└── merged_swagger.json (output)
```

## 🎯 Sử dụng với Gateway Service

Sau khi chạy merge script, file `merged_swagger.json` sẽ được copy vào `gateway-service/docs/`.

### Truy cập Swagger UI:

1. **Khởi động Gateway Service**:
   ```bash
   cd gateway-service
   go run cmd/http/main.go
   ```

2. **Truy cập Merged Swagger UI**:
   ```
   🌐 http://localhost:8080/swagger/merged/index.html
   ```

3. **Truy cập từng Service riêng lẻ**:
   ```
   🌐 http://localhost:8080/swagger/organization/index.html
   🌐 http://localhost:8080/swagger/user/index.html
   🌐 http://localhost:8080/swagger/bdspro/index.html
   🌐 http://localhost:8080/swagger/crm/index.html
   🌐 ... (các service khác)
   ```

4. **Truy cập trực tiếp file JSON**:
   ```
   📄 http://localhost:8080/swagger-doc/merged_swagger.json
   ```

### Sử dụng với các tool khác:

1. **Swagger UI Online**: Upload file lên https://editor.swagger.io/
2. **Swagger UI Local**: Sử dụng với Swagger UI Docker container
3. **API Documentation Tools**: Import vào các tool như Postman, Insomnia

## 🔍 Troubleshooting

### Lỗi thường gặp:

1. **"docs directory not found"**
   - Đảm bảo bạn đang ở thư mục `protobuf`
   - Kiểm tra thư mục `docs` có tồn tại

2. **"Go is not installed"**
   - Cài đặt Go từ https://golang.org/
   - Hoặc sử dụng Python script thay thế

3. **"Python is not installed"**
   - Cài đặt Python 3.6+ từ https://python.org/
   - Hoặc sử dụng Go script thay thế

4. **JSON parsing errors**
   - Kiểm tra các file swagger JSON có valid không
   - Script sẽ skip các file lỗi và tiếp tục với file khác

## 📝 Customization

### Thay đổi output file:
```go
// Trong merge_swagger.go
outputFile := "your_custom_name.json"
```

### Thay đổi base path:
```go
// Trong merge_swagger.go
BasePath: "/your-custom-path",
```

### Thay đổi security scheme:
```go
// Trong merge_swagger.go
mergedSwagger.SecurityDefinitions["BearerAuth"] = Security{
    Type:        "http",
    Scheme:      "bearer",
    BearerFormat: "JWT",
}
```

## 🤝 Contributing

Nếu bạn muốn cải thiện scripts:

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch
5. Tạo Pull Request

## 📄 License

MIT License - xem file LICENSE để biết thêm chi tiết. 