#!/usr/bin/env bash
# Start local BDS Pro: microservice (air) + Admin FE.
# Redis + Postgres dùng remote (config service) — không Docker.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# Hỗ trợ 2 vị trí:
#   bdspro/scripts/run-support-gis.sh
#   bdspro/bdspro-golang-microservices/shared/scripts/run-support-gis.sh
if [[ -d "$SCRIPT_DIR/../../tqd-service" ]]; then
  BE="$(cd "$SCRIPT_DIR/../.." && pwd)"
  ROOT="$(cd "$BE/.." && pwd)"
elif [[ -d "$SCRIPT_DIR/../bdspro-golang-microservices" ]]; then
  ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
  BE="$ROOT/bdspro-golang-microservices"
else
  echo "Cannot resolve BE root from $SCRIPT_DIR"
  exit 1
fi
FE="$ROOT/react-js-admin-bds"
CODE="$BE/shared/code"
LOG_DIR="${TMPDIR:-/tmp}/bdspro-support-gis"
mkdir -p "$LOG_DIR"

CRM_DSN="${CRM_DSN:-}"
TQD_DSN="${TQD_DSN:-}"
AUTH_DSN="${AUTH_DSN:-}"

# Quan trọng: để TRỐNG để GetPrefixProtobufUrl → localhost
# (ENV_RUNTIME=local sẽ dial hostname "user"/"auth" — chỉ dùng trong docker-compose)
unset ENV_RUNTIME || true
export PATH="$PATH:$HOME/go/bin"
if [[ -d "$HOME/.nvm/versions/node" ]]; then
  NODE_DIR="$(ls -1 "$HOME/.nvm/versions/node" | tail -1)"
  export PATH="$PATH:$HOME/.nvm/versions/node/$NODE_DIR/bin"
fi

# name | make_target | port
# Port theo gateway-service/config/local.yml + HTTP phụ (file/user/relay); TQD application APIs đi qua Gateway
ALL_SVCS=(
  "auth|auth-grpc|8216"
  "user|user-grpc|8201"
  "gateway|gateway-http|8000"
  "crm|crm|8206"
  "bdspro|bdspro-grpc|8202"
  "organization|organization-grpc|8207"
  "notification|notification-grpc|8204"
  "social|social-grpc|8209"
  "appointment|appointment-grpc|8213"
  "chat|chat-grpc|8208"
  "chat-http|chat-http|8081"
  "relay|relay-grpc|8308"
  "payment|payment-grpc|8205"
  "membership|membership-grpc|8212"
  "file|file-http|8002"
  "transaction|transaction-grpc|8215"
  "assistant|assistant-grpc|8218"
  "tqd-grpc|tqd-grpc|8219"
)

# Stack tối thiểu: login + ticket/report + upload ảnh + Pro AI
SUPPORT_SVCS=(
  "auth|auth-grpc|8216"
  "user|user-grpc|8201"
  "gateway|gateway-http|8000"
  "crm|crm|8206"
  "file|file-http|8002"
  "assistant|assistant-grpc|8218"
  "tqd-grpc|tqd-grpc|8219"
)

ALL_PORTS=(
  8000 8001 8002 8081 5173 8308 8800
  8201 8202 8203 8204 8205 8206 8207 8208 8209
  8210 8211 8212 8213 8215 8216 8217 8218 8219 8280
)

ALL_MOD_SVCS=(
  auth-service user-service gateway-service crm-service tqd-service
  bdspro-service organization-service notification-service social-service
  appointment-service chat-service relay-service payment-service
  membership-service file-service transaction-service
  assistant-service
)

usage() {
  cat <<'EOF'
Usage: ./scripts/run-support-gis.sh [start|all|support|be|fe|deps|migrate|stop|help]

  start | all   Stop sạch rồi start TẤT CẢ BE + FE   (mặc định)
  support       Stack tối thiểu Support/GIS (auth+user+gateway+crm+tqd) + FE
  be            Chỉ BE full (không FE)
  fe            Chỉ Admin FE (vite)
  deps          Cài migrate/air/protoc + go mod + yarn (không Docker)
  migrate       Apply SQL feature (tickets/reports + auth HT/BC_GIS) lên DB remote
  stop          Kill toàn bộ port/service local (không đụng Redis/DB remote)
EOF
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing command: $1"
    exit 1
  }
}

