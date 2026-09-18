# BDSPro Runtime / Logging R5 — Canonical Logger Adoption Inventory

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PAYMENT + SOCIAL + ASSISTANT + RELAY + NOTIFICATION SLICES PROVED/CLOSED / CHAT SELECTED / R5 PHASE ACTIVE

Exact implementation authority:
- active SHA `381c29365a00f35737bb0b6078cf0973198dc7c2`
- parent `c6a9b121946a360d22759bde6a708bc6a35223c2`
- tree `fe7c5153c65074bf406134e4e966c9d097a2d706`
- latest hosted run `35405779633` completed SUCCESS
- canonical branch remains `a6d0722a36090148829ca5ff02f03f409d753bad`

## Repository reality

`shared/common/logging` remains the single backend logging owner. `shared/common/requestlog` enriches canonical records from existing request/operation/actor/business-operation context. Canonical process adoption before this slice existed in Auth, Search, File, and Gateway.

## Payment source reality before mutation

Payment had 13 logging-debt findings: 4 `go.legacy_std_log` and 9 `go.third_party_logger`. The single process composition root was Payment gRPC runtime/main. Fabric logging was used only as a logging mechanism. Workers carried context. GORM's stdlib writer was operational evidence and could not simply be silenced.

## Bounded Payment implementation

- configure canonical `shared/common/logging` at Payment process boundary;
- replace Fabric logger usage in Payment gRPC interceptors with canonical `slog` projections;
- replace outbox/fulfillment `log.Printf` with component-scoped canonical logging;
- preserve GORM slow/error SQL evidence through one thin Payment-local GORM `logger.Interface` adapter to canonical `slog`;
- keep request correlation via `common/logging.WithComponent(ctx, ...)`;
- move Fabric from direct Payment dependency to indirect-only because shared middleware still requires it;
- add Payment zero-ratchets only after source inventory proved both categories at zero.

## Exact changed paths

- `.github/workflows/refactor-observability-errors.yml`
- `payment-service/cmd/grpc/runtime.go`
- `payment-service/go.mod`
- `payment-service/infra/postgres/gorm_logger.go`
- `payment-service/infra/postgres/gorm_logger_test.go`
- `payment-service/infra/postgres/open.go`
- `payment-service/infra/worker/fulfillment/worker.go`
- `payment-service/internal/server/grpc/interceptors/auth_interceptor.go`
- `payment-service/internal/server/grpc/interceptors/logging_interceptor.go`
- `payment-service/internal/server/grpc/interceptors/logging_test.go`
- `payment-service/internal/server/grpc/interceptors/recovery_interceptor.go`
- `payment-service/main.go`
- `payment-service/worker/outbox.go`
- `shared/code/development/audit-observability-errors.py`
- `shared/code/development/test_audit_observability_errors.py`

## Proof

- published exact-SHA local proof PASS on `fca46825...`;
- R1-R4 regression contracts PASS;
- audit detector 28/28 PASS;
- enforce-ratchets PASS;
- repository debt `734→721`;
- Payment `go.legacy_std_log 4→0`;
- Payment `go.third_party_logger 9→0`;
- Payment source logging debt = 0;
- shared common logging/request/requestlog tests PASS;
- Payment typed admin and handler tests PASS;
- Payment gRPC correlation/recovery behavior tests PASS;
- Payment GORM slow-query canonical projection test PASS;
- Payment process/worker compile gates PASS;
- protected path diff empty;
- `shared/code/deploy.sh` SHA unchanged;
- hosted run #83: inventory SUCCESS, common-contracts SUCCESS, boundary-contracts SUCCESS.

## Baseline issue observed but not caused by R5

`payment-service/infra/postgres/migrate_contract_test.go` globs `../../migrate/*.up.sql` while current repository migrations are under `payment-service/database/migrations`; that pre-existing test failure is unchanged from the base and is not part of the R5 logging slice. Hosted boundary tests intentionally run the logging-specific GORM test instead.

