#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
default_repo_root="$(cd "$tool_dir/../../.." && pwd)"
source_root="${QHPRO_RELEASE_SOURCE_ROOT:-$default_repo_root}"
compose_file="${QHPRO_RELEASE_COMPOSE_FILE:-$source_root/compose.yaml}"
env_file="${QHPRO_RELEASE_ENV_FILE:-$source_root/.env}"
protected_deploy_sha="80b70e3ea1375b4a959438da92011574c084bdb14d6ee2399ff3d7ddf023e56e"

die() {
  echo "release-artifact: $*" >&2
  exit 2
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "missing command: $1"
}

sha256_value() {
  sha256sum "$1" | awk '{print $1}'
}

read_metadata() {
  local dir="$1"
  local file="$dir/release.meta"
  [[ -f "$file" ]] || die "missing release metadata: $file"

  RELEASE_FORMAT=
  RELEASE_COMMIT=
  RELEASE_TREE=
  RELEASE_NAMESPACE=
  RELEASE_TAG=
  RELEASE_MIGRATION_DIGEST=

  while IFS='=' read -r key value; do
    case "$key" in
      format) RELEASE_FORMAT="$value" ;;
      commit) RELEASE_COMMIT="$value" ;;
      tree) RELEASE_TREE="$value" ;;
      namespace) RELEASE_NAMESPACE="$value" ;;
      tag) RELEASE_TAG="$value" ;;
      migration_digest) RELEASE_MIGRATION_DIGEST="$value" ;;
      '') ;;
      *) die "unknown release metadata key: $key" ;;
    esac
  done < "$file"

  [[ "$RELEASE_FORMAT" == "1" ]] || die "unsupported release metadata format: $RELEASE_FORMAT"
  [[ "$RELEASE_COMMIT" =~ ^[0-9a-f]{40}$ ]] || die "invalid release commit"
  [[ "$RELEASE_TREE" =~ ^[0-9a-f]{40}$ ]] || die "invalid release tree"
  [[ -n "$RELEASE_NAMESPACE" ]] || die "release namespace is empty"
  [[ "$RELEASE_TAG" == "$RELEASE_COMMIT" ]] || die "release tag must equal exact commit SHA"
  [[ "$RELEASE_MIGRATION_DIGEST" =~ ^[0-9a-f]{64}$ ]] || die "invalid migration digest"
}

migration_manifest() {
  local commit="$1"
  local output="$2"
  local migration_dir
  : > "$output"

  while IFS= read -r migration_dir; do
    [[ -n "$migration_dir" ]] || continue
    printf '%s\t%s\n' "$(git -C "$source_root" rev-parse "$commit:$migration_dir")" "$migration_dir" >> "$output"
  done < <(
    git -C "$source_root" ls-tree -r --name-only "$commit" |
      awk -F/ '
        $1 ~ /-service$/ && $2 == "database" && $3 == "migrations" {
          print $1 "/database/migrations"
        }
        $1 ~ /-service$/ && $2 == "migrate" {
          print $1 "/migrate"
        }
      ' |
      sort -u
  )
}

verify_protected_source() {
  local commit="$1"
  local deploy_sha
  deploy_sha="$(git -C "$source_root" show "$commit:shared/code/deploy.sh" | sha256sum | awk '{print $1}')"
  [[ "$deploy_sha" == "$protected_deploy_sha" ]] ||
    die "protected shared/code/deploy.sh changed in $commit: $deploy_sha"
}

