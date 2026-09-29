# BDSPro Greenfield Learning / Rebuild

Status: PROVISIONAL learning and architecture laboratory.

This subtree exists to rebuild engineering understanding from the smallest useful
Go program upward while using real BDSPro problems as the curriculum.

It is intentionally separate from legacy implementation areas. Its existence does
not mean this path is the final production topology.

## Learning direction

```text
Go program
  -> package / visibility
  -> values / pointers
  -> errors / context
  -> dependency / interface
  -> composition
  -> domain invariant
  -> application command
  -> transaction
  -> persistence
  -> test
  -> HTTP/gRPC boundary
  -> event / async mechanics
  -> module contract
  -> runtime extraction only when proven
```

## Rules

1. Start from a real problem; do not start from a framework or pattern name.
2. Write the smallest code that exposes one responsibility clearly.
3. Explain the Go mechanism and the business/operational reason together.
4. Record alternatives and trade-offs before promoting an abstraction.
5. Keep business concepts out of generic core.
6. Treat local calls as default until process/network separation proves value.
7. Keep PostgreSQL as the initial durable-truth hypothesis; add Redis/broker only
   when a workload proves the need.
8. Tests and proof are part of learning, not an afterthought.
9. Every material decision updates the durable decision/checkpoint system.
10. Do not migrate legacy code into this subtree merely to make it look complete.

## Current starting point

Before adding architecture abstractions, align the preferred local workspace to the
live `architecture/rebuild-v2` branch.

Then begin at the current authorized design gate:

`R1.1 — What exactly is Core?`

The first exercises should stay small enough that every line and dependency can be
explained from first principles.
