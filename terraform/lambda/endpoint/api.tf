resource "aws_api_gateway_method" "method" {
  rest_api_id   = var.api_gateway_id
  resource_id   = var.root_resource_id
  http_method   = var.http_method
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "integration" {
  rest_api_id      = var.api_gateway_id
  resource_id      = var.root_resource_id
  http_method      = aws_api_gateway_method.method.http_method
  content_handling = var.content_handling

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.aws_region}:${var.aws_account_id}:function:${aws_lambda_function.function.function_name}:$${stageVariables.ENVIRONMENT}/invocations"

  depends_on = [
    aws_lambda_function.function,
    aws_api_gateway_method.method
  ]
}

resource "aws_api_gateway_method_response" "ok" {
  rest_api_id = var.api_gateway_id
  resource_id = var.root_resource_id
  http_method = aws_api_gateway_method.method.http_method
  status_code = "200"
  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true,
    "method.response.header.Access-Control-Allow-Methods" = true,
    "method.response.header.Access-Control-Allow-Origin"  = true
  }

  depends_on = [aws_api_gateway_method.method]
}

resource "aws_api_gateway_integration_response" "integration_response" {
  rest_api_id      = var.api_gateway_id
  resource_id      = var.root_resource_id
  http_method      = aws_api_gateway_method.method.http_method
  content_handling = var.content_handling
  status_code      = aws_api_gateway_method_response.ok.status_code

  depends_on = [
    aws_api_gateway_method_response.ok,
    aws_api_gateway_method.method,
    aws_api_gateway_integration.integration
  ]
}
