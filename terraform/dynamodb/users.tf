resource "aws_dynamodb_table" "users-table" {
  name           = "goidc-users-${var.ENV}"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "userId"

  attribute {
    name = "userId"
    type = "S"
  }

  global_secondary_index {
    name            = "userId"
    hash_key        = "userId"
    write_capacity  = 10
    read_capacity   = 10
    projection_type = "KEYS_ONLY"
  }
}
