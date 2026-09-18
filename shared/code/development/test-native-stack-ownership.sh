#!/usr/bin/env bash
set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
native_stack="$repository_root/shared/code/development/native-stack.sh"
dev_runner="$repository_root/shared/code/development/dev-service.sh"

source "$native_stack"

fail() {
  echo "native ownership contract FAIL: $*" >&2
  exit 1
}

assert_state() {
  local expected=$1 supervised=$2 dev=$3 port_flag=$4
  (
    owned_pid() { [[ "$supervised" == yes ]] && printf '111'; }
    dev_pid() { [[ "$dev" == yes ]] && printf '222'; }
    service_port() { printf '9999'; }
    port_open() { [[ "$port_flag" == yes ]]; }

    local state pid
    read -r state pid < <(ownership_state auth)
    [[ "$state" == "$expected" ]] || fail "expected $expected, got $state"
  )
}

assert_state SUPERVISED yes no no
assert_state DEV no yes no
assert_state FOREIGN no no yes
assert_state DOWN no no no

(
  owned_pid() { return 1; }
  dev_pid() { printf '222'; }
  service_port() { printf '9999'; }
  port_open() { return 1; }
  output=$(status_one auth 2>/dev/null || true)
  [[ "$output" == *"[DEV"* && "$output" == *"ready=no"* ]] || fail "DEV readiness separation missing: $output"
)

(
  mode=supervised
  owned_pid() { [[ "$mode" == supervised ]] && printf '111'; }
  dev_pid() { [[ "$mode" == dev ]] && printf '222'; }
  service_port() { printf '9999'; }
  port_open() { [[ "$mode" == foreign ]]; }
  stop_one() { mode=down; }
  prepare_dev auth
  [[ "$mode" == down ]] || fail "SUPERVISED was not transferred to DEV-ready DOWN"
)

(
  owned_pid() { return 1; }
  dev_pid() { printf '222'; }
  service_port() { printf '9999'; }
  port_open() { return 1; }
  if prepare_dev auth >/dev/null 2>&1; then fail "second DEV was not rejected"; fi
)

(
  owned_pid() { return 1; }
  dev_pid() { return 1; }
  service_port() { printf '9999'; }
  port_open() { return 0; }
  if prepare_dev auth >/dev/null 2>&1; then fail "FOREIGN owner was not rejected"; fi
)

pid_file="$repository_root/.tmp/development/dev-pids/auth.pid"
mkdir -p "$(dirname "$pid_file")"
printf '999999\n' >"$pid_file"

cleanup_actual() {
  if [[ -s "$pid_file" ]]; then
    pid=$(<"$pid_file")
    [[ "$pid" =~ ^[0-9]+$ ]] && kill -TERM "$pid" 2>/dev/null || true
  fi
  rm -f "$pid_file"
}
trap cleanup_actual EXIT

(
  cd "$repository_root/auth-service"
  "$dev_runner" auth bash -c 'while :; do sleep 1; done'
) &
launcher_job=$!

for _ in $(seq 1 100); do
  if [[ -s "$pid_file" ]] && "$native_stack" dev-pid auth >/dev/null 2>&1; then
    break
  fi
  sleep 0.05
done

[[ -s "$pid_file" ]] || fail "DEV runner did not create owner state"
dev_owner=$("$native_stack" dev-pid auth) || fail "DEV owner marker did not validate"

if (
  cd "$repository_root/auth-service"
  "$dev_runner" auth true
) >/dev/null 2>&1; then
  fail "concurrent second DEV launch unexpectedly succeeded"
fi

status_output=$("$native_stack" status-one auth 2>/dev/null || true)
[[ "$status_output" == *"[DEV"* ]] || fail "status did not report DEV: $status_output"

kill -TERM "$dev_owner"
set +e
wait "$launcher_job"
runner_status=$?
set -e
[[ "$runner_status" == 143 || "$runner_status" == 0 ]] || fail "unexpected DEV runner TERM status: $runner_status"

for _ in $(seq 1 100); do
  [[ ! -e "$pid_file" ]] && break
  sleep 0.05
done
[[ ! -e "$pid_file" ]] || fail "DEV owner state was not cleaned after TERM"

trap - EXIT
echo 'native ownership contract PASS'
