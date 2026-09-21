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


## Slice B — Redis/tile-session ownership — PROVED/CLOSED

Local immutable evidence:
- source candidate: `cf5e62e5ece1add747c7be77cbd468128bf4ba27`
- source tree: `804cff3330a568346f2d6824dce755e6623f3f03`
- harness-only child: `7a229202f8a0b31c8470cb14faf3ce2bf8dcfb1e`
- harness tree: `cb357507b528beef25a799aa346278c1d0ed7798`
- detached exact-SHA proof: PASS.

Source outcomes:
- technical Redis lifecycle remains in `common/redis`;
- tile-session keyspace/TTL/persistence moved to narrow `common/tilesession.Store`;
- production key literal `ss:k:` has one owner;
- plaintext `sessionEncryptKey` logging retired and zero-ratcheted;
- User no longer keeps process-global `sessionKey`; each generated tile session is independent;
- TQD tile server now loads typed Redis config and owns exactly one explicit Redis open/close lifecycle;
- dead common Redis token helpers retired;
- Wire provider authority lives in `shared/code/script/injection_clients.yml` and is service-scoped;
- User Wire regeneration is hash-stable;
- protected paths and deploy hash unchanged.

Safe publication mapping:
- remote source commit: `eb745852c7cd888d7298e3e6500d862be780db81`
- remote source tree: `804cff3330a568346f2d6824dce755e6623f3f03` (identical to local source tree)
- remote harness/head: `b7317728721821225c2596a8b0b5597adb08a6fc`
- remote harness tree: `cb357507b528beef25a799aa346278c1d0ed7798` (identical to local harness tree)
- parent remains Slice A remote head `373663d955d2fcffac1cfc96736ce694e8ce5b73`.

Hosted exact-SHA:
- workflow run: `35619824787` / run number `104`
- head SHA: `b7317728721821225c2596a8b0b5597adb08a6fc`
- conclusion: `success`
- `inventory`: SUCCESS
- `common-contracts`: SUCCESS
- `boundary-contracts`: SUCCESS
- exact-SHA artifact: `observability-error-inventory-b7317728721821225c2596a8b0b5597adb08a6fc`
- artifact id: `10646784034`
- expired: false
- digest: `sha256:c8ab6d279734bbcb7f5fe8dba63910dc35c461b316592ec6821b22844d225753`.

`SLICE_B=PROVED/CLOSED`

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
- remote ref was re-read immediately before publication at exact prior head `373663d955d2fcffac1cfc96736ce694e8ce5b73`; it was then fast-forwarded to `b7317728721821225c2596a8b0b5597adb08a6fc` with `force=false`.


### Hosted exact-SHA closure

- workflow: `Refactor Observability and Error Contracts`;
- run number: `104`;
- run ID: `35619824787`;
- branch: `refactor/pre-final-architecture-hardening-ca98ece`;
- exact head SHA: `b7317728721821225c2596a8b0b5597adb08a6fc`;
- status/conclusion: `completed/success`;
- `inventory`: SUCCESS;
- `common-contracts`: SUCCESS, including Redis transport + tile-session semantic ownership;
- `boundary-contracts`: SUCCESS, including canonical `make setup`, Wire service-scoping, User Wire reproducibility, User/TQD tile-session boundaries, ownership-retirement proof, and all existing shared/Gateway/User/Payment/TQD regressions;
- exact-SHA artifact: `observability-error-inventory-b7317728721821225c2596a8b0b5597adb08a6fc`;
- artifact ID: `10646784034`;
- artifact expired: `false`;
- artifact digest: `sha256:c8ab6d279734bbcb7f5fe8dba63910dc35c461b316592ec6821b22844d225753`.

`SLICE_B_SOURCE_MUTATION=CLOSED`
`SLICE_B_FOCUSED_PROOF=PASS`
`SLICE_B_ZERO_REMOVAL_PROOF=PASS`
`SLICE_B_DETACHED_EXACT_SHA_PROOF=PASS`
`SLICE_B_PUBLICATION=PASS`
`SLICE_B_HOSTED_PROOF=PASS/CLOSED`
`SLICE_B=PROVED/CLOSED`

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

