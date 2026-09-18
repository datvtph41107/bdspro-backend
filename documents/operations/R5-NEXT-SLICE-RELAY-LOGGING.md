# BDSPro R5 Next Slice Selection — Relay Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PROVED / CLOSED — HOSTED EXACT-SHA #86 SUCCESS

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

## Relay dependency-retirement probe — PASS

Reports:
- `relay-dependency-retirement-probe-20260918-163541.txt`
- `relay-tidy-diff-20260918-163541.txt`

Verified:
- proven eight-file Relay patch remained unchanged;
- direct Relay Fabric Go imports = 0;
- current `relay-service/go.mod` still declares `github.com/hyperledger/fabric` direct;
- `go mod tidy -diff` returned DIFF_PRESENT and mutated no files;
- exact tidy delta moves:
  - `github.com/hyperledger/fabric v2.1.1+incompatible` from direct -> indirect;
  - `github.com/redis/go-redis/v9 v9.21.0` from direct -> indirect;
- no version changes;
- no go.sum delta was requested by the tidy diff;
- `go mod why -m github.com/hyperledger/fabric` still resolves an indirect need through shared/common (`bdspro/infra/redis -> github.com/hyperledger/fabric/common/flogging`);
- module-file hashes were unchanged after the read-only probe;
- tracked source scope remained unchanged;
- probe FINAL_FAIL_COUNT=0.

Decision:
- authorize exactly the tidy-proposed `relay-service/go.mod` retirement delta;
- do not remove either dependency from the module graph entirely;
- after mutation require `go mod tidy -diff` to be empty/exit 0;
- rerun Relay test/build, canonical audit, zero-debt proof, protected/deploy invariants;
- only after this retirement proof passes may Relay zero-ratchets be added;
- no commit/push yet.

Live GitHub after probe review remains identical to exact authority `492b94a102e26b8d86575d72cca05b57911c745b`.

## Relay module-retirement proof — PASS

Report:
`relay-module-retirement-proof-20260918-164741.txt`

Verified:
- writer remains at exact authority `492b94a102e26b8d86575d72cca05b57911c745b` on `local/r5-relay-canonical-logging`;
- exact tidy-proposed module retirement applied;
- `github.com/hyperledger/fabric` direct -> indirect, exactly once;
- `github.com/redis/go-redis/v9` direct -> indirect, exactly once;
- no version changes;
- `go mod tidy -diff` exit 0 with empty output;
- `go.sum` unchanged;
- final tracked source/module scope exactly nine Relay files;
- canonical protobuf materialization PASS with no tracked expansion;
- canonical Relay test PASS;
- canonical Relay build PASS;
- canonical audit PASS;
- Relay total debt = 0;
- Relay `go.legacy_std_log = 0`;
- Relay `go.third_party_logger = 0`;
- repository debt remains 691;
- Fabric and redis/v9 remain justified transitively via shared/common;
- protected tracked diff empty;
- deploy SHA unchanged;
- RATCHET_ADDED=NO;
- COMMIT=NO;
- PUSH=NO;
- FINAL_FAIL_COUNT=0.

Ratchet decision:
- authorize exactly `go.legacy_std_log@relay-service` and `go.third_party_logger@relay-service`;
- add registration plus regression-enforcement tests for both;
- rerun audit-tool tests, `--enforce-ratchets`, tidy-diff, Relay test/build, zero-debt proof and protected/deploy invariants;
- no commit/push until ratchet proof is reviewed.

Live GitHub after proof review remains identical to exact authority `492b94a102e26b8d86575d72cca05b57911c745b`.

## Relay dual-ratchet proof — PASS

Report:
`relay-ratchet-proof-20260918-165324.txt`

