package main

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/connellrobert/pixelart-serverless-api/functions/lib/process"
	aiTypes "github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
	"github.com/connellrobert/pixelart-serverless-api/functions/result/internal"
)

type subProcess interface {
	ParseSQSEvent(message events.SQSMessage) aiTypes.ResultRequest
	GetAWSConfig() aws.Config
	GetAnalyticsItemInputStruct(id, tableName string) *dynamodb.GetItemInput
	GetDBItem(client *dynamodb.Client, ctx context.Context, getItemInput *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	GetUpdateAnalyticsItemInput(analyticsItem aiTypes.AnalyticsItem, tableName string) *dynamodb.UpdateItemInput
	UpdateDBItem(client *dynamodb.Client, ctx context.Context, updateItemInput *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	SendRequestToQueue(record aiTypes.QueueRequest)
	SubmitXRayTraceSubSegment(traceId, name string)
}

type subproc struct{}

func (s subproc) ParseSQSEvent(message events.SQSMessage) aiTypes.ResultRequest {
	return internal.ParseSQSEvent(message)
}

func (s subproc) GetAWSConfig() aws.Config {
	return process.GetAWSConfig()
}

func (s subproc) GetAnalyticsItemInputStruct(id, tableName string) *dynamodb.GetItemInput {
	return internal.GetAnalyticsItemInputStruct(id, tableName)
}

func (s subproc) GetDBItem(client *dynamodb.Client, ctx context.Context, getItemInput *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return client.GetItem(ctx, getItemInput)
}

func (s subproc) GetUpdateAnalyticsItemInput(analyticsItem aiTypes.AnalyticsItem, tableName string) *dynamodb.UpdateItemInput {
	return internal.GetUpdateAnalyticsItemInput(analyticsItem, tableName)
}

func (s subproc) UpdateDBItem(client *dynamodb.Client, ctx context.Context, updateItemInput *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return client.UpdateItem(ctx, updateItemInput)
}

func (s subproc) SendRequestToQueue(record aiTypes.QueueRequest) {
	process.SendRequestToQueue(record)
}

func (s subproc) SubmitXRayTraceSubSegment(traceId, name string) {
	process.SubmitXRayTraceSubSegment(traceId, name)
}

var subc subProcess = subproc{}

// List of environment variables:
// ANALYTICS_TABLE_NAME

func Handler(ctx context.Context, sqsResult events.SQSEvent) {

	for _, message := range sqsResult.Records {
		rr := subc.ParseSQSEvent(message)
		cfg := subc.GetAWSConfig()

		tableName := os.Getenv("ANALYTICS_TABLE_NAME")
		getItemInput := subc.GetAnalyticsItemInputStruct(rr.Record.Id, tableName)

		client := dynamodb.NewFromConfig(cfg)
		record, err := subc.GetDBItem(client, context.Background(), getItemInput)
		if err != nil {
			panic(err)
		}
		analyticsItem := aiTypes.AnalyticsItem{}
		analyticsItem.FromDynamoDB(record.Item)
		attemptNum := len(analyticsItem.Attempts)
		analyticsItem.Attempts[strconv.Itoa(attemptNum)] = rr.Result
		analyticsItem.Success = rr.Result.Success

		updateItemInput := subc.GetUpdateAnalyticsItemInput(analyticsItem, tableName)

		_, err = subc.UpdateDBItem(client, context.Background(), updateItemInput)

		if err != nil {
			panic(err)
		}

		if attemptNum < 3 && !rr.Result.Success {
			subc.SendRequestToQueue(rr.Record)
		}
		subc.SubmitXRayTraceSubSegment(rr.Record.Metadata.TraceId, "Updated analytics item")
	}
}

func main() {
	lambda.Start(Handler)
}
