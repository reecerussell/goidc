resource "aws_api_gateway_resource" "users_proxy" {
  rest_api_id = var.api_gateway_id
  parent_id   = var.root_resource_id
  path_part   = "users"
}