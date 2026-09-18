# BDSPro R5 Next Slice Selection — Chat Service Canonical Logging

Updated: 2026-09-19 Asia/Bangkok
Status: CHAT SELECTED / WRITER PRECHECK NEXT

## 1. Exact authority

Active refactor branch:
`refactor/canonical-observability-errors-a6d0722a`

Exact authority SHA:
`381c29365a00f35737bb0b6078cf0973198dc7c2`

Exact tree:
`fe7c5153c65074bf406134e4e966c9d097a2d706`

Authority commit:
`feat(notification): adopt canonical logging`

Hosted authority proof:
- run id `35405779633`;
- inventory/common-contracts/boundary-contracts all SUCCESS;
- artifact id `10572063928`;
- artifact digest `sha256:8c8b98a33e56dda9aea3db532b4770437f729aa8d8e8a95988da8e453529c61d`;
- repository debt = `661`;
- Notification debt = `0`.

## 2. Fresh hosted inventory

Chat logging debt at exact authority:
- total = `27`;
- `go.legacy_std_log = 4`;
- `go.third_party_logger = 23`.

Debt-bearing files:
1. `chat-service/config/config.go` — 2 legacy std-log
2. `chat-service/infra/delivery/wshandler/runtime.go` — 2 third-party logger findings
3. `chat-service/infra/handler/chat_handler.go` — 3 third-party logger findings
4. `chat-service/infra/handler/redis_client_worker.go` — 2 third-party logger findings
5. `chat-service/infra/middleware/logging_middleware.go` — 2 third-party logger findings
6. `chat-service/infra/middleware/ws_logging_middleware.go` — 2 third-party logger findings
7. `chat-service/infra/redis/runtime.go` — 3 third-party logger findings
8. `chat-service/infra/server/grpc_server.go` — 2 third-party logger findings
9. `chat-service/infra/server/http_server.go` — 2 third-party logger findings
10. `chat-service/infra/server/swagger.go` — 2 legacy std-log
11. `chat-service/internal/usecases/read_mark_usecase.go` — 3 third-party logger findings
12. `chat-service/scripts/add_swagger_security.go` — 2 third-party logger findings

## 3. Why Chat is selected before Hub

Hub is numerically smaller at 19 logging-debt findings, but exact source trace proves Hub's direct Fabric logger is part of a shared runtime contract:
- `hub-service/initial/startup.go` owns `*flogging.FabricLogger`;
- `hub-service/cmd/grpc/main.go` passes `app.Logger` into `common/middleware.UnaryRecoveryInterceptor`;
- `shared/common/middleware/recovery_interceptor.go` itself accepts `*flogging.FabricLogger`;
- CRM also consumes the same shared recovery contract.

Therefore a true Hub zero-debt migration is not a Hub-only slice; it requires a shared contract decision and at least one additional consumer migration/adapter proof.

Chat has more findings, but its direct Fabric ownership is local to Chat. Exact source trace does not show use of the shared Fabric recovery-interceptor contract in Chat's gRPC server. Its migration shape is structurally similar to the already-proved Relay slice:
- HTTP server;
- gRPC server;
- WebSocket logging;
- Redis runtime;
- middleware;
- handler/background worker logging.

Selection rule:
`owner/coupling boundary > raw occurrence count`.

## 4. Process/root ownership

Primary process root:
`chat-service/main.go`

Cobra subcommands run under that process root:
- HTTP;
- gRPC;
- Swagger;
- prepare-swagger.

Target process ownership:
- configure `shared/common/logging` once in `chat-service/main.go`;
- do not create a second Chat logging facade;
- normal runtime/application code uses `log/slog` and `common/logging.WithComponent(ctx,...)`.

The Swagger-security helper is a package invoked inside the Chat process; it is not a separate standalone `main` at the exact authority.

## 5. Important behavior constraints

Preserve:
- current Cobra command/subcommand behavior;
- current process exit codes and fatal decision sites;
- config-load immediate-failure behavior;
- gRPC/HTTP server lifecycle and shutdown behavior;
- Redis connection/publish/typing semantics;
- WebSocket routing and broadcast behavior;
- worker goroutine lifecycle;
- chat event-channel semantics;
- DB migration/setup behavior;
- outbound RPC client lifecycle;
- Swagger file mutation behavior;
- business/usecase/error-response semantics.

