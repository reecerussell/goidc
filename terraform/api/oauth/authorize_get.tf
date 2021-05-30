module "authorize_get" {
  source = "../../lambda/endpoint"

  name        = "authorize-ui"
  http_method = "GET"

  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.authorize_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region
  aws_account_id   = var.aws_account_id
}

module "authorize_get_dev" {
  source = "../../lambda/alias"

  name                      = "dev"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.authorize_get.function_arn
  function_name             = module.authorize_get.function_name
}

module "authorize_get_test" {
  source = "../../lambda/alias"

  name                      = "test"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.authorize_get.function_arn
  function_name             = module.authorize_get.function_name
}

module "authorize_get_prod" {
  source = "../../lambda/alias"

  name                      = "prod"
  api_gateway_execution_arn = var.api_gateway_execution_arn
  function_arn              = module.authorize_get.function_arn
  function_name             = module.authorize_get.function_name
}

resource "aws_iam_policy" "authorize_s3" {
  name        = "authorize-s3"
  path        = "/"
  description = "IAM policy for s3 for authorize"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::${var.ui_bucket}/*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "authorize_s3_attachment" {
  role       = module.authorize_get.execution_role
  policy_arn = aws_iam_policy.authorize_s3.arn

  depends_on = [aws_iam_policy.authorize_s3, module.authorize_get]
}
