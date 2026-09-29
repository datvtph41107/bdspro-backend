# BDSPro Architecture Rebuild V2

Branch: `architecture/rebuild-v2`
Baseline: `architecture/source-migration-v1@5fa292903f74396a07b9a4eeb21540a76ddbe83c`

This file is the single durable entry point for Architecture Rebuild V2.

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

The execution order is intentionally not the same as that list:

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

## Recovery

A new chat, new account, new machine or context-limited session must not rely on
conversation memory.

Read these files in this exact order:

1. `ARCHITECTURE-REBUILD-V2.md`
2. `docs/architecture/rebuild-v2/CONTINUATION-PROMPT.md`
3. `docs/architecture/rebuild-v2/OPERATING-MODEL.md`
4. `docs/architecture/rebuild-v2/HISTORY.md`
5. `docs/architecture/rebuild-v2/MEMORY.md`
6. `docs/architecture/rebuild-v2/DECISION-REGISTER.md`
7. `docs/architecture/rebuild-v2/CHECKPOINT.md`
8. `docs/architecture/rebuild-v2/WORKLOG.md`
9. `docs/architecture/rebuild-v2/RECOVERY-SCENARIOS.md`

Then reconcile those documents against immutable live Git, current branch HEAD,
the current worktree and completed exact-SHA proof before changing source.

If prose and immutable Git/source/proof disagree, immutable Git/source/proof wins.

## Current phase

```text
Program: BDSPro Architecture Rebuild V2
Branch: architecture/rebuild-v2
Baseline source checkpoint: 5fa292903f74396a07b9a4eeb21540a76ddbe83c
Phase: R0 — Workspace & Architecture Charter
Status: ACTIVE
Next design gate: R1.1 — What exactly is Core?
```

No production/business rewrite is implied by the existence of this branch.
