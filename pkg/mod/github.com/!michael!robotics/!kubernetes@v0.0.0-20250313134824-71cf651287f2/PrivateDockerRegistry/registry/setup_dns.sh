#!/bin/bash

# Set the registry name
REGISTRY_NAME="my-registry"

# Print the registry name
echo "Registry Name: $REGISTRY_NAME"

# Get the Minikube IP address
REGISTRY_IP=$(minikube ip)

# Check if the Minikube IP was found
if [ -z "$REGISTRY_IP" ]; then
  echo "Error: Could not retrieve Minikube IP."
  exit 1
fi

# Print the Minikube IP address
echo "Minikube IP is: $REGISTRY_IP"

# Export the variables for use in other commands or scripts
export REGISTRY_NAME
export REGISTRY_IP

# Update /etc/hosts to associate the registry name with the Minikube IP
echo "$REGISTRY_IP $REGISTRY_NAME" | sudo tee -a /etc/hosts > /dev/null



