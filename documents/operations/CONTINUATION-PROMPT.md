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
- the user's local setup currently restarts from Phase 1 verification in `LOCAL-EXECUTION-RUNBOOK.md`;
- do not assume existing local repo/worktrees/scripts are correct until Phase 1 is reconciled;
- do not destroy existing local files/worktrees merely to restart setup.

## Current checkpoint

- canonical `a6d0722a...`
- active `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- R1-R4 CLOSED
- R5 Payment slice CLOSED / hosted #83 SUCCESS
- Payment logging debt 13→0; repository debt 734→721; Payment logging ratchets active
- R5 phase ACTIVE; next-service inventory authorized after local setup verification
- Error/Response FINAL DESIGN CLOSED / IMPLEMENTATION PENDING
- Auth zero-ratchet PARKED
- FINAL ACCEPTED = NO

## Next authorized sequence

1. complete/reconcile the user's local Phase 1 repository verification;
2. reconcile active exact SHA and clean source;
3. perform read-only next-service R5 logger inventory;
4. choose one bounded non-protected service from owner/consumer/value, not counts alone;
5. bounded migration -> zero proof -> retirement -> ratchet -> local exact-SHA proof -> safe publication -> hosted proof;
6. synchronize durable authority;
7. only after R5 adoption gate is intentionally stable, start Error/Response FINAL implementation.

Do not reopen superseded Kind/Definition/catalog/i18n or message-only Error V1 without a new proven requirement.
