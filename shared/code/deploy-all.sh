#!/bin/bash

# Script deploy toàn bộ microservices lên production
# Usage: ./deploy-all.sh [options]
# Options:
#   --skip-shared     : Bỏ qua việc build shared modules
#   --parallel        : Build song song (nhanh hơn nhưng khó debug)
#   --services="..."  : Deploy chỉ các service cụ thể (cách nhau bởi dấu phẩy)

start=$(date +%s)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Cấu hình
SKIP_SHARED=false
PARALLEL=false
SELECTED_SERVICES=""

# Parse arguments
for arg in "$@"; do
  case $arg in
    --skip-shared)
      SKIP_SHARED=true
      shift
      ;;
    --parallel)
      PARALLEL=true
      shift
      ;;
    --services=*)
      SELECTED_SERVICES="${arg#*=}"
      shift
      ;;
    *)
      ;;
  esac
done

# Danh sách services cần deploy theo thứ tự
# Core services trước, business services sau
CORE_SERVICES=(
  "auth"           # Auth service - phải deploy trước
  "user"           # User service
  "user-grpc"      # User gRPC service
  "gateway"        # API Gateway - deploy sau khi các service khác đã sẵn sàng
)

BUSINESS_SERVICES=(
  "bdspro"         # BDS Pro service
  "crm"            # CRM service
  "notification"   # Notification service
  "social"         # Social service
  "chat"           # Chat service
  "chathttp"       # Chat HTTP service
  "file"           # File service
  "map"            # Map service
  "payment"        # Payment service
  "membership"     # Membership service
  "relay"          # Relay service
  "appointment"    # Appointment service
  "org"            # Organization service
  "transaction"    # Transaction service
  "tqd"            # TQD service
  "marketing"      # Marketing service
  "feedback"       # Feedback service
)

# Nếu có selected services, chỉ deploy những cái đó
if [ -n "$SELECTED_SERVICES" ]; then
  IFS=',' read -ra SERVICES_TO_DEPLOY <<< "$SELECTED_SERVICES"
  echo -e "${YELLOW}=== Deploying selected services: ${SELECTED_SERVICES} ===${NC}"
else
  SERVICES_TO_DEPLOY=("${CORE_SERVICES[@]}" "${BUSINESS_SERVICES[@]}")
  echo -e "${YELLOW}=== Deploying all services ===${NC}"
fi

