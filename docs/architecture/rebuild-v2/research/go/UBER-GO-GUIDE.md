# Uber Go Style Guide — Rebuild V2 Research Reference

Source: https://github.com/uber-go/guide
Observed default branch: `master`
Observed repository commit: `1d60a91aa5e87d443002e23c21903c49489dbde5`
Primary generated guide: `style.md`

## Role in Rebuild V2

This repository is a foundational Go engineering reference, not a BDSPro
architecture authority.

Use it when a real problem appears, not as a checklist to apply blindly.

Decision protocol for every borrowed guideline:

```text
BDSPro PROBLEM
  -> GO LANGUAGE / RUNTIME MECHANISM
  -> UBER GUIDE OBSERVATION
  -> CURRENT BDSPro CONTEXT
  -> VALUE
  -> COST
  -> ADOPT / ADAPT / REJECT / OPEN
  -> CODE / TEST PROOF
```

Do not create abstractions or source structure merely because the guide contains a
rule.

## Why this reference is useful

The guide explicitly covers more than formatting. Relevant areas include:

- pointers to interfaces;
- interface compliance;
- receivers and interfaces;
- zero-value mutexes;
- copying slices/maps at boundaries;
- `defer` for cleanup;
- channels and goroutine lifecycle;
- time handling;
- error types, wrapping, naming and handling;
- type assertion failure;
- panic avoidance;
- mutable globals;
- public struct embedding;
- `init()`;
- process exit ownership in `main`;
- struct/map initialization;
- naming and package conventions;
- scope and nesting;
- table-driven tests;
- functional options;
- linting and vetting.

The guide also points to Effective Go, Go Common Mistakes and Go Code Review
Comments as upstream/general references.

## Problem-indexed learning map

Do not read this guide linearly as a course. Pull the relevant section only when
our code reaches the matching problem.

| Rebuild problem | Guide areas to consult |
| --- | --- |
| package/API boundaries | package names, exported APIs, embedding |
| pointer vs value semantics | receivers and interfaces, struct references |
| interfaces | pointers to interfaces, compliance checks, receivers |
| errors | error types, wrapping, naming, handle once |
| resource lifecycle | defer, cleanup, exit in main |
| concurrency | mutex zero values, channels, goroutine lifecycle |
| mutable state | globals, maps/slices at boundaries |
| construction | struct/map initialization, functional options |
| maintainability | reduce nesting/scope, naming, consistency |
| testing | table tests |
| tooling baseline | gofmt/goimports, vet, linting |

## Rebuild V2 usage rule

For each relevant section:

1. first understand the Go mechanism from the language/runtime;
2. reproduce the problem with the smallest code;
3. read the matching Uber guideline;
4. explain why the guideline exists;
5. test whether the same failure/cost exists in BDSPro;
6. decide ADOPT / ADAPT / REJECT / OPEN;
7. record material decisions in the durable decision system.

This keeps the guide as a high-quality engineering reference while preserving
BDSPro-specific first-principles reasoning.
