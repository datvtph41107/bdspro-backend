# BDSPro Local Execution Runbook — Start From Zero

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: CANONICAL USER SETUP GUIDE

This is the step-by-step guide the user follows on Windows + WSL/Ubuntu. Do not skip verification steps.

The screenshot at adoption time showed:
- the user was already inside tmux session `bdspro`;
- shell path was `~/projects/bdspro-canonical-bootstrap-20260908`;
- a first draft of `~/bdspro-ops/run.sh` had been created;
- local repository identity/SHA had NOT yet been re-verified under this operating model.

Therefore setup restarts from verification, not from destructive cleanup.

## Phase 0 — Understand the three local areas

1. Anchor clone
   `~/projects/bdspro-canonical-bootstrap-20260908`

   Purpose: fetch Git and manage worktrees.
   After setup, do not code here.

2. Worktrees
   `~/bdspro-worktrees`

   Purpose: isolated reader/writer/proof working directories.

3. Operations evidence
   `~/bdspro-ops`

   Purpose: logs, proof output, helper scripts.

Do NOT delete the existing repository or existing ops directory just because setup is restarting.

## Phase 1 — Verify the anchor clone

Run exactly:

```bash
cd ~/projects/bdspro-canonical-bootstrap-20260908

printf '\n== PATH ==\n'
pwd

printf '\n== REPOSITORY ROOT ==\n'
git rev-parse --show-toplevel

printf '\n== REMOTES ==\n'
git remote -v

printf '\n== LOCAL STATUS ==\n'
git status --short --branch

printf '\n== CURRENT HEAD ==\n'
git rev-parse HEAD

printf '\n== FETCH ACTIVE AUTHORITY ==\n'
git fetch origin refactor/canonical-observability-errors-a6d0722a

printf '\n== REMOTE ACTIVE HEAD ==\n'
git rev-parse origin/refactor/canonical-observability-errors-a6d0722a
```

Expected remote active SHA for this checkpoint:

`fca4682587d93cbc436f2466244ee8bd03b8b1b9`

STOP after Phase 1 and send the complete output to the coordinator.

Do not create/reset worktrees until the coordinator confirms:
- correct remote;
- correct repository;
- no dangerous uncommitted state;
- remote SHA agrees with live GitHub.

## Phase 2 — Create persistent local directories

Only after Phase 1 confirmation:

```bash
mkdir -p ~/bdspro-worktrees
mkdir -p ~/bdspro-ops/logs
mkdir -p ~/bdspro-ops/proofs
mkdir -p ~/bdspro-ops/state
```

Meaning:
- `mkdir -p` creates the directory if missing;
- it does not delete existing contents.

## Phase 3 — Install/check tmux

Check:

```bash
tmux -V
```

If missing:

```bash
sudo apt update
sudo apt install -y tmux
```

List sessions:

```bash
tmux ls
```

If `bdspro` already exists, reuse it:

```bash
tmux attach -t bdspro
```

Otherwise:

```bash
tmux new -s bdspro
```

Do not kill an existing session merely to make setup look clean.

## Phase 4 — Replace run.sh with the canonical wrapper

After directories exist:

```bash
cat > ~/bdspro-ops/run.sh <<'EOF'
#!/usr/bin/env bash
set -uo pipefail

REPO="${1:?repo required}"
NAME="${2:?name required}"
shift 2

STAMP="$(date +%Y%m%d-%H%M%S)"
LOG="$HOME/bdspro-ops/logs/${STAMP}-${NAME}.log"

cd "$REPO" || exit 1

{
  echo "========================================"
  echo "name=$NAME"
  echo "started=$(date -Is)"
  echo "repo=$REPO"
  echo "head_before=$(git rev-parse HEAD)"
  echo "branch=$(git branch --show-current)"
  echo "========================================"
  echo

  "$@"
  EXIT=$?

  echo
  echo "========================================"
  echo "finished=$(date -Is)"
  echo "exit=$EXIT"
  echo "head_after=$(git rev-parse HEAD)"
  echo "========================================"

  exit "$EXIT"
} 2>&1 | tee "$LOG"

EXIT=${PIPESTATUS[0]}
echo "log=$LOG"
exit "$EXIT"
EOF

chmod +x ~/bdspro-ops/run.sh
```