# Counters
TOTAL_SERVICES=${#SERVICES_TO_DEPLOY[@]}
SUCCESS_COUNT=0
FAILED_COUNT=0
FAILED_SERVICES=()

# Function to deploy a single service
deploy_service() {
  local service=$1
  local index=$2
  local total=$3
  
  echo -e "\n${BLUE}[$index/$total] Deploying ${service}...${NC}"
  
  if bash ./release.sh "$service"; then
    echo -e "${GREEN}✓ ${service} deployed successfully${NC}"
    return 0
  else
    echo -e "${RED}✗ ${service} deployment failed${NC}"
    return 1
  fi
}

# Function to deploy services sequentially
deploy_sequential() {
  local index=1
  
  for service in "${SERVICES_TO_DEPLOY[@]}"; do
    if deploy_service "$service" "$index" "$TOTAL_SERVICES"; then
      ((SUCCESS_COUNT++))
    else
      ((FAILED_COUNT++))
      FAILED_SERVICES+=("$service")
    fi
    ((index++))
  done
}

# Function to deploy services in parallel
deploy_parallel() {
  local pids=()
  local index=1
  
  # Deploy core services first (sequential)
  echo -e "${YELLOW}=== Phase 1: Deploying core services (sequential) ===${NC}"
  for service in "${CORE_SERVICES[@]}"; do
    if [[ " ${SERVICES_TO_DEPLOY[@]} " =~ " ${service} " ]]; then
      if deploy_service "$service" "$index" "$TOTAL_SERVICES"; then
        ((SUCCESS_COUNT++))
      else
        ((FAILED_COUNT++))
        FAILED_SERVICES+=("$service")
      fi
      ((index++))
    fi
  done
  
  # Wait a bit for core services to be ready
  echo -e "${YELLOW}Waiting 10 seconds for core services to stabilize...${NC}"
  sleep 10
  
  # Deploy business services in parallel
  echo -e "${YELLOW}=== Phase 2: Deploying business services (parallel) ===${NC}"
  for service in "${BUSINESS_SERVICES[@]}"; do
    if [[ " ${SERVICES_TO_DEPLOY[@]} " =~ " ${service} " ]]; then
      echo -e "${BLUE}[$index/$TOTAL_SERVICES] Starting deployment of ${service}...${NC}"
      bash ./release.sh "$service" > "logs/${service}_deploy.log" 2>&1 &
      pids+=($!)
      ((index++))
    fi
  done
  
  # Wait for all parallel deployments
  for i in "${!pids[@]}"; do
    local pid=${pids[$i]}
    local service_index=$((${#CORE_SERVICES[@]} + i + 1))
    
    if wait $pid; then
      ((SUCCESS_COUNT++))
      echo -e "${GREEN}✓ Service #$service_index deployed successfully${NC}"
    else
      ((FAILED_COUNT++))
      echo -e "${RED}✗ Service #$service_index deployment failed${NC}"
    fi
  done
}

# Main deployment logic
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║     MICROSERVICES DEPLOYMENT - PRODUCTION                  ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Total services to deploy: ${TOTAL_SERVICES}"
echo -e "Parallel mode: ${PARALLEL}"
echo -e "Skip shared: ${SKIP_SHARED}"
echo ""

# Pre-deployment validation
echo -e "${YELLOW}=== Running pre-deployment checks ===${NC}"
if [ -f "./pre-deploy-check.sh" ]; then
  if ./pre-deploy-check.sh; then
    echo -e "${GREEN}✓ Pre-deployment validation passed${NC}"
  else
    echo -e "${RED}✗ Pre-deployment validation failed${NC}"
    echo -e "${YELLOW}Fix the issues above before deploying${NC}"
    exit 1
  fi
else
  echo -e "${YELLOW}⚠ pre-deploy-check.sh not found, skipping validation${NC}"
fi
echo ""

# Create logs directory if needed
mkdir -p logs

# Build shared modules first (if not skipped)
if [ "$SKIP_SHARED" = false ]; then
  echo -e "${YELLOW}=== Building shared modules ===${NC}"
  cd ../..
  
  echo "Building shared/common..."
  cd shared/common && go mod tidy && cd ../..
  
  echo "Building shared/code..."
  cd shared/code && go mod tidy && cd ../..
  
  echo "Building shared/protobuf..."
  cd shared/protobuf && go mod tidy && cd ../..
  
  echo -e "${GREEN}✓ Shared modules built successfully${NC}\n"
  cd shared/code
fi

# Deploy services
if [ "$PARALLEL" = true ]; then
  deploy_parallel
else
  deploy_sequential
fi

# Summary
end=$(date +%s)
elapsed=$((end - start))
minutes=$((elapsed / 60))
seconds=$((elapsed % 60))

echo ""
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║                    DEPLOYMENT SUMMARY                      ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo -e "Total services: ${TOTAL_SERVICES}"
echo -e "${GREEN}Successful: ${SUCCESS_COUNT}${NC}"
echo -e "${RED}Failed: ${FAILED_COUNT}${NC}"
echo -e "Duration: ${minutes}m ${seconds}s"
echo ""

if [ ${FAILED_COUNT} -gt 0 ]; then
  echo -e "${RED}Failed services:${NC}"
  for service in "${FAILED_SERVICES[@]}"; do
    echo -e "  ${RED}✗ ${service}${NC}"
  done
  echo ""
  echo -e "${YELLOW}Check logs in ./logs/ directory for details${NC}"
  exit 1
else
  echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║         ALL SERVICES DEPLOYED SUCCESSFULLY! 🎉            ║${NC}"
  echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
  exit 0
fi

