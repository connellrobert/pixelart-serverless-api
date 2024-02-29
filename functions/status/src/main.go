package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/connellrobert/pixelart-serverless-api/functions/lib/process"
	aiTypes "github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
	"github.com/connellrobert/pixelart-serverless-api/functions/status/internal"
)

// List of environment variables:
// ANALYTICS_TABLE_NAME

type subprocess interface {
	GetAWSConfig() aws.Config
	SubmitXRayTraceSubSegment(traceId string, name string)
	GetAnalyticsItem(id, tableName string, client *dynamodb.Client) map[string]types.AttributeValue
	FromDynamoDB(record map[string]types.AttributeValue, obj interface {
		FromDynamoDB(record map[string]types.AttributeValue)
	})
	GetAnalyticsItemAttemptsUrls(attempts map[string]aiTypes.ImageResponseWrapper) []string
	CreateResponse(urls []string, ai aiTypes.AnalyticsItem) (events.APIGatewayProxyResponse, error)
}

type sub struct{}

func (s sub) GetAWSConfig() aws.Config {
	return process.GetAWSConfig()
}

func (s sub) SubmitXRayTraceSubSegment(traceId string, name string) {
	process.SubmitXRayTraceSubSegment(traceId, name)
}

func (s sub) GetAnalyticsItem(id, tableName string, client *dynamodb.Client) map[string]types.AttributeValue {
	return internal.GetAnalyticsItem(id, tableName, client)
}

func (s sub) FromDynamoDB(record map[string]types.AttributeValue, obj interface {
	FromDynamoDB(record map[string]types.AttributeValue)
}) {
	obj.FromDynamoDB(record)
}

func (s sub) GetAnalyticsItemAttemptsUrls(attempts map[string]aiTypes.ImageResponseWrapper) []string {
	return internal.GetAnalyticsItemAttemptsUrls(attempts)
}

func (s sub) CreateResponse(urls []string, ai aiTypes.AnalyticsItem) (events.APIGatewayProxyResponse, error) {
	return internal.CreateResponse(urls, ai)
}

var subc subprocess = sub{}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get id from request path
	id := req.PathParameters["id"]
	tableName := os.Getenv("ANALYTICS_TABLE_NAME")

	// check analytics db for id
	cfg := subc.GetAWSConfig()
	client := dynamodb.NewFromConfig(cfg)
	record := subc.GetAnalyticsItem(id, tableName, client)
	analyticsItem := aiTypes.AnalyticsItem{}
	subc.FromDynamoDB(record, &analyticsItem)
	subc.SubmitXRayTraceSubSegment(analyticsItem.Record.Metadata.TraceId, "Retrieved analytics item from db")

	urls := subc.GetAnalyticsItemAttemptsUrls(analyticsItem.Attempts)
	return subc.CreateResponse(urls, analyticsItem)
}

func main() {
	lambda.Start(Handler)
}
