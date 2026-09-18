# BDSPro — Development Runtime + Logging Operating Contract Checkpoint

Updated: 2026-09-19 Asia/Bangkok
Status: DESIGN DISCUSSION CLOSED FOR THIS SLICE; IMPLEMENTATION NOT YET AUTHORIZED

This document captures the durable architecture reasoning agreed before semantic-error design. It is intentionally stronger than conversational memory. Live Git/source still outranks this text for current implementation facts.

## 1. Governing principle

One truth per concern. One owner per truth. Multiple projections/views are allowed, but projections must not redefine the truth.

Use this reasoning chain for every mechanism:

REALITY -> PROBLEM -> OWNER -> CANONICAL TRUTH -> PROJECTION -> CONSUMER -> VALUE -> TRADE-OFF

Do not let PID files, ports, terminals, JSONL files, HTTP bodies, log levels, Docker images, or helper commands independently redefine the same fact.

## 2. Hierarchy of runtime concepts

Repository -> Service -> Process/runtime role -> Component/actor -> Work item/connection/request.

These are distinct scopes:
- Repository: source collaboration/change boundary.
- Service: business/runtime ownership boundary.
- Process/runtime role: OS execution/lifecycle boundary.
- Component/actor: goroutine/worker/listener/scheduler inside one process when lifecycle is shared.
- Work item/connection/request: one execution unit handled by a component.

Do not promote a worker, socket connection, scheduler tick, or internal actor into a separate process/run unless it truly has an independent OS process/deployment/scaling/lifecycle boundary.

## 3. Local development operating model

### `make up`
Meaning: establish the whole local development runtime from current source.

Expected intent:
- validate environment/toolchain;
- start long-lived backing infrastructure in Docker (PostgreSQL/PostGIS, Redis, RabbitMQ);
- apply required migrations/development bootstrap as owned by the current root orchestration contract;
- build all native application services from current source;
- start application services as background native processes;
- wait for startup/readiness evidence owned by the local runtime tooling;
- return the shell to the developer; no per-service foreground terminal is required.

`make up` is NOT the full application-container development loop. Full Compose remains explicit image/package/integration proof.

### `make dev service=X`
Meaning: transfer exactly service X from supervised background ownership to focused foreground hot-reload ownership.

Expected intent:
- stop only the existing supervised process for X;
- leave all other services running;
- run X with Air/foreground development in the current terminal;
- save -> rebuild -> restart X only;
- terminal shows live Air/build/runtime process output;
- local process output is mirrored to the same repository-owned temporary process-log model used by supervised mode;
- canonical structured application logging remains unchanged between `make up` and `make dev`.

Changing process owner must not change config semantics, application logging schema, correlation semantics, or business behavior.

### Process ownership state
A service process should have one local runtime ownership state at a time:
- SUPERVISED: root local runtime owns the process.
- DEV: foreground Air/developer session owns the process.
- DOWN: no owned process is running.
- FOREIGN: expected port/process resource is occupied by something not owned by the BDSPro local runtime.

`make status` should eventually distinguish runtime ownership from readiness. PID alive, port open, and service ready are different facts.

## 4. Backing dependencies

Backing infrastructure and application services are different dependency classes.

Example: Payment currently owns `DEV_DEPENDENCIES := postgres rabbitmq`. Service-local `deps-up` uses shared `deps.mk` to start only those long-running infrastructure dependencies. Migration is intentionally separate because it is a one-shot state transition, not a healthy long-running dependency.

Canonical ownership:
- Service Makefile declares which backing dependencies it needs.
- Shared `deps.mk` owns how those dependencies are started/stopped.
- Root command may be a thin convenience router, but must not duplicate the dependency list.

Desired developer equivalence:
- from root: `make deps service=payment-service` then `make dev service=payment-service`;
- from service: `cd payment-service && make deps-up && make dev`.

These are two entrypoints into the same underlying service-owned dependency truth, not two implementations.

`make dev` should not silently auto-start backing dependencies because developers may intentionally test dependency-down behavior. It should fail/help clearly when required infrastructure is unavailable.

