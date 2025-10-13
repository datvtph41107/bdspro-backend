#!/bin/bash

# Script để check các dependencies trước khi deploy
# Usage: ./pre-deploy-check.sh [service-name]

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SERVICE_TO_CHECK=$1

echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║              PRE-DEPLOYMENT VALIDATION CHECK               ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

cd ../..

ALL_OK=true

# Check if service specified
if [ -n "$SERVICE_TO_CHECK" ]; then
  echo -e "${BLUE}Checking service: ${SERVICE_TO_CHECK}${NC}"
  echo ""
fi

# Function to check file exists
check_file() {
  local file=$1
  local desc=$2
  
  if [ -f "$file" ]; then
    echo -e "${GREEN}✓${NC} $desc"
    return 0
  else
    echo -e "${RED}✗${NC} $desc"
    echo -e "  ${RED}Missing: $file${NC}"
    ALL_OK=false
    return 1
  fi
}

# Function to check directory exists
check_dir() {
  local dir=$1
  local desc=$2
  
  if [ -d "$dir" ]; then
    echo -e "${GREEN}✓${NC} $desc"
    return 0
  else
    echo -e "${RED}✗${NC} $desc"
    echo -e "  ${RED}Missing: $dir${NC}"
    ALL_OK=false
    return 1
  fi
}

# Check common files
echo -e "${YELLOW}=== Checking common files ===${NC}"
check_file "config/production/auth.yml" "Auth production config"
check_file "config/production/user.yml" "User production config"
check_file "config/production/gateway.yml" "Gateway production config"
check_file "config/production/bdspro.yml" "BDS Pro production config"
check_file "config/production/notification.yml" "Notification production config"
check_file "config/production/bdspro-5e599-firebase-adminsdk-fbsvc-5aff064187.json" "Firebase credentials for notification"
echo ""

# Check service-specific files based on argument
if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "notification" ]; then
  echo -e "${YELLOW}=== Checking notification service ===${NC}"
  check_file "notification-service/config/product.yml" "Notification product.yml"
  check_file "notification-service/Dockerfile" "Notification Dockerfile"
  check_dir "notification-service/config" "Notification config directory"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "chat" ]; then
  echo -e "${YELLOW}=== Checking chat service ===${NC}"
  check_dir "chat-service/configs" "Chat configs directory"
  check_file "chat-service/Dockerfile" "Chat Dockerfile"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "gateway" ]; then
  echo -e "${YELLOW}=== Checking gateway service ===${NC}"
  check_dir "gateway-service/docs" "Gateway docs directory"
  check_file "gateway-service/config/product.yml" "Gateway product.yml"
  check_file "gateway-service/Dockerfile" "Gateway Dockerfile"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "file" ]; then
  echo -e "${YELLOW}=== Checking file service ===${NC}"
  check_dir "file-service/docs" "File service docs directory"
  check_file "file-service/config/product.yml" "File service product.yml"
  check_file "file-service/Dockerfile" "File service Dockerfile"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "relay" ]; then
  echo -e "${YELLOW}=== Checking relay service ===${NC}"
  check_dir "relay-service/config" "Relay config directory"
  check_file "relay-service/Dockerfile" "Relay Dockerfile"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "appointment" ]; then
  echo -e "${YELLOW}=== Checking appointment service ===${NC}"
  check_dir "appointment-service/config" "Appointment config directory"
  check_file "appointment-service/config/product.yml" "Appointment product.yml"
  check_file "appointment-service/Dockerfile" "Appointment Dockerfile"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "org" ] || [ "$SERVICE_TO_CHECK" = "organization" ]; then
  echo -e "${YELLOW}=== Checking organization service ===${NC}"
  check_dir "organization-service/config" "Organization config directory"
  check_file "organization-service/Dockerfile" "Organization Dockerfile"
  echo ""
fi

if [ -z "$SERVICE_TO_CHECK" ] || [ "$SERVICE_TO_CHECK" = "payment" ]; then
  echo -e "${YELLOW}=== Checking payment service ===${NC}"
  check_dir "payment-service/config" "Payment config directory"
  check_file "payment-service/Dockerfile" "Payment Dockerfile"
  echo ""
fi

# Check all services have go.mod
if [ -z "$SERVICE_TO_CHECK" ]; then
  echo -e "${YELLOW}=== Checking go.mod files ===${NC}"
  for service_dir in *-service; do
    if [ -d "$service_dir" ]; then
      service_name=$(echo $service_dir | sed 's/-service//')
      if [ -f "$service_dir/go.mod" ]; then
        echo -e "${GREEN}✓${NC} $service_name has go.mod"
      else
        echo -e "${RED}✗${NC} $service_name missing go.mod"
        ALL_OK=false
      fi
    fi
  done
  echo ""
fi

# Check all Dockerfiles have config copy
if [ -z "$SERVICE_TO_CHECK" ]; then
  echo -e "${YELLOW}=== Checking Dockerfiles for config copy ===${NC}"
  for service_dir in *-service; do
    if [ -f "$service_dir/Dockerfile" ]; then
      service_name=$(echo $service_dir | sed 's/-service//')
      if grep -q "COPY config config" "$service_dir/Dockerfile" || grep -q "COPY \. \." "$service_dir/Dockerfile"; then
        echo -e "${GREEN}✓${NC} $service_name Dockerfile includes config"
      else
        echo -e "${RED}✗${NC} $service_name Dockerfile missing config copy"
        ALL_OK=false
      fi
    fi
  done
  echo ""
fi

# Check shared modules
echo -e "${YELLOW}=== Checking shared modules ===${NC}"
check_file "shared/common/go.mod" "shared/common go.mod"
check_file "shared/protobuf/go.mod" "shared/protobuf go.mod"
check_file "shared/base/go.mod" "shared/base go.mod"
echo ""

# Summary
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║                    VALIDATION SUMMARY                      ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"

if [ "$ALL_OK" = true ]; then
  echo -e "${GREEN}✓ All checks passed! Ready to deploy.${NC}"
  exit 0
else
  echo -e "${RED}✗ Some checks failed. Please fix the issues above before deploying.${NC}"
  echo ""
  echo -e "${YELLOW}Common fixes:${NC}"
  echo -e "  1. Run: ${BLUE}./tidy-all.sh${NC} to ensure all go.mod files are present"
  echo -e "  2. Check if all config/production/*.yml files exist"
  echo -e "  3. Ensure Firebase credentials are in config/production/"
  exit 1
fi

