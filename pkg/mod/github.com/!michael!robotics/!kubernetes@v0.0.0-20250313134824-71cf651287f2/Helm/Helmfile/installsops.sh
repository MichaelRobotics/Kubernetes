#!/bin/bash

# Define the version of SOPS to install
SOPS_VERSION="v3.9.3"

# Download the binary
echo "Downloading SOPS $SOPS_VERSION..."
curl -LO https://github.com/getsops/sops/releases/download/$SOPS_VERSION/sops-$SOPS_VERSION.linux.amd64

# Move the binary into your PATH
echo "Moving SOPS binary to /usr/local/bin..."
sudo mv sops-$SOPS_VERSION.linux.amd64 /usr/local/bin/sops

# Make the binary executable
echo "Making the SOPS binary executable..."
sudo chmod +x /usr/local/bin/sops

# Verify the installation
echo "Verifying SOPS installation..."
sops --version

echo "SOPS installation completed successfully!"
