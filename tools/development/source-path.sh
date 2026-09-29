#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  gateway) echo "api/gateway" ;;
  user) echo "internal/user" ;;
  auth) echo "internal/authentication" ;;
  organization) echo "internal/organization-legacy" ;;
  bdspro) echo "internal/property" ;;
  crm) echo "internal/crm" ;;
  tqd) echo "internal/planning" ;;
  payment) echo "internal/payment" ;;
  file) echo "internal/file" ;;
  notification) echo "internal/notification" ;;
  hub) echo "internal/hub" ;;
  assistant) echo "internal/assistant" ;;
  chat) echo "internal/chat" ;;
  chat-v1) echo "internal/chat-legacy" ;;
  social) echo "internal/social" ;;
  map) echo "internal/map-legacy" ;;
  relay) echo "internal/realtime-relay" ;;
  search) echo "internal/search" ;;
  ai) echo "internal/ai" ;;
  base) echo "infrastructure/base" ;;
  runtime) echo "infrastructure/runtime" ;;
  proto) echo "proto" ;;
  tools) echo "tools/development" ;;
  "") echo "usage: source-path.sh <service>" >&2; exit 2 ;;
  *) echo "${1}-service" ;;
esac
