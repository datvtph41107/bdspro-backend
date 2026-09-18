# BDSPro Canonical Error / Response Architecture — FINAL TARGET

Updated: 2026-09-19 Asia/Bangkok
Status: DESIGN REFINED — IMPLEMENTATION PENDING

This document supersedes `ERROR-RESPONSE-ARCHITECTURE-V1.md` for new implementation decisions. `ERROR-I18N-OPERATING-MODE.md`, `ERROR-SEMANTIC-CATALOG-CHECKPOINT.md`, and `ERROR-RESPONSE-ARCHITECTURE-V1.md` remain historical reasoning/audit evidence only.

2026-09-19 refinement: a real consumer/value has now been identified for a stable machine-readable application `reason`: frontend branching/action, structured log query/support correlation, low-cardinality metrics and cross-projection consistency. The target therefore evolves from `code + safe message` to `code + reason + safe message` without reviving Kind/Definition/catalog/i18n or a second numeric business-code truth.

Live Git/source and completed exact-SHA proof outrank durable prose for current implementation facts.

## 1. Final objective

BDSPro must end with ONE canonical outbound/application error core for all services.

The target is not a catalog, registry, Definition hierarchy, custom Kind taxonomy, or per-service error framework.

Canonical flow:

`SOURCE DECISION -> ONE CORE ERROR -> ONE gRPC NORMALIZATION -> ONE HTTP PROJECTION -> CLIENT`

The same execution context projects structured observability; response/logging are projections and must not redefine the failure truth.

Governing rule:

`ONE FAILURE FACT -> ONE SOURCE OF TRUTH -> MULTIPLE PROJECTIONS`

## 2. Canonical source of truth

The canonical application-failure truth is declared at the occurrence where source code intentionally decides the failure.

Example:

```go
if order == nil {
    return _errors.ReturnError(
        codes.NotFound,
        "ORDER_NOT_FOUND",
        "Không tìm thấy đơn hàng",
    )
}
```

That occurrence contains the minimum facts required by current consumers:
- standard gRPC `codes.Code`: protocol classification;
- stable low-cardinality `reason`: application/business machine identity;
- safe human message: public presentation.

`code`, `reason`, and `message` are different truths and must not redefine each other.

Do not hide these facts behind `Definition`, `Kind`, message catalogs, numeric business-code maps, or runtime registries. Reason governance should prefer static conventions + tests/audit rather than a second runtime catalog.

## 3. One core Error

Final semantic shape:

```go
type Failure struct {
    code    codes.Code
    reason  string
    message string
}

func New(code codes.Code, reason string, safeMessage string) error
func (e *Failure) Error() string
func (e *Failure) Code() codes.Code
func (e *Failure) Reason() string
```

Preferred implementation vocabulary for review is a dedicated owner such as `shared/common/failure` with `failure.New(...)`, because generic `errors` + aliased `ReturnError` obscures the semantic distinction from ordinary Go errors. The exact package rename/migration must still be proven against live compatibility before mutation. The durable semantic invariant is one caller-facing failure owner and one constructor.

The core intentionally does NOT own:
- HTTP status;
- request_id / operation_id / trace_id / run_id;
- service/component identity;
- stack/cause/technical error chain;
- custom numeric business code;
- a second/parallel custom business identity in addition to `reason`;
- custom Kind;
- Definition/catalog/message-key/locale;
- generic metadata map.

Each excluded field has another owner or lacks a proven consumer.

## 4. One write path for intentional outbound failures

For an intentional application/business/request/auth failure that must be returned to a caller, source uses exactly one constructor:

```go
return failure.New(codes.SomeCode, "SOME_STABLE_REASON", "Safe message")
```

Do not create canonical parallel constructors such as:
- `NotFound()`;
- `Conflict()`;
- `Unavailable()`;
- `Validation()`;
- `BusinessError()`;
- service-specific `*_faults.go` factories.

The code, reason and message remain visible at the call-site to optimize VS Code reading/search/debugging.

## 5. Why standard gRPC code belongs at the occurrence

