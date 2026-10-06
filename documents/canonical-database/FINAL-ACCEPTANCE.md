# BDSPro Canonical Database — Final Acceptance

## Acceptance scope

This document is the final architecture/design acceptance for the BDSPro database redesign carried from v0 through v1 and v2.

It accepts:

- the semantic ownership model,
- the canonical logical schema,
- the physical PostgreSQL enforcement reference,
- the service/table disposition,
- the migration and proof strategy.

It does **not** claim that the production databases have already been migrated to the target schema.

Production adoption remains a separate execution program governed by:

```text
OLD
→ EXPAND
→ BACKFILL
→ COMPARE / PROVE
→ SWITCH AUTHORITY
→ CONTRACT
```

---

# 1. Acceptance verdict

## v0 — ACCEPTED

The business facts, identities, lifecycle boundaries, authority relationships, negative schema, and service ownership are coherent enough to stop exploratory schema invention.

Accepted backbone:

```text
Party
├── Person
│   ├── Profile
│   └── Account
└── Organization

Organization
├── Invitation
├── Membership
│   ├── MembershipBlock
│   └── RoleAssignment
├── Ownership
└── Branch
    ├── BranchAssignment
    └── Manager capacity

Group
├── GroupMembership
└── GroupLeadership

Deal
├── Organization OR Group context
├── optional Branch placement
├── Invitation
├── Participation
│   ├── Customer relationship
│   ├── Partner relationship
│   ├── Administrator capacity
│   ├── Commission terms
│   └── Investment commitment
├── exactly one current Lead while open
├── Investments
├── Products
├── Milestones
├── Costs
├── Notes
└── Documents
```

## v1 — ACCEPTED

`V1-CANONICAL-CORE.dbml` is accepted as the logical target ER model.

It intentionally includes logical references across bounded contexts for architecture review. Those references are not instructions to create PostgreSQL foreign keys across service databases.

## v2 — ACCEPTED AS REFERENCE DDL

`V2-PHYSICAL-POSTGRESQL.sql` is accepted as the physical enforcement reference.

It demonstrates:

- primary/foreign keys where facts share a database,
- partial unique indexes,
- lifecycle/outcome checks,
- same-context constraint triggers,
- monetary precision,
- local concurrency coordinators,
- explicit cross-service validation boundaries.

It is **not** a one-shot production migration script.

## final — ACCEPTED TARGET / MIGRATION PENDING

The architecture package is accepted as the target.

Production migration is pending execution of the proof plan and cannot be called complete until exact-SHA migration/runtime evidence exists.

---

# 2. Canonical write ownership

| Business fact | Final canonical owner |
|---|---|
| Party | user-service |
| Person | user-service |
| Profile | user-service |
| Account / credentials / session | user-service |
| Permission | user-service |
| Organization | user-service |
| Organization Invitation | user-service |
| Organization Membership / Block | user-service |
| Organization Role / RoleAssignment | user-service |
| Organization Ownership | user-service |
| Organization Branch / assignment / manager | user-service target |
| Group / GroupMembership / Leadership | organization-service |
| Deal and Deal lifecycle | organization-service |
| Deal Invitation / Participation | organization-service |
| Deal Customer / Partner relationships | organization-service |
| Deal Admin / Lead authority | organization-service |
| Deal commission / investment commitment / investment / cost | organization-service |
| Bank / BankAccount / Wallet / Commerce payment | payment-service |
| Product catalog / subscription | user-service |
| Property | bdspro-service |
| Planning / quota / POI / reports | tqd-service |
| CRM / SEO / support | crm-service |
| Notification delivery/history | notification-service |
| File metadata/access | file-service |
| Hub lookup/config/help | hub-service |

No mutable business fact has two intended write owners in the target.

---

# 3. Major semantic failures repaired

## 3.1 Profile/User identity conflation

Legacy failure:

```text
profile_id / user_id
was used as
human identity
+ login identity
+ participant identity
+ authorization subject
```

Repair:

```text
Person
≠ Profile
≠ Account
≠ runtime Actor
```

Business relationships now terminate at the narrowest correct contextual identity.

---

## 3.2 Organization Invitation/Membership/status conflation

Legacy:

```text
organization_members.status =
invited | active | suspended | removed
```

This mixed:

- proposal,
- effective relationship,
- temporary authority restriction,
- terminal relationship end.

Repair:

- OrganizationInvitation
- OrganizationMembership
- MembershipBlock
- Membership.ended_at

---

## 3.3 Organization ownership duplication

Legacy authority existed simultaneously in:

- `organizations.owner_profile_id`,
- owner role/member rows,
- Organization-service owner scalar/fallbacks.

Repair:

```text
OrganizationOwnership
    organization_id
    membership_id
```

Exactly one current owner for each non-archived Organization.

Ownership is not Role, creator, or Profile.

---

## 3.4 Shared role scope leakage

Legacy reused role identifiers/catal​ogs across Organization, Group, Branch and Deal paths.

Repair:

- OrganizationRole is Organization-scoped only.
- Branch Manager is a Branch relationship.
- Group Leadership is a Group relationship.
- Deal Lead/Admin are Deal capacities.
- no polymorphic `role_id` is used to hide scope.

