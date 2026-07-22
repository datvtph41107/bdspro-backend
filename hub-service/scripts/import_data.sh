#!/bin/bash

# Script to import location data into hub-service database

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Hub Service - Import Location Data${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "scripts/import_location_data.go" ]; then
    echo -e "${RED}Error: Please run this script from hub-service directory${NC}"
    exit 1
fi

# Load environment variables if .env file exists
if [ -f ".env" ]; then
    echo -e "${YELLOW}Loading environment variables from .env file...${NC}"
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo -e "${YELLOW}Warning: DATABASE_URL not set. Using default connection string.${NC}"
    echo -e "${YELLOW}Set DATABASE_URL in .env file or environment variable to override.${NC}"
    echo ""
fi

# Run the import script
echo -e "${GREEN}Running import script...${NC}"
cd scripts
go run import_location_data.go

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Import completed!${NC}"
echo -e "${GREEN}========================================${NC}"

