# BDSPro R5 Next Slice Selection — Assistant Service Canonical Logging

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: DETACHED EXACT-SHA PROOF PASS / SAFE PUBLICATION AUTHORIZED

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


## Writer precheck attempt — blocked by untracked cache

Report:
`assistant-writer-precheck-20260918-144455.txt`

Verified:
- remote active branch = exact authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- writer HEAD = exact authority;
- current writer branch = `local/r5-social-canonical-logging`;
- no source mutation occurred.

Blocker:
- writer is not clean because `shared/code/development/__pycache__/` is untracked;
- this is execution artifact from Python audit tooling, not source authority;
- do not commit it;
- archive metadata/listing, remove only this exact cache directory, then rerun writer clean/branch-preparation gate;
- no Assistant source mutation until that gate passes.

Live GitHub was reconciled after this report and the active refactor branch remains identical to exact authority.


## Cleanup/precheck follow-up — branch already prepared

Report:
`assistant-writer-cleanup-precheck-20260918-145100.txt`

Observed:
- HEAD = exact authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- current branch = `local/r5-assistant-canonical-logging`.

The short report stopped before clean-state proof because the follow-up script expected the prior Social branch. This is not a source failure. The Assistant branch already exists/current, so do not recreate or delete it.

Next gate:
- verify only allowed Python cache may remain;
- archive/remove that cache if still present;
- prove writer tracked/untracked state clean;
- prove ahead=0 / behind=0 versus authority;
- no source mutation until that clean verification passes.

Live GitHub remains identical to exact authority after this report.


## Final writer precheck — PASS

Report:
`assistant-writer-final-precheck-20260918-145549.txt`

Verified:
- writer HEAD = exact authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- writer branch = `local/r5-assistant-canonical-logging`;
- worktree clean;
- ahead = 0;
- behind = 0;
- source mutation = none;
- commit = none;
- push = none;
- FINAL_FAIL_COUNT=0.

Live GitHub was reconciled again immediately after this proof and the active refactor branch remains identical to exact authority.

Next authorized action:
- bounded Assistant source mutation only in `assistant-service/main.go` and `assistant-service/cmd/grpc/main.go`;
- preserve all four former `log.Fatalf` sites as structured error log + immediate `os.Exit(1)`;
- replace non-fatal startup/readiness stdlib logs and the process timezone print with canonical `slog`;
- configure `shared/common/logging` once at process root;
- then gofmt, targeted Assistant test/build, canonical audit and zero-debt proof;
- no ratchet, commit or push until zero proof is reviewed.


## Assistant bounded mutation + first proof attempt

Report:
`assistant-logging-mutation-zero-proof-20260918-150122.txt`

Source result:
- pre-mutation gate PASS at exact authority;
- mutation applied only to `assistant-service/main.go` and `assistant-service/cmd/grpc/main.go`;
- gofmt PASS;
- exact tracked source scope PASS;
- former legacy fatal logger count = 0;
- four immediate `os.Exit(1)` sites present in Assistant gRPC composition;
- protected tracked diff empty;
- deploy SHA unchanged.

Logging result:
- canonical audit PASS;
- repository findings = 2767;
- repository debt = 707;
- Assistant total debt = 0;
- Assistant `go.legacy_std_log = 0`;
- Assistant `go.third_party_logger = 0`;
- direct legacy-log search PASS.

Proof-environment blocker:
- Assistant test/build failed because `assistant-service/wire/wire_gen.go` was absent;
- `assistant-service/wire/wire.go` is `wireinject`-only;
- repository `generate-backend` / `make setup` contract generates Wire for Assistant before repository tests/builds;
- Assistant `.gitignore` ignores `wire_gen.go`, so this generated file is proof-environment material, not source authority for this slice.

Interpretation:
- do not reset the two-file source patch;
- do not add ratchet yet because targeted test/build has not passed in a canonical generated environment;
- materialize Assistant Wire with repository-owned generator, prove no tracked scope expansion, rerun Assistant test/build, rerun audit/zero proof, then review ratchet gate.

Live GitHub was reconciled after this report and remains identical to exact authority `b55ce3c6...`.


## Assistant canonical generated-environment zero proof — PASS

Report:
`assistant-wire-test-build-zero-proof-20260918-150831.txt`

Verified:
- writer remains at exact authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- tracked source diff remains exactly two Assistant files;
- repository-owned Assistant Wire generation PASS;
- generated `assistant-service/wire/wire_gen.go` remains untracked/generated proof material;
- no tracked scope expansion after Wire generation;
- canonical Assistant test PASS;
- canonical Assistant build PASS;
- canonical audit PASS;
- repository findings = 2767;
- repository debt = 707;
- Assistant total debt = 0;
- Assistant `go.legacy_std_log = 0`;
- Assistant `go.third_party_logger = 0`;
- former legacy fatal count = 0;
- exactly four immediate `os.Exit(1)` fatal-site replacements remain;
- protected tracked diff empty;
- deploy SHA unchanged;
- FINAL_FAIL_COUNT=0.

