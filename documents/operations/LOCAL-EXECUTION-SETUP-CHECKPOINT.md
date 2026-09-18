# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 1 RECONCILED / PHASE 1B INSPECTION NEXT

This checkpoint records the user's WSL/local execution reality discovered while adopting Operating Model V2.

## Authority reconciliation

Live GitHub active refactor remains:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- SHA: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`

GitHub compare against that SHA is identical: ahead 0 / behind 0.

## User anchor clone

Path:

`/home/bop/projects/bdspro-canonical-bootstrap-20260908`

Repository root resolves to the same path.

Remote:

`git@github.com:datvtph41107/bdspro-backend.git`

Local branch:

`refactor/canonical-observability-errors-a6d0722a`

Local HEAD:

`e857fd4d0e34d0f61e5beb33b1e849590e78849b`

Relationship to remote active branch:

- local anchor is behind by 5 commits;
- do NOT pull/reset/clean yet.

Local dirty state:

- modified: `shared/code/development/audit-observability-errors.py`
- modified: `shared/code/development/test_audit_observability_errors.py`
- untracked: `shared/code/development/__pycache__/`

These local changes are not yet classified as valuable/stale/generated. Preserve them until inspected.

## Existing worktrees discovered

- anchor: `/home/bop/projects/bdspro-canonical-bootstrap-20260908` @ `e857fd4`
- reader: `/home/bop/bdspro-worktrees/reader` @ `fca4682` detached
- writer: `/home/bop/bdspro-worktrees/writer` @ `fca4682` branch `local/r5-next-candidate`
- several `/tmp/bdspro-chat-*` worktrees @ `9d3d3ce`, marked prunable

The reader/writer worktrees appear to be based on the correct authority SHA, but their cleanliness and exact role state have not yet been verified.

## Safety decision

Until Phase 1B inspection completes:

- do not run `git pull` in the anchor;
- do not run `git reset --hard`;
- do not run `git clean -fd`;
- do not delete the two modified audit files;
- do not delete the untracked `__pycache__` yet;
- do not reuse/push the writer branch;
- do not prune worktrees yet.

## Next authorized local step — Phase 1B

Read-only inspection only:

1. inspect anchor diff for the two modified audit files;
2. inspect anchor untracked `__pycache__` contents;
3. inspect reader HEAD/status/branch;
4. inspect writer HEAD/status/branch and diff from authority;
5. inspect worktree lock/prunable metadata.

After this evidence is returned, the coordinator will decide:
- whether anchor changes should be preserved, archived or discarded;
- whether reader can be reused as-is;
- whether writer can be reused or should be recreated;
- when it is safe to prune stale `/tmp` worktrees;
- when to install the canonical `run.sh` wrapper.

No source mutation is authorized by this checkpoint.
