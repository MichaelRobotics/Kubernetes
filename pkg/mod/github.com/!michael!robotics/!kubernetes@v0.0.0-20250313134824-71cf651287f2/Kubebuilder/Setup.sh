#!/bin/bash

# Download Go 1.19
curl -LO https://go.dev/dl/go1.19.linux-amd64.tar.gz

# Remove any previous Go installation (if it exists)
sudo rm -rf /usr/local/go

# Extract the downloaded Go tarball to /usr/local
sudo tar -C /usr/local -xvzf go1.19.linux-amd64.tar.gz



export PATH=$PATH:/usr/local/go/bin
source ~/.bashrc # For bash

# Define the version and URL for kubebuilder
KUBEBUILDER_VERSION="v3.6.0"
KUBEBUILDER_URL="https://github.com/kubernetes-sigs/kubebuilder/releases/download/${KUBEBUILDER_VERSION}/kubebuilder_linux_amd64"

# Download kubebuilder binary
echo "Downloading kubebuilder version ${KUBEBUILDER_VERSION}..."
curl -L -o kubebuilder "${KUBEBUILDER_URL}"

# Make the binary executable
echo "Making kubebuilder executable..."
chmod +x kubebuilder

# Move the binary to /usr/local/bin
echo "Moving kubebuilder to /usr/local/bin..."
sudo mv kubebuilder /usr/local/bin/

# Verify the installation
echo "Verifying kubebuilder installation..."
kubebuilder version

echo "Kubebuilder version ${KUBEBUILDER_VERSION} has been successfully installed!"

rm go1.19.linux-amd64.tar.gz

go get github.com/nlopes/slack
go mod tidy