Application dependencies such as User/Auth/Notification are NOT `deps-up`; they are other application services. Use `make up` when a full application graph is needed.

## 5. Physical process-output model

stdout and stderr are OS process streams.

Supervised mode (`make up`):
- service stdout/stderr are redirected to a repository-owned local process log;
- service runs in background;
- root command can return while service remains running.

Focused dev mode (`make dev service=X`):
- Air + child process run in foreground;
- stdout/stderr are visible in the current terminal;
- target design mirrors the same Air/process stream to the repository-owned process log (conceptually `tee`, but implementation must preserve signals/exit status correctly).

Thus one process stream has multiple views:
- foreground terminal (DEV mode);
- persisted local process log;
- `make logs` tail/follow view.

This is not duplicate business/application truth.

## 6. Raw process log vs structured application log

Do not collapse these.

### Raw process log
Question answered: "What did this local execution print?"

It may contain:
- Air build/reload messages;
- compiler errors;
- stdout/stderr;
- legacy `log.Printf`/`fmt` output;
- third-party library output;
- Go runtime panic stderr;
- console rendering of canonical slog events.

Recommended terminology/path direction: repository-owned process output, e.g. `.tmp/development/.../<service>/process.log` or an equivalent simple per-service layout. The exact physical layout is an implementation detail; developer commands are the public contract.

### Structured application logs
Question answered: "What canonical application event happened?"

Owned by `shared/common/logging` / `log/slog` target direction.

They provide machine-readable structured records, service identity, correlation, redaction, event fields, and optional local JSONL projections.

A single `slog` event may appear both as console text and structured JSONL. This is one event with multiple projections, not two event truths.

## 7. `.tmp` semantics

`.tmp` is ephemeral local development/runtime state, not business truth and not production retention.

Desired principle:
- one repository-owned `.tmp` root;
- working-directory differences must not create unrelated service-local logging roots;
- developer normally uses commands (`make status`, `make logs`, future structured query command) rather than navigating `.tmp` manually;
- deleting `.tmp` must not delete durable business state.

Do not overdesign the public filesystem contract. The exact subdirectory layout can evolve. The stable developer experience matters more than teaching paths.

A useful simple internal grouping is per-service artifacts under one repository root, with process output and structured runs distinguishable. Build/Air artifacts may also live under the same repository-owned temporary root.

## 8. `make logs`

For the current developer mental model, `make logs [service=...]` should mean process stdout/stderr view, analogous to `docker compose logs`:
- works when service is SUPERVISED;
- should also work when service is DEV because the foreground stream is mirrored to the same local process-log model;
- does not pretend to be the structured query engine.

Structured correlation/search can use a separate explicit view/command later (for example a structured/events view). Do not force every developer to learn JSONL paths for normal work.

## 9. Run identity

Current `shared/common/logging` session mode generates a run ID when a process configures the logger unless `QHPRO_LOG_RUN_ID` is supplied. Current generation is UTC timestamp + PID. It is used to create the structured run directory.

Meaning:
- `service.instance.id`: concrete process instance identity.
- `run_id`: one application logger/process run/generation.
- `request_id`: one HTTP request occurrence.
- `operation_id`: logical operation when applicable.
- `trace_id`: distributed execution identity when tracing is introduced.

A run ID is not a public API identifier and is not normally returned to the client. A process restart/Air child restart should create a new run ID. Old run directories remain as local historical evidence until local cleanup/retention removes them.

If one Air session restarts the service many times:
- process stream may be continuous across the Air session;
- each actual service process generation gets a new structured run.

Do not reuse run ID to group an Air terminal session. If such grouping ever has a proven consumer, create a separate dev-session identity.

## 10. Components/workers/listeners inside a process

A service can own multiple long-lived actors without creating multiple processes.

Current Payment composition root is the model:
- one Payment OS process;
- one PostgreSQL pool;
- outbound clients/resources;
- gRPC server;
- dashboard jobs;
- notification worker;
- fulfillment worker;
- outbox publisher;
- shared cancellation/shutdown lifecycle via one composition root/errgroup.

