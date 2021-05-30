variable "api_gateway_id" {
  type = string
}

variable "root_resource_id" {
  type = string
}

variable "api_gateway_execution_arn" {
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

variable "ui_bucket" {
  type = string
}
