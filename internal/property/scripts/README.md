# BDSPro Service Scripts

## Import Location V2 Data

Script để import dữ liệu tỉnh/thành và phường/xã từ file `locationv2.json` vào database.

### Cách sử dụng

#### 1. Sử dụng script shell (khuyến nghị):

```bash
cd bdspro-service/scripts
./import-location-v2.sh
```

#### 2. Chạy trực tiếp với Go:

```bash
cd bdspro-service/scripts
go run import_location_v2_data.go
```

### Cấu hình Database

Mặc định script sẽ kết nối tới:
```
host=localhost user=postgres password=postgres dbname=bdspro_service port=5432 sslmode=disable
```

Để sử dụng connection string khác, set biến môi trường `DATABASE_URL`:

```bash
export DATABASE_URL='host=localhost user=postgres password=postgres dbname=bdspro_service port=5432 sslmode=disable'
./import-location-v2.sh
```

### Lưu ý

- Script sẽ tự động truncate dữ liệu cũ trong bảng `province_v2` và `ward_v2`
- Dữ liệu được import từ file `../../shared/code/locationv2.json`
- Migration `010_create_location_v2_tables.sql` cần được chạy trước