Why this version:
- it does not use global `set -e`, so a failed tested command can still have its exit code and ending metadata recorded;
- `pipefail` + `PIPESTATUS[0]` preserve the wrapped command/group failure even though output is piped through `tee`;
- the command log survives chat/browser interruption.

Check syntax:

```bash
bash -n ~/bdspro-ops/run.sh
echo $?
```

Expected: `0`.

## Phase 5 — Create reader worktree

Only after coordinator confirmation of the authority SHA:

```bash
cd ~/projects/bdspro-canonical-bootstrap-20260908

git worktree add --detach   ~/bdspro-worktrees/reader   fca4682587d93cbc436f2466244ee8bd03b8b1b9
```

Verify:

```bash
cd ~/bdspro-worktrees/reader

git status --short --branch
git rev-parse HEAD
```

Expected:
- detached HEAD;
- clean status;
- exact authority SHA.

Reader is read-only.

## Phase 6 — Create writer worktree

The writer is created only when the coordinator has named the next bounded slice.

Generic form:

```bash
cd ~/projects/bdspro-canonical-bootstrap-20260908

git worktree add   -b local/<slice-name>   ~/bdspro-worktrees/writer   <authority-sha>
```

Do not reuse a stale writer branch for a new slice without coordinator reconciliation.

## Phase 7 — Proof worktree

After writer produces a committed candidate SHA:

```bash
git worktree add --detach   ~/bdspro-worktrees/proof   <candidate-sha>
```

All parallel proof lanes use that same immutable candidate SHA.

## Phase 8 — Running a user-lane command

Example:

```bash
~/bdspro-ops/run.sh   ~/bdspro-worktrees/reader   inventory   python3 shared/code/development/audit-observability-errors.py
```

Interpretation:
- first argument = repo/worktree;
- second = human task name;
- remaining arguments = command to execute;
- stdout/stderr is shown and saved to a timestamped log.

## Phase 9 — What to send back to ChatGPT

For every user-lane task send:
- the full terminal output if reasonably sized;
- exit code;
- path printed as `log=...`;
- any unexpected warning/error;
- never summarize away an error because it looks unrelated.

## Phase 10 — What NOT to do

Do not:
- edit in `reader`;
- push a local writer without explicit authorization;
- run `git reset --hard` to fix confusion;
- run `git clean -fd`;
- force push;
- migrate another service while one writer slice is active;
- add a ratchet before zero proof;
- treat a local log file as source authority.

## Phase 11 — Recovery after interruption

When terminal/browser/chat is interrupted:
- leave tmux process running;
- do not rerun commands automatically;
- reconnect with `tmux attach -t bdspro`;
- inspect the task log;
- start/continue ChatGPT from durable continuation prompt;
- let coordinator reconcile exact Git SHA/CI before new mutation.

## Current first action

Only Phase 1 is authorized for the user's machine right now.

Run Phase 1 exactly and return the output.


## Phase 12 — tmux option scope correction

Operating Model V2 requires stable tmux role names.

Use option scopes correctly:

- `automatic-rename` is a window option.
- `allow-rename` is a pane option.

Canonical pattern:

```bash
tmux set-window-option -t "$SESSION:$name" automatic-rename off

PANE_ID="$(
  tmux display-message     -p     -t "$SESSION:$name"     '#{pane_id}'
)"

tmux set-option   -p   -t "$PANE_ID"   allow-rename off
```

Verification:

```bash
AUTO="$(
  tmux show-window-options     -v     -t "$SESSION:$name"     automatic-rename
)"

PANE_ID="$(
  tmux display-message     -p     -t "$SESSION:$name"     '#{pane_id}'
)"

ALLOW="$(
  tmux show-options     -p     -v     -t "$PANE_ID"     allow-rename
)"

test "$AUTO" = "off"
test "$ALLOW" = "off"
```

Do not treat an empty `show-window-options ... allow-rename` result as a source/tmux failure; it means the verification used the wrong option scope.
