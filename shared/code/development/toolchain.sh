#!/usr/bin/env bash
set -euo pipefail

tool_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$tool_dir/../../.." && pwd)"
# shellcheck disable=SC1091
source "$tool_dir/toolchain.versions"

toolchain_root="${QHPRO_TOOLCHAIN_ROOT:-${XDG_CACHE_HOME:-$HOME/.cache}/qhpro/toolchain}"
download_dir="$toolchain_root/downloads"
go_root="$toolchain_root/go${GO_VERSION}"
tool_bin="$toolchain_root/bin"
protoc_root="$toolchain_root/protoc${PROTOC_VERSION}"
weasyprint_root="$toolchain_root/weasyprint${WEASYPRINT_VERSION}"

usage() {
  cat <<USAGE
Usage: $0 <print|verify|bootstrap|generate-commerce>

Commands:
  print              Print the pinned versions and installation root.
  verify             Verify every repository-owned tool is available.
  bootstrap          Install the pinned Linux x86-64 toolchain.
  generate-commerce  Compatibility command for the commercial proto subset.
USAGE
}

require_bootstrap_host() {
  [[ "$(uname -s)" == "Linux" && "$(uname -m)" == "x86_64" ]] || {
    echo "QHPRO bootstrap currently supports Linux/WSL2 x86-64." >&2
    echo "Other hosts may provide the same pinned tools and run '$0 verify'." >&2
    return 1
  }
  command -v sha256sum >/dev/null 2>&1 || { echo "sha256sum is required" >&2; return 1; }
  command -v tar >/dev/null 2>&1 || { echo "tar is required" >&2; return 1; }
  command -v unzip >/dev/null 2>&1 || { echo "unzip is required" >&2; return 1; }
  command -v python3 >/dev/null 2>&1 || { echo "python3 is required" >&2; return 1; }
  python3 -m venv --help >/dev/null 2>&1 || { echo "python3 venv support is required" >&2; return 1; }
  command -v curl >/dev/null 2>&1 || command -v wget >/dev/null 2>&1 || {
    echo "curl or wget is required" >&2
    return 1
  }
}

download() {
  local url="$1" target="$2"
  mkdir -p "$(dirname "$target")"
  if command -v curl >/dev/null 2>&1; then
    curl --fail --location --retry 3 --output "$target" "$url"
  else
    wget -O "$target" "$url"
  fi
}

verify_checksum() {
  local expected="$1" file="$2" actual
  actual="$(sha256sum "$file" | awk '{print $1}')"
  [[ "$actual" == "$expected" ]] || {
    echo "Checksum mismatch: $file" >&2
    echo "Expected: $expected" >&2
    echo "Actual:   $actual" >&2
    return 1
  }
}

print_versions() {
  cat <<VERSIONS
Go=$GO_VERSION
Buf=$BUF_VERSION
protoc=$PROTOC_VERSION
protoc-gen-go=$PROTOC_GEN_GO_VERSION
protoc-gen-go-grpc=$PROTOC_GEN_GO_GRPC_VERSION
protoc-gen-grpc-gateway=$GRPC_GATEWAY_VERSION
Wire=$WIRE_VERSION
migrate=$MIGRATE_VERSION
Air=$AIR_VERSION
grpcurl=$GRPCURL_VERSION
WeasyPrint=$WEASYPRINT_VERSION
pydyf=$WEASYPRINT_PYDYF_VERSION
root=$toolchain_root
VERSIONS
}

verify_executable() {
  local name="$1" path="$2"
  if [[ ! -x "$path" ]]; then
    printf '[MISSING] %-20s %s\n' "$name" "$path"
    return 1
  fi
  printf '[PASS]    %-20s %s\n' "$name" "$path"
}

verify_go_binary() {
  local name="$1" module="$2" expected="v$3" path="$4" actual=""
  if [[ ! -x "$path" ]]; then
    printf '[MISSING] %-20s %s\n' "$name" "$path"
    return 1
  fi
  actual="$("$go_root/bin/go" version -m "$path" | awk -v module="$module" '$1 == "mod" && $2 == module { print $3; exit }')"
  if [[ "$actual" != "$expected" ]]; then
    printf '[MISMATCH] %-19s expected %s, found %s\n' "$name" "$expected" "${actual:-unknown}"
    return 1
  fi
  printf '[PASS]    %-20s %s\n' "$name $3" "$path"
}

