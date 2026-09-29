# BDSPro Context - Tài Liệu Dự Án

**Last Updated**: October 18, 2025

---

## 📚 Tổng Quan

Thư mục này chứa **toàn bộ kiến thức** về project BDSPro Microservices - được tổ chức để **tra cứu nhanh** và **hiểu sâu** hệ thống.

---

## ⚡ BẮT ĐẦU NHANH

### 1️⃣ Lần Đầu Dùng Project?

```
Đọc ngay: 23_QUICK_REFERENCE_GUIDE.md
          ↓
Mất 5 phút scan
          ↓  
Biết ngay cách tạo API, patterns, conventions!
```

### 2️⃣ Đang Code, Cần Tra Cứu?

```bash
# Cmd+P (hoặc Ctrl+P) → Gõ:
23_QUICK_REFERENCE_GUIDE.md

# Sau đó Cmd+F để tìm trong file:
- "Repository pattern"
- "Import aliases"  
- "profileId"
- "Payment service"
- v.v...
```

### 3️⃣ Muốn Hiểu Sâu 1 Service?

```
Quick Reference → Tìm service → Link đến doc chi tiết
```

---

## 📂 Cấu Trúc Files

### ⭐ Files Quan Trọng Nhất

| File | Mục Đích | Khi Nào Dùng |
|------|----------|--------------|
| **23_QUICK_REFERENCE_GUIDE.md** | Tra cứu cực nhanh | ✅ **Mở ngay khi code!** |
| **00_INDEX.md** | Navigation, roadmap | Tìm file nào đọc |
| **00_SYSTEM_OVERVIEW.md** | Hiểu tổng quan | Mới vào project |
| **07_DEVELOPMENT_WORKFLOW.md** | Tạo API step-by-step | Lần đầu tạo API |

### 📖 Core Services (Chi Tiết)

| File | Service | Nội Dung |
|------|---------|----------|
| **02_AUTH_SERVICE.md** | Auth (50062) | JWT, OAuth, RBAC |
| **05_ORGANIZATION_SERVICE.md** | Organization (50052) | Org hierarchy, deals |
| **06_BDSPRO_SERVICE.md** | BDSPro (50053) | Property, product, asset |

### 🔧 Foundation

| File | Nội Dung |
|------|----------|
| **03_SHARED_COMMON_LIBRARY.md** | BaseEntity, CRUD, DTOs |
| **04_TECHNOLOGY_STACK.md** | Go, gRPC, PostgreSQL, Redis |
| **09_DATABASE_ARCHITECTURE.md** | DB schemas, tables |

### 🌟 Supporting Services

| File | Services Covered |
|------|------------------|
| **08_OTHER_SERVICES_SUMMARY.md** | CRM, Social, Chat, File, etc. |
| **22_ADDITIONAL_SERVICES_DEEP_DIVE.md** | TQD, Assistant, Payment, Membership |

### 📊 Advanced Topics

| File | Nội Dung |
|------|----------|
| **13_AUTHENTICATION_FLOW.md** | Auth flow chi tiết |
| **14_INTER_SERVICE_COMMUNICATION.md** | gRPC patterns |
| **15_DATA_ACCESS_PATTERNS.md** | Repository, transaction |
| **16_BUSINESS_FLOWS.md** | Business logic flows |

---

## 🎯 Use Cases

### Scenario 1: "Tôi cần tạo API mới"

```bash
1. Mở: 23_QUICK_REFERENCE_GUIDE.md
2. Tìm: "TẠO API MỚI - 7 BƯỚC"
3. Follow từng bước
4. Copy code templates từ file
5. Done! ✅
```

### Scenario 2: "Repository pattern viết sao?"

```bash
1. Mở: 23_QUICK_REFERENCE_GUIDE.md
2. Tìm: "REPOSITORY PATTERN"
3. Copy template
4. Thay tên entity
5. Done! ✅
```

### Scenario 3: "Payment service làm gì?"

```bash
1. Mở: 23_QUICK_REFERENCE_GUIDE.md
2. Tìm: "Payment Service"
3. Xem quick info
4. Muốn hiểu sâu? → Click link 22_ADDITIONAL_SERVICES_DEEP_DIVE.md
```

### Scenario 4: "Service X có port gì?"

```bash
1. Mở: 23_QUICK_REFERENCE_GUIDE.md
2. Tìm: "SERVICES PORTS"
3. Xem bảng
4. Done! ✅
```

### Scenario 5: "Import aliases nào đúng?"

```bash
1. Mở: 23_QUICK_REFERENCE_GUIDE.md  
2. Tìm: "IMPORT ALIASES"
3. Copy paste
4. Done! ✅
```

---

## 🚀 Workflow Thực Tế

### Developer Mới

```
Day 1:
├─ Scan 23_QUICK_REFERENCE_GUIDE.md (5 min)
├─ Read 00_SYSTEM_OVERVIEW.md (10 min)
└─ Practice tạo API đầu tiên (30 min)

Day 2-4:
├─ Đọc chi tiết services cần dùng
└─ Keep Quick Reference open để lookup!
```

### Developer Có Kinh Nghiệm

```
Start:
├─ Scan 23_QUICK_REFERENCE_GUIDE.md (5 min)
└─ Start coding ngay!

Khi Cần:
└─ Quick Reference → Search → Copy → Paste → Done!
```

---

## 📊 Thống Kê

```yaml
Total Files: 23
Total Lines: ~7,500+
Services Documented: 19/19 (100%)
Code Templates: 10+
Quick Reference Lookup: < 30 seconds
Full Read Time: 5-6 hours
```

---

## 💡 Tips & Tricks

### 1. Bookmark Quick Reference

```
# VS Code / Cursor
Cmd+P → 23_QUICK → Enter
Cmd+B để bookmark sidebar
```

### 2. Search Hiệu Quả

```
# Trong Quick Reference, tìm:
- "profileId" → Cách lấy user info
- "CrudRepo" → Repository pattern
- "50055" → Port của service nào?
- "PAYMENT" → Transaction types
```

### 3. Copy-Paste Templates

```
Quick Reference có sẵn templates:
├─ Handler template
├─ Usecase template
├─ Repository template
└─ Copy → Rename → Done!
```

### 4. Debug Nhanh

```
API không chạy?
→ Quick Reference → "DEBUG CHECKLIST"
→ Check từng bước 1-10
```

---

## 🔄 Cập Nhật & Bảo Trì

### Khi Nào Update Context?

- ✅ Thêm service mới
- ✅ Thay đổi pattern quan trọng
- ✅ Thêm best practices mới
- ✅ Phát hiện antipatterns cần tránh

### Cách Update

```bash
1. Edit file .md tương ứng
2. Update 00_INDEX.md nếu thêm file mới
3. Update 23_QUICK_REFERENCE_GUIDE.md nếu thay đổi conventions
4. Commit với message rõ ràng
```

---

## 🎉 Kết Luận

Context này được thiết kế để:

✅ **Tra cứu cực nhanh** (< 30 giây)  
✅ **Hiểu sâu** khi cần (đọc chi tiết)  
✅ **Copy-paste** templates (code nhanh)  
✅ **Tự học** (learning paths rõ ràng)  
✅ **Onboard nhanh** (new devs < 1 day)

---

## 📞 Hỗ Trợ

**Không tìm thấy info cần thiết?**

1. Search trong `23_QUICK_REFERENCE_GUIDE.md`
2. Search trong `00_INDEX.md`
3. Check các file chi tiết services
4. Vẫn không có → Cần bổ sung vào context!

---

**Happy Coding! 🚀**

*Remember: Quick Reference là best friend của bạn!*