Fatal/log-fatal sites are behavioral gates, not mere formatting:
- the two `config.LoadConfig` failures currently terminate immediately;
- gRPC/HTTP database/start/listen/serve fatal paths must retain immediate termination behavior unless a separate lifecycle change is explicitly proved;
- Swagger `log.Fatal(http.ListenAndServe(...))` must retain its current terminal behavior.

Do not silently convert fatal sites into recoverable ordinary returns during the logging slice.

## 6. Existing structural retirement opportunities

Read-only trace proves two Fabric API surfaces carry no independent business value:
- `ReadMarkUsecases.logger` is construction-time Fabric ownership and is not needed for domain behavior;
- `NewRedisClientWorker(..., logger *flogging.FabricLogger)` receives a logger parameter that is not stored or used directly.

These should be removed only as part of the bounded logging migration and covered by compile/tests.

`chatHandler.logger` is used by background Redis workers. Replace that usage with canonical component-scoped logging from the worker context rather than introducing another logger field type.

## 7. Module/dependency reality

At exact authority:
- `github.com/hyperledger/fabric v2.1.1+incompatible` is a direct Chat dependency;
- Chat source directly imports Fabric logging;
- after source zero proof, run canonical generation/materialization and `go mod tidy -diff`;
- retire Fabric from direct Chat requirements only if post-generation tidy proves the exact module delta;
- Fabric may remain indirect because shared/common still uses Fabric elsewhere.

Do not hand-edit module classification by assumption.

## 8. Sensitive-evidence constraint

Chat's current HTTP/WS logging includes raw request/response/body/header-style evidence in some middleware.

R5 scope is canonical logging ownership, not an unreviewed security/privacy semantics rewrite.

Therefore:
- do not accidentally add new sensitive fields;
- do not broaden payload/header logging;
- preserve or intentionally reduce existing evidence only when the exact change is reviewed and tests/proof make the behavior explicit;
- never introduce Authorization/token values into new structured fields;
- flag legacy raw body/header evidence as a separate observability/privacy review concern rather than silently expanding it during mechanical logger replacement.

## 9. No-touch / no-scope-expansion

Do NOT:
- change `shared/protobuf/**`;
- change `organization-service/**`;
- change `map-service/**`;
- change Error/Response semantics;
- change Gateway response semantics;
- change chat business decisions;
- change Redis schemas;
- change WebSocket protocol/event types;
- change database schemas/migrations;
- change generated source unless canonical generation proves it;
- change deploy bytes;
- add zero-ratchets before hosted/local zero proof.

## 10. Expected zero-ratchet policy

Only after Chat logging debt is proved zero:
- add `go.legacy_std_log@chat-service = 0`;
- add `go.third_party_logger@chat-service = 0`;
- add registration/regression tests for both.

Expected repository debt after successful Chat slice:
`661 - 27 = 634`.

## 11. Required proof sequence

1. writer precheck at exact authority `381c29365a00f35737bb0b6078cf0973198dc7c2`;
2. bounded Chat source mutation only;
3. gofmt;
4. canonical protobuf/generation materialization required by Chat build contract;
5. require generation leaves tracked source within reviewed scope;
6. Chat tests/build;
7. canonical audit;
8. prove Chat legacy std-log = 0;
9. prove Chat third-party logger = 0;
10. inspect fatal/process/HTTP/gRPC/WS/Redis behavior shape;
11. only then add both Chat zero-ratchets + regression tests;
12. rerun audit-tool tests + `--enforce-ratchets`;
13. post-generation `go mod tidy -diff`;
14. protected/deploy invariant proof;
15. immutable candidate commit;
16. fresh detached exact-SHA proof;
17. safe non-force publication;
18. hosted exact-SHA proof + inventory artifact;
19. durable checkpoint synchronization.

Writer remains the only source writer.

No source mutation is authorized by this selection document itself. Next gate is writer precheck/branch preparation only.

FINAL ACCEPTED = NO.
