# BDSPro Operating System — One Durable Way of Working

Updated: 2026-09-19 Asia/Bangkok
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
- current remote active source SHA at this operating-model checkpoint: `c6a9b121946a360d22759bde6a708bc6a35223c2`; Notification candidate `381c29365a00f35737bb0b6078cf0973198dc7c2` remains local-only pending corrected detached proof;
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
- Error/Response FINAL design is refined with the proven `reason` requirement; implementation remains pending until R5 is deliberately stable.
- R5 must not be bypassed merely because Error/Response design already exists.

## 17. Working philosophy

`ONE TRUTH PER CONCERN -> ONE OWNER -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`

The goal is not maximum activity. The goal is maximum trustworthy progress per unit of execution.


## 18. External AI code-review lane

OpenCodeReview (Alibaba) is approved for a bounded pilot as an OPTIONAL REVIEW/PROOF LANE.

It is not:
- source authority;
- a replacement for deterministic tests/audit/ratchets;
- permission to auto-fix source;
- a reason to broaden a bounded slice;
- a final acceptance gate by itself.

Intended placement:

```text
ONE WRITER
  -> candidate commit SHA
  -> deterministic local proof
  -> OCR review on exact base SHA -> candidate SHA
  -> coordinator triage of OCR findings
  -> publication only after required deterministic gates
```

Pilot rules:
- install a pinned OCR version, never floating latest in CI;
- start local-only/read-only;
- use exact validated 40-hex Git SHAs for `--from`/`--to`/`--commit`;
- output review artifacts under `~/bdspro-ops/proofs/ocr`;
- do not enable automatic fix;
- do not post PR comments or provide write GitHub tokens during the pilot;
- do not copy a `pull_request_target` workflow blindly;
- keep credentials outside chat and repository;
- review findings are advisory evidence requiring coordinator/source verification.

Only after the local pilot proves useful, reproducible and acceptably low-noise may OCR become a hosted CI review lane.


## 19. Naming, ownership and public-surface doctrine

Durable doctrine:
`ARCHITECTURE-NAMING-OWNERSHIP-MINDSET.md`.

Extend the working philosophy to:

`ONE TRUTH PER CONCERN -> ONE OWNER -> ONE PUBLIC NAME -> ONE OBVIOUS ENTRY POINT -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

Rules:
- do not use one catch-all owner such as `QuotaManager` for pricing, reservation and durable usage when they preserve different invariants;
- split by invariant/reason-to-change, not mechanically by method count;
- prefer role/capability names such as `QuotaPricer`, `QuotaReserver`, `QuotaLedger`;
- scrutinize generic `Manager`, `Helper`, `Util`, `Common`, `Service` names when they obscure ownership;
- consumer-facing interfaces should contain only the behavior that consumer actually needs;
- do not expose two long-lived public commands/APIs for the same caller intent merely because internal implementations differ;
- preserve separate internal mechanics where the underlying truths differ.

Runtime/logging implication:
- raw process output and structured evidence remain different internal sources/projections;
- target daily public vocabulary converges on one `make logs` entry point with filters/options;
- R4 `make log-query` is still closed/proved implementation history and must not be destructively removed or renamed without a bounded compatibility/docs/tests/proof slice.

Error/Response implication:
- one application failure fact is represented by protocol code + stable machine reason + public-safe message;
- response, logging, metrics and traces project the same code/reason;
- no second mandatory numeric business code;
- no generic ErrorManager/failure logging framework;
- ordinary technical Go errors remain ordinary wrapped Go errors;
- target implementation naming should make the semantic distinction obvious, conceptually `failure.New(code, reason, safeMessage)`, subject to live compatibility/source proof before mutation.

## 20. Continuation prompt contract

The short cross-chat bootstrap may remain intentionally small. A new chat does not need the entire project history pasted into the prompt if it is instructed to:
1. read `CONTINUATION-PROMPT.md` from the recovery Sprite;
2. fall back to the mirrored Library file when Sprite is unavailable;
3. follow that document exactly;
4. reconcile live Git/source, local worktree/candidate state and completed exact-SHA proof before mutation.

The durable continuation file, not the short pasted bootstrap text, owns the evolving recovery order and current checkpoint.

A chat interruption must not cause:
- source reset/clean/rebase;
- candidate recreation;
- reopening closed architecture merely because conversational context is missing;
- reliance on Memory as source authority.

A chat is disposable; exact Git state + proof + durable checkpoint are the continuity mechanism.
