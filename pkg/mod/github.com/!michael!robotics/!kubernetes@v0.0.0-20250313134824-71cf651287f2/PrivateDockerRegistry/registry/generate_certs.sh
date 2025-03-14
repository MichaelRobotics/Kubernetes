openssl req -x509 -newkey rsa:4096 -days 365 -nodes -sha256 \
  -keyout certs/tls.key -out certs/tls.crt \
  -subj "/CN=my-registry" \
  -addext "subjectAltName=DNS:my-registry"



# Variables
TLS_CERT_PATH="${BASE_PATH}/certs/tls.crt"  # Path to your tls.crt file
SYSTEM_CA_DIR="/usr/local/share/ca-certificates"  # Directory for system-wide certificates
DOCKER_CERT_DIR="/etc/docker/certs.d/my-registry:30000"  # Docker directory for the registry certificate

# Check if the TLS certificate exists
if [[ ! -f "$TLS_CERT_PATH" ]]; then
  echo "Error: TLS certificate file not found at $TLS_CERT_PATH."
  exit 1
fi

# Step 1: Copy the TLS certificate to /usr/local/share/ca-certificates
echo "Copying TLS certificate to $SYSTEM_CA_DIR..."
sudo cp "$TLS_CERT_PATH" "$SYSTEM_CA_DIR/tls.crt"

# Check if the file was copied successfully
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to copy the certificate to $SYSTEM_CA_DIR."
  exit 1
else
  echo "TLS certificate copied to $SYSTEM_CA_DIR successfully."
fi

# Step 2: Update the CA certificates
echo "Updating CA certificates..."
sudo update-ca-certificates --fresh

# Check if the update was successful
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to update CA certificates."
  exit 1
else
  echo "CA certificates updated successfully."
fi

# Step 3: Copy the TLS certificate to Docker's registry directory
echo "Creating Docker certificate directory and copying the TLS certificate..."
sudo mkdir -p "$DOCKER_CERT_DIR"
sudo cp "$TLS_CERT_PATH" "$DOCKER_CERT_DIR/tls.crt"

# Check if the file was copied successfully to Docker
if [[ $? -ne 0 ]]; then
  echo "Error: Failed to copy the certificate to Docker's registry directory."
  exit 1
else
  echo "TLS certificate copied to Docker registry directory successfully."
fi

echo "Script execution completed."

