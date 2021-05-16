resource "aws_dynamodb_table" "clients-table" {
  name           = "goidc-clients"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "clientId"
  range_key      = "name"

  attribute {
    name = "clientId"
    type = "S"
  }

  attribute {
    name = "name"
    type = "S"
  }

  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }

  global_secondary_index {
    name               = "clientId"
    hash_key           = "clientId"
    range_key          = "name"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "INCLUDE"
    non_key_attributes = ["name"]
  }
}