Ratchet decision:
- authorize only `("go.legacy_std_log", "assistant-service")`;
- do not add an Assistant third-party-logger ratchet because that category was already zero before this migration;
- add registration + regression-enforcement tests using the established Social ratchet pattern;
- rerun audit-tool tests, `--enforce-ratchets`, Assistant test/build, fatal-lifecycle shape, protected/deploy invariants;
- no commit/push until ratchet proof is reviewed.

Live GitHub was reconciled immediately after this proof and remains identical to exact authority.


## Assistant ratchet proof — PASS

Report:
`assistant-ratchet-proof-20260918-151400.txt`

Verified:
- pre-ratchet source scope = exactly two Assistant source files;
- Assistant legacy std-log ratchet added exactly once;
- registration + regression-enforcement tests added;
- audit-tool unit tests: 32 PASS;
- audit with `--enforce-ratchets`: PASS;
- repository debt remains 707;
- `go.legacy_std_log@assistant-service = 0`;
- Assistant total debt = 0;
- repository-owned Assistant Wire generation PASS;
- canonical Assistant test PASS;
- canonical Assistant build PASS;
- former legacy fatal count = 0;
- exactly four immediate `os.Exit(1)` fatal-site replacements remain;
- final tracked scope = exactly four reviewed files;
- protected tracked diff empty;
- deploy SHA unchanged;
- commit = NO;
- push = NO;
- FINAL_FAIL_COUNT=0.

Live GitHub immediately after proof review:
- active refactor branch remains identical to authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`.

Candidate commit is now authorized.

Candidate commit constraints:
- parent must be exact authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- stage only the four reviewed tracked files;
- generated `assistant-service/wire/wire_gen.go`, protobuf outputs, Python cache and build artifacts must not enter the commit;
- do not amend/rebase/merge;
- do not push;
- after commit record candidate SHA/tree and prove exact committed file set before detached exact-SHA proof.


## Assistant candidate commit — CREATED

Report:
`assistant-candidate-commit-20260918-151700.txt`

Immutable candidate:
- SHA: `492b94a102e26b8d86575d72cca05b57911c745b`;
- tree: `ce8ea2254b8715c8467535173f8c5f7182e7b4f2`;
- parent: `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- single parent: PASS;
- commit message: `feat(assistant): adopt canonical logging`;
- committed file set: exactly four reviewed tracked files;
- writer tracked state clean after commit;
- candidate is exactly one commit ahead of authority and zero behind;
- pushed: NO;
- FINAL_FAIL_COUNT=0.

Live GitHub reconciliation after candidate creation:
- active refactor branch remains identical to authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- candidate remains local-only.

Next authorized gate:
fresh detached exact-SHA proof at candidate SHA:
1. verify exact candidate SHA/tree/parent and four-file commit set;
2. materialize canonical protobuf contracts;
3. materialize Assistant Wire via repository-owned generator;
4. verify generation creates no tracked delta;
5. run audit-tool unit tests;
6. run audit with `--enforce-ratchets` and prove Assistant zero ratchet;
7. run canonical Assistant test/build;
8. prove fatal-lifecycle shape;
9. prove protected paths + deploy checksum;
10. leave proof worktree detached and tracked-clean at exact candidate SHA.

Publication remains forbidden until detached exact-SHA proof is reviewed.


## Assistant detached exact-SHA proof — PASS

Report:
`assistant-detached-exact-sha-proof-20260918-152404.txt`

Immutable candidate:
- SHA: `492b94a102e26b8d86575d72cca05b57911c745b`;
- tree: `ce8ea2254b8715c8467535173f8c5f7182e7b4f2`;
- parent: `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- exact four-file commit set PASS;
- fresh proof worktree created detached at exact candidate SHA/tree.

Fresh detached proof verified:
- canonical protobuf materialization PASS;
- repository-owned Assistant Wire generation PASS;
- generated contracts/Wire created no tracked delta;
- audit-tool unit tests: 32 PASS;
- audit with `--enforce-ratchets`: PASS;
- repository debt = 707;
- Assistant total debt = 0;
- Assistant `go.legacy_std_log = 0`;
- Assistant `go.third_party_logger = 0`;
- Assistant zero-ratchet registration count = 1;
- canonical Assistant test PASS;
- canonical Assistant build PASS;
- legacy fatal count = 0;
- exactly four immediate `os.Exit(1)` replacements remain;
- protected candidate diff empty;
- deploy SHA unchanged;
- final proof worktree remains detached, tracked-clean, exact candidate SHA/tree;
- PUSHED=NO;
- FINAL_FAIL_COUNT=0.

Live GitHub reconciliation immediately after proof review:
- active refactor branch remains identical to exact authority `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`;
- candidate remains unpublished.

Safe publication is now authorized:
- target branch: `refactor/canonical-observability-errors-a6d0722a`;
- remote target must still equal exact authority immediately before push;
- publication must be non-force fast-forward only;
- push exact candidate SHA directly to target branch;
- after push require remote target equals exact candidate SHA;
- do not amend, rebase, merge, or force-push.

After publication:
- wait for hosted workflow on exact candidate SHA;
- only completed exact-SHA hosted evidence may close the Assistant slice.

FINAL ACCEPTED = NO.