Verified:
- writer remains at exact authority `492b94a102e26b8d86575d72cca05b57911c745b` on `local/r5-relay-canonical-logging`;
- pre-ratchet tracked scope exactly nine proven Relay source/module files;
- `go.legacy_std_log@relay-service` added exactly once;
- `go.third_party_logger@relay-service` added exactly once;
- registration + regression-enforcement tests added for both;
- final tracked scope exactly eleven reviewed files;
- audit-tool unit tests: 35 PASS;
- audit with `--enforce-ratchets`: PASS;
- Relay total debt = 0;
- Relay legacy std-log = 0;
- Relay third-party logger = 0;
- repository debt remains 691;
- module graph tidy (`go mod tidy -diff` empty);
- canonical protobuf materialization PASS with no tracked expansion;
- canonical Relay test PASS;
- canonical Relay build PASS;
- direct Fabric logger imports = 0;
- direct stdlib log imports = 0;
- functional utility stdout count = 2;
- utility immediate-exit count = 2;
- runtime debug-print count = 0;
- protected tracked diff empty;
- deploy SHA unchanged;
- COMMIT=NO;
- PUSH=NO;
- FINAL_FAIL_COUNT=0.

Candidate commit is now authorized.

Candidate constraints:
- parent must be exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- stage exactly the eleven reviewed files;
- no generated protobuf/build artifacts or Python cache may enter the commit;
- no amend/rebase/merge;
- no push;
- after commit record candidate SHA/tree and prove exact committed file set before detached exact-SHA proof.

Live GitHub after ratchet-proof review remains identical to exact authority.

## Relay candidate commit — CREATED

Report:
`relay-candidate-commit-20260918-170313.txt`

Immutable candidate:
- SHA: `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- tree: `01d08e582fb55759132cc199215303dd65f6d255`;
- parent: `492b94a102e26b8d86575d72cca05b57911c745b`;
- single parent: PASS;
- commit message: `feat(relay): adopt canonical logging`;
- committed file set: exactly eleven reviewed tracked files;
- writer tracked state clean after commit;
- candidate is exactly one commit ahead of authority and zero behind;
- pushed: NO;
- FINAL_FAIL_COUNT=0.

Live GitHub reconciliation after candidate creation:
- active refactor branch remains identical to exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- Relay candidate remains local-only.

Next authorized gate:
fresh detached exact-SHA proof at candidate SHA:
1. verify exact candidate SHA/tree/parent and eleven-file commit set;
2. materialize canonical protobuf contracts;
3. require generation creates no tracked delta;
4. run audit-tool unit tests;
5. run audit with `--enforce-ratchets` and prove both Relay ratchets remain zero;
6. require Relay module graph tidy;
7. run canonical Relay test/build;
8. prove utility/runtime behavior shape;
9. prove protected paths + deploy checksum;
10. leave proof worktree detached and tracked-clean at exact candidate SHA.

Publication remains forbidden until detached exact-SHA proof is reviewed.

FINAL ACCEPTED = NO.

## First Relay detached exact-SHA proof — proof-order defect

Report:
`relay-detached-exact-sha-proof-20260918-170550.txt`

Verified candidate/proof results:
- remote remained at exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- candidate identity/tree/parent PASS;
- exact eleven-file candidate scope PASS;
- fresh detached worktree created at exact candidate `c6a9b121946a360d22759bde6a708bc6a35223c2` / tree `01d08e582fb55759132cc199215303dd65f6d255`;
- canonical protobuf materialization PASS;
- generated protobuf caused no tracked delta;
- audit-tool tests 35 PASS;
- audit with `--enforce-ratchets` PASS;
- Relay total debt = 0;
- both Relay ratchets remain zero;
- repository debt = 691;
- canonical Relay test/build PASS;
- module retirement classification PASS;
- old direct logging mechanisms = 0;
- utility/runtime behavior shape PASS;
- protected diff empty;
- deploy SHA unchanged;
- final proof worktree remains detached, tracked-clean, exact candidate SHA/tree.

Single blocker:
- the proof script executed `go mod tidy -diff` BEFORE canonical protobuf materialization;
- in a fresh proof worktree generated protobuf sources are not present yet;
- that pre-generation tidy probe proposed removing grpc-gateway/googleapis-api module edges and produced `FINAL_FAIL_COUNT=1`;
- after that failed probe, canonical protobuf materialization completed and every subsequent gate passed.

Interpretation:
- candidate source is not rejected by this report;
- the failure is in proof ordering relative to the generated-source contract;
- durable candidate gate already specified protobuf materialization before the module-tidy assertion;
- do not amend/reset/rebase/recreate candidate;
- keep the existing detached proof worktree and rerun the exact-SHA proof with canonical protobuf materialized before `go mod tidy -diff`.

No publication is authorized until the corrected post-generation tidy reproof passes.

Live GitHub after review remains identical to exact authority.

## Relay corrected detached exact-SHA reproof — PASS

Report:
`relay-detached-post-generation-reproof-20260918-171305.txt`

Verified:
- remote remains at exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- writer remains exact candidate `c6a9b121946a360d22759bde6a708bc6a35223c2` / tree `01d08e582fb55759132cc199215303dd65f6d255`;
- proof worktree remains detached at the same exact candidate/tree;
- canonical protobuf materialization runs before module-tidy assertion and leaves tracked state clean;
- post-generation `go mod tidy -diff` exits 0 with empty output and does not mutate go.mod/go.sum;
- audit-tool unit tests: 35 PASS;
- audit with `--enforce-ratchets`: PASS;
- Relay total debt = 0;
- Relay legacy std-log = 0;
- Relay third-party logger = 0;
- repository debt = 691;
- both Relay zero-ratchets registered exactly once;
- canonical Relay test PASS;
- canonical Relay build PASS;
- Fabric and redis/v9 module retirement classifications remain correct;
- functional utility stdout count = 2;
- utility immediate-exit count = 2;
- runtime debug-print count = 0;
- protected candidate diff empty;
- deploy SHA unchanged;
- final proof remains detached at exact candidate SHA/tree;
- PUSHED=NO;
- FINAL_FAIL_COUNT=0.

Safe publication is now authorized:
- target branch: `refactor/canonical-observability-errors-a6d0722a`;
- remote target must still equal exact authority immediately before push;
- publication must be non-force fast-forward only;
- push exact candidate SHA directly to the target branch;
- after push require remote target = exact candidate;
- preserve writer/proof identities;
- no amend/rebase/merge/force-push.

After safe publication, require hosted exact-SHA workflow success before closing the Relay R5 slice.

FINAL ACCEPTED = NO.


## Relay safe publication + hosted exact-SHA proof — PASS

Publication report:
`relay-safe-publication-20260918-171620.txt`

Publication:
- remote before push = exact authority `492b94a102e26b8d86575d72cca05b57911c745b`;
- candidate = `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- tree = `01d08e582fb55759132cc199215303dd65f6d255`;
- non-force exact-SHA fast-forward push PASS;
- remote after push = exact candidate;
- force push = NO;
- writer/proof immutable identities preserved;
- publication FINAL_FAIL_COUNT=0.

