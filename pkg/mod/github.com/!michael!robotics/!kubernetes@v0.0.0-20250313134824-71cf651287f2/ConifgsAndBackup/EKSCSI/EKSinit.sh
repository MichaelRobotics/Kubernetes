#!/bin/bash

set -e  # Exit on error

# Check if the bucket name argument is provided
if [ -z "$1" ]; then
  echo "Error: No bucket name provided."
  echo "Usage: ./EKSinit.sh <bucket_name>"
  exit 1
fi

# Assign the first argument to the bucket_name variable
bucket_name=$1

# Print the bucket name
echo "Bucket name is: $bucket_name"

# Deploy EKS_1
echo "Initializing and applying EKS_1..."
cd EKS_us_east_1
terraform init
terraform validate
terraform plan -var "bucket=$bucket_name" 
terraform apply -var "bucket=$bucket_name" -auto-approve
# Run terraform output to get the outputs in JSON format
terraform output -json > terraform_outputs.json

# Define the output INI file
OUTPUT_FILE="../credentials-velero"

# Extract values from the terraform outputs (adjust output names based on your Terraform configuration)
AWS_VELERO_USER_ACCESS_KEY_ID=$(jq -r '.velero_access_key_id.value' terraform_outputs.json)
AWS_VELERO_USER_ACCESS_KEY=$(jq -r '.velero_secret_access_key.value' terraform_outputs.json)
rm terraform_outputs.json

# Create or overwrite the INI file
echo "[default]" > "$OUTPUT_FILE"
echo "aws_access_key_id = $AWS_VELERO_USER_ACCESS_KEY_ID" >> "$OUTPUT_FILE"
echo "aws_secret_access_key = $AWS_VELERO_USER_ACCESS_KEY" >> "$OUTPUT_FILE"

# Output message to indicate secret file creation completion
echo "Terraform secrets saved to $OUTPUT_FILE"

# Deploy EKS_2
echo "Initializing and applying EKS_2..."
cd ../EKS_us_east_2
terraform init
terraform validate
terraform plan 
terraform apply -auto-approve