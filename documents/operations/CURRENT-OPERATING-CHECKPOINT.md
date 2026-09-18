# BDSPro Final Acceptance — Canonical Checkpoint

Updated: 2026-09-18 Asia/Ho_Chi_Minh

Durable prose is recovery guidance. Immutable live Git/source and completed exact-SHA proof outrank durable text.

## Authority

- canonical acceptance: `final-acceptance/source-canonicalization` @ `a6d0722a36090148829ca5ff02f03f409d753bad`
- active refactor: `refactor/canonical-observability-errors-a6d0722a` @ `fca4682587d93cbc436f2466244ee8bd03b8b1b9`
- active commit: `feat(payment): adopt canonical logging`
- active tree: `ac2ff2a1a98aec5afa647edb1584683aeb5d545f`
- production lineage remains separate
- FINAL ACCEPTED = NO

Protected invariants: `shared/protobuf/**`, `organization-service/**`, `map-service/**` NO-TOUCH; preserve `make wire`; preserve `make buf`; `shared/code/deploy.sh` remains byte-identical SHA256 `80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e`.

## Operating model

- Operating Model V2 is ACTIVE.
- Master authority: `BDSPro-OPERATING-SYSTEM.md`.
- Local setup authority: `LOCAL-EXECUTION-RUNBOOK.md`.
- Delegated execution contract: `USER-LANE-TASK-PROTOCOL.md`.
- Parallelism rule: many read-only inventory/proof lanes, exactly one source writer, exact-SHA aggregation.
- User local setup is intentionally restarted from Phase 1 verification; existing files/processes are preserved until verified.
- Important decisions must live in Git/exact proof + Sprite/Library, not conversational Memory alone.

## Runtime / Logging state

- R1 repository-owned structured log root — PROVED/CLOSED @ `016912fc4a3556e9d36f61fe9dc2605af2b6cccf`; hosted #77 SUCCESS.
- R2 native runtime ownership truth — PROVED/CLOSED @ `bb8b3caa9c8e115f14796516ea718816eda10273`; hosted #78 SUCCESS.
- R3 DEV raw process-stream mirror — PROVED/CLOSED @ `647ad4105540c34121590b6452cf258540f70e8e`; hosted #79 SUCCESS.
- R4 structured runtime query — PROVED/CLOSED @ `e6e0ad271dbe9d5eb2877520ab9189ed366606d9`; hosted #82 SUCCESS.
- R5 canonical logger adoption — ACTIVE service-by-service.
- R5 Payment slice — PROVED/CLOSED @ `fca4682587d93cbc436f2466244ee8bd03b8b1b9`; hosted #83 / `35303545759` SUCCESS.
- Payment logging debt: `go.legacy_std_log 4→0`, `go.third_party_logger 9→0`; repository debt `734→721`; both Payment zero-ratchets active.
- R5 next gate: read-only exact-SHA inventory for the next bounded non-protected service; no bulk migration.

## Error / Response FINAL

DESIGN CLOSED / IMPLEMENTATION PENDING. Canonical caller-facing form remains `return _errors.ReturnError(codes.SomeCode, "Safe message")`. One core Error owns standard gRPC code + safe message; one gRPC normalizer; one Gateway projection. Technical failures remain wrapped Go errors.

## Recovery order

Read `CHECKPOINT.md`, then `BDSPro-OPERATING-SYSTEM.md`, `LOCAL-EXECUTION-RUNBOOK.md`, `USER-LANE-TASK-PROTOCOL.md`, handoff/roadmap, R1-R4 runtime inventories, `SOURCE-INVENTORY-RUNTIME-LOGGER-R5.md`, Error/Response FINAL authority, next-actions, ledger, operating-mode and live reconciliation; then reconcile live Git/GitHub before mutation.

Historical/superseded error design remains historical only. Memory is orientation only.

`ONE TRUTH PER CONCERN -> ONE OWNER -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`
