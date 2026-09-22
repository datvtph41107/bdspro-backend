# BDSPro Authentication, Authorization & Organization Foundation

Updated: 2026-09-22 Asia/Ho_Chi_Minh
Status: ACTIVE DURABLE FOUNDATION
Source authority: `datvtph41107/bdspro-backend` branch `main`, live source wins whenever prose disagrees.

## 1. Purpose of this phase

This phase establishes the complete operating order for:
- human identity;
- login methods;
- sessions;
- organization lifecycle;
- organization ownership;
- invitations;
- membership;
- branches;
- teams;
- roles;
- permissions;
- resource ownership;
- resource relationships;
- authorization decisions;
- revocation;
- audit;
- service responsibility.

The target is not a discussion-only model. The target is a model precise enough to be implemented, tested, migrated and enforced across the repository.

## 2. Required language

Use direct business and engineering terms.

Preferred vocabulary:
- identity;
- profile;
- login method;
- session;
- organization;
- owner;
- member;
- invitation;
- branch;
- team;
- role;
- permission;
- action;
- resource;
- scope;
- relationship;
- state;
- restriction;
- decision;
- allow;
- deny;
- revoke;
- audit;
- source of truth;
- read copy;
- cache;
- database;
- request path;
- service responsibility.

Do not introduce vague abstraction words when a direct business term exists. Existing old names in source may remain temporarily for compatibility, but no new design or code should reproduce them merely because they already exist.

## 3. Governing rule

`ONE BUSINESS FACT -> ONE OWNER -> ONE SOURCE OF TRUTH -> ONE CLEAR WRITE PATH -> MANY READERS MAY EXIST BUT MUST NOT REDEFINE THE FACT`

Every important fact must answer:
1. What is the fact?
2. Who creates it?
3. Who may change it?
4. Which service owns the durable record?
5. Which other services may read it?
6. How quickly must a change take effect?
7. What happens when the owner is unavailable?
8. What audit record must remain?

## 4. Authentication boundary

Authentication answers only:
`Who is making this request, and is the session/token valid?`

Authentication owns or uses:
- AuthID;
- ProfileID;
- login method;
- password or external identity linkage;
- session ID;
- token type;
- token issued time;
- token expiry;
- session state;
- refresh state;
- signing and verification trust.

Authentication must not answer final business authorization questions.

## 5. Trusted request identity

Canonical request identity should be centered on:
- AuthID;
- ProfileID;
- SessionID;
- TokenType;
- OrganizationID when present;
- coarse system role only when a proven system-wide use requires it.

ProfileID is the default human owner/audit identity.
OrganizationID is context, not proof of membership and not proof of permission.
A client must never gain authority by sending its own organization, role or permission headers.

## 6. JWT direction

JWT should primarily carry authenticated identity and request context.

Target minimal direction:
- auth_id;
- profile_id;
- session_id;
- token_type;
- active_organization_id when needed.

A role may remain only when a specific system-wide use proves it is needed.
Do not make current permissions depend on a long-lived list of permissions embedded in the token.
Permission removal, membership suspension and organization removal must not wait indefinitely for token expiry.

## 7. Organization is a first-class business structure

Organization is not merely a field on Profile.

An organization has its own:
- identity;
- owner;
- lifecycle;
- members;
- invitations;
- branches;
- teams;
- roles;
- permission assignment;
- security rules;
- audit history;
- billing relationships;
- resource relationships.

This does not automatically mean Organization must remain a separate network service.

## 8. Organization lifecycle

Minimum lifecycle:
1. create organization;
2. establish first owner;
3. invite member;
4. accept or decline invitation;
5. activate membership;
6. assign role;
7. assign branch/team when needed;
8. suspend member;
9. restore member;
10. member leaves;
11. member removed;
12. transfer resources;
13. transfer organization ownership;
14. archive organization.

