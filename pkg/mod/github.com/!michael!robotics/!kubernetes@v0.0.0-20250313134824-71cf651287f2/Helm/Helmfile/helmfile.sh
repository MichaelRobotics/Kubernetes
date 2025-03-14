#!/bin/bash

# Define the version and architecture
HELMFILE_VERSION="0.169.2"
ARCH="linux_386"

# Download the Helmfile tarball
echo "Downloading Helmfile v${HELMFILE_VERSION}..."
wget -q "https://github.com/helmfile/helmfile/releases/download/v${HELMFILE_VERSION}/helmfile_${HELMFILE_VERSION}_${ARCH}.tar.gz" -O helmfile.tar.gz

# Extract the tarball
echo "Extracting Helmfile..."
tar -xzf helmfile.tar.gz

# Move the binary to /usr/local/bin
echo "Installing Helmfile..."
chmod +x helmfile
sudo mv helmfile /usr/local/bin/

# Verify the installation
echo "Verifying Helmfile installation..."
if helmfile --version; then
  echo "Helmfile v${HELMFILE_VERSION} installed successfully!"
else
  echo "Helmfile installation failed."
  exit 1
fi

# Clean up tarball
echo "Cleaning up tarball..."
rm helmfile.tar.gz

# Delete LICENSE, README-zh_CN.md, and README.md if they exist
echo "Deleting unnecessary files..."
rm -f LICENSE README-zh_CN.md README.md


helmfile init --interactive=false
