# Database Migration and Authority Cutover Plan

## Non-negotiable migration pattern

Every convergence follows:

OLD → EXPAND → BACKFILL → COMPARE / PROVE → SWITCH READS → SWITCH WRITES / AUTHORITY → OBSERVE → CONTRACT.

No reset, destructive rewrite, dual independent writer or big-bang cutover is allowed.

## Phase 0 — Freeze evidence

1. Record exact source SHA and migration checksums.
2. Export schema-only snapshots for every service database.
3. Record row counts and key distributions.
4. Profile nullable/orphan/duplicate contradictions before target constraints are added.
5. Preserve legacy tables read-only until reconciliation evidence is signed off.

Required profiling includes:

- Deal owner_type counts and owner_id=0 counts.
- Organization owner scalar vs owner role disagreement.
- Organization member invited/active/suspended/removed distributions.
- Multiple active Organization members for same org/profile.
- Multiple current owner candidates.
- Branch members without Organization membership.
- Branch managers outside Organization.
- Deal role_key vs is_owner disagreement.
- Deal Customer/Partner/Member collisions for same Deal+Profile.
- Deal rows whose owner context conflicts with deal_of_* rows.
- Duplicate/reopened invitations.
- Product/Asset owner_type populations.
- Property old-generation ↔ new-generation mapping coverage.
- Chat participant ↔ conversation_membership disagreement.
- Notification histories whose canonical domain owner can be reconstructed.
- Hub/TQD province/ward code disagreements.

## Wave 1 — Party / Person / Account identity

Expand parties, persons, profiles and accounts while legacy Profile IDs still work. Prefer deterministic ID mapping and preserve legacy Profile ID as Person/Party ID where collision analysis proves it safe.

Backfill Profile descriptive data into profiles. Backfill authentication methods into provider-specific credential tables. Sessions/devices map to Account.

Compare: every live Profile used by a foreign service resolves to one Person; every active auth method resolves to one Account/Person.

Switch internal contracts from ambiguous user/profile identity to Person/Account semantics before deleting any legacy column.

## Wave 2 — Platform IAM

Split platform system roles from organization roles. Seed permissions by stable code. Backfill current platform assignments. Keep legacy role reads as compatibility projection until authorization matrix proof passes.

## Wave 3 — Organization

Expand canonical invitations, memberships, blocks, role assignments, ownerships, branches and branch assignments.

Backfill order:

1. Organizations and Party mapping.
2. Effective Organization memberships from legacy active rows.
3. Ownership from scalar owner plus owner membership evidence; conflicts enter reconciliation queue.
4. Current role assignments.
5. Pending invitations from rows that are truly proposals, not active participation.
6. Membership blocks from temporary suspension evidence.
7. Branches and current assignments.
8. Branch manager Membership endpoint.

Do not infer Invitation merely because a Membership exists. Do not infer owner from creator when live owner evidence differs.

Switch governance write paths first: CreateOrganization, Invite, Accept, Block/Unblock, Leave/Remove, ChangeRole, TransferOwnership, Archive, Branch assignment/manager changes.

Then switch read/authorization paths. Finally stop writes to owner_profile_id, organization_members.status/role_key and branch user IDs.

## Wave 4 — Group

Create GroupMembership for current group members. Founder/creator becomes initial lead only when live history and current business data support it; otherwise reconcile explicitly.

Introduce GroupLead as the only current leader source. Replace IsLeaderGroup creator shortcut. Remove group member role/status authority once new paths prove equivalent.

## Wave 5 — Deal

Expand explicit primary context, invitations, participations, customer/partner relations, admin capacity, lead, commissions, investments, costs and audit.

Backfill context:

- owner_type=ORGANIZATION → organization_id after validating endpoint.
- owner_type=GROUP → group_id after validating endpoint.
- USER/MEMBER/unknown/zero → reconciliation queue; do not invent Person/Membership Deal mode.
- compare deal_of_organization/deal_of_group/deal_of_branch and record conflicts.

Backfill participation:

- accepted deal_members → DealParticipation.
- invited rows → DealInvitation only.
- rejected rows → terminal Invitation.
- withdrawn rows → historical Participation end when acceptance evidence exists; otherwise terminal proposal evidence is reconciled.

