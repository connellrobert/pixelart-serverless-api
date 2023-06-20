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

	id := uuid.New().String()
	body := make(map[string]interface{})
	err := json.Unmarshal([]byte(request.Body), &body)
	action, err := strconv.Atoi(fmt.Sprintf("%v", body["action"]))
	if err != nil {
		panic(err)
	}
	requestAction := lib.RequestAction(action)
	paramName := lib.mapping[requestAction]["params"]
	param := fmt.Sprintf("%v", body[paramName])
	record := lib.QueueRequest{
		Id:       id,
		Action:   requestAction,
		Priority: 0,
	}
	record[paramName] = param
	tableName := os.Getenv(lib.mapping[requestAction]["TableEnv"])
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)
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
				Value: record,
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
