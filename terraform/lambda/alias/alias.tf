resource "aws_lambda_alias" "alias" {
  name             = var.name
  function_name    = var.function_arn
  function_version = "$LATEST"

  lifecycle {
    ignore_changes = [function_version]
  }
}

resource "aws_lambda_permission" "permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  qualifier     = var.name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${var.api_gateway_execution_arn}/*/*"
}