# Script Export User Service Data

Script này dùng để export toàn bộ dữ liệu từ database của user-service ra file JSON.

## Dữ liệu được export

### Core Entities
- **User Profiles**: Thông tin người dùng (user_profile table)
- **Admin Profiles**: Thông tin admin (admin_profile table)
- **Plans**: Các gói membership (plans table)
- **Memberships**: Thông tin đăng ký gói của user (memberships table)

### Social Entities
- **Groups**: Nhóm bạn bè (db_group table)
- **Friends**: Quan hệ bạn bè (db_friend table)
- **Follows**: Quan hệ theo dõi (db_follow table)
- **Blocks**: Danh sách chặn (db_block table)
- **Contacts**: Danh bạ liên hệ (db_contact table)
- **Bookmark Users**: Người dùng được đánh dấu (db_bookmark_user table)

### Profile Related
- **Professions**: Nghề nghiệp (db_profession table)
- **Certifications**: Chứng chỉ (db_certification table)

### System
- **System Configs**: Cấu hình hệ thống (system_config table)

## Cách sử dụng

### 1. Sử dụng script shell (Khuyến nghị)

```bash
cd user-service/scripts
./run_export.sh
```

### 2. Chạy trực tiếp với Go

```bash
cd user-service/scripts
export DATABASE_URL="host=14.225.210.29 user=hieubap password=Matkhau.1 dbname=db_user port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
go run export_data.go
```

### 3. Build và chạy binary

```bash
cd user-service/scripts

# Build
go build -o export_tool export_data.go

# Chạy
./export_tool
```

## Cấu hình Database

Script sẽ đọc connection string từ biến môi trường `DATABASE_URL`. Nếu không có, sẽ sử dụng giá trị từ `config/app.yml`:

```
host=14.225.210.29 user=hieubap password=Matkhau.1 dbname=db_user port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh
```

Để thay đổi, export biến môi trường:

```bash
export DATABASE_URL="host=your_host user=your_user password=your_password dbname=your_db port=5432 sslmode=disable"
```

## Output

File JSON sẽ được tạo với tên `user_service_export.json` trong thư mục hiện tại.

### Cấu trúc file JSON:

```json
{
  "userProfiles": [...],
  "adminProfiles": [...],
  "plans": [...],
  "memberships": [...],
  "groups": [...],
  "friends": [...],
  "follows": [...],
  "blocks": [...],
  "contacts": [...],
  "bookmarkUsers": [...],
  "professions": [...],
  "certifications": [...],
  "systemConfigs": [...],
  "exportedAt": "2025-10-13T10:30:00+07:00",
  "totalUsers": 1000,
  "totalPlans": 5,
  "totalMemberships": 500,
  "totalFriends": 2000
}
```

## Xem kết quả

Để xem file JSON với format đẹp, sử dụng `jq`:

```bash
cat user_service_export.json | jq
```

Hoặc mở bằng editor:

```bash
code user_service_export.json
# hoặc
open user_service_export.json
```

## Lưu ý

- Script chỉ export các records chưa bị xóa (`deleted_at IS NULL`)
- Tất cả dữ liệu được export trong một file JSON duy nhất
- Script tự động handle lỗi và báo cáo chi tiết
- Nếu một bảng không tồn tại hoặc lỗi, script sẽ warning nhưng vẫn tiếp tục với các bảng khác

## Troubleshooting

### Lỗi kết nối database

```
Failed to connect to database: ...
```

**Giải pháp**: Kiểm tra lại DATABASE_URL và đảm bảo database đang chạy, firewall cho phép kết nối.

### Lỗi không tìm thấy bảng

```
Warning: Failed to fetch user profiles: table "user_profile" does not exist
```

**Giải pháp**: Đảm bảo đã chạy migration và bảng đã được tạo.

### File không được tạo

```
Failed to write file: permission denied
```

**Giải pháp**: Kiểm tra quyền ghi file trong thư mục hiện tại.

### Kết nối bị từ chối

```
dial tcp: connect: connection refused
```

**Giải pháp**: 
- Kiểm tra database server có đang chạy không
- Kiểm tra firewall/security group có cho phép kết nối từ IP của bạn không
- Kiểm tra host và port có đúng không

