module "create_user" {
  source = "../../lambda/endpoint"

  name        = "create-user"
  http_method = "POST"

  aws_account_id = var.aws_account_id
  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.users_proxy.id
  root_resource_path = aws_api_gateway_resource.users_proxy.path
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"]

  depends_on = [aws_api_gateway_resource.users_proxy]
}
