# BDSPro Canonical Database — Final Architecture Package

Status: **architecture acceptance candidate**  
Branch: `architecture/canonical-database-final`  
Authority order used by this package: **live Git/source > completed proof/tests > durable architecture docs > conversation**.

## Purpose

This package closes the database-design reasoning that started at v0 and carries it through v1, v2, and final acceptance.

The goal is not to copy the legacy schema. The goal is to define the smallest responsible set of durable business facts, their owners, invariants, transaction boundaries, and migration path.

## Version ladder

### v0 — Semantic model

Questions answered before physical schema:

- what business fact does each row assert?
- what is the identity of that fact?
- what is proposal vs effective relationship vs outcome?
- what is current authority vs provenance/history?
- what can be derived instead of stored?
- which generic `(type,id)`, `status`, `owner`, or soft-delete constructs are hiding different facts?
- which service owns the write truth?

See: [V0-SEMANTIC-MODEL.md](./V0-SEMANTIC-MODEL.md)

### v1 — Logical relational model

The accepted semantic facts are expressed as a normalized logical ER model. Cross-service references are shown logically even when no physical PostgreSQL foreign key can exist across databases.

See: [V1-CANONICAL-CORE.dbml](./V1-CANONICAL-CORE.dbml)

### v2 — Physical PostgreSQL reference

v2 turns the v1 model into enforceable PostgreSQL DDL patterns:

- keys and foreign keys
- partial unique indexes
- CHECK constraints
- explicit lifecycle/outcome constraints
- monetary precision
- concurrency/locking notes
- same-context invariants that require transactional validation
- clear separation of authoritative business tables from projections/audit/outbox state

See: [V2-PHYSICAL-POSTGRESQL.sql](./V2-PHYSICAL-POSTGRESQL.sql)

### final — Acceptance and migration

Final acceptance does **not** mean "drop all legacy tables now." It means the target authority is explicit and the migration can proceed using:

`OLD → EXPAND → BACKFILL → COMPARE/PROVE → SWITCH AUTHORITY → CONTRACT`

See:

- [SERVICE-TABLE-INVENTORY.md](./SERVICE-TABLE-INVENTORY.md)
- [MIGRATION-AND-PROOF-PLAN.md](./MIGRATION-AND-PROOF-PLAN.md)
- [FINAL-ACCEPTANCE.md](./FINAL-ACCEPTANCE.md)
- [CONTINUATION-PROMPT.md](./CONTINUATION-PROMPT.md)

## Core rules preserved

1. Intent is not a table.
2. Proposal/request is not an effective relationship.
3. Same shape is not the same fact.
4. A foreign ID does not establish business ownership.
5. One mutable fact has one canonical write owner.
6. Current state, history, audit, events, and projections are different records.
7. `status` is accepted only for one coherent mutually-exclusive lifecycle.
8. NULL means absence; defaults are business assertions.
9. Money is amount + currency.
10. Generic `owner_type + owner_id` is rejected for canonical business truth when the endpoint is known.
11. A role/capacity is not Party identity.
12. Creator is provenance, not current authority.
13. Cross-service access is by contract/projection/event; databases are not joined through runtime cross-DB foreign keys.
14. Transaction boundaries are earned by invariants.
15. Unknown is kept explicit rather than converted into invented precision.

## Physical ownership

The target deliberately remains a microservice data architecture rather than one monolithic PostgreSQL schema.

- **user-service**: Party, Person, Account/authentication, presentation Profile, Organization identity/governance/membership/role/branch facts, catalog/subscription authority.
- **organization-service** (deal/collaboration bounded context): Group and Deal facts, Deal participation/governance/financial terms.
- **payment-service**: payment order/attempt/settlement/fulfillment/outbox, wallets, banks/payment methods and canonical bank-account/payment facts.
- **bdspro-service**: Property truth and property projections/partitions.
- **tqd-service**: planning, map, quota, report, AI-job and planning-workflow truth.
- **crm-service**: CRM/support/SEO/measurement truth.
- **notification-service**: notification delivery/history and payment-event inbox.
- **file-service**: file metadata/access truth.
- **hub-service**: platform lookup/config/help/event compatibility state.

A logical relationship across these boundaries is validated through the owning service contract. It is **not** turned into an impossible cross-database FK.