`New(reason,message)` without protocol code is insufficient because a lower layer cannot determine whether the intended classification is NotFound, InvalidArgument, FailedPrecondition, PermissionDenied, ResourceExhausted, Unavailable, etc. without creating a hidden mapper or Definition registry.

A custom Kind is also unnecessary because gRPC already provides the transport-neutral RPC taxonomy used by service boundaries.

Therefore the minimal consumer-proven declaration is:

`failure.New(codes.Code, reason, safeMessage)`.

The gRPC code is NOT the business identity; it classifies protocol behavior. The reason is the stable application identity used when multiple failures can share the same protocol class or when clients/operations need machine-readable branching.

## 6. One stable application reason; no parallel numeric business truth

A real consumer requirement now exists for a stable application `reason`.

Reason rules:
- static, low-cardinality, normally `UPPER_SNAKE_CASE`;
- describes the application/business fact, not service topology;
- no request/user IDs or dynamic values;
- no technical dependency/root-cause text;
- stable enough for frontend branching, tests, structured log queries, support and low-cardinality metrics.

Examples:
- `EXPORT_QUOTA_EXHAUSTED`;
- `PAYMENT_INSUFFICIENT_BALANCE`;
- `EMAIL_ALREADY_EXISTS`;
- `RATE_LIMIT_EXCEEDED`.

Do NOT add a second mandatory numeric business ID/code in parallel. Legacy numeric compatibility may exist only as a bounded adapter with explicit retirement criteria.

Do NOT build a Kind/Definition/catalog hierarchy merely to store reasons. Prefer literals/conventions plus static audit/tests. Typed retry/action details remain separate extensions and require their own proven consumer.

## 7. Technical Go errors remain technical

Unexpected infrastructure/runtime failures remain ordinary Go errors:

```go
result, err := repo.Do(ctx)
if err != nil {
    return fmt.Errorf("do repository operation: %w", err)
}
```

Do NOT convert technical failures into:

```go
ReturnError(codes.Internal, err.Error())
```

Technical error chains are internal diagnostic truth. They feed structured observability and are sanitized at the outer normalization boundary.

## 8. Domain sentinels remain only when an internal consumer needs identity

`errors.New` sentinels and `errors.Is` remain valid for internal/domain communication when they buy value.

Example:

```go
if errors.Is(err, domain.ErrOrderNotFound) {
    return failure.New(
        codes.NotFound,
        "ORDER_NOT_FOUND",
        "Không tìm thấy đơn hàng",
    )
}
```

Do not introduce a sentinel merely as an intermediate hop when the usecase already knows the final condition directly.

## 9. One gRPC normalization boundary

Business/usecase code must not construct `status.Error/status.New` directly for application failures.

At the outer gRPC boundary/interceptor, normalize once:
- canonical application `*failure.Failure` -> its `Code()` + `Reason()` + safe message;
- `context.Canceled` -> `codes.Canceled`;
- `context.DeadlineExceeded` -> `codes.DeadlineExceeded`;
- unknown technical error -> structured technical evidence + `codes.Internal` + generic safe message.

This is ONE normalization path for the repository.

Protocol/framework-specific cases may enter this same normalization boundary, but they must not create a second application-error model.

## 10. One HTTP projection in Gateway

Gateway owns deterministic protocol/public projection only:

`gRPC status code -> HTTP status`

and carries the already-canonical application reason/message into the public response. Gateway must not invent or remap business reasons.

Gateway must not know Payment/User/TQD business vocabularies and must not map human message text to status.

A central mapping may project canonical gRPC classes to HTTP classes, e.g. InvalidArgument->400, Unauthenticated->401, PermissionDenied->403, NotFound->404, AlreadyExists->409, ResourceExhausted->429, Unavailable->503, Internal->500. FailedPrecondition requires one documented repository policy and must not vary by service.

The Gateway then serializes one small public response contract based on actual client needs.

## 11. Public response V1

Default target body:

```json
{
  "reason": "ORDER_NOT_FOUND",
  "message": "Không tìm thấy đơn hàng"
}
```

HTTP status is transport truth and is not duplicated in the body without a consumer requirement.

`X-Request-ID` is the default correlation surface for clients. Do not put stack, DB/provider text, technical cause, generic metadata, service internals, or duplicated status fields in the public body.

