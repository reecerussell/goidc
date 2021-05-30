variable "name" {
  type        = string
  description = "The name of a stage, i.e. dev, test, prod."
}

variable "deployment_id" {
  type = string
}

variable "api_gateway_id" {
  type = string
}

variable "ui_bucket" {
  type        = string
  description = "The value for the UI_BUCKET stage variable"
}
