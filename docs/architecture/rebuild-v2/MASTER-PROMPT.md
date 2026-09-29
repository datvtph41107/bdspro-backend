# MASTER PROMPT — BDSPro Architecture Rebuild V2

This is the single orchestration prompt for the entire Architecture Rebuild V2
program, from recovery through final convergence.

There must be no second competing bootstrap prompt. A new chat/account/session
only needs to be told to read this file and follow it exactly.

---

## 0. Mission

Continue BDSPro Architecture Rebuild V2 from durable state, not conversational
memory.

Repository:

`datvtph41107/bdspro-backend`

Working branch:

`architecture/rebuild-v2`

The project is both:

- a real BDSPro architecture reconstruction/rebuild program;
- a practical Go/backend training corpus;
- a durable engineering decision system;
- a controlled laboratory for source, data, runtime and operational design.

Do not turn it into a tutorial toy, a pattern-copy exercise or a framework-first
rewrite.

---

## 1. Recovery is mandatory before reasoning or mutation

At the start of every fresh context, account, machine or resumed session:

1. resolve the live `architecture/rebuild-v2` branch HEAD;
2. read `ARCHITECTURE-REBUILD-V2.md`;
3. read, in order:
   - `docs/architecture/rebuild-v2/OPERATING-MODEL.md`
   - `docs/architecture/rebuild-v2/HISTORY.md`
   - `docs/architecture/rebuild-v2/MEMORY.md`
   - `docs/architecture/rebuild-v2/DECISION-REGISTER.md`
   - `docs/architecture/rebuild-v2/CHECKPOINT.md`
   - `docs/architecture/rebuild-v2/WORKLOG.md`
   - `docs/architecture/rebuild-v2/RECOVERY-SCENARIOS.md`
4. inspect live source relevant to the current gate;
5. inspect the current worktree when available;
6. inspect completed local/hosted proof for the exact relevant SHA;
7. reconcile durable prose against immutable Git/source/proof;
8. preserve interrupted or unexplained work;
9. determine the current authorized gate and exact next action.

Do not ask the user to reconstruct a previous conversation if the repository is
sufficient.

Do not use chat memory to override Git.

Authority order:

```text
1. immutable live Git/source at exact SHA
2. completed exact-SHA local proof
3. completed exact-SHA hosted CI/proof
4. CHECKPOINT.md
5. DECISION-REGISTER.md
6. WORKLOG.md / MEMORY.md / HISTORY.md
7. external research
8. conversational memory
```

If sources disagree, explicitly record the contradiction and reconcile using the
higher authority.

---

## 1.1 Preferred local workspace

When the user's WSL filesystem is available, the preferred working copy is:

`~/projects/bdspro-canonical-bootstrap-20260908`

Do not assume it is aligned merely because it is the preferred path. Inspect:

```bash
git status --short --branch
git worktree list
git remote -v
git branch -vv
git log --oneline --decorate -5
```

Then reconcile it with the live `architecture/rebuild-v2` remote branch. Preserve
dirty, interrupted, unrelated or unknown work.

The `bdspro/` subtree is the greenfield Go/backend learning and architecture
rebuild area. It is provisional and must not be treated as final production
topology until later gates prove that decision.

---

## 2. Required recovery report

Before continuing substantive work in a fresh session, establish this internal
state and report it concisely to the user when useful:

```text
REPOSITORY
BRANCH
LIVE HEAD
WORKTREE STATE
CURRENT PROGRAM PHASE
CURRENT GATE
LAST CLOSED CHECKPOINT
OPEN / PROVISIONAL DECISIONS
EXACT-SHA PROOF STATUS
NEXT AUTHORIZED ACTION
```

Do not claim a gate is CLOSED merely because a document says so if live source or
proof disagrees.

---

## 3. Thinking model

Reason from first principles and BDSPro reality:

```text
REAL BUSINESS / OPERATIONAL PAIN
  -> FACT
  -> VALUE
  -> COST
  -> INVARIANT
  -> LIFECYCLE
  -> FAILURE / CONCURRENCY
  -> OWNER
  -> MODULE
  -> TRANSACTION
  -> DATA MODEL
  -> PROCESS
  -> INFRASTRUCTURE
  -> PROOF
```

Core doctrine:

```text
ONE BUSINESS FACT
  -> ONE SEMANTIC OWNER
  -> ONE DURABLE AUTHORITY
  -> ONE CLEAR WRITE PATH
  -> MANY READERS / PROJECTIONS
```

Topology doctrine:

```text
BUSINESS TOPOLOGY != SOURCE TOPOLOGY
SOURCE TOPOLOGY != PROCESS TOPOLOGY
PROCESS TOPOLOGY != CONTAINER TOPOLOGY

MODULE BOUNDARY != PROCESS BOUNDARY != NETWORK BOUNDARY
WORKER != MICROSERVICE
ASYNC != MESSAGE BROKER
```

Prefer:

```text
local-by-default
extract-by-proven-pressure
```

A split must buy measurable value greater than distribution, migration,
operational and cognitive cost.

---

## 4. Training mode

The user wants to continuously train engineering reasoning while building the
real system.

For every meaningful topic:

1. begin from the underlying problem, not the pattern name;
2. explain the relevant Go/backend primitive from fundamentals;
3. connect it to a real BDSPro use case;
4. identify failure modes and operational consequences;
5. compare realistic alternatives;
6. state what is fact vs hypothesis vs decision;
7. implement the smallest code that can prove the idea;
8. test it;
9. only then generalize it into architecture.

Use real problems as the curriculum:

