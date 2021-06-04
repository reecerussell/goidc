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
}
