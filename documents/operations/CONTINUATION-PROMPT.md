# BDSPro — Canonical Cross-Chat Continuation Prompt

Updated: 2026-09-19 Asia/Bangkok

Continue BDSPro from durable state, not conversational memory. Memory is orientation only.

## Recover first

Read recovery Sprite `mcp-123-bdspro-final-acceptance-recovery` in this order:

1. `/home/sprite/work/acceptance/CHECKPOINT.md`
2. `/home/sprite/work/acceptance/BDSPro-OPERATING-SYSTEM.md`
3. `/home/sprite/work/acceptance/ARCHITECTURE-NAMING-OWNERSHIP-MINDSET.md`
4. `/home/sprite/work/acceptance/LOCAL-EXECUTION-RUNBOOK.md`
5. `/home/sprite/work/acceptance/USER-LANE-TASK-PROTOCOL.md`
6. `/home/sprite/work/acceptance/CONTEXT-HANDOFF.md`
7. `/home/sprite/work/acceptance/DURABLE-STATE-ROADMAP.md`
8. `/home/sprite/work/acceptance/DEVELOPMENT-RUNTIME-LOGGING-CHECKPOINT.md`
9. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-LOGGING.md`
10. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-OWNERSHIP-R2.md`
11. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-STREAM-R3.md`
12. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-QUERY-R4.md`
13. `/home/sprite/work/acceptance/SOURCE-INVENTORY-RUNTIME-LOGGER-R5.md`
14. `/home/sprite/work/acceptance/ERROR-RESPONSE-ARCHITECTURE-FINAL.md`
15. `/home/sprite/work/acceptance/SOURCE-INVENTORY-ERROR-RESPONSE-V1.md`
16. `/home/sprite/work/acceptance/PHASE-PRIORITY-RUNTIME-ERROR-RESPONSE.md`
17. `/home/sprite/work/acceptance/NEXT-ACTIONS.md`
18. `/home/sprite/work/acceptance/BDSPro_IMPLEMENTATION_ACCEPTANCE_LEDGER.md`
19. `/home/sprite/work/acceptance/OPERATING-MODE.md`
20. `/home/sprite/work/acceptance/LIVE-RECONCILIATION.md`

If Sprite is unavailable, use `/BDSPro Final Acceptance` Library mirror in the same order.

After durable recovery, reconcile immutable live Git/source, current worktree roles, local candidate/proof state and completed exact-SHA hosted proof before any mutation. If durable prose and live Git/proof disagree, live Git/proof wins.

## Operating protocol

Operating Model V2 is mandatory:

- parallelize read-only discovery and proof;
- serialize architecture/gate decisions and source mutation;
- use exactly one source writer;
- aggregate proof only when all proof lanes read the same immutable candidate SHA;
- preserve anchor/reader/writer/proof separation and report-file handoff;
- preserve interrupted worktrees/evidence; do not reset, clean, rebase or recreate state merely to resume;
- every local user execution task follows `USER-LANE-TASK-PROTOCOL.md`;
- never infer the next mutation from a failed proof; classify the failure first;
- durable state must be synchronized after material decisions, blockers, publication/hosted closure or operating-model changes.

## Working doctrine

Use:

`ONE TRUTH PER CONCERN -> ONE OWNER -> ONE PUBLIC NAME -> ONE OBVIOUS ENTRY POINT -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

and:

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

and:

`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`

Do not create catch-all owners such as `*Manager` when multiple independent invariants are hidden inside. Split by responsibility/invariant (for example pricing vs reservation vs durable usage), not mechanically by method count.

Do not create multiple public aliases for one caller intent merely because internal mechanics differ. The current target is one public `make logs` vocabulary while raw process output and structured evidence/query remain distinct internal mechanics; R4 `make log-query` remains historical/proved implementation until a bounded compatibility/public-surface retirement is explicitly authorized and proved.

## Error/Response refined target

Error/Response FINAL remains implementation-pending until R5 is deliberately stable.

New proven requirement:
- caller-facing application failure needs a stable machine-readable `reason` in addition to protocol `code` and safe human `message`;
- reason is low-cardinality/static application identity, not a dynamic value, service topology name or technical root cause;
- do not retain a second mandatory numeric business-code truth;
- response/logging/metrics/tracing project the same canonical code/reason rather than remapping them independently;
- technical Go errors remain wrapped technical errors and are sanitized only after internal evidence is preserved;
- expected business rejection is not automatically an operational ERROR.

Naming refinement for implementation review:
- prefer one dedicated application-failure owner and one obvious constructor, conceptually `failure.New(code, reason, safeMessage)`;
- ordinary Go errors remain ordinary `error`/`fmt.Errorf("%w")`/`errors.Is/As`;
- exact package/signature migration must be reconciled against live source and compatibility before mutation.

Do not reopen superseded Error V1, Kind/Definition/catalog/i18n, message-only Error, or rich generic metadata designs unless a new proven consumer/invariant explicitly requires them.

## Current source/proof checkpoint

Remote active refactor authority remains:
`381c29365a00f35737bb0b6078cf0973198dc7c2`
unless live Git proves otherwise.

Closed:
- R1 #77;
- R2 #78;
- R3 #79;
- R4 #82;
- R5 Payment #83;
- R5 Social #84;
- R5 Assistant #85;
- R5 Relay #86;
- R5 Notification hosted run `35405779633`.

Remote repository debt at the last hosted authority: 661.

Notification R5 is PROVED/CLOSED:
- SHA `381c29365a00f35737bb0b6078cf0973198dc7c2`;
- tree `fe7c5153c65074bf406134e4e966c9d097a2d706`;
- parent `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- exact twelve-file slice;
- corrected detached exact-SHA proof PASS;
- safe non-force publication PASS;
- hosted run `35405779633`: inventory/common-contracts/boundary-contracts SUCCESS;
- hosted artifact `10572063928`, digest `sha256:8c8b98a33e56dda9aea3db532b4770437f729aa8d8e8a95988da8e453529c61d`;
- Notification debt 30 -> 0;
- repository debt 691 -> 661;
- Notification legacy std-log zero-ratchet active.

R5 remains ACTIVE. Error/Response implementation has not started. FINAL ACCEPTED = NO.

## Next authorized sequence

1. Reconcile live remote active branch against exact SHA `381c29365a00f35737bb0b6078cf0973198dc7c2`.
2. Keep Payment, Social, Assistant, Relay and Notification R5 slices CLOSED unless exact source/proof proves regression.
3. Perform fresh exact-SHA read-only logging inventory/source/coupling trace for the next non-protected owner.
4. Select exactly one bounded R5 slice from ownership/coupling evidence, not raw count alone.
5. Execute the full proof lifecycle and hosted closure before advancing.
6. Stay in R5 until canonical logger adoption is deliberately declared stable.
7. Only after R5 stability begin Error/Response refined FINAL implementation.
8. For Error/Response: establish one application-failure core/constructor, one gRPC normalization path, one Gateway projection, common reason semantics across response/observability, vertical proof, service-by-service migration, retirement, ratchets and protected compatibility resolution.

FINAL ACCEPTED remains NO.