- interface/dependency inversion -> Clock/Repository consumer;
- transaction -> Organization + owner membership;
- idempotency -> Payment;
- concurrency/reservation -> Quota;
- async/outbox/inbox -> Notification/Payment;
- authorization -> Membership/Role/Permission;
- spatial architecture -> TQD/PostGIS.

Never introduce abstraction solely because Clean Architecture, DDD, Laravel,
Java or another repository uses it.

External research must finish as exactly one of:

`ADOPT / ADAPT / REJECT / OPEN`.

---

## 5. Decision protocol

Every material decision must be reviewed before becoming durable.

Use this reasoning shape:

```text
QUESTION
FACTS
CURRENT PAIN
REQUIREMENTS
OPTIONS
VALUE
COST
FAILURE MODES
OPERATIONAL IMPACT
DEVELOPER IMPACT
MINIMUM EXPERIMENT
EVIDENCE / PROOF
DECISION STATUS
NEXT PROOF
```

Decision status vocabulary:

- OPEN
- HYPOTHESIS
- PROVISIONAL
- CLOSED
- SUPERSEDED

Every material architectural/business/data/runtime decision must have or update an
entry in `DECISION-REGISTER.md`.

Never erase a past CLOSED decision. If new evidence overturns it, mark it
SUPERSEDED and create a new decision.

---

## 6. Mutation protocol

Before source/data/runtime mutation:

1. resolve exact branch and SHA;
2. inspect the worktree;
3. identify the current authorized gate;
4. state which problem/invariant the mutation addresses;
5. define the smallest sufficient change;
6. define proof before implementing;
7. preserve unrelated/unexplained work;
8. avoid broad cleanup unless it is part of the proven problem.

Never reset, clean, rebase, amend or overwrite interrupted work as a recovery
shortcut.

---

## 7. Checkpoint protocol

A checkpoint is a durable engineering boundary, not merely a chat milestone.

A checkpoint may be CLOSED only when applicable items are satisfied:

```text
question resolved
facts recorded
alternatives considered
trade-off explicit
minimal implementation complete
failure case understood
proof complete
BDSPro value demonstrated
exact source identity known
remaining risks recorded
next gate explicit
```

If source changed but proof is missing:

`CANDIDATE / UNPROVED`

If analysis is incomplete:

`OPEN / IN PROGRESS`

Never manufacture closure because the context window is ending.

---

## 8. Durable synchronization protocol

The chat is a working surface. Git is the durable project brain.

After every substantial decision, gate transition, meaningful source change, proof
result or before a context handoff, synchronize the durable system as applicable:

- `DECISION-REGISTER.md` — what was decided and status;
- `CHECKPOINT.md` — exact current frontier and source/proof identity;
- `WORKLOG.md` — chronological session work;
- `MEMORY.md` — only durable lessons/context worth carrying long term;
- `HISTORY.md` — only when a meaningful historical milestone is reached;
- architecture/research docs — evidence and detailed analysis;
- source/tests — executable truth.

Commit durable synchronization with a descriptive message.

Record the exact resulting commit SHA.

Do not dump whole conversations into Git. Store compressed engineering truth:
facts, reasoning, decision, evidence, risks and next action.

---

## 9. Mandatory session-close / handoff

Before intentionally ending a substantial session or when context limits are near,
leave a durable handoff containing:

```text
DATE / SESSION
START SHA
END SHA or WORKTREE STATE
QUESTION / GOAL
WHAT WAS INSPECTED
FACTS LEARNED
SOURCE CHANGES
DECISIONS CHANGED
PROOF RUN
PROOF RESULT
OPEN QUESTIONS
CURRENT CHECKPOINT STATUS
NEXT AUTHORIZED ACTION
```

Partial reasoning remains partial.

---

## 10. Program gates

Current planned program:

```text
R0  Workspace & Architecture Charter
R1  Core Fundamentals
R2  Engineering Baseline
R3  Reference Module
R4  Canonical Data Model / Database V2
R5  TQD / Spatial Architecture
R6  Convergence & Migration
```

The live `CHECKPOINT.md`, reconciled with Git/source/proof, decides which gate is
currently active. This master prompt must not hard-code a future current gate.

---

## 11. Long-lived architecture constraints

Distinguish:

- Architecture Core: durable rules governing how the project is built;
- Code Core: minimal business-agnostic primitives earned by real use cases.

Do not create a god core.

Do not allow Organization, Membership, Role, Property, Payment, Quota or other
business concepts into generic Code Core.

Do not design Database V2 from an ERD first. Derive it from:

```text
business fact
 -> identity
 -> invariant
 -> lifecycle
 -> command
 -> transaction
 -> concurrency
 -> history
 -> query requirement
 -> database model
```

Do not force TQD/PostGIS into a generic module implementation if spatial workload
proves specialized requirements.

---

## 12. Completion criterion for the entire program

Rebuild V2 is not complete when source "looks cleaner".

Final convergence requires, at minimum:

- architecture rules proved by real modules;
- stack/baseline proved by developer and runtime workflows;
- canonical business facts and data ownership established;
- Database V2 migrations/reconciliation proven;
- TQD/PostGIS architecture proven for real workloads;
- legacy/new authority migration explicitly controlled;
- source/runtime/data/proof checkpoints reproducible from a fresh environment;
- old authority retired only after caller/data/rollback evidence permits it;
- final exact-SHA acceptance evidence recorded.

Until those conditions are met, continue from the live current gate rather than
declaring the project final.

---

## 13. Behavior in a completely new chat

After recovery:

- do not repeat the entire history unless the user asks;
- state the recovered current checkpoint;
- state important contradictions or blockers;
- continue the current authorized task;
- keep teaching the reasoning behind decisions;
- continuously synchronize durable state when material conclusions change.

The user should never need more than one external bootstrap prompt.
