#!/usr/bin/env bash
set -euo pipefail

gateway_url="${QHPRO_GATEWAY_URL:-http://127.0.0.1:8000}"
file_url="${QHPRO_FILE_URL:-http://127.0.0.1:8002}"

probe() {
  local name="$1" url="$2"
  curl --fail --silent --show-error --max-time 10 --output /dev/null "$url"
  printf '[PASS] %-18s %s\n' "$name" "$url"
}

probe "gateway live" "$gateway_url/livez"
probe "gateway ready" "$gateway_url/readyz"
probe "file live" "$file_url/livez"

catalog="$(curl --fail --silent --show-error --max-time 10 \
  "$gateway_url/v2/user/commercial/plans?productCode=qhpro&subjectKind=profile")"
if ! printf '%s' "$catalog" | grep -q 'plans'; then
  echo '[FAIL] public catalog response does not contain plans' >&2
  exit 1
fi
printf '[PASS] %-18s %s\n' "public catalog" "$gateway_url/v2/user/commercial/plans"
