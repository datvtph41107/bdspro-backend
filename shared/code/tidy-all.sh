#!/bin/bash

# Script chạy go mod tidy cho tất cả microservices
# Usage: ./tidy-all.sh [options]
# Options:
#   --services="..."  : Chỉ tidy các service cụ thể (cách nhau bởi dấu phẩy)
#   --skip-shared     : Bỏ qua shared modules

start=$(date +%s)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Cấu hình
SKIP_SHARED=false
SELECTED_SERVICES=""

# Parse arguments
for arg in "$@"; do
  case $arg in
    --skip-shared)
      SKIP_SHARED=true
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

# Function to get service directory
get_service_dir() {
  case "$1" in
    "auth") echo "auth-service" ;;
    "user") echo "user-service" ;;
    "gateway") echo "gateway-service" ;;
    "bdspro") echo "bdspro-service" ;;
    "crm") echo "crm-service" ;;
    "notification") echo "notification-service" ;;
    "social") echo "social-service" ;;
    "chat") echo "chat-service" ;;
    "file") echo "file-service" ;;
    "map") echo "map-service" ;;
    "payment") echo "payment-service" ;;
    "membership") echo "membership-service" ;;
    "relay") echo "relay-service" ;;
    "appointment") echo "appointment-service" ;;
    "org") echo "organization-service" ;;
    "transaction") echo "transaction-service" ;;
    "marketing") echo "marketing-service" ;;
    "tqd") echo "tqd-service" ;;
    "feedback") echo "feedback-service" ;;
    "assistant") echo "assistant-service" ;;
    *) echo "" ;;
  esac
}

# All services
ALL_SERVICES=(
  "auth"
  "user"
  "gateway"
  "bdspro"
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
  "org"
  "transaction"
  "marketing"
  "tqd"
  "feedback"
  "assistant"
)

# Nếu có selected services, chỉ tidy những cái đó
if [ -n "$SELECTED_SERVICES" ]; then
  IFS=',' read -ra SERVICES_TO_TIDY <<< "$SELECTED_SERVICES"
  echo -e "${YELLOW}=== Running go mod tidy for selected services: ${SELECTED_SERVICES} ===${NC}"
else
  SERVICES_TO_TIDY=("${ALL_SERVICES[@]}")
  echo -e "${YELLOW}=== Running go mod tidy for all services ===${NC}"
fi

# Counters
TOTAL_SERVICES=${#SERVICES_TO_TIDY[@]}
SUCCESS_COUNT=0
FAILED_COUNT=0
FAILED_SERVICES=()

# Function to tidy a single service
tidy_service() {
  service=$1
  index=$2
  total=$3
  service_dir=$(get_service_dir "$service")
  
  if [ -z "$service_dir" ]; then
    echo -e "${RED}[$index/$total] Unknown service: ${service}${NC}"
    return 1
  fi
  
  echo -e "${BLUE}[$index/$total] go mod tidy in ${service_dir}...${NC}"
  
  if [ ! -d "../../${service_dir}" ]; then
    echo -e "${RED}  ✗ Directory not found: ${service_dir}${NC}"
    return 1
  fi
  
  if [ ! -f "../../${service_dir}/go.mod" ]; then
    echo -e "${YELLOW}  ⚠ No go.mod found, skipping${NC}"
    return 0
  fi
  
  # Run go mod tidy
  if (cd "../../${service_dir}" && go mod tidy 2>&1); then
    echo -e "${GREEN}  ✓ ${service} tidied successfully${NC}"
    return 0
  else
    echo -e "${RED}  ✗ ${service} tidy failed${NC}"
    return 1
  fi
}

# Main logic
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║           GO MOD TIDY - ALL MICROSERVICES                  ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Total services: ${TOTAL_SERVICES}"
echo -e "Skip shared: ${SKIP_SHARED}"
echo ""

# Tidy shared modules first (if not skipped)
if [ "$SKIP_SHARED" = false ]; then
  echo -e "${YELLOW}=== Tidying shared modules ===${NC}"
  
  echo -e "${BLUE}[shared] go mod tidy in shared/common...${NC}"
  if (cd ../../shared/common && go mod tidy 2>&1); then
    echo -e "${GREEN}  ✓ shared/common tidied successfully${NC}"
  else
    echo -e "${RED}  ✗ shared/common tidy failed${NC}"
  fi
  
  echo -e "${BLUE}[shared] go mod tidy in shared/code...${NC}"
  if (cd ../../shared/code && go mod tidy 2>&1); then
    echo -e "${GREEN}  ✓ shared/code tidied successfully${NC}"
  else
    echo -e "${RED}  ✗ shared/code tidy failed${NC}"
  fi
  
  echo -e "${BLUE}[shared] go mod tidy in shared/protobuf...${NC}"
  if (cd ../../shared/protobuf && go mod tidy 2>&1); then
    echo -e "${GREEN}  ✓ shared/protobuf tidied successfully${NC}"
  else
    echo -e "${RED}  ✗ shared/protobuf tidy failed${NC}"
  fi
  
  echo -e "${BLUE}[shared] go mod tidy in shared/base...${NC}"
  if (cd ../../shared/base && go mod tidy 2>&1); then
    echo -e "${GREEN}  ✓ shared/base tidied successfully${NC}"
  else
    echo -e "${RED}  ✗ shared/base tidy failed${NC}"
  fi
  
  echo ""
fi

# Tidy services
index=1
for service in "${SERVICES_TO_TIDY[@]}"; do
  if tidy_service "$service" "$index" "$TOTAL_SERVICES"; then
    ((SUCCESS_COUNT++))
  else
    ((FAILED_COUNT++))
    FAILED_SERVICES+=("$service")
  fi
  ((index++))
  echo ""
done

# Summary
end=$(date +%s)
elapsed=$((end - start))
minutes=$((elapsed / 60))
seconds=$((elapsed % 60))

echo ""
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║                      TIDY SUMMARY                          ║${NC}"
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
  exit 1
else
  echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║         ALL SERVICES TIDIED SUCCESSFULLY! 🎉              ║${NC}"
  echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
  exit 0
fi

