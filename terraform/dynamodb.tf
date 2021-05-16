module "dynamodb-test" {
  source = "./dynamodb"

  ENV = "test"
}

module "dynamodb-dev" {
  source = "./dynamodb"

  ENV = "dev"
}

module "dynamodb-prod" {
  source = "./dynamodb"

  ENV = "prod"
}
