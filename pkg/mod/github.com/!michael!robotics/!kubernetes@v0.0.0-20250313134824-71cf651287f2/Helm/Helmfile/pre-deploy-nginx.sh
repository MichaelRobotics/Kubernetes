#!/bin/bash

# pre-deploy-nginx.sh
echo "Starting pre-deployment tasks for Nginx..."

# Example task: Checking if a specific Kubernetes namespace exists
kubectl get namespace web || kubectl create namespace web

# Example task: Check if Nginx is already deployed
kubectl get deployments -n web nginx || echo "Nginx is not deployed yet."

echo "Pre-deployment tasks for Nginx completed."
