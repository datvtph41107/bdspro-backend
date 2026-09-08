#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

authorities=(
  "user:user-service/database/migrations"
  "organization:organization-service/migrate"
  "payment:payment-service/migrate"
  "tqd:tqd-service/migrate"
  "notification:notification-service/migrate"
  "file:file-service/migrate"
  "hub:hub-service/migrate"
  "bdspro:bdspro-service/infra/db/migrations/v2"
  "crm:crm-service/infra/db/migrate_v2"
)

fail=0
for entry in "${authorities[@]}"; do
  owner="${entry%%:*}"
  dir="${entry#*:}"
  if [[ ! -d "$dir" ]]; then
    echo "$owner: missing migration authority $dir" >&2
    fail=1
    continue
  fi

  declare -A up_name=()
  declare -A down_name=()
  versions=()

  while IFS= read -r -d '' file; do
    base="${file##*/}"
    if [[ "$base" =~ ^([0-9]{6})_(.+)\.(up|down)\.sql$ ]]; then
      version="${BASH_REMATCH[1]}"
      name="${BASH_REMATCH[2]}"
      direction="${BASH_REMATCH[3]}"
      number=$((10#$version))
      if [[ "$direction" == up ]]; then
        if [[ -n "${up_name[$number]:-}" ]]; then
          echo "$owner: duplicate up migration version $version" >&2
          fail=1
        fi
        up_name[$number]="$name"
      else
        if [[ -n "${down_name[$number]:-}" ]]; then
          echo "$owner: duplicate down migration version $version" >&2
          fail=1
        fi
        down_name[$number]="$name"
      fi
    elif [[ "$base" == *.sql ]]; then
      echo "$owner: non-canonical SQL in migration authority: $dir/$base" >&2
      fail=1
    fi
  done < <(find "$dir" -maxdepth 1 -type f -name '*.sql' -print0 | sort -z)

  mapfile -t versions < <(
    {
      for v in "${!up_name[@]}"; do echo "$v"; done
      for v in "${!down_name[@]}"; do echo "$v"; done
    } | sort -n -u
  )

  if [[ ${#versions[@]} -eq 0 ]]; then
    echo "$owner: no versioned migrations in $dir" >&2
    fail=1
    unset up_name down_name
    continue
  fi

  expected=1
  for version in "${versions[@]}"; do
    printf -v padded '%06d' "$version"
    if (( version != expected )); then
      printf '%s: migration version gap, expected %06d got %06d\n' "$owner" "$expected" "$version" >&2
      fail=1
      expected=$version
    fi
    if [[ -z "${up_name[$version]:-}" ]]; then
      echo "$owner: missing ${padded}_*.up.sql" >&2
      fail=1
    fi
    if [[ -z "${down_name[$version]:-}" ]]; then
      echo "$owner: missing ${padded}_*.down.sql" >&2
      fail=1
    fi
    if [[ -n "${up_name[$version]:-}" && -n "${down_name[$version]:-}" && "${up_name[$version]}" != "${down_name[$version]}" ]]; then
      echo "$owner: version $padded name mismatch: up=${up_name[$version]} down=${down_name[$version]}" >&2
      fail=1
    fi
    expected=$((version + 1))
  done

  echo "$owner: migration contract checked (${#versions[@]} versions)"
  unset up_name down_name
done

if (( fail != 0 )); then
  exit 1
fi

echo 'migration contract PASS'
