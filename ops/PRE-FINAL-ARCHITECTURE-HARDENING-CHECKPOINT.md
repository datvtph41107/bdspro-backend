# BDSPro Pre-Final Architecture Hardening Checkpoint

Updated: 2026-09-21 Asia/Bangkok

## Status

`PRE_FINAL_ARCHITECTURE_HARDENING=AUTHORIZED/ACTIVE`

Production Promotion remains paused. Official production lineage has not been mutated.

Baseline canonical authority before hardening:
- branch: `final-acceptance/source-canonicalization`
- SHA: `ca98eceb276dca8249b2f1d4d73cdce6248ec7dd`
- tree: `76b77be3745d004a7c59cbfcd32cf544568b5000`
- Fresh-clone Reconstruction: PROVED/CLOSED at run `35588223133`
- FINAL ACCEPTED = NO

## Slice A — Error boundary ownership convergence

Classification:
- common `UnaryErrorInterceptor` is the canonical application-error normalization boundary;
- Payment lacked this common boundary in its unary chain;
- Payment manually projected canonical errors with `_errors.ToGRPC`;
- Payment duplicated common recovery;
- Payment completion logging independently owned ERROR severity and technical error payloads.

Bounded change:
- install common `UnaryErrorInterceptor` in Payment;
- use common recovery and retire Payment-local recovery;
- return canonical application failures unchanged from Payment handlers;
- make Payment completion logging an INFO-only outcome projection using canonical code/reason/grpc_code;
- log canonical operational failures once at common boundary for Internal/Unknown/Unavailable/DataLoss/DeadlineExceeded, preserving wrapped technical cause only in internal logs;
- add `go.direct_grpc_error_projection@payment-service=0` detector/ratchet.

Local immutable evidence:
- source candidate: `f7f4c981da29a48a63d430ba7e8afc3312961ce2`
- source tree: `750855edc374d9fc556d9b089a29dd9cef0b9a05`
- harness-only child: `bec7db04aa2502269eed7018d6c4e3898125c979`
- harness tree: `d21e700f2fd09282f3d75727cafcfbc9853ba471`
- harness delta: only `.github/workflows/refactor-observability-errors.yml`

Local proof PASS:
- common middleware/errors focused tests;
- Payment interceptor/handler focused tests;
- Payment compile-only gate;
- detector unit tests: 57 PASS;
- all zero ratchets PASS, including direct Payment gRPC projection = 0;
- Payment direct `ToGRPC(` = 0;
- Payment-local recovery file removed;
- fresh detached worktree reconstruction via `make setup` PASS;
- detached exact-SHA proof on `bec7db0...` PASS;
- protected paths unchanged: `shared/protobuf`, `organization-service`, `map-service`, `shared/code/deploy.sh`;
- deploy hash preserved: `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`.

## Safe publication

Target branch:
- `refactor/pre-final-architecture-hardening-ca98ece`

Pre-publication proof:
- target branch was absent by `git ls-remote`;
- publication was create-only and non-force.

Tree-preserving mapping:
- local source `f7f4c981...` tree `750855edc374d9fc556d9b089a29dd9cef0b9a05`
  -> remote source commit `58a527cafcd5c87256e43ef41dfe1512338280cb`
  -> same tree `750855edc374d9fc556d9b089a29dd9cef0b9a05`;
- local harness `bec7db0...` tree `d21e700f2fd09282f3d75727cafcfbc9853ba471`
  -> remote branch head `373663d955d2fcffac1cfc96736ce694e8ce5b73`
  -> same tree `d21e700f2fd09282f3d75727cafcfbc9853ba471`.

Remote verification:
- `git ls-remote` reports branch head exactly `373663d955d2fcffac1cfc96736ce694e8ce5b73`;
- GitHub commit object reports tree `d21e700f2fd09282f3d75727cafcfbc9853ba471` and parent `58a527cafcd5c87256e43ef41dfe1512338280cb`.

Environment classification:
- recovery Sprite host prerequisites repaired to `HOST_PREREQS=PASS`;
- Docker daemon remains unavailable in recovery Sprite; Docker-dependent proof remains a hosted-CI responsibility.

Current writer:
- `/home/sprite/work/arch-hardening-ca98ece`
- branch: `work/pre-final-hardening`
- HEAD: `bec7db04aa2502269eed7018d6c4e3898125c979`
- clean

Preserved proof worktrees:
- `/home/sprite/work/proof-pre-final-a-f7f4c98`
- `/home/sprite/work/proof-pre-final-a-bec7db0`
- `/home/sprite/work/proof-pre-final-a2-bec7db0`

## Remaining hardening order

1. Prove hosted exact-SHA for Slice A at remote `373663d...`.
2. Slice B — Redis transport/key-policy/secret ownership separation for proven consumers.
3. Slice C — Payment outbox degradation/backoff/health/idempotency proof and bounded fixes.
4. Slice D — static enforcement for permanent canonical APIs where compatibility state permits.
5. Aggregate exact-SHA proof, canonical workflows, fresh-clone reconstruction, release build/verify/rollback proof.
6. Re-open Production Promotion only after hardening is PROVED/CLOSED.

`SLICE_A_SOURCE_MUTATION=CLOSED`
`SLICE_A_FOCUSED_PROOF=PASS`
`SLICE_A_ZERO_RATCHET=PASS`
`SLICE_A_DETACHED_EXACT_SHA_PROOF=PASS`
`SLICE_A_PUBLICATION=PASS`
`SLICE_A_HOSTED_PROOF=PENDING`

`NEXT_GATE=SLICE_A_HOSTED_EXACT_SHA_PROOF`

`FINAL ACCEPTED=NO`
