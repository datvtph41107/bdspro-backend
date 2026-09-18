#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
service=r4-contract-service
service_root="$repo_root/.tmp/development/logs/$service"
runtime="$service_root/runs/r4-contract/runtime.jsonl"
channel="$service_root/runs/r4-contract/channels/http.jsonl"
capture=$(mktemp)

cleanup() {
  rm -f "$capture"
  rm -rf "$service_root"
}
trap cleanup EXIT

mkdir -p "$(dirname "$runtime")" "$(dirname "$channel")"
cat >"$runtime" <<'JSONL'
{"time":"2026-09-18T02:00:03Z","level":"ERROR","msg":"late","service.name":"r4-contract-service","request_id":"r4-request","operation_id":"r4-operation","component":"worker"}
{"time":"2026-09-18T02:00:01Z","level":"INFO","msg":"early","service.name":"r4-contract-service","request_id":"r4-request","operation_id":"r4-operation","event_name":"r4.started"}
JSONL
cat >"$channel" <<'JSONL'
{"time":"2026-09-18T02:00:04Z","level":"INFO","msg":"duplicate-projection","service.name":"r4-contract-service","request_id":"r4-request"}
JSONL

make --no-print-directory -C "$repo_root" log-query \
  request_id=r4-request \
  service="$service" \
  LIMIT=1 \
  FORMAT=json >"$capture"

grep -Fq '"msg":"late"' "$capture"
! grep -Fq '"msg":"early"' "$capture"
! grep -Fq 'duplicate-projection' "$capture"

make --no-print-directory -C "$repo_root" log-query \
  operation_id=r4-operation \
  event_name=r4.started \
  service="$service" >"$capture"

grep -Fq 'event_name=r4.started' "$capture"
grep -Fq 'msg=early' "$capture"

echo 'structured runtime log query contract PASS'
