module "authorize_get" {
  source = "../../lambda/endpoint"

  name        = "authorize-ui"
  http_method = "GET"

  aws_account_id = var.aws_account_id
  api_gateway_id   = var.api_gateway_id
  root_resource_id = aws_api_gateway_resource.authorize_proxy.id
  s3_bucket        = var.s3_bucket
  aws_region       = var.aws_region
  content_handling = "CONVERT_TO_BINARY"

  depends_on = [aws_api_gateway_resource.authorize_proxy]
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
          "s3:GetObject*"
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

resource "aws_api_gateway_resource" "authorize_ui_proxy" {
  rest_api_id = var.api_gateway_id
  parent_id   = aws_api_gateway_resource.authorize_proxy.id
  path_part   = "{proxy+}"

  depends_on = [module.authorize_get, aws_api_gateway_resource.authorize_proxy]
}

resource "aws_api_gateway_method" "authorize_ui_get" {
  rest_api_id   = var.api_gateway_id
  resource_id   = aws_api_gateway_resource.authorize_ui_proxy.id
  http_method   = "GET"
  authorization = "NONE"

  depends_on = [
    aws_api_gateway_resource.authorize_ui_proxy
  ]
}

resource "aws_api_gateway_integration" "authorize_ui_integration" {
  rest_api_id      = var.api_gateway_id
  resource_id      = aws_api_gateway_resource.authorize_ui_proxy.id
  http_method      = aws_api_gateway_method.authorize_ui_get.http_method
  content_handling = "CONVERT_TO_BINARY"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.aws_region}:${var.aws_account_id}:function:${module.authorize_get.function_name}:$${stageVariables.ENVIRONMENT}/invocations"

  depends_on = [
    aws_api_gateway_resource.authorize_ui_proxy,
    aws_api_gateway_method.authorize_ui_get
  ]
}

resource "aws_api_gateway_method_response" "authorize_ui_ok" {
  rest_api_id = var.api_gateway_id
  resource_id = aws_api_gateway_resource.authorize_ui_proxy.id
  http_method = aws_api_gateway_method.authorize_ui_get.http_method
  status_code = "200"
  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Origin"  = true
  }

  depends_on = [
    aws_api_gateway_resource.authorize_ui_proxy,
    aws_api_gateway_method.authorize_ui_get
  ]
}

resource "aws_api_gateway_integration_response" "authorize_ui_response_integration" {
  rest_api_id      = var.api_gateway_id
  resource_id      = aws_api_gateway_resource.authorize_ui_proxy.id
  http_method      = aws_api_gateway_method.authorize_ui_get.http_method
  content_handling = "CONVERT_TO_BINARY"
  status_code      = aws_api_gateway_method_response.authorize_ui_ok.status_code

  depends_on = [
    aws_api_gateway_method_response.authorize_ui_ok,
    aws_api_gateway_resource.authorize_ui_proxy,
    aws_api_gateway_method.authorize_ui_get
  ]
}
