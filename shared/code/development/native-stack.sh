#!/usr/bin/env bash
set -euo pipefail

repository_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
state_dir="$repository_root/.tmp/development"
pid_dir="$state_dir/pids"
log_dir="$state_dir/logs"

# This is process topology, not business ownership. Each entry points to the
# binary and role owned by the corresponding service Makefile.
services=(assistant notification organization payment file user auth hub tqd gateway)

service_dir() {
  printf '%s-service' "$1"
}

service_binary() {
  case "$1" in
    assistant) printf 'assistant-service' ;;
    notification) printf 'notification-service' ;;
    organization) printf 'organization-service' ;;
    payment) printf 'payment-service' ;;
    file) printf 'file-service' ;;
    user) printf 'user-service' ;;
    auth) printf 'auth-service' ;;
    hub) printf 'hub-service' ;;
    tqd) printf 'tqd-service' ;;
    gateway) printf 'gateway-service' ;;
    *) return 1 ;;
  esac
}

service_args() {
  case "$1" in
    notification|organization|payment|hub) printf 'grpc' ;;
    tqd) printf '%s' '-server=grpc' ;;
    gateway) printf 'http' ;;
    assistant|file|user|auth) printf '' ;;
    *) return 1 ;;
  esac
}

service_port() {
  case "$1" in
    gateway) printf '8000' ;;
    file) printf '8002' ;;
    user) printf '8201' ;;
    notification) printf '8204' ;;
    payment) printf '8205' ;;
    organization) printf '8207' ;;
    auth) printf '8216' ;;
    assistant) printf '8218' ;;
    tqd) printf '8219' ;;
    hub) printf '8280' ;;
    *) return 1 ;;
  esac
}

normalize_service() {
  local service=${1%-service}
  for known in "${services[@]}"; do
    if [[ "$known" == "$service" ]]; then
      printf '%s' "$service"
      return 0
    fi
  done
  echo "unknown native service: $1" >&2
  return 2
}

pid_file() {
  printf '%s/%s.pid' "$pid_dir" "$1"
}

log_file() {
  printf '%s/%s.log' "$log_dir" "$1"
}

owned_pid() {
  local file
  file=$(pid_file "$1")
  [[ -s "$file" ]] || return 1
  local pid
  pid=$(<"$file")
  [[ "$pid" =~ ^[0-9]+$ ]] || return 1
  kill -0 "$pid" 2>/dev/null || return 1
  local executable expected
  executable=$(readlink -f "/proc/$pid/exe" 2>/dev/null || true)
  # `make native-restart` builds before stopping so a running executable can
  # legitimately appear as ".../bin/<service> (deleted)" after the atomic
  # replacement. It is still the process recorded in our PID file and must be
  # stopped before the new binary can bind its port.
  executable=${executable% (deleted)}
  expected="$repository_root/$(service_dir "$1")/bin/$(service_binary "$1")"
  [[ "$executable" == "$expected" ]] || return 1
  printf '%s' "$pid"
}

port_open() {
  (exec 3<>"/dev/tcp/127.0.0.1/$1") >/dev/null 2>&1
}

wait_for_port() {
  local service=$1 port=$2 pid=$3
  for _ in $(seq 1 120); do
    if port_open "$port"; then
      return 0
    fi
    if ! kill -0 "$pid" 2>/dev/null; then
      echo "$service exited before port $port became ready" >&2
      tail -n 80 "$(log_file "$service")" >&2 || true
      return 1
    fi
    sleep 0.25
  done
  echo "$service did not open port $port within 30 seconds" >&2
  tail -n 80 "$(log_file "$service")" >&2 || true
  return 1
}

