package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// List of environment variables:
// TABLE_NAME
// EMPTY_DB_ALARM_NAME
// QUEUE_URL

func Handler(ctx context.Context, request events.SNSEvent) (interface{}, error) {
	tableName := os.Getenv("TABLE_NAME")
	emptyDbAlarmName := os.Getenv("EMPTY_DB_ALARM_NAME")
	queueUrl := os.Getenv("QUEUE_URL")
	// get item from dynamodb
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	dbClient := dynamodb.NewFromConfig(cfg)
	//scan dynamodb table for top 10 items
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int32(10),
	}
	result, err := dbClient.Scan(context.Background(), scanInput)
	if err != nil {
		panic(err)
	}
	if len(result.Items) == 0 {
		lib.SetAlarmState(emptyDbAlarmName, "ALARM")
		return request, nil
	}

	sqsClient := sqs.NewFromConfig(cfg)

	// send items to queue
	for _, item := range result.Items {
		j, e := json.Marshal(item)
		if e != nil {
			panic(err)
		}
		// send item to queue
		sendMessageInput := &sqs.SendMessageInput{
			MessageBody: aws.String(string(j)),
			QueueUrl:    aws.String(queueUrl),
		}
		_, err := sqsClient.SendMessage(context.Background(), sendMessageInput)
		if err != nil {
			panic(err)
		}
		// delete item from dynamodb
		deleteItemInput := &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{
					Value: fmt.Sprintf("%v", item["Id"]),
				},
			},
		}
		_, err = dbClient.DeleteItem(context.Background(), deleteItemInput)
		if err != nil {
			panic(err)
		}

	}

	return request, nil
}

func main() {
	lambda.Start(Handler)
}
