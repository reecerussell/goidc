resource "aws_api_gateway_deployment" "default" {
  rest_api_id = aws_api_gateway_rest_api.api.id

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_api_gateway_rest_api.api
  ]
}

module "dev_stage" {
  source = "./stage"

  name           = "dev"
  api_gateway_id = aws_api_gateway_rest_api.api.id
  deployment_id  = aws_api_gateway_deployment.default.id

  depends_on = [
    aws_api_gateway_deployment.default
  ]
}

module "test_stage" {
  source = "./stage"

  name           = "test"
  api_gateway_id = aws_api_gateway_rest_api.api.id
  deployment_id  = aws_api_gateway_deployment.default.id

  depends_on = [
    aws_api_gateway_deployment.default
  ]
}

module "prod_stage" {
  source = "./stage"

  name           = "prod"
  api_gateway_id = aws_api_gateway_rest_api.api.id
  deployment_id  = aws_api_gateway_deployment.default.id

  depends_on = [
    aws_api_gateway_deployment.default
  ]
}
