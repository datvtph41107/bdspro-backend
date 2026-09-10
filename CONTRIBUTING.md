# Contributing to BDSPro Backend

## 1. Decision loop trước khi code

Mọi change bắt đầu bằng 6 câu:

1. **Owner nào?** Service/business, engineering, contract hay root workspace?
2. **Capability nào?** Đặt change gần business vocabulary thay vì global helper.
3. **Execution flow nào đổi?** Handler → business → store/client → state/effect.
4. **State/contract/config có đổi không?** Nếu có, owner phải mang migration/contract/config proof tương ứng.
5. **Proof nhỏ nhất đủ trả lời là gì?** Service test trước, integration chỉ khi boundary cần.
6. **Change có tăng global complexity/coupling không?** Nếu có, phải chứng minh giá trị vận hành.

## 2. Ownership rules

```text
business/application truth     → <owner>-service/
generic runtime mechanics      → shared/common hoặc shared/base
communication contract         → shared/protobuf
repository engineering         → shared/code
whole-workspace public UX      → root
cross-owner runtime proof      → integration-test
```

Admission test cho `shared/code`: nếu thay toàn business BDSPro bằng một Go microservice product khác mà logic vẫn có nghĩa (generate/build/verify/deploy/toolchain), nó có thể thuộc engineering. Payment settlement, TQD quota, entitlement... thì không.

## 3. Daily development

```bash
make setup
make deps service=<owner>-service
make dev service=<owner>-service
```

`make up` chạy toàn bộ application native; `make dev` chuyển riêng owner đang sửa
sang Air trong terminal. Không thêm Docker build vào edit loop chỉ để “đồng bộ
production”. Nếu cần proof image/package, dùng `make integration-up` hoặc
`make rebuild` tường minh.

Service Makefile là local application API và dùng shared contract tại `shared/code/development/service.mk`/`migration.mk`.

## 4. Service design

Một service phải đọc được như một application độc lập:

```text
entrypoint/composition
→ config
→ transport handler
→ business/usecase
→ store/client
→ database/outbox/job/effect
```

Ưu tiên business/capability locality. Không tạo abstraction, process, container hoặc service mới chỉ vì “có thể cần sau này”.

Một process role mới không tự động là service mới. Nếu API/worker/cron/outbox actor vẫn cùng business owner và cùng lifecycle correctness, giữ trong service; chỉ tách khi có lý do ownership/deployment/scaling/failure isolation rõ.

## 5. Config

```text
safe structural defaults → config/runtime.yml hoặc code
runtime/deploy topology   → ENV
secrets                   → ENV/secret store
business data/rules       → domain/DB
```

Application components nhận typed config; business code không tự tìm ENV rải rác.

## 6. Database changes

Mỗi database owner có một canonical migration history. Không sửa schema production bằng startup side effect/ORM auto-update.

```bash
make migrate service=payment-service
make verify-migrations
```

Migration Compose containers chỉ là one-shot executor của cùng authority.

## 7. Contracts

Cross-service communication phải đi qua contract owner. Khi đổi protobuf:

```bash
make generate
make test service=<producer-or-consumer>-service
make verify
```

Generated code không phải source of truth; schema mới là truth.

## 8. Proof ladder

```text
service/unit test → local code/business proof
verify            → source/repository/tooling proof
smoke             → public runtime boundary alive
E2E               → cross-service business flow
accept            → final local candidate proof
```

Luôn chạy proof nhỏ nhất trước. `make accept` không thay thế fast feedback hằng ngày.

## 9. Review standard

Review không chỉ nhìn syntax. Reviewer phải kiểm:

- đúng owner/capability;
- dependency direction có rõ;
- business có leak vào `shared`/engineering không;
- config/migration/contract có đi cùng change;
- test chứng minh behavior thực, không assertion giả;
- change có giảm hoặc ít nhất không làm xấu code health;
- docs/public commands còn đúng.

## 10. Definition of Done

Tối thiểu cho service-local change:

```bash
make test service=<owner>-service
```

Cross-service/shared/config/migration change:

```bash
make verify
```

Release candidate/local acceptance:

```bash
make accept
```
