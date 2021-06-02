variable "name" {
  type        = string
  description = "The name of the alias."
}

variable "function_arn" {
  type = string
}

variable "function_name" {
  type = string
}

variable "api_gateway_execution_arn" {
  type = string
}