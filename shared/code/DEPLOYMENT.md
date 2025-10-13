# 🚀 Deployment Guide

Hướng dẫn deploy microservices lên production server.

## 📋 Prerequisites

- Đã setup SSH key với server production
- Đã cài đặt Docker và Docker Compose trên server
- Go version >= 1.21
- Đã config các file `config/production/*.yml`

## 🎯 Quick Start

### Deploy tất cả services

```bash
cd shared/code
./deploy-all.sh
```

### Deploy với các options

```bash
# Deploy song song (nhanh hơn)
./deploy-all.sh --parallel

# Bỏ qua build shared modules
./deploy-all.sh --skip-shared

# Deploy chỉ một số services cụ thể
./deploy-all.sh --services="auth,user,gateway"
```

### Deploy các nhóm services

```bash
# Deploy chỉ core services (auth, user, gateway)
./deploy-core.sh

# Deploy chỉ business services (exclude core)
./deploy-business.sh

# Deploy business services song song
./deploy-business.sh --parallel
```

### Deploy từng service riêng lẻ

```bash
# Deploy một service cụ thể
./release.sh auth
./release.sh user
./release.sh gateway
```

## 📊 Monitoring

### Check status của các services

```bash
./check-services.sh
```

### Xem logs của một service

```bash
ssh root@103.172.239.89
cd /root/golang-microservice/<service-name>
docker compose logs -f <service-name>
```

### Xem logs deployment

```bash
# Logs được lưu trong thư mục logs/
cat logs/<service-name>_deploy.log
```

## 🛠️ Maintenance Commands

### Pre-deployment Validation

**HIGHLY RECOMMENDED:** Chạy validation check trước khi deploy:

```bash
# Check tất cả services
./pre-deploy-check.sh

# Check một service cụ thể
./pre-deploy-check.sh notification
./pre-deploy-check.sh gateway
```

**Script này sẽ check:**
- ✅ Tất cả config production files (*.yml) tồn tại
- ✅ Firebase credentials cho notification service
- ✅ Các thư mục đặc biệt (configs/, docs/) cho services cần
- ✅ Dockerfile tồn tại cho mỗi service
- ✅ go.mod files tồn tại
- ✅ Shared modules đầy đủ

**Note:** `deploy-all.sh` tự động chạy validation trước khi deploy!

### Go Mod Tidy

Trước khi deploy, nên chạy `go mod tidy` để clean dependencies:

```bash
# Tidy tất cả services
./tidy-all.sh

# Tidy chỉ một số services
./tidy-all.sh --services="auth,user,gateway"

# Tidy bỏ qua shared modules
./tidy-all.sh --skip-shared
```

**Khi nào cần chạy:**
- Sau khi thêm/xóa dependencies
- Trước khi deploy lên production
- Khi có conflict về go.mod/go.sum
- Định kỳ để clean unused dependencies

## 🔧 Available Services

### Core Services (Deploy trước)
- `auth` - Authentication service
- `user` - User profile service  
- `user-grpc` - User gRPC service
- `gateway` - API Gateway (deploy sau cùng)

### Business Services
- `bdspro` - Real estate management service
- `bdspro-grpc` - Real estate gRPC service
- `crm` - Customer relationship management
- `notification` - Notification service
- `social` - Social networking features
- `chat` - Chat messaging service
- `chathttp` - Chat HTTP API
- `file` - File storage service
- `map` - Map and location service
- `payment` - Payment processing
- `membership` - Membership management
- `relay` - WebSocket relay service
- `appointment` - Appointment scheduling
- `org` - Organization management
- `transaction` - Transaction management

## 📝 Deployment Flow

### Sequential Deployment (Default)

```
1. Build shared modules
2. Deploy auth service
3. Deploy user service
4. Deploy user-grpc service
5. Deploy gateway service
6. Deploy business services (one by one)
```

**Thời gian**: ~30-45 phút cho tất cả services

### Parallel Deployment (Faster)

```
1. Build shared modules
2. Phase 1: Deploy core services (sequential)
   - auth → user → user-grpc → gateway
3. Wait 10 seconds
4. Phase 2: Deploy business services (parallel)
   - All business services deploy at same time
```

**Thời gian**: ~15-20 phút cho tất cả services

**⚠️ Lưu ý**: Parallel mode deploy nhanh hơn nhưng khó debug nếu có lỗi.

## 🐛 Troubleshooting

### Service deployment failed

```bash
# Check logs
cat logs/<service-name>_deploy.log

# SSH vào server và check
ssh root@103.172.239.89
cd /root/golang-microservice/<service-name>
docker compose logs <service-name>
```

### Service không start được

```bash
# Restart service
ssh root@103.172.239.89
cd /root/golang-microservice/<service-name>
docker compose restart <service-name>
```

### Port conflict

```bash
# Check port usage
ssh root@103.172.239.89
netstat -tulpn | grep <port-number>

# Stop conflicting service
docker compose stop <service-name>
```

### Memory issues

```bash
# Check memory usage
ssh root@103.172.239.89
free -h
docker stats

# Clean up unused images
docker system prune -a
```

## 🔄 Rollback

```bash
# SSH vào server
ssh root@103.172.239.89
cd /root/golang-microservice/<service-name>

# Stop service hiện tại
docker compose stop <service-name>

# Restore từ backup (nếu có)
# Hoặc re-deploy version cũ từ git tag/commit

# Start lại service
docker compose up -d <service-name>
```

## 📌 Best Practices

1. **Always deploy core services first**: auth, user, gateway
2. **Test locally before deploying**: Ensure service builds and runs locally
3. **Deploy in maintenance window**: Avoid peak hours
4. **Monitor after deployment**: Check logs and metrics
5. **Have rollback plan ready**: Know how to quickly rollback
6. **Use parallel mode carefully**: Only when you're confident
7. **Keep logs**: Check `logs/` directory for deployment history

## 🔐 Security

- Không commit SSH keys vào git
- Không commit production credentials
- Luôn dùng environment variables cho secrets
- Review code trước khi deploy
- Backup database trước khi deploy breaking changes

## 📞 Support

Nếu gặp vấn đề khi deploy, hãy:
1. Check logs trong thư mục `logs/`
2. SSH vào server và check docker logs
3. Review file config trong `config/production/`
4. Check server resources (RAM, CPU, Disk)

---

**Last updated**: October 2025
**Server**: 103.172.239.89
**Workspace**: /root/golang-microservice

