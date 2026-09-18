# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 3A COMPLETE / PHASE 3B TMUX STANDARDIZATION NEXT

## Authority

Live GitHub active refactor remains:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- SHA: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- GitHub comparison: identical, ahead 0 / behind 0.

## Local Git layout

Anchor:
- `/home/bop/projects/bdspro-canonical-bootstrap-20260908`
- branch `refactor/canonical-observability-errors-a6d0722a`
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
- PARKED until next bounded source slice is authorized

Parked Auth zero-ratchet archive remains preserved under:
`/home/bop/bdspro-ops/state/anchor-auth-ratchet-20260918-112954`

## Phase 3A operations wrapper — COMPLETE

Operations directories exist:

- `/home/bop/bdspro-ops/logs`
- `/home/bop/bdspro-ops/proofs`
- `/home/bop/bdspro-ops/state`

Previous runner was preserved at:

`/home/bop/bdspro-ops/state/run.sh.before-v2-20260918-113731`

with SHA256:

`3b2ecb0c98fe2db7ec13ca207d9e89a7c2c30fa7efac1e8b7e2f563408b24063`

Canonical runner installed:

`/home/bop/bdspro-ops/run.sh`

Proof:

- `bash -n` PASS;
- reader authority check PASS;
- first wrapped task `phase3-reader-smoke` exit 0;
- evidence log:
  `/home/bop/bdspro-ops/logs/20260918-113731-phase3-reader-smoke.log`;
- log contains task/start/repo/branch/head_before/exit/head_after;
- head_before/head_after both exact authority SHA;
- reader remained clean.

## tmux reality

The shell used for Phase 3A was outside tmux.

Existing tmux state:

- session: `bdspro`
- one window only
- window 0 name: `bash`
- one pane
- session remains reusable; do not kill/recreate it.

## Phase 3B — next authorized local action

Standardize the existing `bdspro` tmux session into durable role windows without starting source work:

- `0-control` — coordination/status; anchor/control-plane only
- `1-inventory` — read-only work in reader
- `2-proof` — proof/evidence shell; no candidate source until an exact candidate SHA exists
- `3-runtime` — runtime/integration shell; no service start yet
- `4-git` — Git/worktree/checkpoint inspection from anchor

Requirements:

1. reuse existing session;
2. rename existing window rather than killing it;
3. create only missing windows;
4. establish deterministic working directories;
5. verify windows/panes;
6. do not mutate active source;
7. do not use writer yet.

No active source mutation is authorized by this checkpoint.
