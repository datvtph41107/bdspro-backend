# BDSPro R5 Next Slice Selection — Relay Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: SELECTED / WRITER PRECHECK NEXT

## Authority

Active source branch:
`refactor/canonical-observability-errors-a6d0722a`

Exact authority SHA:
`492b94a102e26b8d86575d72cca05b57911c745b`

Exact tree:
`ce8ea2254b8715c8467535173f8c5f7182e7b4f2`

Hosted proof:
- workflow #85 / `35324490969` completed SUCCESS;
- exact head SHA = authority;
- inventory/common-contracts/boundary-contracts all SUCCESS;
- hosted inventory total debt = 707.

## Selected bounded slice

Service:
`relay-service`

Concern:
R5 canonical logging migration only.

Hosted logging debt at exact authority:
- total Relay debt = 16;
- `go.legacy_std_log = 5`;
- `go.third_party_logger = 11`.

Exact debt-bearing source files:
- `relay-service/main.go`
- `relay-service/message/message.go`
- `relay-service/middleware/logging_middleware.go`
- `relay-service/middleware/ws_logging_middleware.go`
- `relay-service/models/models.go`
- `relay-service/redis/runtime.go`
- `relay-service/server/ws_server.go`
- `relay-service/wshandler/runtime.go`

Additional logging-like review findings in the same bounded area:
- `relay-service/wshandler/runtime.go` contains runtime debug `print/fmt.Printf` output that should be evaluated as logging projection;
- `relay-service/message/message.go` contains `fmt.Println` output whose payload is the intended utility result (plain/base64 message), not runtime logging, and should remain data output unless source trace disproves that role.

## Why Relay is selected now

After closing Assistant, Relay is the smallest remaining non-protected service logging-debt owner:
- 16 logging debt findings;
- no Error/Response debt is required to close this logging slice;
- no Organization or Map protected source is involved;
- logging mechanisms are confined to stdlib `log` and Fabric `flogging`;
- the runtime composition is local to Relay and does not depend on the shared Fabric recovery-interceptor logger contract that makes Hub a less isolated next slice.

Relay has one existing service test:
- `relay-service/middleware/jwt_middleware_contract_test.go`.

Canonical service proof contracts:
- `make -C relay-service test`;
- `make -C relay-service build`.

## Ownership and semantics

Runtime process root:
- `relay-service/main.go`.

Canonical runtime logger owner:
- configure `shared/common/logging` once at Relay process root;
- application/runtime code uses `log/slog`.

Standalone utility:
- `relay-service/message/message.go` is a separate `package main` utility;
- its fatal/error logging may migrate to canonical structured logging while preserving immediate exit behavior;
- its plain/base64 stdout is functional tool output, not diagnostic logging.

Lifecycle constraints:
- preserve `run() error` and root `os.Exit(1)` behavior;
- preserve WebSocket server startup/shutdown, signal handling, Redis actor lifecycle and cleanup;
- preserve request middleware behavior and status capture;
- preserve Redis method return semantics;
- preserve DB SaveMessagesToDB return semantics;
- preserve message utility immediate-failure behavior.

## Explicit no-touch / no-scope-expansion

Do NOT:
- change `shared/protobuf/**`;
- change `organization-service/**`;
- change `map-service/**`;
- change Error/Response contracts, manual JSON response behavior, grpc status mapping or HTTP semantics;
- change WebSocket routing/auth behavior;
- change Redis data semantics;
- change chat/user/notification RPC behavior;
- change deploy bytes;
- remove functional stdout from the message utility merely to reduce review counts;
- add Relay zero-ratchets before zero logging debt is proved;
- commit/push before proof gates.

## Expected ratchet policy

Only after zero proof:
- add `go.legacy_std_log@relay-service`;
- add `go.third_party_logger@relay-service`;
- add registration + regression-enforcement tests for both.

Both categories are migration-owned Relay debt and therefore both must reach zero before ratchet authorization.

## Required proof sequence

1. writer precheck at exact authority SHA;
2. bounded Relay logging mutation;
3. gofmt;
4. canonical protobuf materialization if fresh proof environment requires it;
5. Relay test/build;
6. canonical audit;
7. prove Relay legacy std-log = 0 and third-party logger = 0;
8. prove runtime debug process-print decisions are explicit;
9. only then add both Relay zero-ratchets + tests;
10. rerun audit with ratchet enforcement;
11. protected/deploy invariant proof;
12. commit immutable candidate;
13. detached exact-SHA proof;
14. safe non-force publication;
15. hosted exact-SHA proof;
16. durable checkpoint synchronization.

Writer remains the only source writer.

FINAL ACCEPTED = NO.
