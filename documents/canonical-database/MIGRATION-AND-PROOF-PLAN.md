# Migration and Proof Plan

Target: migrate from the current live schemas to the v1/v2 canonical model without losing evidence, breaking rollback, or allowing two mutable authorities to survive indefinitely.

## 1. Governing migration rule

Every canonical replacement follows:

```text
OLD
→ EXPAND
→ BACKFILL
→ COMPARE / PROVE
→ SWITCH AUTHORITY
→ CONTRACT
```

"Deploy new endpoint" is not sufficient. The authority switch is the real migration boundary.

## 2. Global safety rules

1. Never drop or repurpose a legacy column before profiling its live values.
2. Never rewrite historical rows merely to make the new model look clean.
3. Unknown or contradictory rows go to a reconciliation report/quarantine path.
4. During dual-write, define one side as authoritative and the other as compatibility projection.
5. Every backfill is idempotent.
6. Every batch has stable ordering/checkpoints.
7. Every authority switch has:
   - exact source SHA,
   - migration version,
   - row-count/checksum proof,
   - invariant queries,
   - rollback boundary.
8. Do not use one large distributed transaction across service databases.
9. Cross-service propagation uses outbox/idempotent consumer/reconciliation.
10. Contract old schema only after all readers stop treating it as truth.

---

# 3. Phase A — Identity foundation

## 3.1 Expand

Create:

- parties
- persons
- profiles
- accounts
- credential/session tables

Do not remove `user_profile`, `auth_method`, `user_info`, or `user_session` yet.

## 3.2 Backfill identity map

Build a deterministic mapping:

```text
legacy profile_id
→ person.party_id
→ account.id
```

Required checks:

- one legacy Profile never maps to two Persons
- one live auth identity never maps to two open Accounts
- federated provider subject uniqueness
- passkey credential uniqueness
- open Account partial uniqueness

## 3.3 Compare

For every active account/profile:

- login resolves same Person
- current profile display is equivalent
- session revocation behavior is preserved
- authorization context resolves the same business Person

## 3.4 Switch

Runtime principal resolution becomes:

```text
credential
→ Account
→ Person
```

Profile remains presentation only.

## 3.5 Contract

Only then retire authority reads from:

- `auth_method.user_id`
- `user_info.role_id`
- Profile role/status identity shortcuts

---

# 4. Phase B — Organization governance

## 4.1 Expand canonical facts

Create:

- organization_invitations
- organization_memberships
- organization_membership_blocks
- organization_roles
- organization_role_permissions
- organization_role_assignments
- organization_ownerships
- organization_branches
- organization_branch_assignments

Existing `organizations` may be expanded/mapped to the new Party-backed identity rather than physically renamed in the first deployment.

## 4.2 Profile legacy Organization data

Required queries/reports:

### Organization code

Profile:

- distinct count
- mutation evidence if audit exists
- references from other services
- archived code reuse

Do not map legacy code to canonical identity.

### Tax code

Report:

- NULL/blank
- duplicate exact values
- duplicate canonicalized values
- duplicate active+archived values

Conflicts are manually reconciled before enabling final unique index.

### Owner

Compare all available owner truths:

- `organizations.owner_profile_id`
- legacy Owner role/member row
- organization-service `OwnerId`
- admin/full-access fallbacks

Generate:

```text
organization_id
legacy_scalar_owner
legacy_member_owner
candidate_person
candidate_membership
classification
```

No silent winner when values disagree.

## 4.3 Backfill Membership episodes

Legacy `organization_members.status` classification:

- active → current Membership
- suspended → current Membership + current MembershipBlock
- removed → ended Membership
- invited → OrganizationInvitation only, not Membership

If legacy rows cannot distinguish invitation from effective membership, use source timestamps/events where available and report ambiguity.

## 4.4 Backfill roles

Legacy current role behavior is "replace old role with one current role."

For every current Membership:

- resolve exactly one Organization Role
- create one current RoleAssignment

Rows with zero or multiple effective roles fail acceptance until reconciled.

## 4.5 Backfill ownership

For every non-archived Organization:

- exactly one current Membership must be selected
- create exactly one OrganizationOwnership

Proof query must return zero rows for:

- missing owner
- two owners
- owner from different Organization
- ended owner Membership
- blocked owner Membership

