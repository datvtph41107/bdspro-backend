# BDSPro R5 Next Slice Selection — Social Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: SELECTED / WRITER PRECHECK NEXT

## Authority

Active source branch:
`refactor/canonical-observability-errors-a6d0722a`

Exact authority SHA:
`fca4682587d93cbc436f2466244ee8bd03b8b1b9`

Live GitHub reconciliation immediately before selection:
- status: identical
- ahead: 0
- behind: 0

## Selected bounded slice

Service:
`social-service`

Concern:
R5 canonical logging migration only.

Current canonical inventory at the authority SHA:
- Social logging debt: 6
- all 6 are `go.legacy_std_log`
- no direct `go.third_party_logger` finding
- no existing canonical `common/logging` / `log/slog` usage

Finding locations:
- `social-service/cmd/grpc/main.go` — 1
- `social-service/infra/postgre/news_feed_postgre.go` — 2
- `social-service/internal/usecase/news_feed.go` — 3

Composition:
- process root: `social-service/main.go`
- gRPC command: `social-service/cmd/grpc/main.go`
- HTTP command exists but currently has no Run/RunE implementation

Existing tests:
- `social-service/infra/service/comment_error_test.go`
- `social-service/infra/service/news_feed_share_error_test.go`
- `social-service/infra/service/report_error_test.go`

## Why Social is selected before Assistant

Assistant is also small, but 4 of its 8 std-log findings are `log.Fatalf` inside the gRPC composition root.

Changing Fatal behavior into ordinary error returns can alter process-exit/defer/cleanup semantics. R5 Logging must not smuggle lifecycle refactoring into a logging slice.

Social's active gRPC composition already propagates startup/serve failures through Cobra `RunE`, so canonical logging can be introduced at the process root with lower lifecycle risk.

Selection is therefore based on bounded semantic risk, not smallest debt count alone.

## Authorized implementation shape

Expected source owner:
- `social-service/main.go` configures canonical `common/logging` once for the process.
- application code uses `log/slog`.
- existing startup and business flow control remains unchanged.

Expected logging replacements:
- gRPC listen event -> structured slog
- NewsFeed decode failures -> structured slog with error context
- SyncNewsFeed start/error/success -> structured slog; preserve existing control flow

The implementation may use `logging.WithComponent` when it adds useful ownership/component context, but must not create a second logging abstraction.

## Explicit no-touch / no-scope-expansion

Do NOT:
- change `shared/protobuf/**`
- change `organization-service/**`
- change `map-service/**`
- change Error/Response design or migrate ReturnError
- change scheduler timing/lifecycle
- change gRPC/HTTP routing semantics
- change DB query behavior
- change timezone behavior
- change generated Wire output unless `make wire` deterministically requires it
- alter `shared/code/deploy.sh`
- add the Social zero ratchet before zero debt is proved
- publish/push before exact-SHA local proof

## Required proof sequence

1. writer precheck at exact authority SHA
2. bounded source mutation
3. gofmt
4. targeted Social tests/build
5. canonical observability audit
6. prove Social `go.legacy_std_log = 0`
7. only then add Social zero-debt ratchet + ratchet test
8. rerun audit with ratchet enforcement
9. protected-path/deploy invariant proof
10. commit candidate
11. detached exact-SHA proof
12. safe publication
13. hosted exact-SHA proof
14. durable checkpoint sync

Writer must remain the only source writer.

FINAL ACCEPTED = NO.
