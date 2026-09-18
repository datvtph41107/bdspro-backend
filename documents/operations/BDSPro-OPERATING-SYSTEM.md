# BDSPro Operating System — One Durable Way of Working

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: OPERATING MODEL V2 — ACTIVE

This document defines the single long-lived operating model for BDSPro work between the user, ChatGPT coordinator, recovery Sprite, local WSL execution, GitHub, CI, and the ChatGPT Library.

It is control-plane documentation. It does NOT replace source authority, exact-SHA proof, or the acceptance ledger.

## 1. Core objective

The workflow must remain correct even when:
- a chat becomes long or is interrupted/deloaded;
- a cloud execution tool times out;
- a local terminal closes;
- CI takes a long time;
- multiple read-only investigations are useful in parallel;
- a future chat has little conversational context.

The solution is:

ONE SOURCE TRUTH + ONE WRITER + MANY READERS/PROVERS + DURABLE CHECKPOINTS.

## 2. Authority hierarchy

For source and acceptance decisions:

1. immutable live Git/source at an exact SHA;
2. completed local exact-SHA proof;
3. completed hosted exact-SHA CI/proof;
4. recovery Sprite durable documents;
5. ChatGPT Library mirror;
6. local execution logs/evidence;
7. conversational Memory/context.

If durable prose disagrees with immutable live Git or completed exact-SHA evidence, live Git/evidence wins.

Memory is orientation only. Never rely on Memory as the only place that stores an important decision.

## 3. Current control-plane/source separation

Source authority:
- canonical acceptance branch: `final-acceptance/source-canonicalization`;
- active refactor branch: `refactor/canonical-observability-errors-a6d0722a`;
- current active source SHA at this operating-model checkpoint: `fca4682587d93cbc436f2466244ee8bd03b8b1b9`;
- production lineage remains separate.

Control-plane documentation:
- recovery Sprite: `mcp-123-bdspro-final-acceptance-recovery`;
- Sprite directory: `/home/sprite/work/acceptance`;
- Library mirror: `/BDSPro Final Acceptance`;
- Git ops branch: `ops/bdspro-operating-system`.

The Git ops branch is documentation/control-plane history only. Its HEAD must never be confused with active source authority.

## 4. Protected invariants

Unless a later explicit acceptance decision changes them:
- `shared/protobuf/**` NO-TOUCH;
- `organization-service/**` NO-TOUCH;
- `map-service/**` NO-TOUCH;
- preserve `make wire`;
- preserve `make buf`;
- preserve byte-identical `shared/code/deploy.sh`;
- never force-push active acceptance/refactor branches;
- never call FINAL ACCEPTED until every required gate passes.

## 5. Roles

### Coordinator — ChatGPT

Responsibilities:
- recover durable state;
- reconcile live Git/CI;
- reason about architecture and source reality;
- classify REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST;
- choose exactly one bounded mutation slice;
- define exact user execution tasks;
- review evidence;
- decide gate transitions;
- synchronize durable state.

### User local execution lane

Responsibilities:
- run heavy/deterministic commands on WSL;
- run long tests, builds, generation and integration/runtime proof;
- keep execution alive in tmux;
- return exact outputs/exit status to the coordinator;
- never invent or broaden source changes outside the currently authorized writer task.

### Read-only workers

May run concurrently on one immutable SHA:
- source inventories;
- rg/search;
- dependency/consumer tracing;
- audit reports;
- test discovery;
- CI/source comparison.

They must not mutate source.

### Source writer

There is exactly ONE source writer for an authorized bounded slice.

No second service is mutated in parallel.

### Proof workers

After a candidate exact SHA exists, independent proof lanes may run concurrently:
- service/unit tests;
- audit/ratchets;
- common regressions;
- protected-path checks;
- runtime smoke;
- hosted CI jobs.

All proof workers must read the same candidate SHA.

## 6. Parallelism rule

Parallelize:
- discovery;
- inventory;
- search;
- read-only source tracing;
- deterministic proof of one immutable candidate;
- hosted CI jobs.

Serialize:
- architecture decision;
- source mutation;
- deletion/retirement decision;
- ratchet creation;
- branch publication;
- final gate transition;
- durable-state closure.

Mnemonic:

PARALLEL READERS + ONE WRITER + PARALLEL PROVERS + ONE CLOSER.

## 7. Per-slice lifecycle

Every source slice follows this order:

1. RECOVER
2. RECONCILE LIVE GIT
3. READ-ONLY INVENTORY
4. OWNER / CONSUMER / VALUE DECISION
5. AUTHORIZE ONE BOUNDED SLICE
6. ONE WRITER MUTATION
7. ZERO / REMOVAL PROOF where applicable
8. DELETE / RETIRE competing path where justified
9. RATCHET only after zero proof
10. COMMIT CANDIDATE
11. LOCAL EXACT-SHA PROOF
12. SAFE NON-FORCE PUBLICATION
13. HOSTED EXACT-SHA PROOF
14. DURABLE CHECKPOINT SYNCHRONIZATION
15. NEXT SLICE

