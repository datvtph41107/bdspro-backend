#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$tool_dir/../../.." && pwd)"
compose_file="${QHPRO_RELEASE_COMPOSE_FILE:-$repo_root/compose.yaml}"
env_file="${QHPRO_RELEASE_ENV_FILE:-$repo_root/.env}"
protected_deploy_sha="80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e"

die() {
  echo "release-candidate: $*" >&2
  exit 2
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "missing command: $1"
}

read_metadata() {
  local dir="$1"
  local file="$dir/release.meta"
  [[ -f "$file" ]] || die "missing release metadata: $file"

  RELEASE_COMMIT=
  RELEASE_TREE=
  RELEASE_NAMESPACE=
  RELEASE_TAG=
  RELEASE_SERVICES=

  while IFS='=' read -r key value; do
    case "$key" in
      commit) RELEASE_COMMIT="$value" ;;
      tree) RELEASE_TREE="$value" ;;
      namespace) RELEASE_NAMESPACE="$value" ;;
      tag) RELEASE_TAG="$value" ;;
      services) RELEASE_SERVICES="$value" ;;
      '') ;;
      *) die "unknown release metadata key: $key" ;;
    esac
  done < "$file"

  [[ "$RELEASE_COMMIT" =~ ^[0-9a-f]{40}$ ]] || die "invalid release commit"
  [[ "$RELEASE_TREE" =~ ^[0-9a-f]{40}$ ]] || die "invalid release tree"
  [[ -n "$RELEASE_NAMESPACE" ]] || die "release namespace is empty"
  [[ "$RELEASE_TAG" == "$RELEASE_COMMIT" ]] || die "release tag must equal exact commit SHA"
  [[ -n "$RELEASE_SERVICES" ]] || die "release service set is empty"
}

source_manifest_for_commit() {
  local commit="$1"
  local output="$2"
  local source_path digest
  : > "$output"

  while IFS= read -r -d '' source_path; do
    [[ "$source_path" == "SOURCE-MANIFEST.sha256" ]] && continue
    digest="$(git show "$commit:$source_path" | sha256sum | awk '{print $1}')"
    printf '%s  ./%s\n' "$digest" "$source_path" >> "$output"
  done < <(git ls-tree -r -z --name-only "$commit" | sort -z)
}

migration_manifest() {
  local commit="$1"
  local output="$2"
  local migration_dir
  : > "$output"

  while IFS= read -r migration_dir; do
    [[ -n "$migration_dir" ]] || continue
    printf '%s\t%s\n' "$(git rev-parse "$commit:$migration_dir")" "$migration_dir" >> "$output"
  done < <(
    git ls-tree -r --name-only "$commit" |
      awk -F/ '$2 == "database" && $3 == "migrations" { print $1 "/database/migrations" }' |
      sort -u
  )
}

prepare_commit_release() {
  local commit="$1"
  local release_dir="$2"

  require_command git
  require_command sha256sum
  require_command python3

  cd "$repo_root"
  git cat-file -e "$commit^{commit}" 2>/dev/null || die "release commit is unavailable: $commit"

  local tree namespace tag services deploy_sha zip
  tree="$(git rev-parse "$commit^{tree}")"
  namespace="${QHPRO_IMAGE_NAMESPACE:-bdspro-release}"
  tag="$commit"
  services="${QHPRO_RELEASE_SERVICES:-}"
  [[ -n "$services" ]] || die "QHPRO_RELEASE_SERVICES is required"

  deploy_sha="$(git show "$commit:shared/code/deploy.sh" | sha256sum | awk '{print $1}')"
  [[ "$deploy_sha" == "$protected_deploy_sha" ]] ||
    die "protected shared/code/deploy.sh changed in $commit: $deploy_sha"

  mkdir -p "$release_dir"
  zip="$release_dir/bdspro-backend-$commit.zip"

  source_manifest_for_commit "$commit" "$release_dir/SOURCE-MANIFEST.sha256"
  git archive --format=zip --output "$zip" "$commit"
  (cd "$release_dir" && sha256sum "$(basename "$zip")" > source-artifact.sha256)
  migration_manifest "$commit" "$release_dir/migration-trees.tsv"

  cat > "$release_dir/release.meta" <<EOF
commit=$commit
tree=$tree
namespace=$namespace
tag=$tag
services=$(tr ' ' ',' <<<"$services" | sed 's/,,*/,/g; s/^,//; s/,$//')
EOF

  (
    cd "$release_dir"
    sha256sum SOURCE-MANIFEST.sha256 migration-trees.tsv release.meta > release-evidence.sha256
  )

  echo "Release source prepared: $release_dir"
  echo "commit=$commit"
  echo "tree=$tree"
  echo "image_tag=$tag"
}

