# BDSPro Final Acceptance — Canonical Checkpoint

Updated: 2026-09-19 Asia/Bangkok

Durable prose is recovery guidance. Immutable live Git/source and completed exact-SHA proof outrank durable text.

## Authority

- canonical acceptance: `final-acceptance/source-canonicalization` @ `a6d0722a36090148829ca5ff02f03f409d753bad`
- active refactor: `refactor/canonical-observability-errors-a6d0722a` @ `c6a9b121946a360d22759bde6a708bc6a35223c2`
- active commit: `feat(relay): adopt canonical logging`
- active tree: `01d08e582fb55759132cc199215303dd65f6d255`
- production lineage remains separate
- FINAL ACCEPTED = NO

Protected invariants: `shared/protobuf/**`, `organization-service/**`, `map-service/**` NO-TOUCH; preserve `make wire`; preserve `make buf`; `shared/code/deploy.sh` remains byte-identical SHA256 `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`.

## Operating model

- Operating Model V2 is ACTIVE.
- Master authority: `BDSPro-OPERATING-SYSTEM.md`.
- Local setup authority: `LOCAL-EXECUTION-RUNBOOK.md`.
- Delegated execution contract: `USER-LANE-TASK-PROTOCOL.md`.
- Parallelism rule: many read-only inventory/proof lanes, exactly one source writer, exact-SHA aggregation.
- User local execution infrastructure is established through Phase 3C; reader/writer/proof separation and report-file handoff are active.
- Important decisions must live in Git/exact proof + Sprite/Library, not conversational Memory alone.

## Runtime / Logging state

- R1 repository-owned structured log root — PROVED/CLOSED @ `016912fc4a3556e9d36f61fe9dc2605af2b6cccf`; hosted #77 SUCCESS.
- R2 native runtime ownership truth — PROVED/CLOSED @ `bb8b3caa9c8e115f14796516ea718816eda10273`; hosted #78 SUCCESS.
- R3 DEV raw process-stream mirror — PROVED/CLOSED @ `647ad4105540c34121590b6452cf258540f70e8e`; hosted #79 SUCCESS.
- R4 structured runtime query — PROVED/CLOSED @ `e6e0ad271dbe9d5eb2877520ab9189ed366606d9`; hosted #82 SUCCESS.
- R5 canonical logger adoption — ACTIVE service-by-service.
- R5 Payment slice — PROVED/CLOSED @ `fca4682587d93cbc436f2466244ee8bd03b8b1b9`; hosted #83 / `35303545759` SUCCESS.
- Payment logging debt: `go.legacy_std_log 4→0`, `go.third_party_logger 9→0`; repository debt `734→721`; both Payment zero-ratchets active.
- R5 Social slice — PROVED/CLOSED @ `b55ce3c6d8005e5d7375228196a7cafb10f8ddce`; tree `91ccd0d547227191bc20df4498bdbb4ff723c0a2`; hosted #84 / `35319045710` SUCCESS.
- Social logging debt: `go.legacy_std_log 6→0`; repository debt `721→715`; Social legacy-std-log zero-ratchet active.
- R5 Assistant slice — PROVED/CLOSED @ `492b94a102e26b8d86575d72cca05b57911c745b`; hosted #85 / `35324490969` SUCCESS.
- Assistant logging debt: `go.legacy_std_log 8→0`; repository debt `715→707`; Assistant legacy-std-log zero-ratchet active.

## Error / Response FINAL

DESIGN REFINED / IMPLEMENTATION PENDING. The consumer-proven target now carries protocol `code` + stable low-cardinality application `reason` + public-safe `message`. Preferred one-way caller-facing vocabulary for implementation review is conceptually `failure.New(code, reason, safeMessage)`; exact package/signature migration still requires live compatibility proof. One gRPC normalizer and one Gateway projection remain the boundary owners. Technical failures remain wrapped Go errors and are sanitized only after internal evidence is preserved.

## Recovery order

Read `CHECKPOINT.md`, then `BDSPro-OPERATING-SYSTEM.md`, `ARCHITECTURE-NAMING-OWNERSHIP-MINDSET.md`, `LOCAL-EXECUTION-RUNBOOK.md`, `USER-LANE-TASK-PROTOCOL.md`, handoff/roadmap, R1-R4 runtime inventories, `SOURCE-INVENTORY-RUNTIME-LOGGER-R5.md`, Error/Response FINAL authority, next-actions, ledger, operating-mode and live reconciliation; then reconcile live Git/GitHub before mutation.

Historical/superseded error design remains historical only. Memory is orientation only.

`ONE TRUTH PER CONCERN -> ONE OWNER -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`


## Naming / ownership / public-surface refinement

Durable doctrine:
`ARCHITECTURE-NAMING-OWNERSHIP-MINDSET.md`.

Active principle:
`ONE TRUTH PER CONCERN -> ONE OWNER -> ONE PUBLIC NAME -> ONE OBVIOUS ENTRY POINT -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`.

Examples:
- split catch-all `QuotaManager` by invariant into pricing/reservation/durable-ledger ownership where source reality supports it;
- converge long-lived daily logging UX toward one `make logs` public name while preserving raw vs structured internal mechanics and R4 proof history;
- Error/Response refined target uses protocol code + stable application reason + public-safe message;
- preferred one-way caller-facing naming is conceptually `failure.New(code, reason, safeMessage)`, subject to live compatibility/source migration proof;
- do not revive Kind/Definition/catalog/i18n or a parallel numeric business-code truth.

## Current Notification interrupted-state checkpoint

Remote active refactor remains `c6a9b121946a360d22759bde6a708bc6a35223c2` unless live Git proves otherwise.

Local Notification candidate:
- SHA `381c29365a00f35737bb0b6078cf0973198dc7c2`;
- tree `fe7c5153c65074bf406134e4e966c9d097a2d706`;
- parent `c6a9b121946a360d22759bde6a708bc6a35223c2`;
- exact twelve-file candidate;
- not pushed.

The corrected detached exact-SHA proof at the same immutable candidate is now fully green with `FINAL_FAIL_COUNT=0`. Protobuf/Wire, empty post-generation tidy, baseline parity, all non-baseline tests, build, audit-tool tests, enforce-ratchets, Notification zero debt, repository debt 661, logger/process shape, protected paths, deploy hash and detached immutability all PASS.

Do not amend/reset/recreate the candidate. Safe non-force fast-forward publication of exact candidate `381c29365a00f35737bb0b6078cf0973198dc7c2` is PASS. Remote active refactor now points at that exact candidate. Hosted exact-SHA proof for workflow `Refactor Observability and Error Contracts` is the next gate; Notification is not CLOSED until inventory/common-contracts/boundary-contracts and exact-SHA inventory artifact are verified.

FINAL ACCEPTED = NO.
