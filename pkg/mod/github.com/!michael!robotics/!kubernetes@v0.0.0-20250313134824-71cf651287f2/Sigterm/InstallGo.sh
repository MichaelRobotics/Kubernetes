#!/bin/bash

# Download Go 1.19
curl -LO https://go.dev/dl/go1.19.linux-amd64.tar.gz

# Remove any previous Go installation (if it exists)
sudo rm -rf /usr/local/go

# Extract the downloaded Go tarball to /usr/local
sudo tar -C /usr/local -xvzf go1.19.linux-amd64.tar.gz

rm go1.19.linux-amd64.tar.gz

go mod tidy