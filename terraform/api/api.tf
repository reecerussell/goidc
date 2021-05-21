resource "aws_api_gateway_rest_api" "api" {
  name = var.name

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "oauth_proxy" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "oauth"

  depends_on = [aws_api_gateway_rest_api.api]
}

module "oauth_endpoints" {
  source = "./oauth"

  api_gateway_id = aws_api_gateway_rest_api.api.id
  root_resource_id = aws_api_gateway_resource.oauth_proxy.id
  api_gateway_execution_arn = aws_api_gateway_rest_api.api.execution_arn
  s3_bucket = var.s3_bucket
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  execution_role = aws_iam_role.execution.arn

  depends_on = [aws_api_gateway_resource.oauth_proxy]
}