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

## Slice A — Error boundary ownership convergence — PROVED/CLOSED

Local immutable evidence:
- source candidate: `f7f4c981da29a48a63d430ba7e8afc3312961ce2`
- source tree: `750855edc374d9fc556d9b089a29dd9cef0b9a05`
- harness-only child: `bec7db04aa2502269eed7018d6c4e3898125c979`
- harness tree: `d21e700f2fd09282f3d75727cafcfbc9853ba471`

Safe publication mapping:
- local source tree `750855edc...` -> remote source commit `58a527cafcd5c87256e43ef41dfe1512338280cb` -> same tree;
- local harness tree `d21e700f...` -> remote branch head `373663d955d2fcffac1cfc96736ce694e8ce5b73` -> same tree;
- branch: `refactor/pre-final-architecture-hardening-ca98ece`.

Hosted exact-SHA:
- workflow: `Refactor Observability and Error Contracts`
- run number: `103`
- run ID: `35612933219`
- head SHA: `373663d955d2fcffac1cfc96736ce694e8ce5b73`
- status/conclusion: `completed/success`
- `inventory`, `common-contracts`, `boundary-contracts`: SUCCESS
- artifact: `observability-error-inventory-373663d955d2fcffac1cfc96736ce694e8ce5b73`
- artifact ID: `10644069205`
- expired: `false`

`SLICE_A=PROVED/CLOSED`

## Slice B — Redis transport/key-policy/secret ownership

### Classification

Proven pre-mutation:
- `common/redis.RedisService` mixed technical transport/lifecycle with tile-session feature semantics;
- tile-session keyspace `ss:k:`, TTL and AES session-key storage had no semantic owner;
- `sessionEncryptKey` was logged in plaintext;
- User handler retained a process-global `sessionKey`, creating shared mutable session truth;
- common Redis token helpers duplicated User-owned token-cache semantics and had no live external caller;
- TQD tile server used legacy `NewRedisService()` despite having typed Redis runtime configuration.

### Bounded mutation

Implemented:
- `common/redis` remains the technical Redis transport/lifecycle owner;
- new `common/tilesession.Store` owns tile-session keyspace, TTL and Redis persistence semantics;
- production `ss:k:` literal now has one executable owner;
- User creates a fresh tile session per request and no longer stores singleton `sessionKey` state;
- User and TQD consume the semantic tile-session store instead of redefining Redis keys/TTL;
- plaintext tile-session secret logging removed;
- dead common Redis token helpers `SaveToken/GetToken/IsTokenValid/DeleteToken` retired;
- TQD tile server now loads narrow typed Redis runtime config and uses explicit `Open/Close`;
- Wire generation authority extended with service-scoped injection packages; `common/tilesession.NewStore` is scoped to User only;
- detector zero-ratchets added for plaintext tile-session secret logging in `shared/common`, User and TQD.

### Proof

Focused proof PASS:
- Wire generator tests PASS, including existing tests plus service-scope tests;
- `common/redis` compile PASS;
- `common/tilesession` tests PASS;
- User handler tests PASS;
- User Wire + gRPC compile PASS;
- TQD config tests PASS;
- TQD cmd compile PASS;
- observability/error detector tests: 60 PASS;
- all ratchets PASS, including tile-session secret logging = 0;
- zero/removal proof PASS:
  - production `ss:k:` literal count = 1 and owned by `common/tilesession/store.go`;
  - legacy tile-session Redis APIs = 0;
  - dead common Redis token helpers = 0;
  - TQD `_redis.NewRedisService()` = 0;
  - User process-global `sessionKey uint64` = 0;
  - plaintext log of `sessionEncryptKey` = 0;
- User Wire regeneration is hash-stable from generation authority;
- protected paths unchanged: `shared/protobuf`, `organization-service`, `map-service`, `shared/code/deploy.sh`;
- deploy script hash preserved: `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`.

Immutable local evidence:
- source candidate: `cf5e62e5ece1add747c7be77cbd468128bf4ba27`
- source tree: `804cff3330a568346f2d6824dce755e6623f3f03`
- harness-only child: `7a229202f8a0b31c8470cb14faf3ce2bf8dcfb1e`
- harness tree: `cb357507b528beef25a799aa346278c1d0ed7798`
- harness delta from source candidate: exactly `.github/workflows/refactor-observability-errors.yml`.

Detached exact-SHA proof:
- worktree: `/home/sprite/work/proof-pre-final-b-7a22920`
- exact SHA `7a229202...` / tree `cb357507...`;
- canonical `make setup` reconstruction PASS;
- full focused tests/ratchets/removal/Wire reproducibility/protected-invariant/final-cleanliness proof PASS;
- `SLICE_B_DETACHED_EXACT_SHA_PROOF=PASS`.

### Safe publication

Remote branch:
- `refactor/pre-final-architecture-hardening-ca98ece`

Publication mapping:
- previous remote head: `373663d955d2fcffac1cfc96736ce694e8ce5b73`;
- remote source commit: `eb745852c7cd888d7298e3e6500d862be780db81`;
- remote source tree: `804cff3330a568346f2d6824dce755e6623f3f03`, exactly equal to local source tree;
- remote harness/head: `b7317728721821225c2596a8b0b5597adb08a6fc`;
- remote harness tree: `cb357507b528beef25a799aa346278c1d0ed7798`, exactly equal to local harness tree;
- parent chain is `373663d... -> eb745852... -> b7317728...`;
- a concurrent ref advance was detected before any ref update; it was classified as the exact tree-equivalent Slice B publication, so no overwrite/force update was performed.

`SLICE_B_SOURCE_MUTATION=CLOSED`
`SLICE_B_FOCUSED_PROOF=PASS`
`SLICE_B_ZERO_REMOVAL_PROOF=PASS`
`SLICE_B_DETACHED_EXACT_SHA_PROOF=PASS`
`SLICE_B_PUBLICATION=PASS`
`SLICE_B_HOSTED_PROOF=PENDING`

Current writer:
- `/home/sprite/work/arch-hardening-ca98ece`
- branch: `work/pre-final-hardening`
- HEAD: `7a229202f8a0b31c8470cb14faf3ce2bf8dcfb1e`
- clean

Preserved proof worktrees must not be removed/reset:
- `/home/sprite/work/proof-pre-final-a-f7f4c98`
- `/home/sprite/work/proof-pre-final-a-bec7db0`
- `/home/sprite/work/proof-pre-final-a2-bec7db0`
- `/home/sprite/work/proof-pre-final-b-7a22920`

## Remaining hardening order

1. Hosted exact-SHA proof for Slice B at remote head `b7317728...`.
2. Slice C — Payment outbox degradation/backoff/health/idempotency proof and bounded fixes.
3. Slice D — static enforcement for permanent canonical APIs where compatibility state permits.
4. Aggregate exact-SHA proof, canonical workflows, fresh-clone reconstruction, release build/verify/rollback proof.
5. Re-open Production Promotion only after hardening is PROVED/CLOSED.

`NEXT_GATE=SLICE_B_HOSTED_EXACT_SHA_PROOF`

`FINAL ACCEPTED=NO`
