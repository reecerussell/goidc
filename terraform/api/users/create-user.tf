module "create_user" {
  source = "../../lambda/endpoint"

  name        = "create-user"
  http_method = "POST"

  api_gateway_id   = var.api_gateway_id
  root_resource_id = var.root_resource_id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region
  aws_account_id   = var.aws_account_id

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"]
}

module "create_user_dev" {
  source = "../../lambda/alias"

  name                      = "dev"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.create_user.function_arn
  function_name             = module.create_user.function_name
}

module "create_user_test" {
  source = "../../lambda/alias"

  name                      = "test"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.create_user.function_arn
  function_name             = module.create_user.function_name
}

module "create_user_prod" {
  source = "../../lambda/alias"

  name                      = "prod"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.create_user.function_arn
  function_name             = module.create_user.function_name
}
