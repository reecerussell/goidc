resource "aws_api_gateway_resource" "users_proxy" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = var.root_resource_id
  path_part   = "users"
}