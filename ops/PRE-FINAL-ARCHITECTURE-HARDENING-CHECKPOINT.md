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
- branch: `refactor/pre-final-architecture-hardening-ca98ece`;
- branch was absent before publication; publication was create-only/non-force.

Proof:
- focused common middleware/errors tests PASS;
- focused Payment interceptor/handler tests PASS;
- Payment compile gate PASS;
- detector tests: 57 PASS;
- all ratchets PASS, including `go.direct_grpc_error_projection@payment-service=0`;
- Payment direct `ToGRPC(` = 0;
- Payment-local recovery retired;
- detached exact-SHA reconstruction/proof PASS;
- protected paths unchanged and `shared/code/deploy.sh` hash preserved.

Hosted exact-SHA:
- workflow: `Refactor Observability and Error Contracts`
- run number: `103`
- run ID: `35612933219`
- head SHA: `373663d955d2fcffac1cfc96736ce694e8ce5b73`
- status/conclusion: `completed/success`
- `inventory`: SUCCESS
- `common-contracts`: SUCCESS
- `boundary-contracts`: SUCCESS
- exact-SHA artifact: `observability-error-inventory-373663d955d2fcffac1cfc96736ce694e8ce5b73`
- artifact ID: `10644069205`
- expired: `false`
- digest: `sha256:506b53ef09cfd654b856b8b6902563ae254cd0951cb027543dea6d0ac62bdafa`

`SLICE_A=PROVED/CLOSED`

## Slice B — Redis transport/key-policy/secret ownership

Read-only classification already proven before mutation:
- `common/redis.RedisService` currently mixes technical transport/lifecycle with feature semantics;
- tile-session keyspace `ss:k:`, TTL and AES session-key storage are owned inside `common/redis` rather than a semantic concern owner;
- `sessionEncryptKey` is logged in plaintext in both `shared/common/redis/tile_session.go` and User `GenTileSessionToken`;
- logging redaction is key-based and cannot redact secrets interpolated into message text;
- live tile-session consumers are User session generation and TQD tile encryption;
- common Redis token helper methods duplicate User-owned token-cache behavior and have no live external caller;
- TQD standalone tile server still opens Redis through legacy `NewRedisService()`, so process-resource ownership must be traced within this slice.

Authorized design constraints:
- `common/redis` remains technical Redis transport/lifecycle owner, not a generic business manager;
- introduce a narrow tile-session semantic owner only for proven tile-session invariants;
- key prefix/TTL/secret representation have one owner;
- no plaintext secret logging;
- consumers receive semantic capability rather than redefining key/TTL;
- retire proven-dead duplicate common Redis business helpers rather than preserve speculative APIs;
- preserve wire behavior and protobuf schema; `shared/protobuf` remains no-touch.

Current writer:
- `/home/sprite/work/arch-hardening-ca98ece`
- branch: `work/pre-final-hardening`
- HEAD: `bec7db04aa2502269eed7018d6c4e3898125c979`
- clean

`SLICE_A_HOSTED_PROOF=PASS`
`SLICE_A_PUBLICATION=PASS`
`SLICE_B_SOURCE_MUTATION=AUTHORIZED`

`NEXT_GATE=SLICE_B_BOUNDED_SOURCE_MUTATION`

`FINAL ACCEPTED=NO`