## 12. Validation/auth/security use the same core representation

Representation is one even when ownership differs.

Examples:
- request/input validation -> `ReturnError(codes.InvalidArgument, safe message)`;
- authentication -> `codes.Unauthenticated`;
- authorization -> `codes.PermissionDenied`;
- business state rejection -> appropriate standard RPC code such as `FailedPrecondition`;
- resource absence -> `NotFound`.

Ownership of deciding the condition remains with the correct layer. The representation/projection mechanism remains shared and singular.

Structured field violations or retry/action details are added only if a real consumer requires them, as typed extensions to the same protocol—not as a second error framework.

## 13. Observability is a projection, not another truth

The error core does not store correlation/runtime fields.

Structured logging combines:
- the canonical application code/reason/message OR technical Go error;
- request/operation context;
- service/component identity;
- runtime/run identity.

Expected business rejection is not automatically an operational ERROR.

Technical root failure should be logged once at the owner/boundary with enough context; do not log the same root cause as ERROR at repository, usecase, handler, and Gateway simultaneously.

Public sanitization must occur after technical evidence has been preserved.

## 14. Local Runtime/Logging integration

Local `.tmp` is repository-owned ephemeral evidence.

Developer workflow must converge on:

`CLIENT RESPONSE -> X-Request-ID/reason -> make logs filters -> service/component/dependency evidence -> VS Code source`

Raw process output and structured retained records remain distinct internal mechanics/truths. The target daily public vocabulary is one `make logs` entry point with filters/options; the proven R4 query engine may remain internal/compatibility implementation during migration. Do not preserve a second long-lived public `make log-query` name solely because the implementation path differs.

Error/Response acceptance therefore includes response-to-evidence correlation proof; it is not accepted merely because `ReturnError` compiles.

## 15. VS Code / human DX acceptance

A developer should read:

```go
return failure.New(
    codes.FailedPrecondition,
    "PAYMENT_INSUFFICIENT_BALANCE",
    "Số dư không đủ để thực hiện giao dịch",
)
```

and immediately know the protocol class and public-safe meaning without jumping through a Definition/catalog chain.

Expected IDE workflow:
- `Ctrl+Shift+F` message -> real decision occurrences/tests;
- search `failure.New(` / reason literal -> intentional outbound failures;
- search `codes.NotFound` etc. -> protocol-class usage;
- F12 on `failure.New` -> exactly one caller-facing failure implementation;
- request ID from client/runtime -> structured evidence -> source/component.

Architecture optimizes:
`READ -> UNDERSTAND -> SEARCH -> TRACE -> FIX -> TEST`.

## 16. Final source organization

Preferred shared target:

```text
shared/common/failure/
  failure.go   # one Failure + New(code, reason, safeMessage)
  grpc.go      # one normalization/projection path
```

Exact physical rename is migration-reviewable; do not create a second package while the legacy package remains equally valid.

Logging remains a separate mechanics concern under the canonical logging owner.

Gateway owns one response/projection location for gRPC -> HTTP + serialization + request-id header behavior.

Services do not create default `apperror/`, `*_faults.go`, `definition.go`, `catalog.go`, or business status mappers. Domain `errors.go` files exist only where internal sentinels genuinely have consumers.

## 17. Target deletion/retirement set

After usage reaches zero and compatibility gates allow removal, retire/delete active use of:
- `shared/common/fault/**` as an outbound business model;
- `fault.Kind` and Kind mappers;
- `_errors.ThrowError` when redundant;
- direct business/usecase `status.Error/status.New`;
- service `*_faults.go` business factories/mappers;
- custom numeric business-code registries as canonical truth;
- active `sharepb.ErrorResponse` legacy response path;
- oversized `httpresponse.Problem` fields/logic that no real client requires;
- public `err.Error()` leakage;
- duplicate message/code/status mapping tables.

`shared/protobuf/**` remains NO-TOUCH; a protected historical protobuf type may remain physically present while runtime usage becomes zero.

## 18. What is NOT duplicate debt

