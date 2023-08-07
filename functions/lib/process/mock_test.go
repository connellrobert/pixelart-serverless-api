package process

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/stretchr/testify/mock"
)

const (
	region        = "us-east-1"
	testUrl       = "test_url"
	testResultUrl = "test_result_url"
)

func init() {
	os.Setenv("RESULT_QUEUE_URL", testResultUrl)
	os.Setenv("AWS_REGION", region)
	os.Setenv("QUEUE_URL", testUrl)
}

type mockProcess struct {
	mock.Mock
}

func (p mockProcess) LoadDefaultConfig(ctx context.Context, fn ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	args := p.Called(ctx, fn)
	return args.Get(0).(aws.Config), args.Error(1)
}

func (p mockProcess) WithRegion(v string) config.LoadOptionsFunc {
	args := p.Called(v)
	return args.Get(0).(config.LoadOptionsFunc)
}

func (p mockProcess) SendMessage(sqsClient *sqs.Client, ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	args := p.Called(sqsClient, ctx, params, optFns)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

func (p mockProcess) PutTraceSegments(xrayClient *xray.Client, ctx context.Context, params *xray.PutTraceSegmentsInput, optFns ...func(*xray.Options)) (*xray.PutTraceSegmentsOutput, error) {
	args := p.Called(xrayClient, ctx, params, optFns)
	return args.Get(0).(*xray.PutTraceSegmentsOutput), args.Error(1)
}

func (p mockProcess) PutItem(client *dynamodb.Client, ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := p.Called(client, ctx, params, optFns)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (p mockProcess) GetSecretValue(client secretsmanager.Client, ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	args := p.Called(client, ctx, params, optFns)
	return args.Get(0).(*secretsmanager.GetSecretValueOutput), args.Error(1)
}

func (p mockProcess) GetS3Object(client s3.Client, ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := p.Called(client, ctx, params, optFns)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func (p mockProcess) NewFile(fd uintptr, name string) *os.File {
	args := p.Called(fd, name)
	return args.Get(0).(*os.File)
}

func (p mockProcess) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	args := p.Called(dst, src)
	return args.Get(0).(int64), args.Error(1)
}

func init() {
	p = &mockProcess{}
}