## Bounded Social implementation

Before Social mutation:
- Social `go.legacy_std_log = 6`;
- Social `go.third_party_logger = 0`;
- repository debt = 721.

Implementation:
- configure canonical `shared/common/logging` once at `social-service/main.go`;
- replace one gRPC listen log, two NewsFeed decode-failure logs and three SyncNewsFeed logs with structured `log/slog`;
- preserve scheduler flow, DB fallback behavior, routing and timezone semantics;
- add `go.legacy_std_log@social-service` zero-ratchet only after zero proof.

After:
- Social total logging debt = 0;
- repository debt = 715;
- Social legacy std-log ratchet = 0.

Exact-SHA proof:
- candidate `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- tree `91ccd0d547227191bc20df4498bdbb4ff723c0a2`;
- parent `fca4682587d93cbc436f2466244ee8bd03b8b1b9`;
- detached local proof PASS;
- protected diff empty and deploy SHA unchanged;
- non-force publication PASS;
- hosted run #84 / `35319045710` inventory/common/boundary SUCCESS.

## Next gate

R5 remains active service-by-service. Perform read-only exact-SHA inventory/trace for the next bounded non-protected logging owner. Do not start Error/Response FINAL implementation until R5 adoption is deliberately declared stable.

FINAL ACCEPTED = NO.


## Next selected slice — Assistant

Authority: `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`.

Hosted inventory:
- Assistant debt = 8;
- all 8 = `go.legacy_std_log`;
- direct Assistant third-party logger debt = 0.

Source concentration:
- `assistant-service/main.go` — 1 legacy std log plus one logging-like process print review;
- `assistant-service/cmd/grpc/main.go` — 7 legacy std logs, including four `log.Fatalf` sites.

Selection constraint:
fatal sites must remain immediate process exits; do not convert them into ordinary returned errors. Canonical migration may use structured `slog.Error` immediately followed by `os.Exit(1)` at the same decision sites.

Durable slice authority: `R5-NEXT-SLICE-ASSISTANT-LOGGING.md`.


## Bounded Assistant implementation

Before Assistant mutation:
- Assistant `go.legacy_std_log = 8`;
- Assistant `go.third_party_logger = 0`;
- repository debt = 715.

Implementation:
- configure canonical `shared/common/logging` once at `assistant-service/main.go`;
- migrate startup/runtime/readiness stdlib logs to structured `log/slog`;
- preserve four fatal startup decision sites as structured error + immediate `os.Exit(1)`;
- preserve provider selection, runtime config semantics, gRPC registration and listen/serve lifecycle;
- add `go.legacy_std_log@assistant-service` zero-ratchet only after zero proof.

After:
- Assistant total logging debt = 0;
- repository debt = 707;
- Assistant legacy std-log ratchet = 0.

Exact-SHA proof:
- candidate `492b94a102e26b8d86575d72cca05b57911c745b`;
- tree `ce8ea2254b8715c8467535173f8c5f7182e7b4f2`;
- parent `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- detached local proof PASS;
- protected diff empty and deploy SHA unchanged;
- non-force publication PASS;
- hosted run #85 / `35324490969` inventory/common/boundary SUCCESS.


## Next selected slice — Relay

Authority: `492b94a102e26b8d86575d72cca05b57911c745b`.

Hosted inventory:
- Relay total logging debt = 16;
- `go.legacy_std_log = 5`;
- `go.third_party_logger = 11`.

Debt-bearing files:
- `relay-service/main.go`
- `relay-service/message/message.go`
- `relay-service/middleware/logging_middleware.go`
- `relay-service/middleware/ws_logging_middleware.go`
- `relay-service/models/models.go`
- `relay-service/redis/runtime.go`
- `relay-service/server/ws_server.go`
- `relay-service/wshandler/runtime.go`

