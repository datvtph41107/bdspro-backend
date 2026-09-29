#!/usr/bin/env bash
set -u

COMMAND="${1:-}"
SERVICE="${2:-}"
TOOL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$TOOL_DIR/../.." && pwd)"
SOURCE_PATH="$("$TOOL_DIR/source-path.sh" "$SERVICE" 2>/dev/null || true)"

case "$COMMAND" in
  swag)
    if [[ -z "$SERVICE" ]]; then
      echo "Usage: $0 swag <service>" >&2
      exit 2
    fi
    cd "$REPO_ROOT/$SOURCE_PATH" || exit 1
    swag init --parseDependency --output "$REPO_ROOT/api/gateway/docs/${SERVICE}"
    ;;

  wire)
    if [[ -z "$SERVICE" ]]; then
      echo "Usage: $0 wire <service>" >&2
      exit 2
    fi
    cd "$REPO_ROOT/$SOURCE_PATH" || exit 1
    if [[ "$SERVICE" == "file" || "$SERVICE" == "payment" ]]; then
      for retired in wire wire.go wire_gen.go; do
        if [[ -e "$retired" ]]; then
          echo "FAIL: File Wire authority is retired but $retired exists" >&2
          exit 1
        fi
      done
      echo "${SERVICE}-service uses explicit process composition; Wire authority is retired."
      exit 0
    fi
    if [[ "$SERVICE" == "notification" || "$SERVICE" == "tqd" || "$SERVICE" == "hub" || "$SERVICE" == "assistant" ]]; then
      echo "Generating Wire graph from ${SERVICE}'s process-input authority"
      (cd wire && wire gen)
      exit $?
    fi
    echo "Generating Wire graph in $(pwd)"
    go run "$TOOL_DIR/script/gen_wire.go" "$SERVICE"
    ;;

  social)
    cd "$REPO_ROOT/internal/social" || exit 1
    swag init --parseDependency --output "$REPO_ROOT/api/gateway/docs/social"
    go run "$TOOL_DIR/script/gen_wire.go" social
    ;;

  *)
    echo "Usage: $0 {wire|swag} <service>" >&2
    exit 2
    ;;
esac
