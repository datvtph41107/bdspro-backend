# v0 — Canonical Semantic Model

This document is the v0 closure: business facts and invariants before physical schema.

## 1. Semantic backbone

The canonical reasoning chain is:

```text
BUSINESS INTENT
→ ACTOR / CAPACITY / CONTEXT
→ PRECONDITION / AUTHORITY
→ BUSINESS TRANSITION
→ DURABLE FACTS
→ INVARIANTS
→ TRANSACTION BOUNDARY
→ EFFECTS / EVENTS / PROJECTIONS
→ RELATIONAL MODEL
```

A table is accepted only when a row can be read as one durable business sentence.

---

## 2. Identity model

### Party

A Party is an independently identifiable business participant capable of rights, obligations, capacities, or relationships.

Accepted subtypes in current scope:

- Person
- Organization

`Party` is not Actor, Principal, Account, Profile, Customer, Partner, Member, Owner, or Role.

### Person

A Person is the durable human business identity.

A Person may exist without an open Account. A Person may participate in many Organizations, Groups, and Deals.

### Account

An Account is the durable local digital identity used to authenticate a Person.

Rules:

- Account is not Person.
- Authenticator is not Account.
- Session is not Account.
- Closing an Account does not erase the Person.
- At most one open local Account per Person is the target invariant.
- Re-registration after closure may create a new Account episode if product policy permits.

### Profile

Profile is presentation/descriptive state for a Person.

It may contain display name, avatar, public bio, contact/display URLs and preferences. It is never the authorization endpoint.

Legacy failure repaired:

```text
profile_id / user_id
was used interchangeably as
human identity + login identity + business participant + authorization subject
```

### Authentication facts

Canonical durable authentication facts are separated:

- password credential
- federated provider identity
- passkey credential
- account session
- temporary account block/restriction

OTP/challenge data is operational security state, not Person/Account identity.

---

## 3. Permission and runtime Actor

### Permission

Permission is a stable system-defined authorization primitive, identified by a code.

Permission is **not** a final allow/deny decision.

### Actor

Actor is runtime context, not a table:

```text
request
→ caller
→ authenticated principal
→ account/service identity
→ person/party
→ selected organization/group/deal capacity
→ membership/participation
→ role/capability
→ resource relationships
→ business state
→ risk/freshness/policy
→ ALLOW / DENY
```

No generic `actor_id` foreign key is introduced into all business tables.

---

## 4. Organization

### Organization identity

A canonical Organization row asserts:

> Party O is an Organization with current descriptive business data.

Accepted core properties:

- name
- optional tax code
- phone
- email
- website
- description
- terminal archive timestamp

Rejected from core:

- generic `code`
- generic `type = business/team/other`
- generic `status`
- `owner_profile_id`
- verification/warning fields
- address whose meaning is not known
- created_by/updated_by as authority
- generic soft-delete flags

### Tax code

Tax code is an external business identifier attributed to Organization.

It is not canonical identity. `party_id` remains canonical identity.

If present, one normalized tax code must identify at most one Organization across active and archived history. Archiving does not release identifier uniqueness.

### Organization address

Legacy `address/detail_address` is semantically unresolved.

Possible meanings include registered address, head office, mailing address, primary operating site, or display address. No canonical address fact is created until the business meaning is explicit.

### Organization Invitation

Invitation is a proposal, not Membership.

One invitation targets exactly one of:

- known Person
- email
- phone

Invitation target is immutable.

Lifecycle:

- unresolved: `outcome IS NULL`
- terminal: ACCEPTED / DECLINED / REVOKED
- `resolved_at` exists iff terminal
- expiration is derived from `expires_at`
- resend after a terminal outcome creates a new Invitation

Accepted Invitation materializes exactly one Membership.

### Organization Membership

One row is one effective Person↔Organization participation episode.

Rules:

- `joined_at` means effective participation.
- `ended_at` terminates the episode.
- leave and admin removal are different intents but the same resulting core fact: Membership ends.
- ended Membership is never restored.
- rejoin creates a new Membership.
- at most one current Membership per Organization+Person.
- Invitation is optional origin because founder/direct establishment exists.

### Membership Block

A Block is a temporary authority restriction on a still-current Membership.

Rules:

