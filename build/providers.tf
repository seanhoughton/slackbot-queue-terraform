provider "aws" {}

terraform {
  backend "s3" {
    bucket = "ct-tf-state"
    key    = "slack-relay/terraform.tfstate"
    region = "us-east-1"
  }
}