prepare_release() {
  local release_dir="$1"

  require_command git
  require_command sha256sum
  require_command python3

  git -C "$source_root" rev-parse --is-inside-work-tree >/dev/null 2>&1 ||
    die "release source root is not a Git worktree: $source_root"
  git -C "$source_root" diff --quiet -- ||
    die "tracked worktree changes are not allowed"
  git -C "$source_root" diff --cached --quiet -- ||
    die "staged changes are not allowed"

  local commit tree namespace tag source_zip migration_digest
  commit="$(git -C "$source_root" rev-parse HEAD)"
  tree="$(git -C "$source_root" rev-parse HEAD^{tree})"
  namespace="${QHPRO_IMAGE_NAMESPACE:-bdspro-release}"
  tag="$commit"

  [[ "$commit" =~ ^[0-9a-f]{40}$ ]] || die "invalid source commit"
  verify_protected_source "$commit"
  [[ -f "$source_root/shared/code/development/source-manifest.sh" ]] ||
    die "source manifest owner is missing"
  [[ -f "$compose_file" ]] || die "compose file is missing: $compose_file"
  [[ -f "$env_file" ]] || die "runtime env is missing: $env_file"

  if [[ -e "$release_dir" ]] && [[ -n "$(find "$release_dir" -mindepth 1 -maxdepth 1 -print -quit 2>/dev/null)" ]]; then
    die "release directory is not empty: $release_dir"
  fi
  mkdir -p "$release_dir"

  bash "$source_root/shared/code/development/source-manifest.sh"     "$release_dir/SOURCE-MANIFEST.sha256"

  source_zip="$release_dir/bdspro-backend-$commit.zip"
  git -C "$source_root" archive --format=zip --output "$source_zip" "$commit"
  (
    cd "$release_dir"
    sha256sum "$(basename "$source_zip")" > source-artifact.sha256
  )

  migration_manifest "$commit" "$release_dir/migration-trees.tsv"
  (
    cd "$release_dir"
    sha256sum migration-trees.tsv > migration-trees.sha256
  )
  migration_digest="$(sha256_value "$release_dir/migration-trees.tsv")"

  cat > "$release_dir/release.meta" <<EOF
format=1
commit=$commit
tree=$tree
namespace=$namespace
tag=$tag
migration_digest=$migration_digest
EOF

  (
    cd "$release_dir"
    sha256sum       SOURCE-MANIFEST.sha256       source-artifact.sha256       migration-trees.tsv       migration-trees.sha256       release.meta       > source-evidence.sha256
  )

  echo "Release source prepared: $release_dir"
  echo "commit=$commit"
  echo "tree=$tree"
  echo "namespace=$namespace"
  echo "tag=$tag"
  echo "migration_digest=$migration_digest"
}

compose_images() {
  read_metadata "$1"
  require_command docker
  [[ -f "$compose_file" ]] || die "compose file is missing: $compose_file"
  [[ -f "$env_file" ]] || die "runtime env is missing: $env_file"

  QHPRO_IMAGE_NAMESPACE="$RELEASE_NAMESPACE"   QHPRO_IMAGE_TAG="$RELEASE_TAG"     docker compose --env-file "$env_file" -f "$compose_file" config --images |
    sed '/^[[:space:]]*$/d' |
    sort -u
}

capture_images() {
  local dir="$1"
  require_command docker
  require_command sha256sum
  read_metadata "$dir"

  local -a images=()
  mapfile -t images < <(compose_images "$dir")
  [[ "${#images[@]}" -gt 0 ]] || die "compose image set is empty"

  : > "$dir/image-manifest.tsv"
  local image image_id
  for image in "${images[@]}"; do
    image_id="$(docker image inspect --format '{{.Id}}' "$image" 2>/dev/null)" ||
      die "release image is not available locally: $image"
    [[ "$image_id" =~ ^sha256:[0-9a-f]{64}$ ]] ||
      die "invalid image ID for $image: $image_id"
    printf '%s\t%s\n' "$image" "$image_id" >> "$dir/image-manifest.tsv"
  done

  docker save --output "$dir/images.tar" "${images[@]}"
  [[ -s "$dir/images.tar" ]] || die "Docker image archive is empty"

  (
    cd "$dir"
    sha256sum image-manifest.tsv > image-manifest.tsv.sha256
    sha256sum images.tar > images.tar.sha256
    sha256sum       SOURCE-MANIFEST.sha256       source-artifact.sha256       migration-trees.tsv       migration-trees.sha256       release.meta       source-evidence.sha256       image-manifest.tsv       image-manifest.tsv.sha256       images.tar.sha256       > release-evidence.sha256
  )

  echo "Release image archive captured: $dir/images.tar"
  echo "image_count=${#images[@]}"
}