Required invariants:
- an active organization always has at least one owner;
- the last owner cannot be removed without ownership transfer;
- suspended/removed members lose organization authority within the agreed revocation time;
- historical records remain;
- ownership transfer is explicit and auditable.

## 9. Membership

Membership is the relationship between a Profile and an Organization.

`Profile != Membership`

One Profile may have many memberships with different roles.

Minimum state direction:
`INVITED -> ACTIVE -> SUSPENDED -> REMOVED`

Membership must answer organization, profile, state, join/remove time, active role assignment and branch/team relationships.

## 10. Invitation

Invitation is not the same record as active membership.

Invitation needs:
- inviter;
- organization;
- target identity;
- intended role;
- intended branch/team when relevant;
- expiry;
- accepted/declined/cancelled state;
- accepted Profile;
- audit timestamps.

## 11. Role

Role means:
`a reusable group of permissions assignable inside a defined scope`

Examples:
- Owner;
- Organization Admin;
- Branch Manager;
- Team Manager;
- Broker;
- Billing;
- Legal;
- Viewer.

Role is not the final authorization decision.

## 12. Permission

Permission names one allowed action class.

Examples:
- organization.member.invite;
- organization.member.suspend;
- organization.role.assign;
- listing.create;
- listing.update;
- lead.view;
- deal.manage;
- billing.checkout.

Permission alone does not prove access to every resource of that type.

## 13. Scope

Scope answers where a role or permission applies:
- whole system;
- one organization;
- one branch;
- one team;
- one resource.

System roles and organization roles must not be silently interchangeable.

## 14. Resource relationship

Business services own relationships specific to their resources.

Examples:
- Listing -> organization/branch/assigned broker;
- Lead -> organization/team/owner;
- Deal -> participants;
- Conversation -> members;
- Payment -> payer organization;
- Report -> requester/organization.

These relationships must not all be copied into one central authorization database merely to simplify checks.

## 15. Authorization question

Every authorization decision should be reducible to:
`WHO wants to perform WHAT ACTION on WHICH RESOURCE inside WHICH SCOPE under WHICH CURRENT STATE and RELATIONSHIPS?`

Conceptual order:
1. authenticated identity exists;
2. session is valid;
3. organization context is valid when required;
4. membership is active;
5. required permission exists;
6. resource is inside an allowed scope;
7. required resource relationship exists;
8. resource state permits the action;
9. no explicit restriction blocks the action;
10. allow or deny.

## 16. Responsibility direction

### Gateway
- verify external token;
- build trusted request identity;
- send signed service identity downstream;
- do not own business resource decisions.

### User / identity owner
- Profile;
- login methods;
- durable account identity;
- system roles only when truly system-wide;
- system permission definitions;
- session-related identity facts where current source proves ownership.

### Organization owner
- organization lifecycle;
- owner;
- invitations;
- membership;
- organization roles;
- organization role assignment;
- branch/team structure;
- ownership transfer;
- organization audit.

Whether this lives inside User or an independent Organization process must be decided from source, data, runtime dependency and business lifecycle evidence.

### Auth process
May remain only if its read/caching function has measurable value.
It must not become a second durable owner of organization membership, roles or permissions.

### Business resource services
Each resource-owning service keeps final checks that depend on its own data.

## 17. Allow/Deny ownership rule

For every protected action identify:
- human identity owner;
- organization membership owner;
- permission source;
- resource owner;
- final service returning allow/deny;
- audit owner;
- revocation path.

## 18. Revocation

Define revocation behavior for:
- session revoked;
- account locked;
- membership suspended;
- member removed;
- role changed;
- permission removed;
- branch/team assignment removed;
- ownership transferred.

For each define maximum delay, cache/read copy invalidation, failure behavior, fail-closed behavior and audit evidence.

## 19. Audit

Authorization audit is not ordinary logging.

Important events:
- organization created;
- owner transferred;
- invitation lifecycle;
- member suspended/restored/removed;
- role changes;
- permission assignment changes;
- branch/team assignment changes;
- resource ownership changes;
- high-risk deny;
- emergency system access.

