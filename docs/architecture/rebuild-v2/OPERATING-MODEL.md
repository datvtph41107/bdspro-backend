# Operating Model

## 1. Fundamental doctrine

Architecture is derived from reality:

```text
REAL BUSINESS
  -> BUSINESS FACT
  -> INVARIANT
  -> LIFECYCLE
  -> FAILURE / CONCURRENCY
  -> OWNER
  -> MODULE
  -> TRANSACTION
  -> DATA MODEL
  -> PROCESS
  -> INFRASTRUCTURE
```

Do not start with framework, folder, database table, message broker or service count.

Core doctrine:

```text
ONE BUSINESS FACT
  -> ONE SEMANTIC OWNER
  -> ONE DURABLE AUTHORITY
  -> ONE CLEAR WRITE PATH
  -> MANY READERS / PROJECTIONS
```

Topology doctrine:

```text
BUSINESS TOPOLOGY != SOURCE TOPOLOGY != PROCESS TOPOLOGY != CONTAINER TOPOLOGY
MODULE BOUNDARY != PROCESS BOUNDARY != NETWORK BOUNDARY
WORKER != MICROSERVICE
ASYNC != MESSAGE BROKER
```

Default strategy:

```text
local-by-default
extract-by-proven-pressure
```

A split is justified only when concrete value exceeds distribution, migration,
operational and cognitive cost.

## 2. Architecture Core vs Code Core

Architecture Core is the long-lived rule set for:

- module ownership;
- dependency direction;
- command/query boundaries;
- transaction ownership;
- error propagation;
- event creation/delivery;
- configuration;
- composition;
- test/proof strategy;
- runtime extraction criteria;
- admission into shared/core.

Code Core is only business-agnostic primitives that have been proven necessary.

Expected initial candidates:

```text
core/
  identity/
  clock/
  id/
  transaction/
  event/
  errors/
```

Business terms such as Organization, Membership, Role, Property, Payment and Quota
must not enter generic core.

## 3. Decision lifecycle

Every architectural subject follows:

```text
QUESTION
  -> FACTS
  -> OPTIONS
  -> VALUE
  -> COST
  -> FAILURE MODES
  -> MINIMUM EXPERIMENT
  -> BDSPro EVIDENCE
  -> DECISION
  -> PROOF
  -> CHECKPOINT
```

Research must end in one of:

- ADOPT
- ADAPT
- REJECT
- OPEN

External projects are evidence/counterexamples, never BDSPro authority.

## 4. Change classes

Classify work before mutation:

- RESEARCH: no source behavior change.
- EXPERIMENT: isolated/reversible candidate, not canonical.
- ARCHITECTURE: changes durable source/ownership/dependency structure.
- BUSINESS: changes business semantics or state transitions.
- DATA: changes durable schema/migration/data ownership.
- RUNTIME: changes process/network/deployment topology.
- PROOF: only tests, verification or evidence.

Do not combine unrelated classes in one checkpoint without explicit reason.

## 5. Checkpoint closure rule

A checkpoint is CLOSED only when all applicable items exist:

```text
question resolved
facts recorded
alternatives considered
trade-off explicit
minimal implementation complete
tests/proof complete
failure case understood
BDSPro value demonstrated
exact source identity known
remaining risks/open questions recorded
next gate explicit
```

"Looks good" is not closure.

## 6. Session close protocol

Before ending any substantial work session:

1. reconcile live branch/worktree;
2. update `WORKLOG.md`;
3. update `DECISION-REGISTER.md` for changed decisions;
4. update `CHECKPOINT.md`;
5. update `MEMORY.md` only for durable lessons/context;
6. commit the durable state with a descriptive message;
7. record exact commit SHA and proof status;
8. leave the next authorized action explicit.

A session may end with OPEN work. Preserve it; never manufacture closure.

## 7. Source mutation guard

Before mutating source:

- identify current exact branch/SHA;
- determine the current authorized gate;
- state which invariant/problem the mutation addresses;
- preserve unrelated work;
- choose the smallest sufficient change;
- define proof before broadening scope.

Never perform cleanup merely because a path looks old.
