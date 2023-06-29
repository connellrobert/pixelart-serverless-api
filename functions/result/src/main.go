package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// List of environment variables:
// ANALYTICS_TABLE_NAME

func Handler(ctx context.Context, sqsResult events.SQSEvent) {

	for _, message := range sqsResult.Records {
		var result lib.ResultRequest
		err := json.Unmarshal([]byte(message.Body), &result)
		if err != nil {
			panic(err)
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"),
		)
		if err != nil {
			panic(err)
		}

		tableName := os.Getenv("ANALYTICS_TABLE_NAME")

		getItemInput := &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{
					Value: result.Record.Id,
				},
			},
		}

		client := dynamodb.NewFromConfig(cfg)
		record, err := client.GetItem(context.Background(), getItemInput)
		if err != nil {
			panic(err)
		}
		analyticsItem := lib.AnalyticsItem{}
		analyticsItem.FromDynamoDB(record.Item)
		attemptNum := len(analyticsItem.Attempts)
		analyticsItem.Attempts[strconv.Itoa(attemptNum)] = result.Result
		analyticsItem.Success = result.Result.Success

		updateItemInput := &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{
					Value: result.Record.Id,
				},
			},
			ExpressionAttributeNames: map[string]string{
				"#A": "Record",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":a": &types.AttributeValueMemberM{Value: analyticsItem.ToDynamoDB()},
			},
			UpdateExpression: aws.String("SET #A = :a"),
		}

		_, err = client.UpdateItem(context.Background(), updateItemInput)
		if err != nil {
			panic(err)
		}

		if attemptNum < 3 || !result.Result.Success {
			lib.SendRetrySignal(result.Record)
		}
		lib.SubmitXRayTraceSubSegment(result.Record.Metadata.TraceId, "Updated analytics item")
	}
}

func main() {
	lambda.Start(Handler)
}
