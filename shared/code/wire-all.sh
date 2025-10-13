#!/bin/bash

# Script chạy wire generation cho tất cả microservices
# Usage: ./wire-all.sh [options]
# Options:
#   --services="..."  : Chỉ generate wire cho các service cụ thể (cách nhau bởi dấu phẩy)

start=$(date +%s)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Cấu hình
SELECTED_SERVICES=""

# Parse arguments
for arg in "$@"; do
  case $arg in
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

# Nếu có selected services, chỉ wire những cái đó
if [ -n "$SELECTED_SERVICES" ]; then
  IFS=',' read -ra SERVICES_TO_WIRE <<< "$SELECTED_SERVICES"
  echo -e "${YELLOW}=== Running wire generation for selected services: ${SELECTED_SERVICES} ===${NC}"
else
  SERVICES_TO_WIRE=("${ALL_SERVICES[@]}")
  echo -e "${YELLOW}=== Running wire generation for all services ===${NC}"
fi

# Counters
TOTAL_SERVICES=${#SERVICES_TO_WIRE[@]}
SUCCESS_COUNT=0
FAILED_COUNT=0
SKIPPED_COUNT=0
FAILED_SERVICES=()
SKIPPED_SERVICES=()

# Function to generate wire for a single service
wire_service() {
  service=$1
  index=$2
  total=$3
  service_dir=$(get_service_dir "$service")
  
  if [ -z "$service_dir" ]; then
    echo -e "${RED}[$index/$total] Unknown service: ${service}${NC}"
    return 1
  fi
  
  echo -e "${BLUE}[$index/$total] Generating wire in ${service_dir}...${NC}"
  
  if [ ! -d "../../${service_dir}" ]; then
    echo -e "${RED}  ✗ Directory not found: ${service_dir}${NC}"
    return 1
  fi
  
  if [ ! -d "../../${service_dir}/wire" ]; then
    echo -e "${YELLOW}  ⚠ No wire/ directory found, skipping${NC}"
    return 2
  fi
  
  # Run wire generation
  if (cd "../../${service_dir}" && go run github.com/google/wire/cmd/wire@latest ./wire 2>&1); then
    echo -e "${GREEN}  ✓ ${service} wire generated successfully${NC}"
    return 0
  else
    echo -e "${RED}  ✗ ${service} wire generation failed${NC}"
    return 1
  fi
}

# Main logic
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║         WIRE GENERATION - ALL MICROSERVICES                ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Total services: ${TOTAL_SERVICES}"
echo ""

# Generate wire for services
index=1
for service in "${SERVICES_TO_WIRE[@]}"; do
  result=$(wire_service "$service" "$index" "$TOTAL_SERVICES")
  exit_code=$?
  
  echo "$result"
  
  if [ $exit_code -eq 0 ]; then
    ((SUCCESS_COUNT++))
  elif [ $exit_code -eq 2 ]; then
    ((SKIPPED_COUNT++))
    SKIPPED_SERVICES+=("$service")
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
echo -e "${YELLOW}║                   WIRE GENERATION SUMMARY                  ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo -e "Total services: ${TOTAL_SERVICES}"
echo -e "${GREEN}Successful: ${SUCCESS_COUNT}${NC}"
echo -e "${YELLOW}Skipped: ${SKIPPED_COUNT}${NC}"
echo -e "${RED}Failed: ${FAILED_COUNT}${NC}"
echo -e "Duration: ${minutes}m ${seconds}s"
echo ""

if [ ${SKIPPED_COUNT} -gt 0 ]; then
  echo -e "${YELLOW}Skipped services (no wire/ directory):${NC}"
  for service in "${SKIPPED_SERVICES[@]}"; do
    echo -e "  ${YELLOW}⚠ ${service}${NC}"
  done
  echo ""
fi

if [ ${FAILED_COUNT} -gt 0 ]; then
  echo -e "${RED}Failed services:${NC}"
  for service in "${FAILED_SERVICES[@]}"; do
    echo -e "  ${RED}✗ ${service}${NC}"
  done
  echo ""
  exit 1
else
  echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║      ALL WIRE GENERATIONS COMPLETED SUCCESSFULLY! 🎉       ║${NC}"
  echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
  exit 0
fi
