module github.com/reecerussell/goidc/cmd/authorize

go 1.15

replace github.com/reecerussell/goidc v0.0.0 => ../../

require (
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.45
	github.com/golang/mock v1.5.0
	github.com/reecerussell/goidc v0.0.0
	github.com/reecerussell/gojwt v0.4.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
