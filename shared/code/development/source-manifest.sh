#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$tool_dir/../../.." && pwd)"
output="${1:-SOURCE-MANIFEST.sha256}"

cd "$repo_root"

case "$output" in
  /*) output_path="$output" ;;
  *) output_path="$repo_root/$output" ;;
esac

mkdir -p "$(dirname "$output_path")"
tmp="$(mktemp "${output_path}.tmp.XXXXXX")"
trap 'rm -f "$tmp"' EXIT

while IFS= read -r -d '' path; do
  # The final manifest is generated evidence, never an input to its own identity.
  [[ "$path" == "SOURCE-MANIFEST.sha256" ]] && continue
  digest="$(git show ":$path" | sha256sum | awk '{print $1}')"
  printf '%s  ./%s\n' "$digest" "$path" >>"$tmp"
done < <(git ls-files -z | sort -z)

mv "$tmp" "$output_path"
trap - EXIT
printf 'Source manifest: %s\n' "$output_path"
