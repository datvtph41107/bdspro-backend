#!/bin/bash
# Deploy userv4 (HTTP by default) to the production host.
# Usage:
#   ./deploy.sh              # build cmd/http → remote compose service `v4`
#   ./deploy.sh http         # same as default
#   ./deploy.sh grpc         # build cmd/grpc as main
#   ./deploy.sh compose      # push shared/code/docker-compose.yml only
#   ./deploy.sh logs         # tail remote container logs
#   ./deploy.sh restart      # recreate remote container without rebuild
set -euo pipefail

start=$(date +%s)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
SHARED_CODE="${REPO_ROOT}/shared/code"

REMOTE_USER="${REMOTE_USER:-root}"
REMOTE_HOST="${REMOTE_HOST:-14.225.210.29}"
WORKSPACE_NAME="${WORKSPACE_NAME:-backend-microservice}"
REMOTE_PATH="/root/${WORKSPACE_NAME}"
SERVICE_PATH="v4"
PAYLOAD_BUILD="v4"
REMOTE_DIR="${REMOTE_PATH}/${SERVICE_PATH}"

TARGET="${1:-http}"

usage() {
  cat <<'EOF'
Usage: ./deploy.sh [http|grpc|compose|logs|restart]

  http     Build ./cmd/http, ship binary+config+migrate, recreate compose service v4 (default)
  grpc     Build ./cmd/grpc as main (same compose service; override command if needed)
  compose  Push shared/code/docker-compose.yml to the server
  logs     Tail remote v4 container logs
  restart  Force-recreate remote v4 without rebuilding
EOF
}

remote() {
  ssh "${REMOTE_USER}@${REMOTE_HOST}" "$@"
}

case "${TARGET}" in
  -h|--help|help)
    usage
    exit 0
    ;;
  compose)
    echo "Pushing docker-compose.yml → ${REMOTE_PATH}"
    remote "mkdir -p ${REMOTE_PATH}"
    scp "${SHARED_CODE}/docker-compose.yml" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/docker-compose.yml"
    echo "Done."
    exit 0
    ;;
  logs)
    remote "cd ${REMOTE_PATH} && docker compose logs -f --tail=200 ${PAYLOAD_BUILD}"
    exit 0
    ;;
  restart)
    remote "cd ${REMOTE_PATH} && docker compose up -d --force-recreate ${PAYLOAD_BUILD}"
    exit 0
    ;;
  http)
    BUILD_PACKAGE="./cmd/http"
    ;;
  grpc)
    BUILD_PACKAGE="./cmd/grpc"
    ;;
  *)
    echo "Unknown target: ${TARGET}"
    usage
    exit 1
    ;;
esac

cd "${SCRIPT_DIR}"

echo "--- GO BUILD (${BUILD_PACKAGE} → linux/amd64 main) ---"
if ! CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOWORK=off go build -trimpath -ldflags='-s -w' -o main "${BUILD_PACKAGE}"; then
  echo "Go build failed — aborting (old binary not pushed)."
  exit 1
fi
echo "--- BUILD OK ---"

echo "Preparing remote dirs..."
remote "
mkdir -p ${REMOTE_DIR}
mkdir -p ${REMOTE_DIR}/config
mkdir -p ${REMOTE_DIR}/migrate
rm -f ${REMOTE_DIR}/config/local.yml ${REMOTE_DIR}/config/product.yml ${REMOTE_DIR}/config/production.yml
"

echo "Shipping binary, Dockerfile, config, migrate..."
scp main "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/main"
scp "${SCRIPT_DIR}/Dockerfile.develop" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/Dockerfile"

if [ ! -d migrate ]; then
  mkdir -p migrate
fi
rsync -avz --delete migrate/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/migrate/"

rsync -avz --exclude='local.yml' --exclude='product.yml' --exclude='production.yml' \
  config/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/config/"

echo "--- REMOTE docker compose build/up (${PAYLOAD_BUILD}) ---"
remote "
set -e
cd ${REMOTE_PATH}
if ! docker compose config --services 2>/dev/null | grep -qx '${PAYLOAD_BUILD}'; then
  echo \"Compose service '${PAYLOAD_BUILD}' missing on server.\" >&2
  echo \"Run: ./deploy.sh compose   (after v4 is in shared/code/docker-compose.yml)\" >&2
  exit 1
fi
docker compose build --no-cache ${PAYLOAD_BUILD}
docker compose up -d --force-recreate ${PAYLOAD_BUILD}
docker compose ps ${PAYLOAD_BUILD}
"

rm -f main

end=$(date +%s)
echo "Deploy OK (${TARGET}) in $((end - start))s → ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}"
