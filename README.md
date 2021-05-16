# Goidc

Goidc is a basic implemetation of an Open ID Connect API, using serverless functions in AWS.

## Stack

- Go
- DynamoDB (NoSQL)
- AWS
    - KMS
    - S3
    - API Gateway

# Folder structure

- `cmd/` - contains the serverless functions
- `dal/` - holds the data-access and persistence logic
    - `dynamodb/` - an implementation of the DAL for DynamoDB
- `ui/` - contains the UI aspects, such as the login page
- `terraform/` - holds the terraform IaaC for the API
- `scripts/` - contains scripts for deploying, setting up, etc
