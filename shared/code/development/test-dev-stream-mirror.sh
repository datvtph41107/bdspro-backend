#!/usr/bin/env bash
set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
native_stack="$repository_root/shared/code/development/native-stack.sh"
dev_runner="$repository_root/shared/code/development/dev-service.sh"

fail() {
  echo "DEV raw-stream mirror contract FAIL: $*" >&2
  exit 1
}

service=auth
marker="r3-$PPID-$$"
capture=$(mktemp)
logs_view=$(mktemp)
pid_file="$repository_root/.tmp/development/dev-pids/$service.pid"

cleanup() {
  rm -f "$capture" "$logs_view"
  if [[ -s "$pid_file" ]]; then
    pid=$(<"$pid_file")
    [[ "$pid" =~ ^[0-9]+$ ]] && kill -TERM "$pid" 2>/dev/null || true
  fi
  rm -f "$pid_file"
}
trap cleanup EXIT

set +e
(
  cd "$repository_root/auth-service"
  "$dev_runner" "$service" bash -c 'printf "%s-out\n" "$1"; printf "%s-err\n" "$1" >&2; exit 7' _ "$marker"
) >"$capture" 2>&1
status=$?
set -e

[[ "$status" == 7 ]] || fail "child exit status changed: expected 7, got $status"
grep -Fq "$marker-out" "$capture" || fail "DEV stdout disappeared from foreground/caller view"
grep -Fq "$marker-err" "$capture" || fail "DEV stderr disappeared from foreground/caller view"

TAIL=400 FOLLOW=0 "$native_stack" logs "$service" >"$logs_view"
grep -Fq "$marker-out" "$logs_view" || fail "DEV stdout missing from repository raw log view"
grep -Fq "$marker-err" "$logs_view" || fail "DEV stderr missing from repository raw log view"

trap - EXIT
cleanup
echo 'DEV raw-stream mirror contract PASS'
