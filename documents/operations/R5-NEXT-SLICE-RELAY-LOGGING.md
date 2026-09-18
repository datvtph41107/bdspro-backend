# BDSPro R5 Next Slice Selection — Relay Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: RELAY ZERO PROOF PASS / DEPENDENCY RETIREMENT PROBE NEXT

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


## Relay writer precheck — PASS

Report:
`relay-writer-precheck-20260918-155824.txt`

Verified:
- remote active branch = exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- writer began at exact authority on prior Assistant branch;
- tracked writer state clean;
- Python execution cache gate clean;
- new writer branch `local/r5-relay-canonical-logging` created successfully;
- final writer HEAD = exact authority;
- ahead = 0;
- behind = 0;
- final worktree clean;
- source mutation = NONE;
- commit = NONE;
- push = NONE;
- FINAL_FAIL_COUNT=0.

Live GitHub was reconciled immediately after this proof and active refactor remains identical to exact authority.

Next authorized action:
- bounded Relay logging mutation only;
- migrate the 5 legacy std-log and 11 direct Fabric logger debt findings;
- preserve functional stdout in the standalone message utility;
- convert the two runtime debug process prints in `wshandler/runtime.go` into canonical logging projections;
- preserve WebSocket/Redis/DB/RPC behavior;
- then gofmt, canonical protobuf materialization, Relay test/build and canonical audit;
- no ratchet, commit or push until Relay zero logging debt is proved and reviewed.


## First Relay mutation attempt — script bug, partial local patch preserved

Report:
`relay-logging-mutation-zero-proof-20260918-161203.txt`

Observed:
- pre-mutation gate PASS at exact authority on `local/r5-relay-canonical-logging`;
- deterministic mutation stopped while processing `relay-service/wshandler/runtime.go`;
- stop message: `unsupported active wsLogger line: wsLogger.Errorf("Failed to upgrade to WebSocket: %v", err)`;
- this is a mutation-script regex escaping bug, not a source/architecture failure.

Important local-state interpretation:
- the mutation script writes each earlier Relay file immediately;
- therefore files processed before `wshandler/runtime.go` may already contain the intended logging migration;
- `wshandler/runtime.go` write happens only after its conversion loop, so the failed loop should have left that file at authority state;
- do not reset, clean, delete, or rerun the whole mutation blindly;
- first verify the exact partial tracked scope, then apply only the missing wshandler continuation with corrected matching and continue zero proof.

No ratchet, commit, or push is authorized.

Live GitHub was reconciled after the failed attempt and remains identical to exact authority `492b94a102e26b8d86575d72cca05b57911c745b`.


## Relay continuation zero-proof — logging zero, test blocked by format vet

Report:
`relay-mutation-continuation-zero-proof-20260918-162216.txt`

Verified source/proof results:
- remote/writer remain at exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- expected seven-file partial patch verified;
- `wshandler/runtime.go` was authority-clean before continuation;
- corrected continuation converted 38 active Fabric logger calls;
- gofmt PASS;
- final tracked scope exactly eight Relay source files;
- direct Relay Fabric logger imports = 0;
- direct Relay stdlib log imports = 0;
- standalone utility functional stdout count = 2;
- standalone utility immediate-exit count = 2;
- runtime debug-print count = 0;
- canonical protobuf materialization PASS with no tracked scope expansion;
- Relay build PASS;
- canonical audit PASS;
- Relay total debt = 0;
- Relay `go.legacy_std_log = 0`;
- Relay `go.third_party_logger = 0`;
- repository debt `707→691`;
- protected diff empty;
- deploy SHA unchanged.

Proof blocker:
- canonical Relay test failed because Go test/vet now sees four `fmt.Sprintf` format/type mismatches introduced while mechanically preserving old Fabric format strings:
  - connected user: `%s` with `uint64 userId`;
  - send-to-user failure: `%s` with `uint64 userId`;
  - room evidence: `%s` with `uint64 roomID`;
  - room-not-found evidence: `%s` with `uint64 roomID`.
- Relay build still passed; this is a logging-format proof defect, not a runtime ownership or architecture change.

Decision:
- preserve the eight-file patch;
- repair only those four format verbs from `%s` to `%d`;
- rerun Relay test/build, canonical audit and all zero/protected invariants;
- do not add ratchets until canonical Relay test also passes;
- no commit/push.

Live GitHub after proof review remains identical to exact authority.


## Relay complete zero proof — PASS

Report:
`relay-format-repair-zero-proof-20260918-162803.txt`

Verified:
- remote and writer remain at exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- tracked mutation scope remains exactly eight Relay source files;
- the four uint64 logging format defects were repaired exactly once each;
- direct Fabric logger imports = 0;
- direct stdlib log imports = 0;
- standalone message utility keeps exactly two functional stdout records and two immediate exits;
- runtime debug-print count = 0;
- canonical protobuf materialization PASS with no tracked scope expansion;
- canonical Relay test PASS;
- canonical Relay build PASS;
- canonical audit PASS;
- Relay total debt = 0;
- Relay `go.legacy_std_log = 0`;
- Relay `go.third_party_logger = 0`;
- repository debt `707→691`;
- protected tracked diff empty;
- deploy SHA unchanged;
- RATCHET_ADDED=NO;
- COMMIT=NO;
- PUSH=NO;
- FINAL_FAIL_COUNT=0.

Retirement decision before ratchet:
- Relay source no longer imports Fabric logging, but authority `relay-service/go.mod` still declares `github.com/hyperledger/fabric` as a direct requirement;
- Payment R5 precedent moved Fabric from direct to indirect when only shared/common still required it;
- do not guess the Relay module mutation;
- first run read-only `go mod tidy -diff` on the proven local Relay patch and capture the exact module retirement delta;
- no source mutation, ratchet, commit or push during this probe.

Live GitHub was reconciled after zero-proof review and remains identical to exact authority.