Never reorder proof -> retirement -> ratchet merely to make a counter green.

## 8. Local filesystem model

The existing local repository may be used as the anchor clone only after verification.

Preferred long-lived layout:

```text
~/projects/bdspro-canonical-bootstrap-20260908/   # anchor clone; fetch/worktree manager
~/bdspro-worktrees/
  reader/                                        # detached immutable inventory
  writer/                                        # the only mutable candidate
  proof/                                         # detached candidate proof when needed
~/bdspro-ops/
  logs/                                          # command logs
  proofs/                                        # summarized local proof artifacts
  state/                                         # convenience snapshots, NOT authority
  run.sh                                         # command wrapper
```

The anchor clone is not a coding workspace once worktrees are established.

## 9. Worktree roles

### reader

Detached at the current active source SHA.

Allowed:
- grep/search;
- audit;
- read-only tests;
- inventory.

Not allowed:
- source edits;
- commits;
- pushes.

### writer

Local branch created from the exact active SHA for ONE authorized slice.

Allowed:
- only files inside the approved scope;
- commit after proof prerequisites.

Not allowed:
- unrelated cleanup;
- second service migration;
- push before coordinator reconciliation.

### proof

Detached at the exact committed candidate SHA.

Used for:
- parallel proof without touching writer state;
- runtime/integration tests;
- heavy verification.

## 10. tmux operating model

Use one persistent session:

`tmux new -s bdspro`

Recommended windows:
- 0-control
- 1-inventory-a
- 2-inventory-b
- 3-proof
- 4-runtime
- 5-git

Useful keys:
- Ctrl+b c: create window;
- Ctrl+b n / p: next / previous;
- Ctrl+b ,: rename current window;
- Ctrl+b d: detach;
- `tmux attach -t bdspro`: return later.

tmux protects local processes from browser/chat/terminal interruption.

## 11. User-lane task contract

Every coordinator request for local execution must contain:

```text
[USER-LANE]
MODE: READ_ONLY | WRITER | PROOF | RUNTIME
AUTHORITY_SHA: <sha>
WORKDIR: <absolute/~/path>
PURPOSE: <why this command exists>
COMMAND:
  <copy-paste command>
EXPECTED:
  <what proves success>
STOP IF:
  <conditions that require returning to coordinator>
RETURN:
  <exact output/exit code/files>
```

The user should not infer the next mutation from a failed command. Return evidence first.

## 12. Evidence wrapper

Long-running local commands should be run through `~/bdspro-ops/run.sh`.

The wrapper records:
- task name;
- start/end time;
- repo;
- branch;
- HEAD before/after;
- stdout/stderr;
- exit code.

Logs are evidence, not authority.

## 13. Interruption / deload recovery

If ChatGPT/browser/tooling is interrupted:

1. do not guess where work stopped;
2. preserve any running local process in tmux;
3. preserve local logs;
4. start a new chat with the canonical continuation prompt;
5. coordinator reads durable Sprite/Library;
6. coordinator reconciles live active Git SHA and completed CI;
7. only then resume mutation.

A chat is disposable. Git exact SHA + completed proof + durable checkpoint are not.

## 14. Checkpoint policy

Create/synchronize a durable checkpoint after:
- a material architecture/operating decision;
- a bounded slice is published and hosted proof completes;
- an important failure changes next actions;
- the operating model itself changes.

Do not create false CLOSED state while hosted proof is still pending.

## 15. Git safety

Never run without explicit reason:
- `git reset --hard`;
- `git clean -fd`;
- force push;
- mass checkout/restore over unknown local changes.

Before mutation:
- `git status --short --branch`;
- verify exact authority SHA;
- verify correct worktree role.

Before publication:
- verify parent SHA;
- verify exact changed paths;
- `git diff --check`;
- protected-path diff;
- deploy hash;
- required local proof.

## 16. Current runtime/error phase

At this checkpoint:
- Runtime R1-R4 are closed.
- R5 canonical logger adoption is active service-by-service.
- R5 Payment slice is closed at `fca4682587d93cbc436f2466244ee8bd03b8b1b9`.
- Error/Response FINAL design is closed but implementation remains pending.
- R5 must not be bypassed merely because Error/Response design already exists.

## 17. Working philosophy

`ONE TRUTH PER CONCERN -> ONE OWNER -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`

The goal is not maximum activity. The goal is maximum trustworthy progress per unit of execution.
