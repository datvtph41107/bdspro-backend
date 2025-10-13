#!/bin/bash

# Script để check status của các services trên production server
# Usage: ./check-services.sh

REMOTE_USER="root"
REMOTE_HOST="103.172.239.89"
WORKSPACE_NAME=golang-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║           CHECKING SERVICES STATUS ON SERVER               ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# List of all services
SERVICES=(
  "auth"
  "user"
  "user-grpc"
  "gateway"
  "bdspro"
  "bdspro-grpc"
  "crm"
  "notification"
  "social"
  "chat"
  "file"
  "map"
  "payment"
  "membership"
  "relay"
  "appointment"
  "organization"
  "transaction"
)

echo "Connecting to server ${REMOTE_HOST}..."
echo ""

# Check docker containers
ssh ${REMOTE_USER}@${REMOTE_HOST} << 'EOF'
cd /root/golang-microservice

echo "=== Docker Containers Status ==="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(NAMES|golang-microservice)"

echo ""
echo "=== Resource Usage ==="
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | grep -E "(NAME|golang-microservice)" | head -20

echo ""
echo "=== Disk Usage ==="
df -h | grep -E "(Filesystem|/dev/)"

echo ""
echo "=== Memory Usage ==="
free -h
EOF

echo ""
echo -e "${GREEN}Check completed!${NC}"

