# BDSPro Canonical Database — Final Acceptance Report

Date: 2026-10-06
Architecture status: ACCEPTED
Runtime migration status: NOT STARTED
Source baseline reviewed: main @ b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd

## Executive decision

The v0 research phase is closed. The target database architecture is accepted through v1 relational design and v2 operational hardening. The next phase is implementation by bounded cutover waves, not more unconstrained schema brainstorming.

Acceptance means:

- canonical business facts have explicit durable owners;
- the major duplicated/generic truth patterns have a defined replacement;
- system-wide logical relationships are represented in CANONICAL-DATABASE-FINAL.dbml;
- invariants and transaction/concurrency rules are explicit;
- every destructive legacy decision is gated by backfill/compare proof;
- operational projections, audit, inbox/outbox and worker state are not confused with core truth.

It does NOT mean production tables have been altered by this documentation change.

## v0 acceptance

| Question | Result |
| --- | --- |
| Party / Person / Profile / Account separated? | PASS |
| Runtime Actor kept out of durable identity tables? | PASS |
| Organization Invitation separated from Membership? | PASS |
| Membership temporary block separated from termination? | PASS |
| Role/Permission separated from final authorization decision? | PASS |
| Organization ownership separated from Role/creator? | PASS |
| Branch assignment/manager endpoints narrowed to Membership? | PASS |
| Branch/Department/Store legacy conflation rejected? | PASS |
| Group creator shortcut rejected as authority? | PASS |
| Deal context separated from Deal participant governance? | PASS |
| Deal Invitation separated from Participation? | PASS |
| Deal Member/Customer/Partner/Admin/Lead dimensions separated? | PASS |
| Property / Listing / Asset separated? | PASS |
| Generic core owner/type pattern rejected? | PASS |
| Current truth / history / audit / integration effects separated? | PASS |

## v1 relational acceptance

### Identity / IAM

Accepted: Party → Person; Profile is one-to-one presentation; Account is digital identity. Authentication credentials/sessions are Account-scoped. Platform roles are independent from Organization roles.

Legacy fields such as user_profile.role_key/role_type/plan_id are not allowed to remain hidden authority after cutover.

### Organization

Accepted source-of-truth chain:

Organization → Membership → RoleAssignment → Role → Permission.

Current ownership is Organization → Membership. Temporary MembershipBlock removes effective authority without ending participation. BranchAssignment and manager capacity also target Membership.

Exactly one current owner for a non-archived Organization remains a hard invariant.

### Group

Accepted: GroupMembership is current participation; GroupLead is current governance. CreatedBy is audit provenance only.

### Deal

Accepted primary context is Organization XOR Group. Branch is secondary Organization scope.

Accepted participant model:

Person → DealParticipation, with independent Customer, Partner, Admin and Lead facts.

Legacy MEMBER role is removed because Participation already means membership. Legacy OWNER role is promoted into explicit DealLead. Legacy is_owner is discarded as shadow truth.

Commission and Investment are separate economic facts. Money uses minor units + currency.

### Property / Listing / Asset

Accepted:

- Property = identifiable real estate object.
- Listing = market offer to sell/rent a Property.
- Asset = owned business asset tied to Property.
- Asset ownership points to Party.
- Listing business context is explicit and is not called legal ownership.

Parallel Property generations must converge to one authority.

### Commercial / Payment

Current versioned catalog, checkout snapshot, payment attempt, provider settlement, fulfillment and outbox concepts are accepted because they already model separate failure/authority boundaries well.

Target improvement is mainly Party subject identity and Wallet ledger truth.

### Planning / TQD

Planning and report/quota schemas are already relatively mature. Final preserves immutable usage events, quota pool state, report jobs with claim fencing, planning source data and legal references.

Hub owns platform administrative geography; TQD owns planning jurisdiction/geometry.

### Chat / Social / CRM

Accepted convergence removes duplicate Chat memberships and generic Social owner/target patterns. CRM customers table is recognized as Opportunity, not Person identity. SEO remains the deliberately simplified page-publishing system established by migration 000024.

