# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 2A ARCHIVE VERIFIED / PHASE 2B CLEAN + FAST-FORWARD NEXT

## Authority

Live GitHub active refactor remains:

- branch: `refactor/canonical-observability-errors-a6d0722a`
- SHA: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- GitHub comparison: identical, ahead 0 / behind 0.

## Anchor clone

Path:

`/home/bop/projects/bdspro-canonical-bootstrap-20260908`

Local HEAD before cleanup:

`e857fd4d0e34d0f61e5beb33b1e849590e78849b`

Local branch:

`refactor/canonical-observability-errors-a6d0722a`

Relationship to authority:

- behind by 5 commits;
- two modified tracked audit files;
- one generated Python bytecode artifact.

## Parked Auth ratchet archive — VERIFIED

Archive directory:

`/home/bop/bdspro-ops/state/anchor-auth-ratchet-20260918-112954`

The repository remained unchanged after archive creation.

Archived evidence:

- `local-working-tree.patch` — 88 lines, exact local working diff vs local HEAD;
- `vs-authority.patch` — 133 lines, local working state vs current authority;
- exact copies of both modified audit files;
- `git-status.txt`;
- `metadata.txt`;
- `pycache-list.txt`;
- `SHA256SUMS`.

Recorded archive checksums:

- audit file: `4cc4b368f17dec31f1a9ed12a30aa41e54c49387948cd43e96377a1fe98c4ed6`
- audit test file: `1889050770094eb00c5f917a72b3636ae07cb9f890a3f9163e7c95b58ffd51dc`
- git-status: `3191540ad924e503e38f410a90690221083b7cdd47345ecce1e604cafb8254f3`
- local patch: `a7a0dfa47c13f70c6f2543663961f9b75a437c44913d664d39afd0756f15af58`
- metadata: `a2ad1b1912fd6d905ed5c6816b1e7e324b53f6cb4f44e1581b8cce9baf0137d7`
- pycache list: `b084430c24fbe9ba673c2e674c38f08e4ceb200c39bc0bfca0237603dc526f4d`
- authority comparison patch: `ab2afc1e43dbc2192f3e5e8ecf6718a0fdd5c03480088b977e2860f680d5b773`

Classification:

- tracked local diff = PARKED / UNPUBLISHED AUTH ZERO-RATCHET CANDIDATE;
- do not publish during local setup;
- generated `__pycache__` is not source authority.

## Existing worktrees

Reader:

- `/home/bop/bdspro-worktrees/reader`
- detached at `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- clean
- REUSE.

Writer:

- `/home/bop/bdspro-worktrees/writer`
- branch `local/r5-next-candidate`
- HEAD `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- no commits/diff
- keep PARKED until next bounded slice is selected.

Stale temp worktree metadata:

- four `/tmp/bdspro-chat-*` entries are prunable;
- no worktree source directories remain.

## Phase 2B — authorized local action

Now that parked local work is archived and checksummed, it is authorized to:

1. restore ONLY the two tracked audit files to the local anchor HEAD;
2. remove ONLY the generated audit `__pycache__` directory;
3. verify anchor becomes clean but remains behind 5;
4. fetch and verify remote authority still equals `fca46825...`;
5. fast-forward anchor using `git merge --ff-only origin/refactor/canonical-observability-errors-a6d0722a`;
6. verify exact authority SHA and clean working tree;
7. prune stale worktree registry metadata with `git worktree prune --verbose`;
8. verify reader/writer registrations remain intact.

Do NOT use:
- `git reset --hard`;
- `git clean -fd`;
- force push;
- plain `git pull`.

No active source mutation is authorized by this checkpoint.
