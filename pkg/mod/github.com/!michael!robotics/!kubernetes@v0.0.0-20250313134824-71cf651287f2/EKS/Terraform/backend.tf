#terraform {
#  backend "s3" {
#    bucket         = "michael-s3-demo-xyz" # change this
#    key            = "michael/terraform.tfstate"
#    region         = "us-east-1"
#    encrypt        = true
#    dynamodb_table = "terraform-lock"
#  }
#}