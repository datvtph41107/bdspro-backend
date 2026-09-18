# BDSPro Implementation & Acceptance Ledger — Current Canonical Ledger

Updated: 2026-09-18 Asia/Ho_Chi_Minh

Immutable live Git and completed exact-SHA evidence outrank durable text. Canonical authority remains `a6d0722a36090148829ca5ff02f03f409d753bad`; FINAL ACCEPTED = NO.

## Closed acceptance / observability slices

- Search semantic + ratchet CLOSED/PASS.
- File-service semantic CLOSED/PASS at `9302e357...`.
- File legacy std-log ratchet FULL CLOSED/PASS at `93321c11...`; hosted run #75 PASS.
- Auth process logging semantic FULL CLOSED/PASS at `e857fd4d0e34d0f61e5beb33b1e849590e78849b`; local checks PASS; inventory `2793/734`; Auth legacy std-log `0`; hosted run #76 PASS.

## Parked Auth legacy std-log ratchet

Base at prior checkpoint: `e857fd4...`.

Prior intended audit-only mutation was limited to:
- `shared/code/development/audit-observability-errors.py`
- `shared/code/development/test_audit_observability_errors.py`

Prior local evidence: detector 27/27 PASS; enforce-ratchets PASS; inventory unchanged `2793/734`; `go.legacy_std_log@auth-service = 0`; diff check PASS. Generated `__pycache__/` is execution artifact and must not be committed.

No ratchet commit existed. Gate remains PARKED, not rejected.

## Historical ADR A — Kind/Definition/catalog/i18n — SUPERSEDED

Earlier reasoning explored:
- shared `fault.Kind`;
- stable custom machine `Code` / Definitions;
- semantic catalogs;
- localization resources/message keys;
- rich fault Cause/metadata;
- richer public projections.

Useful lesson: separate business meaning, transport, presentation and diagnostics. Rejected as V1 implementation because abstractions/fields did not have sufficiently proven consumers/value and would add navigation/governance/runtime complexity.

Historical files:
- `ERROR-I18N-OPERATING-MODE.md`
- `ERROR-SEMANTIC-CATALOG-CHECKPOINT.md`

## Historical ADR B — ReturnError(message) V1 — SUPERSEDED

`ERROR-RESPONSE-ARCHITECTURE-V1.md` simplified the design to one expected-business write path with a safe message and trigger-based future identity/localization.

This corrected the over-modeling problem but still left an unresolved transport-classification problem: a lower boundary cannot determine NotFound/InvalidArgument/FailedPrecondition/PermissionDenied/Unavailable from message alone without reintroducing a hidden mapper/catalog or another source of truth.

Therefore V1 remains historical reasoning and is superseded for implementation by the Final ADR below.

## Current ADR — Canonical Error / Response FINAL

Authority: `ERROR-RESPONSE-ARCHITECTURE-FINAL.md`.
Status: DESIGN CLOSED; IMPLEMENTATION PENDING.

### Reality

Exact-SHA source inventory at active snapshot `e857fd4...` proved multiple overlapping mechanisms:
- 689 production `_errors.ReturnError` occurrences;
- 12 production `_errors.ThrowError`;
- 33 `fault.New`, 28 `fault.Wrap`, 47 `fault.Validation` production occurrences;
- 779 production `status.Error`;
- 28 production `sharepb.ErrorResponse` references;
- 633 production `err.Error()` occurrences;
- 33 production files importing `common/fault`.

Counts are not automatically debt. They prove overlapping mechanisms and require classification by owner/consumer.

Current `ReturnError(int32,message)` integer is mixed legacy data: HTTP-like values plus custom constants/DTO/runtime values. It is not a clean canonical business identity.

### Problem

Repository-wide developers currently face multiple valid-looking ways to express/transport failures. Business meaning, protocol status, public response, technical cause and observability are repeatedly remapped across helpers/frameworks, producing multiple sources of truth, IDE navigation overhead, inconsistent behavior and long-term technical debt.

