#!/bin/bash

# Define paths for the SSH key pair
SSH_KEY_PATH="$HOME/.ssh/id_rsa_minikube"
SSH_PUBLIC_KEY_PATH="$HOME/.ssh/id_rsa_minikube.pub"

# Step 1: Generate SSH key pair if it doesn't exist
if [ ! -f "$SSH_KEY_PATH" ]; then
  echo "SSH key not found, generating a new key pair..."
  ssh-keygen -t rsa -b 4096 -f "$SSH_KEY_PATH" -N ""
else
  echo "SSH key already exists at $SSH_KEY_PATH"
fi

# Step 2: Manually copy the public key to Minikube VM
echo "Copying SSH public key to Minikube VM..."
minikube ssh "mkdir -p ~/.ssh && echo $(cat $SSH_PUBLIC_KEY_PATH) >> ~/.ssh/authorized_keys"

# Step 3: Allow passwordless sudo for the 'docker' user
echo "Allowing passwordless sudo for the 'docker' user..."

# SSH into Minikube and modify the sudoers file to allow passwordless sudo
minikube ssh <<EOF
  # Add the 'docker' user to sudoers for passwordless sudo
  echo "docker ALL=(ALL) NOPASSWD:ALL" | sudo tee -a /etc/sudoers > /dev/null
EOF

# Step 4: Verify passwordless SSH and sudo
echo "Testing passwordless SSH and sudo..."

# Test SSH connection
ssh -i "$SSH_KEY_PATH" docker@$(minikube ip) "echo 'SSH is working without password'"

# Test sudo access
minikube ssh <<EOF
  sudo echo 'Sudo is working without password'
EOF

echo "Passwordless SSH and sudo setup complete!"
