# BDSPro R5 Next Slice Selection — Notification Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: NOTIFICATION MUTATION PARTIAL / FIREBASE DUPLICATE-MATCH SCRIPT BUG / CONTINUATION REPAIR NEXT

## Authority

Active source branch:
`refactor/canonical-observability-errors-a6d0722a`

Exact authority SHA:
`c6a9b121946a360d22759bde6a708bc6a35223c2`

Exact tree:
`01d08e582fb55759132cc199215303dd65f6d255`

Hosted proof:
- workflow #86 / `35333855680` completed SUCCESS;
- exact head SHA = authority;
- inventory/common-contracts/boundary-contracts all SUCCESS;
- hosted inventory total debt = 691.

## Selected bounded slice

Service:
`notification-service`

Concern:
R5 canonical logging migration only.

Hosted logging debt at exact authority:
- total Notification debt = 30;
- `go.legacy_std_log = 30`;
- `go.third_party_logger = 0`.

Exact debt-bearing source files:
- `notification-service/cmd/delivery-worker/main.go` — 1
- `notification-service/cmd/grpc_server.go` — 4
- `notification-service/cmd/payment-event-worker/main.go` — 1
- `notification-service/infra/broker/rabbitmq/payment_completed_supervisor.go` — 1
- `notification-service/infra/firebase/provider.go` — 5
- `notification-service/infra/handler/internal_handler.go` — 2
- `notification-service/infra/worker/delivery/worker.go` — 1
- `notification-service/internal/usecase/notification_usecase.go` — 15

## Why Notification is selected now

The smallest remaining non-protected owner by count is Hub (19), but Hub still carries direct Fabric logger ownership through `hub-service/initial/startup.go` and passes that logger into the shared recovery-interceptor contract. Moving Hub cleanly would therefore require either shared middleware contract migration or a bounded adapter decision.

Chat has 27 debt findings but spans 12 files and mixes 4 legacy std-log findings with 23 direct Fabric logger findings across HTTP, gRPC, WebSocket, Redis and utility/script roots.

Notification has 30 debt findings, all in one category (`go.legacy_std_log`) across eight files, with no direct third-party logger debt. The migration can stay inside Notification source plus the audit ratchet/test after zero proof. This is a cleaner ownership boundary than Hub/Chat even though its raw count is slightly higher.

## Process roots / logging ownership

Primary process root:
- `notification-service/main.go` owns the Cobra process that runs `cmd/grpc_server.go`.
- Configure `shared/common/logging` once there for the canonical service process.

Compatibility/recovery worker process roots:
- `notification-service/cmd/delivery-worker/main.go`
- `notification-service/cmd/payment-event-worker/main.go`

Each standalone worker is independently buildable and must configure the same canonical logging owner at its own process root.

Runtime/application code uses `log/slog` / `common/logging.WithComponent(ctx, ...)` as appropriate. Do not introduce a second logging facade.

## Behavior constraints

Preserve:
- Cobra command behavior and root process exit behavior;
- gRPC listener/server lifecycle;
- `common/process.Run` actor ownership and cancellation behavior;
- RabbitMQ reconnect/backoff behavior;
- delivery worker retry/poll behavior;
- Firebase send result/error semantics;
- database history writes;
- Redis publish semantics;
- online/offline detection behavior;
- Auth token lookup behavior;
- notification merge/removal semantics;
- public/internal gRPC response semantics;
- compatibility worker immediate exit behavior.

Two current `log.Println("dto", dto)` and two `log.Println("ismerge=", ...)` sites are diagnostic process output and should become structured debug/info projections without changing handler/usecase decisions.

## Explicit no-touch / no-scope-expansion

Do NOT:
- change `shared/protobuf/**`;
- change `organization-service/**`;
- change `map-service/**`;
- change Error/Response contracts, gRPC status mapping, manual JSON responses or HTTP semantics;
- change notification delivery durability/retry classification;
- change RabbitMQ exchange/queue/consumer semantics;
- change Firebase effect classification;
- change Redis channel/payload schema;
- change notification business decisions;
- change deploy bytes;
- change Wire-generated source unless canonical `make wire` materialization proves tracked regeneration is required;
- add Notification zero-ratchet before zero logging debt is proved;
- commit/push before proof gates.

## Module/dependency reality

`notification-service/go.mod` already has Fabric only as indirect. This slice has no direct Fabric logger source debt, so there is no Fabric direct->indirect retirement step analogous to Payment/Relay.

Any module-file change must be justified by a read-only post-generation `go mod tidy -diff` proof; do not change go.mod by assumption.

