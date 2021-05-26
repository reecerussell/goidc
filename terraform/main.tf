provider "aws" {
  region     = var.AWS_REGION
  access_key = var.AWS_ACCESS_KEY
  secret_key = var.AWS_SECRET_KEY
}

terraform {
  backend "s3" {
    bucket = "goidc-terraform-state"
    key    = ".tfstate"
    region = "eu-west-2"
  }
}

module "api" {
  source = "./api"

  name           = "goidc"
  s3_bucket      = var.S3_SOURCE_BUCKET
  aws_region     = var.AWS_REGION
  aws_account_id = var.AWS_ACCOUNT_ID
}
