# Durable Project Memory

This file stores durable context that a new session needs but that should not be
left only in conversational memory.

## Current intent

Architecture Rebuild V2 is both:

- an architecture reconstruction lane for BDSPro;
- a practical Go/backend learning corpus;
- a controlled refactor/rebuild laboratory;
- a source of durable engineering decisions and evidence.

It is not a tutorial project and not a pattern-copy exercise.

## Working mindset

Prefer:

- business value over architectural fashion;
- explicit ownership over shared ambiguity;
- module boundaries before network boundaries;
- simple local calls before remote calls;
- PostgreSQL durable truth before adding cache/broker authority;
- measured operational pressure before extraction;
- deterministic tests;
- intent-shaped commands over generic CRUD;
- vertical business slices over technical-type folder scattering;
- small, earned abstractions over speculative framework building.

## Known source debt motivating V2

After Source Migration V1, physical topology improved, but inner areas still
contain inconsistent organizations. For example User has overlapping concepts such
as domain/dto/models/interface/usecase/usecases/utils/infra/db/database while
Payment and Property use different inner structures.

Therefore V2 must define and prove a module contract before broadly migrating
inner source structure.

## Current architecture hypotheses — not yet final decisions

These are starting hypotheses to challenge with evidence:

```text
one repository
one Go module initially for the rebuild baseline
modular business source
manual composition initially
PostgreSQL as durable truth
local call by default
network call by proven pressure
no Redis by default
no message broker by default
no ORM by default
no framework abstraction by default
one executable initially, additional process roles only when justified
```

Each may be adopted, adapted or rejected.

## Initial module model hypothesis

A business module should make it easy to answer:

- what business fact does it own?
- what invariant does it enforce?
- what commands/queries does it expose?
- what state/write path does it own?
- what external capabilities does it require?
- what events does it create?
- how is it tested?

Candidate anatomy:

```text
module/
  domain/
  application/
  adapters (e.g. postgres, transport)
  module.go / composition
```

Repository interfaces belong with the consuming business/application need, not in
a global generic BaseRepository.

## Core principle

"Core" has two meanings:

1. Architecture Core: rules for how the project works.
2. Code Core: minimal business-agnostic primitives proven useful by real slices.

Do not allow business concepts into Code Core.

## Database principle

Do not start Database V2 from an ERD. Derive it from:

```text
business fact
 -> identity
 -> invariant
 -> lifecycle
 -> command
 -> transaction
 -> concurrency
 -> history
 -> query requirement
 -> database model
```

## TQD principle

Do not force TQD into a generic module shape when spatial workload proves special
needs. Standardize ownership/dependency/testing/composition, while allowing
PostGIS geometry/geography, spatial indexes, bulk work, workers, projections or
specialized caching when evidence justifies them.

## Learning principle

Learn Go/backend through real problems:

```text
Go concept
 + real BDSPro problem
 + architecture decision
 + production failure mode
 + proof
```

Examples:

- interface through Clock/Repository consumers;
- transaction through Organization + owner membership;
- idempotency through Payment;
- concurrency through Quota;
- asynchronous delivery through Notification/outbox;
- spatial design through TQD/PostGIS;
- authorization through Membership/Role/Permission.


## Preferred local development workspace

The user's primary WSL working copy for Rebuild V2 is:

`~/projects/bdspro-canonical-bootstrap-20260908`

A durable recovery must treat local state as something to reconcile, not assume.
The initial observed local snapshot before alignment to Rebuild V2 was:

```text
working directory: ~/projects/bdspro-canonical-bootstrap-20260908
HEAD branch: refactor/canonical-observability-errors-a6d0722a
HEAD commit: fca4682
```

The repository also contains multiple historical/local branches and at least one
branch checked out in another worktree. Therefore recovery must inspect
`git status`, `git worktree list`, branch tracking and the live remote before
switching or deleting anything.

Target local working branch for the rebuild program:

`architecture/rebuild-v2`

The `bdspro/` subtree on that branch is the greenfield learning/rebuild area. It
is a training and architecture laboratory until a later decision promotes any of
its structure into canonical production topology.