- Membership stays current.
- Role assignment stays current.
- max one current Block per Membership.
- unblock ends the Block.
- ending Membership ends any current Block.
- no generic Membership `ACTIVE/SUSPENDED/REMOVED` status.

Effective organization authority requires:

```text
current Membership
AND no current MembershipBlock
AND current RoleAssignment
AND required Permission
AND resource/business policy
```

### Organization Roles

Organization Role is a reusable permission bundle scoped to one Organization.

Role is not Ownership.

A current Membership has exactly one current Organization RoleAssignment in the current business model.

Role change:

```text
end current assignment
+ create new assignment
at the same effective transition
```

### Organization Ownership

Ownership is current Organization governance truth.

Canonical fact:

> Organization O is currently owned by Membership M.

Current evidence supports exactly one owner for each non-archived Organization.

Rules:

- owning Membership belongs to same Organization
- owning Membership is current
- owning Membership is usable under current policy
- creator/founder does not imply current ownership forever
- Owner Role is not the source of truth
- no `is_owner` boolean on Membership
- transfer is a dedicated business transition
- ownership history belongs to audit/events unless a real temporal query requirement emerges

Governance mutations serialize on the Organization row.

### Organization Branch

Branch is accepted as an operational subdivision of Organization, not merely an address.

Core facts:

- stable branch identity
- exactly one Organization
- current name
- optional current manager capacity
- terminal retirement

Rejected legacy representation:

- `type = branch|department|store` as proof that all are one ontology
- `is_active`
- `is_deleted/deleted_at`
- direct `user_id` member relation
- direct `manager_id → user/profile`
- branch `role_id` copied from generic IAM

Department and Store remain separate unresolved concepts. They are not silently converted into Branch or a generic OrganizationUnit.

### Branch Assignment

A Branch Assignment says:

> current Organization Membership M is assigned to Branch B.

It is a structural relationship, not a permission grant.

Rules:

- target is OrganizationMembership, never User/Profile
- Branch and Membership must belong to the same Organization
- same Branch/Membership pair exists at most once
- Membership→Branch maximum cardinality remains open until business evidence proves 1 or N
- ending Membership removes current Branch assignments
- retiring Branch removes current Branch assignments
- assignment history is audit/event unless a historical query requirement emerges

### Branch Manager

Branch Manager is a real Branch-scoped business capacity.

Representation target:

> Branch B is currently managed by Organization Membership M.

Rules:

- same Organization
- ended Membership cannot remain current manager
- manager relationship is not a global Organization Role
- manager does not bypass all Branch permissions
- policy explicitly decides which Branch intents a manager may execute
- no separate temporal manager table until historical manager queries are required

---

## 5. Group

Group is an ad-hoc collaboration context, not an Organization subtype and not an Organization Branch.

Legacy evidence:

- Group has its own name/description/avatar.
- Group has members and chat integration.
- Deal may belong to Group context.
- historical code intended one creator/leader, but `IsLeaderGroup` currently has the real check commented out.
- `created_by` therefore cannot remain current authority.

### Group Membership

A Group Membership is one Person↔Group participation episode.

Rules:

- target Person, not Profile/User
- at most one current Membership per Group+Person
- leave/remove ends episode
- rejoin creates a new episode

### Group Leadership

Group leadership is a current governance relationship from Group to one Group Membership.

Current v0 decision:

- exactly one leader for each non-disbanded Group
- creator becomes initial Group Membership + leader in one transaction
- `created_by` is provenance only
- transfer leadership is a dedicated intent
- leadership is not inferred from a generic role string

### Group lifecycle

Legacy `active/paused/disbanded` mixes reversible availability with terminal lifecycle.

Final core stores terminal `disbanded_at`. A future temporary pause model is not invented until its behavioral consequences are defined.

### Group role debt

Legacy `group_members.role` and `role_id` do not have coherent scoped semantics. User IAM currently resolves them through a role repository also used for Organization/Deal roles.

Target rule:

> Organization role, Group authority, and Deal authority do not share a polymorphic `role_id` merely because they all eventually influence authorization.

Group leader is explicit. Additional reusable Group roles are deferred until concrete permissions require them.

---

## 6. Deal context

A Deal has one primary business context:

- Organization; or
- Group.

Legacy `owner_id + owner_type` is rejected.