prepare_release() {
  local release_dir="$1"

  cd "$repo_root"
  git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "release must run from a Git worktree"
  git diff --quiet -- || die "tracked worktree changes are not allowed"
  git diff --cached --quiet -- || die "staged changes are not allowed"

  prepare_commit_release "$(git rev-parse HEAD)" "$release_dir"
}

capture_images() {
  local dir="$1"
  require_command docker
  require_command sha256sum
  read_metadata "$dir"

  : > "$dir/image-manifest.tsv"
  local csv service image image_id
  csv="$RELEASE_SERVICES"
  IFS=',' read -r -a release_services <<< "$csv"
  for service in "${release_services[@]}"; do
    [[ -n "$service" ]] || continue
    image="$RELEASE_NAMESPACE/$service-service:$RELEASE_TAG"
    image_id="$(docker image inspect --format '{{.Id}}' "$image" 2>/dev/null)" ||
      die "missing prebuilt release image: $image"
    [[ "$image_id" =~ ^sha256:[0-9a-f]{64}$ ]] || die "invalid image ID for $image"
    printf '%s\t%s\t%s\n' "$service" "$image" "$image_id" >> "$dir/image-manifest.tsv"
  done

  (cd "$dir" && sha256sum image-manifest.tsv > image-manifest.tsv.sha256)
  echo "Release images captured: $dir/image-manifest.tsv"
}

verify_source_artifact() {
  local dir="$1"
  local zip
  zip="$dir/bdspro-backend-$RELEASE_COMMIT.zip"
  [[ -f "$zip" ]] || die "missing source artifact: $zip"
  [[ -f "$dir/source-artifact.sha256" ]] || die "missing source artifact checksum"
  (cd "$dir" && sha256sum -c source-artifact.sha256 >/dev/null)

  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN
  python3 - "$zip" "$tmp" <<'PY'
import sys, zipfile
with zipfile.ZipFile(sys.argv[1]) as archive:
    archive.extractall(sys.argv[2])
PY
  (cd "$tmp" && sha256sum -c "$dir/SOURCE-MANIFEST.sha256" >/dev/null)
  rm -rf "$tmp"
  trap - RETURN
}

verify_release() {
  local dir="$1"
  require_command git
  require_command sha256sum
  require_command python3

  read_metadata "$dir"
  cd "$repo_root"

  git cat-file -e "$RELEASE_COMMIT^{commit}" 2>/dev/null ||
    die "release commit is unavailable in this repository"
  [[ "$(git rev-parse "$RELEASE_COMMIT^{tree}")" == "$RELEASE_TREE" ]] ||
    die "release tree does not match commit"

  [[ -f "$dir/release-evidence.sha256" ]] || die "missing release evidence checksum"
  (cd "$dir" && sha256sum -c release-evidence.sha256 >/dev/null)

  local expected_migrations
  expected_migrations="$(mktemp)"
  migration_manifest "$RELEASE_COMMIT" "$expected_migrations"
  cmp -s "$expected_migrations" "$dir/migration-trees.tsv" ||
    die "migration state does not match release commit"
  rm -f "$expected_migrations"

  verify_source_artifact "$dir"

  if [[ -f "$dir/image-manifest.tsv" || -f "$dir/image-manifest.tsv.sha256" ]]; then
    [[ -f "$dir/image-manifest.tsv" && -f "$dir/image-manifest.tsv.sha256" ]] ||
      die "image evidence is incomplete"
    (cd "$dir" && sha256sum -c image-manifest.tsv.sha256 >/dev/null)
  fi

  echo "Release evidence verified: commit=$RELEASE_COMMIT tree=$RELEASE_TREE"
}

