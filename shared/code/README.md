# 🚀 Deployment Scripts

Scripts để deploy microservices lên production server.

## 📖 Quick Reference

### 🔨 Development Scripts

```bash
# Generate tất cả protobuf files (buf/protoc)
./buf-all.sh

# Generate protobuf cho specific modules
./buf-all.sh --modules="auth,user,bdspro"

# Generate protobuf với clean trước
./buf-all.sh --clean

# Wire dependency injection cho tất cả services
./wire-all.sh

# Wire chỉ một số services
./wire-all.sh --services="auth,user,gateway"

# Go mod tidy tất cả services
./tidy-all.sh

# Go mod tidy chỉ một số services
./tidy-all.sh --services="auth,user,gateway"

# Go mod tidy bỏ qua shared modules
./tidy-all.sh --skip-shared
```

### 🚀 Deployment Scripts

```bash
# Pre-deployment validation (recommended)
./pre-deploy-check.sh

# Deploy tất cả services (sequential - an toàn)
# Note: Tự động chạy pre-deploy-check trước khi deploy
./deploy-all.sh

# Deploy tất cả services (parallel - nhanh)
./deploy-all.sh --parallel

# Deploy chỉ core services
./deploy-core.sh

# Deploy chỉ business services
./deploy-business.sh --parallel

# Deploy services cụ thể
./deploy-all.sh --services="auth,user,gateway"

# Deploy một service đơn lẻ
./release.sh <service-name>

# Check status các services trên server
./check-services.sh

# Restart services
./restart-service.sh <service-name> [service-name2 ...]

# Stop tất cả services
./stop-all.sh
```

## 📚 Documentation

Xem file [DEPLOYMENT.md](./DEPLOYMENT.md) để biết thêm chi tiết về:
- Deployment flow
- Troubleshooting
- Best practices
- Rollback procedures

## 🔄 Development Workflow

Quy trình phát triển chuẩn khi thay đổi proto hoặc code:

```bash
# 1. Chỉnh sửa proto files trong shared/protobuf/schema/
# 2. Generate protobuf
./buf-all.sh --clean

# 3. Implement code trong services (handler, usecase, repo...)
# 4. Generate wire dependency injection
./wire-all.sh

# 5. Tidy dependencies
./tidy-all.sh

# 6. Test local...
# 7. Deploy to production
./deploy-all.sh
```

**Quick workflow cho một service:**
```bash
# Generate proto cho module cụ thể
./buf-all.sh --modules="bdspro" --clean

# Generate wire cho service cụ thể  
./wire-all.sh --services="bdspro"

# Tidy cho service cụ thể
./tidy-all.sh --services="bdspro,gateway"

# Deploy service cụ thể
./release.sh bdspro
```

## 🎯 Available Services

### Core Services
- `auth` - Authentication
- `user` - User profile  
- `user-grpc` - User gRPC
- `gateway` - API Gateway

### Business Services
- `bdspro` - Real estate
- `crm` - CRM
- `notification` - Notifications
- `social` - Social features
- `chat`, `chathttp` - Messaging
- `file` - File storage
- `map` - Maps
- `tqd` - TQD (Directory/POI service)
- `payment` - Payments
- `membership` - Memberships
- `relay` - WebSocket relay
- `appointment` - Appointments
- `org` - Organizations
- `transaction` - Transactions
- `marketing` - Marketing campaigns
- `feedback` - Feedback & reports

## 📁 Available Scripts

### Development Scripts
| Script | Description | Usage |
|--------|-------------|-------|
| `buf-all.sh` | Generate protobuf cho tất cả modules | `./buf-all.sh [--options]` |
| `wire-all.sh` | Wire DI cho tất cả services | `./wire-all.sh [--options]` |
| `tidy-all.sh` | Go mod tidy tất cả services | `./tidy-all.sh [--options]` |

### Deployment Scripts
| Script | Description | Usage |
|--------|-------------|-------|
| `pre-deploy-check.sh` | Validate trước khi deploy | `./pre-deploy-check.sh [service]` |
| `deploy-all.sh` | Deploy tất cả services | `./deploy-all.sh [--options]` |
| `deploy-core.sh` | Deploy core services | `./deploy-core.sh` |
| `deploy-business.sh` | Deploy business services | `./deploy-business.sh [--parallel]` |
| `release.sh` | Deploy một service | `./release.sh <service-name>` |
| `check-services.sh` | Check status trên server | `./check-services.sh` |
| `restart-service.sh` | Restart services | `./restart-service.sh <services...>` |
| `stop-all.sh` | Stop tất cả services | `./stop-all.sh [--remove]` |

## 🔧 Server Info

- **Host**: 103.172.239.89
- **User**: root
- **Path**: /root/golang-microservice
- **ENV**: production (ENV_RUNTIME=release)

## 📝 Notes

- Sequential deployment: ~30-45 phút
- Parallel deployment: ~15-20 phút
- Logs được lưu trong `logs/` directory
- Core services phải deploy trước business services
- `deploy-all.sh` tự động chạy `pre-deploy-check.sh`

