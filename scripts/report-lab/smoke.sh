#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/env.sh"

for command_name in curl jq node; do
  command -v "$command_name" >/dev/null 2>&1 || {
    echo "$command_name is required" >&2
    exit 1
  }
done

token="$(node "$REPORT_LAB_DIR/create-token.mjs")"
suffix="$(date +%s)"
operation_id="op-report-smoke-$suffix"
idempotency_key="idem-report-smoke-$suffix"
temp_dir="$(mktemp -d)"
trap 'rm -rf "$temp_dir"' EXIT

create_response="$temp_dir/create.json"
http_code="$(curl --silent --show-error --output "$create_response" --write-out '%{http_code}' \
  'http://127.0.0.1:8000/v2/tqd/map-workspace/reports' \
  --header "Authorization: Bearer $token" \
  --header "X-Operation-ID: $operation_id" \
  --header "Idempotency-Key: $idempotency_key" \
  --header 'Content-Type: application/json' \
  --data "{\"reportType\":1,\"entityType\":1,\"entityId\":1,\"parcelId\":1,\"title\":\"QHPro local smoke\",\"clientRequestId\":\"$idempotency_key\",\"metadataJson\":\"{\\\"smoke\\\":true}\"}")"

jq . "$create_response"
[[ "$http_code" == "200" ]] || {
  echo "create report returned HTTP $http_code" >&2
  exit 1
}

report_id="$(jq -r '.data.item.reportId // empty' "$create_response")"
[[ -n "$report_id" ]] || {
  echo 'response does not contain data.item.reportId' >&2
  exit 1
}

get_response="$temp_dir/get.json"
pdf_url=''
for _ in $(seq 1 20); do
  curl --silent --show-error \
    "http://127.0.0.1:8000/v2/tqd/map-workspace/reports/$report_id" \
    --header "Authorization: Bearer $token" \
    --header "X-Operation-ID: get-report-$suffix" \
    --output "$get_response"
  status="$(jq -r '.data.status // 0' "$get_response")"
  pdf_url="$(jq -r '.data.assets.pdfUrl // empty' "$get_response")"
  if [[ "$status" == "10" && -n "$pdf_url" ]]; then
    break
  fi
  sleep 1
done

jq . "$get_response"
[[ -n "$pdf_url" ]] || {
  echo 'report did not become ready with a PDF URL' >&2
  exit 1
}

output_file="$REPORT_LAB_DIR/qhpro-smoke-output.pdf"
curl --silent --show-error "$pdf_url" --output "$output_file"
[[ "$(head -c 5 "$output_file")" == '%PDF-' ]] || {
  echo 'downloaded output is not a PDF' >&2
  exit 1
}

echo "PASS report_id=$report_id"
echo "PASS pdf_url=$pdf_url"
echo "PASS output=$output_file"
