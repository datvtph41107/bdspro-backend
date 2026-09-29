# Continuation Prompt — BDSPro Architecture Rebuild V2

Use this file to recover the project in a completely new ChatGPT conversation,
account, machine or context window.

## Bootstrap instruction

Continue BDSPro Architecture Rebuild V2 from durable state, not conversational
memory.

Repository: `datvtph41107/bdspro-backend`
Working branch: `architecture/rebuild-v2`
Historical source-topology checkpoint:
`architecture/source-migration-v1@5fa292903f74396a07b9a4eeb21540a76ddbe83c`

First read, in order:

1. `ARCHITECTURE-REBUILD-V2.md`
2. `docs/architecture/rebuild-v2/OPERATING-MODEL.md`
3. `docs/architecture/rebuild-v2/HISTORY.md`
4. `docs/architecture/rebuild-v2/MEMORY.md`
5. `docs/architecture/rebuild-v2/DECISION-REGISTER.md`
6. `docs/architecture/rebuild-v2/CHECKPOINT.md`
7. `docs/architecture/rebuild-v2/WORKLOG.md`
8. `docs/architecture/rebuild-v2/RECOVERY-SCENARIOS.md`

After reading them:

- inspect the live branch HEAD;
- inspect the exact worktree state if available;
- compare live source to the checkpoint;
- inspect completed proof for the exact current SHA;
- preserve uncommitted/interrupted work;
- do not reset, clean, rebase, amend or overwrite unexplained work;
- do not infer current state from chat memory;
- do not redesign from patterns before recovering BDSPro facts and decisions.

Authority order:

```text
1. immutable live Git/source at exact SHA
2. completed exact-SHA local proof
3. completed exact-SHA hosted CI/proof
4. current CHECKPOINT.md
5. DECISION-REGISTER.md
6. WORKLOG.md / MEMORY.md / HISTORY.md
7. external research
8. conversational memory
```

If two sources disagree, stop the disputed conclusion, record the contradiction,
and reconcile from higher-authority evidence.

## Current frontier at publication

```text
Architecture Rebuild V2
R0 — Workspace & Architecture Charter
  -> recovery ecosystem publication
  -> architecture operating model
  -> checkpoint protocol
NEXT:
R1.1 — define Architecture Core before building Code Core
```

The goal is not to build a framework. Every abstraction must be earned by a real
problem, a trade-off and proof.
