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