`ReturnError(message)` alone also cannot provide deterministic protocol classification without another mapper.

### Decision

Final target has ONE canonical outbound/application Error core and ONE intentional caller-facing constructor:

```go
return _errors.ReturnError(
    codes.SomeCode,
    "Safe message",
)
```

Canonical Error owns only:
- standard gRPC `codes.Code`;
- safe human message.

It does NOT own HTTP status, request/operation/run identity, technical cause/stack, custom Kind, Definition/catalog, locale/message key, generic metadata, or mandatory custom business identity.

Standard gRPC code is declared directly at the occurrence. This avoids both custom Kind and hidden message/status mapping.

### Propagation / projection

- Intentional caller-facing failure -> canonical Error.
- Unexpected technical failure -> ordinary wrapped Go error.
- Domain sentinel -> only when an internal typed consumer needs identity.
- One outer gRPC normalization path handles canonical Error/context/unknown technical failures.
- One Gateway projection maps gRPC status -> HTTP status and serializes the public response.
- Logging/metrics are projections combining the same outcome with execution context; they do not redefine error semantics.
- Expected business rejection is not automatically operational ERROR.

### Public contract

Default public body remains minimal and consumer-driven, e.g. `{ "message": "..." }`.

Correlation is primarily `X-Request-ID` plus structured observability. Technical DB/provider/stack/cause data never leaks publicly.

### DX value

Architecture optimizes:
`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`.

Developer should:
- read code + message directly at occurrence;
- search message for actual decision sites/tests;
- search `ReturnError(` for intentional outbound failures;
- F12 to one core implementation;
- trace request/operation identity from client/runtime response into structured evidence and source.

### Rejected alternatives

- custom `Kind`: duplicates standard gRPC classification;
- `Definition`/catalog registry: adds indirection without current consumer value;
- message-only ReturnError: insufficient protocol semantics without hidden mapper;
- many constructors (`NotFound`, `Conflict`, etc.): enlarges API/mental surface without adding truth;
- permanent `...any` compatibility overload: type-unsafe debt;
- global replacement of every `status.Error`/`err.Error()`: invalid because inventory includes legitimate protocol/internal uses.

### Migration finish line

For each scope:
`INVENTORY -> CLASSIFY OWNER -> CONVERT -> PROVE -> DELETE/RETIRE -> RATCHET`.

Final active business-error architecture must have:
- ONE canonical Error core;
- ONE `ReturnError(codes.Code,message)` caller-facing path;
- ONE gRPC normalization path;
- ONE Gateway HTTP projection/serializer;
- ZERO equally-valid competing outbound business error frameworks.

Retire/delete after zero-use/compatibility proof:
- active business `shared/common/fault/**` path;
- `fault.Kind`/business mappers;
- `_errors.ThrowError` when redundant;
- direct business/usecase `status.Error/status.New`;
- service `*_faults.go` factories/mappers;
- active legacy numeric/custom protobuf response path;
- oversized Problem behavior not required by real clients;
- public technical `err.Error()` leakage;
- duplicate message/code/status mapping tables.

Legitimate Go mechanics remain: ordinary `error`, `%w`, `errors.Is/As`, context errors, internal sentinels with real consumers.

### Compatibility constraint

CURRENT signature: `ReturnError(code int32, message string)`.
FINAL signature: `ReturnError(code codes.Code, message string)`.

Any migration bridge is bounded and must have retirement criteria.

`organization-service/**` is currently NO-TOUCH but contains legacy callers. Full one-core completion requires a later explicit acceptance decision allowing a mechanical protected migration or an acknowledged bounded shim. Do not claim final one-core state while that conflict remains unresolved.

### Proof requirement

Error/Response acceptance requires vertical proof, not only unit tests:
- expected business rejection;
- technical DB/provider failure;
- validation failure;
- auth/security failure;
- background job/work-item failure;
with correct public projection, no leak, correlation to structured evidence, and correct operational severity.

