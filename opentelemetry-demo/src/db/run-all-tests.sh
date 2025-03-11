#!/bin/bash

# Set color variables for better output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Running tests in the main db module..."
go test ./... -v
MAIN_RESULT=$?

echo -e "\n${GREEN}=================================${NC}"
echo -e "${GREEN}Running tests in postgres module...${NC}"
echo -e "${GREEN}=================================${NC}\n"
cd postgres && go test ./... -v
POSTGRES_RESULT=$?

# Return non-zero exit code if any test suite failed
if [ $MAIN_RESULT -ne 0 ] || [ $POSTGRES_RESULT -ne 0 ]; then
  echo -e "\n${RED}Some tests failed!${NC}"
  exit 1
else
  echo -e "\n${GREEN}All tests passed!${NC}"
  exit 0
fi 