## Slice C — Payment outbox degradation/backoff/idempotency operational hardening

Read-only classification from immutable Slice B authority:
- delivery semantics are intentionally at-least-once, not exactly-once;
- producer outbox has durable unique `event_id`, `SKIP LOCKED` claim, lease recovery, claim version and attempt count;
- Notification consumer has durable Inbox dedupe by `event_id` and malformed-event DLQ;
- Rabbit publisher uses mandatory routing + publisher confirms; unknown confirm is retried and consumer Inbox absorbs duplicates;
- Payment currently uses one fixed retry delay for both durable message retry and transport reconnect even though these own different invariants;
- initial Rabbit connection failure can terminate `OutboxSupervisor` before establishment and therefore cancel the whole Payment process despite durable outbox safety;
- reconnect and message retry have no exponential cap/jitter policy;
- Payment compose health is TCP-port liveness only and the repository has no established gRPC-health/metrics consumer for this concern, so broker degradation must not redefine whole-service readiness without a consumer contract;
- current logs do not expose explicit degraded/recovered state, failure count, retry delay or sustained-failure escalation.

Authorized bounded direction:
- preserve at-least-once + Inbox idempotency;
- broker outage, including startup outage, is a recoverable actor degradation and must not terminate the Payment API process;
- transport reconnect timing is owned by `OutboxSupervisor`;
- durable message retry timing is owned by `internal/usecase/outbox`;
- use exponential bounded backoff with jitter for both concerns, but keep separate owners and settings;
- canonical publisher-unavailable classification belongs to the outbox port; Rabbit adapter wraps/projects it;
- transport unavailability releases the durable claim without stacking a second message-level delay; reconnect policy controls transport recovery;
- non-transport publish failures use attempt-aware durable retry backoff;
- structured logs project degraded/retrying/recovered state with retry duration and attempt/failure count; sustained capped failure escalates to ERROR;
- do not add producer terminal/DLQ state or a new health/metrics abstraction without a proven remediation/consumer contract;
- preserve transaction boundaries, unique event identity, Notification Inbox semantics, protobuf, Organization/Map no-touch and deploy-script invariants.

`SLICE_C_READ_ONLY_CLASSIFICATION=PASS/CLOSED`
`SLICE_C_SOURCE_MUTATION=AUTHORIZED`



## Slice C — Payment outbox operational hardening — PROVED/CLOSED

Local immutable evidence:
- source candidate: `70a59d3d2a60b16afb450d04118fcfcd9408798c`
- source tree: `a0a446f013199e4608718ffbb8e9972dd497ffb5`
- harness-only child: `0f40c937798f10f6e9c88f2f2f673ee84d26cce7`
- harness tree: `c6885fc2df1f7ac80f7e414c32832efc38d10bbf`
- harness delta from source candidate: exactly `.github/workflows/refactor-observability-errors.yml`.

Implemented bounded outcomes:
- at-least-once + durable Inbox semantics preserved;
- canonical `ErrPublisherUnavailable` ownership moved to `internal/usecase/outbox`; Rabbit adapter only wraps/projects it;
- durable event retry and broker reconnect now have separate typed policies/config ownership;
- both use bounded exponential delay with deterministic equal jitter and caps;
- transport unavailability releases the durable claim immediately so reconnect backoff is not stacked with message retry delay;
- non-transport publish failures use attempt-aware durable retry timing;
- broker unavailability at process startup is recoverable degradation and no longer terminates Payment API ownership;
- reconnect logs project retry delay, failure count, degraded duration, one capped ERROR escalation and recovered state;
- existing `PAYMENT_OUTBOX_RETRY_DELAY` remains a tested compatibility bridge while canonical base/max envs own the new policies;
- no new health/metrics abstraction, producer terminal/DLQ state, event schema or delivery guarantee was introduced.