Do not delete legitimate Go mechanics:
- ordinary `error`;
- `fmt.Errorf("...: %w", err)`;
- `errors.Is` / `errors.As`;
- context cancellation/deadline errors;
- domain sentinels with real internal consumers.

These are internal error propagation/identity mechanics, not competing outbound response architectures.

## 19. Migration constraint and final signature

CURRENT legacy signature:

```go
ReturnError(code int32, message string)
```

FINAL semantic target:

```go
failure.New(code codes.Code, reason string, safeMessage string)
```

The implementation may require a bounded transition through the existing package/symbol, but the transition must converge to one public caller-facing constructor rather than leave aliases permanently.

The current integer is mixed legacy compatibility data and must not be reinterpreted as canonical business identity.

A bounded compatibility bridge may exist during migration, but it must have an explicit retirement condition. Do not use a permanent `...any`/type-unsafe overload.

`organization-service/**` is currently NO-TOUCH while it still contains legacy calls. Therefore final one-core cleanup requires an explicit later acceptance decision: either authorize a mechanical no-business-semantics migration of protected legacy call-sites, or acknowledge a bounded compatibility shim remains. Do not claim full one-core completion while both constraints conflict.

## 20. Migration method

For each service/scope:

`INVENTORY -> CLASSIFY OWNER -> CONVERT -> PROVE -> DELETE/RETIRE -> RATCHET`

Classify each current occurrence as:
- intentional outbound application failure -> canonical one-way `failure.New(codes.Code,reason,safeMessage)`;
- technical error -> ordinary wrapped Go error;
- domain internal condition -> sentinel only if consumed;
- context/protocol concern -> normalized by the proper boundary;
- legacy compatibility -> bounded outer adapter until consumer removal.

Do not globally replace all `status.Error` or `err.Error()` occurrences; source inventory proved they contain both legitimate and debt usages.

## 21. Acceptance scenarios

Error/Response is not accepted only by unit-testing Error methods. Prove vertical behavior including:

1. Expected business rejection:
   - correct gRPC/HTTP class;
   - safe message;
   - request ID correlation;
   - not falsely logged as system ERROR.

2. Technical DB/provider failure:
   - technical evidence preserved internally;
   - safe generic public response;
   - no internal leak;
   - same request/operation identity reconstructs timeline.

3. Request validation:
   - canonical InvalidArgument projection;
   - one public response mechanism.

4. Auth/security failure:
   - canonical protocol semantics;
   - no service-specific response framework.

5. Background job failure:
   - operation/job identity correlation;
   - job failure distinguished from worker/process failure.

## 22. Final invariant

A developer asking, “How do I intentionally return a caller-facing failure?” has exactly one answer:

```go
return failure.New(
    codes.SomeCode,
    "SOME_STABLE_REASON",
    "Safe message",
)
```

A developer asking, “How do I propagate an unexpected technical failure?” has exactly one answer:

```go
return fmt.Errorf("operation context: %w", err)
```

There is no third competing business-error architecture.

## 23. Protected acceptance constraints

Preserve until explicitly superseded:
- `shared/protobuf/**` NO-TOUCH;
- `organization-service/**` NO-TOUCH;
- `map-service/**` NO-TOUCH;
- `make wire`;
- `make buf`;
- byte-identical `shared/code/deploy.sh` with required SHA256;
- production-lineage separation.

FINAL ACCEPTED remains NO until source/runtime/E2E/recovery/release/promotion gates pass.


## 24. Naming and public-surface invariant

Error/Response follows `ARCHITECTURE-NAMING-OWNERSHIP-MINDSET.md`.

Rules:
- ordinary Go `error` and intentional caller-facing application `failure` are distinct concepts and should not be obscured behind one generic name;
- one application failure constructor is the only normal business/application write path;
- do not add `ErrorManager`, `FailureManager`, `NewAndLog`, per-service failure factories or synonym constructors;
- gRPC normalization, Gateway HTTP projection, structured logging and metrics remain different owners/projections, but none may create a second reason/code truth;
- reason is shared across public response and observability; technical root cause remains internal;
- target naming must make the ownership obvious at the call site.

FINAL ACCEPTED remains NO.
