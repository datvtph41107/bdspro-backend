# BDSPro Repository Engineering Application

`shared/code` là owner của **engineering workflow**, không phải global business layer.

## Owns

- pinned Go/Buf/protoc/Wire/migrate/Air/grpcurl toolchain;
- local `.env` generation và config guards;
- shared service Make contracts;
- protobuf/Wire/Swagger/source generation mechanics;
- repository verification/build helpers;
- native application-process supervision and local backing-infrastructure mechanics;
- explicit Compose image/package integration mechanics;
- deployment/release scripts that are still part of repository operations;
- legacy/manual engineering utilities kept outside the root.

## Does not own

- Payment/TQD/User/Organization business rules;
- service database schema ownership;
- service runtime config semantics;
- protobuf contract semantics;
- cross-service authorization/business policy.

## Public vs internal API

Normal developer starts at repository root:

```bash
make setup
make deps service=payment-service
make dev service=payment-service
make test service=payment-service
make verify
```

Root `Makefile` is intentionally a tiny front door and includes `shared/code/development/root.mk`. Service Makefiles reuse `shared/code/development/service.mk` and `migration.mk`.

Repository maintainers may work here directly for generation/release compatibility targets, but there must be only **one implementation** behind root/service wrappers.

## Layout

```text
shared/code/
├── development/    root/service Make implementation, toolchain, config, verification
├── script/         generation helpers
├── gen/            engineering generators
├── monitor/        development process utility
├── legacy/         retired compatibility utilities
├── manual/         manually-invoked diagnostic utilities
├── Makefile        generation/engineering targets
├── deploy.sh
└── release.sh
```

## Admission rule

Nếu thay BDSPro bằng một Go microservice product khác mà logic vẫn có ý nghĩa (build/generate/dev/verify/release), nó có thể thuộc engineering. Nếu logic nói về order, quota, entitlement, settlement, organization permission... nó phải ở business owner.
