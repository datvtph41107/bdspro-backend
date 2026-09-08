# BDSPro Development Architecture Standard

## Purpose

Repo phải tự dạy developer cách làm việc. Một developer học một mental model rồi áp dụng cho mọi service.

## Canonical repository shape

```text
bdspro-backend/
├── README.md
├── Makefile                 # stable front door only
├── compose.yaml             # local integration topology
├── go.work
├── .vscode/
├── documents/
├── shared/
│   ├── code/                # repository engineering application
│   │   ├── development/
│   │   ├── script/
│   │   ├── legacy/
│   │   └── manual/
│   ├── common/
│   ├── base/
│   └── protobuf/
├── integration-test/
└── *-service/
```

Root không có `tools/` song song với `shared/code`; engineering có một physical owner.

## Stable developer vocabulary

```text
setup deps dev test migrate generate
up rebuild status logs smoke test-e2e
audit/verify accept down reset
```

Một developer intent phải có một obvious command. Compatibility targets không được tạo implementation thứ hai.

## Service standard

Service là application owner. Entry/composition phải đọc được; config có một path; business capability giữ locality; persistence/client/worker nằm sau business boundary; migration và tests thuộc owner.

Không bắt buộc mọi service dùng cùng tên folder nếu runtime/business khác nhau. Chuẩn hóa cognition và command, không ép abstraction giả.

## Engineering standard

`shared/code` sở hữu *how repository engineering works*. Root sở hữu *what developer asks*. Service sở hữu *what product does*.

## Daily vs integration

```text
Daily:       all applications native + Docker backing dependencies
Focused:     one native Air process in the IDE terminal
Integration: explicit full Compose image/package topology
Acceptance:  native runtime + recovery + E2E; image proof is independent
```

Dev/prod parity nằm ở code/config semantics/contracts/migrations/artifact proof; không yêu cầu edit loop dùng cùng packaging mechanics production.

## Config standard

```text
ENV/YAML inputs → loader → typed Config → application components
```

Không để business package tự tìm ENV.

## State standard

Mỗi DB owner có một versioned migration authority. Serving process không tự mutation schema mặc định.

## Documentation standard

- root README: 5-minute orientation;
- service README: owner + entrypoints + commands + dependencies;
- CONTRIBUTING: change/review constitution;
- tutorials/how-to/reference/explanation tách theo nhu cầu;
- docs phải phản chiếu source/public commands, không trở thành nguồn truth thứ hai.

## Architectural change gate

Trước khi thêm service/process/container/shared package/tool/config layer, trả lời: nó mua giá trị gì? Chấp nhận khi giảm ambiguity/duplication/blast radius, tăng feedback/reproducibility/correctness/onboarding. Không chấp nhận chỉ vì “trông chuẩn” hoặc “có thể sau này cần”.
