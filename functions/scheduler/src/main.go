package main

import (
	"context"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib/process"
	aiTypes "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aimless-it/ai-canvas/functions/scheduler/internal"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

// List of environment variables:
// ANALYTICS_TABLE_NAME
// EMPTY_DB_ALARM_NAME

// lambda handler
/*
{
	"action": 0,
	"params": {
		"image": "https://aimless.ai/images/ai-canvas-logo.png",
		"size": "512x512",
		"prompt": "something simple",
		"n": 1,
		"responseFormat": "URL",
		"user": "user-id"
	}
}
"{\"action\":0,\"params\":{\"image\":\"https://aimless.ai/images/ai-canvas-logo.png\",\"size\":\"512x512\",\"prompt\":\"something simple\",\"n\":1,\"responseFormat\":\"URL\",\"user\":\"user-id\"}}"
*/

type subprocess interface {
	ParseApiRequest(request events.APIGatewayProxyRequest) map[string]interface{}
	ParseRequestAction(body map[string]interface{}) aiTypes.RequestAction
	ConstructQueueRequest(args internal.QueueRequestArgs) aiTypes.QueueRequest
	PutAnalyticsItem(ai aiTypes.AnalyticsItem, tableName string, dbClient *dynamodb.Client)
	SubmitXRayTraceSubSegment(traceId string, name string)
	GetAWSConfig() aws.Config
	SendRequestToQueue(record aiTypes.QueueRequest)
	ApiResponse(ai aiTypes.AnalyticsItem) (events.APIGatewayProxyResponse, error)
}

type sub struct{}

func (s sub) ParseApiRequest(request events.APIGatewayProxyRequest) map[string]interface{} {
	return internal.ParseApiRequest(request)
}

func (s sub) ParseRequestAction(body map[string]interface{}) aiTypes.RequestAction {
	return internal.ParseRequestAction(body)
}

func (s sub) ConstructQueueRequest(args internal.QueueRequestArgs) aiTypes.QueueRequest {
	return internal.ConstructQueueRequest(args)
}

func (s sub) PutAnalyticsItem(ai aiTypes.AnalyticsItem, tableName string, dbClient *dynamodb.Client) {
	internal.PutAnalyticsItem(ai, tableName, dbClient)
}

func (s sub) SubmitXRayTraceSubSegment(traceId string, name string) {
	process.SubmitXRayTraceSubSegment(traceId, name)
}

func (s sub) GetAWSConfig() aws.Config {
	return process.GetAWSConfig()
}

func (s sub) SendRequestToQueue(record aiTypes.QueueRequest) {
	process.SendRequestToQueue(record)
}

func (s sub) ApiResponse(ai aiTypes.AnalyticsItem) (events.APIGatewayProxyResponse, error) {
	return internal.ApiResponse(ai)
}

var subProcess subprocess = sub{}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	traceId := request.Headers["x-amzn-trace-id"]
	id := uuid.New().String()
	body := subProcess.ParseApiRequest(request)
	requestAction := subProcess.ParseRequestAction(body)

	record := subProcess.ConstructQueueRequest(internal.QueueRequestArgs{
		Id:      id,
		Action:  requestAction,
		Params:  body["params"],
		TraceId: traceId,
	})

	cfg := subProcess.GetAWSConfig()
	client := dynamodb.NewFromConfig(cfg)
	// Create Analytics Item
	analyticsItem := aiTypes.AnalyticsItem{
		Id:       record.Id,
		Record:   record,
		Attempts: make(map[string]aiTypes.ImageResponseWrapper),
	}
	// Store Analytics Item
	analyticsTable := os.Getenv("ANALYTICS_TABLE_NAME")

	subProcess.PutAnalyticsItem(analyticsItem, analyticsTable, client)

	subProcess.SendRequestToQueue(record)
	subProcess.SubmitXRayTraceSubSegment(traceId, "Added item "+record.Id+" to queue")

	return subProcess.ApiResponse(analyticsItem)
}

func main() {
	lambda.Start(Handler)
}