kill_port() {
  local port="$1"
  if command -v lsof >/dev/null 2>&1; then
    local pids
    pids="$(lsof -ti "tcp:$port" 2>/dev/null || true)"
    if [[ -n "$pids" ]]; then
      # shellcheck disable=SC2086
      kill -9 $pids 2>/dev/null || true
    fi
  fi
  if command -v fuser >/dev/null 2>&1; then
    fuser -k "${port}/tcp" >/dev/null 2>&1 || true
  fi
}

port_listening() {
  local port="$1"
  if command -v lsof >/dev/null 2>&1; then
    lsof -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1 && return 0
  fi
  if command -v ss >/dev/null 2>&1; then
    ss -ltn 2>/dev/null | grep -q ":${port} " && return 0
  fi
  return 1
}

wait_port() {
  local port="$1"
  local name="$2"
  local tries="${3:-40}"
  local i
  for i in $(seq 1 "$tries"); do
    if port_listening "$port"; then
      echo "    ✓ $name listening :$port"
      return 0
    fi
    sleep 0.5
  done
  echo "    ✗ WARN: $name chưa listen :$port — xem log: $LOG_DIR/$name.log"
  tail -n 20 "$LOG_DIR/$name.log" 2>/dev/null || true
  return 1
}

start_svc() {
  local name="$1"
  local make_target="$2"
  local port="$3"
  kill_port "$port"
  echo "==> Start $name (make $make_target) → :$port  log: $LOG_DIR/$name.log"
  (cd "$CODE" && make "$make_target") >"$LOG_DIR/$name.log" 2>&1 &
  echo $! >"$LOG_DIR/$name.pid"
  # Stagger nhẹ để nhiều air không tranh bind cùng lúc
  sleep 0.3
}

start_svc_list() {
  local row name target port
  for row in "$@"; do
    IFS='|' read -r name target port <<<"$row"
    start_svc "$name" "$target" "$port"
  done
}

kill_matching_processes() {
  local pattern="$1"
  local pids
  pids="$(pgrep -f "$pattern" 2>/dev/null || true)"
  if [[ -n "$pids" ]]; then
    # shellcheck disable=SC2086
    kill -9 $pids 2>/dev/null || true
  fi
}

ensure_migrate_cli() {
  if command -v migrate >/dev/null 2>&1; then
    return 0
  fi
  need_cmd go
  echo "==> Installing golang-migrate..."
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
}

ensure_air() {
  if command -v air >/dev/null 2>&1; then
    return 0
  fi
  need_cmd go
  echo "==> Installing air..."
  go install github.com/air-verse/air@latest
}

ensure_protoc_plugins() {
  need_cmd go
  local need=0
  command -v protoc-gen-go >/dev/null 2>&1 || need=1
  command -v protoc-gen-go-grpc >/dev/null 2>&1 || need=1
  command -v protoc-gen-grpc-gateway >/dev/null 2>&1 || need=1
  if [[ "$need" -eq 1 ]]; then
    echo "==> Installing protoc plugins..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  fi
}

do_deps() {
  need_cmd go
  need_cmd make
  need_cmd yarn

  ensure_migrate_cli
  ensure_air
  ensure_protoc_plugins

  echo "==> go mod download..."
  for svc in "${ALL_MOD_SVCS[@]}"; do
    if [[ -d "$BE/$svc" ]]; then
      echo "    $svc"
      (cd "$BE/$svc" && go mod download)
    fi
  done

  need_cmd buf
  echo "==> generate Proto bằng Makefile chuẩn..."
  (
    cd "$CODE"
    make buf-crm
    make buf-auth
    make buf-user
    make buf-tqd
    make buf-hub
  )

  echo "==> yarn install (admin FE)..."
  (cd "$FE" && yarn install)

  echo "Deps done (Redis/DB remote — không Docker)."
}

