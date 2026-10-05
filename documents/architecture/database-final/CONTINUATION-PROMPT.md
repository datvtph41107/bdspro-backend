# Continuation Prompt — Canonical Database Final → Implementation

Continue BDSPro from durable state, not conversational memory.

Canonical repository:
datvtph41107/bdspro-backend

Architecture branch:
docs/canonical-database-final-20261006

Baseline live source used for the architecture package:
main @ b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd

First recover the database architecture package in this exact order:

1. documents/architecture/database-final/README.md
2. documents/architecture/database-final/DATABASE-FINAL-ACCEPTANCE.md
3. documents/architecture/database-final/CANONICAL-DATABASE-FINAL.dbml
4. documents/architecture/database-final/DATABASE-INVARIANTS.md
5. documents/architecture/database-final/LEGACY-DISPOSITION.md
6. documents/architecture/database-final/MIGRATION-CUTOVER-PLAN.md
7. documents/architecture/database-final/DATABASE-VERSION-EVOLUTION.md

Then read the earlier architecture foundations only where a decision needs provenance:

8. documents/architecture/ARCHITECTURE-NAMING-OWNERSHIP-MINDSET.md
9. documents/architecture/AUTHENTICATION-AUTHORIZATION-ORGANIZATION-FOUNDATION.md
10. documents/SERVICE-MAP.md
11. documents/QHPRO-COMMERCIAL-OPERATIONS.md

## Authority rule

Immutable live Git/source and completed exact-SHA runtime/CI proof win over prose. The Final DB package is the target architecture owner, but before source mutation reconcile it against the current main HEAD because main may have advanced after the baseline SHA.

Do not infer current implementation progress from this prompt. Inspect:

- current main SHA;
- architecture branch HEAD / merged status;
- relevant service migration heads;
- worktree/branch role if local workspace information is available;
- CI/migration-contract evidence for the exact candidate SHA.

Preserve interrupted worktrees, candidates, reports and evidence. Do not reset, clean, rebase, amend or rewrite unrelated work.

## Closed decisions that must not be casually reopened

- Party / Person / Profile / Account are distinct.
- Actor is runtime, not table.
- Invitation != Membership/Participation.
- Creator != authority.
- Organization current owner is explicit Membership relation.
- MembershipBlock is separate from Membership termination.
- Organization RoleAssignment targets Membership.
- Branch assignment/manager target OrganizationMembership.
- Department/Store are not automatically Branch types.
- Group creator is not leader authority; GroupLead is explicit.
- Deal primary context is Organization XOR Group; Branch is secondary.
- Deal participant identity is Person.
- Deal Customer/Partner relationship can overlap; Admin/Lead are separate authority axes.
- Deal Member role is redundant with Participation.
- Generic owner_id/type_id is rejected for canonical core.
- Property != Listing != Asset.
- Asset ownership attaches to Party.
- Money uses minor units + currency.
- Cache/projection/audit/event transport do not become canonical business truth.
- Cross-service consistency uses local ACID + outbox/inbox/idempotency, not distributed DB transaction.

## Next authorized phase

Start implementation with MIGRATION-CUTOVER-PLAN.md.

Do NOT begin by generating every target migration at once.

Start Wave 0 evidence freeze and Wave 1 Party/Person/Account only:

1. reconcile current main with the architecture target;
2. inventory exact Profile/Auth consumers;
3. write the smallest expand migration;
4. add migration contract tests;
5. add deterministic backfill/compare tooling;
6. prove on an empty DB and representative legacy DB;
7. commit the slice;
8. only then advance to the next gate.

At every slice explain:

business failure → target fact → invariant → transaction/constraint → migration mapping → proof.

Do not add generic status, owner, type/id, actor_id, soft delete, timestamps or JSON simply because another table has them.

## Final objective

Finish all cutover waves until the production databases implement the accepted Final model, old duplicate authorities are contracted, exact-SHA proof is green, and the durable package is updated to the implementation SHA.