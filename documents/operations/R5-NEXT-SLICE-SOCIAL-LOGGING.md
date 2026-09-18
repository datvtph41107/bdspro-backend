# BDSPro R5 Next Slice Selection — Social Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: DETACHED EXACT-SHA PROOF PASS / SAFE PUBLICATION AUTHORIZED

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


## Writer precheck proof

User-local writer report:
`writer-social-precheck-20260918-133831.txt`

Verified:
- HEAD exact authority PASS;
- branch `local/r5-social-canonical-logging` PASS;
- working tree clean PASS;
- ahead=0 / behind=0 versus authority;
- no worktree diff;
- no index diff;
- final fail count 0.

Live GitHub was reconciled again immediately after this report and the active refactor branch remained identical to authority SHA.

Next authorized transition:
bounded Social source mutation + gofmt + targeted local proof + canonical audit zero-proof.

No commit, ratchet, or push is authorized until that proof is reviewed.


## Mutation attempt and proof-environment finding

Primary report:
`social-logging-mutation-zero-proof-20260918-134454.txt`

Second accidental rerun:
`social-logging-mutation-zero-proof-20260918-134542.txt`

The second report correctly stopped at the pre-mutation dirty-tree guard because the first run had already applied the intended four-file mutation. Do not reset or replay the mutation.

First-run source result:
- exact pre-mutation writer gate PASS;
- mutation applied only to the four authorized Social files;
- gofmt PASS;
- exact changed-file scope PASS;
- 4 files changed, 44 insertions, 16 deletions;
- no ratchet, commit, or push occurred.

Proof failure classification:
- `go test ./...` failed because ignored generated protobuf packages under `shared/protobuf/types/**` were absent in the fresh writer worktree;
- `go build ./...` failed for the same generated-contract absence;
- canonical audit was interrupted by `KeyboardInterrupt` before completion;
- this is NOT evidence that the logging patch failed compilation after canonical protobuf materialization.

Repository contract reconciliation:
- root `.gitignore` intentionally ignores `/shared/protobuf/types/`;
- tracked authority SHA does not contain generated `shared/protobuf/types/**`;
- canonical CI/acceptance workflows materialize protobuf before Go tests;
- repository-owned generator is `make -C shared/code buf-all`;
- `shared/code/Makefile` generates shared, crm, bdspro, hub, chat, organization, social, assistant, user/auth, payment, tqd/operation, notification and transaction contracts;
- generated protobuf output is proof-environment material, not tracked source authority.

Next authorized action:
1. preserve the current four-file Social patch;
2. materialize canonical ignored protobuf contracts with `make -C shared/code buf-all`;
3. verify no tracked path outside the four authorized Social files changed;
4. rerun canonical Social test/build;
5. rerun canonical audit to completion;
6. prove Social legacy std-log and third-party logger debt are both zero;
7. only after successful zero proof may the Social ratchet be added.

Still forbidden:
- reset/clean of the writer;
- tracked mutation under `shared/protobuf/**`;
- ratchet before zero proof;
- commit or push before proof review.


## Social zero proof — PASS

Primary proof report:
`social-proof-after-buf-20260918-135244.txt`

Duplicate uploaded report:
`social-proof-after-buf-20260918-135353.txt`

The two uploaded reports are byte-identical:
SHA256 `d6a4bf82c73c0dea4188ff63563d171a263e459a515a90461bc0a83cbd381847`.

Verified result:
- canonical protobuf materialization PASS;
- required generated contracts exist;
- generated protobuf produced no tracked source mutation;
- Social test PASS;
- Social build PASS;
- canonical observability audit PASS;
- total repository debt moved 721 -> 715;
- Social total debt = 0;
- Social legacy std-log debt = 0;
- Social third-party logger debt = 0;
- protected tracked diff empty;
- deploy.sh SHA unchanged;
- final tracked source scope remains exactly four Social source files;
- FINAL_FAIL_COUNT=0.

Ratchet decision:
- authorize only `("go.legacy_std_log", "social-service")` in this slice;
- do not add a Social third-party-logger ratchet because that category was already zero before this migration and was not retired debt in this slice;
- add registration + regression-enforcement test;
- rerun audit with `--enforce-ratchets`;
- rerun Social test/build and protected/deploy invariants;
- no commit or push until ratchet proof is reviewed.


## Social ratchet proof — PASS

Proof report:
`social-ratchet-proof-20260918-135813.txt`

