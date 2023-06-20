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

// lambda handler
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) {

	/*
		{
			body: {
				action: 0,
				params: {
					<type of image create, edit, or variate>
				}
			}
		}
	*/

	id := uuid.New().String()
	body := make(map[string]interface{})
	err := json.Unmarshal([]byte(request.Body), &body)
	action, err := strconv.Atoi(fmt.Sprintf("%v", body["action"]))
	if err != nil {
		panic(err)
	}
	requestAction := lib.RequestAction(action)
	record := lib.QueueRequest{
		Id:       id,
		Action:   requestAction,
		Priority: 0,
	}
	record.MapParams(requestAction, body["params"])
	tableName := os.Getenv(lib.ActionToTableEnvMapping[requestAction])
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
			"PK": &types.AttributeValueMemberS{
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
			"PK": &types.AttributeValueMemberS{
				Value: record.Id,
			},
			"Priority": &types.AttributeValueMemberN{
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

	fmt.Println(request)
}

func main() {
	lambda.Start(Handler)
}
