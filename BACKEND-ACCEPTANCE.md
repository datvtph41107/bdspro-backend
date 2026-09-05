# BDSPro Backend Acceptance

Ngày kiểm chứng: 2026-09-03 (Asia/Bangkok)  
Source revision nền: `c9a06c8cad26b0e65cb3d64d8d596687c41a6192`  
Trạng thái: **CORE COMMERCIAL ACCEPTED; WHOLE BACKEND NOT CLOSED**

Working tree đang có thay đổi lớn chưa commit. Revision trên chỉ nhận diện nền,
không phải checksum của deliverable cuối.

## 1. Environment

| Thành phần | Phiên bản đã dùng |
| --- | --- |
| Go | `go1.25.10 linux/amd64` từ pinned repository toolchain |
| Docker client/server | `26.1.1 / 26.1.1` |
| Docker Compose | `2.27.0-desktop.2` |
| Compose project | `bdspro` |
| Image convention | `bdspro/<service>-service:local` |
| Local infrastructure | 1 PostgreSQL/PostGIS + 1 Redis + 1 RabbitMQ; logical DB/Redis ownership tách riêng |

## 2. Bằng chứng đã PASS

| Gate | Kết quả |
| --- | --- |
| Root/service env ownership | PASS |
| Runtime config isolation | PASS |
| Architecture/business owner guards | PASS |
| Core `go test ./...` | PASS: Auth, User, Organization, Payment, TQD, Notification, File, Gateway, Hub |
| Core `go vet ./...` | PASS |
| Core `go test -race ./...` | PASS |
| Core binary build | PASS |
| Remaining repository modules | PASS vet/test/race/build: Assistant, BDSPro, Chat, Chat-v1, CRM, Map, Relay, Search, Social, Shared |
| Compose render/start/health | PASS; application và infrastructure container đang healthy |
| Public smoke | PASS: Gateway live/ready, File live, commercial catalog |
| Customer business E2E | PASS |
| Admin projection E2E | PASS |
| Fresh workspace onboarding | PASS: copied working truth without `.git/.env/.tmp/bin`; `make setup` 10m42s |
| Fresh database/bootstrap | PASS: new disposable volumes, all versioned migrations, health, smoke và full E2E |
| Recovery labs | PASS: TQD stale claim, File response loss, Payment Outbox/RabbitMQ redelivery, Notification Inbox dedupe, Redis loss/rebuild |
| Shared local infrastructure cutover | PASS: old volumes retained; topology mới start/health/smoke/full E2E lại |
| `git diff --check` | PASS |

Customer E2E không ghi business state trực tiếp. Fixture chỉ tạo credential, sau
đó test gọi public HTTP qua Gateway:

```text
login -> paid catalog plan -> checkout -> payment attempt -> Sepay webhook
-> ACTIVE order -> subscription -> entitlement/quota -> report acceptance
-> idempotent replay -> report actor -> File PDF -> durable usage +1
-> Payment Outbox -> RabbitMQ -> Notification Inbox -> customer notification
```

Admin E2E đã đọc thành công catalog, subscription, usage, generated report,
payment order và fulfillment projection qua JWT + IAM.

Hai lỗi acceptance đã được phát hiện và sửa theo contract thật:

- PDF proof chấp nhận revision hợp lệ `%PDF-*`, không khóa renderer vào 1.4.
- Notification proof dùng page zero-based giống protobuf/handler và React Native.

Các lỗi workflow/runtime được phát hiện trong vòng fresh proof và đã sửa:

- RabbitMQ healthcheck chạy bằng user `rabbitmq`, không tạo root-owned Erlang cookie.
- Migration container retry khi PostgreSQL chuyển từ init server sang final server.
- Service Makefile source `.env` bằng absolute path; không phụ thuộc `.` có trong `PATH`.
- Pinned `migrate` được build với PostgreSQL driver và toolchain verify kiểm tra driver.
- Existing database có migration repair idempotent cho quyền Catalog/Usage Admin.

## 3. Runtime classification

| Nhóm | Trạng thái |
| --- | --- |
| Gateway/User/Auth/Organization/Payment/TQD/Notification/File/Hub/Assistant | Runtime-proven trong Compose core |
| BDSPro/CRM/Chat/Chat-v1/Map/Relay/Search/Social | Source-healthy; chưa có full runtime/E2E proof trong canonical Compose |
| AI | Chưa phân loại deployable/asset/dead source |
| Appointment | Owner replacement tại CRM đã có code evidence; chưa có production data cutover proof |
| Membership/Marketing/Task/Transaction production deployables | Chưa có retirement/replacement proof đầy đủ |

## 4. Gate còn thiếu trước `WHOLE BACKEND CLOSED`

1. Lập denominator chi tiết endpoint/table/event/deployment cho các owner chưa
   resolved trong [`documents/PRODUCTION-COMPATIBILITY.md`](documents/PRODUCTION-COMPATIBILITY.md).
2. Chứng minh hoặc phục hồi owner của Membership, Marketing, Task và legacy
   Transaction API; không suy diễn từ việc source compile.
3. Phân loại AI và quyết định rõ Chat/Chat-v1.
4. Audit production Compose/config và thay hard-coded host/secret behavior trong
   `shared/code/deploy.sh` bằng operator input, đồng thời giữ semantics deploy
   một service. Việc tăng quyền cho script deployment cần approval riêng.
5. Quyết định có retire vật lý các Notification worker compatibility binary hay giữ như
   recovery/scaling entrypoint; mặc định Compose hiện không chạy chúng.
6. Runtime-proof BDSPro/CRM/Chat/Chat-v1/Map/Relay/Search/Social trong topology
   canonical và bổ sung health/business smoke tương ứng.
7. Tạo immutable source manifest/checksum chỉ sau khi working tree được chốt.
8. React Admin commercial/quota đã được sync sớm với accepted core; React Native
   chỉ chốt sau backend whole-product proof.

## 5. Lệnh developer đã được kiểm chứng

```bash
make setup
make up
make status
make smoke
make test service=payment-service
make test
make test-e2e
make logs service=payment-service
make down
```

Không dùng trạng thái PASS của tài liệu này để deploy production. Cho đến khi
mục 4 được đóng, kết luận đúng là commercial/QHPro core đã vận hành end-to-end,
không phải toàn bộ BDSPro backend đã đạt production parity.
