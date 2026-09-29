# Go Package Naming — First-Principles Research

Status: research foundation for R1.1. Not a final BDSPro package topology.

## Core theorem

In Go, the package name is part of the caller-facing API.

A caller imports by package path but normally refers to exported identifiers with
the package name:

```go
bytes.Buffer
http.Server
time.Duration
```

Therefore package naming and exported-symbol naming must be designed together.

If a package name cannot act as a meaningful, natural prefix for what it exports,
treat that as evidence that the package abstraction/boundary may be wrong.

## Language mechanics

The Go specification distinguishes:

- package path: used to locate/import a package;
- package name: declared by `package name`, used as the default local qualifier;
- exported identifiers: referenced as `pkg.Identifier` by clients.

A set of source files sharing the package name forms a package implementation,
normally within one directory. Import declarations establish package dependency
relations and cycles are illegal.

## Naming principles

Prefer names that are:

- lowercase;
- short but clear;
- usually one word;
- meaningful at call sites;
- unlikely to need import aliases;
- unlikely to steal a useful local variable name;
- focused on what the package provides rather than a generic technical bucket.

Avoid reflexive names such as:

- util / utils;
- common / shared;
- helper / helpers;
- misc;
- types;
- interfaces;
- generic api/model buckets when they do not express a real responsibility.

Pluralization is generally avoided when the singular concept reads naturally.

## Call-site test

Design from the consumer's point of view.

Prefer:

```go
http.Server
time.Now()
stringset.New(...)
organization.Create(...)
payment.Charge(...)
```

over repetitive/stuttering APIs:

```go
http.HTTPServer
organization.OrganizationService
payment.PaymentService
```

The package name already supplies context.

## Boundary test

Before creating a package, ask:

1. Can its responsibility be described in one sentence?
2. Does its proposed name make its exported API read naturally?
3. What belongs here?
4. What clearly does not belong here?
5. Who are its consumers?
6. What does it depend on?
7. Will consumers commonly import it together with another same-named package?
8. Would callers need aliases because the name collides with ordinary variables?
9. If the package were removed, what cohesive capability disappears?

If these are hard to answer, do not manufacture a package just to move files.

## Performance theorem

The spelling of a package name itself has no meaningful runtime-performance
benefit.

Performance effects come from package/dependency topology, not from choosing the
word `organization` instead of `org`.

However naming can indirectly influence build performance because vague packages
such as `util` and `common` tend to become dependency magnets. The official Go
package-name guidance explicitly warns that these packages can accumulate
dependencies and make compilation unnecessarily slower in large programs.

Go builds packages together with their dependencies, so a poor boundary can widen
the dependency graph and increase rebuild/analysis cost.

Therefore distinguish:

```text
PACKAGE NAME SPELLING
  -> readability/API effect
  -> essentially no runtime speed effect

PACKAGE BOUNDARY + DEPENDENCY GRAPH
  -> compilation/cache/rebuild/tooling effect
  -> architecture and maintainability effect
  -> can indirectly affect runtime through initialization/dependency choices
```

Never justify a package rename by claiming runtime speed unless measurement proves
a separate mechanism.

## Developer-performance value

A good package name reduces cognitive translation at:

- reading call sites;
- code review;
- autocomplete;
- go-to-definition/references;
- grep/search;
- debugging dependency direction;
- onboarding;
- deciding where new code belongs;
- refactoring;
- avoiding import aliases and name collisions.

The practical objective is not fewer characters. It is fewer decisions the reader
must reconstruct mentally.

## Community evidence

Community discussion is not uniform about every naming convention, especially
receiver/variable abbreviation. A recurring theme is that consistency and
readability matter more than treating conventions as dogma.

For package organization specifically, community discussions repeatedly surface:

- generic `utils` packages become grab bags;
- if a helper is used only locally, keeping it unexported in the consuming package
  is often simpler than inventing a shared package;
- if no meaningful package name is available, the abstraction may not be mature
  enough to extract;
- domain/responsibility-oriented organization can make related code easier to
  reason about, though teams differ on exact layout;
- import collisions from generic names such as `client`, `server`, or
  `errors` create aliasing and consistency costs when such packages meet at the
  same call site.

Use community feedback as experience evidence, not language law.

## BDSPro provisional policy

For Rebuild V2:

1. name packages from responsibility/capability and caller language;
2. design call sites before finalizing names;
3. do not create `utils`, `common`, `shared`, `types`, `interfaces` or
   generic `service` packages merely as storage;
4. keep one-off helpers local/unexported until a stable shared responsibility
   emerges;
5. avoid package names that force routine import aliases;
6. avoid stuttering exported names;
7. treat inability to name a package clearly as a boundary-design signal;
8. allow renaming while learning; early names are hypotheses, not sacred API;
9. benchmark/profile performance claims separately from naming/style decisions.

This policy remains PROVISIONAL until exercised by the first real BDSPro packages.