The words "owner" and "ownership" are removed from this context relationship because the relationship means containment/business context, not legal/resource ownership.

Target logical invariant:

```text
exactly one of:
  organization_id
  group_id
is present
```

Personal/User/Member-owned Deal is **not earned** by current live behavior.

Migration must still profile legacy owner types before contracting old columns.

### Branch placement

An Organization Deal may optionally be associated with one Branch.

Current evidence supports at-most-one Branch per Deal strongly enough for the final reference model.

Rules:

- Group Deal cannot have Organization Branch.
- Branch must belong to the same Organization as the Deal.
- Branch relationship is resource/operational scope, not automatic authorization.

---

## 7. Deal lifecycle and classification

Legacy lifecycle values PROCESSING / COMPLETED / CANCELED form one coherent lifecycle, but v0 stores the more precise facts:

- unresolved/current processing: no terminal outcome
- COMPLETED terminal outcome
- CANCELED terminal outcome
- terminal `ended_at`
- cancellation reason only for CANCELED

Deal kind is a separate classification:

- BROKERAGE
- INVESTMENT
- JOINT

Deal kind is not lifecycle state.

---

## 8. Deal Invitation and Participation

### Deal Invitation

Invitation says:

> Deal D issued one offer to Person P to participate.

Rules:

- target is Person
- Invitation is not Participation
- unresolved uses NULL outcome
- terminal outcomes: ACCEPTED / DECLINED / REVOKED
- terminal invitation is immutable
- resend after terminal outcome creates a new Invitation
- legacy WITHDRAWN is not Invitation state

Deal invitations may explicitly offer independent terms:

- customer relationship
- partner relationship
- administrator capacity

Ordinary participation needs no "member role" because Participation itself means membership in the Deal.

### Deal Participation

One row is one effective Person↔Deal participation episode.

Rules:

- acceptance creates Participation
- direct founder/creator participation may have no originating Invitation
- max one current Participation per Deal+Person
- leave/remove ends Participation
- rejoin creates a new Participation
- ended Participation is preserved because investments/contracts/financial records may reference that exact episode

Legacy `deal_members.status` is rejected because it mixed Invitation and Participation.

---

## 9. Deal participant dimensions

Legacy `role_key` mixed independent dimensions:

```text
OWNER / ADMIN / MEMBER / CUSTOMER / PARTNER
```

This loses information when one Person is both Partner and Admin, or Customer and Partner.

Canonical separation:

### Participation

Existence of DealParticipation already means "member".

No `MEMBER` role is stored.

### Customer relationship

Customer is a Deal-scoped business relationship, not Person identity.

Evidence:

- customer IDs resolve through the same Profile/Person universe
- Deal contract/payment records carry customer ID
- a Person may be customer in one Deal and something else in another

A Customer relationship is attached to DealParticipation.

### Partner relationship

Partner is a Deal-scoped business relationship, not Person identity.

Evidence:

- same Person/Profile identity space
- product distribution and partner commission behavior exist elsewhere
- source explicitly notes a user may overlap member/customer/partner

A Partner relationship is attached to DealParticipation.

Customer and Partner can coexist for the same Participation.

### Administrator capacity

Deal Admin is authorization capacity, separate from Customer/Partner business relationship.

Multiple admins are permitted in the target unless product requirements later impose another cardinality.

### Lead / legacy "Deal Owner"

Legacy `RoleKeyDealOwner` is not the same fact as Deal Organization/Group context.

It represents primary Deal governance/responsibility.

Final v0 decision:

- each current/open Deal has exactly one current Lead
- creator becomes initial Participation + Lead
- transfer Lead is a dedicated transition
- generic role update cannot create or remove Lead
- creator remains provenance, not authority
- `deal_members.is_owner` is rejected duplicate truth
- lead references DealParticipation, the narrowest correct contextual identity

### charge_person

Legacy `charge_person_id` is seeded as an invited ADMIN, not as Deal Lead.

It is therefore not kept as a second current responsibility scalar. The intended admin is represented by Invitation terms and, after acceptance, Deal Administrator assignment.

---

## 10. Deal financial facts

### Commission

Commission terms are not Participation identity.

They are a separate optional financial agreement attached to a Participation.

Supported canonical forms:

- PERCENT
- FIXED_AMOUNT