Verified:
- Social legacy std-log zero ratchet registered;
- registration + regression-enforcement tests added;
- audit-tool unit tests: 30 PASS;
- `audit-observability-errors.py --enforce-ratchets`: PASS;
- `go.legacy_std_log@social-service = 0`;
- Social canonical test PASS;
- Social canonical build PASS;
- protected tracked diff empty;
- `shared/code/deploy.sh` SHA unchanged;
- exact tracked scope is six files: four Social source files plus audit script and audit test;
- FINAL_FAIL_COUNT=0.

Ratchet proof closes the pre-commit implementation gate.

Candidate commit is now authorized with exact expected parent:
`fca4682587d93cbc436f2466244ee8bd03b8b1b9`.

Commit must:
- stage only the six reviewed files;
- have the authority SHA as its single parent;
- not include ignored generated protobuf/build artifacts;
- not push;
- record candidate SHA and tree for detached exact-SHA proof.

Next gate after commit:
detached proof worktree at the immutable candidate SHA, canonical generation/test/build/audit/ratchet/protected proof, then publication only if that exact-SHA proof passes.


## Candidate commit — CREATED

Commit report:
`social-candidate-commit-20260918-140314.txt`

Immutable candidate:
- SHA: `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- tree: `91ccd0d547227191bc20df4498bdbb4ff723c0a2`;
- parent: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`;
- single parent: PASS;
- commit message: `feat(social): adopt canonical logging`;
- committed file set: exactly the six reviewed files;
- writer tracked state clean after commit;
- candidate is exactly one commit ahead of authority;
- pushed: NO;
- FINAL_FAIL_COUNT=0.

Live GitHub reconciliation after candidate creation:
- active refactor branch remains identical to authority SHA;
- candidate is still local-only and therefore must be proven detached before publication.

Next authorized gate:
create a fresh detached proof worktree at candidate SHA and run exact-SHA canonical proof:
1. verify detached HEAD = candidate;
2. materialize canonical protobuf with `make -C shared/code buf-all`;
3. verify generation creates no tracked source delta;
4. run audit-tool unit tests;
5. run audit with `--enforce-ratchets`;
6. prove `go.legacy_std_log@social-service = 0`;
7. run canonical Social test/build;
8. prove protected paths and deploy checksum;
9. leave proof HEAD detached and tracked-clean.

Publication remains forbidden until detached exact-SHA proof is reviewed.


## Detached proof partial upload

Uploaded report:
`social-detached-exact-sha-proof-20260918-141507.txt`

Verified before truncation:
- candidate SHA and tree identity PASS;
- candidate parent and six-file set PASS;
- fresh proof worktree created detached at candidate SHA;
- canonical protobuf generation PASS;
- tracked state after generation CLEAN;
- audit-tool unit tests PASS;
- audit with zero-ratchet enforcement PASS;
- Social legacy std-log ratchet count = 0;
- total debt = 715.

The uploaded report ends while canonical Social test is starting and contains no test exit, build result, protected/deploy proof, final detached-state proof, or FINAL_FAIL_COUNT.

Do not rerun the mutation or detached-proof setup. Determine whether the original proof process is still running; if it completes, re-upload the same report file after completion. If the process has stopped, diagnose from the existing worktree/report without resetting or deleting proof state.


## Detached exact-SHA proof — PASS

Completed proof report:
`social-detached-exact-sha-proof-20260918-141507(1).txt`

Immutable candidate:
- SHA: `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- tree: `91ccd0d547227191bc20df4498bdbb4ff723c0a2`;
- parent: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`.

Verified in fresh detached proof worktree:
- candidate identity and six-file set PASS;
- detached HEAD exactly at candidate PASS;
- canonical protobuf materialization PASS;
- generation left tracked state clean;
- audit-tool unit tests: 30 PASS;
- audit with `--enforce-ratchets`: PASS;
- `go.legacy_std_log@social-service = 0`;
- total repository debt = 715;
- canonical Social test PASS;
- canonical Social build PASS;
- protected candidate diff empty;
- deploy.sh byte identity PASS;
- final proof worktree remained detached and tracked-clean at exact candidate SHA;
- PUSHED=NO;
- FINAL_FAIL_COUNT=0.

Live GitHub immediately after proof review:
- active refactor branch remains identical to authority SHA;
- candidate is still unpublished.

Safe publication is now authorized:
- target branch: `refactor/canonical-observability-errors-a6d0722a`;
- publication must be a non-force fast-forward from authority SHA to candidate SHA;
- before push, re-read remote target with `git ls-remote` and require exact authority SHA;
- push exact candidate SHA directly to target branch;
- after push, require remote target equals exact candidate SHA;
- do not force, rebase, amend, merge, or mutate candidate.

Hosted proof after publication:
- workflow `.github/workflows/refactor-observability-errors.yml` is configured to run on pushes to the active refactor branch;
- only completed exact-SHA hosted evidence may close this Social slice.
