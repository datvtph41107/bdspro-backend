#!/bin/bash

# Script để stop tất cả services trên production server
# Usage: ./stop-all.sh [--remove]
# Options:
#   --remove    : Remove containers sau khi stop

REMOTE_USER="root"
REMOTE_HOST="103.172.239.89"
WORKSPACE_NAME=golang-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

REMOVE_CONTAINERS=false

# Parse arguments
if [[ "$1" == "--remove" ]]; then
  REMOVE_CONTAINERS=true
fi

echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║         STOPPING ALL SERVICES ON PRODUCTION                ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

if [ "$REMOVE_CONTAINERS" = true ]; then
  echo -e "${RED}⚠️  WARNING: This will STOP and REMOVE all containers${NC}"
else
  echo -e "${YELLOW}This will STOP all services (containers will remain)${NC}"
fi

echo ""
read -p "Are you sure? (yes/no): " -r
echo ""

if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
  echo "Cancelled."
  exit 0
fi

echo "Connecting to server..."

if [ "$REMOVE_CONTAINERS" = true ]; then
  ssh ${REMOTE_USER}@${REMOTE_HOST} << EOF
cd ${REMOTE_PATH}
echo "Stopping and removing all containers..."
docker compose down
echo "Cleaning up unused images..."
docker image prune -f
EOF
else
  ssh ${REMOTE_USER}@${REMOTE_HOST} << EOF
cd ${REMOTE_PATH}
echo "Stopping all containers..."
docker compose stop
EOF
fi

echo ""
echo -e "${GREEN}✓ All services stopped successfully${NC}"

if [ "$REMOVE_CONTAINERS" = true ]; then
  echo -e "${YELLOW}Note: Containers have been removed. Use deploy scripts to recreate them.${NC}"
fi

