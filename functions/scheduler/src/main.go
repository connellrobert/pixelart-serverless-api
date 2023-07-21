package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aimless-it/ai-canvas/functions/lib/process"
	aiTypes "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	traceId := request.Headers["x-amzn-trace-id"]
	id := uuid.New().String()
	body := make(map[string]interface{})
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		panic(err)
	}
	action, err := strconv.Atoi(fmt.Sprintf("%v", body["action"]))
	if err != nil {
		panic(err)
	}
	requestAction := aiTypes.RequestAction(action)
	fmt.Print(requestAction)
	record := aiTypes.QueueRequest{
		Metadata: aiTypes.CommonMetadata{
			TraceId: traceId,
		},
		Id:       id,
		Action:   requestAction,
		Priority: 0,
	}
	record.MapParams(requestAction, body["params"])
	cfg := process.GetAWSConfig()
	client := dynamodb.NewFromConfig(cfg)
	// Create Analytics Item
	analyticsItem := aiTypes.AnalyticsItem{
		Id:       record.Id,
		Record:   record,
		Attempts: make(map[string]aiTypes.ImageResponseWrapper),
	}
	// Store Analytics Item
	analyticsTable := os.Getenv("ANALYTICS_TABLE_NAME")
	putAItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(analyticsTable),
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: record.Id,
			},
			"record": &types.AttributeValueMemberM{
				Value: analyticsItem.ToDynamoDB(),
			},
		},
	}
	_, err = client.PutItem(context.Background(), putAItemInput)
	if err != nil {
		panic(err)
	}
	process.SendRequestToQueue(record)
	process.SubmitXRayTraceSubSegment(traceId, "Added item "+record.Id+" to queue")
	responseBody := map[string]interface{}{
		"message": fmt.Sprintf("Successfully added %v", record.Id),
		"id":      record.Id,
	}
	re, err := json.Marshal(responseBody)
	if err != nil {
		panic(err)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(re),
	}
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