Backfill relationship/capacity using numeric source + creation history, not display labels, because Customer/Partner legacy labels are swapped in RoleKeyMap.

Lead backfill uses RoleKeyDealOwner plus creation/removal history. is_owner is comparison evidence only, never authority.

Switch dedicated intents. Generic UpdateDealMemberRole is not allowed to create/remove Lead after cutover.

## Wave 6 — Property / Asset / Listing

Choose the newer property identity as authority and build an explicit mapping from legacy property_identify/property_lineage. Preserve external reference and evidence history.

Backfill legacy Product into Listing, preserving price history. Backfill Asset tied to Property, then map owner to Party. Generic owner types that cannot be resolved become reconciliation records.

After compare proof, remove property_relation polymorphism and direct generic Product/Asset ownership.

## Wave 7 — Commercial / Payment

Commercial tables are already relatively mature. Re-key profile/organization subjects to Party without changing plan/version/order semantics.

Preserve immutable PlanVersion, checkout snapshots, settlements and outbox IDs. Never regenerate provider evidence.

For Wallet, introduce immutable ledger entries and dual-calculate balance until ledger sum equals current balance for every wallet. Only then can mutable balance cease to be authority.

## Wave 8 — Chat / Social / CRM

Chat: backfill one ConversationMembership from the union of participants and conversation_memberships. Conflicts are resolved before participants is retired.

Social: backfill Post author/context from owner fields and news_feed_of_* links; detect contradictions. Split generic reactions by target table.

CRM: rename Customers semantics to Opportunities without rewriting Person identity. Convert Profile endpoints to Person. Preserve SEO page state established by migration 000024.

## Wave 9 — File / Notification

File: preserve stable storage key, hashes and idempotent owner effect keys; do not restore absolute_path. Convert access IDs to explicit Party grants where access is human/business identity.

Notification: create Party-recipient Notifications. Keep event inbox/delivery fencing. Stop creating foreign-domain audit rows once each owning domain writes its own audit.

## Wave 10 — Hub / TQD reference convergence

Hub becomes platform administrative province/ward authority. Use stable public codes to map current province_v2/ward_v2 and TQD copies.

TQD retains planning jurisdiction and imported geometry, which are different facts from the platform administrative catalog.

Reports rename/re-key without changing command idempotency. Report + quota usage + job atomicity must remain intact.

## v2 physical constraints to add after clean backfill

- Partial unique current OrganizationMembership on organization_id+person_id where ended_at IS NULL.
- Partial unique current MembershipBlock per membership where ended_at IS NULL.
- Partial unique current RoleAssignment per membership where ended_at IS NULL.
- Current Organization ownership PK on organization_id and unique membership_id.
- Partial unique current GroupMembership per group+person.
- Current GroupLead PK group_id and unique membership_id.
- Exactly-one Deal context check: organization_id XOR group_id.
- Branch check implemented by transaction validation/constraint strategy: branch requires organization_id and same Organization.
- Partial unique current DealParticipation per deal+person.
- DealLead one per Deal.
- Commission-term check: percentage XOR fixed amount; percent in [0,100].
- Listing context exactly one of person/organization/group.
- Listing branch same Organization invariant.
- Partial unique current AssetOwnership until co-ownership is explicitly introduced.
- Message sequence unique per conversation.
- Checkout/order idempotency uniqueness and fingerprint mismatch rejection.
- Job/outbox claim_version monotonic and lease state constraints.

## Compare proof

Each wave produces:

- old row count;
- target row count;
- mapped count;
- deliberately excluded count with reason;
- orphan count;
- semantic conflict count;
- sample old→new mappings;
- authorization behavior comparison;
- exact source SHA and migration SHA.

Cutover requires zero unexplained conflicts, not merely equal row counts.

## Rollback

Before authority switch, rollback means return reads/writes to old owner while target remains shadow-updated/backfilled.

After authority switch, old tables are read-only compatibility sources. Rollback must not reactivate old writers without reconciling new writes first.

## Contract phase

Drop/rename legacy columns/tables only after:

1. no production caller reads them;
2. no writer mutates them;
3. compare remains clean for the agreed observation window;
4. backup/export exists;
5. rollback procedure no longer requires the old representation;
6. source tests prohibit reintroduction of the retired pattern.