Live GitHub:
- active refactor branch now exactly `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- branch vs candidate = identical, ahead 0 / behind 0.

Hosted exact-SHA proof:
- workflow: `Refactor Observability and Error Contracts`;
- run #86 / ID `35333855680`;
- exact head SHA `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- workflow status: completed;
- workflow conclusion: success;
- inventory: completed/success;
- common-contracts: completed/success;
- boundary-contracts: completed/success.
- hosted inventory artifact: `observability-error-inventory-c6a9b121946a360d22759bde6a708bc6a35223c2`;
- artifact digest: `sha256:464b2003d943b22a22b0bec6d9db29be61e3c0657d27d5857d5cce421eeacaf8`;
- hosted repository debt = 691;
- hosted Relay debt_by_owner = 0;
- `go.legacy_std_log@relay-service = 0`;
- `go.third_party_logger@relay-service = 0`.

Acceptance decision:
Relay R5 canonical logging slice is PROVED/CLOSED.

Result:
- Relay logging debt `16→0`;
- repository debt `707→691`;
- Relay legacy std-log + third-party logger zero-ratchets active;
- Fabric/redis-v9 direct module ownership retired to indirect-only as justified by shared/common;
- functional utility stdout and immediate-exit semantics preserved;
- runtime debug process prints retired into canonical logging;
- protected invariants preserved;
- Error/Response scope not touched.

R5 remains ACTIVE service-by-service.
Next authorized action: read-only exact-SHA inventory/trace for one next bounded non-protected logging owner.

FINAL ACCEPTED = NO.
