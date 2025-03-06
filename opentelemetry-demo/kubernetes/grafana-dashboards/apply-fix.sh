#!/bin/bash
# Script to apply the fix for the Grafana dashboards ConfigMap size issue

set -e  # Exit on any error

echo "Applying the fix for the Grafana dashboards ConfigMap size issue..."

# Step 1: Apply the modified manifest without the large ConfigMap
echo "Step 1: Applying the modified manifest..."
kubectl apply -f ../opentelemetry-demo-modified.yaml

# Step 2: Apply the split dashboard ConfigMaps
echo "Step 2: Applying the split dashboard ConfigMaps..."
kubectl apply -f dashboard-config.yaml

# Step 3: Apply the Grafana deployment patch
echo "Step 3: Applying the Grafana deployment patch..."
kubectl patch deployment grafana -n otel-demo --patch-file grafana-volume-patch.yaml

echo "Fix applied successfully!"
echo "You can now access Grafana to verify that the dashboards are loaded correctly."
echo "Grafana should be available at: http://localhost:30300/grafana (if using NodePort)" 