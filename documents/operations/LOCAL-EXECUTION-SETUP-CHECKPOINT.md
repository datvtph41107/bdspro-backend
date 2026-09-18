# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 3B PARTIAL PASS / TMUX OPTION-SCOPE FIX NEXT

## Authority

Live GitHub active refactor remains:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- SHA: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- GitHub comparison: identical, ahead 0 / behind 0.

## Phase 3A

Operations wrapper remains PROVED locally:

- canonical `~/bdspro-ops/run.sh` installed;
- syntax PASS;
- first wrapped reader smoke task exit 0;
- evidence fields and exact SHA preservation PASS.

## Phase 3B current tmux reality

The existing `bdspro` session was reused successfully.

Window structure created successfully:

- `0 control`
- `1 inventory`
- `2 proof`
- `3 runtime`
- `4 git`

All five windows have one pane.

Working-directory verification PASS:

- control -> `/home/bop/projects/bdspro-canonical-bootstrap-20260908`
- inventory -> `/home/bop/bdspro-worktrees/reader`
- proof -> `/home/bop/bdspro-ops/proofs`
- runtime -> `/home/bop/bdspro-ops`
- git -> `/home/bop/projects/bdspro-canonical-bootstrap-20260908`

## Verification defect discovered

Phase 3B stopped at the stable-window-name verification:

`control automatic-rename=off allow-rename=`

This is a control-plane script verification defect, not a source or tmux-layout failure.

Reason:

- `automatic-rename` is a window option and was correctly set/read as `off`;
- `allow-rename` is pane-scoped behavior and should be set/read through pane options;
- the canonical layout script incorrectly used window-option commands for `allow-rename`;
- with `set -e`, the empty value caused the equality test to stop the block.

No source mutation occurred. The active Git authority remains unchanged.

## Next authorized action

Repair only the local control-plane script:

1. preserve the current `~/bdspro-ops/tmux-layout.sh` as a pre-fix backup;
2. change `allow-rename` handling to pane scope;
3. syntax-check;
4. reapply the existing idempotent layout;
5. verify automatic-rename at window scope;
6. verify allow-rename at pane scope;
7. verify all paths and Git worktrees remain clean/exact authority.

Do not recreate the session and do not start source work yet.
