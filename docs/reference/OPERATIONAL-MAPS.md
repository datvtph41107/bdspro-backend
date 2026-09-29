# BDSPro Local Operational Maps

Tài liệu này mô tả **developer/local operation**: command flow, local runtime topology, owner process và đường developer lần theo source.

## Daily control flow

```text
make setup        one-time/fresh source
    |
make up           Docker infra + native Go processes
    |
make status       health/process view
    |
make dev service=<owner>-service
    |
make test service=<owner>-service
    |
make smoke        public boundary
    |
make verify       repository proof when needed
    |
make accept       full local candidate proof
```

## Runtime topology

```text
PostgreSQL/PostGIS (one physical local server)
  |- user_service
  |- organization
  |- payment_service
  |- qhpro_tqd
  |- qhpro_notification
  |- file_service
  |- hub_service
  |- db_bdspro        source-health migration DB
  `- db_crm           source-health migration DB

Redis
  |- DB 0 User
  |- DB 1 TQD
  |- DB 2 Notification
  `- DB 3 Hub

RabbitMQ
  `- Payment -> consumer owners
```

Standing canonical applications:

```text
gateway
user
auth
organization
payment
tqd
notification
file
hub
assistant
```

Các application trên là process host được build từ working tree. Runtime state
và log có owner tại `.tmp/development/{pids,dev-pids,logs}`. `pids` nhận diện
SUPERVISED, `dev-pids` nhận diện foreground DEV; port/readiness được kiểm tra riêng
và port mở không có owner hợp lệ được báo FOREIGN. Compose application services
chỉ được gọi bằng `make integration-up` để kiểm tra image/package, không thuộc
default development runtime.

`*-migrate` và acceptance fixtures là one-shot operations, không phải business
services.

## Owned actors

```text
Payment process
  |- gRPC/API
  |- fulfillment actor
  `- outbox publisher

Notification process
  |- gRPC/API
  |- consumer of Payment events (reconnecting RabbitMQ supervisor)
  `- delivery actor when configured

TQD process
  |- gRPC/API
  `- report actor
```

Compatibility binaries có thể còn trong source để recovery/caller proof, nhưng
không tạo standing process riêng trong canonical development topology.

## Developer source path

```text
root Makefile
 -> service Makefile
 -> main/composition root
 -> config
 -> handler/transport
 -> business/usecase
 -> DB/client
 -> owned actor
 -> tests
```
