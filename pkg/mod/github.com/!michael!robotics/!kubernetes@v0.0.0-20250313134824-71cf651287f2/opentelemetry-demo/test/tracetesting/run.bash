# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0
#/bin/bash

# This script set up how to run Tracetest and which test files 
# be executed

set -e

# Availalble services to test
ALL_SERVICES=("ad" "cart" "currency" "checkout" "frontend" "email" "payment" "product-catalog" "recommendation" "shipping")

## Script variables
# Will contain the list of services to test
chosen_services=()
# Array to hold process IDs
pids=()
# Array to hold exit codes
exit_codes=()

## Script functions
check_if_tracetest_is_installed() {
  if ! command -v tracetest &> /dev/null
  then
      echo "tracetest CLI could not be found"
      exit -1
  fi
}

create_env_file() {
  cat << EOF > tracetesting-vars.yaml
type: VariableSet
spec:
  id: tracetesting-vars
  name: tracetesting-vars
  values:
    - key: AD_ADDR
      value: $AD_ADDR
    - key: CART_ADDR
      value: $CART_ADDR
    - key: CHECKOUT_ADDR
      value: $CHECKOUT_ADDR
    - key: CURRENCY_ADDR
      value: $CURRENCY_ADDR
    - key: EMAIL_ADDR
      value: $EMAIL_ADDR
    - key: FRONTEND_ADDR
      value: $FRONTEND_ADDR
    - key: PAYMENT_ADDR
      value: $PAYMENT_ADDR
    - key: PRODUCT_CATALOG_ADDR
      value: $PRODUCT_CATALOG_ADDR
    - key: RECOMMENDATION_ADDR
      value: $RECOMMENDATION_ADDR
    - key: SHIPPING_ADDR
      value: $SHIPPING_ADDR
    - key: KAFKA_ADDR
      value: $KAFKA_ADDR
EOF
}

# New function to wait for services to be fully initialized
wait_for_services() {
  echo "Waiting for services to fully initialize..."
  
  # First, wait a fixed amount of time to allow services to start up
  echo "Initial delay to ensure services have time to initialize..."
  sleep 30
  
  # Check if frontend is responding
  echo "Checking if frontend service is fully ready..."
  max_retries=10
  retry_count=0
  while [ $retry_count -lt $max_retries ]; do
    if curl -s -o /dev/null -w "%{http_code}" http://$FRONTEND_ADDR | grep -q "200"; then
      echo "Frontend service is responding!"
      break
    else
      retry_count=$((retry_count+1))
      if [ $retry_count -eq $max_retries ]; then
        echo "Warning: Frontend service did not respond after $max_retries attempts, but continuing anyway..."
      else
        echo "Waiting for frontend service to respond... (attempt $retry_count/$max_retries)"
        sleep 5
      fi
    fi
  done
  
  # Check if ad service is ready by testing a simple request
  echo "Checking if ad service is fully ready..."
  retry_count=0
  while [ $retry_count -lt $max_retries ]; do
    if curl -s -X GET http://$FRONTEND_ADDR/api/data -H "Content-Type: application/json" -d '{"contextKeys":["test"]}' | grep -q "redirectUrl"; then
      echo "Ad service integration is working!"
      break
    else
      retry_count=$((retry_count+1))
      if [ $retry_count -eq $max_retries ]; then
        echo "Warning: Ad service integration did not respond correctly after $max_retries attempts, but continuing anyway..."
      else
        echo "Waiting for ad service integration to be ready... (attempt $retry_count/$max_retries)"
        sleep 5
      fi
    fi
  done
  
  # Final delay to ensure everything is stable
  echo "Final stabilization delay..."
  sleep 5
  echo "Services should be ready now. Starting tests..."
}

run_tracetest() {
  service_name=$1
  testsuite_file=./$service_name/all.yaml

  tracetest --config ./cli-config.yml run testsuite --file $testsuite_file --vars ./tracetesting-vars.yaml &
  pids+=($!)
}

## Script execution
while [[ $# -gt 0 ]]; do
  chosen_services+=("$1")
  shift
done

if [ ${#chosen_services[@]} -eq 0 ]; then
  for service in "${ALL_SERVICES[@]}"; do
    chosen_services+=("$service")
  done
fi

check_if_tracetest_is_installed
create_env_file

# Wait for services to be ready before running tests
wait_for_services

echo "Starting tests..."
echo "Running trace-based tests for: ${chosen_services[*]} ..."
echo ""

for service in "${chosen_services[@]}"; do
  run_tracetest $service
done

# Wait for processes to finish and capture their exit codes
for pid in "${pids[@]}"; do
    wait $pid
    exit_codes+=($?)
done

# Find the maximum exit code
max_exit_code=0
for code in "${exit_codes[@]}"; do
    if [[ $code -gt $max_exit_code ]]; then
        max_exit_code=$code
    fi
done

echo ""
echo "Tests done! Exit code: $max_exit_code"

exit $max_exit_code
