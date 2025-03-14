#!/bin/bash

# Path to your local TLS certificate
TLS_CERT_PATH="${BASE_PATH}/certs/tls.crt"

# Ensure the certificate file exists
if [ ! -f "$TLS_CERT_PATH" ]; then
  echo "TLS certificate not found at $TLS_CERT_PATH"
  exit 1
fi

# Copy the certificate to the Minikube VM
echo "Copying certificate to Minikube..."
scp "$TLS_CERT_PATH" docker@$(minikube ip):/tmp/

# SSH into Minikube and move the certificate to the trusted directory
echo "Moving certificate to trusted directory on Minikube..."
minikube ssh <<EOF
  sudo mv /tmp/tls.crt /usr/local/share/ca-certificates/tls.crt
  sudo update-ca-certificates --fresh
EOF

# Restart Minikube to apply the changes
echo "Restarting Minikube..."
minikube stop
minikube start

echo "TLS certificate successfully added to Minikube!"