### File / Notification

File durable identity is storage metadata, never physical machine path. Notification owns delivery effects only; foreign-domain audit returns to the domain that understands the event.

## v2 operational acceptance

| Concern | Accepted mechanism |
| --- | --- |
| Retry of command | stable command key + fingerprint |
| Duplicate provider callback | provider transaction uniqueness / inbox event ID |
| Distributed event publication | local transaction + outbox |
| Consumer duplicate | inbox/effect unique identity |
| Worker concurrency | database claim + lease + monotonic claim_version |
| Governance race | stable aggregate row lock + fresh invariant re-read |
| Historical preservation | terminal timestamps/immutable event facts; no resurrection |
| Search authorization | filter at owner before projection/response |
| Cross-service changes | local ACID + idempotent distributed effects |
| Cache loss | reconstruct from durable owner; cache never authority |

## Deliberate absences from the Final core

These are not omissions:

- Generic organization type business/team/other.
- Generic organization code as identity.
- Generic organization address whose business meaning is unknown.
- Department and Store as automatic Branch subtypes.
- Generic owner_id/owner_type on Deal, Listing, Asset, Post, Notification.
- Generic resource/type authorization table.
- Generic active/inactive/removed status for Membership.
- Generic Deal participant role combining Customer/Partner/Admin/Lead.
- Customer or Partner as global Person subtype.
- Creator/created_by as current authority.
- Mutable wallet balance as sole financial truth.

Each may be introduced later only through a new business fact with its own evidence and invariant.

## Risks that remain implementation risks, not architecture blockers

1. Legacy data may contain contradictory owner/context rows.
2. RoleKey Customer/Partner display labels are swapped; migration must rely on numeric/source behavior plus history.
3. Some auth/profile and social leaf tables still use Profile IDs and require deterministic Person mapping.
4. Property has two historical schema generations; mapping coverage must be proven.
5. Hub and TQD contain overlapping administrative geography copies.
6. Legacy Group leader authorization is effectively bypassed because its CreatedBy check is commented.
7. Some service schemas originated from GORM AutoMigrate and need schema snapshot comparison before migration.

These risks are exactly why contraction is forbidden before compare proof.

## Acceptance gates before implementation can be called production-complete

### Gate A — Schema proof

- target migrations generated from accepted owner slices;
- migration contract tests pass on an empty database;
- migrations apply on representative legacy snapshot;
- all v2 constraints are present;
- no serving process calls AutoMigrate as schema authority.

### Gate B — Data proof

- every legacy authoritative row maps or has an explicit reconciliation reason;
- zero unexplained ownership/context conflicts;
- one-current invariants hold after backfill;
- monetary values/currencies reconcile;
- provider/payment evidence hashes/IDs remain unchanged.

### Gate C — Behavior proof

- old/new authorization decisions compared for intended compatible cases;
- known legacy security holes are intentionally denied by target;
- retries produce one durable effect;
- concurrent governance tests preserve invariants;
- worker crash/reclaim tests prove fencing.

### Gate D — Authority cutover

- reads use target;
- writes use target only;
- old writers blocked by tests/permissions;
- outbox/inbox delivery monitored;
- rollback procedure verified.

### Gate E — Contract

- compatibility views/endpoints retired;
- legacy columns/tables dropped only after observation window;
- docs and continuation state updated to exact migration SHA.

## Final verdict

ARCHITECTURE: PASS.

V0 SEMANTICS: CLOSED.

V1 RELATIONAL TARGET: CLOSED.

V2 OPERATIONAL MODEL: CLOSED AS DESIGN.

PRODUCTION MIGRATION: NOT EXECUTED; NEXT AUTHORIZED PHASE.

The next work must start from MIGRATION-CUTOVER-PLAN.md Wave 0/1 and live repository reconciliation. Do not reopen a closed semantic decision merely because a legacy table has a convenient shape; reopen only when new business evidence invalidates an invariant.