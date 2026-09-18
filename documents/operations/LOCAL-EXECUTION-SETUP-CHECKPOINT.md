# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 2B COMPLETE / PHASE 3 OPERATIONS WRAPPER + TMUX STANDARDIZATION NEXT

## Authority

Live GitHub active refactor remains:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- SHA: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- GitHub comparison: identical, ahead 0 / behind 0.

## Phase 2A archive — VERIFIED

Parked unpublished Auth zero-ratchet work is preserved outside the repository at:

`/home/bop/bdspro-ops/state/anchor-auth-ratchet-20260918-112954`

Archive includes exact patch, authority comparison patch, both modified files, metadata, status, pycache inventory and SHA256 checksums.

The parked Auth ratchet remains PARKED and is not active source.

## Phase 2B local cleanup / synchronization — COMPLETE

Anchor:

`/home/bop/projects/bdspro-canonical-bootstrap-20260908`

Final state:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- HEAD: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- working tree clean
- tracking remote active branch without ahead/behind divergence.

Cleanup performed safely:

- restored ONLY the two archived tracked audit files to the old local HEAD before fast-forward;
- removed ONLY generated `shared/code/development/__pycache__`;
- verified clean anchor;
- fetched and verified exact remote authority;
- fast-forwarded with `git merge --ff-only`;
- did not use reset-hard, clean -fd, force-push or plain pull.

Reader:

`/home/bop/bdspro-worktrees/reader`

- detached at exact authority SHA;
- clean;
- canonical READ_ONLY / inventory worktree.

Writer:

`/home/bop/bdspro-worktrees/writer`

- branch `local/r5-next-candidate`;
- exact authority SHA;
- clean;
- no commits/diff;
- remains PARKED until a bounded R5 slice is authorized.

Stale temporary worktree registry metadata was pruned successfully. Only anchor, reader and writer remain registered.

## Benign post-block command error

After the successful Phase 2B subshell finished, the user manually ran:

`git merge --ff-only "origin/$BRANCH"`

outside the subshell.

Because `BRANCH` was defined only inside `( ... )`, it no longer existed in the parent shell. The command therefore became effectively:

`git merge --ff-only origin/`

and failed with:

`merge: origin/ - not something we can merge`

This did NOT change Git/source state and does NOT invalidate Phase 2B. Do not rerun the merge; the anchor is already at the exact authority SHA.

## Phase 3 — next authorized local action

Standardize the user's execution control plane:

1. create/verify `~/bdspro-ops/logs`, `proofs`, `state`;
2. replace `~/bdspro-ops/run.sh` with the canonical evidence wrapper;
3. syntax-test the wrapper;
4. execute one harmless READ_ONLY command through the wrapper;
5. verify the log records repo, branch, head_before, head_after and exit code;
6. inspect/reuse the existing `tmux` session rather than killing it;
7. create/rename stable tmux windows for control, inventory, proof, runtime and git;
8. do not start the next R5 service source mutation yet.

No active source mutation is authorized by this checkpoint.
