package main

import (
	"context"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	lambda.Start(Handler)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get id from request path
	id := req.PathParameters["id"]
	tableName := os.Getenv("ANALYTICS_TABLE_NAME")
	// check analytics db for id
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	}
	record, err := client.GetItem(context.Background(), getItemInput)
	if err != nil {
		panic(err)
	}
	analyticsItem := lib.AnalyticsItem{}
	analyticsItem.FromDynamoDB(record.Item)
	// return analytics item
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{\"url\": \"" + analyticsItem.Attempts["0"].Response.Data[0].URL + "\"}",
	}, nil

}
