resource "aws_api_gateway_resource" "token_proxy" {
  rest_api_id = var.api_gateway_id
  parent_id   = var.root_resource_id
  path_part   = "token"
}

module "generate_token" {
  source = "../../lambda/endpoint"

  name        = "generate-token"
  http_method = "POST"

  aws_account_id = var.aws_account_id
  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.token_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"]

  depends_on = [
    aws_api_gateway_resource.token_proxy
  ]
}

resource "aws_iam_policy" "generate_token_kms" {
  name        = "generate-token-kms"
  path        = "/"
  description = "IAM policy for kms for generate-token"

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

resource "aws_iam_role_policy_attachment" "generate_token_kms_attachment" {
  role       = module.generate_token.execution_role
  policy_arn = aws_iam_policy.generate_token_kms.arn

  depends_on = [aws_iam_policy.generate_token_kms, module.generate_token]
}