## Runtime/Logging R1 — Repository-owned structured log root

Status: PROVED / CLOSED.

Commit:
- SHA `016912fc4a3556e9d36f61fe9dc2605af2b6cccf`
- message `fix(dev): canonicalize structured log root`
- parent `e857fd4d0e34d0f61e5beb33b1e849590e78849b`
- tree `86df2033d3aa686d65e3e7be4f9c0a0cf7ec7d3f`

Problem proved by source: structured logger default was relative while both native supervisor and Air operate from service CWD, so structured evidence could land under `<service>/.tmp/...` instead of one repository evidence root.

Change:
- native supervisor exports repository-owned absolute `QHPRO_LOG_ROOT`;
- generic service dev mechanics inject the same absolute root;
- source-layout verification guards both launch paths;
- developer operating guide documents raw-vs-structured behavior.

Changed paths only:
- `documents/DEVELOPER-OPERATING-GUIDE.md`
- `shared/code/development/native-stack.sh`
- `shared/code/development/service.mk`
- `shared/code/development/verify-source-layout.sh`

Published exact-SHA proof on recovery Sprite:
- clean checkout at `016912fc...`;
- exact changed-path set confirmed;
- protected paths unchanged;
- deploy SHA byte-identical;
- `git diff --check` PASS;
- shell syntax PASS;
- service-local `make -n dev` resolves absolute `<repo>/.tmp/development/logs`;
- source-layout verification PASS;
- `go test -count=1 ./logging` in `shared/common` PASS.

Hosted proof:
- workflow run #77 / ID `35246142700`;
- exact head SHA `016912fc4a3556e9d36f61fe9dc2605af2b6cccf`;
- workflow COMPLETED / SUCCESS;
- inventory PASS;
- common-contracts PASS;
- boundary-contracts PASS.

Acceptance decision: R1 PROVED / CLOSED. Repository-owned absolute structured log root is now the active development-launch contract. R2 ownership-state inventory is authorized next.

## Runtime/Logging R2 — Native runtime ownership state

Status: PROVED / CLOSED.

Commit:
- SHA `bb8b3caa9c8e115f14796516ea718816eda10273`
- message `fix(dev): model native runtime ownership`
- parent `016912fc4a3556e9d36f61fe9dc2605af2b6cccf`
- tree `27d2c11d337a2c99f15a448c1600905f55cf0cd6`

Reality/problem: prior runtime recognized only supervised PID ownership; foreground DEV appeared DOWN; unrelated port owners were not FOREIGN; ownership/readiness were conflated; root and service-local DEV did not share one lifecycle owner.

Change/value:
- DEV PID truth under `.tmp/development/dev-pids`;
- shared `dev-service.sh` owns foreground DEV lifecycle/signal/PID cleanup;
- `native-stack.sh` owns `SUPERVISED / DEV / FOREIGN / DOWN` classification;
- readiness/port reported separately;
- second DEV fails fast;
- SUPERVISED→DEV transfer defined once;
- `native-up` preflights DEV/FOREIGN;
- `make down` does not kill DEV/FOREIGN;
- R3 mirroring and R4 query intentionally excluded.

Published exact-SHA proof PASS. Hosted run #78 / `35297061769` COMPLETED / SUCCESS: inventory/native ownership contract PASS; common-contracts PASS; boundary-contracts PASS.

Acceptance decision: R2 PROVED / CLOSED. R3 read-only inventory is authorized next.

## Development Runtime + Logging architecture checkpoint

Status: DESIGN CLOSED; IMPLEMENTATION PENDING.
Authority: `DEVELOPMENT-RUNTIME-LOGGING-CHECKPOINT.md`.

Core decisions:
- `make up` whole native supervised local runtime + Docker backing infra;
- `make dev service=X` transfers only X to foreground DEV/Air ownership;
- raw process output distinct from structured application records;
- repository-owned `.tmp` ephemeral evidence;
- structured query separate from `make logs`;
- request/operation correlation bridges runtime symptom to evidence/source;
- local/integration/production share semantic contracts, not identical supervision mechanics.