---

## 3.5 Branch polymorphic bucket

Legacy `organization_branches.type` stored:

```text
branch | department | store
```

while all rows shared generic CRUD and quota behavior.

Repair:

- Branch accepted only from evidence of real operational subdivision.
- Department and Store are not silently treated as Branch or OrganizationUnit.
- Branch member endpoint becomes OrganizationMembership.
- Branch manager endpoint becomes OrganizationMembership.
- `is_active/is_deleted/deleted_at` are not copied as canonical lifecycle.

---

## 3.6 Group creator-as-authority

Legacy intended creator as leader, but current leader enforcement is commented/broken.

Repair:

- GroupMembership
- GroupLeadership
- creator is provenance only
- leader transfer is explicit

---

## 3.7 Deal generic owner

Legacy:

```text
deals.owner_type
deals.owner_id
```

coexisted with:

- `deal_of_organization`
- `deal_of_group`
- `deal_of_branch`

Repair:

- Deal has exactly one primary business context: Organization XOR Group.
- Branch is optional secondary operational placement.
- User/Member-owned Deal is not invented without live business evidence.
- duplicate relation tables are retired after migration proof.

---

## 3.8 `deal_members` semantic blob

Legacy one row mixed:

- Invitation
- Participation
- Member/Customer/Partner classification
- Owner/Admin authority
- commission
- investment commitment
- investment completion
- withdrawal/removal history

Repair:

- DealInvitation
- DealParticipation
- Customer relationship
- Partner relationship
- Admin assignment
- Lead relationship
- Commission terms
- Investment commitment
- Investment

No generic `role_key` or participant `status` remains as canonical truth.

---

## 3.9 Deal Owner naming collision

Legacy "owner" meant two different facts:

1. Organization/Group context of Deal.
2. participant designated as "Deal Owner."

Repair:

- first is Deal business context;
- second is Deal Lead/governance capacity.

Creator becomes initial Lead but does not remain authority forever.

---

## 3.10 charge person duplication

Legacy `charge_person_id` looked like a responsible-person scalar but creation code seeds it as an invited Deal ADMIN.

Repair:

- no independent charge-person authority scalar;
- Invitation may offer Admin capacity;
- accepted Participation receives Admin assignment.

---

## 3.11 Financial duplication

Legacy stored financial state across:

- `deal_members.amount_commit`
- commission fields
- `deal_members.done_investment`
- `investments.confirmed`
- investment status
- costs without guaranteed currency

Repair:

- InvestmentCommitment
- ParticipantCommission
- Investment
- DealCost
- all money carries currency
- duplicate summary booleans are not canonical truth

---

## 3.12 Generic document ownership

Legacy:

```text
attach_documents(owner_id, doc_owner)
```

Repair:

- explicit DealDocument
- explicit InvestmentDocument
- File service owns file metadata

The endpoint becomes structurally knowable.

---

## 3.13 Payment subject polymorphism

Legacy commerce uses:

```text
subject_kind = profile | organization
subject_id
```

Repair target:

- both are Party,
- commerce Order references logical Party identity.

Existing commerce Order/Attempt/Settlement/Fulfillment/Outbox separation is retained.

---

# 4. Final canonical invariants

## Identity

- at most one open Account per Person
- provider+subject unique
- passkey credential unique
- at most one current AccountBlock

## Organization

- at most one current Membership per Organization+Person
- exactly one current Organization RoleAssignment per current Membership
- at most one current MembershipBlock
- current Role must belong to Membership Organization
- exactly one current usable Ownership per non-archived Organization
- Branch manager is current Membership of same Organization
- Branch assignment uses current Membership of same Organization
- terminal Membership/Branch cannot remain in current structural authority

## Group

- at most one current Membership per Group+Person
- exactly one current Leader per non-disbanded Group
- Leader is current Membership of same Group

## Deal

- exactly one of Organization/Group context
- Branch only on Organization Deal and same Organization
- at most one unresolved Invitation per Deal+Person
- at most one current Participation per Deal+Person
- Lead is current Participation of same Deal
- one current Lead while Deal is open
- Customer and Partner relationships may coexist
- Admin capacity is independent of Customer/Partner
- terminal Deal cannot accept forbidden current-state mutations
- fixed commission requires currency
- percent commission is bounded
- investment/cost amounts cannot be negative
- Investment references exact Participation episode

## Payment

- one personal Wallet per Person under current model
- unique Bank+Account number
- Order command idempotency
- provider transaction idempotency
- one Fulfillment per Order
- Outbox event identity unique

---

# 5. Transaction boundaries accepted

### Create Organization

One user-service transaction:

```text
Party
+ Organization
+ founder Membership
+ initial RoleAssignment
+ Ownership
```

### Accept Organization Invitation

One transaction:

```text
lock Invitation
→ validate
→ Membership
→ RoleAssignment
→ Invitation ACCEPTED
```

### Transfer Organization Ownership

Organization row is the governance lock coordinator.

### End Organization Membership

One transition coordinates:

