package main

import (
	"context"
	"fmt"
	"os"

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
func Handler(ctx context.Context, request events.ApiGatewayProxyRequest) {
	mapping := map[string]map[string]string{
		"GENERATE_IMAGE": {
			"params":   "generateImage",
			"TableEnv": "GI_TABLE_NAME",
		},
		"EDIT_IMAGE": {
			"params":   "createImageEdit",
			"TableEnv": "EI_TABLE_NAME",
		},
		"VARIATE_IMAGE": {
			"params":   "createImageVariation",
			"TableEnv": "VI_TABLE_NAME",
		},
	}
	id := uuid.New().String()
	body := request.Body
	action := body["action"]
	paramName := mapping[action]["params"]
	param := body[paramName]
	record := lib.QueueRequest{
		Id:       id,
		Action:   action,
		Priority: 0,
	}
	record[paramName] = param
	tableName := os.Getenv(mapping[action]["TableEnv"])
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
				id: record.Id,
			},
			"Priority": &types.AttributeValueMemberN{
				id: record.Priority,
			},
			"request": &types.AttributeValueMemberM{
				id: record,
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