start_one() {
  local service
  service=$(normalize_service "$1")
  mkdir -p "$pid_dir" "$log_dir"
  local existing
  if existing=$(owned_pid "$service"); then
    printf '[READY] %-14s pid=%s port=%s\n' "$service" "$existing" "$(service_port "$service")"
    return 0
  fi

  rm -f "$(pid_file "$service")"
  local port
  port=$(service_port "$service")
  if port_open "$port"; then
    echo "$service cannot start: port $port is owned by another process/container" >&2
    echo "Run 'make integration-down' if an old Compose application stack is still active." >&2
    return 1
  fi

  local directory binary args log
  directory=$(service_dir "$service")
  binary=$(service_binary "$service")
  args=$(service_args "$service")
  log=$(log_file "$service")
  {
    printf '\n[%s] starting %s in host mode\n' "$(date -Iseconds)" "$service"
  } >>"$log"

  (
    cd "$repository_root/$directory"
    set -a
    # Root owns environment/trust; the service file owns localhost resources.
    source "$repository_root/.env"
    source "$repository_root/$directory/.env"
    set +a
    export QHPRO_EXECUTION_MODE=host
    # Local structured evidence belongs to the repository, never to the service CWD.
    export QHPRO_LOG_ROOT="$log_dir"
    if [[ -n "$args" ]]; then
      exec "./bin/$binary" $args
    else
      exec "./bin/$binary"
    fi
  ) >>"$log" 2>&1 &

  local pid=$!
  printf '%s\n' "$pid" >"$(pid_file "$service")"
  if ! wait_for_port "$service" "$port" "$pid"; then
    rm -f "$(pid_file "$service")"
    return 1
  fi
  printf '[START] %-14s pid=%s port=%s log=%s\n' "$service" "$pid" "$port" "$log"
}

stop_one() {
  local service
  service=$(normalize_service "$1")
  local pid
  if ! pid=$(owned_pid "$service"); then
    rm -f "$(pid_file "$service")"
    printf '[DOWN]  %-14s\n' "$service"
    return 0
  fi

  kill -TERM "$pid" 2>/dev/null || true
  for _ in $(seq 1 40); do
    if ! kill -0 "$pid" 2>/dev/null; then
      rm -f "$(pid_file "$service")"
      printf '[STOP]  %-14s pid=%s\n' "$service" "$pid"
      return 0
    fi
    sleep 0.25
  done
  kill -KILL "$pid" 2>/dev/null || true
  rm -f "$(pid_file "$service")"
  printf '[KILL]  %-14s pid=%s\n' "$service" "$pid"
}

start_all() {
  local started=()
  for service in "${services[@]}"; do
    if start_one "$service"; then
      started+=("$service")
      continue
    fi
    for ((index=${#started[@]}-1; index>=0; index--)); do
      stop_one "${started[$index]}" || true
    done
    return 1
  done
}

stop_all() {
  local index
  for ((index=${#services[@]}-1; index>=0; index--)); do
    stop_one "${services[$index]}"
  done
}

status_all() {
  local status=0
  for service in "${services[@]}"; do
    local pid port
    port=$(service_port "$service")
    if pid=$(owned_pid "$service") && port_open "$port"; then
      printf '[READY] %-14s pid=%s port=%s\n' "$service" "$pid" "$port"
    else
      printf '[DOWN]  %-14s port=%s\n' "$service" "$port"
      status=1
    fi
  done
  return "$status"
}

show_logs() {
  local service=${1:-}
  local tail_lines=${TAIL:-200}
  local follow=${FOLLOW:-1}
  local files=()
  mkdir -p "$log_dir"
  if [[ -n "$service" ]]; then
    service=$(normalize_service "$service")
    files+=("$(log_file "$service")")
  else
    for service in "${services[@]}"; do
      files+=("$(log_file "$service")")
    done
  fi
  touch "${files[@]}"
  if [[ "$follow" == 0 ]]; then
    tail -n "$tail_lines" "${files[@]}"
  else
    tail -n "$tail_lines" -F "${files[@]}"
  fi
}

command=${1:-}
case "$command" in
  start) start_all ;;
  stop) stop_all ;;
  restart)
    service=$(normalize_service "${2:?service is required}")
    stop_one "$service"
    start_one "$service"
    ;;
  start-one) start_one "${2:?service is required}" ;;
  stop-one) stop_one "${2:?service is required}" ;;
  status) status_all ;;
  pid) owned_pid "$(normalize_service "${2:?service is required}")" ;;
  logs) show_logs "${2:-}" ;;
  *)
    echo "usage: $0 {start|stop|restart|start-one|stop-one|status|pid|logs} [service]" >&2
    exit 2
    ;;
esac