- owner invariant
- RoleAssignment end
- MembershipBlock end
- Branch assignments
- Branch manager consequences
- Membership end

### Create Group

One organization-service transaction:

```text
Group
+ founder GroupMembership
+ GroupLeadership
```

External chat creation is an effect/integration concern, not a reason to hold a distributed SQL transaction.

### Create Deal

One organization-service transaction:

```text
Deal
+ creator Participation
+ DealLeadership
+ invitation offers
+ local Deal relations
```

### Accept Deal Invitation

One transaction:

```text
lock Invitation
→ verify target Person
→ create Participation
→ materialize offered Customer/Partner/Admin facts
→ mark ACCEPTED
```

### Transfer Deal Lead

Deal row is the local governance coordinator.

### Payment completion

Retain local durable payment commit + outbox/event propagation. Consumer updates are idempotent.

---

# 6. Deferred decisions — intentionally NOT invented

These are not blockers for the accepted target but remain explicitly unresolved:

1. exact Organization address semantics:
   - registered/legal address?
   - head office?
   - mailing/display address?
2. Department ontology.
3. Store ontology.
4. whether one Membership may belong to multiple Branches as a product invariant.
5. temporary Branch disable/reactivate model.
6. temporary Group pause model.
7. reusable Group Role system beyond explicit Leadership.
8. more granular Deal authority roles beyond Lead/Admin.
9. whether customer/partner relationships later need their own attributes/lifecycle.
10. historical commission/commitment versions if contractual history becomes required.
11. exact tax-code normalization policy by jurisdiction.
12. Organization/Group/Deal address/location expansion if future operations require physical sites.

Rule for all deferred items:

> add a fact only after a real business transition/invariant requires it.

---

# 7. Production-readiness gates

Architecture acceptance does not authorize dropping legacy data.

Production cutover requires all gates below.

## Gate A — schema expand

- additive target migrations reviewed
- no destructive legacy mutation
- rollback route documented

## Gate B — data profile

Required reports include:

- Organization owner disagreements
- tax-code collisions
- Organization member state classification
- Branch/Department/Store classification
- invalid Branch membership/manager endpoints
- Group leader ambiguity
- Deal owner_type population
- Deal context bridge conflicts
- Deal member lifecycle split
- Deal role/capacity collisions
- Customer/Partner label drift
- currency gaps
- investment duplicate-state disagreements

## Gate C — backfill

- deterministic
- idempotent
- checkpointed
- ambiguous rows isolated, never silently skipped

## Gate D — compare/prove

- row counts
- key checksums
- business query equivalence
- authorization equivalence
- negative invariant tests
- concurrent race tests

## Gate E — authority switch

- new write path is the only source of truth
- legacy writes disabled
- projections/events demonstrably converge

## Gate F — exact-SHA runtime proof

- migrations run from exact source SHA
- target service starts against migrated schema
- acceptance/E2E tests pass
- event/outbox consumers pass
- reconciliation reports clean

## Gate G — contract

Only after all readers stop treating old structures as authority may the old columns/tables be dropped.

---

# 8. Required implementation order

```text
1. Party / Person / Account / Profile
2. Organization
3. Organization Invitation / Membership / Block
4. Organization Role / Permission / RoleAssignment
5. Organization Ownership
6. Branch / Assignment / Manager
7. Group / GroupMembership / Leadership
8. Deal context/lifecycle
9. Deal Invitation / Participation
10. Deal Customer/Partner/Admin/Lead
11. Deal financial facts
12. Deal Product/Milestone/Cost/Document
13. Payment Party subject + Payment-owned BankAccount
14. cross-service event/projection reconciliation
15. compatibility contraction
```

---

# 9. Package artifacts

Canonical package:

- `documents/canonical-database/README.md`
- `documents/canonical-database/V0-SEMANTIC-MODEL.md`
- `documents/canonical-database/V1-CANONICAL-CORE.dbml`
- `documents/canonical-database/V2-PHYSICAL-POSTGRESQL.sql`
- `documents/canonical-database/SERVICE-TABLE-INVENTORY.md`
- `documents/canonical-database/MIGRATION-AND-PROOF-PLAN.md`
- `documents/canonical-database/FINAL-ACCEPTANCE.md`
- `documents/canonical-database/CONTINUATION-PROMPT.md`

Working branch:

`architecture/canonical-database-final`

---

# 10. Final acceptance statement

The BDSPro database design is accepted as a **canonical target architecture**.

The redesign has moved the system from:

```text
generic IDs
+ duplicated owners
+ generic statuses
+ profile/user identity ambiguity
+ invitation/effective-relation conflation
+ mixed authority scopes
+ large semantic CRUD rows
```

to:

```text
explicit identity
+ explicit capacity
+ explicit effective relationships
+ explicit current authority
+ narrow lifecycle facts
+ enforceable local invariants
+ service-owned write truth
+ migration/proof boundaries
```

The next phase is no longer "design another database."

The next phase is:

> implement the accepted target slice-by-slice, prove exact migration and runtime behavior, switch authority, then contract the legacy schema.

That distinction is the final acceptance boundary.
