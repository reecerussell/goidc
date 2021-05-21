variable "name" {
  type        = string
  description = "The name of the lambda function."
}

variable "execution_role" {
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

variable "http_method" {
  type = string
}

variable "api_gateway_id" {
  type = string
}

variable "root_resource_id" {
  type = string
}