## Active priority / next acceptance sequence

Current mode: IMPLEMENTATION / PROOF.

1. R1-R4 CLOSED.
2. R5 Payment canonical logging slice CLOSED at `fca4682587d93cbc436f2466244ee8bd03b8b1b9`; hosted #83 SUCCESS.
3. R5 Social canonical logging slice CLOSED at `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`; hosted #84 / `35319045710` SUCCESS.
4. R5 Assistant canonical logging slice CLOSED at `492b94a102e26b8d86575d72cca05b57911c745b`; hosted #85 / `35324490969` SUCCESS.
5. R5 Relay canonical logging slice CLOSED at `c6a9b121946a360d22759bde6a708bc6a35223c2`; hosted #86 / `35333855680` SUCCESS.
6. R5 remains ACTIVE service-by-service; next authorized action is read-only exact-SHA inventory/trace for one bounded non-protected logging owner.
7. For each R5 slice: owner/consumer/value -> bounded change -> zero proof -> retirement -> ratchet -> local exact-SHA proof -> safe publication -> hosted proof.
8. Error/Response FINAL remains implementation-pending until R5 is deliberately stable.
9. Then prove vertical failures, migrate Error/Response service-by-service, retire competing paths at zero, remove compatibility only after consumer proof, resolve protected Organization conflict.
10. FINAL ACCEPTED remains NO.
## Protected invariants

- `shared/protobuf/**` NO-TOUCH
- `organization-service/**` NO-TOUCH until explicit later acceptance decision
- `map-service/**` NO-TOUCH
- preserve `make wire`
- preserve `make buf`
- `shared/code/deploy.sh` byte-identical SHA256 `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`
- production lineage remains separate

## Runtime/Logging R3 — DEV raw process-stream mirror

Status: PROVED / CLOSED.

Commit:
- SHA `647ad4105540c34121590b6452cf258540f70e8e`
- message `fix(dev): mirror foreground DEV process output`
- parent `bb8b3caa9c8e115f14796516ea718816eda10273`
- tree `497b85c90a3ec339d0b5ec071d7be89b715237d3`

Reality:
- foreground DEV Air/child stdout/stderr was terminal-only;
- `make logs` read the repository raw log owned by `native-stack.sh`;
- therefore DEV and SUPERVISED did not share one raw process-stream view.

Decision:
- keep `dev-service.sh` as the direct lifecycle/signal/exit owner;
- append DEV stdout/stderr to the existing repository raw log through Bash process substitution and `tee -a`;
- do not create another supervisor, new log format, or pipeline owner.

Value:
- foreground terminal remains immediate;
- `make logs service=X` now sees the same raw stream in DEV and SUPERVISED;
- child PID/exit ownership is preserved;
- R2 ownership semantics remain unchanged.

Exact-SHA local proof PASS:
- mirror contract PASS;
- stdout/stderr visible to foreground/caller and raw log view;
- non-zero child exit 7 remains 7;
- R2 ownership/TERM regression PASS;
- source-layout PASS;
- common/logging tests PASS;
- protected/deploy invariants PASS.

Hosted run #79 / `35299399757`:
- inventory PASS including R3 mirror contract;
- common-contracts PASS;
- boundary-contracts SUCCESS.

R3 is PROVED/CLOSED; hosted run #79 completed successfully.


## Runtime/Logging R4 — Structured runtime query

Status: PROVED / CLOSED.

Commit:
- SHA `e6e0ad271dbe9d5eb2877520ab9189ed366606d9`
- message `feat(dev): query structured runtime evidence`
- parent `647ad4105540c34121590b6452cf258540f70e8e`
- tree `af02aec4a392a2097618b7e0f97f03d6bfd318a7`

