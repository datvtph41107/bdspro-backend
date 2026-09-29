# Durable History

This is the compact history required to understand why Rebuild V2 exists.
It is not a substitute for Git history.

## H0 — Existing BDSPro

BDSPro evolved as a Go multi-service system with Gateway, User, Auth,
Organization, Property/BDSPro, CRM, TQD/planning, Payment, Notification, File,
Hub, Assistant, Chat, Search and other runtime areas.

The system paid real distributed-system costs: many modules/processes, gRPC,
PostgreSQL, Redis, RabbitMQ, Docker/Compose, per-service config/migrations,
cross-service debugging, CI, deployment and recovery.

## H1 — Final Acceptance

The prior Final Acceptance program closed the repository/implementation acceptance
lineage. Historical accepted production-promotion source was
`3c8e341211dccff653b040466fec3480c5f7c8d5`.

This meant the accepted repository state was proved; it did not mean an external
production deployment was performed.

## H2 — Post-acceptance architecture learning

The program then shifted from "make existing services pass" to first-principles
architecture analysis:

- distinguish business/source/runtime/infrastructure topology;
- measure cognitive and operational cost;
- reconstruct ownership and business facts;
- study external systems such as Laravel/Java/open-source repositories as evidence;
- prefer local/module boundaries until extraction pressure is proven.

Business-fact reconstruction began around Organization, Ownership, Invitation,
Membership, Role and authorization.

## H3 — Main foundation

Two documentation/foundation commits followed Final Acceptance:

```text
3c8e3412...
  -> eed3e151... docs(auth): establish naming and ownership mindset
  -> b2a7a816... docs(auth): define organization and authorization foundation
```

## H4 — Architecture Source Migration V1

Branch `architecture/source-migration-v1` materialized a source-topology
checkpoint at:

`5fa292903f74396a07b9a4eeb21540a76ddbe83c`

Commit:

`refactor(architecture): migrate repository source topology`

V1 changed physical source topology while deliberately preserving current runtime
service identity, business semantics and Go module identities.

Major source zones became:

```text
api/
internal/
infrastructure/
proto/
tools/
integration-test/
docs/
```

V1 explicitly did not claim that large areas such as `internal/user`,
`internal/property` or `internal/planning` were final bounded contexts.

Hosted source-integrity / Redis-ownership / observability-error workflows passed
on that exact V1 SHA, but that fact alone was not treated as a new full Final
Acceptance closure.

## H5 — Architecture Rebuild V2

Branch `architecture/rebuild-v2` was created directly from
`5fa292903f74396a07b9a4eeb21540a76ddbe83c`.

Purpose:

- treat old BDSPro as real evidence/reference;
- start from first principles as if a senior engineer were establishing a new repo;
- build a reusable architecture/engineering baseline;
- make the branch a real Go learning corpus tied to production problems;
- redesign database only after business facts/invariants/lifecycle are proven;
- handle TQD/PostGIS as a specialized workload after the baseline is stable;
- checkpoint decisions only after analysis + implementation + proof meet the gate.

Current phase at creation: `R0 — Workspace & Architecture Charter`.
