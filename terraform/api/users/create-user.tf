module "create_user" {
  source = "../../lambda/endpoint"

  name        = "create-user"
  http_method = "POST"

  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.users_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region
  aws_account_id   = var.aws_account_id

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"]

  depends_on = [aws_api_gateway_resource.users_proxy]
}
