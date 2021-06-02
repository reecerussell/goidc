module github.com/reecerussell/goidc/cmd/create-user

go 1.15

replace github.com/reecerussell/goidc v0.0.0 => ../../

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.40
	github.com/golang/mock v1.4.4
	github.com/google/uuid v1.2.0
	github.com/reecerussell/adaptive-password-hasher v1.0.1
	github.com/reecerussell/goidc v0.0.0
	github.com/stretchr/testify v1.7.0
)
