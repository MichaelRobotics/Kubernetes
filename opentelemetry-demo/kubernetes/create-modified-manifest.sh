#!/bin/bash
# This script creates a modified version of the OpenTelemetry demo Kubernetes manifest
# It removes the large grafana-dashboards ConfigMap and adds a reference to use our split ConfigMaps

# Create a temporary file
cat kubernetes/opentelemetry-demo.yaml | \
  # Remove the large grafana-dashboards ConfigMap section
  sed '/# Source: opentelemetry-demo\/templates\/grafana-dashboards.yaml/,/---/d' > \
  kubernetes/opentelemetry-demo-modified.yaml

# Add a note at the top of the file
sed -i '1i# Modified version of opentelemetry-demo.yaml with grafana-dashboards ConfigMap removed\n# Use this file along with the split dashboard ConfigMaps in grafana-dashboards/dashboard-config.yaml' \
  kubernetes/opentelemetry-demo-modified.yaml

echo "Modified manifest created at kubernetes/opentelemetry-demo-modified.yaml"
echo "Use this file along with the split dashboard ConfigMaps in grafana-dashboards/dashboard-config.yaml" 