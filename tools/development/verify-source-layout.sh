#!/usr/bin/env bash
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

required=(
  README.md CONTRIBUTING.md BACKEND-ARCHITECTURE.md BACKEND-ACCEPTANCE.md
  Makefile compose.yaml go.work bdspro.code-workspace
  docs/reference/DEVELOPER-OPERATING-GUIDE.md
  docs/reference/DEVELOPMENT-ARCHITECTURE-STANDARD.md
  docs/reference/SERVICE-MAP.md docs/reference/MIGRATIONS.md
  tools/development/service.mk tools/development/migration.mk
  tools/development/native-stack.sh tools/development/dev-service.sh
  tools/development/query-structured-logs.py
  tools/development/test_query_structured_logs.py
  tools/development/test-structured-log-query.sh
  tools/development/test-native-stack-ownership.sh
  tools/development/test-dev-stream-mirror.sh
  .cursor/rules/bdspro-development.mdc
)
for path in "${required[@]}"; do
  [[ -e "$path" ]] || { echo "missing canonical source file: $path" >&2; exit 1; }
done

retired_root=(
  qhpro.code-workspace
  generate_bdspro_proto.sh generate_hub_proto.sh generate_tqd_proto.sh generate_utility_proto.sh
  kill-all-services.sh export-selected-to-md.js
  test_api_analyze.py test_api_from_json.py test_compare_apis.py test_data_template.json
  TQD_HTTP_HANDLERS_ANALYSIS.md TQD_PARCEL_OVERVIEW_API_OPTIMIZATION.md TQD_USECASE_ANALYSIS.md
  promt-guildeline.md system-evaluation.md API_USER_GUIDE.md
)
for path in "${retired_root[@]}"; do
  [[ ! -e "$path" ]] || { echo "legacy/ad-hoc file returned to repository root: $path" >&2; exit 1; }
done

legacy_cursor_rules=(
  .cursor/rules/api.mdc .cursor/rules/common.mdc .cursor/rules/danhmuc.mdc
  .cursor/rules/domain.mdc .cursor/rules/transaction.mdc
)
for path in "${legacy_cursor_rules[@]}"; do
  [[ ! -e "$path" ]] || { echo "legacy Cursor rule returned to active rules: $path" >&2; exit 1; }
done

while IFS= read -r -d '' mod; do
  dir="${mod%/go.mod}"
  dir="${dir#./}"
  [[ "$dir" == infrastructure/* || "$dir" == proto || "$dir" == tools/development || "$dir" == integration-test ]] && continue
  [[ -f "$dir/README.md" ]] || { echo "$dir: missing service README.md" >&2; exit 1; }
  [[ -f "$dir/Makefile" ]] || { echo "$dir: missing service Makefile" >&2; exit 1; }
  grep -Eq 'tools/development/(service|go)\.mk' "$dir/Makefile" || {
    echo "$dir: Makefile does not use repository Go/dev contract" >&2; exit 1;
  }
done < <(find api internal infrastructure proto tools/development integration-test -name go.mod -print0 2>/dev/null | sort -z)

[[ -f internal/ai/README.md ]] || { echo 'ai-service: missing README.md' >&2; exit 1; }

# Development launchers must agree on one repository-owned structured-log root.
# This keeps `make up` and service-local/root-routed `make dev` independent of CWD.
expected_log_root="$repo_root/.tmp/development/logs"
dev_recipe="$(make --no-print-directory -C internal/authentication -n dev)"
grep -Fq "QHPRO_LOG_ROOT=\"$expected_log_root\"" <<<"$dev_recipe" || {
  echo "service dev does not inject repository-owned QHPRO_LOG_ROOT=$expected_log_root" >&2
  exit 1
}
grep -Fq 'export QHPRO_LOG_ROOT="$log_dir"' tools/development/native-stack.sh || {
  echo 'native supervisor does not export the repository-owned structured log root' >&2
  exit 1
}
bash -n tools/development/native-stack.sh
bash -n tools/development/dev-service.sh
bash -n tools/development/test-native-stack-ownership.sh
bash -n tools/development/test-dev-stream-mirror.sh
bash -n tools/development/test-structured-log-query.sh
python3 tools/development/verify-source-topology.py
python3 tools/development/verify-docker-topology.py

echo 'source layout verification PASS'
