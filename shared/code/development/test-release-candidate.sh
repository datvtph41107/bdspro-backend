#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$tool_dir/../../.." && pwd)"
release_tool="$repo_root/shared/code/development/release-candidate.sh"
protected_deploy="$repo_root/shared/code/deploy.sh"

fail() {
  echo "release-candidate-test: $*" >&2
  exit 1
}

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

repo="$tmp/repo"
mkdir -p "$repo/shared/code/development" "$repo/shared/code" "$repo/payment-service/database/migrations" "$tmp/fake-bin"
cp "$release_tool" "$repo/shared/code/development/release-candidate.sh"
cp "$protected_deploy" "$repo/shared/code/deploy.sh"
chmod +x "$repo/shared/code/development/release-candidate.sh"

cat > "$repo/compose.yaml" <<'EOF'
services: {}
EOF
: > "$repo/.env"
printf 'v1\n' > "$repo/app.txt"
printf 'create table release_probe(id bigint primary key);\n' > "$repo/payment-service/database/migrations/000001_release_probe.up.sql"
printf 'drop table release_probe;\n' > "$repo/payment-service/database/migrations/000001_release_probe.down.sql"

git -C "$repo" init -q
git -C "$repo" config user.name "BDSPro Release Contract"
git -C "$repo" config user.email "release-contract@local.invalid"
git -C "$repo" add .
git -C "$repo" commit -q -m "fixture: previous release"
previous="$(git -C "$repo" rev-parse HEAD)"

printf 'v2\n' > "$repo/app.txt"
git -C "$repo" add app.txt
git -C "$repo" commit -q -m "fixture: current release"
current="$(git -C "$repo" rev-parse HEAD)"

cat > "$tmp/fake-bin/docker" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
log="${QHPRO_FAKE_DOCKER_LOG:?}"
if [[ "${1:-}" == image && "${2:-}" == inspect ]]; then
  image="${@: -1}"
  digest="$(printf '%s' "$image" | sha256sum | awk '{print $1}')"
  printf 'sha256:%s\n' "$digest"
  exit 0
fi
if [[ "${1:-}" == compose ]]; then
  printf 'namespace=%s tag=%s args=%s\n'     "${QHPRO_IMAGE_NAMESPACE:-}" "${QHPRO_IMAGE_TAG:-}" "$*" >> "$log"
  exit 0
fi
echo "unexpected fake docker invocation: $*" >&2
exit 97
EOF
chmod +x "$tmp/fake-bin/docker"
export PATH="$tmp/fake-bin:$PATH"
export QHPRO_FAKE_DOCKER_LOG="$tmp/docker.log"
export QHPRO_RELEASE_SERVICES="auth payment"
export QHPRO_IMAGE_NAMESPACE="fixture-registry/bdspro"
export QHPRO_RELEASE_ENV_FILE="$repo/.env"
export QHPRO_RELEASE_COMPOSE_FILE="$repo/compose.yaml"

previous_dir="$tmp/releases/$previous"
current_dir="$tmp/releases/$current"

(
  cd "$repo"
  bash shared/code/development/release-candidate.sh prepare-commit "$previous" "$previous_dir"
  bash shared/code/development/release-candidate.sh prepare-commit "$current" "$current_dir"
  bash shared/code/development/release-candidate.sh verify "$previous_dir"
  bash shared/code/development/release-candidate.sh verify "$current_dir"

  grep -Fxq "commit=$previous" "$previous_dir/release.meta"
  grep -Fxq "tag=$previous" "$previous_dir/release.meta"
  grep -Fxq "commit=$current" "$current_dir/release.meta"
  grep -Fxq "tag=$current" "$current_dir/release.meta"

  bash shared/code/development/release-candidate.sh capture-images "$previous_dir"
  bash shared/code/development/release-candidate.sh capture-images "$current_dir"
  bash shared/code/development/release-candidate.sh verify-images "$previous_dir"
  bash shared/code/development/release-candidate.sh verify-images "$current_dir"

  : > "$QHPRO_FAKE_DOCKER_LOG"
  bash shared/code/development/release-candidate.sh activate "$current_dir"
  grep -Fq "namespace=fixture-registry/bdspro tag=$current" "$QHPRO_FAKE_DOCKER_LOG"
  grep -Fq -- "--no-build" "$QHPRO_FAKE_DOCKER_LOG"

  : > "$QHPRO_FAKE_DOCKER_LOG"
  bash shared/code/development/release-candidate.sh rollback "$current_dir" "$previous_dir"
  grep -Fq "namespace=fixture-registry/bdspro tag=$previous" "$QHPRO_FAKE_DOCKER_LOG"
  grep -Fq -- "--no-build" "$QHPRO_FAKE_DOCKER_LOG"

  # HEAD-based prepare must reject tracked drift.
  git checkout -q "$current"
  printf 'dirty\n' >> app.txt
  if bash shared/code/development/release-candidate.sh prepare "$tmp/dirty-release" >"$tmp/dirty.out" 2>&1; then
    fail "prepare accepted tracked worktree drift"
  fi
  grep -Fq "tracked worktree changes are not allowed" "$tmp/dirty.out"
  git restore app.txt

  # A later release with a different migration tree must not auto-rollback.
  git checkout -q "$current"
  printf 'alter table release_probe add column note text;\n' > payment-service/database/migrations/000002_release_probe.up.sql
  printf 'alter table release_probe drop column note;\n' > payment-service/database/migrations/000002_release_probe.down.sql
  git add payment-service/database/migrations
  git commit -q -m "fixture: migration-changing release"
  migration_current="$(git rev-parse HEAD)"
  migration_dir="$tmp/releases/$migration_current"
  bash shared/code/development/release-candidate.sh prepare-commit "$migration_current" "$migration_dir"
  bash shared/code/development/release-candidate.sh capture-images "$migration_dir"

  if bash shared/code/development/release-candidate.sh rollback "$migration_dir" "$current_dir" >"$tmp/migration.out" 2>&1; then
    fail "rollback accepted mismatched migration trees"
  fi
  grep -Fq "automatic rollback refused: migration state differs" "$tmp/migration.out"

  # Rollback direction must be ancestral, not arbitrary.
  if bash shared/code/development/release-candidate.sh rollback "$previous_dir" "$current_dir" >"$tmp/ancestry.out" 2>&1; then
    fail "rollback accepted non-ancestor previous release"
  fi
  grep -Fq "previous release is not an ancestor of current release" "$tmp/ancestry.out"
)

echo "release candidate contract PASS"