apply_sql_file() {
  local dsn="$1"
  local sql_file="$2"
  local workdir="$3"
  need_cmd go
  echo "    apply: $(basename "$sql_file")"
  local tmp="$LOG_DIR/apply_sql.go"
  cat >"$tmp" <<'GO'
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Args[1]
	path := os.Args[2]
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if _, err := db.Exec(string(b)); err != nil {
		panic(err)
	}
	fmt.Println("ok")
}
GO
  (cd "$workdir" && go run "$tmp" "$dsn" "$sql_file")
}

do_migrate() {
  : "${CRM_DSN:?Set CRM_DSN before running migrate}"
  : "${TQD_DSN:?Set TQD_DSN before running migrate}"
  : "${AUTH_DSN:?Set AUTH_DSN before running migrate}"

  echo "==> CRM: apply support_tickets SQL..."
  apply_sql_file "$CRM_DSN" \
    "$BE/crm-service/infra/db/migrate_v2/000017_support_tickets.up.sql" \
    "$BE/crm-service"
  apply_sql_file "$CRM_DSN" \
    "$BE/crm-service/infra/db/migrate_v2/000018_support_ticket_images.up.sql" \
    "$BE/crm-service"
  apply_sql_file "$CRM_DSN" \
    "$BE/crm-service/infra/db/migrate_v2/000020_support_ticket_handling_team.up.sql" \
    "$BE/crm-service"
  apply_sql_file "$CRM_DSN" \
    "$BE/crm-service/infra/db/migrate_v2/000021_admin_opportunity_events.up.sql" \
    "$BE/crm-service"
  apply_sql_file "$CRM_DSN" \
    "$BE/crm-service/infra/db/migrate_v2/000022_admin_opportunity_fields.up.sql" \
    "$BE/crm-service"

  echo "==> TQD: apply report QA + images SQL..."
  apply_sql_file "$TQD_DSN" \
    "$BE/tqd-service/infra/db/migrations/v2/000021_report_qa_admin_fields.up.sql" \
    "$BE/tqd-service"
  apply_sql_file "$TQD_DSN" \
    "$BE/tqd-service/infra/db/migrations/v2/000022_report_images.up.sql" \
    "$BE/tqd-service"

  echo "==> Auth: seed ADMIN_HT_* / ADMIN_BC_GIS_* + attach SUPER_ADMIN..."
  # workdir = crm-service vì auth-service go.mod không có github.com/lib/pq
  apply_sql_file "$AUTH_DSN" \
    "$BE/auth-service/scripts/attach_admin_ht_bc_gis_permissions.sql" \
    "$BE/crm-service" || echo "WARN: auth permission seed failed (check AUTH_DSN / HE_THONG exists)"
  apply_sql_file "$AUTH_DSN" \
    "$BE/auth-service/scripts/attach_admin_kd_permissions.sql" \
    "$BE/crm-service" || echo "WARN: auth ADMIN_KD permission seed failed"

  echo "Migrate (feature SQL) done."
  echo "NOTE: nếu menu/API vẫn 403 — xóa Redis ROLE_PERMS_* rồi login lại."
}

print_be_urls() {
  echo
  echo "BE started (ENV_RUNTIME unset → dial localhost)."
  echo "  Redis/DB: remote (theo config service — không Docker)"
  echo "  Gateway:  http://localhost:8000"
  echo "  Auth:     :8216"
  echo "  User:     :8201"
  echo "  CRM:      :8206"
  echo "  File:     http://localhost:8002  (upload /v1/file/upload)"
  echo "  TQD API:  http://localhost:8000/v2/tqd/... (Gateway)"
  echo "  TQD gRPC: :8219                 (Pro AI / planning)"
  echo "  Logs:     $LOG_DIR/*.log"
  echo "  Login API: POST http://localhost:8000/v2/auth/admin/login"
}