Audit must answer who, what, target, organization, time, session/request and reason when needed.

## 20. Search and list access

Authorization applies before returning collections/search results.

Never rely on:
`fetch everything -> return everything -> frontend hides rows`

Query construction must include allowed scope/resource relationships when exposure would leak unauthorized data.

## 21. System administrator vs organization administrator

These are different roles.

System administrator operates the BDSPro platform only where explicitly granted.
Organization administrator operates one organization and gains no authority over other organizations or system administration merely because the role name contains `admin`.

## 22. Admin surface

An Admin web application does not imply an Admin service.
Administrative operations should call the owning service directly unless a real cross-domain process has its own durable state, lifecycle, invariants and audit requirements.

## 23. Current repository questions to resolve

1. User-owned organizations/members versus Organization-service records;
2. User roles/permissions versus Organization-service roles/permissions;
3. RoleProfile account assignment versus organization-scoped membership roles;
4. JWT role/roleIds/authorities versus current durable authorization state;
5. Auth role cache keyed by ProfileID versus organization-scoped authority;
6. organization-service ownership of group/deal/investment behavior;
7. final owner of branch/team structure;
8. final owner of invitation lifecycle;
9. final owner of organization audit;
10. final owner of member removal and resource transfer.

No new authorization framework is selected before these are resolved from business behavior and source evidence.

## 24. Implementation order

### Phase A — Source map
Trace current tables, RPCs, handlers, writes, reads and consumers for Gateway, User, Auth and Organization.

### Phase B — Organization operating model
Define create, owner, invitation, membership, roles, branch/team, suspension, removal, transfer and audit end-to-end.

### Phase C — Authorization matrix
For every protected action define actor, action, resource, scope, required permission, required relationship, state checks, restrictions, final decision owner, revocation requirement and audit requirement.

### Phase D — Ownership convergence
Choose one durable owner for each fact. Remove duplicate writes and conflicting read paths.

### Phase E — Authentication cleanup
Reduce JWT to required identity/context evidence, establish session revocation, signing trust and organization-context binding.

### Phase F — Authorization implementation
Implement organization checks and resource-owner checks using the simplest structure that satisfies proven requirements.

### Phase G — Cache and performance
Only after decision correctness is proved, measure hot paths and add caching/read copies where justified.

### Phase H — Repository-wide migration
Move callers one group at a time, prove behavior, remove superseded paths and add tests that prevent ownership from splitting again.

## 25. Evidence standard

Use evidence in this order:
1. exact live source and migrations;
2. executable runtime wiring;
3. registered protobuf/RPC paths;
4. tests and hosted proof;
5. reconciled durable documents;
6. comments near current code;
7. older documentation only for tracing why something exists.

## 26. Completion definition

Complete only when:
- Authentication responsibility is explicit and enforced;
- Organization lifecycle is fully defined and implemented;
- membership transitions are implemented and tested;
- role/permission meanings are unambiguous;
- system and organization roles are separated;
- multi-organization users work correctly;
- active organization context cannot grant authority by itself;
- resource relationships are enforced by their owners;
- revocation times are defined and proved;
- search/list paths do not leak unauthorized records;
- audit events are durable and queryable;
- duplicate organization/IAM sources are removed or explicitly read-only;
- Auth process has a proven reason to remain or is simplified/retired;
- JWT contains only justified identity/context fields;
- repository tests prevent ownership from splitting again.

## 27. Working scenario

`Phúc creates Organization ABC -> becomes Owner -> invites Lan -> Lan accepts -> Lan becomes ACTIVE -> Lan joins Branch Hà Nội -> Lan joins Team Sales A -> Lan receives Broker role -> Lan creates Listing -> Lan is later suspended -> her access is revoked within the agreed time -> her historical Listing/Lead/Deal records remain correctly owned/audited.`

Every design decision must be testable against this scenario and counterexamples.
