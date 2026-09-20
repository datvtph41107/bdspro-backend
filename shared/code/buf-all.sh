#!/bin/bash

# Script chạy buf generate cho tất cả protobuf schemas
# Usage: ./buf-all.sh [options]
# Options:
#   --modules="..."   : Chỉ generate các module cụ thể (cách nhau bởi dấu phẩy)
#   --clean           : Xóa các file generated cũ trước khi generate

start=$(date +%s)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Cấu hình
SELECTED_MODULES=""
CLEAN_FIRST=false

# Parse arguments
for arg in "$@"; do
  case $arg in
    --modules=*)
      SELECTED_MODULES="${arg#*=}"
      shift
      ;;
    --clean)
      CLEAN_FIRST=true
      shift
      ;;
    *)
      ;;
  esac
done

# Directories
PROTOBUF_DIR="../../shared/protobuf"
SCHEMA_DIR="${PROTOBUF_DIR}/schema"
OUTPUT_DIR="${PROTOBUF_DIR}/types"
DOCS_DIR="${PROTOBUF_DIR}/docs"

# All proto modules
ALL_MODULES=(
  "appointment"
  "assistant"
  "auth"
  "bdspro"
  "chat"
  "crm"
  "feedback"
  "generic"
  "notification"
  "operation"
  "organization"
  "payment"
  "shared"
  "social"
  "spec"
  "task"
  "tqd"
  "transaction"
  "user"
)

# Nếu có selected modules, chỉ generate những cái đó
if [ -n "$SELECTED_MODULES" ]; then
  IFS=',' read -ra MODULES_TO_GENERATE <<< "$SELECTED_MODULES"
  echo -e "${YELLOW}=== Running buf generate for selected modules: ${SELECTED_MODULES} ===${NC}"
else
  MODULES_TO_GENERATE=("${ALL_MODULES[@]}")
  echo -e "${YELLOW}=== Running buf generate for all modules ===${NC}"
fi