verify_protoc() {
  local actual=""
  if [[ ! -x "$protoc_root/bin/protoc" ]]; then
    printf '[MISSING] %-20s %s\n' "protoc $PROTOC_VERSION" "$protoc_root/bin/protoc"
    return 1
  fi
  actual="$("$protoc_root/bin/protoc" --version)"
  if [[ "$actual" != "libprotoc $PROTOC_VERSION" ]]; then
    printf '[MISMATCH] %-19s expected %s, found %s\n' "protoc" "$PROTOC_VERSION" "$actual"
    return 1
  fi
  printf '[PASS]    %-20s %s\n' "protoc $PROTOC_VERSION" "$protoc_root/bin/protoc"
}

verify_weasyprint() {
  local path="$tool_bin/weasyprint" actual="" pydyf_version=""
  if [[ ! -x "$path" ]]; then
    printf '[MISSING] %-20s %s\n' "WeasyPrint $WEASYPRINT_VERSION" "$path"
    return 1
  fi
  actual="$($path --version 2>/dev/null || true)"
  if [[ "$actual" != "WeasyPrint version $WEASYPRINT_VERSION" ]]; then
    printf '[MISMATCH] %-19s expected %s, found %s\n' "WeasyPrint" "$WEASYPRINT_VERSION" "${actual:-unknown}"
    return 1
  fi
  pydyf_version="$($weasyprint_root/bin/python -c 'import pydyf; print(pydyf.__version__)' 2>/dev/null || true)"
  if [[ "$pydyf_version" != "$WEASYPRINT_PYDYF_VERSION" ]]; then
    printf '[MISMATCH] %-19s expected %s, found %s\n' "pydyf" "$WEASYPRINT_PYDYF_VERSION" "${pydyf_version:-unknown}"
    return 1
  fi
  printf '[PASS]    %-20s %s\n' "WeasyPrint $WEASYPRINT_VERSION" "$path"
  printf '[PASS]    %-20s %s\n' "pydyf $WEASYPRINT_PYDYF_VERSION" "$weasyprint_root"
}

verify_migrate_drivers() {
  local help=""
  if [[ ! -x "$tool_bin/migrate" ]]; then
    return 1
  fi
  help="$("$tool_bin/migrate" -help 2>&1)"
  if ! grep -Eq 'Database drivers:.*(^|,| )postgres(,| |$)' <<<"$help"; then
    printf '[MISMATCH] %-19s PostgreSQL database driver is missing\n' "migrate drivers"
    return 1
  fi
  printf '[PASS]    %-20s %s\n' "migrate postgres" "$tool_bin/migrate"
}

verify() {
  local status=0 selected=""
  verify_executable "go $GO_VERSION" "$go_root/bin/go" || status=1
  if [[ -x "$go_root/bin/go" ]]; then
    selected="$(GOROOT="$go_root" GOTOOLCHAIN=local "$go_root/bin/go" env GOVERSION)"
    [[ "$selected" == "go$GO_VERSION" ]] || {
      echo "Pinned Go mismatch: expected go$GO_VERSION, found $selected" >&2
      status=1
    }
    verify_go_binary "buf" "github.com/bufbuild/buf" "$BUF_VERSION" "$tool_bin/buf" || status=1
    verify_go_binary "protoc-gen-go" "google.golang.org/protobuf" "$PROTOC_GEN_GO_VERSION" "$tool_bin/protoc-gen-go" || status=1
    verify_go_binary "protoc-gen-go-grpc" "google.golang.org/grpc/cmd/protoc-gen-go-grpc" "$PROTOC_GEN_GO_GRPC_VERSION" "$tool_bin/protoc-gen-go-grpc" || status=1
    verify_go_binary "protoc-gen-grpc-gateway" "github.com/grpc-ecosystem/grpc-gateway/v2" "$GRPC_GATEWAY_VERSION" "$tool_bin/protoc-gen-grpc-gateway" || status=1
    verify_go_binary "wire" "github.com/google/wire" "$WIRE_VERSION" "$tool_bin/wire" || status=1
    verify_go_binary "migrate" "github.com/golang-migrate/migrate/v4" "$MIGRATE_VERSION" "$tool_bin/migrate" || status=1
    verify_migrate_drivers || status=1
    verify_go_binary "air" "github.com/air-verse/air" "$AIR_VERSION" "$tool_bin/air" || status=1
    verify_go_binary "grpcurl" "github.com/fullstorydev/grpcurl" "$GRPCURL_VERSION" "$tool_bin/grpcurl" || status=1
  fi
  verify_protoc || status=1
  verify_weasyprint || status=1
  return "$status"
}