verify_source_artifact() {
  local dir="$1"
  local source_zip="$dir/bdspro-backend-$RELEASE_COMMIT.zip"
  [[ -f "$source_zip" ]] || die "missing source artifact: $source_zip"

  (
    cd "$dir"
    sha256sum -c source-artifact.sha256 >/dev/null
  ) || die "source artifact checksum mismatch"

  local tmp
  tmp="$(mktemp -d)"
  if ! python3 - "$source_zip" "$tmp" <<'PY'
import sys
import zipfile

with zipfile.ZipFile(sys.argv[1]) as archive:
    archive.extractall(sys.argv[2])
PY
  then
    rm -rf "$tmp"
    die "cannot extract source artifact"
  fi

  if ! (cd "$tmp" && sha256sum -c "$dir/SOURCE-MANIFEST.sha256" >/dev/null); then
    rm -rf "$tmp"
    die "source manifest does not match source artifact"
  fi
  rm -rf "$tmp"
}

verify_release() {
  local dir="$1"

  require_command git
  require_command sha256sum
  require_command python3
  read_metadata "$dir"

  git -C "$source_root" cat-file -e "$RELEASE_COMMIT^{commit}" 2>/dev/null ||
    die "release commit is unavailable in source repository"
  [[ "$(git -C "$source_root" rev-parse "$RELEASE_COMMIT^{tree}")" == "$RELEASE_TREE" ]] ||
    die "release tree does not match release commit"
  verify_protected_source "$RELEASE_COMMIT"

  for file in     SOURCE-MANIFEST.sha256     source-artifact.sha256     migration-trees.tsv     migration-trees.sha256     release.meta     source-evidence.sha256     image-manifest.tsv     image-manifest.tsv.sha256     images.tar     images.tar.sha256     release-evidence.sha256; do
    [[ -f "$dir/$file" ]] || die "missing release artifact file: $file"
  done

  (
    cd "$dir" || exit 1
    sha256sum -c source-evidence.sha256 >/dev/null || exit 1
    sha256sum -c image-manifest.tsv.sha256 >/dev/null || exit 1
    sha256sum -c images.tar.sha256 >/dev/null || exit 1
    sha256sum -c release-evidence.sha256 >/dev/null || exit 1
  ) || die "release evidence checksum mismatch"

  local actual_migration_digest
  actual_migration_digest="$(sha256_value "$dir/migration-trees.tsv")"
  [[ "$actual_migration_digest" == "$RELEASE_MIGRATION_DIGEST" ]] ||
    die "release migration digest mismatch"

  local expected_migrations
  expected_migrations="$(mktemp)"
  migration_manifest "$RELEASE_COMMIT" "$expected_migrations"
  if ! cmp -s "$expected_migrations" "$dir/migration-trees.tsv"; then
    rm -f "$expected_migrations"
    die "migration manifest does not match release commit"
  fi
  rm -f "$expected_migrations"

  verify_source_artifact "$dir"

  local expected_images
  expected_images="$(mktemp)"
  compose_images "$dir" > "$expected_images"
  if ! diff -u     "$expected_images"     <(cut -f1 "$dir/image-manifest.tsv" | sort -u); then
    rm -f "$expected_images"
    die "image manifest does not match Compose image set"
  fi
  rm -f "$expected_images"

  echo "Release artifact verified: commit=$RELEASE_COMMIT tree=$RELEASE_TREE"
}

