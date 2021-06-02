resource "aws_api_gateway_resource" "authorize_proxy" {
  rest_api_id = var.api_gateway_id
  parent_id   = var.root_resource_id
  path_part   = "authorize"
}

module "authorize_post" {
  source = "../../lambda/endpoint"

  name        = "authorize"
  http_method = "POST"

  aws_account_id = var.aws_account_id
  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.authorize_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"]

  depends_on = [
    aws_api_gateway_resource.authorize_proxy
  ]
}

resource "aws_iam_policy" "authorize_kms" {
  name        = "authorize-kms"
  path        = "/"
  description = "IAM policy for kms for authorize"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "kms:GetPublicKey",
          "kms:Sign"
      ],
      "Resource": "arn:aws:kms:${var.aws_region}:${var.aws_account_id}:key/*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "authorize_kms_attachment" {
  role       = module.authorize_post.execution_role
  policy_arn = aws_iam_policy.authorize_kms.arn

  depends_on = [aws_iam_policy.authorize_kms, module.authorize_post]
}

module "authorize_post_dev" {
  source = "../../lambda/alias"

  name                      = "dev"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.authorize_post.function_arn
  function_name             = module.authorize_post.function_name
}

module "authorize_post_test" {
  source = "../../lambda/alias"

  name                      = "test"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.authorize_post.function_arn
  function_name             = module.authorize_post.function_name
}

module "authorize_post_prod" {
  source = "../../lambda/alias"

  name                      = "prod"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.authorize_post.function_arn
  function_name             = module.authorize_post.function_name
}