These actors share:
- PID/process instance;
- run ID;
- process stdout/stderr stream;
- process log;
- structured run directory.

Structured events distinguish internal ownership using fields such as component, worker ID, job/event/payment ID, etc. Do not create separate process files/run IDs just because multiple goroutines exist.

Current Relay/WebSocket pattern is similar:
- one Relay process;
- HTTP/WebSocket listener actor;
- Redis subscription actor;
- many WebSocket connections.

A WebSocket connection is a connection/session scope, not a process. Use connection/session identity only where useful. Do not create a run or file per connection.

## 11. Independent runtime roles

Split an internal actor into an independent runtime role/process only when real value requires it, e.g.:
- independent scaling;
- independent deployment/lifecycle;
- different runtime dependencies;
- different resource/security profile;
- failure of the actor should not share process lifecycle;
- separate availability/restart requirements.

If Payment fulfillment later needs 20 worker replicas while API needs 2, a future design could use the same Payment image with separate commands/roles (`grpc`, `fulfillment-worker`, etc.) or separate images if runtime dependencies warrant it.

Only then should each role have its own process instance/run/process stream. `service.name` can remain payment-service while `process.role` differentiates roles.

Do not introduce role complexity before source/runtime requirements prove it.

## 12. Worker/work-item semantics

Long-lived worker identity and work-item identity are separate:
- worker ID identifies the actor instance;
- event/job/message/payment ID identifies work being processed.

A job failure is not automatically a worker failure; a worker failure is not automatically a process failure. Preserve scope:
- work-item outcome;
- worker health;
- process health.

Similarly, WebSocket connection failure != listener failure != process failure.

## 13. Error/recovery relation (pre-semantic-error boundary)

This checkpoint intentionally stops before final semantic-error design, but establishes the observability substrate:
- compile/build failures live primarily in dev/process stream;
- startup failures must remain visible in process output and structured run if logger initialized;
- expected validation/business rejection is not automatically an ERROR log;
- dependency failures should emit structured owner evidence before public sanitization loses technical cause;
- retries are events/policy, not automatically final errors;
- panic/internal failures need rich operator evidence while public response remains safe;
- request correlation should allow Gateway -> service -> dependency reconstruction;
- background work uses operation/job/event identity when no request ID exists.

Do not log the same root failure as ERROR at every layer. Wrapping/classifying and logging are different responsibilities.

## 14. Local -> integration -> remote environments

Local daily development:
- Docker for backing infrastructure;
- native application processes;
- `make up` for whole-system supervised runtime;
- `make dev service=X` for one foreground hot-reload service;
- local `.tmp` artifacts for convenience/debugging.

Image/package integration:
- explicit full Compose integration topology;
- per-service images;
- shared source is build-time dependency, not a shared runtime container.

Remote development/staging/production:
- immutable per-service container images;
- no Air;
- no repository-local `.tmp` contract;
- structured stdout/telemetry collected by platform;
- production orchestrator owns restart/resource/health behavior.

Local and production should share application/config/error/log semantic contracts, not necessarily identical process supervision behavior.

## 15. Developer-experience principles

Public mental surface should stay small. A new developer should quickly answer:
- setup: `make setup`, `make doctor`;
- whole system: `make up`;
- state: `make status`;
- focused service: `make dev service=X`;
- backing infra for one service: `make deps service=X` or service-local `deps-up`;
- process output: `make logs [service=X]`;
- migration/test/integration through explicit commands;
- teardown: `make down`;
- image proof: explicit integration command.

Progressive disclosure:
1. Does it run? -> status.
2. What is it printing? -> terminal/process logs.
3. Which request/operation failed? -> correlation identity.
4. Which structured events belong to it? -> structured observability.
5. Which owner/dependency caused it? -> component/dependency evidence.
6. Which source logic is wrong? -> code/stack/debugger.

Do not require new developers to understand PID files, Air internals, JSONL paths, Docker topology, or logging handlers before they can change one endpoint.

