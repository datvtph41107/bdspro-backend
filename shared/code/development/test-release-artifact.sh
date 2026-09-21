#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$tool_dir/../../.." && pwd)"
release_tool="$repo_root/shared/code/development/release-artifact.sh"
source_manifest_tool="$repo_root/shared/code/development/source-manifest.sh"
protected_deploy="$repo_root/shared/code/deploy.sh"

fail() {
  echo "release-artifact-test: $*" >&2
  exit 1
}

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

repo="$tmp/repo"
previous_worktree="$tmp/previous"
mkdir -p   "$repo/shared/code/development"   "$repo/shared/code"   "$repo/payment-service/database/migrations"   "$repo/organization-service/migrate"   "$tmp/fake-bin"

cp "$release_tool" "$repo/shared/code/development/release-artifact.sh"
cp "$source_manifest_tool" "$repo/shared/code/development/source-manifest.sh"
cp "$protected_deploy" "$repo/shared/code/deploy.sh"
chmod +x   "$repo/shared/code/development/release-artifact.sh"   "$repo/shared/code/development/source-manifest.sh"

cat > "$repo/compose.yaml" <<'EOF'
services:
  auth:
    image: ${QHPRO_IMAGE_NAMESPACE:-bdspro}/auth-service:${QHPRO_IMAGE_TAG:-local}
  payment:
    image: ${QHPRO_IMAGE_NAMESPACE:-bdspro}/payment-service:${QHPRO_IMAGE_TAG:-local}
  redis:
    image: redis:7.4-alpine
EOF

cat > "$repo/.env" <<'EOF'
QHPRO_ENVIRONMENT=development
EOF

printf 'v1\n' > "$repo/app.txt"
printf 'create table release_probe(id bigint primary key);\n'   > "$repo/payment-service/database/migrations/000001_release_probe.up.sql"
printf 'drop table release_probe;\n'   > "$repo/payment-service/database/migrations/000001_release_probe.down.sql"
printf 'create table organization_probe(id bigint primary key);\n'   > "$repo/organization-service/migrate/000001_organization_probe.up.sql"
printf 'drop table organization_probe;\n'   > "$repo/organization-service/migrate/000001_organization_probe.down.sql"

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

git -C "$repo" worktree add -q --detach "$previous_worktree" "$previous"
cp "$repo/.env" "$previous_worktree/.env"

cat > "$tmp/fake-bin/docker" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

log="${QHPRO_FAKE_DOCKER_LOG:?}"
printf 'docker %s\n' "$*" >> "$log"

if [[ "${1:-}" == "image" && "${2:-}" == "inspect" ]]; then
  image="${@: -1}"
  digest="$(printf '%s' "$image" | sha256sum | awk '{print $1}')"
  printf 'sha256:%s\n' "$digest"
  exit 0
fi

if [[ "${1:-}" == "save" ]]; then
  shift
  [[ "${1:-}" == "--output" ]] || exit 91
  output="$2"
  shift 2
  printf '%s\n' "$@" > "$output"
  exit 0
fi

if [[ "${1:-}" == "load" ]]; then
  [[ "${2:-}" == "--input" ]] || exit 92
  [[ -s "${3:-}" ]] || exit 93
  printf 'Loaded image archive\n'
  exit 0
fi

if [[ "${1:-}" == "compose" ]]; then
  args=" $* "
  if [[ "$args" == *" config --images "* ]]; then
    printf '%s\n'       "${QHPRO_IMAGE_NAMESPACE:?}/auth-service:${QHPRO_IMAGE_TAG:?}"       "${QHPRO_IMAGE_NAMESPACE:?}/payment-service:${QHPRO_IMAGE_TAG:?}"       "redis:7.4-alpine"
    exit 0
  fi
  if [[ "$args" == *" up "* ]]; then
    exit 0
  fi
fi

echo "unexpected fake docker invocation: $*" >&2
exit 97
EOF
chmod +x "$tmp/fake-bin/docker"

export PATH="$tmp/fake-bin:$PATH"
export QHPRO_FAKE_DOCKER_LOG="$tmp/docker.log"
export QHPRO_IMAGE_NAMESPACE="fixture-registry/bdspro"

previous_dir="$tmp/releases/$previous"
current_dir="$tmp/releases/$current"

: > "$QHPRO_FAKE_DOCKER_LOG"
QHPRO_RELEASE_SOURCE_ROOT="$previous_worktree"   bash "$release_tool" prepare "$previous_dir"
QHPRO_RELEASE_SOURCE_ROOT="$previous_worktree"   bash "$release_tool" capture-images "$previous_dir"
QHPRO_RELEASE_SOURCE_ROOT="$previous_worktree"   bash "$release_tool" verify "$previous_dir"

(
  cd "$repo"
  bash shared/code/development/release-artifact.sh prepare "$current_dir"
  bash shared/code/development/release-artifact.sh capture-images "$current_dir"
  bash shared/code/development/release-artifact.sh verify "$current_dir"
)

grep -Fxq "commit=$previous" "$previous_dir/release.meta"
grep -Fxq "tag=$previous" "$previous_dir/release.meta"
grep -Fxq "commit=$current" "$current_dir/release.meta"
grep -Fxq "tag=$current" "$current_dir/release.meta"
test "$(wc -l < "$current_dir/migration-trees.tsv" | tr -d ' ')" = 2
grep -Fq $'organization-service/migrate' "$current_dir/migration-trees.tsv"
grep -Fq $'payment-service/database/migrations' "$current_dir/migration-trees.tsv"
test -s "$current_dir/images.tar"
test -s "$current_dir/images.tar.sha256"

