# Database Version Evolution — v0 → v1 → v2 → Final

## v0 — Semantic closure

Goal: discover what the business facts actually are before deciding tables.

v0 rejected the main legacy failure patterns:

- Profile/User/Person/Account conflation.
- Invitation stored as Membership/Participation status.
- Owner scalar duplicated with owner role.
- Creator used as current authority.
- Generic owner_id + owner_type.
- Generic status mixing temporary restriction and terminal lifecycle.
- Branch/Department/Store collapsed by a type flag.
- Deal member row mixing invitation, participation, role, customer/partner relation, commission and investment state.
- Generic target/type structures where a narrow FK exists.
- Soft delete used as a substitute for business lifecycle.

v0 established the canonical facts:

- Party, Person, Profile and Account are distinct.
- Organization Membership is an effective participation episode.
- Organization Invitation is a proposal; accepting it creates Membership.
- Organization Ownership is current governance relation to Membership.
- Membership Block is temporary authority restriction, not Membership termination.
- RoleAssignment is temporal; Permission is an action primitive, not final allow/deny.
- Branch is a real operational subdivision. Department and Store are not automatically Branch subtypes.
- Branch assignment and Branch manager capacity attach to OrganizationMembership.
- Group has participation and current lead; creator is not the authority source.
- Deal belongs to exactly one primary business context: Organization or Group. Branch is secondary operational scope.
- DealInvitation and DealParticipation are distinct.
- Deal participant identity is Person. Customer and Partner are scoped relationships; Admin and Lead are authority capacities.
- Property, Listing and Asset are different facts.
- Asset ownership is genuine Party ownership; Listing context is not legal ownership.
- Commercial payment facts, planning usage facts, event inbox/outbox and job claim facts remain separate from domain identity.

## v1 — Relational target

Goal: express each accepted fact with the narrowest relational representation.

Key v1 transformations:

| Legacy | v1 target |
| --- | --- |
| user_profile as identity + presentation + authority | persons + profiles + accounts + explicit IAM |
| roles with organization_id/domain_type magic | system_roles OR organization_roles |
| organization_members status invited/active/suspended/removed | organization_invitations + organization_memberships + organization_membership_blocks |
| organizations.owner_profile_id + role owner | organization_ownerships |
| organization_branch_members.user_id | organization_branch_assignments.membership_id |
| organization_branches.manager_id user/profile | manager_membership_id |
| groups.created_by as leader | group_memberships + group_leads |
| deals.owner_id/owner_type | deals.organization_id OR deals.group_id |
| deal_of_organization / deal_of_group duplicate | direct primary Deal context |
| deal_of_branch duplicate/ambiguous | deals.branch_id when applicable |
| deal_members blob | invitations + participations + relationship/capacity/economic facts |
| deal_members.role_key OWNER | deal_leads |
| deal_members.role_key MEMBER | participation existence |
| deal_members role CUSTOMER/PARTNER | separate scoped relationship tables |
| deal_members commission columns | deal_commission_terms |
| deal_members.done_investment | derived from deal_investments |
| Product/Asset generic owner | explicit Listing context / Asset Party ownership |
| property generations in parallel | one Property authority plus explicit Listing/Asset relations |
| chat participants + conversation_memberships | one conversation_memberships truth |
| social owner_of + owner_id + news_feed_of_* | explicit post author/context |
| notification as foreign-domain audit warehouse | notification delivery only; audits return to owning domains |
| mutable wallet balance as sole truth | immutable wallet_entries; balance becomes projection |
| profile/organization subject_kind+subject_id in target logical model | Party subject where the subject is truly Person-or-Organization |

v1 permits explicit nullable alternatives only when they are real closed-world business alternatives and a database check can enforce exactly one. It does not permit open-ended polymorphic owner/type pairs.

## v2 — Operational hardening

Goal: prove invariants under retries, concurrency, crashes and distributed delivery.

v2 requires:

1. Partial unique indexes for one-current-episode invariants.
2. Same-context checks in write paths where PostgreSQL cannot express a cross-row invariant with a simple FK.
3. Transaction boundaries around facts that must become visible atomically.
4. Explicit lock coordinators for low-frequency governance transitions.
5. Idempotency keys/fingerprints for externally retried commands.
6. Durable outbox for state-change events that must cross service boundaries.
7. Durable inbox/deduplication at consumers.
8. Claim version/lease fencing for workers.
9. Audit owned by the domain that understands the event.
10. Search/list authorization before data leaves the owner.

Examples:

- Accept Organization Invitation: lock invitation → validate target and current membership absence → create Membership + initial RoleAssignment → resolve Invitation in one transaction.
- Transfer Organization Ownership: lock Organization first → re-read owner/current memberships/blocks → update one current relation → audit/outbox.
- End Membership: end role assignment/block, remove branch assignments, clear/replace manager/owner relations as required, then set ended_at in one transaction.
- Accept Deal Invitation: lock invitation → validate target/no current participation → create DealParticipation → resolve invitation.
- Create report: Report + quota usage + ReportJob in one transaction; Redis remains runtime projection only.
- Payment: provider settlement evidence commits before downstream fulfillment; outbox publication is retriable.

## Final — Accepted system model

Final combines:

- v0 semantic boundaries;
- v1 normalized logical schema;
- v2 enforcement/operational rules;
- explicit owner-by-service;
- legacy disposition and cutover order.

Final deliberately does not unify all databases. The architecture converges truth, not physical deployment.

## Final non-negotiables

- No new canonical owner_id + owner_type pairs.
- No new status column that mixes independent lifecycle dimensions.
- No creator field used as authority.
- No proposal row mutated into an effective relationship.
- No cross-service write transaction.
- No event consumer allowed to become a second writer of producer truth.
- No cache/projection used as durable truth.
- No destructive contraction before compare proof.
- No migration that silently discards unresolved legacy data.