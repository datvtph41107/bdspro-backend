# BDSPro Architecture Rebuild V2

Branch: `architecture/rebuild-v2`
Baseline: `architecture/source-migration-v1@5fa292903f74396a07b9a4eeb21540a76ddbe83c`

This file is the durable project entry point.

## One-prompt rule

There is exactly one canonical orchestration prompt for the complete program:

`docs/architecture/rebuild-v2/MASTER-PROMPT.md`

A completely new chat/account/session needs only:

```text
Continue BDSPro Architecture Rebuild V2 from durable state.
Repository: datvtph41107/bdspro-backend
Branch: architecture/rebuild-v2
Read docs/architecture/rebuild-v2/MASTER-PROMPT.md and follow it exactly.
```

Do not create competing continuation prompts. The master prompt owns recovery,
reasoning, training, decision review, mutation guards, proof, checkpointing,
durable synchronization and final convergence.

## Purpose

Rebuild the BDSPro engineering baseline from first principles while preserving the
existing system as evidence, not as an unquestioned template.

The program has four long-lived pillars:

1. architecture core: the rules that determine how code, modules, dependencies,
   transactions, errors, events, tests and runtime roles are designed;
2. canonical database redesign derived from business facts, invariants, lifecycle,
   transaction and query requirements;
3. stack, source structure and engineering baseline chosen by measured value and
   operational cost rather than convention;
4. careful reconstruction of TQD/planning/spatial workloads, including PostGIS,
   after the baseline is proven.

Execution direction:

```text
FIRST PRINCIPLES
  -> CORE CONTRACT
  -> ENGINEERING BASELINE
  -> REFERENCE MODULE
  -> BUSINESS FACT MODEL
  -> DATABASE V2
  -> TQD / POSTGIS
  -> SYSTEM CONVERGENCE
```

## Durable state

The master prompt recovers and reconciles:

- operating model;
- history;
- durable project memory;
- decision register;
- current checkpoint;
- worklog;
- recovery scenarios;
- live Git/source/worktree;
- exact-SHA proof.

If prose and immutable Git/source/proof disagree, immutable Git/source/proof wins.

No production/business rewrite is implied merely by the existence of this branch.
