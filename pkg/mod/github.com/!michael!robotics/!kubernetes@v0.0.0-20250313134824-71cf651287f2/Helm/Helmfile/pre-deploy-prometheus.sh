#!/bin/bash

# pre-deploy-prometheus.sh
echo "Starting pre-deployment tasks for Prometheus..."

# Example task: Checking if a specific Kubernetes namespace exists
kubectl get namespace monitoring || kubectl create namespace monitoring

# Example task: Check if Prometheus is already deployed
kubectl get deployments -n monitoring prometheus || echo "Prometheus is not deployed yet."

echo "Pre-deployment tasks for Prometheus completed."
