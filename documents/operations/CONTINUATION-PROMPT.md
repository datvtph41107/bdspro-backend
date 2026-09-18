# BDSPro — Canonical Cross-Chat Continuation Prompt

Continue BDSPro from durable state, not conversational memory. Memory is orientation only.

## Recover first

Read recovery Sprite `mcp-123-bdspro-final-acceptance-recovery` in order:

1. `/home/sprite/work/acceptance/CHECKPOINT.md`
2. `/home/sprite/work/acceptance/BDSPro-OPERATING-SYSTEM.md`
3. `/home/sprite/work/acceptance/LOCAL-EXECUTION-RUNBOOK.md`
4. `/home/sprite/work/acceptance/USER-LANE-TASK-PROTOCOL.md`
5. `/home/sprite/work/acceptance/CONTEXT-HANDOFF.md`
6. `/home/sprite/work/acceptance/DURABLE-STATE-ROADMAP.md`
7. `/home/sprite/work/acceptance/DEVELOPMENT-RUNTIME-LOGGING-CHECKPOINT.md`
8. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-LOGGING.md`
9. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-OWNERSHIP-R2.md`
10. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-STREAM-R3.md`
11. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-QUERY-R4.md`
12. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-LOGGER-R5.md`
13. `/home/sprite/work/acceptance/ERROR-RESPONSE-ARCHITECTURE-FINAL.md`
14. `/home/sprite/work/acceptance/SOURCE-INVENTORY-ERROR-RESPONSE-V1.md`
15. `/home/sprite/work/acceptance/PHASE-PRIORITY-RUNTIME-ERROR-RESPONSE.md`
16. `/home/sprite/work/acceptance/NEXT-ACTIONS.md`
17. `/home/sprite/work/acceptance/BDSPro_IMPLEMENTATION_ACCEPTANCE_LEDGER.md`
18. `/home/sprite/work/acceptance/OPERATING-MODE.md`
19. `/home/sprite/work/acceptance/LIVE-RECONCILIATION.md`

If Sprite is unavailable, use `/BDSPro Final Acceptance` Library mirror in the same order. Then reconcile live Git and completed exact-SHA CI before mutation.

## Operating protocol

Operating Model V2 is mandatory:

- parallelize read-only discovery and proof;
- serialize source mutation and gate decisions;
- use exactly one source writer;
- aggregate proof only when all proof lanes read the same immutable candidate SHA;
- every local user execution task must use `USER-LANE-TASK-PROTOCOL.md`;
- the user's local execution infrastructure is established through Phase 3C in `LOCAL-EXECUTION-RUNBOOK.md`;
- preserve the established anchor/reader/writer/proof separation and report-file handoff;
- do not destroy existing local files/worktrees merely to restart setup.

## Current checkpoint

- canonical `a6d0722a...`
- active `492b94a102e26b8d86575d72cca05b57911c745b`
- R1-R4 CLOSED
- R5 Payment slice CLOSED / hosted #83 SUCCESS
- R5 Social slice CLOSED / hosted #84 `35319045710` SUCCESS
- Payment logging debt 13→0; Social legacy std-log 6→0; repository debt 734→721→715
- Payment logging ratchets and Social legacy-std-log ratchet active
- R5 Assistant slice CLOSED / hosted #85 `35324490969` SUCCESS
- repository debt now 707; Assistant legacy std-log 8→0 with zero-ratchet active
- R5 phase ACTIVE; next-service read-only inventory authorized
- Error/Response FINAL DESIGN CLOSED / IMPLEMENTATION PENDING
- Auth zero-ratchet PARKED
- FINAL ACCEPTED = NO

## Next authorized sequence

1. reconcile active exact SHA `b55ce3c6...` and current worktree roles;
2. preserve Social CLOSED status unless live source proves a regression;
3. keep Assistant CLOSED unless live source proves a regression;
4. perform read-only exact-SHA next-service R5 logger inventory;
5. choose one bounded non-protected service from owner/consumer/value, not counts alone; then follow the full proof/publication sequence.
6. synchronize durable authority;
7. only after R5 adoption gate is intentionally stable, start Error/Response FINAL implementation.

Do not reopen superseded Kind/Definition/catalog/i18n or message-only Error V1 without a new proven requirement.