For fixed amount, amount and currency are stored together. For percent, currency is absent.

### Investment commitment

Legacy `amount_commit` is separated from actual investments.

A current Investment Commitment is an optional amount+currency term attached to Participation.

Historical revisions may remain audit/events unless product rules require contractual commitment history.

### Investment

One Investment row is one submitted contribution/payment evidence by a Deal Participation.

Core facts:

- Participation
- amount + currency
- transfer time
- proof file reference when present
- processing result
- note when business-relevant

Legacy duplicate flags are rejected:

- `confirmed`
- `done_investment` on deal member

Investment outcome uses a coherent workflow such as unresolved/PENDING → APPROVED or REJECTED.

### Costs

Deal Cost is a monetary business fact classified by Cost Type.

Money always carries currency.

### Target profit

Legacy source treats `target_profit` as an amount in financial summaries/commission limits, not a percent.

Target therefore becomes:

- target_profit_amount
- currency_code

### Bank account

Legacy Organization service `bank_accounts.owner_id + owner_type` is rejected.

Bank-account ownership belongs to Payment and points to Party. Deal references a selected bank account by external Payment-owned identifier.

---

## 11. Deal products, milestones, notes, documents

### DealProduct

Deal↔Product/Property association is many-to-many current business scope.

The external Product/Property authority remains BDSPro service; no runtime cross-database FK is assumed.

### Milestone

Milestone is a Deal child fact.

Core:

- title
- description
- expected date
- completed date
- ordering

No duplicate milestone status is required when completion is represented by `completed_at`.

### Internal note

Internal note is a record, not a string duplicated on Deal.

Legacy `deals.internal_note` is rejected in favor of note rows.

### Documents

Generic `attach_documents(owner_id, doc_owner)` is rejected.

Target uses explicit relationships:

- DealDocument
- InvestmentDocument

File bytes/metadata remain owned by file-service.

---

## 12. Deal authorization principle

Deal relationships do not automatically grant authority.

Example policy:

```text
required Permission/capacity
AND current DealParticipation
AND optional Branch scope match
AND Deal business state
AND operation-specific rules
```

Being assigned to Branch does not by itself grant access to all Branch Deals.

Being Partner or Customer does not automatically imply Admin permissions.

---

## 13. Property

Property truth remains owned by bdspro-service.

Current canonical source already has a dedicated `property` authority with partition/projection support.

v0 does not merge Property into Party, Deal, CRM, or Planning.

Cross-context rules:

- DealProduct references Property/Product through the BDSPro contract.
- CRM may project/reference Property but does not own Property truth.
- TQD planning/location data may enrich Property but does not own Property identity.
- `source_type` partitioning is physical storage policy, not a new Party subtype.

---

## 14. Commerce / Payment

Payment is a separate bounded context.

Existing commerce flow is retained because it already separates:

- Order
- PaymentAttempt
- Settlement
- Fulfillment
- CommandEffect
- OutboxEvent

Important target correction:

Legacy/request-facing `subject_kind + subject_id` for Profile/Organization is logically replaced by one Party reference where the billing subject is a business Party.

This is a logical external reference; payment-service cannot hold a PostgreSQL FK to user-service.

Wallet remains a personal financial account unless product evidence adds Organization wallets.

Bank account ownership references Party, removing generic owner type/id.

Money in Payment uses minor units + currency, which is retained.

---

## 15. Catalog / Subscription

Catalog and Subscription authority remains user-service-owned.

Existing separation is retained:

- Product
- Plan
- PlanVersion
- Entitlement
- OperationPolicy
- PriceItem
- Subscription
- SubscriptionEvent
- CheckoutSnapshot
- SettlementReceipt

Logical subject identity converges on Party rather than profile/organization type pairs.

Commercial snapshots remain snapshots: they do not become the authority for Organization/Person identity.

---

## 16. TQD / Planning

TQD remains an independent bounded context.

Its current canonical migration authority contains planning layers, regions, labels, planning projects/documents, legal documents, parcels/POIs, reporting, audit/lifecycle, quota, AI jobs, saved/followed planning workflows and map points.

No attempt is made to merge these tables into the Organization/Deal core merely because they may be displayed together.

Quota is explicitly separated into:

- durable usage/event facts
- quota pools/current projection
- report jobs/process state