Focused/local proof:
- full `payment-service: go test -count=1 ./...`: PASS;
- observability/error detector tests: 63 PASS;
- all ratchets PASS, including `go.payment_parallel_publisher_unavailable@payment-service = 0`;
- zero/removal proof PASS:
  - exactly one `ErrPublisherUnavailable` definition, in canonical outbox owner;
  - `rabbit.ErrPublisherUnavailable` refs = 0;
  - legacy `PublisherRetry` owner = 0;
  - legacy fixed `reconnectDelay time.Duration` owner = 0;
  - canonical retry/reconnect base/max envs present;
  - legacy retry env confined to compatibility owner/tests;
- protected paths unchanged: `shared/protobuf`, `organization-service`, `map-service`, `shared/code/deploy.sh`;
- deploy hash preserved: `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`.

Detached exact-SHA proof:
- worktree: `/home/sprite/work/proof-pre-final-c-0f40c93`;
- exact local harness SHA/tree: `0f40c937...` / `c6885fc2...`;
- canonical `make setup` reconstruction PASS;
- full Payment suite PASS;
- detector/ratchets PASS;
- zero/removal + protected invariants + final cleanliness PASS;
- `SLICE_C_DETACHED_EXACT_SHA_PROOF=PASS`.

Safe publication mapping:
- previous remote head: `b7317728721821225c2596a8b0b5597adb08a6fc`;
- remote source commit: `e32cc99da6135d236123f9cf3e2d871e319e2023`;
- remote source tree: `a0a446f013199e4608718ffbb8e9972dd497ffb5` (identical to local source tree);
- remote harness/head: `4c4db1c60f044dea57df61b423c9c5b7b9620341`;
- remote harness tree: `c6885fc2df1f7ac80f7e414c32832efc38d10bbf` (identical to local harness tree);
- parent chain: `b7317728... -> e32cc99d... -> 4c4db1c...`.

Hosted exact-SHA:
- workflow: `Refactor Observability and Error Contracts`;
- run number: `105`;
- run ID: `35623544994`;
- branch: `refactor/pre-final-architecture-hardening-ca98ece`;
- exact head SHA: `4c4db1c60f044dea57df61b423c9c5b7b9620341`;
- status/conclusion: `completed/success`;
- `inventory`: SUCCESS;
- `common-contracts`: SUCCESS;
- `boundary-contracts`: SUCCESS;
- Payment-specific hosted steps `Test Payment durable outbox recovery contract` and `Verify Payment outbox ownership retirement`: SUCCESS;
- exact-SHA artifact: `observability-error-inventory-4c4db1c60f044dea57df61b423c9c5b7b9620341`;
- artifact ID: `10649819592`;
- expired: false;
- digest: `sha256:f90154075def4828bc724da234d02c6e1e69f40860f6acfa47206a0880ec9f89`.

`SLICE_C_SOURCE_MUTATION=CLOSED`
`SLICE_C_FOCUSED_PROOF=PASS`
`SLICE_C_ZERO_REMOVAL_PROOF=PASS`
`SLICE_C_DETACHED_EXACT_SHA_PROOF=PASS`
`SLICE_C_PUBLICATION=PASS`
`SLICE_C_HOSTED_PROOF=PASS/CLOSED`
`SLICE_C=PROVED/CLOSED`

## Remaining hardening order

1. Slice D — static enforcement for permanent canonical APIs where compatibility state permits.
2. Aggregate exact-SHA proof, canonical workflows, fresh-clone reconstruction, release build/verify/rollback proof.
3. Re-open Production Promotion only after hardening is PROVED/CLOSED.

`SLICE_D_READ_ONLY_CLASSIFICATION=AUTHORIZED`
`SLICE_D_SOURCE_MUTATION=NOT_AUTHORIZED`

`NEXT_GATE=SLICE_D_STATIC_ENFORCEMENT_READ_ONLY_CLASSIFICATION`

`FINAL ACCEPTED=NO`