## 16. Acceptance scenarios for future implementation

Implementation should not be considered closed until behavior is proven across at least:
- fresh clone -> setup/doctor;
- first `make up` -> whole stack background and terminal returned;
- startup failure -> service and evidence clearly identified;
- `make status` -> correct SUPERVISED/DEV/DOWN/FOREIGN semantics plus readiness distinction;
- whole-system and single-service `make logs`;
- `make dev service=X` stops only X supervised process;
- Air rebuild/restart visible in terminal and mirrored process log;
- second `make logs service=X` sees the same dev process stream;
- service process restart -> new structured run ID, old run retained;
- Ctrl+C focused dev -> X down, others remain;
- second DEV owner attempt -> fail clearly rather than hidden bind conflict;
- missing backing dependency -> actionable failure, no hidden auto-start;
- request correlation across Gateway/service;
- DB/provider/retry/panic/background-worker cases remain diagnosable;
- full image integration remains separate from daily native loop;
- production uses immutable service images/platform log collection.

## 17. Current implementation facts vs target

Current source already proves:
- root `make up` maps to native-up;
- native-up builds application services and starts native process graph;
- backing dependencies use Docker;
- native supervisor redirects stdout/stderr to root `.tmp/development/logs/<service>.log`;
- root `make dev service=X` stops the supervised X then runs service Air in foreground;
- service-local `deps-up` is owned by shared deps mechanics and service declaration;
- common logging supports structured console + JSONL outputs and per-run session directories;
- Payment currently runs multiple long-lived actors in one process;
- Relay currently runs WebSocket/HTTP + Redis actor in one process lifecycle.

Current gaps relative to the agreed target include:
- foreground `make dev` is not yet mirrored to the same root process-log view;
- structured log root is relative by default and can depend on service working directory;
- local status does not yet model DEV as a first-class ownership state;
- repository logging adoption is transitional (canonical slog and legacy logging coexist);
- `.tmp` naming/layout is not yet normalized to the final simplified developer model;
- structured query UX is not yet a public stable command.

Do not implement these gaps until semantic-error/code-organization discussion is closed and implementation is explicitly resumed from the durable checkpoint.

## 18. Next discussion gate

Next design topic: semantic errors and code organization.

The runtime/logging decisions above are constraints for that discussion:
- semantic error identity must not depend on process mode (`make up` vs `make dev`);
- public response and operational log are two projections of one outcome;
- event_name and semantic error code are distinct dimensions;
- technical cause stays internal;
- correlation must connect public occurrence to internal owner evidence;
- business/application failure scope must not be confused with worker/process/system health.


## 19. 2026-09-19 public naming refinement

R1-R4 implementation/proof history is not reopened by this refinement. R4's structured query capability remains valid and proved.

The public developer vocabulary is refined by the naming/ownership doctrine:

`ONE TRUTH PER CONCERN -> ONE OWNER -> ONE PUBLIC NAME -> ONE OBVIOUS ENTRY POINT`

Raw process output and structured application evidence are still different internal mechanics/truths. Do not collapse their storage/semantics.

However, the developer's top-level intent is one: inspect logs/evidence. Therefore the long-lived target public surface should converge on:

```bash
make logs
make logs service=payment-service
make logs request_id=req_...
make logs operation_id=op_...
make logs reason=EXPORT_QUOTA_EXHAUSTED
```

Internal routing may select:
- raw process-log tail/follow when no structured filters are supplied;
- the existing structured query engine when correlation/semantic filters are supplied.

The existing R4 `make log-query` command is not retroactively invalid. It is a proved compatibility/public alias until a bounded cleanup slice:
- updates docs/tests/help;
- proves `make logs` preserves both use cases;
- retires/deprecates the extra public name safely.

Do not delete the query engine. Do not merge raw and structured storage. Converge public vocabulary, not internal truth.

Error/Response integration adds a stable low-cardinality application `reason` as a structured query dimension. The same reason must originate from the canonical application failure truth and be projected into response/logging/metrics without message parsing or a second mapper.

FINAL ACCEPTED remains NO.
