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

  api_gateway_id            = aws_api_gateway_rest_api.api.id
  root_resource_id          = aws_api_gateway_resource.oauth_proxy.id
  api_gateway_execution_arn = aws_api_gateway_rest_api.api.execution_arn
  ui_bucket                 = aws_s3_bucket.ui_bucket.bucket
  s3_bucket                 = var.s3_bucket
  aws_region                = var.aws_region
  aws_account_id            = var.aws_account_id

  depends_on = [aws_api_gateway_resource.oauth_proxy]
}
