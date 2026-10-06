Continue BDSPro Canonical Database from durable Git state, not conversational memory.

Repository:
`datvtph41107/bdspro-backend`

Canonical database architecture branch:
`architecture/canonical-database-final`

The database redesign has completed its architecture/design acceptance from v0 → v1 → v2 → final. Do not restart the schema design from scratch and do not ask me to reconstruct the previous chat unless live Git proves the durable package insufficient.

## Authority order

Use this order whenever sources disagree:

1. immutable/live Git source and migrations
2. completed exact-SHA tests/runtime/migration proof
3. the canonical database package on `architecture/canonical-database-final`
4. older architecture/history documents
5. conversational memory

Never let stale prose override current source or completed proof.

## Read these files first, in this exact order

1. `documents/canonical-database/README.md`
2. `documents/canonical-database/V0-SEMANTIC-MODEL.md`
3. `documents/canonical-database/V1-CANONICAL-CORE.dbml`
4. `documents/canonical-database/V2-PHYSICAL-POSTGRESQL.sql`
5. `documents/canonical-database/SERVICE-TABLE-INVENTORY.md`
6. `documents/canonical-database/MIGRATION-AND-PROOF-PLAN.md`
7. `documents/canonical-database/FINAL-ACCEPTANCE.md`
8. this `CONTINUATION-PROMPT.md`

Then reconcile the package against the current branch head and current `main` before changing source.

## Current accepted phase boundary

The database architecture/design is finished.

Accepted states:

- v0 semantic model: ACCEPTED
- v1 logical DBML: ACCEPTED
- v2 physical PostgreSQL reference: ACCEPTED AS REFERENCE
- final architecture acceptance: ACCEPTED TARGET
- production migration/cutover: NOT YET EXECUTED

Do not claim production schema migration is finished merely because the target DDL exists.

The next work is implementation/proof:

`OLD → EXPAND → BACKFILL → COMPARE/PROVE → SWITCH AUTHORITY → CONTRACT`

## Core accepted semantics

### Identity

- Party is business identity root.
- Person and Organization are Party specializations.
- Person != Profile != Account != runtime Actor.
- Profile is presentation, never authorization identity.
- Account is local digital identity of Person.
- authenticator/session/account are separate facts.
- Actor is runtime context, not a table.

### Organization

- Organization canonical identity = `organizations.party_id`.
- current core descriptive fields: name, optional tax_code, phone, email, website, description, archived_at.
- legacy code/type/status/address/owner scalar/verification/warning are not canonical core.
- Organization address remains deliberately DEFERRED because its business meaning is unresolved.
- Invitation != Membership.
- Membership is an effective Person↔Organization episode.
- rejoin = new Membership.
- MembershipBlock is temporary authority restriction while Membership remains current.
- one current Organization RoleAssignment per current Membership under the accepted current model.
- Organization Role != Ownership.
- exactly one current usable OrganizationOwnership per non-archived Organization.
- Ownership points to OrganizationMembership.
- creator/founder is provenance, not current authority.
- Branch is a real operational subdivision.
- Department/Store are NOT proven to be Branch subtypes and remain deferred.
- Branch assignment points to OrganizationMembership, not user/profile.
- Branch manager is current same-Organization Membership.
- Branch assignment/manager relationship is not automatic authorization.

### Group

- Group is an ad-hoc collaboration context, not Organization subtype/unit.
- GroupMembership = Person↔Group episode.
- exactly one current GroupLeadership for non-disbanded Group.
- creator becomes initial member/leader but created_by is not current authority.
- generic group role_id/shared role-catalog coupling is rejected.
- temporary Group pause semantics deferred.

### Deal

