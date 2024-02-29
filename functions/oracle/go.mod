module github.com/connellrobert/pixelart-serverless-api/functions/oracle

go 1.20

require (
	github.com/connellrobert/pixelart-serverless-api/functions/lib v0.0.0
	github.com/aws/aws-lambda-go v1.41.0
	github.com/sashabaranov/go-openai v1.14.1
	github.com/stretchr/testify v1.8.4
)

replace github.com/connellrobert/pixelart-serverless-api/functions/lib => ../lib
replace github.com/connellrobert/pixelart-serverless-api/functions/oracle/internal => ./internal
require (
	github.com/aws/aws-sdk-go-v2 v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.33 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.32 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.39 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.33 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.38.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.20.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.24.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/xray v1.17.2 // indirect
	github.com/aws/smithy-go v1.14.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
