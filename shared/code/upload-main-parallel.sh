#!/bin/bash
# Upload binary (main) theo từng part, song song + resume.
# Part đã có trên remote đúng size sẽ bị bỏ qua — kill giữa chừng chạy lại được.
#
# Usage:
#   ./upload-main-parallel.sh gateway
#   ./upload-main-parallel.sh gateway /path/to/main
#   PARTS=8 PARALLEL=8 ./upload-main-parallel.sh gateway
#   REMOTE_HOST=103.172.239.89 ./upload-main-parallel.sh tqd

set -euo pipefail

REMOTE_USER="${REMOTE_USER:-root}"
REMOTE_HOST="${REMOTE_HOST:-103.172.239.89}"
WORKSPACE_NAME="${WORKSPACE_NAME:-golang-microservice}"
REMOTE_PATH="${REMOTE_PATH:-/root/$WORKSPACE_NAME}"
PARTS="${PARTS:-8}"
PARALLEL="${PARALLEL:-8}"

if [ -z "${1:-}" ]; then
  echo "Usage: $0 <service> [local_binary]"
  echo "  service: gateway | user | tqd | ..."
  echo "  env: PARTS=8 PARALLEL=8 REMOTE_HOST=... REMOTE_USER=..."
  exit 1
fi

service="$1"
case "$service" in
  gateway) service_path=gateway-service ;;
  user|user-grpc) service_path=user-service ;;
  auth) service_path=auth-service ;;
  file) service_path=file-service ;;
  map) service_path=map-service ;;
  chat|chathttp|swagger-chat) service_path=chat-service ;;
  payment) service_path=payment-service ;;
  notification) service_path=notification-service ;;
  crm) service_path=crm-service ;;
  membership) service_path=membership-service ;;
  bdspro|bdspro-grpc) service_path=bdspro-service ;;
  relay) service_path=relay-service ;;
  social) service_path=social-service ;;
  appointment) service_path=appointment-service ;;
  assistant) service_path=assistant-service ;;
  org) service_path=organization-service ;;
  tqd|pmtiles) service_path=tqd-service ;;
  feedback) service_path=feedback-service ;;
  marketing) service_path=marketing-service ;;
  transaction) service_path=transaction-service ;;
  hub) service_path=hub-service ;;
  *)
    echo "Unknown service: $service"
    exit 1
    ;;
esac

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOCAL_BIN="${2:-$REPO_ROOT/$service_path/main}"
REMOTE_DIR="$REMOTE_PATH/$service_path"
REMOTE="${REMOTE_USER}@${REMOTE_HOST}"

if [ ! -f "$LOCAL_BIN" ]; then
  echo "Binary not found: $LOCAL_BIN"
  echo "Build trước: (cd $REPO_ROOT/$service_path && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .)"
  exit 1
fi

if ! [[ "$PARTS" =~ ^[0-9]+$ ]] || [ "$PARTS" -lt 1 ]; then
  echo "PARTS must be a positive integer"
  exit 1
fi

file_size() {
  local path="$1"
  if stat -f%z "$path" >/dev/null 2>&1; then
    stat -f%z "$path"
  else
    stat -c%s "$path"
  fi
}

TOTAL_SIZE="$(file_size "$LOCAL_BIN")"
if stat -f%m "$LOCAL_BIN" >/dev/null 2>&1; then
  MTIME="$(stat -f%m "$LOCAL_BIN")"
else
  MTIME="$(stat -c%Y "$LOCAL_BIN")"
fi
# Gồm PARTS để đổi số part không reuse nhầm split cũ
FINGERPRINT="${TOTAL_SIZE}-${MTIME}-p${PARTS}"
WORK_DIR="${TMPDIR:-/tmp}/upload-main-${service}-${FINGERPRINT}"
PART_PREFIX="$WORK_DIR/part."

mkdir -p "$WORK_DIR"

EXISTING_PARTS="$(find "$WORK_DIR" -maxdepth 1 -type f -name 'part.*' 2>/dev/null | wc -l | tr -d ' ')"
if [ "$EXISTING_PARTS" -eq "$PARTS" ]; then
  echo "--- REUSE local parts in $WORK_DIR ($EXISTING_PARTS files) ---"
else
  if [ "$EXISTING_PARTS" -ne 0 ]; then
    echo "--- RESPLIT (had $EXISTING_PARTS parts, want $PARTS) ---"
    rm -f ${PART_PREFIX}*
  else
    echo "--- SPLIT ($PARTS parts, ${TOTAL_SIZE} bytes) ---"
  fi
  PART_SIZE=$(( (TOTAL_SIZE + PARTS - 1) / PARTS ))
  split -b "$PART_SIZE" "$LOCAL_BIN" "$PART_PREFIX"
fi

PART_FILES=()
while IFS= read -r line; do
  PART_FILES+=("$line")
done <<EOF
$(ls -1 ${PART_PREFIX}* | sort)
EOF

if [ "${#PART_FILES[@]}" -eq 0 ]; then
  echo "No parts created"
  exit 1