- Deal has exactly one primary context: Organization XOR Group.
- generic owner_type + owner_id is rejected.
- personal/member-owned Deal is not currently earned.
- Organization Deal may optionally have one Branch placement; Branch must belong to same Organization.
- Branch placement != permission.
- Deal kind = BROKERAGE / INVESTMENT / JOINT.
- current processing is represented by no terminal outcome.
- terminal outcome = COMPLETED or CANCELED with ended_at.
- Invitation != Participation.
- Deal Invitation target = Person.
- Deal Participation = Person↔Deal episode; max one current participation per Deal+Person.
- ordinary "member" requires no separate role: Participation itself means membership.
- Customer and Partner are independent Deal-scoped business relationships and may overlap.
- Admin is Deal authority capacity, separate from Customer/Partner.
- legacy RoleKeyDealOwner is primary Deal governance/responsibility, accepted as Deal Lead semantics.
- each open Deal has exactly one current Lead.
- creator becomes initial Participation + Lead.
- generic role update cannot create/destroy Lead.
- legacy is_owner is rejected duplicate truth.
- charge_person_id is not Lead; live create flow seeds charge person as invited ADMIN. Target uses Invitation admin term and AdministratorAssignment after acceptance.
- role_key/status/member_type blob is rejected.
- Commission, InvestmentCommitment, Investment and Cost are separate financial facts.
- all money carries currency.
- done_investment/confirmed duplicate flags are not canonical truth.
- generic attach_documents(owner_id,doc_owner) is rejected; use explicit Deal/Investment document relations.

### Payment

- existing Order / Attempt / Settlement / Fulfillment / CommandEffect / Outbox separation is retained.
- commercial subject converges from profile/organization type+id to logical Party reference.
- personal Wallet targets Person.
- BankAccount is Payment-owned and targets Party.
- no cross-service PostgreSQL foreign keys.

### Other bounded contexts

Retain their own schema/write authority:

- bdspro-service: Property
- tqd-service: planning/quota/report/POI
- crm-service: CRM/SEO/support
- notification-service: notification delivery/history
- file-service: file metadata/access
- hub-service: platform lookup/help/config

Do not pull those tables into a monolithic schema just because UI/business flows connect them.

## Important live-source failures already proven

Do not reintroduce these:

- organization_members mixed invited/active/suspended/removed.
- duplicate Organization owner truth existed in scalar owner + member/role/fallback.
- legacy Branch type bucket mixed branch/department/store.
- Branch member/manager used global user/profile IDs instead of OrganizationMembership.
- Branch manager used magic permission bypass.
- Group leader check is currently commented/broken.
- User IAM resolves Group/Deal role IDs through a shared role repository, leaking scope.
- Deal owner_id/owner_type duplicates context bridge tables.
- Deal generic owner enum includes modes not coherently implemented.
- Deal creation dedupes member/customer/partner/admin by person ID and therefore loses overlapping relationships.
- Customer/Partner labels are swapped in one legacy RoleKey map.
- creator gets Deal OWNER; charge person gets invited ADMIN.
- deal_members.is_owner is not coherently maintained.
- Deal invitation accept target check was commented out.
- resend mutates a terminal invitation back to invited.
- withdraw/remove uses the same deal_members row as both invitation and participation.
- investment state is duplicated between investment rows and deal-member booleans.
- generic attach-document owner typing obscures endpoints.

## Deferred facts

Do not invent these without new business/source evidence:

- Organization registered/head-office/mailing address semantics
- Department model
- Store model
- max number of Branches per OrganizationMembership
- temporary Branch disable/reactivate model
- temporary Group pause model
- reusable Group Role system
- additional reusable Deal roles beyond accepted Lead/Admin
- Customer/Partner attribute/lifecycle expansion
- commission/commitment historical versioning
- jurisdiction-specific tax-code normalization

## Next authorized execution gate

Start implementation with the first dependency-safe migration slice:

### Phase A — Party / Person / Account / Profile

Before source mutation:

1. inspect current `main` and architecture branch heads;
2. inspect user-service identity migrations/models/write paths;
3. define additive migration only;
4. define deterministic legacy Profile/Auth → Person/Account mapping;
5. write invariant/backfill/reconciliation tests;
6. do not delete or repurpose legacy columns;
7. produce exact-SHA proof;
8. only then proceed to Organization Phase B.

If the user instead asks for a final review first, review the DBML/DDL against current live source and correct the architecture branch before implementation.

## Working discipline

- Vietnamese explanation.
- Explain from business fact → invariant → transaction → storage.
- Use concrete counterexamples when rejecting a schema.
- Keep one mutable truth owner.
- Preserve old data until migration proof says it can be contracted.
- Do not add generic type/id, owner, status, soft-delete, actor, resource, or transaction abstractions to solve uncertainty.
- Cross-service references are validated through authoritative APIs/events/projections, not cross-database joins.
- Use the narrowest correct contextual identity.
- If live source contradicts an accepted design fact, show the evidence and reopen only that fact.
- If no contradiction exists, do not re-litigate accepted GREEN facts.

When I say **“Tiếp tục”**, continue from the exact current execution checkpoint without asking me to reconstruct prior context.