## Existing proof contracts

Notification has existing tests covering:
- runtime config;
- migration contract;
- RabbitMQ payment-completed supervisor;
- delivery store;
- delivery usecase;
- payment-completed eventing usecase.

Canonical service proof contracts:
- `make -C notification-service test`;
- `make -C notification-service build`;
- `make -C notification-service wire`;
- canonical protobuf materialization before fresh-proof module assertions.

## Expected ratchet policy

Only after zero proof:
- add `go.legacy_std_log@notification-service`;
- add registration + regression-enforcement tests.

No third-party logger ratchet is needed because hosted Notification third-party logger debt is already zero and this slice does not own a third-party logger migration.

## Required proof sequence

1. writer precheck at exact authority SHA;
2. bounded Notification logging mutation;
3. gofmt;
4. canonical protobuf + Wire materialization;
5. require generation leaves tracked source authoritative/clean except explicitly reviewed generated deltas;
6. Notification test/build;
7. canonical audit;
8. prove Notification legacy std-log = 0;
9. inspect/retain all runtime behavior constraints;
10. only then add Notification legacy std-log zero-ratchet + tests;
11. rerun audit-tool tests + `--enforce-ratchets`;
12. post-generation `go mod tidy -diff`;
13. protected/deploy invariant proof;
14. immutable candidate commit;
15. detached exact-SHA proof;
16. safe non-force publication;
17. hosted exact-SHA proof;
18. durable checkpoint synchronization.

Writer remains the only source writer.

FINAL ACCEPTED = NO.


## Notification writer precheck — PASS

Report:
`notification-writer-precheck-20260918-212335.txt`

Verified:
- remote active branch = exact authority `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- writer HEAD = exact authority;
- writer branch = `local/r5-notification-canonical-logging`;
- final writer HEAD = exact authority;
- final writer branch = expected Notification branch;
- source mutation = NONE;
- commit = NONE;
- push = NONE;
- FINAL_FAIL_COUNT=0.

Live GitHub was reconciled immediately after the precheck and the active refactor branch remains identical to the exact authority.

Next authorized action:
- bounded Notification logging mutation only;
- configure canonical logging at the primary process root and both compatibility/recovery worker process roots;
- migrate the 30 hosted legacy std-log findings to canonical slog projections;
- migrate the primary-root diagnostic `fmt.Println(err)` into canonical logging while preserving immediate exit behavior;
- preserve all Notification business, delivery, retry, Firebase, Redis, RPC and process-lifecycle semantics;
- then gofmt, canonical protobuf/Wire materialization, Notification test/build and canonical audit;
- no ratchet, commit or push until Notification legacy std-log is proved zero and reviewed.


## First Notification mutation attempt — duplicate-match script bug

Report:
`notification-logging-mutation-zero-proof-20260918-213924.txt`

Observed:
- pre-mutation authority gate PASS at exact authority `c6a9b121946a360d22759bde6a708bc6a35223c2` on `local/r5-notification-canonical-logging`;
- deterministic mutation stopped in `notification-service/infra/firebase/provider.go`;
- stop message: `firebase topic send error: expected exactly 1 match, found 2`;
- the source contains two identical active `log.Println("Error sending message:", err)` sites (topic and token paths), while the script incorrectly required a unique match for the first replacement;
- this is a mutation-script cardinality bug, not a source/architecture failure.

Partial local-state interpretation:
- files written before the Firebase step may already contain the intended migration:
  - `notification-service/main.go`
  - `notification-service/cmd/grpc_server.go`
  - `notification-service/cmd/delivery-worker/main.go`
  - `notification-service/cmd/payment-event-worker/main.go`
  - `notification-service/infra/broker/rabbitmq/payment_completed_supervisor.go`
- `notification-service/infra/firebase/provider.go` is written only after all of its replacements complete, so the failed duplicate check should have left that file at authority state;
- files after Firebase in the mutation sequence should also remain at authority state;
- do not reset/clean/restart the whole mutation blindly.

Next authorized action:
- verify the exact five-file partial scope;
- verify Firebase and the four later files are still authority-clean;
- continue only the missing Firebase/handler/delivery-loop/usecase mutation with corrected cardinality handling;
- then gofmt all nine files and run protobuf/Wire/test/build/audit/zero-debt/protected proofs;
- no ratchet, commit or push.

Live GitHub was reconciled after the failed attempt and remains identical to exact authority `c6a9b121946a360d22759bde6a708bc6a35223c2`.
