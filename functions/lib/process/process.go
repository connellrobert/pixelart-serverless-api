package process

// environment variables required for this file:
// RESULT_FUNCTION_ARN - the ARN of the result function
// GENERATE_IMAGE_TABLE_NAME - the name of the generate image table
// EDIT_IMAGE_TABLE_NAME - the name of the edit image table
// VARIATE_IMAGE_TABLE_NAME - the name of the variate image table

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	. "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/google/uuid"
)

type process interface {
	LoadDefaultConfig(ctx context.Context, fn ...func(*config.LoadOptions) error) (cfg aws.Config, err error)
	WithRegion(v string) config.LoadOptionsFunc
	SendMessage(sqsClient *sqs.Client, ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	PutTraceSegments(xrayClient *xray.Client, ctx context.Context, params *xray.PutTraceSegmentsInput, optFns ...func(*xray.Options)) (*xray.PutTraceSegmentsOutput, error)
	PutItem(client *dynamodb.Client, ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetSecretValue(client secretsmanager.Client, ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
	GetS3Object(client s3.Client, ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	NewFile(fd uintptr, name string) *os.File
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

type proc struct{}

func (p proc) LoadDefaultConfig(ctx context.Context, fn ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
	return config.LoadDefaultConfig(ctx, fn...)
}

func (p proc) WithRegion(v string) config.LoadOptionsFunc {
	return config.WithRegion(v)
}

func (p proc) SendMessage(sqsClient *sqs.Client, ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return sqsClient.SendMessage(ctx, params, optFns...)
}

func (p proc) PutTraceSegments(xrayClient *xray.Client, ctx context.Context, params *xray.PutTraceSegmentsInput, optFns ...func(*xray.Options)) (*xray.PutTraceSegmentsOutput, error) {
	return xrayClient.PutTraceSegments(ctx, params, optFns...)
}

func (p proc) PutItem(client *dynamodb.Client, ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return client.PutItem(ctx, params, optFns...)
}

func (p proc) GetSecretValue(client secretsmanager.Client, ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return client.GetSecretValue(ctx, params, optFns...)
}

func (p proc) GetS3Object(client s3.Client, ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return client.GetObject(ctx, params, optFns...)
}

func (p proc) NewFile(fd uintptr, name string) *os.File {
	return os.NewFile(fd, name)
}

func (p proc) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

var p process

func init() {
	p = proc{}
}

func GetAWSConfig() aws.Config {
	cfg, err := p.LoadDefaultConfig(context.TODO(),
		p.WithRegion(Region()),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Region() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	return region
}

func SendResult(record QueueRequest, response ImageResponseWrapper) {
	// Create a Lambda client
	// invoke lambda
	tmp := ResultRequest{
		Record: record,
		Result: response,
	}
	// Send req to sqs queue
	queue := os.Getenv("RESULT_QUEUE_URL")
	// send item to queue
	messageInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(tmp.Record.ToString())),
		QueueUrl:    aws.String(queue),
	}
	sendSqsMessage(*messageInput)
}

func SendRequestToQueue(record QueueRequest) {
	queueUrl := os.Getenv("QUEUE_URL")
	// send item to queue
	sendMessageInput := &sqs.SendMessageInput{
		MessageBody:            aws.String(record.ToString()),
		QueueUrl:               aws.String(queueUrl),
		MessageGroupId:         aws.String("1"),
		MessageDeduplicationId: aws.String(record.Id),
	}
	sendSqsMessage(*sendMessageInput)
}

func sendSqsMessage(message sqs.SendMessageInput) {
	sqsClient := sqs.NewFromConfig(GetAWSConfig())
	// send item to queue
	_, err := p.SendMessage(sqsClient, context.Background(), &message)
	if err != nil {
		panic(err)
	}
}

func GetTraceId() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)[:16]
}

func SubmitXRayTraceSubSegment(parentSegmentId, name string) {
	id := GetTraceId()
	xrayClient := xray.NewFromConfig(GetAWSConfig())
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	doc := XrayTraceSegmentDocument{
		TraceId:   parentSegmentId,
		Id:        id,
		Name:      name,
		StartTime: float64(startTime),
		EndTime:   float64(startTime + 10000),
	}
	submitSubSegmentInput := &xray.PutTraceSegmentsInput{
		TraceSegmentDocuments: []string{
			doc.ToString(),
		},
	}
	_, err := p.PutTraceSegments(xrayClient, context.Background(), submitSubSegmentInput)
	if err != nil {
		panic(err)
	}
}

func StoreAnalyticsItem(ai AnalyticsItem) {
	client := dynamodb.NewFromConfig(GetAWSConfig())

	analyticsTable := os.Getenv("ANALYTICS_TABLE_NAME")
	putAItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(analyticsTable),
		Item:      ai.ToDynamoDB(),
	}
	_, err := p.PutItem(client, context.Background(), putAItemInput)
	if err != nil {
		panic(err)
	}
}

func GetSecretValue(envVar string) string {

	secretName := os.Getenv(envVar)

	config := GetAWSConfig()

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := p.GetSecretValue(*svc, context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		panic(err)
	}

	// Decrypts secret using the associated KMS key.
	return *result.SecretString
}

func GetImageFromS3(imageName string) io.ReadCloser {
	// create s3 client
	config := GetAWSConfig()
	svc := *s3.NewFromConfig(config)
	obj := &s3.GetObjectInput{
		Bucket: aws.String("openai-image-storage"),
		Key:    aws.String(imageName),
	}
	result, err := p.GetS3Object(svc, context.Background(), obj)
	if err != nil {
		panic(err)
	}
	return result.Body
}

func SaveFile(fileName string, fileContents io.ReadCloser) *os.File {
	file := p.NewFile(4, fileName)
	_, err := p.Copy(file, fileContents)
	if err != nil {
		panic(err)
	}
	return file
}
