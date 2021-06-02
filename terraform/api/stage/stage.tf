resource "aws_api_gateway_stage" "stage" {
  deployment_id = var.deployment_id
  rest_api_id   = var.api_gateway_id
  stage_name    = var.name

  variables = {
    CLIENTS_TABLE_NAME = "goidc-clients-${var.name}"
    USERS_TABLE_NAME   = "goidc-users-${var.name}"
    JWT_KEY_ID         = aws_kms_key.jwt.key_id
    UI_BUCKET          = var.ui_bucket
  }

  lifecycle {
    ignore_changes = [
      deployment_id,
      variables
    ]
  }
}
