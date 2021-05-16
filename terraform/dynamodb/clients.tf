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

  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }

  global_secondary_index {
    name            = "clientId"
    hash_key        = "clientId"
    write_capacity  = 10
    read_capacity   = 10
    projection_type = "KEYS_ONLY"
  }
}
