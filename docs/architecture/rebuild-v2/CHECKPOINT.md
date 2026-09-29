# Current Checkpoint

Updated: 2026-09-29

## Source identity

Repository: `datvtph41107/bdspro-backend`

Working branch: `architecture/rebuild-v2`

Creation baseline:

`5fa292903f74396a07b9a4eeb21540a76ddbe83c`

Historical parent checkpoint:

`architecture/source-migration-v1@5fa292903f74396a07b9a4eeb21540a76ddbe83c`

On recovery, resolve the current branch HEAD from Git. Never assume this document's
publication SHA remains the branch HEAD forever.

## Program state

```text
BDSPro Architecture Rebuild V2
R0 — Workspace & Architecture Charter
STATUS: ACTIVE
```

Completed in R0:

- dedicated rebuild branch created;
- V1 baseline preserved;
- durable recovery ecosystem created;
- one canonical `MASTER-PROMPT.md` established for all new contexts;
- authority hierarchy defined;
- reasoning/training protocol defined;
- decision lifecycle and review protocol defined;
- source mutation guard defined;
- checkpoint closure protocol defined;
- session-close/handoff synchronization defined;
- recovery scenarios defined.

Not yet done:

- Architecture Core theorem;
- Code Core implementation;
- stack/baseline selection proof;
- reference business module;
- Database V2;
- TQD/PostGIS redesign;
- convergence/migration into canonical production line.

## Four program pillars

1. Core: long-term architecture rules + earned reusable primitives.
2. Database V2: canonical durable model derived from proven business truth.
3. Stack/baseline: Go/project structure/runtime/dev/proof model selected from
   value and operational cost.
4. TQD/PostGIS: specialized reconstruction after the baseline is proven.

## Planned gates

```text
R0  Workspace & Architecture Charter          <- CURRENT
R1  Core Fundamentals
R2  Engineering Baseline
R3  Reference Module
R4  Canonical Data Model / Database V2
R5  TQD / Spatial Architecture
R6  Convergence & Migration
```

## Immediate recovery/alignment gate

`R0.1 — Align preferred local workspace to live architecture/rebuild-v2` — CLOSED

Preferred local workspace:

`~/projects/bdspro-canonical-bootstrap-20260908`

Local branch alignment was verified by the user:

```text
LOCAL  = 033f0964956cc049e3237e343bc8fcc037312015
REMOTE = 033f0964956cc049e3237e343bc8fcc037312015
branch = architecture/rebuild-v2 tracking origin/architecture/rebuild-v2
MASTER-PROMPT.md = present
bdspro/README.md = present
```

After switching, these residual untracked paths appeared and must be preserved
until classified:

```text
assistant-service/
file-service/
hub-service/
shared/
tqd-service/
```

They do not invalidate branch alignment, but root-wide cleanup is not authorized.
R1.1 work may proceed inside the tracked `bdspro/` subtree while the residuals
remain untouched.

## Next authorized design gate

`R1.1 — What exactly is Core?`

R1.1 defines Architecture Core before Code Core. Do not create a broad shared
framework or migrate all existing modules.

The first code must be intentionally small and make responsibilities visible.
Every abstraction must be earned by a demonstrated problem.

The branch now contains `bdspro/` as a greenfield learning/rebuild area. This is
not yet a claim about final production topology.

## Only bootstrap prompt

```text
Continue BDSPro Architecture Rebuild V2 from durable state.
Repository: datvtph41107/bdspro-backend
Branch: architecture/rebuild-v2
Read docs/architecture/rebuild-v2/MASTER-PROMPT.md and follow it exactly.
```
