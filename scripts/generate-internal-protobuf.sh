#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
schema_dir="${repo_root}/shared/protobuf/schema"

cd "${schema_dir}"
buf generate \
  --template buf.internal.gen.yaml \
  --path user/access_internal.proto \
  --path file/internal.proto

for forbidden in \
  ../types/user/access_internal.pb.gw.go \
  ../types/file/internal.pb.gw.go; do
  if [[ -f "${forbidden}" ]]; then
    echo "internal proto unexpectedly generated HTTP gateway code: ${forbidden}" >&2
    exit 1
  fi
done

echo "Generated internal user/file Go and gRPC types without HTTP handlers."
