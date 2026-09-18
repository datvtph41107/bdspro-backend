# BDSPro User-Lane Task Protocol

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: ACTIVE

Every task delegated from ChatGPT to the user's WSL machine must use this contract.

## Task block

```text
[USER-LANE]
ID: <stable short id>
MODE: READ_ONLY | WRITER | PROOF | RUNTIME
AUTHORITY_SHA: <sha>
WORKDIR: <path>
PURPOSE: <why>

PRECHECK:
  <commands/conditions>

COMMAND:
  <copy-paste block>

EXPECTED:
  <success facts>

STOP IF:
  <conditions requiring coordinator reconciliation>

RETURN:
  - exit code
  - relevant output
  - log path
  - git HEAD/status if source-related
```

## Mode rules

READ_ONLY:
- may search/test/inventory;
- must not edit/commit/push.

WRITER:
- only one active writer;
- only authorized paths/scope;
- stop on unexpected source reality.

PROOF:
- detached immutable candidate SHA;
- no source mutation;
- may run parallel proof lanes.

RUNTIME:
- may start local dependencies/services;
- must record exact SHA/config and preserve logs;
- runtime observations are evidence, not source authority.

## Failure rule

A failed command is evidence.

Do not repair by guessing.
Do not broaden scope.
Do not rerun with changed flags unless instructed.

Return the failure to the coordinator.

## Completion rule

A user-lane task is not a gate transition by itself.

Only the coordinator can combine:
- exact Git/source state;
- local evidence;
- hosted CI;
- durable phase rules;

and declare the next gate authorized.