# Image archive tampering must fail closed.
cp "$current_dir/images.tar" "$tmp/images.tar.good"
printf 'tamper\n' >> "$current_dir/images.tar"
if (cd "$current_dir" && sha256sum -c images.tar.sha256 >/dev/null 2>&1); then
  fail "direct checksum accepted a tampered image archive"
fi
if (
  cd "$repo"
  bash shared/code/development/release-artifact.sh verify "$current_dir"
) >"$tmp/tamper.out" 2>&1; then
  fail "verify accepted a tampered image archive"
fi
grep -Fq "release evidence checksum mismatch" "$tmp/tamper.out"
cp "$tmp/images.tar.good" "$current_dir/images.tar"
(
  cd "$current_dir"
  sha256sum images.tar > images.tar.sha256
  sha256sum     SOURCE-MANIFEST.sha256     source-artifact.sha256     migration-trees.tsv     migration-trees.sha256     release.meta     source-evidence.sha256     image-manifest.tsv     image-manifest.tsv.sha256     images.tar.sha256     > release-evidence.sha256
)

# Restoring an artifact must load the archive and verify every recorded image ID.
: > "$QHPRO_FAKE_DOCKER_LOG"
(
  cd "$repo"
  bash shared/code/development/release-artifact.sh restore-images "$current_dir"
)
grep -Fq "docker load --input $current_dir/images.tar" "$QHPRO_FAKE_DOCKER_LOG"

# Activation must load the archive and explicitly forbid builds and pulls.
: > "$QHPRO_FAKE_DOCKER_LOG"
(
  cd "$repo"
  bash shared/code/development/release-artifact.sh activate "$current_dir"
)
grep -Fq "docker load --input $current_dir/images.tar" "$QHPRO_FAKE_DOCKER_LOG"
grep -Fq -- "--no-build" "$QHPRO_FAKE_DOCKER_LOG"
grep -Fq -- "--pull never" "$QHPRO_FAKE_DOCKER_LOG"

# Production activation must fail before image loading or Compose activation.
printf 'QHPRO_ENVIRONMENT=production\n' > "$repo/.env"
: > "$QHPRO_FAKE_DOCKER_LOG"
if (
  cd "$repo"
  bash shared/code/development/release-artifact.sh activate "$current_dir"
) >"$tmp/production.out" 2>&1; then
  fail "production activation was accepted"
fi
grep -Fq "production activation is forbidden" "$tmp/production.out"
! grep -Fq "docker load" "$QHPRO_FAKE_DOCKER_LOG"
! grep -Fq " up " "$QHPRO_FAKE_DOCKER_LOG"
printf 'QHPRO_ENVIRONMENT=development\n' > "$repo/.env"

# HEAD-based prepare must reject tracked source drift.
printf 'dirty\n' >> "$repo/app.txt"
if (
  cd "$repo"
  bash shared/code/development/release-artifact.sh prepare "$tmp/dirty-release"
) >"$tmp/dirty.out" 2>&1; then
  fail "prepare accepted tracked worktree drift"
fi
grep -Fq "tracked worktree changes are not allowed" "$tmp/dirty.out"
git -C "$repo" restore app.txt

# Rollback must use the already-built previous archive and never rebuild/pull.
: > "$QHPRO_FAKE_DOCKER_LOG"
(
  cd "$repo"
  bash shared/code/development/release-artifact.sh rollback "$current_dir" "$previous_dir"
)
grep -Fq "docker load --input $previous_dir/images.tar" "$QHPRO_FAKE_DOCKER_LOG"
grep -Fq -- "--no-build" "$QHPRO_FAKE_DOCKER_LOG"
grep -Fq -- "--pull never" "$QHPRO_FAKE_DOCKER_LOG"

# A release with different migration identity cannot auto-rollback.
printf 'alter table release_probe add column note text;\n'   > "$repo/payment-service/database/migrations/000002_release_probe.up.sql"
printf 'alter table release_probe drop column note;\n'   > "$repo/payment-service/database/migrations/000002_release_probe.down.sql"
git -C "$repo" add payment-service/database/migrations
git -C "$repo" commit -q -m "fixture: migration-changing release"
migration_current="$(git -C "$repo" rev-parse HEAD)"
migration_dir="$tmp/releases/$migration_current"
(
  cd "$repo"
  bash shared/code/development/release-artifact.sh prepare "$migration_dir"
  bash shared/code/development/release-artifact.sh capture-images "$migration_dir"
)
if (
  cd "$repo"
  bash shared/code/development/release-artifact.sh rollback "$migration_dir" "$current_dir"
) >"$tmp/migration.out" 2>&1; then
  fail "rollback accepted mismatched migration identity"
fi
grep -Fq "automatic rollback refused: migration state differs" "$tmp/migration.out"

# Rollback direction must be ancestral.
if (
  cd "$repo"
  bash shared/code/development/release-artifact.sh rollback "$previous_dir" "$current_dir"
) >"$tmp/ancestry.out" 2>&1; then
  fail "rollback accepted a non-ancestor previous release"
fi
grep -Fq "previous release is not an ancestor of current release" "$tmp/ancestry.out"

echo "release artifact contract PASS"
