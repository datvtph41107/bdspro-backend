#!/bin/bash

# Script to import location data (provinces and districts) into database
# Usage: ./scripts/import-location.sh [develop|product]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get environment from argument or default to develop
ENV=${1:-develop}

echo -e "${GREEN}=== Location Data Import Script ===${NC}"
echo -e "${YELLOW}Environment: ${ENV}${NC}"
echo ""

# Export environment variable
export ENV_RUNTIME=$ENV

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Navigate to project directory
cd "$PROJECT_DIR"

echo -e "${YELLOW}Building import tool...${NC}"
go build -o bin/import-location cmd/import-location/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
    echo ""
    echo -e "${YELLOW}Running import...${NC}"
    ./bin/import-location
    
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ Import completed successfully!${NC}"
        # Clean up binary
        rm -f bin/import-location
    else
        echo ""
        echo -e "${RED}✗ Import failed!${NC}"
        rm -f bin/import-location
        exit 1
    fi
else
    echo -e "${RED}✗ Build failed!${NC}"
    exit 1
fi

