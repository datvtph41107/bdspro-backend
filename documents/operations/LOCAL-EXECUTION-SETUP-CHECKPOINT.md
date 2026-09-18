# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 3B COMPLETE / PHASE 3C END-TO-END USER-LANE NEXT

## Authority

Live GitHub active refactor remains:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- SHA: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- GitHub comparison: identical, ahead 0 / behind 0.

## Local Git layout

Anchor:
- `/home/bop/projects/bdspro-canonical-bootstrap-20260908`
- active refactor branch
- exact authority SHA
- clean

Reader:
- `/home/bop/bdspro-worktrees/reader`
- detached exact authority SHA
- clean
- canonical READ_ONLY/inventory worktree

Writer:
- `/home/bop/bdspro-worktrees/writer`
- branch `local/r5-next-candidate`
- exact authority SHA
- clean
- PARKED until a bounded R5 source slice is authorized

Parked unpublished Auth ratchet remains archived at:
`/home/bop/bdspro-ops/state/anchor-auth-ratchet-20260918-112954`

## Operations control plane

Canonical runner:
`/home/bop/bdspro-ops/run.sh`

Runner proof:
- syntax PASS
- wrapped reader smoke exit 0
- evidence log persisted
- exact head_before/head_after preserved
- reader remained clean

Current runner SHA256:
`cb5218995c8e5595286d70300411e234f3f563e41fab4d377c6b24440665ce16`

Canonical tmux layout:
`/home/bop/bdspro-ops/tmux-layout.sh`

Current layout SHA256:
`aceb22ffb70b07699c28863a9d501ec7cc9785203ce162b23389bbe099b97eb0`

Pre-fix layout backup:
`/home/bop/bdspro-ops/state/tmux-layout.sh.before-pane-fix-20260918-115031`

## Phase 3B tmux standardization — COMPLETE

Reusable session:
`bdspro`

Stable role windows:
- `0 control` -> anchor
- `1 inventory` -> reader
- `2 proof` -> operations proofs directory
- `3 runtime` -> operations directory
- `4 git` -> anchor

All windows:
- one pane each;
- deterministic working directory PASS;
- `automatic-rename=off` at window scope PASS;
- `allow-rename=off` at pane scope PASS;
- source remained clean/exact authority after layout application.

The previous option-scope defect is CLOSED.

## Phase 3C — next authorized local action

Prove the full user-lane operating loop without source mutation:

1. attach to existing `bdspro` session;
2. practice role navigation;
3. run one real R5 READ_ONLY inventory task from the `inventory` window through `~/bdspro-ops/run.sh`;
4. persist the inventory report under `~/bdspro-ops/proofs`;
5. verify reader SHA/cleanliness before and after;
6. return task output, exit code, log path and proof artifact path to the coordinator;
7. coordinator uses that evidence to choose the next bounded R5 service.

No source mutation, writer use, publication, or ratchet creation is authorized in Phase 3C.