fi

echo "--- PREPARE remote $REMOTE:$REMOTE_DIR ---"
ssh "$REMOTE" "mkdir -p '$REMOTE_DIR/.upload-parts/$FINGERPRINT'"

remote_part_size() {
  local name="$1"
  ssh "$REMOTE" "stat -c%s '$REMOTE_DIR/.upload-parts/$FINGERPRINT/$name' 2>/dev/null || echo 0"
}

upload_one() {
  local local_part="$1"
  local name want have
  name="$(basename "$local_part")"
  want="$(file_size "$local_part")"
  have="$(remote_part_size "$name")"

  if [ "$have" = "$want" ]; then
    echo "[skip] $name (${want} bytes already on remote)"
    return 0
  fi

  echo "[up]   $name (${have} -> ${want} bytes)"
  # không dùng -q: vẫn im khá lâu nếu scp buffer; progress chính nằm ở monitor bên dưới
  scp "$local_part" "$REMOTE:$REMOTE_DIR/.upload-parts/$FINGERPRINT/${name}.partial"
  ssh "$REMOTE" "mv -f '$REMOTE_DIR/.upload-parts/$FINGERPRINT/${name}.partial' '$REMOTE_DIR/.upload-parts/$FINGERPRINT/$name'"
  echo "[ok]   $name"
}

# In tiến độ remote mỗi vài giây (4 scp song song dễ làm stdout scp rối)
monitor_progress() {
  local interval="${PROGRESS_INTERVAL:-5}"
  while true; do
    sleep "$interval" || true
    ssh "$REMOTE" "bash -s" <<EOF || true
cd '$REMOTE_DIR/.upload-parts/$FINGERPRINT' 2>/dev/null || exit 0
echo "--- progress \$(date '+%H:%M:%S') ---"
shopt -s nullglob
for f in part.*; do
  sz=\$(stat -c%s "\$f" 2>/dev/null || echo 0)
  printf '  %-24s %10s bytes\n' "\$f" "\$sz"
done
EOF
  done
}

echo "--- UPLOAD (parallel=$PARALLEL, resume=on) ---"
echo "(im lâu là bình thường; mỗi ~${PROGRESS_INTERVAL:-5}s sẽ in progress remote)"
monitor_progress &
MONITOR_PID=$!
trap 'kill "$MONITOR_PID" 2>/dev/null || true' EXIT

UPLOAD_PIDS=()
fail=0
active=0
for local_part in "${PART_FILES[@]}"; do
  (
    upload_one "$local_part"
  ) &
  UPLOAD_PIDS+=($!)
  active=$((active + 1))
  if [ "$active" -ge "$PARALLEL" ]; then
    # Chỉ chờ 1 upload job, không wait monitor
    pid="${UPLOAD_PIDS[0]}"
    UPLOAD_PIDS=("${UPLOAD_PIDS[@]:1}")
    if ! wait "$pid"; then
      fail=1
    fi
    active=$((active - 1))
  fi
done
for pid in "${UPLOAD_PIDS[@]}"; do
  if ! wait "$pid"; then
    fail=1
  fi
done

kill "$MONITOR_PID" 2>/dev/null || true
wait "$MONITOR_PID" 2>/dev/null || true
trap - EXIT
if [ "$fail" -ne 0 ]; then
  echo "Upload failed — chạy lại script để resume các part còn thiếu"
  exit 1
fi
echo "--- UPLOAD done ---"

echo "--- VERIFY remote parts ---"
for local_part in "${PART_FILES[@]}"; do
  name="$(basename "$local_part")"
  want="$(file_size "$local_part")"
  have="$(remote_part_size "$name")"
  if [ "$have" != "$want" ]; then
    echo "Missing/incomplete part: $name (remote=$have want=$want)"
    echo "Chạy lại script để resume."
    exit 1
  fi
done

echo "--- ASSEMBLE remote main ---"
ssh "$REMOTE" "bash -s" <<EOF
set -euo pipefail
cd '$REMOTE_DIR/.upload-parts/$FINGERPRINT'
cat \$(ls -1 part.* | grep -v '\\.partial\$' | sort) > '$REMOTE_DIR/main.new'
chmod +x '$REMOTE_DIR/main.new'
remote_size=\$(stat -c%s '$REMOTE_DIR/main.new')
if [ "\$remote_size" != "$TOTAL_SIZE" ]; then
  echo "Size mismatch after assemble: \$remote_size != $TOTAL_SIZE"
  exit 1
fi
mv -f '$REMOTE_DIR/main.new' '$REMOTE_DIR/main'
find '$REMOTE_DIR/.upload-parts' -mindepth 1 -maxdepth 1 ! -name '$FINGERPRINT' -exec rm -rf {} +
ls -lh '$REMOTE_DIR/main'
EOF

echo "Done: $REMOTE:$REMOTE_DIR/main ($TOTAL_SIZE bytes)"
echo "Local parts: $WORK_DIR"