# Counters
TOTAL_MODULES=${#MODULES_TO_GENERATE[@]}
SUCCESS_COUNT=0
FAILED_COUNT=0
SKIPPED_COUNT=0
FAILED_MODULES=()
SKIPPED_MODULES=()

# Function to clean old generated files for a module
clean_module() {
  module=$1
  
  echo -e "${BLUE}  Cleaning old generated files for ${module}...${NC}"
  
  # Clean pb.go files
  if [ -d "${OUTPUT_DIR}/${module}" ]; then
    find "${OUTPUT_DIR}/${module}" -name "*.pb.go" -type f -delete 2>/dev/null || true
    find "${OUTPUT_DIR}/${module}" -name "*.pb.gw.go" -type f -delete 2>/dev/null || true
    echo -e "${GREEN}    ✓ Cleaned ${OUTPUT_DIR}/${module}${NC}"
  fi
  
  # Clean swagger files
  if [ -d "${DOCS_DIR}/${module}" ]; then
    find "${DOCS_DIR}/${module}" -name "*.swagger.json" -type f -delete 2>/dev/null || true
    echo -e "${GREEN}    ✓ Cleaned ${DOCS_DIR}/${module}${NC}"
  fi
}

# Function to generate buf for a single module
generate_module() {
  module=$1
  index=$2
  total=$3
  
  echo -e "${BLUE}[$index/$total] Generating protobuf for ${module}...${NC}"
  
  # Check if module directory exists
  if [ ! -d "${SCHEMA_DIR}/${module}" ]; then
    echo -e "${RED}  ✗ Module directory not found: ${module}${NC}"
    return 1
  fi
  
  # Check if module has proto files
  proto_count=$(find "${SCHEMA_DIR}/${module}" -name "*.proto" -type f | wc -l | tr -d ' ')
  if [ "$proto_count" -eq 0 ]; then
    echo -e "${YELLOW}  ⚠ No proto files found in ${module}, skipping${NC}"
    return 2
  fi
  
  echo -e "${BLUE}  Found ${proto_count} proto file(s)${NC}"
  
  # Clean old files if requested
  if [ "$CLEAN_FIRST" = true ]; then
    clean_module "$module"
  fi
  
  # Create output directories
  mkdir -p "${OUTPUT_DIR}/${module}"
  mkdir -p "${DOCS_DIR}/${module}"
  
  # Run buf generate for specific module using protoc
  echo -e "${BLUE}  Running protoc...${NC}"
  if (cd "${SCHEMA_DIR}" && protoc \
      --proto_path=. \
      --proto_path=../../ \
      --go_out=../types \
      --go_opt=paths=source_relative \
      --go-grpc_out=../types \
      --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out=../types \
      --grpc-gateway_opt=paths=source_relative \
      --grpc-gateway_opt=generate_unbound_methods=true \
      --openapiv2_out=../docs \
      --openapiv2_opt=logtostderr=true \
      ${module}/*.proto 2>&1); then
    
    # Format generated code
    if [ -d "${OUTPUT_DIR}/${module}" ]; then
      gofmt -w "${OUTPUT_DIR}/${module}/" 2>/dev/null || true
    fi
    
    echo -e "${GREEN}  ✓ ${module} generated successfully${NC}"
    return 0
  else
    echo -e "${RED}  ✗ ${module} generation failed${NC}"
    return 1
  fi
}

# Function to generate all at once using buf
generate_all_with_buf() {
  echo -e "${YELLOW}=== Running buf generate for entire workspace ===${NC}"
  echo ""
  
  if [ "$CLEAN_FIRST" = true ]; then
    echo -e "${BLUE}Cleaning all old generated files...${NC}"
    find "${OUTPUT_DIR}" -name "*.pb.go" -type f -delete 2>/dev/null || true
    find "${OUTPUT_DIR}" -name "*.pb.gw.go" -type f -delete 2>/dev/null || true
    find "${DOCS_DIR}" -name "*.swagger.json" -type f -delete 2>/dev/null || true
    echo -e "${GREEN}✓ Cleaned all generated files${NC}"
    echo ""
  fi
  
  echo -e "${BLUE}Running buf generate...${NC}"
  if (cd "${SCHEMA_DIR}" && buf generate 2>&1); then
    echo -e "${GREEN}✓ Buf generate completed successfully${NC}"
    
    # Format all generated code
    echo -e "${BLUE}Formatting generated code...${NC}"
    find "${OUTPUT_DIR}" -name "*.pb.go" -type f -exec gofmt -w {} \; 2>/dev/null || true
    echo -e "${GREEN}✓ Code formatted${NC}"
    
    return 0
  else
    echo -e "${RED}✗ Buf generate failed${NC}"
    return 1
  fi
}

# Main logic
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║           BUF GENERATE - ALL PROTOBUF SCHEMAS              ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Schema directory: ${SCHEMA_DIR}"
echo -e "Output directory: ${OUTPUT_DIR}"
echo -e "Docs directory: ${DOCS_DIR}"
echo -e "Clean first: ${CLEAN_FIRST}"
echo ""

# Check if buf is available
if command -v buf &> /dev/null; then
  BUF_VERSION=$(buf --version 2>/dev/null || echo "unknown")
  echo -e "${GREEN}✓ buf is installed: ${BUF_VERSION}${NC}"
  USE_BUF=true
else
  echo -e "${YELLOW}⚠ buf not found, will use protoc directly${NC}"
  USE_BUF=false
fi

# Check if protoc is available
if command -v protoc &> /dev/null; then
  PROTOC_VERSION=$(protoc --version 2>/dev/null || echo "unknown")
  echo -e "${GREEN}✓ protoc is installed: ${PROTOC_VERSION}${NC}"
else
  echo -e "${RED}✗ protoc not found, please install it first${NC}"
  exit 1
fi

echo ""

# If no specific modules selected and buf is available, use buf generate
if [ -z "$SELECTED_MODULES" ] && [ "$USE_BUF" = true ]; then
  if generate_all_with_buf; then
    SUCCESS_COUNT=${#ALL_MODULES[@]}
  else
    FAILED_COUNT=${#ALL_MODULES[@]}
  fi
else
  # Generate per module
  echo -e "${YELLOW}Total modules: ${TOTAL_MODULES}${NC}"
  echo ""
  
  index=1
  for module in "${MODULES_TO_GENERATE[@]}"; do
    result=$(generate_module "$module" "$index" "$TOTAL_MODULES")
    exit_code=$?
    
    echo "$result"
    
    if [ $exit_code -eq 0 ]; then
      ((SUCCESS_COUNT++))
    elif [ $exit_code -eq 2 ]; then
      ((SKIPPED_COUNT++))
      SKIPPED_MODULES+=("$module")
    else
      ((FAILED_COUNT++))
      FAILED_MODULES+=("$module")
    fi
    
    ((index++))
    echo ""
  done
fi

# Summary
end=$(date +%s)
elapsed=$((end - start))
minutes=$((elapsed / 60))
seconds=$((elapsed % 60))

echo ""
echo -e "${YELLOW}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${YELLOW}║                  BUF GENERATION SUMMARY                    ║${NC}"
echo -e "${YELLOW}╚════════════════════════════════════════════════════════════╝${NC}"

if [ -z "$SELECTED_MODULES" ] && [ "$USE_BUF" = true ]; then
  echo -e "Mode: Full workspace (buf generate)"
  if [ $SUCCESS_COUNT -gt 0 ]; then
    echo -e "${GREEN}Status: Success${NC}"
  else
    echo -e "${RED}Status: Failed${NC}"
  fi
else
  echo -e "Mode: Per module (protoc)"
  echo -e "Total modules: ${TOTAL_MODULES}"
  echo -e "${GREEN}Successful: ${SUCCESS_COUNT}${NC}"
  echo -e "${YELLOW}Skipped: ${SKIPPED_COUNT}${NC}"
  echo -e "${RED}Failed: ${FAILED_COUNT}${NC}"
fi

echo -e "Duration: ${minutes}m ${seconds}s"
echo ""

if [ ${SKIPPED_COUNT} -gt 0 ]; then
  echo -e "${YELLOW}Skipped modules (no proto files):${NC}"
  for module in "${SKIPPED_MODULES[@]}"; do
    echo -e "  ${YELLOW}⚠ ${module}${NC}"
  done
  echo ""
fi

if [ ${FAILED_COUNT} -gt 0 ]; then
  echo -e "${RED}Failed modules:${NC}"
  for module in "${FAILED_MODULES[@]}"; do
    echo -e "  ${RED}✗ ${module}${NC}"
  done
  echo ""
  exit 1
else
  echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║      ALL PROTOBUF GENERATIONS COMPLETED! 🎉                ║${NC}"
  echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
  echo ""
  echo -e "${YELLOW}Next steps:${NC}"
  echo "  1. Review generated files in: ${OUTPUT_DIR}"
  echo "  2. Review swagger docs in: ${DOCS_DIR}"
  echo "  3. Run go mod tidy: ./tidy-all.sh"
  echo "  4. Run wire generate: ./wire-all.sh"
  echo ""
  exit 0
fi

