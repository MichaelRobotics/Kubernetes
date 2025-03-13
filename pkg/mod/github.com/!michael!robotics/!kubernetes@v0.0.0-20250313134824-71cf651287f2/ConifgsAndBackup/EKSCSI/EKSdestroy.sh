#!/bin/bash

set -e  # Exit on error
# Destroy EKS_2
cd EKS_us_east_2
terraform destroy -auto-approve

# Destroy EKS_1
echo "Destroy_EKS_1..."
cd ../EKS_us_east_1
terraform destroy -var "bucket=$bucket_name" -auto-approve
