resource "aws_api_gateway_resource" "authorize_proxy" {
  rest_api_id = var.api_gateway_id
  parent_id   = var.root_resource_id
  path_part   = "authorize"
}

module "authorize_post" {
  source = "../../lambda/endpoint"

  name        = "authorize"
  http_method = "POST"

  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.authorize_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region
  aws_account_id   = var.aws_account_id

  iam_policies = ["arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess"]
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
