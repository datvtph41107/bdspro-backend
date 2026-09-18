# BDSPro Architecture Naming, Ownership & Public Surface Doctrine

Updated: 2026-09-19 Asia/Bangkok
Status: ACTIVE DURABLE WORKING DOCTRINE

This document refines how BDSPro names abstractions, exposes developer/public entry points, and decides when to split or converge concepts. It is an operating/design doctrine, not source authority. Immutable live Git/source and completed exact-SHA proof remain higher authority.

## 1. Governing invariant

`ONE TRUTH PER CONCERN -> ONE OWNER -> ONE PUBLIC NAME -> ONE OBVIOUS ENTRY POINT -> MULTIPLE PROJECTIONS THAT NEVER REDEFINE TRUTH`

This does NOT mean one giant object.

Different invariants or reasons-to-change must have different owners. The same caller intent must not be exposed through multiple equally-valid names merely because internal mechanics differ.

## 2. Split by invariant, not by file or method count

Do not create catch-all abstractions such as `QuotaManager` when they own multiple independent truths.

Example separation:
- pricing truth -> `QuotaPricer` / `quota.Pricer`;
- runtime reservation/admission truth -> `QuotaReserver` / `quota.Reserver`;
- durable usage/consumption truth -> `QuotaLedger` / `quota.Ledger`.

`Reserve` and `Release` may remain together when they are one reservation lifecycle. One method is not automatically one interface; one invariant is one abstraction.

Consumer-facing interfaces should remain as small as the consumer actually needs. Do not create implementation-owned wide interfaces merely to support mocking.

## 3. Naming test

Every exported abstraction/name must answer:
1. What truth does it own?
2. What invariant does it preserve?
3. Who consumes it?
4. Why would it change?
5. What capability disappears if it is removed?

Names such as `Manager`, `Helper`, `Util`, `Common`, `Service` require special scrutiny when they hide multiple responsibilities.

Prefer capability/role names that expose the real responsibility.

## 4. Public surface vs internal mechanics

Internal implementations may remain separate when mechanics differ, while the public caller surface stays singular if caller intent is singular.

Runtime/logging example:
- raw process output and structured retained evidence remain distinct internal truths/mechanics;
- developer intent is still "inspect logs/evidence";
- target public vocabulary converges on one `make logs` entry point with filters/options;
- the existing R4 `make log-query` implementation remains historical/proved capability and may be retained internally during migration, but should not remain a second long-lived daily public command if `make logs` can own the same developer intent without ambiguity.

Do not reopen R4 behavior or delete the query engine merely for naming cosmetics. Any public-surface retirement requires compatibility/docs/tests and proof.

## 5. Error/failure naming refinement

The proven caller-facing requirement now includes a stable machine-readable application identity (`reason`) in addition to protocol classification and safe presentation.

Target conceptual truth:
- `code` = standard gRPC protocol class;
- `reason` = stable low-cardinality application/business machine identity;
- `message` = public-safe human presentation.

Do not retain a second mandatory numeric business-code truth in parallel.

Preferred target vocabulary:
- ordinary Go technical/domain propagation remains ordinary `error`, `fmt.Errorf("%w")`, `errors.Is/As`;
- intentional caller-facing application failure uses one obvious application-failure constructor;
- naming target for implementation review is `failure.New(code, reason, safeMessage)` under a dedicated failure owner rather than a generic `errors` package plus aliased `ReturnError`.

The exact package rename/signature is implementation-pending and must be reconciled against live source/compatibility before mutation. The semantic requirement is durable now: one caller-facing failure owner, one constructor, no competing business-error frameworks.

## 6. Reason contract

A reason:
- is stable and machine-readable;
- is low-cardinality;
- is a static semantic identity, normally `UPPER_SNAKE_CASE`;
- describes the application/business fact, not microservice topology;
- contains no user/request IDs or dynamic values;
- is not a technical root cause;
- may be consumed by frontend behavior, tests, logs, metrics, support tooling and analytics.

Examples:
- `EXPORT_QUOTA_EXHAUSTED`;
- `PAYMENT_INSUFFICIENT_BALANCE`;
- `EMAIL_ALREADY_EXISTS`;
- `RATE_LIMIT_EXCEEDED`.

Avoid:
- `BUSINESS_ERROR`;
- `FAILED`;
- `ERROR_1007`;
- `PAYMENT_SERVICE_ERROR`;
- dynamic values such as `USER_123_QUOTA_EXCEEDED`.

## 7. Error/observability integration

One failure fact has multiple projections:
- response projection;
- structured log projection;
- metric projection;
- trace projection.

They must reuse the same canonical `code/reason`; they must not derive separate mappings.

Expected application rejection is not automatically an operational ERROR.

Unexpected technical failures remain wrapped Go errors. Their technical chain is preserved internally and logged once at the appropriate owner/boundary before public sanitization. Public technical reasons remain bounded/generic unless a real client consumer requires more detail.

The canonical request/outcome logging owner should record final outcome fields such as:
- request/operation identity;
- service/component;
- outcome;
- protocol code;
- application reason;
- duration.

Do not add `failure.NewAndLog`, `ErrorManager`, or a parallel error-logging framework.

## 8. One-name review examples

Converge public names when caller intent is one:
- `make logs` + long-lived `make log-query` -> target one public `make logs` surface while retaining distinct internal readers/query engine.

Split owners when invariants differ:
- `QuotaManager` -> pricing / reservation / durable ledger owners.

Converge competing semantic constructors:
- legacy `ReturnError`, `ThrowError`, outbound `fault.New`, direct business `status.Error`, service `*_faults.go` -> one intentional caller-facing failure constructor after bounded migration/proof.

## 9. Review discipline

Before adding a new abstraction, command, helper or exported name, apply:

`REALITY -> PROBLEM -> OWNER -> CONSUMER -> VALUE -> COST -> PROOF -> REMOVAL TEST`

Reject the addition when:
- an existing owner already expresses the same truth;
- it creates a synonym rather than a new capability;
- its consumer cannot be identified;
- its name hides multiple independent invariants;
- removing it would not remove a real capability.

## 10. Continuity

Future chats must read this doctrine before making new naming/public-surface decisions.

Do not reopen superseded Error V1 / Kind / Definition / catalog / i18n designs merely to obtain naming structure. A new mechanism is justified only by a newly proven consumer or invariant.

FINAL ACCEPTED = NO.
