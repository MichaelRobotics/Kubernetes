#!/bin/bash

# Prompt for user input to pass to GPG key generation
echo "Enter your name for the GPG key:"
read GPG_NAME

echo "Enter your email for the GPG key:"
read GPG_EMAIL

# Step 1: Generate a GPG key without specifying a passphrase
echo "Generating GPG key..."
cat <<EOF | gpg --batch --generate-key
%no-protection
Key-Type: 1
Key-Length: 2048
Name-Real: $GPG_NAME
Name-Email: $GPG_EMAIL
Expire-Date: 0
EOF

# Step 2: List secret keys and extract the GPG Key ID
echo "Listing GPG secret keys..."
KEY_ID=$(gpg --list-secret-keys --keyid-format LONG | grep sec | awk '{print $2}' | cut -d'/' -f2)

# Step 3: Display the generated Key ID
echo "Generated GPG Key ID: $KEY_ID"

# Step 4: Encrypt the values-secrets-to-encrypt.yaml file using the generated GPG key
echo "Encrypting the values-secrets-to-encrypt.yaml file..."
sops -e --pgp "$KEY_ID" values-secrets.yaml > values-secrets-encrypted.yaml

echo "Encryption complete. Encrypted file saved as values-secrets-encrypted.yaml"
