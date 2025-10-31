# Hub Service Scripts

## Import Location Data

Scripts để import dữ liệu địa điểm vào database.

### 1. Import Location V1 (provinces, districts, wards)

```bash
cd hub-service/scripts
./import-location.sh
```

### 2. Import Location V2 (province_v2, ward_v2)

```bash
cd hub-service/scripts
./import-location-v2.sh
```

## Import Location V2 Data

Script để import dữ liệu tỉnh/thành và phường/xã từ file `locationv2.json` vào database.

### Cách sử dụng

#### 1. Sử dụng script shell (khuyến nghị):

```bash
cd hub-service/scripts
./import-location-v2.sh
```

#### 2. Chạy trực tiếp với Go:

```bash
cd hub-service/scripts
go run import_location_v2_data.go
```

### Cấu hình Database

Mặc định script sẽ kết nối tới:
```
host=localhost user=postgres password=postgres dbname=hub_service port=5432 sslmode=disable
```

Để sử dụng connection string khác, set biến môi trường `DATABASE_URL`:

```bash
export DATABASE_URL='host=localhost user=postgres password=postgres dbname=hub_service port=5432 sslmode=disable'
./import-location-v2.sh
```

### Lưu ý

- Script sẽ tự động truncate dữ liệu cũ trong bảng `province_v2` và `ward_v2`
- Dữ liệu được import từ file `../../shared/code/locationv2.json`
- Migration `003_create_location_v2_tables.sql` cần được chạy trước
