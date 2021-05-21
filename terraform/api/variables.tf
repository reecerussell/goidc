variable "name" {
  type = string
}

variable "s3_bucket" {
  type        = string
  description = "The S3 Bucket with the source code."
}

variable "aws_region" {
  type = string
}

variable "aws_account_id" {
  type = string
}