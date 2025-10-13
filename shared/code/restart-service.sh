#!/bin/bash

# Script để restart một hoặc nhiều services trên production
# Usage: ./restart-service.sh <service-name> [service-name2 ...]
# Example: ./restart-service.sh auth user gateway

if [ $# -eq 0 ]; then
  echo "Usage: $0 <service-name> [service-name2 ...]"
  echo ""
  echo "Available services:"
  echo "  auth, user, user-grpc, gateway"
  echo "  bdspro, bdspro-grpc, crm, notification, social"
  echo "  chat, chathttp, file, map, payment, membership"
  echo "  relay, appointment, org, transaction"
  exit 1
fi

REMOTE_USER="root"
REMOTE_HOST="103.172.239.89"
WORKSPACE_NAME=golang-microservice
REMOTE_PATH=/root/$WORKSPACE_NAME

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║              RESTARTING SERVICES ON SERVER                 ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

SERVICES_TO_RESTART=("$@")

echo "Services to restart: ${SERVICES_TO_RESTART[*]}"
echo ""

# Map service names to docker compose service names
declare -A SERVICE_MAP
SERVICE_MAP["user-grpc"]="user-grpc"
SERVICE_MAP["bdspro-grpc"]="bdspro-grpc"
SERVICE_MAP["chathttp"]="http-server"
SERVICE_MAP["org"]="organization"
SERVICE_MAP["swagger-chat"]="swagger-chat"

for service in "${SERVICES_TO_RESTART[@]}"; do
  # Get docker compose service name
  compose_service="${SERVICE_MAP[$service]:-$service}"
  
  echo -e "${BLUE}Restarting ${service} (${compose_service})...${NC}"
  
  ssh ${REMOTE_USER}@${REMOTE_HOST} << EOF
cd ${REMOTE_PATH}
docker compose restart ${compose_service}
EOF
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ ${service} restarted successfully${NC}"
  else
    echo -e "${RED}✗ Failed to restart ${service}${NC}"
  fi
  echo ""
done

echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    RESTART COMPLETED                       ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"

