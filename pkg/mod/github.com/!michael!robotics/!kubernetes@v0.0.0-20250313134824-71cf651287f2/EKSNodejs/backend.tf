terraform {
  backend "s3" {
    bucket         = "dot-20250117203252"
    region         = "us-east-1"
    key            = "EKSNodejs/.terraform/terraform.tfstate"
    dynamodb_table = "state-lock"
    encrypt        = true
  }
  required_version = ">=0.13.0"
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source  = "hashicorp/aws"
    }
  }
}