Decision/value:
- keep `make logs` as raw process output;
- add one read-only `make log-query` structured-evidence workflow;
- read only canonical retained `runtime.jsonl`, never channel projections;
- exact-field request/operation/service/event/level/component filters;
- chronological human timeline default, optional JSON, no root-cause inference;
- no logger schema/storage or service business source changes.

Exact-SHA local proof on published SHA PASS:
- query unit contract PASS;
- root Make query contract PASS;
- source-layout PASS;
- common logging/request regressions PASS;
- audit detector + ratchets PASS; debt 734;
- protected diff empty;
- deploy SHA unchanged.

Hosted run #82 / `35301114591`:
- inventory SUCCESS including structured runtime log query contract;
- common-contracts SUCCESS;
- boundary-contracts SUCCESS;
- overall completed/success.

Acceptance decision: R4 PROVED/CLOSED. R5 exact-SHA canonical logger-adoption inventory is authorized next.

## Runtime/Logging R5 — Payment canonical logging adoption

Status: PROVED / CLOSED bounded service slice; R5 phase remains ACTIVE.

Commit:
- SHA `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- parent `e6e0ad271dbe9d5eb2877520ab9189ed366606d9`
- tree `ac2ff2a1a98aec5afa647edb1584683aeb5d545f`
- message `feat(payment): adopt canonical logging`

Before:
- Payment `go.legacy_std_log = 4`
- Payment `go.third_party_logger = 9`
- repository debt = 734.

After:
- Payment `go.legacy_std_log = 0`
- Payment `go.third_party_logger = 0`
- repository debt = 721
- both Payment logging zero-ratchets active.

Implementation:
- canonical process logger via `shared/common/logging`;
- Payment gRPC interceptors/workers project through canonical `slog` with correlation/component fields;
- Payment-local GORM adapter preserves slow/error SQL evidence;
- direct Fabric logger usage retired; Fabric remains indirect through shared middleware and therefore is not ratcheted as module-absent.

Published exact-SHA local proof PASS:
- R1-R4 regressions;
- detector 28/28 + ratchets;
- shared logging/request/requestlog;
- Payment typed boundary tests;
- Payment correlation/recovery/GORM behavior tests;
- Payment compile gates;
- protected diff empty;
- deploy SHA unchanged.

Hosted run #83 / `35303545759`: inventory SUCCESS; common-contracts SUCCESS; boundary-contracts SUCCESS.

Acceptance decision: Payment R5 slice PROVED/CLOSED. Continue R5 with read-only inventory for the next bounded non-protected service.


## Runtime/Logging R5 — Social canonical logging adoption

Status: PROVED / CLOSED bounded service slice; R5 phase remains ACTIVE.

Commit:
- SHA `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`
- parent `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- tree `91ccd0d547227191bc20df4498bdbb4ff723c0a2`
- message `feat(social): adopt canonical logging`

Before:
- Social `go.legacy_std_log = 6`
- Social `go.third_party_logger = 0`
- repository debt = 721.

After:
- Social total logging debt = 0
- repository debt = 715
- `go.legacy_std_log@social-service = 0` ratchet active.

Implementation:
- canonical process logger configured through `shared/common/logging`;
- direct gRPC listen, NewsFeed decode and SyncNewsFeed stdlib logs migrated to structured `slog`;
- scheduler control flow, DB fallback behavior, routing, timezone and protected paths unchanged;
- ratchet added only after zero proof.

Detached exact-SHA local proof PASS:
- canonical protobuf materialization;
- audit detector unit tests 30/30 PASS;
- `--enforce-ratchets` PASS;
- Social test/build PASS;
- protected diff empty;
- deploy SHA unchanged;
- detached tracked-clean exact candidate identity PASS.

Publication:
- remote branch precondition required exact parent authority;
- non-force fast-forward only;
- remote after push exact candidate SHA.

