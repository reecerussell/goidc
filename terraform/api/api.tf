resource "aws_api_gateway_rest_api" "api" {
  name = var.name

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  binary_media_types = [
    "*/*"
  ]
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

  depends_on = [
    aws_api_gateway_resource.oauth_proxy,
    aws_api_gateway_rest_api.api,
    aws_s3_bucket.ui_bucket
  ]
}

resource "aws_api_gateway_resource" "api_proxy" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "api"

  depends_on = [aws_api_gateway_rest_api.api]
}

module "users_endpoints" {
  source = "./users"

  api_gateway_id            = aws_api_gateway_rest_api.api.id
  root_resource_id          = aws_api_gateway_resource.api_proxy.id
  api_gateway_execution_arn = aws_api_gateway_rest_api.api.execution_arn
  ui_bucket                 = aws_s3_bucket.ui_bucket.bucket
  s3_bucket                 = var.s3_bucket
  aws_region                = var.aws_region
  aws_account_id            = var.aws_account_id

  depends_on = [
    aws_api_gateway_rest_api.api,
    aws_api_gateway_resource.api_proxy,
    aws_s3_bucket.ui_bucket
  ]
}