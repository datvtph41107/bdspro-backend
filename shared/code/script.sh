#!/usr/bin/env bash
set -u

COMMAND="${1:-}"
SERVICE="${2:-}"

case "$COMMAND" in
  swag)
    if [[ -z "$SERVICE" ]]; then
      echo "Usage: $0 swag <service>" >&2
      exit 2
    fi
    cd "../../${SERVICE}-service" || exit 1
    swag init --parseDependency --output "../gateway-service/docs/${SERVICE}"
    ;;

  wire)
    if [[ -z "$SERVICE" ]]; then
      echo "Usage: $0 wire <service>" >&2
      exit 2
    fi
    cd "../../${SERVICE}-service" || exit 1
    echo "Generating Wire graph in $(pwd)"
    go run ../shared/code/script/gen_wire.go "$SERVICE"
    ;;

  social)
    cd social-service || exit 1
    swag init --parseDependency --output ../../gateway-service/docs/social
    go run script/gen_wire.go social
    ;;

  *)
    echo "Usage: $0 {wire|swag} <service>" >&2
    exit 2
    ;;
esac