Hosted proof:
- workflow run #84 / `35319045710`;
- exact head SHA `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- workflow COMPLETED / SUCCESS;
- inventory SUCCESS;
- common-contracts SUCCESS;
- boundary-contracts SUCCESS.

Acceptance decision: Social R5 slice PROVED/CLOSED. Continue R5 with read-only inventory for the next bounded non-protected service.


## Runtime/Logging R5 — Assistant canonical logging adoption

Status: PROVED / CLOSED bounded service slice; R5 phase remains ACTIVE.

Commit:
- SHA `492b94a102e26b8d86575d72cca05b57911c745b`
- parent `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`
- tree `ce8ea2254b8715c8467535173f8c5f7182e7b4f2`
- message `feat(assistant): adopt canonical logging`

Before:
- Assistant `go.legacy_std_log = 8`
- Assistant `go.third_party_logger = 0`
- repository debt = 715.

After:
- Assistant total logging debt = 0
- repository debt = 707
- `go.legacy_std_log@assistant-service = 0` ratchet active.

Semantics:
- canonical process logger configured through `shared/common/logging`;
- startup/readiness logging migrated to structured `slog`;
- four former `log.Fatalf` sites preserve immediate-exit behavior through structured error + immediate `os.Exit(1)`;
- protected paths and Error/Response scope unchanged.

Proof:
- detached exact-SHA local proof PASS;
- publication non-force fast-forward only;
- hosted #85 / `35324490969` exact candidate SUCCESS;
- inventory/common-contracts/boundary-contracts all SUCCESS.

Acceptance decision: Assistant R5 slice PROVED/CLOSED. Continue R5 with read-only inventory for the next bounded non-protected service.


## Runtime/Logging R5 — Relay canonical logging adoption

Status: PROVED / CLOSED bounded service slice; R5 phase remains ACTIVE.

Commit:
- SHA `c6a9b121946a360d22759bde6a708bc6a35223c2`
- parent `492b94a102e26b8d86575d72cca05b57911c745b`
- tree `01d08e582fb55759132cc199215303dd65f6d255`
- message `feat(relay): adopt canonical logging`

Before:
- Relay `go.legacy_std_log = 5`
- Relay `go.third_party_logger = 11`
- Relay total logging debt = 16
- repository debt = 707.

After:
- Relay total logging debt = 0
- repository debt = 691
- both Relay logging ratchets active at zero.

Semantics:
- canonical process logger configured through `shared/common/logging`;
- runtime logs migrated to canonical `slog`;
- functional standalone utility stdout preserved;
- immediate failure semantics preserved;
- runtime debug process prints retired into canonical logging;
- Fabric and redis/v9 direct module declarations retired to indirect-only;
- protected paths and Error/Response scope unchanged.

Proof:
- corrected detached exact-SHA local proof PASS;
- publication non-force fast-forward only;
- hosted #86 / `35333855680` exact candidate SUCCESS;
- inventory/common-contracts/boundary-contracts all SUCCESS;
- hosted inventory debt = 691 and Relay debt owner = 0.

Acceptance decision: Relay R5 slice PROVED/CLOSED. Continue R5 with read-only inventory/trace for the next bounded non-protected logging owner.


## Notification R5 hosted closure

Status: PROVED/CLOSED

- exact SHA: `381c29365a00f35737bb0b6078cf0973198dc7c2`
- tree: `fe7c5153c65074bf406134e4e966c9d097a2d706`
- parent: `c6a9b121946a360d22759bde6a708bc6a35223c2`
- hosted run: `35405779633`
- jobs: inventory/common-contracts/boundary-contracts = SUCCESS
- hosted artifact: `10572063928`
- digest: `sha256:8c8b98a33e56dda9aea3db532b4770437f729aa8d8e8a95988da8e453529c61d`
- Notification logging debt: `30 -> 0`
- repository debt: `691 -> 661`
- Notification legacy std-log zero-ratchet: active
- safe non-force publication: PASS
- protected/deploy invariants: PASS

This closes the Notification R5 logging slice only. R5 overall remains ACTIVE. Error/Response implementation remains PENDING.

FINAL ACCEPTED = NO.
