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
    if [[ "$SERVICE" == "notification" || "$SERVICE" == "tqd" || "$SERVICE" == "hub" ]]; then
      echo "Generating Wire graph from ${SERVICE}'s process-input authority"
      (cd wire && wire gen)
      exit $?
    fi
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
