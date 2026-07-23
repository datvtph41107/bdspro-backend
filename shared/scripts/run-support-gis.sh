#!/usr/bin/env bash
# Chạy local tối thiểu: login (auth+user+redis+gateway) + ticket/report (crm+tqd) + FE
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BE="$ROOT/bdspro-golang-microservices"
FE="$ROOT/react-js-admin-bds"
CODE="$BE/shared/code"
LOG_DIR="${TMPDIR:-/tmp}/bdspro-support-gis"
mkdir -p "$LOG_DIR"

CRM_DSN="${CRM_DSN:-postgresql://hieubap:Matkhau.1@14.225.210.29:5432/db_crm?sslmode=disable}"
TQD_DSN="${TQD_DSN:-postgresql://hieubap:Matkhau.1@14.225.210.29:5433/db_tqd?sslmode=disable}"
REDIS_NAME="${REDIS_NAME:-bdspro-redis-local}"
REDIS_PASS="${REDIS_PASS:-protected_redis_secure_password}"

# Quan trọng: để TRỐNG để GetPrefixProtobufUrl → localhost
# (ENV_RUNTIME=local sẽ dial hostname "user"/"auth" — chỉ dùng trong docker-compose)
unset ENV_RUNTIME || true
export PATH="$PATH:$HOME/go/bin"
# nvm node nếu có
if [[ -d "$HOME/.nvm/versions/node" ]]; then
  NODE_DIR="$(ls -1 "$HOME/.nvm/versions/node" | tail -1)"
  export PATH="$PATH:$HOME/.nvm/versions/node/$NODE_DIR/bin"
fi

usage() {
  cat <<'EOF'
Usage: ./scripts/run-support-gis.sh [deps|migrate|be|fe|all|stop]

  deps     Cài migrate/air/protoc + go mod + yarn + redis docker
  migrate  Apply SQL feature (CRM support_tickets + TQD report QA)
  be       Redis + auth + user + gateway + crm + tqd (đủ login + feature)
  fe       Admin FE (vite)
  all      deps + migrate + be + fe  (mặc định)
  stop     Kill ports + redis container
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
    lsof -ti "tcp:$port" 2>/dev/null | xargs -r kill -9 2>/dev/null || true
  fi
}

start_svc() {
  local name="$1"
  local make_target="$2"
  local port="$3"
  kill_port "$port"
  echo "==> Start $name (make $make_target) → :$port  log: $LOG_DIR/$name.log"
  (cd "$CODE" && make "$make_target") >"$LOG_DIR/$name.log" 2>&1 &
  echo $! >"$LOG_DIR/$name.pid"
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

ensure_redis() {
  need_cmd docker
  if docker ps --format '{{.Names}}' | grep -qx "$REDIS_NAME"; then
    echo "==> Redis already running ($REDIS_NAME)"
    return 0
  fi
  if docker ps -a --format '{{.Names}}' | grep -qx "$REDIS_NAME"; then
    echo "==> Starting existing Redis container..."
    docker start "$REDIS_NAME" >/dev/null
    return 0
  fi
  echo "==> Creating Redis container ($REDIS_NAME) on :6379..."
  docker run -d --name "$REDIS_NAME" -p 6379:6379 \
    redis:7-alpine redis-server --requirepass "$REDIS_PASS" >/dev/null
}

do_deps() {
  need_cmd go
  need_cmd make
  need_cmd yarn
  need_cmd docker

  ensure_migrate_cli
  ensure_air
  ensure_protoc_plugins
  ensure_redis

  echo "==> go mod download (auth / user / gateway / crm / tqd)..."
  for svc in auth-service user-service gateway-service crm-service tqd-service; do
    echo "    $svc"
    (cd "$BE/$svc" && go mod download)
  done

  if command -v buf >/dev/null 2>&1 && [[ -d "$BE/shared/protobuf/schema" ]]; then
    echo "==> buf generate crm + auth + user..."
    (
      cd "$BE/shared/protobuf/schema"
      buf generate --template ./buf.gen.yaml --path crm || true
      buf generate --template ./buf.gen.yaml --path auth || true
      buf generate --template ./buf.gen.yaml --path user || true
    )
  else
    echo "warn: buf chưa sẵn — bỏ qua generate proto"
  fi

  echo "==> yarn install (admin FE)..."
  (cd "$FE" && yarn install)

  echo "Deps done."
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
  echo "==> CRM: apply support_tickets SQL..."
  apply_sql_file "$CRM_DSN" \
    "$BE/crm-service/infra/db/migrate_v2/000012_support_tickets.up.sql" \
    "$BE/crm-service"

  echo "==> TQD: apply report QA fields SQL..."
  apply_sql_file "$TQD_DSN" \
    "$BE/tqd-service/infra/db/migrations/v2/000021_report_qa_admin_fields.up.sql" \
    "$BE/tqd-service"

  echo "Migrate (feature SQL) done."
}

do_be() {
  need_cmd air
  need_cmd make
  ensure_redis

  # Login stack
  start_svc auth auth-grpc 8216
  start_svc user user-grpc 8201
  start_svc gateway gateway-http 8000

  # Feature stack
  start_svc crm crm 8206
  start_svc tqd tqd-http 8230

  echo
  echo "BE started (ENV_RUNTIME unset → dial localhost)."
  echo "  Redis:    localhost:6379 ($REDIS_NAME)"
  echo "  Gateway:  http://localhost:8000"
  echo "  Auth:     :8216"
  echo "  User:     :8201"
  echo "  CRM:      :8206"
  echo "  TQD HTTP: http://localhost:8230"
  echo "  Logs:     $LOG_DIR/*.log"
  echo "  Login API: POST http://localhost:8000/v2/auth/admin/login"
}

do_fe() {
  need_cmd yarn
  echo "==> Start Admin FE → log: $LOG_DIR/fe.log"
  (cd "$FE" && yarn dev) >"$LOG_DIR/fe.log" 2>&1 &
  echo $! >"$LOG_DIR/fe.pid"
  echo "FE started. URL thường http://localhost:5173/login"
}

do_stop() {
  echo "==> Stopping services..."
  for port in 8000 8201 8206 8216 8230 5173; do
    kill_port "$port"
  done
  for f in auth user gateway crm tqd fe; do
    if [[ -f "$LOG_DIR/$f.pid" ]]; then
      kill -9 "$(cat "$LOG_DIR/$f.pid")" 2>/dev/null || true
      rm -f "$LOG_DIR/$f.pid"
    fi
  done
  if docker ps --format '{{.Names}}' 2>/dev/null | grep -qx "$REDIS_NAME"; then
    docker stop "$REDIS_NAME" >/dev/null || true
  fi
  echo "Stopped."
}

cmd="${1:-all}"
case "$cmd" in
  -h|--help|help) usage ;;
  deps) do_deps ;;
  migrate) do_migrate ;;
  be) do_be ;;
  fe) do_fe ;;
  stop) do_stop ;;
  all)
    do_deps
    do_migrate
    do_be
    do_fe
    echo
    echo "UI:"
    echo "  Login:    http://localhost:5173/login"
    echo "  Tickets:  http://localhost:5173/admin/support/tickets"
    echo "  Reports:  http://localhost:5173/admin/gis-reports"
    ;;
  *)
    usage
    exit 1
    ;;
esac
