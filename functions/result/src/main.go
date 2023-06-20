package main

import (
	"context"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Handler(ctx context.Context, result lib.ResultRequest) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}

	tableName := os.Getenv("TABLE_NAME")

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: result.Record.Id,
			},
		},
	}

	client := dynamodb.NewFromConfig(cfg)
	record, err := client.GetItem(context.Background(), getItemInput)
	if err != nil {
		panic(err)
	}
	attemptNum := len(record.Item["Attempts"]) + 1
	record.Item["Attempts"] = append(record.Item["Attempts"], result.Result)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: result.Record.Id,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#A": "Attempts",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberL{
				Value: record.Item["Attempts"],
			},
		},
	}

	_, err = client.UpdateItem(context.Background(), updateItemInput)
	if err != nil {
		panic(err)
	}

	if attemptNum < 3 {
		lib.SendRetrySignal(result.Record)
	}

}

func main() {
	lambda.Start(Handler)
}