wait_critical_ports() {
  echo "==> Chờ critical ports..."
  # Auth: Listen() trước AutoMigrate+Serve → port mở sớm nhưng chưa sẵn sàng gRPC.
  # Đợi log "gRPC server listening" (tối đa ~90s) thay vì chỉ check TCP.
  local i
  for i in $(seq 1 90); do
    if grep -q 'gRPC server listening' "$LOG_DIR/auth.log" 2>/dev/null; then
      echo "    ✓ auth ready (:8216)"
      break
    fi
    if [[ "$i" -eq 90 ]]; then
      echo "    ✗ WARN: auth chưa ready — xem log: $LOG_DIR/auth.log"
      tail -n 30 "$LOG_DIR/auth.log" 2>/dev/null || true
    fi
    sleep 1
  done
  wait_port 8201 user 40 || true
  wait_port 8000 gateway 40 || true
}

do_be_full() {
  need_cmd air
  need_cmd make
  # Dọn instance cũ trước khi start full (tránh air chồng / port in use)
  do_stop_soft
  start_svc_list "${ALL_SVCS[@]}"
  wait_critical_ports
  print_be_urls
}

do_be_support() {
  need_cmd air
  need_cmd make
  do_stop_soft
  start_svc_list "${SUPPORT_SVCS[@]}"
  wait_critical_ports
  print_be_urls
}

do_fe() {
  need_cmd yarn
  kill_port 5173
  # Dọn vite/yarn sót (không dùng pkill -f 'yarn' vì dễ match nhầm)
  kill_matching_processes "$FE/node_modules/.bin/vite"
  echo "==> Start Admin FE → log: $LOG_DIR/fe.log"
  # Polling tránh EMFILE khi inotify instances gần đầy (nhiều air)
  (
    cd "$FE"
    export CHOKIDAR_USEPOLLING=1
    export WATCHPACK_POLLING=true
    ulimit -n 65536 2>/dev/null || true
    yarn dev
  ) >"$LOG_DIR/fe.log" 2>&1 &
  echo $! >"$LOG_DIR/fe.pid"
  # Đợi vite bind :5173 (tối đa ~20s)
  local i
  for i in $(seq 1 40); do
    if port_listening 5173; then
      echo "FE started. URL: http://localhost:5173/login"
      return 0
    fi
    sleep 0.5
  done
  echo "WARN: FE chưa listen :5173 — xem log: $LOG_DIR/fe.log"
  tail -n 30 "$LOG_DIR/fe.log" || true
}

do_stop_soft() {
  echo "==> Stopping existing services (ports + pids)..."
  local port
  for port in "${ALL_PORTS[@]}"; do
    kill_port "$port"
  done
  local f
  for f in "$LOG_DIR"/*.pid; do
    [[ -f "$f" ]] || continue
    kill -9 "$(cat "$f")" 2>/dev/null || true
    rm -f "$f"
  done
  # Dọn air + tmp/main (tránh chồng watcher → FE EMFILE)
  kill_matching_processes '(^|/)air($| )'
  kill_matching_processes "$BE/.*/tmp/main"
  kill_matching_processes "$FE/node_modules/.bin/vite"
  sleep 1
}

do_stop() {
  do_stop_soft
  echo "Stopped (không đụng Redis/DB remote)."
}

print_ui() {
  echo
  echo "UI:"
  echo "  Login:    http://localhost:5173/login"
  echo "  Tickets:  http://localhost:5173/admin/support/tickets"
  echo "  Reports:  http://localhost:5173/admin/gis-reports"
}

cmd="${1:-start}"
case "$cmd" in
  -h|--help|help) usage ;;
  deps) do_deps ;;
  migrate) do_migrate ;;
  be) do_be_full ;;
  support)
    if command -v buf >/dev/null 2>&1 && [[ -d "$BE/shared/protobuf/schema" ]]; then
      echo "==> buf generate crm + tqd proto (gateway phụ thuộc cả hai)..."
      (
        cd "$BE/shared/protobuf/schema"
        buf generate --template ./buf.gen.yaml --path crm || true
        buf generate --template ./buf.gen.yaml --path tqd || true
      )
    else
      echo "warn: buf chưa sẵn — bỏ qua generate proto"
    fi
    do_be_support
    do_fe
    print_ui
    ;;
  fe) do_fe ;;
  stop) do_stop ;;
  start|start-all|all)
    do_deps
    do_migrate
    do_be_full
    do_fe
    print_ui
    ;;
  *)
    usage
    exit 1
    ;;
esac
