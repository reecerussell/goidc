resource "aws_dynamodb_table" "clients-table" {
  name           = "goidc-clients-${var.ENV}"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "clientId"

  attribute {
    name = "clientId"
    type = "S"
  }
}
