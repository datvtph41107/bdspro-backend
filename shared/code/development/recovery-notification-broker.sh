#!/usr/bin/env bash
set -euo pipefail

compose_file=${1:?compose file is required}
env_file=${2:?env file is required}
script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
native_stack="$script_dir/native-stack.sh"

compose=(docker compose --env-file "$env_file" -f "$compose_file")
notification_pid=$("$native_stack" pid notification 2>/dev/null || true)
if [[ -z "$notification_pid" ]]; then
  echo 'Native Notification process is not running' >&2
  exit 1
fi

"${compose[@]}" restart rabbitmq
"${compose[@]}" up -d --no-deps --no-build --wait --wait-timeout 120 rabbitmq

# A broker restart must not restart the native Notification API. The consumer
# actor owns reconnection and must establish its named AMQP connection again.
for _ in $(seq 1 60); do
  current_pid=$("$native_stack" pid notification 2>/dev/null || true)
  if [[ "$current_pid" == "$notification_pid" ]] && kill -0 "$current_pid" 2>/dev/null && \
    "${compose[@]}" exec -T rabbitmq rabbitmqctl -q list_connections name 2>/dev/null | rg -q 'notification-service'; then
    echo 'Notification broker recovery PASS (native API stayed alive; consumer reconnected)'
    exit 0
  fi
  sleep 1
done

"${compose[@]}" ps rabbitmq >&2 || true
FOLLOW=0 TAIL=120 "$native_stack" logs notification >&2 || true
"${compose[@]}" logs --tail=120 rabbitmq >&2 || true
echo 'Notification broker recovery FAIL' >&2
exit 1
