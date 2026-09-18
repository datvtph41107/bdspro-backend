# Developer Operating Guide

## 1. Clone → ready

```bash
git clone <repository>
cd bdspro-backend
make help
make setup
make up
make status
```

`setup` tạo local env còn thiếu và bootstrap pinned engineering toolchain.
`up` chạy PostgreSQL/PostGIS, Redis, RabbitMQ bằng Docker rồi build/migrate/start
toàn bộ Go service trực tiếp từ source. Không có application container trong
vòng lặp này.

Trong environment `development`, `make up` tạo hai identity idempotent trong
User PostgreSQL để developer vào thẳng UI mà không cần nhà cung cấp SMS:

```text
Admin:  admin / admin123
Client: 0900000000 (hoặc client) / client123
```

Chúng chỉ là identity; plan, payment, subscription, quota và report vẫn phải
được tạo qua API/state machine thật.

## 2. Nhận task

Ví dụ “Payment settlement retry”:

```text
TASK → Payment owner → settlement capability → execution flow → change → proof
```

Mở `payment-service/README.md`, entrypoint/config và capability liên quan; không đọc toàn repository.

## 3. Daily edit loop

```bash
make deps service=payment-service
make dev service=payment-service
```

`make dev` chạy native Air. Logs stdout/stderr nằm ngay terminal. Save source → Go build → process restart. Đây là loop mặc định vì tối ưu feedback/debug.

Runtime ownership có bốn trạng thái: `SUPERVISED`, `DEV`, `FOREIGN`, `DOWN`. `make status` hiển thị ownership và `ready=yes/no` riêng; port chỉ là readiness signal, không phải ownership truth. Root-routed và service-local `make dev` cùng dùng một DEV lifecycle owner, nên second DEV fail fast và SUPERVISED → DEV transfer không bị định nghĩa hai lần.

Structured application logs của cả `make up` và `make dev` dùng cùng một root tuyệt đối do repository sở hữu: `<repo>/.tmp/development/logs`. Logger tiếp tục tự tách theo service/run bên dưới root này; vị trí evidence không phụ thuộc service CWD. `make logs` vẫn là raw process-stream view, không phải structured-query command. Khi service ở DEV, Air/child stdout và stderr vẫn hiện trực tiếp trong terminal đồng thời được mirror vào cùng raw process-log view mà `make logs` đọc.

VS Code: `Ctrl+Shift+P` → `Tasks: Run Task` → `BDSPro · Dev · <Service>`.

## 4. Proof

```bash
make test service=payment-service
```

Nếu schema đổi:

```bash
make migrate service=payment-service
make verify-migrations
```

Nếu contract/shared/cross-service behavior đổi:

```bash
make verify
```

## 5. Whole-system local runtime

Chỉ khi cần chứng minh phối hợp/packaging:

```bash
make up
make status
make smoke
make test-e2e
```

Source vừa đổi và muốn refresh một process native ổn định:

```bash
make native-restart service=payment-service
```

Hoặc dùng `make dev service=payment-service` để xem log và hot reload ngay trong
terminal. `make logs service=payment-service` đọc cùng raw stdout/stderr stream
dù service đang ở SUPERVISED hay DEV.

Full Compose chỉ còn là proof riêng cho image/package:

```bash
make integration-up
make integration-status
make integration-down
```

## 6. Final acceptance

```bash
make accept
```

Đây là gate đắt nhất: full source proof + image rebuild/runtime + E2E.

## 7. Stop condition

Không refactor global structure chỉ để đẹp. Chỉ thay architecture khi task thực chứng minh ownership ambiguity, duplicate workflow, hidden dependency, slow feedback hoặc correctness risk.
