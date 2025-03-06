# Fixing Grafana Dashboards ConfigMap Size Issue

This directory contains a solution for the error:
```
The ConfigMap "grafana-dashboards" is invalid: metadata.annotations: Too long: must have at most 262144 bytes
```

## Problem

The original `opentelemetry-demo.yaml` file contains a very large ConfigMap for Grafana dashboards that exceeds Kubernetes' size limit of 262,144 bytes for annotations.

## Solution

We've split the large ConfigMap into multiple smaller ConfigMaps:

1. `grafana-dashboards-config` - Contains the dashboard provider configuration
2. `grafana-dashboard-demo` - Contains the demo dashboard

## How to Use

1. Apply the modified manifest that doesn't include the large ConfigMap:
   ```bash
   kubectl apply -f kubernetes/opentelemetry-demo-modified.yaml
   ```

2. Apply the split dashboard ConfigMaps:
   ```bash
   kubectl apply -f kubernetes/grafana-dashboards/dashboard-config.yaml
   ```

3. Apply the Grafana deployment patch to mount the new ConfigMaps:
   ```bash
   kubectl patch deployment grafana -n otel-demo --patch-file kubernetes/grafana-dashboards/grafana-volume-patch.yaml
   ```

4. Verify that Grafana can find and load the dashboards by accessing the Grafana UI.

## How This Works

- The `grafana-dashboards-config` ConfigMap contains the provider configuration that tells Grafana where to look for dashboards.
- The `grafana-dashboard-demo` ConfigMap contains the actual dashboard definition.
- The patch adds the necessary volume mounts to the Grafana deployment to use these ConfigMaps.

This approach follows Kubernetes best practices by splitting large ConfigMaps into smaller, more manageable pieces. 