# BDSPro R5 Next Slice Selection — Assistant Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: SELECTED / WRITER PRECHECK NEXT

## Authority

Active source branch:
`refactor/canonical-observability-errors-a6d0722a`

Exact authority SHA:
`b55ce3c6d8005e5d7375228196a7cafb10f8ddce`

Exact tree:
`91ccd0d547227191bc20df4498bdbb4ff723c0a2`

Live reconciliation before selection:
- active branch vs exact authority: identical;
- ahead: 0;
- behind: 0;
- hosted workflow #84 / `35319045710`: completed SUCCESS;
- inventory/common-contracts/boundary-contracts: all SUCCESS.

Hosted inventory at this exact SHA:
- total debt: 715;
- Assistant debt: 8;
- all 8 are `go.legacy_std_log`;
- direct Assistant `go.third_party_logger` debt: 0.

## Selected bounded slice

Service:
`assistant-service`

Concern:
R5 canonical logging migration only.

Exact debt locations:
- `assistant-service/main.go` — 1 `log.Println`;
- `assistant-service/cmd/grpc/main.go` — 3 `log.Printf` and 4 `log.Fatalf`.

Composition:
- process root: `assistant-service/main.go`;
- active server entrypoint: `assistant-service/cmd/grpc/main.go`;
- `RunGRPCServer` has exactly one production caller: process root;
- no existing Assistant `common/logging` / `log/slog` usage.

Existing Assistant tests:
- `assistant-service/config/runtime_test.go`;
- `assistant-service/infra/client/provider_mode_test.go`;
- `assistant-service/infra/handler/assistant_product_suggest_error_test.go`.

Dependency reality:
- `common` is a direct module dependency;
- Fabric and zap are indirect only;
- no direct third-party logger usage is present in Assistant source.

## Why Assistant is selected now

Assistant is the smallest remaining non-protected service logging-debt scope after Social:
- 8 findings;
- all are one mechanism: stdlib `log`;
- all are concentrated in two composition/startup files;
- no worker graph, DB logging adapter, direct Fabric logger retirement, or broad usecase sweep is required.

Assistant was intentionally deferred before Social because four findings use `log.Fatalf`, and converting those sites into ordinary returns would change exit/defer/cleanup semantics.

The bounded implementation is now constrained to preserve that lifecycle behavior:
- do NOT convert fatal sites into ordinary returned errors;
- replace each fatal log with a canonical structured error record followed immediately by `os.Exit(1)` at the same decision site;
- this preserves the current immediate process termination behavior of `log.Fatalf`, including the fact that deferred cleanup does not run after those fatal sites;
- non-fatal startup/readiness records become canonical `slog` events;
- process root configures `shared/common/logging` once.

This resolves the earlier semantic-risk objection without smuggling a lifecycle refactor into R5.

## Authorized implementation shape

Expected source owner:
- `assistant-service/main.go` configures canonical `common/logging` once;
- application code uses `log/slog`;
- fallback failure to configure logging may write directly to stderr because canonical logging is unavailable.

Expected changes:
- replace process startup `log.Println` with structured `slog`;
- replace the timezone startup print with canonical structured startup evidence rather than leave a competing logging-like process print;
- replace runtime-config/listen/readiness `log.Printf` events with structured `slog`;
- replace each `log.Fatalf` with `slog.Error(...)` then immediate `os.Exit(1)`;
- preserve `RunGRPCServer()` signature and caller shape;
- preserve current cleanup/defer behavior at fatal sites.

The implementation must not introduce a second logging abstraction.

## Explicit no-touch / no-scope-expansion

Do NOT:
- change `shared/protobuf/**`;
- change `organization-service/**`;
- change `map-service/**`;
- change Error/Response design or migrate `ReturnError`;
- change Assistant provider selection or runtime config semantics;
- change gRPC message limits/interceptors/reflection/service registration;
- change listen/serve lifecycle;
- convert fatal startup sites to recoverable return paths;
- change Wire output unless `make wire` deterministically requires it;
- alter `shared/code/deploy.sh`;
- add Assistant zero-ratchet before zero debt is proved;
- publish/push before exact-SHA local proof.

## Required proof sequence

1. writer precheck at exact authority SHA;
2. bounded two-file Assistant source mutation;
3. gofmt;
4. canonical protobuf materialization when needed by fresh proof environment;
5. targeted Assistant test/build;
6. canonical observability audit;
7. prove Assistant `go.legacy_std_log = 0`;
8. only then add `go.legacy_std_log@assistant-service` zero-ratchet + registration/regression test;
9. rerun audit with ratchet enforcement;
10. protected-path/deploy invariant proof;
11. commit immutable candidate;
12. detached exact-SHA proof;
13. safe non-force publication;
14. hosted exact-SHA proof;
15. durable checkpoint synchronization.

Writer must remain the only source writer.

FINAL ACCEPTED = NO.
