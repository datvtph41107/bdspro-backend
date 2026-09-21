# BDSPro Pre-Final Architecture Hardening Checkpoint

Updated: 2026-09-21 Asia/Bangkok

## Status

`PRE_FINAL_ARCHITECTURE_HARDENING=AUTHORIZED/ACTIVE`

Production Promotion remains paused. The official production-lineage repository has not been mutated.

Baseline canonical authority before hardening:
- branch: `final-acceptance/source-canonicalization`
- SHA: `ca98eceb276dca8249b2f1d4d73cdce6248ec7dd`
- tree: `76b77be3745d004a7c59cbfcd32cf544568b5000`
- Fresh-clone Reconstruction: PROVED/CLOSED at run `35588223133`
- FINAL ACCEPTED = NO

## Operating constraints

Operating Model V2 remains mandatory:
- many read-only discovery/proof lanes;
- exactly one source writer;
- architecture/gate decisions and source mutation are serialized;
- classify failed proof before mutation;
- protected/no-touch invariants remain in force;
- preserve all interrupted worktrees/evidence;
- do not mutate official production lineage during hardening.

Doctrine:
`ONE TRUTH PER CONCERN -> ONE OWNER -> ONE PUBLIC NAME -> ONE OBVIOUS ENTRY POINT -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

## Slice A — Error boundary ownership convergence

Classification from exact `ca98ece...` source:
- common `UnaryErrorInterceptor` was already the canonical application-error normalization boundary;
- Payment had no common error interceptor in its unary chain;
- Payment handler manually called `_errors.ToGRPC` for canonical errors;
- Payment had a private recovery interceptor duplicating common recovery;
- Payment completion logger also owned ERROR severity and technical error payloads, creating a second operational-error policy.

Bounded mutation:
- install common `UnaryErrorInterceptor` in Payment unary chain;
- retire Payment-local recovery and use common recovery;
- canonical application errors now return unchanged from Payment handler so the common boundary owns gRPC projection;
- Payment completion logger is INFO-only outcome projection and emits canonical code/reason/grpc_code without owning cause/severity;
- common error boundary logs canonical operational failures for Internal/Unknown/Unavailable/DataLoss/DeadlineExceeded and includes wrapped technical cause only in internal logs;
- expected canonical business errors remain unpromoted;
- add detector + zero ratchet `go.direct_grpc_error_projection@payment-service=0`.

Evidence already PASS:
- focused common middleware/errors tests;
- focused Payment interceptor/handler tests;
- Payment compile gate;
- detector unit tests (57 tests);
- observability/error ratchets, including new Payment direct-projection ratchet = 0;
- zero/removal proof: no direct `ToGRPC(` remains in Payment; Payment-local recovery file retired;
- protected paths unchanged: `shared/protobuf`, `organization-service`, `map-service`, `shared/code/deploy.sh`;
- deploy script hash remains `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`.

Immutable commits:
- source candidate: `f7f4c981da29a48a63d430ba7e8afc3312961ce2`
- source candidate tree: `750855edc374d9fc556d9b089a29dd9cef0b9a05`
- harness-only child: `bec7db04aa2502269eed7018d6c4e3898125c979`
- harness-only child tree: `d21e700f2fd09282f3d75727cafcfbc9853ba471`
- harness delta is exactly one file: `.github/workflows/refactor-observability-errors.yml`
- harness change only enables branch `refactor/pre-final-architecture-hardening-ca98ece` and aligns Payment hosted test names.

Environment classification:
- recovery Sprite originally lacked Pango/FFmpeg prerequisites; host prerequisite installation reached `HOST_PREREQS=PASS`;
- `make doctor` reports Docker daemon unavailable on recovery Sprite. This is an environment capability gap, not a source regression. Hosted/fresh-clone proof must cover Docker-dependent gates.

Current writer:
- `/home/sprite/work/arch-hardening-ca98ece`
- branch `work/pre-final-hardening`
- HEAD `bec7db04aa2502269eed7018d6c4e3898125c979`
- clean

Preserved proof worktrees:
- `/home/sprite/work/proof-pre-final-a-f7f4c98`
- `/home/sprite/work/proof-pre-final-a-bec7db0`
- do not delete/reset/recreate them.

## Remaining hardening order

1. Complete Slice A detached exact-SHA proof on `bec7db0...`.
2. Safe non-force publication to `refactor/pre-final-architecture-hardening-ca98ece`.
3. Hosted exact-SHA proof for Slice A.
4. Slice B — Redis transport/key-policy/secret ownership separation for proven consumers only.
5. Slice C — Payment outbox degradation/backoff/health/idempotency proof and bounded fixes.
6. Slice D — static enforcement for permanent canonical APIs where compatibility state permits.
7. Aggregate exact-SHA proof, canonical workflows, fresh-clone reconstruction, release build/verify/rollback proof.
8. Re-open Production Promotion only after hardening is PROVED/CLOSED.

`SLICE_A_SOURCE_MUTATION=CLOSED`
`SLICE_A_FOCUSED_PROOF=PASS`
`SLICE_A_ZERO_RATCHET=PASS`
`SLICE_A_DETACHED_EXACT_SHA_PROOF=PENDING`
`SLICE_A_PUBLICATION=PENDING`
`SLICE_A_HOSTED_PROOF=PENDING`

`NEXT_GATE=SLICE_A_DETACHED_EXACT_SHA_PROOF`

`FINAL ACCEPTED=NO`