verify_local_images() {
  local dir="$1"
  require_command docker
  read_metadata "$dir"

  [[ -f "$dir/image-manifest.tsv" ]] || die "release image manifest is missing"
  local service image expected_id actual_id
  while IFS=$'\t' read -r service image expected_id; do
    [[ -n "$service" && -n "$image" && -n "$expected_id" ]] || die "invalid image manifest row"
    actual_id="$(docker image inspect --format '{{.Id}}' "$image" 2>/dev/null)" ||
      die "release image is not loaded: $image"
    [[ "$actual_id" == "$expected_id" ]] ||
      die "release image identity mismatch: $image expected=$expected_id actual=$actual_id"
  done < "$dir/image-manifest.tsv"

  echo "Release image identities verified: $dir"
}

activate_release() {
  local dir="$1"
  verify_release "$dir"
  verify_local_images "$dir"
  read_metadata "$dir"

  [[ -f "$env_file" ]] || die "missing runtime env: $env_file"
  [[ -f "$compose_file" ]] || die "missing compose file: $compose_file"

  QHPRO_IMAGE_NAMESPACE="$RELEASE_NAMESPACE"   QHPRO_IMAGE_TAG="$RELEASE_TAG"     docker compose --env-file "$env_file" -f "$compose_file"       up -d --no-build --wait --wait-timeout 300

  echo "Release activated without rebuild: commit=$RELEASE_COMMIT"
}

rollback_release() {
  local current_dir="$1"
  local previous_dir="$2"

  verify_release "$current_dir"
  read_metadata "$current_dir"
  local current_commit="$RELEASE_COMMIT"

  verify_release "$previous_dir"
  read_metadata "$previous_dir"
  local previous_commit="$RELEASE_COMMIT"

  [[ "$current_commit" != "$previous_commit" ]] || die "rollback requires a different previous release"
  git merge-base --is-ancestor "$previous_commit" "$current_commit" ||
    die "previous release is not an ancestor of current release"
  cmp -s "$current_dir/migration-trees.tsv" "$previous_dir/migration-trees.tsv" ||
    die "automatic rollback refused: migration state differs"

  echo "Rollback compatibility verified: $current_commit -> $previous_commit"
  activate_release "$previous_dir"
}

usage() {
  cat <<'EOF'
Usage:
  release-candidate.sh prepare <release-dir>
  release-candidate.sh prepare-commit <commit-sha> <release-dir>
  release-candidate.sh capture-images <release-dir>
  release-candidate.sh verify <release-dir>
  release-candidate.sh verify-images <release-dir>
  release-candidate.sh activate <release-dir>
  release-candidate.sh rollback <current-release-dir> <previous-release-dir>
EOF
}

cmd="${1:-}"
case "$cmd" in
  prepare)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    prepare_release "$2"
    ;;
  prepare-commit)
    [[ $# -eq 3 ]] || { usage >&2; exit 2; }
    prepare_commit_release "$2" "$3"
    ;;
  capture-images)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    capture_images "$2"
    ;;
  verify)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    verify_release "$2"
    ;;
  verify-images)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    verify_release "$2"
    verify_local_images "$2"
    ;;
  activate)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    activate_release "$2"
    ;;
  rollback)
    [[ $# -eq 3 ]] || { usage >&2; exit 2; }
    rollback_release "$2" "$3"
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