bootstrap_go() {
  local archive="$download_dir/go${GO_VERSION}.linux-amd64.tar.gz"
  if [[ ! -f "$archive" ]]; then
    download "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" "$archive"
  fi
  verify_checksum "$GO_LINUX_AMD64_SHA256" "$archive"
  if [[ -x "$go_root/bin/go" ]] &&
     [[ "$(GOROOT="$go_root" GOTOOLCHAIN=local "$go_root/bin/go" env GOVERSION)" == "go$GO_VERSION" ]]; then
    echo "[KEEP] Go $GO_VERSION"
    return
  fi
  rm -rf "$go_root"
  mkdir -p "$go_root"
  tar -xzf "$archive" --strip-components=1 -C "$go_root"
}

bootstrap_protoc() {
  local archive="$download_dir/protoc-${PROTOC_VERSION}-linux-x86_64.zip"
  if [[ ! -f "$archive" ]]; then
    download "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" "$archive"
  fi
  verify_checksum "$PROTOC_LINUX_X86_64_SHA256" "$archive"
  if [[ -x "$protoc_root/bin/protoc" ]] &&
     "$protoc_root/bin/protoc" --version | grep -q "libprotoc ${PROTOC_VERSION}"; then
    echo "[KEEP] protoc $PROTOC_VERSION"
    return
  fi
  rm -rf "$protoc_root"
  mkdir -p "$protoc_root"
  unzip -q "$archive" -d "$protoc_root"
}

install_go_tools() {
  local go="$go_root/bin/go"
  mkdir -p "$tool_bin"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "github.com/google/wire/cmd/wire@v${WIRE_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install -tags postgres "github.com/golang-migrate/migrate/v4/cmd/migrate@v${MIGRATE_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "github.com/air-verse/air@v${AIR_VERSION}"
  GOROOT="$go_root" GOTOOLCHAIN=local GOBIN="$tool_bin" "$go" install "github.com/fullstorydev/grpcurl/cmd/grpcurl@v${GRPCURL_VERSION}"
}

bootstrap_weasyprint() {
  local binary="$weasyprint_root/bin/weasyprint" pydyf_version=""
  if [[ -x "$weasyprint_root/bin/python" ]]; then
    pydyf_version="$($weasyprint_root/bin/python -c 'import pydyf; print(pydyf.__version__)' 2>/dev/null || true)"
  fi
  if [[ -x "$binary" ]] && \
     [[ "$($binary --version 2>/dev/null || true)" == "WeasyPrint version $WEASYPRINT_VERSION" ]] && \
     [[ "$pydyf_version" == "$WEASYPRINT_PYDYF_VERSION" ]]; then
    echo "[KEEP] WeasyPrint $WEASYPRINT_VERSION"
  else
    rm -rf "$weasyprint_root"
    python3 -m venv "$weasyprint_root"
    "$weasyprint_root/bin/python" -m pip install \
      --disable-pip-version-check \
      "weasyprint==$WEASYPRINT_VERSION" \
      "pydyf==$WEASYPRINT_PYDYF_VERSION"
  fi
  ln -sfn "$binary" "$tool_bin/weasyprint"
}

bootstrap() {
  require_bootstrap_host
  mkdir -p "$download_dir" "$tool_bin"
  bootstrap_go
  bootstrap_protoc
  install_go_tools
  bootstrap_weasyprint
  verify
}

generate_commerce() {
  verify
  local schema="$repo_root/shared/protobuf/schema"
  (
    cd "$schema"
    PATH="$tool_bin:$protoc_root/bin:$go_root/bin:$PATH" "$tool_bin/buf" generate \
      --template buf.internal.gen.yaml \
      --path payment/commerce_internal.proto \
      --path user/subscription_internal.proto \
      --path user/checkout.proto
  )
}

case "${1:-}" in
  print) print_versions ;;
  verify) verify ;;
  bootstrap) bootstrap ;;
  generate-commerce) generate_commerce ;;
  *) usage; exit 2 ;;
esac
