resource "aws_api_gateway_resource" "token_proxy" {
  rest_api_id = var.api_gateway_id
  parent_id   = var.root_resource_id
  path_part   = "token"
}

module "generate_token" {
  source = "../../lambda/endpoint"

  name        = "generate-token"
  http_method = "POST"

  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.token_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region
  aws_account_id   = var.aws_account_id

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"]
}

module "generate_token_dev" {
  source = "../../lambda/alias"

  name                      = "dev"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.generate_token.function_arn
  function_name             = module.generate_token.function_name
}

module "generate_token_test" {
  source = "../../lambda/alias"

  name                      = "test"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.generate_token.function_arn
  function_name             = module.generate_token.function_name
}

module "generate_token_prod" {
  source = "../../lambda/alias"

  name                      = "prod"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.generate_token.function_arn
  function_name             = module.generate_token.function_name
}