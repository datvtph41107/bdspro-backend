# Recovery Scenarios

The system must survive loss of chat context, account, machine and interrupted work.

## Scenario A — New Chat / Context Window Exhausted

Do not reconstruct from memory.

1. Open repository branch `architecture/rebuild-v2`.
2. Read `ARCHITECTURE-REBUILD-V2.md`.
3. Follow the exact recovery order.
4. Resolve current branch HEAD.
5. Compare live source/proof with `CHECKPOINT.md`.
6. Continue only from the current authorized gate.

## Scenario B — Different ChatGPT Account

Assume there is zero usable memory.

Provide only:

- repository: `datvtph41107/bdspro-backend`;
- branch: `architecture/rebuild-v2`;
- instruction: read `ARCHITECTURE-REBUILD-V2.md` and follow it exactly.

The repository must contain enough durable state for recovery without the old
account.

## Scenario C — New Machine / Fresh Clone

```bash
git clone git@github.com:datvtph41107/bdspro-backend.git
cd bdspro-backend
git fetch --all --prune
git switch architecture/rebuild-v2
git status
git log -1 --oneline
```

Read the durable recovery documents before running broad setup or mutation.

## Scenario D — Existing Machine, Separate Architecture Workspace

Prefer a worktree to avoid contaminating another active BDSPro workspace:

```bash
cd ~/projects/bdspro-backend
git fetch origin
git worktree add ~/projects/bdspro-rebuild-v2 architecture/rebuild-v2
cd ~/projects/bdspro-rebuild-v2
git status
```

If the worktree already exists, inspect it; do not recreate/delete blindly.

## Scenario E — Uncommitted / Interrupted Work

Never run reset/clean/rebase/amend first.

Capture:

```bash
git status --short --branch
git diff
git diff --staged
git log --oneline --decorate -10
```

Classify each change:

- intended current-gate work;
- generated artifact;
- unrelated user work;
- unknown.

Preserve unknown/unrelated changes until ownership is understood.

Update WORKLOG/CHECKPOINT before any destructive action.

## Scenario F — Remote Branch Advanced Since Checkpoint

Live immutable Git wins.

1. inspect new commits;
2. classify source/doc/proof changes;
3. compare them to the durable checkpoint;
4. update durable docs to reflect verified reality;
5. do not silently pretend the old checkpoint is current.

## Scenario G — Durable Docs Claim More Than Source/Proof

Source and completed exact-SHA proof win.

Downgrade the checkpoint, record the contradiction and repair documentation.
Never use prose to override executable reality.

## Scenario H — Source Changed but Proof Is Missing

State is CANDIDATE/UNPROVED, not CLOSED.

Record:

- exact SHA;
- what changed;
- missing proof;
- failure/unknown risk;
- next proof command.

Do not advance the next gate.

## Scenario I — CI Failed

Do not rewrite history or broaden scope.

Keep the failing SHA identifiable. Record the failing proof and root cause.
Create the smallest corrective commit and re-run proof on the new exact SHA.

## Scenario J — External Research Conflicts with BDSPro Needs

External code is evidence only.

Record:

```text
OBSERVATION
BDSPro REQUIREMENT
VALUE
COST
DIFFERENCE
ADOPT / ADAPT / REJECT / OPEN
```

Do not copy a Laravel/Java/OSS pattern merely because it is established elsewhere.

## Scenario K — Need to Change a CLOSED Decision

Do not silently edit history.

1. identify the decision ID;
2. introduce new evidence;
3. mark old decision SUPERSEDED;
4. create a new decision ID;
5. record migration/compatibility consequences;
6. reopen only the required gate.

## Scenario L — Prompt/Token Limit While Mid-Analysis

Before continuing elsewhere, publish a recovery checkpoint containing:

- exact SHA/worktree identity;
- what was inspected;
- partial conclusions labeled as such;
- hypotheses not yet proved;
- files changed;
- proof already completed;
- next exact action.

Do not convert partial reasoning into a CLOSED decision merely to create a neat handoff.

## Scenario M — Catastrophic Loss of Conversation History

Git is sufficient by design.

Minimum recovery token:

```text
Continue BDSPro Architecture Rebuild V2 from durable state.
Repository: datvtph41107/bdspro-backend
Branch: architecture/rebuild-v2
Read ARCHITECTURE-REBUILD-V2.md and follow it exactly.
```

No previous chat transcript should be required.
