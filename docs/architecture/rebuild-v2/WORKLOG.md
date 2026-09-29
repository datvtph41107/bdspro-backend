# Worklog

Append substantial sessions. This is a chronological operational log, not the
canonical decision authority.

## 2026-09-29 — Rebuild V2 bootstrap

### Starting point

- Repository: `datvtph41107/bdspro-backend`
- Source Migration V1 checkpoint:
  `5fa292903f74396a07b9a4eeb21540a76ddbe83c`
- New branch: `architecture/rebuild-v2`

### Why

The project had already completed a long evolution from service-level refactoring
toward first-principles thinking about value, operational cost, ownership,
business facts, source/process/runtime separation and simpler architecture.

A separate lane was needed to:

- build from fundamentals rather than keep patching old inner structures;
- act as a Go/backend learning corpus tied to real BDSPro problems;
- preserve old source as evidence;
- allow analysis/experimentation before checkpoint closure;
- prevent chat-memory loss from destroying project continuity.

### Durable work completed

Created the Rebuild V2 continuity ecosystem:

- root entry point;
- operating model;
- durable history;
- durable memory;
- decision register;
- current checkpoint;
- recovery scenarios;
- worklog.

### Current frontier

`R0 — Workspace & Architecture Charter`

Next: `R1.1 — What exactly is Core?`

---

## 2026-09-29 — Single master prompt

### Goal

Remove the need for multiple recovery prompts and make one zero-context prompt
sufficient from the start of Rebuild V2 through final convergence.

### Changes

- added `docs/architecture/rebuild-v2/MASTER-PROMPT.md`;
- made it the only canonical orchestration prompt;
- converted the old continuation prompt to a compatibility pointer;
- added mandatory recovery report, thinking/training model, decision review,
  mutation guard, checkpoint closure and durable synchronization protocols;
- added rules for context-limit handoff and final-program completion;
- updated the decision register and current checkpoint.

### Decision

A new chat/account/session must bootstrap from the master prompt, then dynamically
resolve the current gate from durable state and live Git. The bootstrap prompt
must not encode a current SHA or phase that will become stale.

### Current frontier

Still `R0 — Workspace & Architecture Charter`.

Next authorized design gate remains:

`R1.1 — What exactly is Core?`

### Session-close template

For future entries append:

```text
DATE / SESSION
START SHA:
WORKTREE STATE:
QUESTION:
FACTS LEARNED:
SOURCE CHANGES:
DECISIONS:
PROOF RUN:
PROOF RESULT:
OPEN QUESTIONS:
CHECKPOINT STATUS:
END SHA:
NEXT AUTHORIZED ACTION:
```


---

## 2026-09-29 — Preferred local workspace and greenfield learning area

### Local state supplied for reconciliation

Preferred WSL workspace:

`~/projects/bdspro-canonical-bootstrap-20260908`

Observed before Rebuild V2 alignment:

```text
branch: refactor/canonical-observability-errors-a6d0722a
HEAD: fca4682
```

Historical/local branches are present and another worktree is in use, so no
destructive cleanup is authorized.

### Durable decisions

- make this path the preferred daily Rebuild V2 workspace;
- require a fetch + worktree/status reconciliation before switching it;
- create `bdspro/` as an intentionally greenfield training/rebuild subtree;
- keep `bdspro/` provisional until R1/R2 prove its role and shape.

### Immediate gate

`R0.1 — Local Workspace Alignment`

After the user verifies the local checkout tracks the live
`architecture/rebuild-v2` branch, continue:

`R1.1 — What exactly is Core?`


---

## 2026-09-29 — Local Rebuild V2 alignment verified

The user aligned the preferred WSL workspace:

`~/projects/bdspro-canonical-bootstrap-20260908`

Verified:

```text
branch: architecture/rebuild-v2
local HEAD: 033f0964956cc049e3237e343bc8fcc037312015
remote HEAD: 033f0964956cc049e3237e343bc8fcc037312015
LOCAL == REMOTE: OK
MASTER PROMPT: OK
BDSPro learning workspace: OK
```

After the branch switch, residual untracked directories appeared:

```text
assistant-service/
file-service/
hub-service/
shared/
tqd-service/
```

They are preserved as unknown/historical local residue. No cleanup is authorized
until classified.

R0.1 branch alignment is CLOSED. R1.1 may proceed inside `bdspro/`.

The user also explicitly chose a hands-on learning style: personally type/run the
commands and code, return real terminal output, and build engineering reflexes
through small production-relevant steps. This style is now a durable working rule.


---

## 2026-09-29 — Natural-reflex training rule strengthened

The user requested that every command and code fragment be explained from:
immediate purpose, underlying mechanism, and relationship to the current and
program-level goals.

The working style now explicitly maintains a goal hierarchy and derives the next
question from observed fundamentals rather than from a predetermined pattern.
