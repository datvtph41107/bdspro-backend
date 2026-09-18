# BDSPro Local Execution Setup Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PHASE 3B COMPLETE / PHASE 3C FIRST ATTEMPT RECORDED / RETRY NEXT

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


## Phase 3C first attempt — RECORDED, NOT CLOSED

Two interface/tooling failures occurred without source mutation.

### Metadata pasted as shell

The user pasted the descriptive `[USER-LANE]` metadata block into Bash and received `command not found` for labels such as `ID:`, `MODE:`, and Vietnamese purpose text.

Classification:
- presentation/protocol usability defect;
- not a Git/source failure;
- USER-LANE protocol updated so metadata is explicitly DO NOT COPY and only fenced Bash is executable.

### Audit external-output path failure

The real R5 inventory ran on exact reader SHA and measured:
- findings 2783;
- debt 721;
- `go.legacy_std_log` debt 617;
- `go.third_party_logger` debt 99.

It then exited 1 at the script's final reporting step:

`args.output.relative_to(ROOT)`

because Phase 3C passed an output path under `~/bdspro-ops/proofs`, which lies outside repository ROOT.

The runner correctly recorded:
- exit=1;
- head_before=head_after=`fca4682587d93cbc436f2466244ee8bd03b8b1b9`;
- log path `/home/bop/bdspro-ops/logs/20260918-115445-r5-next-service-inventory.log`.

Classification:
- audit CLI output-path UX limitation;
- inventory scan itself completed and printed counts;
- Phase 3C is NOT closed because artifact/evidence verification did not finish;
- do not patch source merely to finish local setup.

Retry strategy:
- run audit using its canonical repository-local ignored default output paths;
- copy completed TSV/JSON artifacts to `~/bdspro-ops/proofs` only after command exit 0;
- build the human report from the copied evidence;
- verify reader remains clean/exact SHA.

## OpenCodeReview evaluation state

Alibaba OpenCodeReview is approved for evaluation only, not installation/source integration yet.

Decision:
- potentially strong fit as a parallel review/proof lane;
- not source authority and not a replacement for exact-SHA deterministic proof;
- local pinned-version pilot comes after Phase 3C retry closes;
- no GitHub Action/write-token integration during the first pilot.