---

## 17. CRM / SEO / Support

CRM owns CRM/SEO/support semantics.

Current canonical tables remain separate for:

- domain registry
- SEO links/relations/trends/generation
- search/traffic/ranking/crawl measurements
- objectives/evidence/metrics
- audience/user need
- intent/query clusters
- content taxonomy/placement
- import/source state/page score
- support ticket/note/event
- opportunity/customer workflow
- saved planning news
- payment event inbox + CRM activity projection

CRM `customer` is not Person identity. It is CRM/business-process state that may refer to a Person/Profile.

---

## 18. Notification, File, Hub

Notification owns delivery/history and event consumption.

File owns file metadata/access.

Hub owns platform lookup/help/config/compatibility data.

These services may reference business IDs but do not become the write authority for the referenced business facts.

---

## 19. Events, audit and projections

Do not confuse:

- current domain state
- domain history
- audit/security evidence
- integration event
- read projection
- snapshot

Outbox/inbox rows may legitimately use generic event metadata because they are operational envelopes, not canonical business ownership.

Audit may record actor IDs, reason and before/after data without those columns being copied onto every domain table.

---

## 20. Concurrency rules

Use the smallest coordinator that protects the invariant.

### Organization governance

Lock Organization first for:

- ownership transfer
- owner leave/remove
- block owner
- archive Organization

### Invitation acceptance

Lock Invitation before creating Membership/Participation.

### Branch governance

Lock Branch for:

- manager change
- Branch assignment if retirement may race
- Branch retirement

### Deal governance

Lock Deal for:

- Lead transfer
- terminal completion/cancel
- operations that must reject after terminal outcome

### Participation

Lock Participation/Deal when terminating a Participant if current Lead/Admin/financial invariants depend on that Participant.

Serializable isolation is not the default. Prefer constraints + atomic statements + row locks where the invariant requires them.

---

## 21. Negative schema — explicitly rejected

The final model does not introduce:

- generic `party_type + party_id` ownership
- generic `owner_type + owner_id`
- global `customer` or `partner` identity subtypes
- global Actor table
- generic Resource table
- generic Transaction table for unrelated domains
- EAV
- one `status` column mixing proposal/current/block/termination
- soft delete as business lifecycle
- `created_by` as current authority
- `profile_id` as universal business identity
- reusable Organization Role as Branch/Group/Deal ownership truth
- one generic OrganizationUnit for Branch/Department/Store without business proof
- direct DB joins/FKs across microservice databases

---

## 22. v0 closure matrix

| Concept | v0 result |
|---|---|
| Party | GREEN |
| Person | GREEN |
| Account | GREEN |
| Profile as presentation | GREEN |
| Account credentials/session separation | GREEN |
| Permission | GREEN |
| Runtime Actor | GREEN, not a table |
| Organization | GREEN |
| Organization Invitation | GREEN |
| Organization Membership | GREEN |
| Membership Block | GREEN |
| Organization Role/Assignment | GREEN |
| Organization Ownership | GREEN |
| Branch | GREEN |
| Branch Assignment | GREEN |
| Branch Manager capacity | GREEN |
| Department/Store ontology | DEFER |
| Organization address meaning | DEFER |
| Group | GREEN |
| Group Membership | GREEN |
| Group Leadership | GREEN |
| Group generic role_id | REJECT |
| Deal Organization/Group context | GREEN |
| Generic Deal owner type/id | REJECT |
| Deal Branch association | GREEN |
| Deal Invitation | GREEN |
| Deal Participation | GREEN |
| Customer relationship | GREEN |
| Partner relationship | GREEN |
| Deal Admin capacity | GREEN |
| Deal Lead capacity | GREEN |
| deal_members role_key/status/is_owner blob | REJECT |
| Commission terms | GREEN separate fact |
| Investment commitment | GREEN separate fact |
| Investment | GREEN separate fact |
| Deal Costs | GREEN |
| Generic attach_document owner | REJECT |
| Property authority | RETAIN BDSPro bounded context |
| Payment commerce model | RETAIN with Party subject correction |
| TQD | RETAIN bounded context |
| CRM/SEO/support | RETAIN bounded context |
| Notification/File/Hub | RETAIN bounded contexts |

v0 is closed enough to proceed to v1.
