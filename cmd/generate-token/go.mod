module github.com/reecerussell/goidc/cmd/generate-token

go 1.15

replace github.com/reecerussell/goidc v0.0.0 => ../../

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.42
	github.com/golang/mock v1.4.4
	github.com/reecerussell/goidc v0.0.0
	github.com/reecerussell/gojwt v0.4.0
	github.com/stretchr/testify v1.7.0
)