## 4.6 Branch migration

Legacy Branch member mapping:

```text
organization_branch_members.user_id
→ Person
→ current OrganizationMembership
→ organization_branch_assignments.membership_id
```

Reject/quarantine rows where:

- Person has no current Membership in Branch Organization
- duplicate branch/person relation exists
- Branch is already terminally removed
- manager cannot resolve to current Membership

Legacy `type=department/store` is not automatically converted to canonical Branch.

Create a classification report:

- clear operational Branch
- Department candidate
- Store candidate
- ambiguous

Only clear Branch rows enter canonical `organization_branches`.

## 4.7 Authority switch

Organization create becomes one local transaction:

```text
Party
+ Organization
+ founder Membership
+ initial RoleAssignment
+ Ownership
```

Invitation acceptance becomes one local transaction:

```text
lock Invitation
→ validate
→ create Membership
→ create RoleAssignment
→ mark Invitation ACCEPTED
```

Membership termination atomically:

- ends current RoleAssignment
- ends current MembershipBlock
- removes Branch assignments
- clears/replaces Branch manager relationship if applicable
- rejects if owner unless Ownership transfer is part of the same transaction
- ends Membership

## 4.8 Contract

After proof:

- stop writes to `organization_members`
- stop writes to `owner_profile_id`
- stop role/profile authorization fallback
- stop organization-service duplicate Organization authority
- retain compatibility projections only while readers remain

---

# 5. Phase C — Group

## 5.1 Expand

Create:

- groups target columns
- group_memberships
- group_leaderships
- explicit Group↔File and Group↔Chat integration links

## 5.2 Backfill

Legacy `group_members.user_id` resolves to Person.

Legacy creator:

- create creator GroupMembership if absent
- seed GroupLeadership from creator only when source evidence has no stronger current leader fact

Because current `IsLeaderGroup` enforcement is commented out, this backfill is explicitly classified as a migration assumption and must be reviewed against production data/activity.

## 5.3 Lifecycle mapping

- disbanded → `disbanded_at` if timestamp evidence exists
- paused → do not invent a permanent target flag unless real behavior is confirmed
- soft deleted → classify whether terminal disbandment or accidental/technical deletion

## 5.4 Contract

Remove current authority dependence on:

- `created_by`
- generic `group_members.role`
- `group_members.role_id` pointing into Organization/Deal role catalog

---

# 6. Phase D — Deal

This is the highest-risk semantic migration.

## 6.1 Profile owner/context data

Before writing any target Deal context:

```sql
SELECT owner_type, COUNT(*)
FROM deals
GROUP BY owner_type
ORDER BY owner_type;
```

Also report:

- owner_id = 0
- unknown owner_type
- Organization owner whose Organization does not exist
- Group owner whose Group does not exist
- MEMBER/USER owner rows
- conflicts with `deal_of_organization`
- conflicts with `deal_of_group`
- conflicts between `deal_of_organization.branch_id` and `deal_of_branch`

Classification:

- Organization context
- Group context
- possible personal/member legacy case requiring business review
- corrupt/unknown

No MEMBER/USER row is silently converted to Person-owned Deal.

## 6.2 Expand Deal target

Create:

- canonical Deal context columns
- deal_invitations
- deal_participations
- deal_customer_relationships
- deal_partner_relationships
- deal_administrator_assignments
- deal_leaderships
- deal_participant_commissions
- deal_investment_commitments
- canonical deal_investments
- deal_milestones
- deal_cost_types/deal_costs
- deal_internal_notes
- explicit document relations

## 6.3 Split `deal_members`

For each legacy row classify by lifecycle.

### INVITED

Creates DealInvitation only.

Map legacy role/capacity intent:

- MEMBER → ordinary invitation
- CUSTOMER → `as_customer=true`
- PARTNER → `as_partner=true`
- ADMIN → `as_administrator=true`
- OWNER → reconciliation; do not create lead from an unresolved invite

### ACCEPTED

Creates effective DealParticipation.

If there is credible invitation provenance:

- create/migrate Invitation as ACCEPTED
- link originating_invitation_id

Otherwise:

- create direct Participation with NULL origin

Then materialize independent facts:

- CUSTOMER → customer relationship
- PARTNER → partner relationship
- ADMIN → administrator assignment
- OWNER → lead candidate
- MEMBER → no extra row; Participation already means member

### REJECTED

Invitation terminal outcome DECLINED. No Participation.

### WITHDRAWN

Do not rewrite Invitation.

Reconstruct:

- accepted Invitation if evidence exists
- ended Participation

Legacy withdrawal/remove reason and actor stay audit/history unless a future business invariant depends on them.

## 6.4 Detect role-key information loss

Legacy create deduped by `member_id`, so a Person appearing as both Partner and Admin may have only one surviving role.

Reconcile using:

- request/event history
- Deal history
- CRM/contract data
- product distribution data
- commission records
- investment records

Do not infer "only one capacity" from the single surviving role_key.

Also account for the legacy Customer/Partner display-label swap.

## 6.5 Deal Lead

For each non-terminal Deal:

- identify creator/legacy OWNER candidates
- require exactly one current lead candidate
- create DealLeadership

Conflicts:

- no owner
- multiple OWNER rows
- owner row ended/withdrawn
- `is_owner` disagreement

must be explicit reconciliation findings.

Legacy `is_owner` is not trusted as primary truth because normal creation paths did not maintain it coherently.

## 6.6 charge_person

Legacy charge person was seeded as invited ADMIN.

Migration behavior:

- if accepted/current participation exists → AdministratorAssignment
- if pending invitation exists → Invitation `as_administrator=true`
- do not create a second Lead
- remove scalar `charge_person_id` authority after switch

## 6.7 Financial split

### amount_commit

Backfill to `deal_investment_commitments`.

Require a currency source. If historical currency cannot be proven, preserve the row in reconciliation state instead of inventing VND silently.

### commission

Backfill to `deal_participant_commissions`.

Validate:

- percent 0..100
- fixed amount has currency
- values are not negative

Legacy commission type/value contradictions are reported.

### investments

Backfill to canonical Investment rows.

Map:

- investor → exact DealParticipation episode
- amount + currency
- transfer timestamp
- proof file
- outcome from current workflow state

Reject duplicate truth:

- `confirmed`
- `deal_members.done_investment`

These become comparison signals, not canonical fields.

### costs

Backfill amount+currency; no currency guessing without an explicit migration policy approved for the dataset.

## 6.8 Deal lifecycle

Map:

- PROCESSING → outcome NULL
- COMPLETED → outcome COMPLETED + ended_at
- CANCELED → outcome CANCELED + ended_at + cancel reason

If terminal timestamp cannot be reconstructed, preserve source timestamp evidence and flag the migration assumption.

## 6.9 Deal authority switch

Stop authorization through:

- `role_key`
- `role_id` from shared IAM role catalog
- `is_owner`
- generic owner type/id
- Branch membership magic allow

New policy consumes:

- Deal current Participation
- Lead/Admin capacity
- Customer/Partner business relationships where relevant
- Organization/Group authority needed for the intent
- Branch/resource scope when explicitly required
- Deal lifecycle
- system Permission

---

# 7. Phase E — Payment and commercial subject identity

Current commerce tables remain authoritative during migration.

Expand Party-based subject columns beside legacy `subject_kind/subject_id`.

Backfill mapping:

- profile → Person Party
- organization → Organization Party

Proof:

- every active Order maps to exactly one Party
- command uniqueness preserved after key change
- no reference collision
- amount/currency unchanged
- settlement/fulfillment links unchanged

Switch writes to Party reference, then remove legacy subject pair.

For wallets:

- resolve legacy user/profile ID to Person
- verify one wallet per Person
- retain balance invariants

For Organization-service bank accounts:

- map owner to Party
- insert Payment-owned BankAccount
- switch Deal to Payment BankAccount ID
- stop writing generic `owner_id/owner_type` bank accounts

---

# 8. Phase F — Cross-service event boundaries

When authority moves across a service boundary:

1. source commits local business truth
2. source commits outbox row in same local transaction
3. publisher emits event
4. consumer inbox/idempotency guards duplicates
5. consumer updates local projection
6. reconciliation compares source and projection

Required events include at least:

- Person/Account changed where consumers need projection
- Organization created/archived
- Membership started/ended/blocked/unblocked
- Organization Role changed
- Branch created/retired/manager changed
- Group disbanded
- Deal created/completed/canceled
- Deal Participation started/ended
- Deal Lead transferred
- BankAccount changed when Deal consumers cache display data
- Payment completed/funds confirmed

Events are integration facts, not a second source of business truth.

---

# 9. Required invariant test matrix

## Identity

- two open Accounts for one Person → rejected
- duplicate provider+subject → rejected
- duplicate passkey credential → rejected
- two current AccountBlocks → rejected

## Organization

- two current Memberships same Org+Person → rejected
- RoleAssignment from another Org → rejected
- two current RoleAssignments → rejected
- two current MembershipBlocks → rejected
- Ownership to ended Membership → rejected
- Ownership to blocked Membership → rejected
- Ownership from another Organization → rejected
- owner leave/remove without transfer → rejected
- archive with stale concurrent governance write → one wins, invariant preserved
- Branch manager from another Org → rejected
- Branch assignment from another Org → rejected
- assign to retired Branch → rejected

## Group

- two current memberships same Group+Person → rejected
- leadership from another Group → rejected
- leader membership termination without transfer → rejected
- disbanded Group receives new member/Deal → rejected

## Deal

- Organization and Group both set → rejected
- neither context set → rejected
- Branch on Group Deal → rejected
- Branch belongs to another Organization → rejected by authoritative write validation
- two unresolved invitations same Deal+Person → rejected
- accept Invitation as wrong Person → rejected
- accept twice concurrently → only one Participation
- two current Participations same Deal+Person → rejected
- Lead points to another Deal Participation → rejected
- second current Lead → rejected
- end Lead Participation without transfer → rejected
- ADMIN relationship on ended Participation → rejected
- Customer+Partner simultaneously → allowed
- Customer+Admin simultaneously → allowed
- commission percent >100 → rejected
- fixed commission without currency → rejected
- negative investment/cost → rejected
- terminal Deal receives forbidden mutation → rejected

## Payment

Retain existing proof plus:

- one wallet per Person
- bank account duplicate bank+account number rejected
- one Order command per Party/product/command key
- provider transaction idempotency
- payment attempt command idempotency
- fulfillment single row/order
- outbox event uniqueness

---

# 10. Concurrency proof

Use deterministic race tests, not only happy-path unit tests.

Required concurrent pairs:

- AcceptOrganizationInvitation × AcceptOrganizationInvitation
- JoinOrganization × JoinOrganization same Person
- TransferOwnership × OwnerLeave
- TransferOwnership × BlockOwner
- ArchiveOrganization × AddMembership
- RetireBranch × AssignMembershipToBranch
- ChangeBranchManager × EndMembership
- AcceptDealInvitation × AcceptDealInvitation
- TransferDealLead × EndLeadParticipation
- CompleteDeal × AddInvestment
- Payment settlement duplicate delivery
- fulfillment lease claim by two workers

Proof condition is final committed invariant, not which request happens to win.

---

# 11. Observability and reconciliation

Every migration slice exposes:

- legacy row count
- target row count
- classified row count
- ambiguous row count
- invalid row count
- checksum/key comparison
- last processed key
- exact source SHA
- migration build/version

No slice reaches CONTRACT while ambiguous/invalid rows are silently ignored.

---

# 12. Rollback boundaries

Before authority switch:

- rollback may disable target writes and continue old authority.

After authority switch but before contract:

- old structures are read-only compatibility projections.
- rollback requires replay/reconciliation from target truth, never blindly making both sides writable again.

After contract:

- rollback is a forward-fix or restore from proven backup/event history, not resurrection of removed dual-write architecture.

---

# 13. Final migration order

Recommended order:

1. Identity / Party / Person / Account
2. Organization identity
3. Organization Invitation / Membership / Block
4. Role / Permission / RoleAssignment
5. Organization Ownership
6. Branch / Assignment / Manager
7. Group / GroupMembership / Leadership
8. Deal context and lifecycle
9. Deal Invitation / Participation
10. Deal Customer/Partner/Admin/Lead
11. Deal financial facts
12. Deal Product/Milestone/Cost/Document
13. Payment Party subject + BankAccount authority
14. consumer projections/events
15. compatibility contraction

This order follows dependency direction and minimizes temporary semantic adapters.
