# BDSPro Architecture Naming & Ownership Mindset

Updated: 2026-09-22 Asia/Ho_Chi_Minh
Status: ACTIVE DURABLE WORKING MINDSET

Live Git/source and completed exact-SHA proof remain higher authority than this document.

## 1. Governing rule

`ONE BUSINESS FACT -> ONE OWNER -> ONE SOURCE OF TRUTH -> ONE CLEAR WRITE PATH -> MANY READERS MAY EXIST BUT MUST NOT REDEFINE THE FACT`

Split responsibility by business fact and invariant, not by file count, method count or fashionable architecture terms.

## 2. Required language

Use direct terms that name the actual business or runtime fact.

Prefer:
- user / profile;
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
- allow / deny;
- session;
- token;
- cache;
- database;
- source of truth;
- service owner;
- consumer;
- invariant;
- evidence.

Do not introduce vague abstraction words when a direct business term exists. Existing old names may be retained temporarily only when changing them would break a proven contract. New work must not copy those names without a concrete reason.

## 3. Naming test

Every exported abstraction or public entry point must answer:
1. What exact fact does it own?
2. What invariant does it protect?
3. Who writes it?
4. Who reads it?
5. Why would it change?
6. What real behavior disappears if it is removed?

Names such as `Manager`, `Helper`, `Util`, `Common` and generic `Service` require additional scrutiny when they hide more than one responsibility.

## 4. One fact, one owner

Do not keep two durable records that can disagree about the same business fact.

Read copies and caches are allowed only when:
- their source is explicit;
- refresh/invalidation is explicit;
- maximum staleness is explicit;
- failure behavior is explicit;
- they cannot become an independent write authority.

## 5. Public surface

If callers express one intent, expose one clear entry point even when internal mechanics differ.

If two operations preserve different invariants, keep them separate even when their code looks similar.

## 6. Authentication / Organization / Authorization

### Authentication
Proves who the caller is and whether the session/token is valid.

### Organization
Owns the business structure: organization identity, owner, invitations, membership, branch/team structure, organization roles, organization role assignment and organization audit.

### Authorization
Decides whether an authenticated actor may perform an action on a resource inside a specific scope under the current relationships and state.

Role alone is never assumed to be the final decision.
Organization context alone is never assumed to prove membership or permission.

## 7. Evidence discipline

Use:
`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

and:
`READ -> UNDERSTAND -> SEARCH -> TRACE -> CHANGE -> TEST`

Prefer evidence in this order:
1. exact live source/migrations;
2. executable runtime wiring;
3. registered RPC/API contracts;
4. tests and hosted proof;
5. reconciled durable documents;
6. current code comments;
7. older documents only for tracing origin.

Do not infer ownership from directory names.

## 8. Architecture change rule

Before adding a service, cache, table, abstraction or external authorization engine, identify:
- the business problem;
- current owner;
- current consumer;
- failure mode;
- required correctness;
- required revocation time;
- actual load;
- measurable value of the new component.

Reject the addition when it merely creates another name or another copy of existing truth.

## 9. Authorization implementation order

Do not optimize authorization checks before proving the decision itself is correct.

Order:
1. define organization operating model;
2. define roles and permissions;
3. define scopes and resource relationships;
4. define allow/deny ownership;
5. define revocation and audit;
6. implement durable owners;
7. implement resource checks;
8. prove search/list isolation;
9. then measure and add cache/read copies if needed.

## 10. Durable phase reference

The detailed active foundation is:
`documents/architecture/AUTHENTICATION-AUTHORIZATION-ORGANIZATION-FOUNDATION.md`

Future work on User/Auth/Organization/Gateway must read that file before design or mutation.
