# BDSPro Canonical Database — Final Architecture Package

Date: 2026-10-06
Status: FINAL TARGET ARCHITECTURE / IMPLEMENTATION NOT YET CUT OVER
Base source: main @ b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd

## Purpose

This package closes the database-architecture research pass that started at v0. It is the durable target used to review and implement the database convergence of BDSPro.

The master DBML is a logical system diagram. It is not a proposal for one shared physical PostgreSQL database. Each bounded context keeps one durable owner. Cross-owner references show semantic identity only; implementation must use service contracts/events and must not introduce cross-database transactions.

## Authority order

1. Immutable live Git/source and completed exact-SHA proof.
2. CANONICAL-DATABASE-FINAL.dbml plus DATABASE-INVARIANTS.md.
3. DATABASE-FINAL-ACCEPTANCE.md.
4. MIGRATION-CUTOVER-PLAN.md and LEGACY-DISPOSITION.md.
5. DATABASE-VERSION-EVOLUTION.md.
6. Conversation memory only for orientation.

If prose conflicts with live source before cutover, reconcile the conflict explicitly. After a target fact has switched authority, the canonical owner defined here wins.

## Package

- CANONICAL-DATABASE-FINAL.dbml — complete logical target diagram across durable owners.
- DATABASE-VERSION-EVOLUTION.md — v0 → v1 → v2 → Final decisions.
- DATABASE-INVARIANTS.md — invariants, transaction boundaries, concurrency and authorization consequences.
- LEGACY-DISPOSITION.md — current tables/patterns and KEEP / REPLACE / PROJECTION / OPERATIONAL / RETIRE decisions.
- MIGRATION-CUTOVER-PLAN.md — expand/backfill/compare/switch/contract plan.
- DATABASE-FINAL-ACCEPTANCE.md — acceptance report and implementation gates.
- CONTINUATION-PROMPT.md — bootstrap for the next chat.

## Governing rules

- One business fact → one owner → one source of truth → one clear write path.
- Same shape does not imply same fact.
- Invitation is not participation.
- Creator is not current authority.
- Relationship is not authorization.
- Current relation and historical fact are different concerns.
- Generic type/id ownership is forbidden for canonical core truth.
- Generic references remain allowed in audit/event/projection records where the reference is descriptive rather than relational authority.
- Money is integer minor units plus currency.
- Unknown is kept explicit instead of inventing precision.
- A table/column exists only when a business fact, invariant, operational proof or required projection earns it.

## Bounded-context owners

| Area | Durable owner | Direction |
| --- | --- | --- |
| Person, account, authenticator, session, platform IAM | User | canonical identity/auth truth |
| Organization, membership, org IAM, branch | User/Organization convergence owner | one owner after cutover |
| Group, Deal, Deal participation/economics | Organization/Deal | business truth |
| Property, Asset, Listing | BDSPro | business truth |
| CRM, Support, SEO | CRM | business truth / projections where noted |
| Conversation, Message | Chat | business truth |
| Post, Comment, Reaction | Social | business truth |
| Catalog, Subscription | User Commercial | commercial truth |
| Order, Attempt, Settlement, Fulfillment, Wallet ledger | Payment | payment truth |
| File metadata/access | File | storage metadata truth |
| Notification inbox/delivery | Notification | effect truth, not foreign-domain audit |
| Administrative geography/config | Hub | platform reference truth |
| Planning, Report, Quota | TQD | planning/report/usage truth |

## What Final means

Final means the target semantic boundaries and owner model are closed enough to implement. It does not mean the production databases have already been migrated. Runtime schemas remain unchanged by this documentation commit.

Implementation starts only through MIGRATION-CUTOVER-PLAN.md. No old table may be dropped because a target table appears in the DBML; authority must be switched only after backfill and compare proof.

## Recovery order for a new chat

Read CONTINUATION-PROMPT.md first and follow its exact recovery order. Do not reconstruct state from conversational memory.