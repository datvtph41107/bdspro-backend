# Decision Register

Status vocabulary:

- OPEN — insufficient evidence.
- HYPOTHESIS — candidate direction to test.
- PROVISIONAL — supported enough to guide limited implementation.
- CLOSED — decision proved for current scope.
- SUPERSEDED — replaced by a later decision.

| ID | Decision / question | Status | Evidence / rationale | Next proof |
| --- | --- | --- | --- | --- |
| ARV2-001 | Use a dedicated `architecture/rebuild-v2` lane instead of overwriting Source Migration V1 | CLOSED | Preserves V1 as immutable historical checkpoint and isolates rebuild work | Keep lineage explicit |
| ARV2-002 | Continuity must be Git/durable-document based, not chat-memory based | CLOSED | Account/context/machine changes are expected recovery cases | Validate continuously through recovery |
| ARV2-003 | Architecture Core and Code Core are distinct | PROVISIONAL | Prevents "god core" while preserving long-lived engineering rules | Prove during R1 |
| ARV2-004 | Every abstraction must be earned by a real use case | PROVISIONAL | Reduces speculative framework cost | First executable vertical slice |
| ARV2-005 | Prefer business-first vertical modules over global technical-type folders | HYPOTHESIS | Current User/Property/Payment inner structures show cognitive fragmentation | Reference module comparison |
| ARV2-006 | Local/module call is default; network split requires proven pressure | HYPOTHESIS | Distributed operation carries real cost; source/process/network are distinct | Evaluate Payment and TQD exceptions |
| ARV2-007 | PostgreSQL is default durable authority; Redis/broker are not default truth | HYPOTHESIS | Simplifies ownership/recovery unless workload requires projection/async mechanics | Reference module + workload proof |
| ARV2-008 | Database V2 is derived after business fact/lifecycle/transaction theorem | PROVISIONAL | Avoids encoding ambiguous business semantics in schema | Org/Membership fact model |
| ARV2-009 | TQD/PostGIS gets specialized treatment after baseline proof | PROVISIONAL | Spatial workload is materially different and already exists | TQD workload/query inventory |
| ARV2-010 | External Laravel/Java/OSS sources are evidence, never architectural authority | CLOSED | BDSPro-specific value/cost must govern decisions | Apply ADOPT/ADAPT/REJECT/OPEN protocol |
| ARV2-011 | Rebuild V2 starts at R0 charter/recovery before new business implementation | CLOSED | Continuity and governance are prerequisites for durable experimentation | Complete R0 |
| ARV2-012 | One canonical master prompt orchestrates recovery through final convergence | CLOSED | Eliminates competing bootstrap procedures and lets a zero-context session self-recover from durable Git state | Validate in a future fresh context |
| ARV2-013 | Every material decision must be reviewed and synchronized into durable state | CLOSED | Long-lived architecture training requires decisions, proof and checkpoint history to survive chat/account/context loss | Enforce through MASTER-PROMPT session-close protocol |

| ARV2-014 | Use `~/projects/bdspro-canonical-bootstrap-20260908` as the preferred local WSL working copy for Rebuild V2, with mandatory live-state reconciliation before mutation | CLOSED | The existing local repository carries historical branches/worktrees and is the intended daily workspace | Validate alignment to live `architecture/rebuild-v2` before R1.1 |
| ARV2-015 | Create a root `bdspro/` greenfield learning/rebuild area inside Rebuild V2 | PROVISIONAL | Separates from legacy implementation while allowing fundamentals-to-production learning on real BDSPro problems | Prove during R1/R2; do not treat as final topology yet |
| ARV2-016 | Local workspace alignment is an explicit R0.1 gate before new Rebuild V2 source work | CLOSED | Prevents remote durable state and local historical checkout from silently diverging | Verify local branch/head/worktree after fetch/switch |

| ARV2-017 | Use `uber-go/guide` as a foundational Go engineering reference, consulted problem-first and never as BDSPro architecture authority | CLOSED | The guide covers interfaces, receivers, errors, lifecycle, concurrency, construction, testing, naming and tooling; Rebuild V2 already requires external sources to be evaluated through BDSPro value/cost and ADOPT/ADAPT/REJECT/OPEN | Apply only when a matching real problem appears and prove the choice in BDSPro code/tests |
When a decision changes, do not erase the old row. Mark it SUPERSEDED and add a
new decision with the reason and evidence.