restore_images() {
  local dir="$1"
  require_command docker
  verify_release "$dir"
  read_metadata "$dir"

  docker load --input "$dir/images.tar" >/dev/null

  local image expected_id actual_id
  while IFS=$'\t' read -r image expected_id; do
    [[ -n "$image" && -n "$expected_id" ]] || die "invalid image manifest row"
    actual_id="$(docker image inspect --format '{{.Id}}' "$image" 2>/dev/null)" ||
      die "archived release image did not load: $image"
    [[ "$actual_id" == "$expected_id" ]] ||
      die "loaded image identity mismatch: $image expected=$expected_id actual=$actual_id"
  done < "$dir/image-manifest.tsv"

  echo "Release image archive restored and verified: $dir"
}

runtime_environment() {
  [[ -f "$env_file" ]] || die "runtime env is missing: $env_file"
  (
    set -a
    # shellcheck disable=SC1090
    source "$env_file"
    printf '%s' "${QHPRO_ENVIRONMENT:-}"
  )
}

assert_nonproduction_environment() {
  local environment
  environment="$(runtime_environment)"
  environment="$(tr '[:upper:]' '[:lower:]' <<<"$environment" | xargs)"
  case "$environment" in
    development|test|ci|acceptance|staging)
      ;;
    production|prod)
      die "production activation is forbidden from the acceptance repository"
      ;;
    *)
      die "unsupported or missing QHPRO_ENVIRONMENT for release activation: $environment"
      ;;
  esac
  echo "Release environment guard PASS: $environment"
}

activate_release() {
  local dir="$1"
  verify_release "$dir"
  assert_nonproduction_environment
  restore_images "$dir"
  read_metadata "$dir"

  QHPRO_IMAGE_NAMESPACE="$RELEASE_NAMESPACE"   QHPRO_IMAGE_TAG="$RELEASE_TAG"     docker compose --env-file "$env_file" -f "$compose_file"       up -d --no-build --pull never --wait --wait-timeout 300

  echo "Release activated without build/pull: commit=$RELEASE_COMMIT"
}

rollback_release() {
  local current_dir="$1"
  local previous_dir="$2"

  verify_release "$current_dir"
  read_metadata "$current_dir"
  local current_commit="$RELEASE_COMMIT"
  local current_migration_digest="$RELEASE_MIGRATION_DIGEST"

  verify_release "$previous_dir"
  read_metadata "$previous_dir"
  local previous_commit="$RELEASE_COMMIT"
  local previous_migration_digest="$RELEASE_MIGRATION_DIGEST"

  [[ "$current_commit" != "$previous_commit" ]] ||
    die "rollback requires a different previous release"
  git -C "$source_root" merge-base --is-ancestor "$previous_commit" "$current_commit" ||
    die "previous release is not an ancestor of current release"
  [[ "$current_migration_digest" == "$previous_migration_digest" ]] ||
    die "automatic rollback refused: migration state differs"

  echo "Rollback compatibility verified: $current_commit -> $previous_commit"
  activate_release "$previous_dir"
}

usage() {
  cat <<'EOF'
Usage:
  release-artifact.sh prepare <release-dir>
  release-artifact.sh capture-images <release-dir>
  release-artifact.sh verify <release-dir>
  release-artifact.sh restore-images <release-dir>
  release-artifact.sh activate <release-dir>
  release-artifact.sh rollback <current-release-dir> <previous-release-dir>

Optional environment:
  QHPRO_RELEASE_SOURCE_ROOT   source worktree to package/verify
  QHPRO_RELEASE_COMPOSE_FILE  Compose file for that worktree
  QHPRO_RELEASE_ENV_FILE      runtime env for that worktree
  QHPRO_IMAGE_NAMESPACE       release image namespace during prepare
EOF
}

cmd="${1:-}"
case "$cmd" in
  prepare)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    prepare_release "$2"
    ;;
  capture-images)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    capture_images "$2"
    ;;
  verify)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    verify_release "$2"
    ;;
  restore-images)
    [[ $# -eq 2 ]] || { usage >&2; exit 2; }
    restore_images "$2"
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
