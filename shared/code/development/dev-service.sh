#!/usr/bin/env bash
set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
native_stack="$repository_root/shared/code/development/native-stack.sh"

service=${1:?service is required}
shift
[[ $# -gt 0 ]] || { echo 'DEV command is required' >&2; exit 2; }

"$native_stack" prepare-dev "$service"

dev_pid_dir="$repository_root/.tmp/development/dev-pids"
log_dir="$repository_root/.tmp/development/logs"
pid_file="$dev_pid_dir/${service%-service}.pid"
log_file="$log_dir/${service%-service}.log"
mkdir -p "$dev_pid_dir" "$log_dir"

if [[ -e "$pid_file" ]] && ! "$native_stack" dev-pid "$service" >/dev/null 2>&1; then
  rm -f "$pid_file"
fi

if ! (set -o noclobber; printf '%s\n' "$$" >"$pid_file") 2>/dev/null; then
  echo "${service%-service} could not claim DEV ownership; another launcher won the race" >&2
  exit 2
fi

child_pid=
cleanup() {
  if [[ -f "$pid_file" ]] && [[ "$(<"$pid_file")" == "$$" ]]; then
    rm -f "$pid_file"
  fi
}
forward_and_exit() {
  local signal=$1 code=$2
  if [[ -n "$child_pid" ]]; then
    kill "-$signal" "$child_pid" 2>/dev/null || true
  fi
  exit "$code"
}
trap cleanup EXIT
trap 'forward_and_exit HUP 129' HUP
trap 'forward_and_exit INT 130' INT
trap 'forward_and_exit TERM 143' TERM

"$@" > >(tee -a "$log_file") 2> >(tee -a "$log_file" >&2) &
child_pid=$!

set +e
wait "$child_pid"
status=$?
set -e
exit "$status"