Selection constraints:
- preserve WebSocket lifecycle/routing/auth behavior;
- preserve Redis/DB/RPC semantics;
- preserve standalone message utility stdout as functional output;
- migrate stdlib `log` and direct Fabric `flogging` only;
- evaluate runtime debug `print/fmt.Printf` as logging projections;
- add both Relay logging zero-ratchets only after zero proof.

Durable slice authority: `R5-NEXT-SLICE-RELAY-LOGGING.md`.


## Bounded Relay implementation

Before Relay mutation:
- Relay `go.legacy_std_log = 5`;
- Relay `go.third_party_logger = 11`;
- Relay total logging debt = 16;
- repository debt = 707.

Implementation:
- configure canonical logging at Relay process root;
- migrate direct stdlib/Fabric runtime logging to canonical structured logging;
- retire two runtime debug process prints into canonical logging;
- preserve standalone message utility stdout as functional output;
- preserve immediate failure behavior;
- retire Fabric and redis/v9 direct module declarations to indirect-only as confirmed by post-generation `go mod tidy -diff`;
- add both Relay logging zero-ratchets only after zero proof.

After:
- Relay logging debt = 0;
- repository debt = 691;
- Relay legacy std-log ratchet = 0;
- Relay third-party logger ratchet = 0.

Exact-SHA proof:
- candidate `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- tree `01d08e582fb55759132cc199215303dd65f6d255`;
- parent `492b94a102e26b8d86575d72cca05b57911c745b`;
- corrected detached local proof PASS with protobuf materialized before tidy assertion;
- protected diff empty and deploy SHA unchanged;
- non-force publication PASS;
- hosted run #86 / `35333855680` inventory/common/boundary SUCCESS;
- hosted artifact confirms repository debt 691 and Relay debt owner 0.


## Notification closure

Before:
- Notification total logging debt = 30;
- all 30 = `go.legacy_std_log`;
- repository debt = 691.

After:
- Notification logging debt = 0;
- repository debt = 661;
- Notification legacy std-log zero-ratchet active.

Exact-SHA closure:
- candidate `381c29365a00f35737bb0b6078cf0973198dc7c2`;
- tree `fe7c5153c65074bf406134e4e966c9d097a2d706`;
- corrected detached local proof PASS;
- non-force publication PASS;
- hosted run `35405779633`: inventory/common/boundary SUCCESS;
- hosted artifact `10572063928`, digest `sha256:8c8b98a33e56dda9aea3db532b4770437f729aa8d8e8a95988da8e453529c61d`;
- hosted inventory confirms Notification debt owner 0 and repository debt 661.

## Next selected slice — Chat

Authority:
`381c29365a00f35737bb0b6078cf0973198dc7c2`.

Hosted inventory:
- Chat total logging debt = 27;
- `go.legacy_std_log = 4`;
- `go.third_party_logger = 23`;
- 12 debt-bearing files.

Selection:
- Hub is numerically smaller at 19, but its direct Fabric logger is coupled to `shared/common/middleware.UnaryRecoveryInterceptor(*flogging.FabricLogger)`, and CRM also consumes that shared contract;
- Chat does not use that shared Fabric recovery contract in its gRPC server;
- Chat's HTTP/gRPC/WS/Redis/middleware structure resembles the already-proved Relay migration;
- choose clean ownership boundary over raw occurrence count.

Important Chat constraints:
- preserve config/server fatal termination semantics;
- configure canonical logging at `chat-service/main.go`;
- remove Fabric-only fields/parameters that carry no independent behavior;
- do not broaden sensitive request/response/header/body evidence;
- direct Fabric module retirement only after post-generation `go mod tidy -diff` proves it;
- expected repository debt after full Chat zero proof = 634.

Durable slice authority:
`R5-NEXT-SLICE-CHAT-LOGGING.md`.

Next gate:
writer precheck / branch preparation only. No source mutation yet.

FINAL ACCEPTED = NO.
