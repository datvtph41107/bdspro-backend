# QHPRO Quota V2 — Group 1 Task 1.2 Implementation Report

## Status

```text
Contract design                 PASS
Backend request context         PASS
Gateway HTTP wiring             PASS
HTTP → gRPC metadata            PASS
Service context                 PASS
Service → downstream gRPC       PASS
Frontend ID generation          PASS
Auth retry preservation         IMPLEMENTED
Go requestcontext tests         PASS
Frontend identity typecheck     PASS
Go middleware tests             PENDING Go 1.24.1
Frontend Jest tests             PENDING node_modules
Final Decision Gate             OPEN until runtime tests pass
```

## Decision

`operation_id` identifies one logical business-operation instance.

It is not:

- a request attempt ID;
- an operation type/code;
- an idempotency guarantee;
- a quota charge key;
- a security credential.

## Contract

```text
HTTP:        X-Operation-ID
gRPC:        x-operation-id
Max length:  128 bytes
Safe chars:  A-Z a-z 0-9 - _ . :
Generated:   op_<32 hex chars>
```

## Retry semantics

```text
request attempt 1 → request_id A → operation_id X
request attempt 2 → request_id B → operation_id X
```

Frontend creates the ID before the first request and preserves it in the auth-refresh queue. Gateway generates only a backward-compatible fallback.

## One Source of Truth

```text
Format/generation/context → shared/common/requestcontext
HTTP ownership            → frontend request client + Gateway boundary
Transport mapping         → shared/common/middleware
Business operation type   → future Operation Registry
Idempotency enforcement   → future Task 1.3+
```

## Source audit findings

The following existing names are not system `operation_id`:

- Swagger `operationId`: documentation endpoint name.
- TQD `ClientRequestID`: local generated-report idempotency field.
- Core `trace_id`: technical telemetry.
- CRM/payment `transaction_id`: business record identifier.

## Backend files

```text
shared/common/requestcontext/identifier.go
shared/common/requestcontext/request_id.go
shared/common/requestcontext/operation_id.go
shared/common/requestcontext/operation_id_test.go
shared/common/middleware/grpc_metadata_interceptor.go
shared/common/middleware/grpc_metadata_interceptor_test.go
shared/common/middleware/client_log.go
shared/common/middleware/client_log_test.go
gateway-service/cmd/http/main.go
docs/quota/group-1-task-1-2-operation-id-contract.md
```

## Frontend files

```text
src/Services/requestIdentity.ts
src/Services/request.ts
src/Services/__tests__/requestIdentity.test.ts
src/Services/__tests__/request.operationId.test.ts
```

## Comment convention

New or changed contract comments use concise blocks:

```text
/**
 * Description.
 */
```

## Validation evidence

Available in the analysis environment:

```text
Static structure/wiring                 PASS
Go formatting                           PASS
GO111MODULE=off requestcontext tests     PASS
requestIdentity isolated TypeScript check PASS
Middleware tests                         SKIP: Go 1.23.2, repo requires 1.24.1
Frontend Jest tests                      SKIP: node_modules absent
```

## Local validation

```bash
ALLOW_SKIPS=0 bash validate-group1-task1-2-operation-id.sh \
  ~/projects/bdspro-golang-microservices \
  ~/projects/react-native-quy-hoach
```

Final expected output:

```text
PASS: Group 1 Task 1.2 operation_id contract and propagation
```

## Manual smoke tests

Preserve a client ID:

```bash
curl -sS -D - -o /dev/null \
  -H 'X-Operation-ID: parcel-preview-001' \
  http://127.0.0.1:8000/v2/tqd/public/enums \
  | grep -i '^x-operation-id:'
```

Generate a fallback:

```bash
curl -sS -D - -o /dev/null \
  http://127.0.0.1:8000/v2/tqd/public/enums \
  | grep -i '^x-operation-id:'
```

Replace an unsafe value:

```bash
curl -sS -D - -o /dev/null \
  -H 'X-Operation-ID: unsafe operation id' \
  http://127.0.0.1:8000/v2/tqd/public/enums \
  | grep -i '^x-operation-id:'
```

## Next gate

Task 1.3 must not start until full runtime validation passes. Its likely subject is `idempotency_key`, but its contract must be analyzed separately from `operation_id`.
