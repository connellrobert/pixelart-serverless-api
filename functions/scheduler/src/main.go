package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
	fmt.Print(action)
	if err != nil {
		panic(err)
	}
	requestAction := lib.RequestAction(action)
	fmt.Print(requestAction)
	record := lib.QueueRequest{
		Metadata: lib.CommonMetadata{
			TraceId: traceId,
		},
		Id:       id,
		Action:   requestAction,
		Priority: 0,
	}
	record.MapParams(requestAction, body["params"])
	tableName := os.Getenv(lib.ActionToTableEnvMapping[requestAction])
	fmt.Println(lib.ActionToTableEnvMapping[requestAction])
	fmt.Println(tableName)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	// Create Analytics Item
	analyticsItem := lib.AnalyticsItem{
		Id:       record.Id,
		Record:   record,
		Attempts: make(map[string]lib.ImageResponseWrapper),
	}
	// Store Analytics Item
	analyticsTable := os.Getenv("ANALYTICS_TABLE_NAME")
	putAItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(analyticsTable),
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: record.Id,
			},
			"Record": &types.AttributeValueMemberM{
				Value: analyticsItem.ToDynamoDB(),
			},
		},
	}
	_, err = client.PutItem(context.Background(), putAItemInput)
	if err != nil {
		panic(err)
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: record.Id,
			},
			"priority": &types.AttributeValueMemberN{
				Value: strconv.Itoa(record.Priority),
			},
			"request": &types.AttributeValueMemberM{
				Value: record.ToDynamoDB(),
			},
		},
	}
	_, err = client.PutItem(context.Background(), putItemInput)
	if err != nil {
		panic(err)
	}
	alarmName := lib.ActionToAlarmMapping[requestAction]
	lib.SetAlarmState(alarmName, "OK")
	lib.SubmitXRayTraceSubSegment(traceId, "Added item to "+tableName)
	responseBody := map[string]interface{}{
		"message": fmt.Sprintf("Successfully added %v to %v", record.Id, tableName),
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
